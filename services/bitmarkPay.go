// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/utils"
	"github.com/bitmark-inc/logger"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type BitmarkPay struct {
	sync.RWMutex
	initialised bool
	bin         string
	log         *logger.L
	command *exec.Cmd
}

func (bitmarkPay *BitmarkPay) Initialise(binFile string) error {
	bitmarkPay.Lock()
	defer bitmarkPay.Unlock()

	if bitmarkPay.initialised {
		return fault.ErrAlreadyInitialised
	}

	bitmarkPay.log = logger.New("service-bitmarkPay")
	if nil == bitmarkPay.log {
		return fault.ErrInvalidLoggerChannel
	}

	// Check bitmarkPay bin exists
	if !utils.EnsureFileExists(binFile) {
		bitmarkPay.log.Errorf("cannot find bitmarkPay bin: %s", binFile)
		return fault.ErrNotFoundBinFile
	}
	bitmarkPay.bin = binFile

	bitmarkPay.initialised = true

	return nil
}

func (bitmarkPay *BitmarkPay) Finalise() error {
	bitmarkPay.Lock()
	defer bitmarkPay.Unlock()

	if !bitmarkPay.initialised {
		return fault.ErrNotInitialised
	}

	if nil != bitmarkPay.command {
		if err := bitmarkPay.Kill(); nil != err {
			return err
		}
	}

	bitmarkPay.initialised = false
	return nil
}

type BitmarkPayType struct {
	Net       string   `json:"net"`
	Config    string   `json:"config"`
	Password  string   `json:"password"`
	Txid      string   `json:"txid"`
	Addresses []string `json:"addresses"`
}

func (bitmarkPay *BitmarkPay) Encrypt(bitmarkPayType BitmarkPayType) ([]byte, error) {
	//check command process is not running
	oldCmd := bitmarkPay.command
	if nil != oldCmd && nil == oldCmd.ProcessState {
		return nil, fault.ErrBitmarkPayIsRunning
	}

	// check config, net, password
	if err := checkRequireStringParameters(bitmarkPayType.Config, bitmarkPayType.Net, bitmarkPayType.Password); nil != err {
		return nil, err
	}

	// check cmd process is finish
	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayType.Net,
		"--config="+bitmarkPayType.Config,
		"--password="+bitmarkPayType.Password,
		"encrypt")

	bitmarkPay.command = cmd
	return getCmdOutput(cmd, "encrypt", bitmarkPay.log)
}

func (bitmarkPay *BitmarkPay) Info(bitmarkPayType BitmarkPayType) ([]byte, error) {
	//check command process is not running
	oldCmd := bitmarkPay.command
	if nil != oldCmd && nil == oldCmd.ProcessState {
		return nil, fault.ErrBitmarkPayIsRunning
	}

	// check config, net
	if err := checkRequireStringParameters(bitmarkPayType.Config, bitmarkPayType.Net); nil != err {
		return nil, err
	}

	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayType.Net,
		"--config="+bitmarkPayType.Config,
		"--json",
		"info")

	bitmarkPay.command = cmd
	return getCmdOutput(cmd, "info", bitmarkPay.log)
}

func (bitmarkPay *BitmarkPay) Pay(bitmarkPayType BitmarkPayType) ([]byte, error) {
	//check command process is not running
	oldCmd := bitmarkPay.command
	if nil != oldCmd && nil == oldCmd.ProcessState {
		return nil, fault.ErrBitmarkPayIsRunning
	}

	// check config, net, password, txid, addresses
	addresses := strings.Join(bitmarkPayType.Addresses, " ")
	if err := checkRequireStringParameters(bitmarkPayType.Config, bitmarkPayType.Net, bitmarkPayType.Password, bitmarkPayType.Txid, addresses); nil != err {
		return nil, err
	}

	bitmarkPay.log.Tracef("txid: %s", bitmarkPayType.Txid)
	bitmarkPay.log.Tracef("addresses: %s", addresses)
	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayType.Net,
		"--config="+bitmarkPayType.Config,
		"--password="+bitmarkPayType.Password,
		"pay",
		bitmarkPayType.Txid,
		addresses,
	)


	bitmarkPay.command = cmd
	return getCmdOutput(cmd, "pay", bitmarkPay.log)
}

func (bitmarkPay *BitmarkPay) Status() string{
	cmd := bitmarkPay.command
	if nil != cmd{
		if  nil != cmd.ProcessState {
			if cmd.ProcessState.Exited() {
				if cmd.ProcessState.Success(){
					return "success"
				}
				return "fail"
			}else {
				return cmd.ProcessState.String()
			}
		} else{
			return "running"
		}


	}
	return "stopped"
}

func (bitmarkPay *BitmarkPay) Kill() error {
	cmd := bitmarkPay.command
	if nil != cmd {
		bitmarkPay.log.Debugf("killing process: %d", cmd.Process.Pid)

		err := cmd.Process.Signal(syscall.SIGINT)
		if nil != err {
			bitmarkPay.log.Errorf("Failed to Kill process: %d", cmd.Process.Pid)
			return err
		}
	}
	bitmarkPay.command = nil
	return nil
}
