// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"encoding/hex"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/utils"
	"github.com/bitmark-inc/logger"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type BitmarkPay struct {
	sync.RWMutex
	initialised bool
	bin         string
	log         *logger.L
	asyncJob    BitmarkPayJob
	// command *exec.Cmd
}

type BitmarkPayJob struct {
	hash    string
	command *exec.Cmd
	cmdType string
	result  []byte
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

	if nil != bitmarkPay.asyncJob.command {
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
	JobHash   string   `json:"job_hash"`
}

func (bitmarkPay *BitmarkPay) Encrypt(bitmarkPayType BitmarkPayType) error {
	//check command process is not running
	oldCmd := bitmarkPay.asyncJob.command
	if nil != oldCmd && nil == oldCmd.ProcessState {
		return fault.ErrBitmarkPayIsRunning
	}

	// check config, net, password
	if err := checkRequireStringParameters(bitmarkPayType.Config, bitmarkPayType.Net, bitmarkPayType.Password); nil != err {
		return err
	}

	// check cmd process is finish
	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayType.Net,
		"--config="+bitmarkPayType.Config,
		"--password="+bitmarkPayType.Password,
		"encrypt")

	return bitmarkPay.runBitmarkPayJob(cmd, "encrypt")
}

func (bitmarkPay *BitmarkPay) Info(bitmarkPayType BitmarkPayType) error {
	//check command process is not running
	oldCmd := bitmarkPay.asyncJob.command
	if nil != oldCmd && nil == oldCmd.ProcessState {
		return fault.ErrBitmarkPayIsRunning
	}

	// check config, net
	if err := checkRequireStringParameters(bitmarkPayType.Config, bitmarkPayType.Net); nil != err {
		return err
	}

	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayType.Net,
		"--config="+bitmarkPayType.Config,
		"--json",
		"info")

	return bitmarkPay.runBitmarkPayJob(cmd, "info")
}

func (bitmarkPay *BitmarkPay) Pay(bitmarkPayType BitmarkPayType) error {
	//check command process is not running
	oldCmd := bitmarkPay.asyncJob.command
	if nil != oldCmd && nil == oldCmd.ProcessState {
		return fault.ErrBitmarkPayIsRunning
	}

	// check config, net, password, txid, addresses
	addresses := strings.Join(bitmarkPayType.Addresses, " ")
	if err := checkRequireStringParameters(bitmarkPayType.Config, bitmarkPayType.Net, bitmarkPayType.Password, bitmarkPayType.Txid, addresses); nil != err {
		return err
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

	return bitmarkPay.runBitmarkPayJob(cmd, "pay")
}

func (bitmarkPay *BitmarkPay) Status() string {
	cmd := bitmarkPay.asyncJob.command
	if nil != cmd {
		if nil != cmd.ProcessState {
			if cmd.ProcessState.Exited() {
				if cmd.ProcessState.Success() {
					return "success"
				}
				return "fail"
			} else {
				return cmd.ProcessState.String()
			}
		} else {
			return "running"
		}

	}
	return "stopped"
}

func (bitmarkPay *BitmarkPay) Kill() error {
	cmd := bitmarkPay.asyncJob.command
	if nil != cmd {
		bitmarkPay.log.Debugf("killing process: %d", cmd.Process.Pid)

		err := cmd.Process.Signal(syscall.SIGINT)
		if nil != err {
			bitmarkPay.log.Errorf("Failed to Kill process: %d", cmd.Process.Pid)
			return err
		}
	}

	bitmarkPay.asyncJob.hash = ""
	bitmarkPay.asyncJob.command = nil
	bitmarkPay.asyncJob.cmdType = ""
	bitmarkPay.asyncJob.result = nil
	return nil
}

func (bitmarkPay *BitmarkPay) runBitmarkPayJob(cmd *exec.Cmd, cmdType string) error {
	byteHash, err := time.Now().MarshalText()
	if nil != err {
		return err
	}

	hash := hex.EncodeToString(byteHash)
	bitmarkPay.asyncJob.hash = hash
	bitmarkPay.asyncJob.command = cmd
	bitmarkPay.asyncJob.cmdType = cmdType

	go func() {
		if result, err := getCmdOutput(cmd, cmdType, bitmarkPay.log); nil != err {
			bitmarkPay.log.Errorf("job fail: %s", bitmarkPay.asyncJob.hash)
		} else {
			bitmarkPay.asyncJob.result = result
		}
	}()

	return nil

}

func (bitmarkPay *BitmarkPay) GetBitmarkPayJobHash() string {
	return bitmarkPay.asyncJob.hash
}

func (bitmarkPay *BitmarkPay) GetBitmarkPayJobResult(bitmarkPayType BitmarkPayType) ([]byte, error) {

	// check config, net, password
	if err := checkRequireStringParameters(bitmarkPayType.JobHash); nil != err {
		return nil, err
	}

	if bitmarkPayType.JobHash != bitmarkPay.asyncJob.hash {
		return nil, fault.ErrNotFoundBitmarkPayJob
	}

	if bitmarkPay.Status() == "running" {
		return nil, fault.ErrInvalidAccessBitmarkPayJobResult
	}

	return bitmarkPay.asyncJob.result, nil
}

func (bitmarkPay *BitmarkPay) GetBitmarkPayJobType(hashString string) string {
	return bitmarkPay.asyncJob.cmdType
}
