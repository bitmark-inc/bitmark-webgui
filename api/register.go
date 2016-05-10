// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"github.com/bitmark-inc/bitmark-webgui/services"
)

var bitmarkService *services.Bitmarkd
var bitmarkPayService *services.BitmarkPay
var bitmarkCliService *services.BitmarkCli

func Register(service interface{}) {
	switch service.(type) {
	case *services.Bitmarkd:
		bitmarkService = service.(*services.Bitmarkd)
	case *services.BitmarkPay:
		bitmarkPayService = service.(*services.BitmarkPay)
	case *services.BitmarkCli:
		bitmarkCliService = service.(*services.BitmarkCli)
	}
}
