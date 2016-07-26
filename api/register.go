// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"github.com/bitmark-inc/bitmark-webgui/services"
)

var bitcoinService *services.Bitcoind
var bitmarkService *services.Bitmarkd
var bitmarkPayService services.BitmarkPayInterface
var bitmarkCliService *services.BitmarkCli
var bitmarkConsoleService *services.BitmarkConsole

func Register(service interface{}) {
	switch t := service.(type) {
	case *services.Bitcoind:
		bitcoinService = service.(*services.Bitcoind)
	case *services.Bitmarkd:
		bitmarkService = service.(*services.Bitmarkd)
	case services.BitmarkPayInterface:
		bitmarkPayService = service.(services.BitmarkPayInterface)
	case *services.BitmarkCli:
		bitmarkCliService = service.(*services.BitmarkCli)
	case *services.BitmarkConsole:
		bitmarkConsoleService = service.(*services.BitmarkConsole)
	default:
		fmt.Printf("Undefined type: %v\n", t)
	}
}
