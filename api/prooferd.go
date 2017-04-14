// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/logger"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type prooferdRequest struct {
	Option  string `json:"option"`
	Network string `json:"network"`
}

// POST /api/bitmarkd
func Prooferd(w http.ResponseWriter, req *http.Request, webguiFilePath string, webguiConfig *configuration.Configuration, log *logger.L) {

	log.Info("POST /api/prooferd")
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	decoder := json.NewDecoder(req.Body)
	var request prooferdRequest
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
		if prooferdService.IsRunning() {
			response.Result = prooferdAlreadyStartErr
		} else {
			prooferdConfigFile := filepath.Join(webguiConfig.DataDirectory, fmt.Sprintf("prooferd-%s", request.Network), "prooferd.conf")
			if err := prooferdService.Setup(prooferdConfigFile, webguiFilePath, webguiConfig); nil != err {
				if os.IsNotExist(err) {
					response.Result = "prooferd config not found"
				} else {
					response.Result = err.Error()
				}
			} else {
				response.Ok = true
				response.Result = nil
			}
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
