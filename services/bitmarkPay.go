// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/utils"
	"github.com/bitmark-inc/logger"
	"os/exec"
	"sync"
)

type BitmarkPay struct {
	sync.RWMutex
	initialised bool
	bin         string
	log         *logger.L
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

	bitmarkPay.initialised = false
	return nil
}

type BitmarkPayPwdType struct {
	Net       string
	Config    string
	Password  string
	Txid      string
	Addresses []string
}

func (bitmarkPay *BitmarkPay) Encrypt(bitmarkPayPwd BitmarkPayPwdType) ([]byte, error) {
	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayPwd.Net,
		"--config="+bitmarkPayPwd.Config,
		"--password="+bitmarkPayPwd.Password,
		"encrypt")

	return getCmdOutput(cmd, "encrypt", bitmarkPay.log)
}

func (bitmarkPay *BitmarkPay) Info(bitmarkPayPwd BitmarkPayPwdType) ([]byte, error) {

	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayPwd.Net,
		"--config="+bitmarkPayPwd.Config,
		"--json",
		"info")

	return getCmdOutput(cmd, "info", bitmarkPay.log)
}

func (bitmarkPay *BitmarkPay) Pay(bitmarkPayPwd BitmarkPayPwdType) ([]byte, error) {
	addresses := ""
	for _, payAddress := range bitmarkPayPwd.Addresses {
		addresses = addresses + " " + payAddress
	}

	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayPwd.Net,
		"--config="+bitmarkPayPwd.Config,
		"--password="+bitmarkPayPwd.Password,
		"pay",
		bitmarkPayPwd.Txid,
		addresses,
	)

	return getCmdOutput(cmd, "pay", bitmarkPay.log)

}
