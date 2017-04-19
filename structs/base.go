package structs

import (
	"github.com/bitmark-inc/logger"
)

// to hold log levels
type LoglevelMap map[string]string

// names of all chains
const (
	Bitmark = "bitmark"
	Testing = "testing"
	Local   = "local"

	defaultDataDirectory = "" // this will error; use "." for the same directory as the config file

	defaultLogDirectory = "log"
	defaultLogCount     = 10          //  number of log files retained
	defaultLogSize      = 1024 * 1024 // rotate when <logfile> exceeds this size
)

var (
	defaultLogLevels = LoglevelMap{
		"main":            "info",
		logger.DefaultTag: "info",
	}
)

// validate a chain name
func IsValidChain(name string) bool {
	switch name {
	case Bitmark, Testing, Local:
		return true
	default:
		return false
	}
}

type LoggerType struct {
	Directory string            `libucl:"directory" json:"directory"`
	File      string            `libucl:"file" json:"file"`
	Size      int               `libucl:"size" json:"size"`
	Count     int               `libucl:"count" json:"count"`
	Levels    map[string]string `libucl:"levels" json:"levels"`
}
