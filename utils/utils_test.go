// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package utils_test

import(
	"github.com/bitmark-inc/bitmark-mgmt/utils"
	"testing"
)

func TestCheckBitcoinAddress(t *testing.T) {
	invalidBitcoinAddressArr := []string{
		"123", //less than 26
		"qwertasdfg123456789olijsdkvae0293j5nfjkso", // greater than 35
		"123456789012345678901234567890@", // contains invalid character
	}

	for _, addr := range invalidBitcoinAddressArr {
		err := 	utils.CheckBitcoinAddress(addr)
		if nil == err {
			t.Errorf("CheckBitcoinAddress pass invalid addr: %s\n", addr)
		}
	}


	validBitcoinAddressArr := []string{
		"123456789012345678901234567890",
		"SD2gsret23SDGAdg46sdf35GDSR2AER3AE",
	}

	for _, addr := range validBitcoinAddressArr {
		err := 	utils.CheckBitcoinAddress(addr)
		if nil != err {
			t.Errorf("CheckBitcoinAddress unpass valid addr: %s\n", addr)
		}
	}
}
