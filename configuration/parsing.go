// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package configuration

import (
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"os"
	"path/filepath"
)

type Configuration struct {
	Port              int    `libucl:"port"`
	Password          string `libucl:"password"`
	EnableHttps       bool   `libucl:"enable_https"`
	BitmarkConfigFile string `libucl:"bitmark_config_file"`
}

func GetConfigPath(dir string) (string, error) {
	fileInfo, err := os.Stat(dir)
	if nil != err {
		return "", err
	}
	if !fileInfo.IsDir() {
		return "", fault.ErrConfigDirPath
	}

	path := dir + "/bitmark-mgmt.conf"

	return path, nil
}

func GetConfiguration(configurationFileName string) (*Configuration, error) {

	configurationFileName, err := filepath.Abs(filepath.Clean(configurationFileName))
	if nil != err {
		return nil, err
	}

	options := &Configuration{}
	if err := readConfigurationFile(configurationFileName, options); err != nil {
		return nil, err
	}

	return options, nil
}
