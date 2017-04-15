package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitmark-inc/bitmarkd/configuration"
	"github.com/bitmark-inc/bitmarkd/util"
	"github.com/bitmark-inc/logger"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultPublicKeyFile  = "prooferd.private"
	defaultPrivateKeyFile = "prooferd.public"

	defaultProoferdLogFile = "prooferd.log"
)

// to hold log levels
type LoglevelMap map[string]string

// path expanded or calculated defaults
var (
	defaultLogLevels = LoglevelMap{
		"main":            "info",
		logger.DefaultTag: "critical",
	}
)

type Connection struct {
	PublicKey string `libucl:"public_key" json:"public_key"`
	Blocks    string `libucl:"blocks" json:"blocks"`
	Submit    string `libucl:"submit" json:"submit"`
}

type PeerType struct {
	PrivateKey string       `libucl:"private_key" json:"private_key"`
	PublicKey  string       `libucl:"public_key" json:"public_key"`
	Connect    []Connection `libucl:"connect" json:"connect"`
}

type ProoferdConfiguration struct {
	DataDirectory string     `libucl:"data_directory" json:"data_directory"`
	PidFile       string     `libucl:"pidfile" json:"pidfile"`
	Chain         string     `libucl:"chain" json:"chain"`
	Threads       int        `libucl:"threads" json:"threads"`
	Peering       PeerType   `libucl:"peering" json:"peering"`
	Logging       LoggerType `libucl:"logging" json:"logging"`
}

func (p *ProoferdConfiguration) SaveToJson(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	p.Logging.File = filepath.Base(p.Logging.File)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(p)
}

// will read decode and verify the configuration
func NewProoferdConfiguration(configurationFileName string) (*ProoferdConfiguration, error) {

	configurationFileName, err := filepath.Abs(filepath.Clean(configurationFileName))
	if nil != err {
		return nil, err
	}

	// absolute path to the main directory
	dataDirectory, _ := filepath.Split(configurationFileName)

	options := &ProoferdConfiguration{

		DataDirectory: dataDirectory,
		PidFile:       "", // no PidFile by default
		Chain:         Bitmark,
		Threads:       0,

		Peering: PeerType{
			Connect:    make([]Connection, 0),
			PrivateKey: defaultProofPrivateKeyFile,
			PublicKey:  defaultProofPublicKeyFile,
		},

		Logging: LoggerType{
			Directory: defaultLogDirectory,
			File:      defaultProoferdLogFile,
			Size:      defaultLogSize,
			Count:     defaultLogCount,
			Levels:    defaultLogLevels,
		},
	}

	_ = configuration.ParseConfigurationFile(configurationFileName, options)

	// if any test mode and the database file was not specified
	// switch to appropriate default.  Abort if then chain name is
	// not recognised.
	options.Chain = strings.ToLower(options.Chain)
	if !IsValidChain(options.Chain) {
		return nil, errors.New(fmt.Sprintf("Chain: %q is not supported", options.Chain))
	}

	// if threads invalid set number of CPUs
	if options.Threads <= 0 {
		options.Threads = runtime.NumCPU()
	}

	// ensure absolute data directory
	if "" == options.DataDirectory || "~" == options.DataDirectory {
		return nil, errors.New(fmt.Sprintf("Path: %q is not a valid directory", options.DataDirectory))
	} else if "." == options.DataDirectory {
		options.DataDirectory = dataDirectory // same directory as the configuration file
	}
	options.DataDirectory = filepath.Clean(options.DataDirectory)

	// this directory must exist - i.e. must be created prior to running
	if fileInfo, err := os.Stat(options.DataDirectory); nil != err {
		return nil, err
	} else if !fileInfo.IsDir() {
		return nil, errors.New(fmt.Sprintf("Path: %q is not a directory", options.DataDirectory))
	}

	// force all relevant items to be absolute paths
	// if not, assign them to the data directory
	mustBeAbsolute := []*string{
		&options.Peering.PublicKey,
		&options.Peering.PrivateKey,
		&options.Logging.Directory,
	}
	for _, f := range mustBeAbsolute {
		*f = util.EnsureAbsolute(options.DataDirectory, *f)
	}

	// optional absolute paths i.e. blank or an absolute path
	optionalAbsolute := []*string{
		&options.PidFile,
	}
	for _, f := range optionalAbsolute {
		if "" != *f {
			*f = util.EnsureAbsolute(options.DataDirectory, *f)
		}
	}

	// fail if any of these are not simple file names i.e. must not contain path seperator
	// then add the correct directory prefix, file item is first and corresponding directory is second
	mustNotBePaths := [][2]*string{
		{&options.Logging.File, &options.Logging.Directory},
	}
	for _, f := range mustNotBePaths {
		switch filepath.Dir(*f[0]) {
		case "", ".":
			*f[0] = util.EnsureAbsolute(*f[1], *f[0])
		default:
			return nil, errors.New(fmt.Sprintf("Files: %q is not plain name", *f[0]))
		}
	}

	// make absolute and create directories if they do not already exist
	for _, d := range []*string{
		&options.Logging.Directory,
	} {
		*d = util.EnsureAbsolute(options.DataDirectory, *d)
		if err := os.MkdirAll(*d, 0700); nil != err {
			return nil, err
		}
	}

	// done
	return options, nil
}
