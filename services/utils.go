// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"errors"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/logger"
	"io/ioutil"
	"os/exec"
)

func getCmdOutput(cmd *exec.Cmd, cmdType string, log *logger.L) ([]byte, error) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	if err := cmd.Start(); nil != err {
		return nil, err
	}

	stde, err := ioutil.ReadAll(stderr)
	if nil != err {
		log.Errorf("Error: %v", err)
	}

	stdo, err := ioutil.ReadAll(stdout)
	if nil != err {
		log.Errorf("Error: %v", err)
	}

	log.Errorf("%s %s stderr: %s", cmd.Path, cmdType, stde)
	log.Infof("%s %s stdout: %s", cmd.Path, cmdType, stdo)
	if len(stde) != 0 {
		return nil, errors.New(string(stde))
	}

	if err := cmd.Wait(); nil != err {
		log.Errorf("%s %s failed: %v", cmd.Path, cmdType, err)
		return nil, err
	}

	return stdo, nil
}

func checkRequireStringParameters(params ...string) error {

	for _, param := range params {
		if "" == param || "0" == param {
			return fault.ErrInvalidCommandParams
		}
	}
	return nil
}
