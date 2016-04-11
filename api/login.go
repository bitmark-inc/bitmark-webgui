// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/logger"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// GET /api/login
func LoginStatus(w http.ResponseWriter, log *logger.L) {

	log.Info("GET /api/login: check login status")
	response := &Response{
		Ok:     true,
		Result: nil,
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

type loginRequset struct {
	Password string
}

// POST /api/login
func LoginBitmarkMgmt(w http.ResponseWriter, req *http.Request, password string, log *logger.L) {

	log.Info("POST /api/login")
	response := &Response{
		Ok:     false,
		Result: fault.ApiErrLogin,
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

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(request.Password)); nil != err {
		log.Errorf("Login failed: %v, Host: %v, User-Agent: %v", fault.ErrWrongPassword, req.Host, req.Header.Get("User-Agent"))
		if err := writeApiResponse(w, response); nil != err {
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
	request.Password = "0000000000000000"

}

// POST /api/logout
func LogoutBitmarkMgmt(w http.ResponseWriter, log *logger.L) {

	log.Info("POST /api/logout")
	cookie := &http.Cookie{
		Name:   cookieName,
		Secure: true,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	response := &Response{
		Ok:     true,
		Result: nil,
	}

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if b, err := json.MarshalIndent(response, "", "  "); nil != err {
		log.Errorf("Error: %v", fault.ErrJsonParseFail)
	} else {
		w.Write(b)
	}
}
