// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/logger"
	"net/http"
)



type bitmarkdRequest struct {
	Option string
}

// POST /api/bitmarkd
func Bitmarkd(w http.ResponseWriter, req *http.Request, bitmarkConfigFile string, log *logger.L) {

	log.Info("POST /api/bitmarkd")
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	decoder := json.NewDecoder(req.Body)
	var request bitmarkdRequest
	err := decoder.Decode(&request)
	if nil != err {
		log.Errorf("Error: %v", err)
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}
	log.Infof("bitmarkd option: %s", request.Option)


	apiErr := fault.ApiErrInvalidValue
	switch request.Option {
	case `start`:
		// Check if bitmarkd is running
		if bitmarkService.IsRunning() {
			response.Result = fault.ApiErrAlreadyStartBitmarkd
		} else {
			bitmarkService.ModeStart <- true
			if !bitmarkService.IsRunning() {
				response.Result = fault.ApiErrStartBitmarkd
			}else {
				response.Ok = true
				response.Result = "bitmarkd is running"
			}
		}
	case `stop`:
		if !bitmarkService.IsRunning() {
			response.Result = "bitmarkd is not running"
		} else {
			bitmarkService.ModeStart <- false
			if bitmarkService.IsRunning() {
				response.Result = fault.ApiErrStopBitmarkd
			}else {
				response.Ok = true
				response.Result = "stop running bitmarkd"
			}

		}
	case `status`:
		response.Ok = true
		if bitmarkService.IsRunning() {
			response.Result = "bitmarkd is running"
		} else {
			response.Result = "bitmarkd is not running"
		}
	default:
		response.Result = apiErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}

}
