// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"github.com/bitmark-inc/bitmark-webgui/services"
)

var bitmarkService *services.Bitmarkd
var bitmarkPayService services.BitmarkPayInterface
var bitmarkCliService *services.BitmarkCli

func Register(service interface{}) {
	switch t := service.(type) {
	case *services.Bitmarkd:
		bitmarkService = service.(*services.Bitmarkd)
	case services.BitmarkPayInterface:
		bitmarkPayService = service.(services.BitmarkPayInterface)
	case *services.BitmarkCli:
		bitmarkCliService = service.(*services.BitmarkCli)
	default:
		fmt.Printf("Undefined type: %v\n", t)
	}
}
