package structs

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
	Directory string            `libucl:"directory"`
	File      string            `libucl:"file"`
	Size      int               `libucl:"size"`
	Count     int               `libucl:"count"`
	Levels    map[string]string `libucl:"levels"`
}
