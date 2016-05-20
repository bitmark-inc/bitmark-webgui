// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/services"
	"github.com/bitmark-inc/logger"
	"net/http"
)

type BitmarkPayInfoResponse struct {
	Address           string
	Estimated_balance float32
	Available_balance float32
}

// POST /api/bitmarkPay/*
func BitmarkPayEncrypt(w http.ResponseWriter, req *http.Request, log *logger.L, command string) {
	log.Info("POST /api/bitmarkPay/encrypt")
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	var decoder *json.Decoder
	var request services.BitmarkPayType

	if "status" != command {
		decoder = json.NewDecoder(req.Body)
		err := decoder.Decode(&request)
		if nil != err {
			log.Errorf("Error: %v", err)
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}

	}

	switch command {
	case "info":
		output, err := bitmarkPayService.Info(request)
		if nil != err {
			response.Result = "bitmark-pay info error"
		} else {

			var jsonInfo BitmarkPayInfoResponse
			if err := json.Unmarshal(output, &jsonInfo); nil != err {
				log.Errorf("Error: %v", err)
			} else {
				response.Ok = true
				response.Result = jsonInfo
			}
		}

	case "pay":
		_, err := bitmarkPayService.Pay(request)
		if nil != err {
			response.Result = "bitmark-pay pay error"
		} else {
			response.Ok = true
			response.Result = "Success"
		}
	case "encrypt":
		_, err := bitmarkPayService.Encrypt(request)
		if nil != err {
			response.Result = "bitmark-pay encrypt error"
		} else {
			response.Ok = true
			response.Result = "Success"
		}
	case "status":
		status := bitmarkPayService.Status()
		response.Ok = true
		response.Result = status
	}


	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}
