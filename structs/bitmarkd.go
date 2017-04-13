package structs

import (
	"github.com/bitmark-inc/bitmarkd/payment/bitcoin"
	"github.com/bitmark-inc/bitmarkd/peer"
	"github.com/bitmark-inc/bitmarkd/proof"
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
