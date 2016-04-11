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
	"time"
)

const (
	cookieName = "bitmark-mgmt"
)

type Response struct {
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

func writeApiResponse(w http.ResponseWriter, response *Response) error {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if b, err := json.MarshalIndent(response, "", "  "); nil != err {
		return fault.ErrJsonParseFail
	} else {
		w.Write(b)
	}

	return nil
}

func writeApiResponseAndSetCookie(w http.ResponseWriter, response *Response) error {
	// set cookie
	if err := setCookie(w); nil != err {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	return writeApiResponse(w, response)
}

func setCookie(w http.ResponseWriter) error {

	cookie := &http.Cookie{
		Name:   cookieName,
		Secure: true,
	}

	timeStr := time.Now().String()
	cookiePlain = cookieName + ":" + timeStr
	cookieCipher, err := bcrypt.GenerateFromPassword([]byte(cookiePlain), bcrypt.DefaultCost)
	if nil != err {
		//fmt.Printf("Error: %v\n", err)
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		return fault.ApiErrSetAuthorize
	}

	cookie.Value = string(cookieCipher)
	cookie.MaxAge = 10 * 60 //seconds
	http.SetCookie(w, cookie)
	return nil
}

func WriteGlobalErrorResponse(w http.ResponseWriter, err error, log *logger.L) error {
	response := &Response{
		Ok:     false,
		Result: err,
	}
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("writeApiResponseAndSetCookie error: %v", err)
		return err
	}

	return nil
}

var cookiePlain string

func GetAndCheckCookie(w http.ResponseWriter, req *http.Request, log *logger.L) error {
	reqCookie, err := req.Cookie(cookieName)
	if nil != err {
		log.Errorf("request cookie error: %v", err)
		return fault.ApiErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(reqCookie.Value), []byte(cookiePlain)); nil != err {
		log.Errorf("decrypt cookie error: %v", err)
		return fault.ApiErrChekAuthorize
	}

	return nil
}

func SetCORSHeader(w http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
