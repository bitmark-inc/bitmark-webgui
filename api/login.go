// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/services"
	"github.com/bitmark-inc/logger"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type LoginResponse struct {
	Chain                string `json:"chain"`
	BitmarkCliConfigFile string `json:"bitmark_cli_config_file"`
}

// GET /api/login
func LoginStatus(w http.ResponseWriter, configuration *configuration.Configuration, log *logger.L) {

	log.Info("GET /api/login: check login status")
	loginResponse := &LoginResponse{
		Chain:                configuration.BitmarkChain,
		BitmarkCliConfigFile: configuration.BitmarkCliConfigFile,
	}
	response := &Response{
		Ok:     true,
		Result: loginResponse,
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

type loginRequset struct {
	Password string
}

// POST /api/login
func LoginBitmarkWebgui(w http.ResponseWriter, req *http.Request, configuration *configuration.Configuration, log *logger.L) {

	log.Info("POST /api/login")
	response := &Response{
		Ok:     false,
		Result: loginErr,
	}

	decoder := json.NewDecoder(req.Body)
	var request loginRequset
	err := decoder.Decode(&request)
	if nil != err {
		log.Errorf("Error: %v", err)
		if err := writeApiResponse(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(configuration.Password), []byte(request.Password)); nil != err {
		log.Errorf("Login failed: %v, Host: %v, User-Agent: %v", fault.ErrWrongPassword, req.Host, req.Header.Get("User-Agent"))
		if err := writeApiResponse(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	loginResponse := &LoginResponse{
		Chain:                configuration.BitmarkChain,
		BitmarkCliConfigFile: configuration.BitmarkCliConfigFile,
	}

	response.Ok = true
	response.Result = loginResponse
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}

	// clean request password
	request.Password = "0000000000000000"

}

type logoutRequset struct {
	Password             string `json:"password"`
	BitmarkPayConfigFile string `json:"bitmark_pay_config_file"`
}

// POST /api/logout
func LogoutBitmarkWebgui(w http.ResponseWriter, req *http.Request, filePath string, webguiConfiguration *configuration.Configuration, log *logger.L) {
	response := &Response{
		Ok:     false,
		Result: "logout error",
	}

	cookie := &http.Cookie{
		Name:   CookieName,
		Secure: true,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	response.Ok = true
	response.Result = nil

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if b, err := json.MarshalIndent(response, "", "  "); nil != err {
		log.Errorf("Error: %v", fault.ErrJsonParseFail)
	} else {
		w.Write(b)
	}
}

// POST /api/logoutOnestep
func LogoutBitmarkWebguiOnestep(w http.ResponseWriter, req *http.Request, filePath string, webguiConfiguration *configuration.Configuration, log *logger.L) {

	log.Info("POST /api/logout")
	response := &Response{
		Ok:     false,
		Result: "logout error",
	}

	decoder := json.NewDecoder(req.Body)
	var request logoutRequset
	err := decoder.Decode(&request)
	if nil != err {
		log.Errorf("Error: %v", err)
		if err := writeApiResponse(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// get privateKey from bitmark-cli
	var keyPair BitmarkCliGenerateResponse
	bitmarkCliKeyPair := services.BitmarkCliKeyPairType{
		Password: request.Password,
	}
	output, err := bitmarkCliService.KeyPair(bitmarkCliKeyPair, webguiConfiguration.BitmarkCliConfigFile)
	if nil != err {
		log.Errorf("Error: %v", err)

		response.Result = "password error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	} else {
		if err := json.Unmarshal(output, &keyPair); nil != err {
			log.Errorf("Error: %v", err)

			response.Result = "parse json error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}
	}

	// decrypt bitmark-wallet
	if _, err := os.Stat(request.BitmarkPayConfigFile); nil != err {
		response.Result = "bitmarkPayConfigFile not existed"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	net := webguiConfiguration.BitmarkChain
	if "local" == net {
		net = "local_bitcoin_reg"
	}
	decryptConfig := services.BitmarkPayType{
		Net:      net,
		Config:   request.BitmarkPayConfigFile,
		Password: keyPair.PrivateKey,
	}
	if err := bitmarkPayService.Decrypt(decryptConfig); nil != err {
		log.Errorf("decrypt bitmarkPay error: %v\n", err)
		response.Result = "internal error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	ticker := time.NewTicker(time.Millisecond * 500)
loop:
	for range ticker.C {
		status, err := bitmarkPayService.Status(bitmarkPayService.GetBitmarkPayJobHash())
		if nil != err {
			log.Errorf("decrypt bitmarkPay error: %v\n", err)
			response.Result = "internal error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			break loop
		}

		switch status {
		case "success":
			ticker.Stop()
			break loop
		case "fail":
			ticker.Stop()
			log.Errorf("decrypt bitmarkPay job error")
			response.Result = "internal error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			break loop
		case "stopped":
			ticker.Stop()
			break loop
		}
	}
	removeBitmarkCliConfigAndCookie(w, filePath, webguiConfiguration, log)
}

func removeBitmarkCliConfigAndCookie(w http.ResponseWriter, filePath string, webguiConfiguration *configuration.Configuration, log *logger.L) {
	response := &Response{
		Ok:     false,
		Result: "logout error",
	}

	// remove bitmark-cli config file
	log.Infof("removing file: %v\n", webguiConfiguration.BitmarkCliConfigFile)
	if err := os.Remove(webguiConfiguration.BitmarkCliConfigFile); nil != err {
		response.Result = "delete bitmark-cli config file error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	webguiConfiguration.BitmarkCliConfigFile = ""
	if err := configuration.UpdateConfiguration(filePath, webguiConfiguration); nil != err {
		response.Result = "update webgui config file error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// remove cookie
	cookie := &http.Cookie{
		Name:   CookieName,
		Secure: true,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	response.Ok = true
	response.Result = nil

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if b, err := json.MarshalIndent(response, "", "  "); nil != err {
		log.Errorf("Error: %v", fault.ErrJsonParseFail)
	} else {
		w.Write(b)
	}
}
