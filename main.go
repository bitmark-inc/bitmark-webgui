// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/bitmark-inc/bitmark-webgui/api"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/templates"
	"github.com/bitmark-inc/bitmark-webgui/utils"
	"github.com/bitmark-inc/exitwithstatus"
	"github.com/bitmark-inc/logger"
	"github.com/codegangsta/cli"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"
)

var GlobalConfig *configuration.Configuration
var BitmarkWebguiConfigFile string

var mainLog *logger.L

func main() {
	// ensure exit handler is first
	defer exitwithstatus.Handler()

	var configFile string

	app := cli.NewApp()
	app.Name = "bitmark-webgui"
	app.Usage = "Configuration program for bitmarkd"
	app.Version = Version()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config-file, c",
			Value:       "",
			Usage:       "*bitmark-webgui config file",
			Destination: &configFile,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "setup",
			Usage: "Initialise bitmark-webgui configuration",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hostname, H",
					Value: "",
					Usage: "generate server certificate with the hostname [localhost]",
				},
				cli.StringFlag{
					Name:  "data-directory, d",
					Value: "",
					Usage: "the direcotry of web and log",
				},
			},
			Action: func(c *cli.Context) error {
				runSetup(c, configFile)
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "start bitmark-webgui",
			Action: func(c *cli.Context) error {
				runStart(c, configFile)
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func runSetup(c *cli.Context, configFile string) {

	// set data-directory
	dataDir := c.String("data-directory")
	defaultConfig, err := configuration.GetDefaultConfiguration(dataDir)
	if nil != err {
		exitwithstatus.Message("Error: %v\n", err)
	}

	// set logger
	setupLogger(&defaultConfig.Logging)
	defer logger.Finalise()

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

		encryptPassword, err := bcrypt.GenerateFromPassword([]byte(defaultConfig.Password), bcrypt.DefaultCost)
		if nil != err {
			mainLog.Errorf("Encrypt password failed: %v", err)
			exitwithstatus.Message("Error: %v\n", err)
		}

		defaultConfig.Password = string(encryptPassword)

		// generate config file
		confTemp := template.Must(template.New("config").Parse(templates.ConfigurationTemplate))
		if err := confTemp.Execute(file, defaultConfig); nil != err {
			mainLog.Errorf("Generate config template failed: %v", err)
			exitwithstatus.Message("Error: %v\n", err)
		}
		mainLog.Info("Successfully setup bitmark-webgui configuration file")

		// gen certificate
		hostname := c.String("hostname")
		if "" != hostname {
			// gen certs
			cert, key, newCreate, err := utils.GetTLSCertFile(defaultConfig.DataDirectory)
			if nil != err {
				mainLog.Errorf("get TLS file failed: %v", err)
				exitwithstatus.Message("get TLS file failed: %v\n", err)
			}

			if newCreate {
				mainLog.Infof("Generate self signed certificate for hostname: %s", hostname)
				hostnames := []string{hostname}
				if err := utils.MakeSelfSignedCertificate("bitmark-webgui", cert, key, false, hostnames); nil != err {
					mainLog.Errorf("generate TLS file failed: %v", err)
					exitwithstatus.Message("generate TLS file failed: %v\n", err)
				}
			} else {
				mainLog.Error("TLS file existed")
				exitwithstatus.Message("TLS file existed\n")
			}
			mainLog.Info("Successfully generate TLS files")
		}
	} else {
		mainLog.Errorf("config file %s existed", configFile)
		exitwithstatus.Message("Error: %s existed\n", configFile)
	}

}

func runStart(c *cli.Context, configFile string) {

	if !utils.EnsureFileExists(configFile) {
		exitwithstatus.Message("Error: %v\n", fault.ErrNotFoundConfigFile)
	}

	BitmarkWebguiConfigFile = configFile

	// read bitmark-webgui config file
	if configs, err := configuration.GetConfiguration(configFile); nil != err {
		exitwithstatus.Message("Error: %v\n", err)
	} else {
		GlobalConfig = configs

		setupLogger(&configs.Logging)
		defer logger.Finalise()

		// initialise services
		if err := InitialiseService(configs); nil != err {
			mainLog.Criticalf("initialise background services failed: %v", err)
			exitwithstatus.Exit(1)
		}
		defer FinaliseBackgroundService()

		go func() {
			if err := startWebServer(GlobalConfig); err != nil {
				mainLog.Criticalf("%s", err)
				exitwithstatus.Message("Error: %v\n", err)
			}
		}()

		// turn Signals into channel messages
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-ch
		mainLog.Infof("received signal: %v", sig)
		mainLog.Info("shutting down...")
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
	baseWebDir := configs.DataDirectory + "/webpages"
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir(baseWebDir+"/lib/"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir(baseWebDir+"/scripts/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(baseWebDir+"/images/"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir(baseWebDir+"/styles/"))))
	http.Handle("/", http.FileServer(http.Dir(baseWebDir+"/")))

	// serve api
	mainLog.Info("Set up server api")
	http.HandleFunc("/api/config", handleConfig)
	http.HandleFunc("/api/password", handleSetPassword)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/logout", handleLogout)
	http.HandleFunc("/api/bitmarkd", handleBitmarkd)

	http.HandleFunc("/api/bitmarkPay/encrypt", handleBitmarkPay)
	http.HandleFunc("/api/bitmarkPay/info", handleBitmarkPay)
	http.HandleFunc("/api/bitmarkPay/pay", handleBitmarkPay)

	http.HandleFunc("/api/bitmarkCli/info", handleBitmarkCli)
	http.HandleFunc("/api/bitmarkCli/setup", handleBitmarkCli)
	http.HandleFunc("/api/bitmarkCli/issue", handleBitmarkCli)
	http.HandleFunc("/api/bitmarkCli/transfer", handleBitmarkCli)

	http.HandleFunc("/api/onestep/status", handleOnestep)
	http.HandleFunc("/api/onestep/setup", handleOnestep)
	http.HandleFunc("/api/onestep/issue", handleOnestep)
	http.HandleFunc("/api/onestep/transfer", handleOnestep)

	server := &http.Server{
		Addr:           host + ":" + port,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if configs.EnableHttps {
		mainLog.Info("Starting https server...")
		// gen certs
		cert, key, newCreate, err := utils.GetTLSCertFile(configs.DataDirectory)
		if nil != err {
			return err
		}

		if newCreate {
			mainLog.Info("Generate self signed certificate...")
			if err := utils.MakeSelfSignedCertificate("bitmark-webgui", cert, key, false, nil); nil != err {
				return err
			}
		}

		if err := server.ListenAndServeTLS(cert, key); nil != err {
			return err
		}
	} else {
		// turn Signals into channel messages
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-ch
			mainLog.Infof("received signal: %v", sig)
		}()

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
			cookie := &http.Cookie{
				Name:   api.CookieName,
				Secure: true,
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)
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
	log := logger.New("api-bitmarkWebgui")
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true, log) {
		return
	}

	switch req.Method {
	case `POST`:
		if !utils.EnsureFileExists(BitmarkWebguiConfigFile) {
			exitwithstatus.Message("Error: %s\n", fault.ErrNotFoundConfigFile)
		}
		if configs, err := configuration.GetConfiguration(BitmarkWebguiConfigFile); nil != err {
			exitwithstatus.Message("Error: %v\n", err)
		} else {
			GlobalConfig = configs
			api.SetBitmarkWebguiPassword(w, req, BitmarkWebguiConfigFile, GlobalConfig.Password, log)
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
		api.LoginBitmarkWebgui(w, req, GlobalConfig.Password, log)
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
		api.LogoutBitmarkWebgui(w, log)
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

func handleBitmarkPay(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-bitmarkPay")
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true, log) {
		return
	}

	switch req.Method {
	case `POST`:
		reqUriArr := strings.Split(req.RequestURI, "/")
		api.BitmarkPayEncrypt(w, req, log, reqUriArr[3])
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}

func handleBitmarkCli(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-bitmarkCli")
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true, log) {
		return
	}

	switch req.Method {
	case `POST`:
		reqUriArr := strings.Split(req.RequestURI, "/")
		api.BitmarkCliExec(w, req, log, reqUriArr[3])
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}

func handleOnestep(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-onestep")
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true, log) {
		return
	}

	switch req.Method {
	case `POST`:
		reqUriArr := strings.Split(req.RequestURI, "/")
		api.OnestepExec(w, req, log, reqUriArr[3])
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}
