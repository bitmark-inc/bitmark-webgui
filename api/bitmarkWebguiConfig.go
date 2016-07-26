// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/logger"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
)

type bWebguiPasswordRequset struct {
	Origin string
	New    string
}

// POST /api/password
func SetBitmarkWebguiPassword(w http.ResponseWriter, req *http.Request, bitmarkWebguiConfigFile string, configs *configuration.Configuration, log *logger.L) {

	log.Info("POST /api/password")
	response := &Response{
		Ok:     false,
		Result: setPasswordErr,
	}

	decoder := json.NewDecoder(req.Body)
	var request bWebguiPasswordRequset
	if err := decoder.Decode(&request); nil != err {
		log.Errorf("Error:%v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	if configs.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(configs.Password), []byte(request.Origin)); nil != err {
			log.Errorf("Error: %v", fault.ErrWrongPassword)
			response.Result = err
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}
	}

	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(request.New), bcrypt.DefaultCost)
	if nil != err {
		log.Errorf("Error: %v", err)
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// write new password to bitmark-webgui config file
	input, err := ioutil.ReadFile(bitmarkWebguiConfigFile)
	if nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Index(line, "password") == 0 {
			lines[i] = `password = "` + string(encryptPassword) + `"`
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(bitmarkWebguiConfigFile, []byte(output), 0644)
	if nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	configs.Password = string(encryptPassword)

	response.Ok = true
	response.Result = nil
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}

	// clear server cookie
	globalCookie[0] = cookie{}
	globalCookie[1] = cookie{}

	// clean request password
	request.New = "0000000000000000"

}
