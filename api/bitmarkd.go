// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/bitmark-mgmt/utils"
	"github.com/bitmark-inc/logger"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type bitmarkdRequest struct {
	Option string
}

var cmd *exec.Cmd

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

	if nil == cmd {
		cmd = exec.Command("bitmarkd", "--config-file="+bitmarkConfigFile)
	}

	apiErr := fault.ApiErrInvalidValue
	switch request.Option {
	case `start`:
		// Check if bitmarkd is running
		if bitmarkdIsRunning() {
			response.Result = fault.ApiErrAlreadyStartBitmarkd
		} else {
			// Check bitmarkConfigFile exists
			if !utils.EnsureFileExists(bitmarkConfigFile) {
				log.Errorf("Error: %v", fault.ErrNotFoundConfigFile)
				response.Result = fault.ApiErrStartBitmarkd
				if err := writeApiResponseAndSetCookie(w, response); nil != err {
					log.Errorf("Error: %v", err)
				}
				return
			}

			// start bitmarkd as sub process
			stderr, err := cmd.StderrPipe()
			if err != nil {
				log.Errorf("Error: %v", err)
			}
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Errorf("Error: %v", err)
			}

			if err := cmd.Start(); nil != err {
				log.Errorf("Error: %v", err)
				response.Result = fault.ApiErrStartBitmarkd
				if err := writeApiResponseAndSetCookie(w, response); nil != err {
					log.Errorf("Error: %v", err)
				}
				cmd = nil
				return
			}

			log.Errorf("running bitmarkd in pid: %d\n", cmd.Process.Pid)

			bitmarkdProcessErr := runBitmarkdProcess(cmd, stderr, stdout, log)
			if nil != <-bitmarkdProcessErr {
				log.Errorf("Exited: %v", cmd.ProcessState.Exited())
				response.Result = fault.ApiErrStartBitmarkd
				if err := writeApiResponseAndSetCookie(w, response); nil != err {
					log.Errorf("Error: %v", err)
				}
				cmd = nil
				return
			}

			response.Ok = true
			response.Result = "start running bitmarkd"
		}

	case `stop`:
		if !bitmarkdIsRunning() {
			response.Result = "bitmarkd is not running"
		} else {
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

		} else {
			response.Result = "bitmarkd is not running"
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

func bitmarkdIsRunning() bool {
	if nil == cmd.Process {
		return false
	}

	return true
}

func runBitmarkdProcess(cmd *exec.Cmd, stderr io.ReadCloser, stdout io.ReadCloser, log *logger.L) <-chan error {
	err := make(chan error, 1)

	go func() {
		stde, e := ioutil.ReadAll(stderr)
		if nil != e {
			log.Errorf("Error: %v", err)
		}

		stdo, e := ioutil.ReadAll(stdout)
		if nil != e {
			log.Errorf("Error: %v", err)
		}

		log.Errorf("Error: %s\n", stde)
		log.Errorf("Out: %s\n", stdo)

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
