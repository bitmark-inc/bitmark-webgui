// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"errors"
	"github.com/bitmark-inc/bitmarkd/configuration"
	"github.com/bitmark-inc/bitmarkd/rpc"
	"github.com/bitmark-inc/logger"
	"net/http"
	"net/rpc/jsonrpc"
	"time"
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

	apiErr := invalidValueErr
	switch request.Option {
	case `start`:
		// Check if bitmarkd is running
		if bitmarkService.IsRunning() {
			response.Result = bitmarkdAlreadyStartErr
		} else {
			bitmarkService.ModeStart <- true
			// wait one second to get correct result
			time.Sleep(time.Second * 1)
			if !bitmarkService.IsRunning() {
				response.Result = bitmarkdStartErr
			} else {
				response.Ok = true
				response.Result = bitmarkdStartSuccess
			}
		}
	case `stop`:
		if !bitmarkService.IsRunning() {
			response.Result = bitmarkdAlreadyStopErr
		} else {
			bitmarkService.ModeStart <- false
			time.Sleep(time.Second * 1)
			if bitmarkService.IsRunning() {
				response.Result = bitmarkdStopErr
			} else {
				response.Ok = true
				response.Result = bitmarkdStopSuccess
			}

		}
	case `status`:
		response.Ok = true
		if bitmarkService.IsRunning() {
			response.Result = bitmarkdStarted
		} else {
			response.Result = bitmarkdStopped
		}
	case `info`:
		if !bitmarkService.IsRunning() {
			response.Result = bitmarkdAlreadyStopErr
		} else {
			if info, err := getBitmarkdInfo(bitmarkConfigFile, log); nil != err {
				response.Result = err
			} else {
				response.Ok = true
				response.Result = info
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

func getBitmarkdInfo(bitmarkConfigFile string, log *logger.L) (*rpc.InfoReply, error) {
	bitmarkConfig, err := configuration.GetConfiguration(bitmarkConfigFile)
	if nil != err {
		log.Errorf("Failed to get bitmarkd configuration: %v", err)
		return nil, errors.New(bitmarkdGetConfigErr)
	}

	conn, err := bitmarkService.Connect(bitmarkConfig.ClientRPC.Listen[0])
	if nil != err {
		log.Errorf("Failed to connect to bitmarkd: %v", err)
		return nil, errors.New(bitmarkdConnectErr)
	}
	defer conn.Close()

	// create a client
	client := jsonrpc.NewClient(conn)
	defer client.Close()

	info, err := bitmarkService.GetInfo(client)
	if nil != err {
		log.Errorf("Failed to get bitmark info: %v", err)
		return nil, errors.New(bitmarkdGetInfoErr)
	}

	return info, nil
}
