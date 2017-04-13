package structs

// names of all chains
const (
	Bitmark = "bitmark"
	Testing = "testing"
	Local   = "local"
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
