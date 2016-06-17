// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package services

import (
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/logger"
	"os/exec"
	"strconv"
	"sync"
)

type BitmarkCli struct {
	sync.RWMutex
	initialised bool
	log         *logger.L
}

func (bitmarkCli *BitmarkCli) Initialise() error {
	bitmarkCli.Lock()
	defer bitmarkCli.Unlock()

	if bitmarkCli.initialised {
		return fault.ErrAlreadyInitialised
	}

	bitmarkCli.log = logger.New("service-bitmarkCli")
	if nil == bitmarkCli.log {
		return fault.ErrInvalidLoggerChannel
	}

	bitmarkCli.initialised = true

	return nil
}

func (bitmarkCli *BitmarkCli) Finalise() error {
	bitmarkCli.Lock()
	defer bitmarkCli.Unlock()

	if !bitmarkCli.initialised {
		return fault.ErrNotInitialised
	}

	bitmarkCli.initialised = false
	return nil
}

func (bitmarkCli *BitmarkCli) Generate() ([]byte, error) {
	out, err := exec.Command("bitmark-cli", "generate").Output()
	if err != nil {
		bitmarkCli.log.Infof("fail to generate bitmark keypair")
		return nil, err
	}

	return out, nil
}

type BitmarkCliInfoType struct {
	Config string `json:"config"`
}

func (bitmarkCli *BitmarkCli) Info(bitmarkCliInfo BitmarkCliInfoType) ([]byte, error) {
	if err := checkRequireStringParameters(bitmarkCliInfo.Config); nil != err {
		return nil, err
	}

	cmd := exec.Command("bitmark-cli",
		"--config", bitmarkCliInfo.Config,
		"info")

	return getCmdOutput(cmd, "setup", bitmarkCli.log, true)
}

type BitmarkCliSetupType struct {
	Config      string `json:"config"`
	Identity    string `json:"identity"`
	Password    string `json:"password"`
	Network     string `json:"network"`
	Connect     string `json:"connect"`
	Description string `json:"description"`
	PrivateKey  string `json:"private_key"`
}

func (bitmarkCli *BitmarkCli) Setup(bitmarkCliSetup BitmarkCliSetupType, filePath string, bitmarkWebguiConfig *configuration.Configuration) ([]byte, error) {
	if err := checkRequireStringParameters(bitmarkCliSetup.Config, bitmarkCliSetup.Identity, bitmarkCliSetup.Password, bitmarkCliSetup.Network, bitmarkCliSetup.Connect, bitmarkCliSetup.Description); nil != err {
		return nil, err
	}

	cmd := exec.Command("bitmark-cli",
		"--config", bitmarkCliSetup.Config,
		"--identity", bitmarkCliSetup.Identity,
		"--password", bitmarkCliSetup.Password,
		"setup",
		"--network", bitmarkCliSetup.Network,
		"--connect", bitmarkCliSetup.Connect,
		"--description", bitmarkCliSetup.Description,
		"--privateKey", bitmarkCliSetup.PrivateKey)

	output, err := getCmdOutput(cmd, "setup", bitmarkCli.log, true)
	if nil == err {
		bitmarkWebguiConfig.BitmarkCliConfigFile = bitmarkCliSetup.Config
		configuration.UpdateConfiguration(filePath, bitmarkWebguiConfig)
	}
	return output, err
}

type BitmarkCliIssueType struct {
	Config      string `json:"config"`
	Identity    string `json:"identity"`
	Password    string `json:"password"`
	Asset       string `json:"asset"`
	Description string `json:"description"`
	Fingerprint string `json:"fingerprint"`
	Quantity    int    `json:"quantity"`
}

func (bitmarkCli *BitmarkCli) Issue(bitmarkCliIssue BitmarkCliIssueType) ([]byte, error) {
	quantity := strconv.Itoa(bitmarkCliIssue.Quantity)
	if err := checkRequireStringParameters(bitmarkCliIssue.Config, bitmarkCliIssue.Identity, bitmarkCliIssue.Password, bitmarkCliIssue.Asset, bitmarkCliIssue.Description, bitmarkCliIssue.Fingerprint, quantity); nil != err {
		return nil, err
	}

	cmd := exec.Command("bitmark-cli",
		"--config", bitmarkCliIssue.Config,
		"--identity", bitmarkCliIssue.Identity,
		"--password", bitmarkCliIssue.Password,
		"issue",
		"--asset", bitmarkCliIssue.Asset,
		"--description", bitmarkCliIssue.Description,
		"--fingerprint", bitmarkCliIssue.Fingerprint,
		"--quantity", quantity)

	return getCmdOutput(cmd, "issue", bitmarkCli.log, true)
}

type BitmarkCliTransferType struct {
	Config   string `json:"config"`
	Identity string `json:"identity"`
	Password string `json:"password"`
	Txid     string `json:"txid"`
	Receiver string `json:"receiver"`
}

func (bitmarkCli *BitmarkCli) Transfer(bitmarkCliTransfer BitmarkCliTransferType) ([]byte, error) {
	if err := checkRequireStringParameters(bitmarkCliTransfer.Config, bitmarkCliTransfer.Identity, bitmarkCliTransfer.Password, bitmarkCliTransfer.Txid, bitmarkCliTransfer.Receiver); nil != err {
		return nil, err
	}

	cmd := exec.Command("bitmark-cli",
		"--config", bitmarkCliTransfer.Config,
		"--identity", bitmarkCliTransfer.Identity,
		"--password", bitmarkCliTransfer.Password,
		"transfer",
		"--txid", bitmarkCliTransfer.Txid,
		"--receiver", bitmarkCliTransfer.Receiver)

	return getCmdOutput(cmd, "transfer", bitmarkCli.log, true)
}

type BitmarkCliKeyPairType struct {
	Password string `json:"password"`
}

func (bitmarkCli *BitmarkCli) KeyPair(bitmarkCliKeyPair BitmarkCliKeyPairType, bitmarkWebguiConfig string) ([]byte, error) {
	if err := checkRequireStringParameters(bitmarkCliKeyPair.Password); nil != err {
		return nil, err
	}

	cmd := exec.Command("bitmark-cli",
		"--config", bitmarkWebguiConfig,
		"--password", bitmarkCliKeyPair.Password,
		"keypair")

	return getCmdOutput(cmd, "keypair", bitmarkCli.log, false)
}
