// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/bitmark-mgmt/utils"
	"net/http"
	"os/exec"
)

type bitmarkdRequest struct {
	Option string
}

// POST /api/bitmarkd
func Bitmarkd(w http.ResponseWriter, req *http.Request, bitmarkConfigFile string) {

	fmt.Println("POST /api/bitmarkd")
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	decoder := json.NewDecoder(req.Body)
	var request bitmarkdRequest
	err := decoder.Decode(&request)
	if nil != err {
		fmt.Printf("Error:%v\n", err)
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	apiErr := fault.ApiErrInvalidValue
	switch request.Option {
	case `start`:
		// Check bitmarkConfigFile exists
		if !utils.EnsureFileExists(bitmarkConfigFile) {
			fmt.Printf("Error: %v\n", fault.ErrNotFoundConfigFile)
			response.Result = fault.ApiErrStartBitmarkd
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				fmt.Printf("Error: %v\n", err)
			}
			return
		}
		apiErr = fault.ApiErrStartBitmarkd
	case `stop`:
		apiErr = fault.ApiErrStopBitmarkd
	case `status`:
		apiErr = fault.ApiErrStatusBitmarkd
	default:
		response.Result = apiErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	// Use service to run bitmarkd
	cmd := exec.Command("sudo", "service", "bitmarkd", request.Option)
	out, err := cmd.Output()
	if nil != err && request.Option != `status` {
		fmt.Printf("exec service command fail: %v\n", err)
		response.Result = apiErr
	} else {
		fmt.Printf("output: %v\n", string(out))
		response.Ok = true
		response.Result = string(out)
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		fmt.Printf("Error: %v\n", err)
	}

}
