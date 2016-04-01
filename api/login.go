// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// GET /api/login
func LoginStatus(w http.ResponseWriter) {

	fmt.Println("GET /api/login")
	response := &Response{
		Ok:     true,
		Result: nil,
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		fmt.Printf("Error: %v\n", err)
	}
}

type loginRequset struct {
	Password string
}

// POST /api/login
func LoginBitmarkMgmt(w http.ResponseWriter, req *http.Request, password string) {

	fmt.Println("POST /api/login")
	response := &Response{
		Ok:     false,
		Result: fault.ApiErrLogin,
	}

	decoder := json.NewDecoder(req.Body)
	var request loginRequset
	err := decoder.Decode(&request)
	if nil != err {
		fmt.Printf("Error:%v\n", err)
		if err := writeApiResponse(w, response); nil != err {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(request.Password)); nil != err {
		fmt.Printf("Error: %v\n", fault.ErrWrongPassword)
		if err := writeApiResponse(w, response); nil != err {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	response.Ok = true
	response.Result = nil
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		fmt.Printf("Error: %v\n", err)
	}

	// clean request password
	request.Password = "0000000000000000"

}

// POST /api/logout
func LogoutBitmarkMgmt(w http.ResponseWriter) {

	fmt.Println("POST /api/logout")
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
		fmt.Printf("Error: %v\n", fault.ErrJsonParseFail)
	} else {
		w.Write(b)
	}
}
