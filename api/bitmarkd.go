// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/bitmark-webgui/structs"
	bitmarkdConfig "github.com/bitmark-inc/bitmarkd/configuration"
	"github.com/bitmark-inc/bitmarkd/rpc"
	"github.com/bitmark-inc/logger"
	"net/http"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"time"
)

type bitmarkdRequest struct {
	Option  string `json:"option"`
	Network string `json:"network"`
	// ConfigFile string `json:"config_file"`
}

// POST /api/bitmarkd
func Bitmarkd(w http.ResponseWriter, req *http.Request, webguiFilePath string, webguiConfig *configuration.Configuration, log *logger.L) {

	log.Info("POST /api/bitmarkd")
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
			if info, err := getBitmarkdInfo(webguiConfig.BitmarkConfigFile, log); "" != err {
				response.Result = err
			} else {
				response.Ok = true
				response.Result = info
			}
		}
	case `setup`:
		if bitmarkService.IsRunning() {
			response.Result = bitmarkdAlreadyStartErr
		} else {
			bitmarkConfigFile := filepath.Join(webguiConfig.DataDirectory, fmt.Sprintf("bitmarkd-%s", request.Network), "bitmarkd.conf")
			if err := bitmarkService.Setup(bitmarkConfigFile, request.Network, webguiFilePath, webguiConfig); nil != err {
				if os.IsNotExist(err) {
					response.Result = "bitmarkd config not found"
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

func getBitmarkdInfo(bitmarkConfigFile string, log *logger.L) (*rpc.InfoReply, string) {
	bitmarkConfig := &structs.BitmarkdConfiguration{}
	err := bitmarkdConfig.ParseConfigurationFile(bitmarkConfigFile, bitmarkConfig)
	if nil != err {
		log.Errorf("Failed to get bitmarkd configuration: %v", err)
		return nil, bitmarkdGetConfigErr
	}

	conn, err := bitmarkService.Connect(bitmarkConfig.ClientRPC.Listen[0])
	if nil != err {
		log.Errorf("Failed to connect to bitmarkd: %v", err)
		return nil, bitmarkdConnectErr
	}
	defer conn.Close()

	// create a client
	client := jsonrpc.NewClient(conn)
	defer client.Close()

	info, err := bitmarkService.GetInfo(client)
	if nil != err {
		log.Errorf("Failed to get bitmark info: %v", err)
		return nil, bitmarkdGetInfoErr
	}

	return info, ""
}
