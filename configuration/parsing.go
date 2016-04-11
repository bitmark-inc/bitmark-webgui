// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package configuration

import (
	"errors"
	"fmt"
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"os"
	"path/filepath"
)

const (
	defaultLogDirectory = "log"
	defaultLogFile      = "bitmark-mgmt.log"
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
	Port              int        `libucl:"port"`
	Password          string     `libucl:"password"`
	EnableHttps       bool       `libucl:"enable_https"`
	BitmarkConfigFile string     `libucl:"bitmark_config_file"`
	Logging           LoggerType `libucl:"logging"`
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

func GetDefaultLogger(baseDir string) *LoggerType {
	setLoggerPath(baseDir, defaultLogger)
	return defaultLogger
}

func GetConfiguration(baseDir string, configurationFileName string) (*Configuration, error) {

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

	setLoggerPath(baseDir, &options.Logging)
	// // force all relevant items to be absolute paths
	// // if not, assign them to the dsts directory
	// mustBeAbsolute := []*string{
	// 	&options.Logging.Directory,
	// }

	// for _, f := range mustBeAbsolute {
	// 	*f = ensureAbsolute(baseDir, *f)
	// }

	// // fail if any of these are not simple file names i.e. must not contain path seperator
	// // then add the correct directory prefix, file item is first and corresponding directory is second
	// mustNotBePaths := [][2]*string{
	// 	{&options.Logging.File, &options.Logging.Directory},
	// }
	// for _, f := range mustNotBePaths {
	// 	switch filepath.Dir(*f[0]) {
	// 	case "", ".":
	// 		*f[0] = ensureAbsolute(*f[1], *f[0])
	// 	default:
	// 		return nil, errors.New(fmt.Sprintf("Files: %q is not plain name", *f[0]))
	// 	}
	// }

	// // make absolute and create directories if they do not already exist
	// for _, d := range []*string{&options.Logging.Directory} {
	// 	*d = ensureAbsolute(baseDir, *d)
	// 	if err := os.MkdirAll(*d, 0700); nil != err {
	// 		return nil, err
	// 	}
	// }

	// done
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
	for _, d := range []*string{&logging.Directory} {
		*d = ensureAbsolute(baseDir, *d)
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
