// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"github.com/bitmark-inc/exitwithstatus"
	"github.com/bitmark-inc/bitmark-mgmt/api"
	"github.com/bitmark-inc/bitmark-mgmt/configuration"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/bitmark-mgmt/templates"
	"github.com/bitmark-inc/bitmark-mgmt/utils"
	"github.com/codegangsta/cli"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

var GlobalConfig *configuration.Configuration
var BitmarkMgmtConfigFile string

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
	fmt.Println("configDir", configDir)

	if nil != err {
		exitwithstatus.Message("Error: %s\n", err)
	}

	if !utils.EnsureFileExists(configDir) {
		if err := os.MkdirAll(configDir, 0755); nil != err {
			exitwithstatus.Message("Error: %v\n", err)
		}
	}

	configFile, err := configuration.GetConfigPath(configDir)
	if nil != err {
		exitwithstatus.Message("Error: %v\n", err)
	}

	// Check if file exist
	if !utils.EnsureFileExists(configFile) {
		file, err := os.Create(configFile)
		if nil != err {
			exitwithstatus.Message("Error: %v\n", err)
		}

		encryptPassword, err := bcrypt.GenerateFromPassword([]byte("bitmark-mgmt"), bcrypt.DefaultCost)
		if nil != err {
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
			exitwithstatus.Message("Error: %v\n", err)
		}
	} else {
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

	// read bitmark-mgmt config file
	if configs, err := configuration.GetConfiguration(configFile); nil != err {
		exitwithstatus.Message("Error: %v\n", err)
	} else {
		GlobalConfig = configs
		if err := startWebServer(GlobalConfig); err != nil {
			exitwithstatus.Message("Error: %v\n", err)
		}

	}
}

func startWebServer(configs *configuration.Configuration) error {
	fmt.Println("bitmark-mgmt web start...")

	host := "127.0.0.1"
	port := strconv.Itoa(configs.Port)

	// set up webpages config.js
	if err := setupWebpagesConfig(host, port, configs.EnableHttps); nil != err {
		return err
	}

	// serve web pages
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("./webpages/lib/"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./webpages/scripts/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./webpages/images/"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./webpages/styles/"))))
	http.Handle("/", http.FileServer(http.Dir("./webpages/")))

	// serve api
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
		// TODO: gen certs
		cert, key, newCreate, err := utils.GetTLSCertFile()
		if nil != err {
			return err
		}

		if newCreate {
			if err := utils.MakeSelfSignedCertificate("bitmark-mgmt", cert, key, false, nil); nil != err {
				return err
			}
		}

		if err := server.ListenAndServeTLS(cert, key); nil != err {
			log.Fatal(err)
			return err
		}
	} else {
		if err := server.ListenAndServe(); nil != err {
			log.Fatal(err)
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

func setupWebpagesConfig(host string, port string, enableHttps bool) error {
	configFile := "./webpages/scripts/config.js"
	if !utils.EnsureFileExists(configFile) {
		if _, err := os.Create(configFile); nil != err {
			return err
		}
	}

	webPagesConfig := &webPagesConfigType{
		Host:        host,
		Port:        port,
		EnableHttps: enableHttps,
	}
	configTemp := template.Must(template.New("config").Parse(templates.WebpagesConfigTemplate))
	configBuffer := new(bytes.Buffer)
	if err := configTemp.Execute(configBuffer, webPagesConfig); nil != err {
		return err
	}
	if err := ioutil.WriteFile(configFile, []byte(configBuffer.String()), 0644); nil != err {
		return err
	}
	return nil
}

func checkAuthorization(w http.ResponseWriter, req *http.Request, writeHeader bool) bool {
	if GlobalConfig.EnableHttps {
		if err := api.GetAndCheckCookie(w, req); nil != err {
			fmt.Printf("Error: %v\n", err)
			if writeHeader {
				w.WriteHeader(http.StatusUnauthorized)
			}
			return false
		}
	}

	return true
}

func handleConfig(w http.ResponseWriter, req *http.Request) {

	api.SetCORSHeader(w, req)

	switch req.Method {
	case `GET`: // list bitmark config
		if !checkAuthorization(w, req, true) {
			return
		}
		api.ListConfig(w, req, GlobalConfig.BitmarkConfigFile)
	case `POST`:
		if !checkAuthorization(w, req, true) {
			return
		}
		api.UpdateConfig(w, req, GlobalConfig.BitmarkConfigFile)
	case `OPTIONS`:
		return
	default:
		fmt.Println("Error: Unknow method")
	}
}

func handleSetPassword(w http.ResponseWriter, req *http.Request) {
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true) {
		return
	}

	switch req.Method {
	case `POST`:
		if !utils.EnsureFileExists(BitmarkMgmtConfigFile) {
			exitwithstatus.Message("Error: %s\n", fault.ErrNotFoundConfigFile)
		}
		if configs, err := configuration.GetConfiguration(BitmarkMgmtConfigFile); nil != err {
			exitwithstatus.Message("Error: %v\n", err)
		} else {
			GlobalConfig = configs
			api.SetBitmarkMgmtPassword(w, req, BitmarkMgmtConfigFile, GlobalConfig.Password)
		}
	case `OPTIONS`:
		return
	default:
		fmt.Println("Error: Unknow method")
	}
}

func handleLogin(w http.ResponseWriter, req *http.Request) {
	api.SetCORSHeader(w, req)

	switch req.Method {
	case `GET`:
		if !checkAuthorization(w, req, true) {
			return
		}
		api.LoginStatus(w)
	case `POST`:
		if GlobalConfig.EnableHttps && checkAuthorization(w, req, false) {
			if err := api.WriteGlobalErrorResponse(w, fault.ApiErrAlreadyLoggedIn); nil != err {
				fmt.Printf("Error: %v\n", err)
			}
			return
		}
		api.LoginBitmarkMgmt(w, req, GlobalConfig.Password)
	case `OPTIONS`:
		return
	default:
		fmt.Println("Error: Unknow method")
	}
}

func handleLogout(w http.ResponseWriter, req *http.Request) {
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true) {
		return
	}

	switch req.Method {
	case `POST`:
		api.LogoutBitmarkMgmt(w)
	case `OPTIONS`:
		return
	default:
		fmt.Println("Error: Unknow method")
	}
}

func handleBitmarkd(w http.ResponseWriter, req *http.Request) {
	api.SetCORSHeader(w, req)

	if req.Method == "OPTIONS" || !checkAuthorization(w, req, true) {
		return
	}

	switch req.Method {
	case `POST`:
		api.Bitmarkd(w, req, GlobalConfig.BitmarkConfigFile)
	case `OPTIONS`:
		return
	default:
		fmt.Println("Error: Unknow method")
	}
}
