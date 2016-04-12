// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"github.com/bitmark-inc/bitmark-mgmt/services"
)

var bitmarkService *services.Bitmarkd

func RegisterBitmarkd(bitmarkd *services.Bitmarkd) {
	bitmarkService = bitmarkd
}