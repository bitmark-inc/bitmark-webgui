// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/logger"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sync"
	"time"
)

const (
	CookieName = "bitmark-webgui"
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

type cookie struct {
	sync.RWMutex
	expireTime time.Time
	cipher     string
}

var globalCookie [2]cookie // 0 is recent cookie, 1 is previous
var cookieExpireDuration = 10

func setCookie(w http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   CookieName,
		Secure: true,
	}

	globalCookie[0].Lock()
	defer globalCookie[0].Unlock()
	globalCookie[1].Lock()
	defer globalCookie[1].Unlock()

	// check if globalCookie need to be updated
	localTime := time.Now().Add(time.Duration(2) * time.Minute)
	if localTime.After(globalCookie[0].expireTime) {
		// cookie will be expired in 2 minutes
		globalCookie[1].expireTime = globalCookie[0].expireTime
		globalCookie[0].expireTime = time.Now().Add(time.Duration(cookieExpireDuration) * time.Minute)

		cookieCipher, err := bcrypt.GenerateFromPassword([]byte(globalCookie[0].expireTime.String()), bcrypt.DefaultCost)
		if nil != err {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			return fault.ApiErrSetAuthorize
		}
		globalCookie[0].cipher = string(cookieCipher)

	}

	cookie.Value = globalCookie[0].cipher
	cookie.MaxAge = cookieExpireDuration * 60 //seconds
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

func GetAndCheckCookie(w http.ResponseWriter, req *http.Request, log *logger.L) error {
	reqCookie, err := req.Cookie(CookieName)
	if nil != err {
		log.Errorf("request cookie error: %v", err)
		return fault.ApiErrUnauthorized
	}

	globalCookie[0].Lock()
	defer globalCookie[0].Unlock()
	globalCookie[1].Lock()
	defer globalCookie[1].Unlock()

	for _, c := range globalCookie {
		log.Infof("decrypt cookie: %s", c.cipher)
		log.Infof("request cookie: %s", reqCookie.Value)
		if err := bcrypt.CompareHashAndPassword([]byte(reqCookie.Value), []byte(c.expireTime.String())); nil == err {
			log.Infof("pass cookie")
			return nil
		}
	}

	log.Errorf("decrypt cookie error: %v", err)
	return fault.ApiErrChekAuthorize
}

func SetCORSHeader(w http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
