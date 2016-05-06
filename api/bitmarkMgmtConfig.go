// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/logger"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
)

type bMgmtPasswordRequset struct {
	Origin string
	New    string
}

// POST /api/password
func SetBitmarkMgmtPassword(w http.ResponseWriter, req *http.Request, bitmarkMgmtConfigFile string, password string, log *logger.L) {

	log.Info("POST /api/password")
	response := &Response{
		Ok:     false,
		Result: setPasswordErr,
	}

	decoder := json.NewDecoder(req.Body)
	var request bMgmtPasswordRequset
	err := decoder.Decode(&request)
	if nil != err {
		log.Errorf("Error:%v", err)
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	if password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(request.Origin)); nil != err {
			log.Errorf("Error: %v", fault.ErrWrongPassword)
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

	// write new password to bitmark-mgmt config file
	input, err := ioutil.ReadFile(bitmarkMgmtConfigFile)
	if nil != err {
		log.Errorf("Error: %v", err)
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
	err = ioutil.WriteFile(bitmarkMgmtConfigFile, []byte(output), 0644)
	if nil != err {
		log.Errorf("Error: %v", err)
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	response.Ok = true
	response.Result = nil
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}

	// clean request password
	request.New = "0000000000000000"
}
