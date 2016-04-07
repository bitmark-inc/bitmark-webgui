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
	"os"
	"io/ioutil"
	"io"
	"time"
)

type bitmarkdRequest struct {
	Option string
}

var cmd *exec.Cmd

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

	if nil == cmd {
		cmd = exec.Command("bitmarkd", "--config-file="+bitmarkConfigFile)
	}

	apiErr := fault.ApiErrInvalidValue
	switch request.Option {
	case `start`:
		// Check if bitmarkd is running
		if bitmarkdIsRunning() {
			response.Result = fault.ApiErrAlreadyStartBitmarkd
		}else{
			// Check bitmarkConfigFile exists
			if !utils.EnsureFileExists(bitmarkConfigFile) {
				fmt.Printf("Error: %v\n", fault.ErrNotFoundConfigFile)
				response.Result = fault.ApiErrStartBitmarkd
				if err := writeApiResponseAndSetCookie(w, response); nil != err {
					fmt.Printf("Error: %v\n", err)
				}
				return
			}

			// start bitmarkd as sub process
			stderr, err := cmd.StderrPipe()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			if err := cmd.Start(); nil != err {
				fmt.Printf("Error: %v\n", err)
				response.Result = fault.ApiErrStartBitmarkd
				if err := writeApiResponseAndSetCookie(w, response); nil != err {
					fmt.Printf("Error: %v\n", err)
				}
				cmd = nil
				return
			}

			fmt.Printf("running bitmarkd in pid: %d\n",cmd.Process.Pid)

			bitmarkdProcessErr := runBitmarkdProcess(cmd, stderr, stdout)
			if nil != <- bitmarkdProcessErr {
				fmt.Printf("Exited: %v\n", cmd.ProcessState.Exited())
				response.Result = fault.ApiErrStartBitmarkd
				if err := writeApiResponseAndSetCookie(w, response); nil != err {
					fmt.Printf("Error: %v\n", err)
				}
				cmd = nil
				return
			}

			response.Ok = true
			response.Result = "start running bitmarkd"
		}


	case `stop`:
		if !bitmarkdIsRunning(){
			response.Result = "bitmarkd is not running"
		}else{
			err := cmd.Process.Signal(os.Interrupt)
			if nil != err {
				cmd.Process.Signal(os.Kill)
			}

			response.Ok = true
			response.Result = "stop running bitmarkd"
			cmd = nil
		}
	case `status`:
		response.Ok = true
		if bitmarkdIsRunning() {
			response.Result = "bitmarkd is running"

		}else{
			response.Result = "bitmarkd is not running"
		}
	default:
		response.Result = apiErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		fmt.Printf("Error: %v\n", err)
	}

}

func bitmarkdIsRunning() bool {
	if nil == cmd.Process {
		return false
	}

	return true
}

func runBitmarkdProcess(cmd *exec.Cmd, stderr io.ReadCloser, stdout io.ReadCloser) <-chan error {
	err := make(chan error, 1)

	go func() {
		stde, e := ioutil.ReadAll(stderr)
		if nil != e {
			fmt.Printf("Error: %v\n", err)
		}

		stdo, e := ioutil.ReadAll(stdout)
		if nil != e {
			fmt.Printf("Error: %v\n", err)
		}

		fmt.Printf("Error: %s\n", stde)
		fmt.Printf("Out: %s\n", stdo)

		if e := cmd.Wait(); nil != e {
			err <- e
		}
	}()


	// wait for 1 second if cmd has no error then return nil
	go func() {
		time.Sleep(time.Second * 1)
		err <- nil
		close(err)
	}()

	return err
}
