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

		if err := bitmarkPayService.Info(*request); nil != err {
			response.Result = err
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}

	case "pay":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		if err := bitmarkPayService.Pay(*request); nil != err {
			response.Result = err
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}
	case "restore":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		if err := bitmarkPayService.Restore(*request); nil != err {
			response.Result = err
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}
	case "encrypt":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		if err := bitmarkPayService.Encrypt(*request); nil != err {
			response.Result = err
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}
	case "decrypt":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		if err := bitmarkPayService.Decrypt(*request); nil != err {
			response.Result = err
		} else {
			response.Ok = true
			response.Result = bitmarkPayService.GetBitmarkPayJobHash()
		}
	case "status":
		request := bitmarkPayParseRequest(w, req, response, log)
		if status, err := bitmarkPayService.Status(request.JobHash); nil != err {
			response.Result = err
		} else {
			response.Ok = true
			response.Result = status
		}

	case "result":
		request := bitmarkPayParseRequest(w, req, response, log)
		if nil == request {
			return
		}

		if output, err := bitmarkPayService.GetBitmarkPayJobResult(*request); nil != err {
			response.Result = err
		} else {
			jobType := bitmarkPayService.GetBitmarkPayJobType(request.JobHash)
			switch jobType {
			case "info":
				var jsonInfo BitmarkPayInfoResponse
				if err := json.Unmarshal(output, &jsonInfo); nil != err {
					response.Result = err
					log.Errorf("Error: %v", err)
				} else {
					response.Ok = true
					response.Result = jsonInfo
				}
			default:
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
	if err := decoder.Decode(&request); nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return nil
	}

	return &request
}

// Get /api/bitmarkPay
func BitmarkPayJobHash(w http.ResponseWriter, req *http.Request, log *logger.L) {
	log.Infof("GET /api/bitmarkPay")

	response := &Response{
		Ok:     true,
		Result: bitmarkPayService.GetBitmarkPayJobHash(),
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

// DELETE /api/bitmarkPay
func BitmarkPayKill(w http.ResponseWriter, req *http.Request, log *logger.L) {
	log.Infof("DELETE /api/bitmarkPay")
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	request := bitmarkPayParseRequest(w, req, response, log)
	if nil == request {
		return
	}

	if request.JobHash != bitmarkPayService.GetBitmarkPayJobHash() {
		response.Result = "invalid job hash"
		log.Errorf("Invalid job hash, requset: %s, in service: %s", request.JobHash, bitmarkPayService.GetBitmarkPayJobHash())
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	log.Infof("Delete job hash: %s", request.JobHash)
	err := bitmarkPayService.Kill()
	if nil != err {
		response.Result = err
	} else {
		response.Ok = true
		response.Result = "success"
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}
