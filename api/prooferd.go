// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	// bitmarkdConfig "github.com/bitmark-inc/bitmarkd/configuration"
	"github.com/bitmark-inc/logger"
	"net/http"
	"time"
)

type prooferdRequest struct {
	Option     string `json:"option"`
	ConfigFile string `json:"config_file"`
}

// POST /api/bitmarkd
func Prooferd(w http.ResponseWriter, req *http.Request, webguiFilePath string, webguiConfig *configuration.Configuration, log *logger.L) {

	log.Info("POST /api/prooferd")
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	decoder := json.NewDecoder(req.Body)
	var request bitmarkdRequest
	if err := decoder.Decode(&request); nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	log.Infof("prooferd option: %s", request.Option)

	apiErr := invalidValueErr
	switch request.Option {
	case `start`:
		// Check if prooferd is running
		if prooferdService.IsRunning() {
			response.Result = prooferdAlreadyStartErr
		} else {
			prooferdService.ModeStart <- true
			// wait one second to get correct result
			time.Sleep(time.Second * 1)
			if !prooferdService.IsRunning() {
				response.Result = prooferdStartErr
			} else {
				response.Ok = true
				response.Result = prooferdStartSuccess
			}
		}
	case `stop`:
		if !prooferdService.IsRunning() {
			response.Result = prooferdAlreadyStopErr
		} else {
			prooferdService.ModeStart <- false
			time.Sleep(time.Second * 1)
			if prooferdService.IsRunning() {
				response.Result = prooferdStopErr
			} else {
				response.Ok = true
				response.Result = prooferdStopSuccess
			}

		}
	case `status`:
		response.Ok = true
		if prooferdService.IsRunning() {
			response.Result = prooferdStarted
		} else {
			response.Result = prooferdStopped
		}
	case `setup`:
		// if prooferdService.IsRunning() {
		// 	response.Result = prooferdAlreadyStartErr
		// } else {
		// 	if err := prooferdService.Setup(request.ConfigFile, webguiFilePath, webguiConfig); nil != err {
		// 		response.Result = err
		// 	} else {
		// 		response.Ok = true
		// 		response.Result = nil
		// 	}
		// }
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
