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

type BitmarkPayInterface interface {
	Encrypt(BitmarkPayType) error
	Decrypt(BitmarkPayType) error
	Info(BitmarkPayType) error
	Pay(BitmarkPayType) error
	Status(string) (string, error)
	Kill() error
	GetBitmarkPayJobHash() string
	GetBitmarkPayJobResult(BitmarkPayType) ([]byte, error)
	GetBitmarkPayJobType(string) string
}

type BitmarkPay struct {
	sync.RWMutex
	initialised bool
	bin         string
	log         *logger.L
	asyncJob    BitmarkPayJob
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

	return bitmarkPay.runBitmarkPayJob(cmd, "encrypt", true)
}

func (bitmarkPay *BitmarkPay) Decrypt(bitmarkPayType BitmarkPayType) error {
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
		"decrypt")

	return bitmarkPay.runBitmarkPayJob(cmd, "decrypt", false)
}

func (bitmarkPay *BitmarkPay) Info(bitmarkPayType BitmarkPayType) error {
	//check command process is not running
	oldCmd := bitmarkPay.asyncJob.command
	if nil != oldCmd && nil == oldCmd.ProcessState {
		bitmarkPay.log.Infof("asyncJob: %v", bitmarkPay.asyncJob)
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

	return bitmarkPay.runBitmarkPayJob(cmd, "info", true)
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

	return bitmarkPay.runBitmarkPayJob(cmd, "pay", true)
}

func (bitmarkPay *BitmarkPay) Status(hashString string) (string, error) {
	if hashString != bitmarkPay.asyncJob.hash {
		return "", fault.ErrNotFoundBitmarkPayJob
	}

	return bitmarkPay.status(), nil
}

func (bitmarkPay *BitmarkPay) status() string {
	bitmarkPay.Lock()
	defer bitmarkPay.Unlock()

	cmd := bitmarkPay.asyncJob.command
	if nil != cmd {
		if nil != cmd.Process {
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
		} else { // process is nil means the command is fail or not start yet
			return "stopped"
		}
	}
	return "stopped"
}

func (bitmarkPay *BitmarkPay) Kill() error {

	cmd := bitmarkPay.asyncJob.command
	if nil != cmd {
		bitmarkPay.log.Infof("killing process: %d", cmd.Process.Pid)

		err := cmd.Process.Signal(syscall.SIGINT)
		if nil != err {
			bitmarkPay.log.Errorf("Failed to Kill process: %d", cmd.Process.Pid)
			return err
		}
	}

	go func() {
		// waiting for process killing done and set the variable to nil
		var waitTime = 30 // 30s
		var count = 0
	loop:
		for {
			select {
			case <-time.After(1 * time.Second):
				if nil != cmd.ProcessState {
					bitmarkPay.log.Infof("get signal: %s", cmd.ProcessState.String())
					break loop
				}
				count++
				if count > waitTime {
					bitmarkPay.log.Infof("force kill process: %d", cmd.Process.Pid)
					err := cmd.Process.Signal(syscall.SIGKILL)
					if nil != err {
						bitmarkPay.log.Errorf("Failed to Kill process: %d", cmd.Process.Pid)
					}
					break loop
				}
			}
		}

		bitmarkPay.Lock()
		defer bitmarkPay.Unlock()

		bitmarkPay.asyncJob.hash = ""
		bitmarkPay.asyncJob.command = nil
		bitmarkPay.asyncJob.cmdType = ""
		bitmarkPay.asyncJob.result = nil
	}()

	return nil
}

func (bitmarkPay *BitmarkPay) runBitmarkPayJob(cmd *exec.Cmd, cmdType string, logStdOut bool) error {
	byteHash, err := time.Now().MarshalText()
	if nil != err {
		bitmarkPay.log.Errorf("get error: %v\n", err)
		return err
	}
	hash := hex.EncodeToString(byteHash)

	bitmarkPay.Lock()
	defer bitmarkPay.Unlock()

	bitmarkPay.asyncJob.hash = hash
	bitmarkPay.asyncJob.command = cmd
	bitmarkPay.asyncJob.cmdType = cmdType

	go func() {
		if result, err := getCmdOutput(cmd, cmdType, bitmarkPay.log, logStdOut); nil != err {
			bitmarkPay.Lock()
			defer bitmarkPay.Unlock()

			bitmarkPay.log.Errorf("job fail: %s", bitmarkPay.asyncJob.hash)
			bitmarkPay.asyncJob.result = nil
			bitmarkPay.asyncJob.command.Process.Kill()
			bitmarkPay.asyncJob.command = nil
		} else {
			bitmarkPay.Lock()
			defer bitmarkPay.Unlock()

			// make sure result is no nil, otherwise the api will consider the result command is failed
			if nil == result {
				result = []byte("")
			}
			bitmarkPay.asyncJob.result = result
		}
	}()

	return nil

}

func (bitmarkPay *BitmarkPay) GetBitmarkPayJobHash() string {
	return bitmarkPay.asyncJob.hash
}

func (bitmarkPay *BitmarkPay) GetBitmarkPayJobResult(bitmarkPayType BitmarkPayType) ([]byte, error) {

	// check jobhash
	if err := checkRequireStringParameters(bitmarkPayType.JobHash); nil != err {
		return nil, err
	}

	if bitmarkPayType.JobHash != bitmarkPay.asyncJob.hash {
		return nil, fault.ErrNotFoundBitmarkPayJob
	}

	if bitmarkPay.status() == "running" || bitmarkPay.status() == "stopped" {
		return nil, fault.ErrInvalidAccessBitmarkPayJobResult
	}

	if nil == bitmarkPay.asyncJob.result {
		return nil, fault.ErrExecBitmarkPayJob
	}

	return bitmarkPay.asyncJob.result, nil
}

func (bitmarkPay *BitmarkPay) GetBitmarkPayJobType(hashString string) string {
	return bitmarkPay.asyncJob.cmdType
}
