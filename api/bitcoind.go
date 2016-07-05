// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/logger"
	"net/http"
	"time"
)

type bitcoindRequest struct {
	Option string `json:"option"`
}

type bitcoindInfoResponse struct {
	Version         int     `json:"version"`
	Protocolversion int     `json:"protocolversion"`
	Walletversion   int     `json:"walletversion"`
	Balance         float32 `json:"balance"`
	Blocks          int     `json:"blocks"`
	Timeoffset      int     `json:"timeoffset"`
	Connections     int     `json:"connections"`
	Proxy           string  `json:"proxy"`
	Difficulty      float32 `json:"difficulty"`
	Testnet         bool    `json:"testnet"`
	Keypoololdest   int     `json:"keypoololdest"`
	Keypoolsize     int     `json:"keypoolsize"`
	Paytxfee        float32 `json:"paytxfee"`
	Relayfee        float32 `json:"relayfee"`
	Errors          string  `json:"errors"`
}

// POST /api/bitcoind
func Bitcoind(w http.ResponseWriter, req *http.Request, log *logger.L) {

	response := &Response{
		Ok:     false,
		Result: nil,
	}

	decoder := json.NewDecoder(req.Body)
	var request bitcoindRequest
	err := decoder.Decode(&request)
	if nil != err {
		log.Errorf("Error: %v", err)
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	log.Infof("POST /api/bitcoind/%s", request.Option)

	apiErr := invalidValueErr
	switch request.Option {
	case `start`:
		// Check if bitcoind is running
		if bitcoinService.IsRunning() {
			response.Result = bitcoindAlreadyStartErr
		} else {
			bitcoinService.ModeStart <- true
			// wait one second to get correct result
			time.Sleep(time.Second * 1)
			if !bitcoinService.IsRunning() {
				response.Result = bitcoindStartErr
			} else {
				response.Ok = true
				response.Result = bitcoindStartSuccess
			}
		}
	case `stop`:
		if !bitcoinService.IsRunning() {
			response.Result = bitcoindAlreadyStopErr
		} else {
			bitcoinService.ModeStart <- false
			time.Sleep(time.Second * 1)
			if bitcoinService.IsRunning() {
				response.Result = bitcoindStopErr
			} else {
				response.Ok = true
				response.Result = bitcoindStopSuccess
			}
		}
	case `status`:
		response.Ok = true
		if bitcoinService.IsRunning() {
			response.Result = bitcoindStarted
		} else {
			response.Result = bitcoindStopped
		}
	case `info`:
		if !bitcoinService.IsRunning() {
			response.Result = bitcoindAlreadyStopErr
		} else {
			if info, err := bitcoinService.GetInfo(); nil != err {
				response.Result = err
			} else {
				var jsonInfo bitcoindInfoResponse
				if err := json.Unmarshal(info, &jsonInfo); nil != err {
					log.Errorf("Error: %v", err)
					if err := writeApiResponseAndSetCookie(w, response); nil != err {
						log.Errorf("Error: %v", err)
					}
					return
				}

				response.Ok = true
				response.Result = jsonInfo
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
