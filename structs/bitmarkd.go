package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitmark-inc/bitmarkd/configuration"
	"github.com/bitmark-inc/bitmarkd/payment/bitcoin"
	"github.com/bitmark-inc/bitmarkd/peer"
	"github.com/bitmark-inc/bitmarkd/proof"
	"github.com/bitmark-inc/bitmarkd/util"
	"os"
	"path/filepath"
	"strings"
)

// basic defaults (directories and files are relative to the "DataDirectory" from Configuration file)
const (
	defaultPeerPublicKeyFile   = "peer.private"
	defaultPeerPrivateKeyFile  = "peer.public"
	defaultProofPublicKeyFile  = "proof.private"
	defaultProofPrivateKeyFile = "proof.public"
	defaultProofSigningKeyFile = "proof.sign"
	defaultKeyFile             = "rpc.key"
	defaultCertificateFile     = "rpc.crt"

	defaultLevelDBDirectory = "data"
	defaultBitmarkDatabase  = Bitmark + ".leveldb"
	defaultTestingDatabase  = Testing + ".leveldb"
	defaultLocalDatabase    = Local + ".leveldb"

	defaultRPCClients = 10
	defaultPeers      = 125
	defaultMines      = 125

	defaultBitmarkdLogFile = "bitmarkd.log"
)

type RPCType struct {
	MaximumConnections int      `libucl:"maximum_connections"`
	Listen             []string `libucl:"listen"`
	Certificate        string   `libucl:"certificate"`
	PrivateKey         string   `libucl:"private_key"`
	Announce           []string `libucl:"announce"`
}

type DatabaseType struct {
	Directory string `libucl:"directory"`
	Name      string `libucl:"name"`
}

type BitmarkdConfiguration struct {
	DataDirectory string       `libucl:"data_directory"`
	PidFile       string       `libucl:"pidfile"`
	Chain         string       `libucl:"chain"`
	Nodes         string       `libucl:"nodes"`
	Database      DatabaseType `libucl:"database"`

	ClientRPC RPCType               `libucl:"client_rpc"`
	Peering   peer.Configuration    `libucl:"peering"`
	Proofing  proof.Configuration   `libucl:"proofing"`
	Bitcoin   bitcoin.Configuration `libucl:"bitcoin"`
	Logging   LoggerType            `libucl:"logging"`
}

func (b *BitmarkdConfiguration) SaveToJson(file *os.File) error {
	b.Database.Name = filepath.Base(b.Database.Name)
	b.Logging.File = filepath.Base(b.Logging.File)
	encoder := json.NewEncoder(file)
	return encoder.Encode(b)
}

// will read decode and verify the configuration
func NewBitmarkdConfiguration(configurationFileName string) (*BitmarkdConfiguration, error) {

	configurationFileName, err := filepath.Abs(filepath.Clean(configurationFileName))
	if nil != err {
		return nil, err
	}

	// absolute path to the main directory
	dataDirectory, _ := filepath.Split(configurationFileName)

	options := &BitmarkdConfiguration{

		DataDirectory: dataDirectory,
		PidFile:       "", // no PidFile by default
		Chain:         Bitmark,

		Database: DatabaseType{
			Directory: defaultLevelDBDirectory,
			Name:      defaultBitmarkDatabase,
		},

		ClientRPC: RPCType{
			MaximumConnections: defaultRPCClients,
			Certificate:        defaultCertificateFile,
			PrivateKey:         defaultKeyFile,
		},

		Peering: peer.Configuration{
			//MaximumConnections: defaultPeers,
			PublicKey:  defaultPeerPublicKeyFile,
			PrivateKey: defaultPeerPrivateKeyFile,
			Announce: peer.Announce{
				Broadcast: []string{},
				Listen:    []string{},
			},
		},

		Proofing: proof.Configuration{
			//MaximumConnections: defaultProofers,
			PublicKey:  defaultProofPublicKeyFile,
			PrivateKey: defaultProofPrivateKeyFile,
			SigningKey: defaultProofSigningKeyFile,
		},

		Logging: LoggerType{
			Directory: defaultLogDirectory,
			File:      defaultBitmarkdLogFile,
			Size:      defaultLogSize,
			Count:     defaultLogCount,
			Levels:    defaultLogLevels,
		},
	}

	if err := configuration.ParseConfigurationFile(configurationFileName, options); err != nil {
		return nil, err
	}

	// if any test mode and the database file was not specified
	// switch to appropriate default.  Abort if then chain name is
	// not recognised.
	options.Chain = strings.ToLower(options.Chain)
	if !IsValidChain(options.Chain) {
		return nil, errors.New(fmt.Sprintf("Chain: %q is not supported", options.Chain))
	}

	// if database was not changed from default
	if options.Database.Name == defaultBitmarkDatabase {
		switch options.Chain {
		case Bitmark:
			// already correct default
		case Testing:
			options.Database.Name = defaultTestingDatabase
		case Local:
			options.Database.Name = defaultLocalDatabase
		default:
			return nil, errors.New(fmt.Sprintf("Chain: %s no default database setting", options.Chain))
		}
	}

	// ensure absolute data directory
	if "" == options.DataDirectory || "~" == options.DataDirectory {
		return nil, errors.New(fmt.Sprintf("Path: %q is not a valid directory", options.DataDirectory))
	} else if "." == options.DataDirectory {
		options.DataDirectory = dataDirectory // same directory as the configuration file
	} else {
		options.DataDirectory = filepath.Clean(options.DataDirectory)
	}

	// this directory must exist - i.e. must be created prior to running
	if fileInfo, err := os.Stat(options.DataDirectory); nil != err {
		return nil, err
	} else if !fileInfo.IsDir() {
		return nil, errors.New(fmt.Sprintf("Path: %q is not a directory", options.DataDirectory))
	}

	// force all relevant items to be absolute paths
	// if not, assign them to the data directory
	mustBeAbsolute := []*string{
		&options.Database.Directory,
		&options.ClientRPC.Certificate,
		&options.ClientRPC.PrivateKey,
		&options.Peering.PublicKey,
		&options.Peering.PrivateKey,
		&options.Proofing.PublicKey,
		&options.Proofing.PrivateKey,
		&options.Proofing.SigningKey,
		&options.Logging.Directory,
	}
	for _, f := range mustBeAbsolute {
		*f = util.EnsureAbsolute(options.DataDirectory, *f)
	}

	// optional absolute paths i.e. blank or an absolute path
	optionalAbsolute := []*string{
		&options.PidFile,
		&options.Bitcoin.CACertificate,
		&options.Bitcoin.Certificate,
		&options.Bitcoin.PrivateKey,
	}
	for _, f := range optionalAbsolute {
		if "" != *f {
			*f = util.EnsureAbsolute(options.DataDirectory, *f)
		}
	}

	// fail if any of these are not simple file names i.e. must not contain path seperator
	// then add the correct directory prefix, file item is first and corresponding directory is second
	mustNotBePaths := [][2]*string{
		{&options.Database.Name, &options.Database.Directory},
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
		&options.Database.Directory,
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
