// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package configuration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultDataDirectory        = "."
	defaultPort                 = 2150
	defaultPassword             = "bitmark-webgui"
	defaultEnableHttps          = true
	defaultBitmarkChain         = "testing"
	defaultBitmarkConfigFile    = "/etc/bitmarkd.conf"
	defaultProoferdConfigFile   = "/etc/prooferd.conf"
	defaultBitmarkCliConfigFile = ""
	defaultBitmarkPayServiceBin = "./bin/bitmarkPayService"
	defaultBitmarkConsoleBin    = "./bin/gotty"

	defaultLogDirectory = "log"
	defaultLogFile      = "bitmark-webgui.log"
	defaultLogCount     = 10          //  number of log files retained
	defaultLogSize      = 1024 * 1024 // rotate when <logfile> exceeds this size
)

var defaultLogger = &LoggerType{
	Directory: defaultLogDirectory,
	File:      defaultLogFile,
	Size:      defaultLogSize,
	Count:     defaultLogCount,
	Levels: map[string]string{
		"main": "info",
		"api":  "info",
		"*":    "info",
	},
}

type LoggerType struct {
	Directory string            `libucl:"directory"`
	File      string            `libucl:"file"`
	Size      int               `libucl:"size"`
	Count     int               `libucl:"count"`
	Levels    map[string]string `libucl:"levels"`
}

type Configuration struct {
	DataDirectory        string     `libucl:"data_directory"`
	Port                 int        `libucl:"port"`
	Password             string     `libucl:"password"`
	EnableHttps          bool       `libucl:"enable_https"`
	BitmarkChain         string     `libucl:"bitmark_chain"`
	BitmarkConfigFile    string     `libucl:"bitmark_config_file"`
	ProoferdConfigFile   string     `libucl:"prooferd_config_file"`
	BitmarkCliConfigFile string     `libucl:"bitmark_cli_config_file"`
	BitmarkPayServiceBin string     `libucl:"bitmark_pay_service_bin"`
	BitmarkConsoleBin    string     `libucl:"bitmark_console_bin"`
	Logging              LoggerType `libucl:"logging"`
}

func GetDefaultConfiguration(dataDirectory string) (*Configuration, error) {
	config := Configuration{
		DataDirectory:        defaultDataDirectory,
		Port:                 defaultPort,
		Password:             defaultPassword,
		EnableHttps:          defaultEnableHttps,
		BitmarkChain:         defaultBitmarkChain,
		BitmarkConfigFile:    defaultBitmarkConfigFile,
		ProoferdConfigFile:   defaultProoferdConfigFile,
		BitmarkCliConfigFile: defaultBitmarkCliConfigFile,
		BitmarkPayServiceBin: defaultBitmarkPayServiceBin,
		BitmarkConsoleBin:    defaultBitmarkConsoleBin,
		Logging:              *defaultLogger,
	}

	if "" != dataDirectory {
		config.DataDirectory = dataDirectory
	}

	if err := setLoggerPath(config.DataDirectory, &config.Logging); nil != err {
		return nil, err
	}

	return &config, nil
}

func GetConfiguration(configurationFileName string) (*Configuration, error) {

	configurationFileName, err := filepath.Abs(filepath.Clean(configurationFileName))
	if nil != err {
		return nil, err
	}

	options := &Configuration{
		Logging: *defaultLogger,
	}

	if err := readConfigurationFile(configurationFileName, options); err != nil {
		return nil, err
	}

	setLoggerPath(options.DataDirectory, &options.Logging)
	return options, nil
}

func setLoggerPath(baseDir string, logging *LoggerType) error {
	// force all relevant items to be absolute paths
	// if not, assign them to the dsts directory
	mustBeAbsolute := []*string{
		&logging.Directory,
	}

	for _, f := range mustBeAbsolute {
		*f = ensureAbsolute(baseDir, *f)
	}

	// fail if any of these are not simple file names i.e. must not contain path seperator
	// then add the correct directory prefix, file item is first and corresponding directory is second
	mustNotBePaths := [][2]*string{
		{&logging.File, &logging.Directory},
	}
	for _, f := range mustNotBePaths {
		switch filepath.Dir(*f[0]) {
		case "", ".":
			*f[0] = ensureAbsolute(*f[1], *f[0])
		default:
			return errors.New(fmt.Sprintf("Files: %q is not plain name", *f[0]))
		}
	}

	// make absolute and create directories if they do not already exist
	for _, d := range mustBeAbsolute {
		if err := os.MkdirAll(*d, 0700); nil != err {
			return err
		}
	}

	return nil

}

// ensure the path is absolute
func ensureAbsolute(directory string, filePath string) string {
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(directory, filePath)
	}
	return filepath.Clean(filePath)
}
