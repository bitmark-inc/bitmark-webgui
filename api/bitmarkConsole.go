// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"github.com/bitmark-inc/logger"
	"net/http"
)

//POST /api/bitmarkConsole
func BitmarkConsole(w http.ResponseWriter, req *http.Request, log *logger.L) {
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	if err := bitmarkConsoleService.StartBitmarkConsole(); nil != err {
		log.Errorf("start bitmark console failed: %v", err)
		response.Result = err

		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	response.Ok = true
	response.Result = bitmarkConsoleService.GetBitmarkConsoleUrl()

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}

}
