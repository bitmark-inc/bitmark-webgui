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
	Address          string  `json:"address"`
	EstimatedBalance float64 `json:"estimated_balance"`
	AvailableBalance float64 `json:"available_balance"`
}

// POST /api/bitmarkPay/info, pay, encrypt, status, result
// for info, pay, encrypt. They are calling bitmark-pay, will return a job hash,
// use status to know if the job is still running, when the status shows success, fail, stopped,
// use the job hash and result api to get the result
// if the result is nil, it means the async job failed.
func BitmarkPay(w http.ResponseWriter, req *http.Request, log *logger.L, command string) {
	log.Infof("POST /api/bitmarkPay/%s", command)
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	switch command {
	case "info":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		err := bitmarkPayService.Info(*request)
		if nil != err {
			response.Result = "bitmark-pay info error"
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}

	case "pay":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		err := bitmarkPayService.Pay(*request)
		if nil != err {
			response.Result = "bitmark-pay pay error"
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}
	case "encrypt":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		err := bitmarkPayService.Encrypt(*request)
		if nil != err {
			response.Result = "bitmark-pay encrypt error"
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}
	case "status":
		status := bitmarkPayService.Status()
		response.Ok = true
		response.Result = status
	case "result":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		output, err := bitmarkPayService.GetBitmarkPayJobResult(*request)
		if nil != err {
			response.Result = "bitmark-pay result error"
		} else {
			if bitmarkPayService.GetBitmarkPayJobType(request.JobHash) == "info" {
				var jsonInfo BitmarkPayInfoResponse
				if err := json.Unmarshal(output, &jsonInfo); nil != err {
					log.Errorf("Error: %v", err)
				} else {
					response.Ok = true
					response.Result = jsonInfo
				}
			} else {
				response.Ok = true
				response.Result = "success"
			}
		}
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

func bitmarkPayParseRequest(w http.ResponseWriter, req *http.Request, response *Response, log *logger.L) *services.BitmarkPayType {
	var decoder *json.Decoder
	var request services.BitmarkPayType

	decoder = json.NewDecoder(req.Body)
	err := decoder.Decode(&request)
	if nil != err {
		log.Errorf("Error: %v", err)
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return nil
	}

	return &request
}
