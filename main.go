// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	// "fmt"
	"github.com/bitmark-inc/bitmark-mgmt/api"
	"github.com/bitmark-inc/bitmark-mgmt/configuration"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/bitmark-mgmt/templates"
	"github.com/bitmark-inc/bitmark-mgmt/utils"
	"github.com/bitmark-inc/exitwithstatus"
	"github.com/bitmark-inc/logger"
	"github.com/codegangsta/cli"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

var GlobalConfig *configuration.Configuration
var BitmarkMgmtConfigFile string
var ConfigDir string

var mainLog *logger.L

func main() {
	// ensure exit handler is first
	defer exitwithstatus.Handler()

	var configDir string

	app := cli.NewApp()
	app.Name = "bitmark-mgmt"
	app.Usage = "Configuration program for bitmarkd"
	app.Version = Version()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       "",
			Usage:       "*bitmark-mgmt config file",
			Destination: &configDir,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "setup",
			Usage: "Initialise bitmark-mgmt configuration",
			Action: func(c *cli.Context) {
				runSetup(c, configDir)
			},
		},
		{
			Name:  "start",
			Usage: "start bitmark-mgmt",
			Action: func(c *cli.Context) {
				runStart(c, configDir)
			},
		},
	}

	app.Run(os.Args)
}

func runSetup(c *cli.Context, configDir string) {

	configDir, err := utils.CheckConfigDir(configDir)
	if nil != err {
		exitwithstatus.Message("Error: %s\n", err)
	}

	if !utils.EnsureFileExists(configDir) {
		if err := os.MkdirAll(configDir, 0755); nil != err {
			exitwithstatus.Message("Error: %v\n", err)
		}
	}

	// set logger
	setupLogger(configuration.GetDefaultLogger(configDir))
	defer logger.Finalise()

	configFile, err := configuration.GetConfigPath(configDir)
	if nil != err {
		mainLog.Errorf("get config file path: %s error: %v", configFile, err)
		exitwithstatus.Message("Error: %v\n", err)
	}

	// Check if file exist
	if !utils.EnsureFileExists(configFile) {
		file, err := os.Create(configFile)
		if nil != err {
			mainLog.Errorf("create config file: %s failed: %v", configFile, err)
			exitwithstatus.Message("Error: %v\n", err)
		}

		encryptPassword, err := bcrypt.GenerateFromPassword([]byte("bitmark-mgmt"), bcrypt.DefaultCost)
		if nil != err {
			mainLog.Errorf("Encrypt password failed: %v", err)
			exitwithstatus.Message("Error: %v\n", err)
		}

		configData := configuration.Configuration{
			Port:              2150,
			Password:          string(encryptPassword),
			EnableHttps:       true,
			BitmarkConfigFile: "/etc/bitmarkd.conf",
		}

		confTemp := template.Must(template.New("config").Parse(templates.ConfigurationTemplate))
		if err := confTemp.Execute(file, configData); nil != err {
			mainLog.Errorf("Generate config template failed: %v", err)
			exitwithstatus.Message("Error: %v\n", err)
		}
		mainLog.Info("Successfully setup bitmark-mgmt configuration file")
	} else {
		mainLog.Errorf("config file %s existed", configFile)
		exitwithstatus.Message("Error: %s existed\n", configFile)
	}

}

func runStart(c *cli.Context, configDir string) {

	configDir, err := utils.CheckConfigDir(configDir)
	if nil != err {
		exitwithstatus.Message("Error: %s\n", err)
	}

	configFile, err := configuration.GetConfigPath(configDir)
	if nil != err {
		exitwithstatus.Message("Error: %v\n", err)
	}

	BitmarkMgmtConfigFile = configFile
	ConfigDir = configDir

	// read bitmark-mgmt config file
	if configs, err := configuration.GetConfiguration(configDir, configFile); nil != err {
		exitwithstatus.Message("Error: %v\n", err)
	} else {
		GlobalConfig = configs

		setupLogger(&configs.Logging)
		defer logger.Finalise()

		if err := startWebServer(GlobalConfig); err != nil {
			mainLog.Criticalf("%s", err)
			exitwithstatus.Message("Error: %v\n", err)
		}

	}
}

func setupLogger(logging *configuration.LoggerType) {
	// start logging
	if err := logger.Initialise(logging.File, logging.Size, logging.Count); nil != err {
		exitwithstatus.Message("%s: logger setup failed with error: %v", err)
	}

	logger.LoadLevels(logging.Levels)

	// create a logger channel for the main program
	mainLog = logger.New("main")
	mainLog.Info("starting…")
	mainLog.Debugf("loggerType: %v", logging)
}

func startWebServer(configs *configuration.Configuration) error {
	host := "0.0.0.0"
	port := strconv.Itoa(configs.Port)

	// serve web pages
	mainLog.Info("Set up server files")
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("./webpages/lib/"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./webpages/scripts/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./webpages/images/"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./webpages/styles/"))))
	http.Handle("/", http.FileServer(http.Dir("./webpages/")))

	// serve api
	mainLog.Info("Set up server api")
	http.HandleFunc("/api/config", handleConfig)
	http.HandleFunc("/api/password", handleSetPassword)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/logout", handleLogout)
	http.HandleFunc("/api/bitmarkd", handleBitmarkd)

	server := &http.Server{
		Addr:           host + ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if configs.EnableHttps {
		mainLog.Info("Starting https server...")
		// gen certs
		cert, key, newCreate, err := utils.GetTLSCertFile()
		if nil != err {
			return err
		}

		if newCreate {
			mainLog.Info("Generate self signed certificate...")
			if err := utils.MakeSelfSignedCertificate("bitmark-mgmt", cert, key, false, nil); nil != err {
				return err
			}
		}

		if err := server.ListenAndServeTLS(cert, key); nil != err {
			return err
		}
	} else {
		mainLog.Info("Starting http server...")
		if err := server.ListenAndServe(); nil != err {
			return err
		}
	}

	return nil
}

type webPagesConfigType struct {
	Host        string
	Port        string
	EnableHttps bool
}

func checkAuthorization(w http.ResponseWriter, req *http.Request, writeHeader bool, log *logger.L) bool {
	if GlobalConfig.EnableHttps {
		if err := api.GetAndCheckCookie(w, req, log); nil != err {
			log.Errorf("Error: %v", err)
			if writeHeader {
				w.WriteHeader(http.StatusUnauthorized)
			}
			return false
		}
	}

	return true
}

func handleConfig(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-config")

	api.SetCORSHeader(w, req)

	switch req.Method {
	case `GET`: // list bitmark config
		if !checkAuthorization(w, req, true, log) {
			return
		}
		api.ListConfig(w, req, GlobalConfig.BitmarkConfigFile, log)
	case `POST`:
		if !checkAuthorization(w, req, true, log) {
			return
		}
		api.UpdateConfig(w, req, GlobalConfig.BitmarkConfigFile, log)
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}

func handleSetPassword(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-bitmarkmgmt")
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true, log) {
		return
	}

	switch req.Method {
	case `POST`:
		if !utils.EnsureFileExists(BitmarkMgmtConfigFile) {
			exitwithstatus.Message("Error: %s\n", fault.ErrNotFoundConfigFile)
		}
		if configs, err := configuration.GetConfiguration(ConfigDir, BitmarkMgmtConfigFile); nil != err {
			exitwithstatus.Message("Error: %v\n", err)
		} else {
			GlobalConfig = configs
			api.SetBitmarkMgmtPassword(w, req, BitmarkMgmtConfigFile, GlobalConfig.Password, log)
		}
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-login")
	api.SetCORSHeader(w, req)

	switch req.Method {
	case `GET`:
		if !checkAuthorization(w, req, true, log) {
			return
		}
		api.LoginStatus(w, log)
	case `POST`:
		if GlobalConfig.EnableHttps && checkAuthorization(w, req, false, log) {
			if err := api.WriteGlobalErrorResponse(w, fault.ApiErrAlreadyLoggedIn, log); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}
		api.LoginBitmarkMgmt(w, req, GlobalConfig.Password, log)
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}

func handleLogout(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-logout")
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true, log) {
		return
	}

	switch req.Method {
	case `POST`:
		api.LogoutBitmarkMgmt(w, log)
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}

func handleBitmarkd(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-bitmarkd")
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true, log) {
		return
	}

	switch req.Method {
	case `POST`:
		api.Bitmarkd(w, req, GlobalConfig.BitmarkConfigFile, log)
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}
