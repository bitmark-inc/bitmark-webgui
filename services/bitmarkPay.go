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

type BitmarkPayType struct {
	Net       string   `json:"net"`
	Config    string   `json:"config"`
	Password  string   `json:"password"`
	Txid      string   `json:"txid"`
	Addresses []string `json:"addresses"`
}

func (bitmarkPay *BitmarkPay) Encrypt(bitmarkPayType BitmarkPayType) ([]byte, error) {
	// check config, net, password
	if err := checkRequireStringParameters(bitmarkPayType.Config, bitmarkPayType.Net, bitmarkPayType.Password); nil != err {
		return nil, err
	}

	cmd := exec.Command("java", "-jar",
		"-Dorg.apache.logging.log4j.simplelog.StatusLogger.level=OFF",
		bitmarkPay.bin,
		"--net="+bitmarkPayType.Net,
		"--config="+bitmarkPayType.Config,
		"--password="+bitmarkPayType.Password,
		"encrypt")

	return getCmdOutput(cmd, "encrypt", bitmarkPay.log)
}

func (bitmarkPay *BitmarkPay) Info(bitmarkPayType BitmarkPayType) ([]byte, error) {
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

	return getCmdOutput(cmd, "info", bitmarkPay.log)
}

func (bitmarkPay *BitmarkPay) Pay(bitmarkPayType BitmarkPayType) ([]byte, error) {
	addresses := strings.Join(bitmarkPayType.Addresses, " ")

	// check config, net, password, txid, addresses
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

	return getCmdOutput(cmd, "pay", bitmarkPay.log)
}
