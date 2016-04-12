// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package main

import (
	"github.com/bitmark-inc/bitmark-mgmt/api"
	"github.com/bitmark-inc/bitmark-mgmt/services"
	"github.com/bitmark-inc/bitmarkd/background"
)

var backgroundService *background.T
var bitmarkService services.Bitmarkd

// start service
func InitialiseBackgroundService(configFile string) error {

	// initialise all services
	if err := bitmarkService.Initialise(configFile); nil != err {
		return err
	}

	// create and start all background service
	var processes = background.Processes{
		bitmarkService.BitmarkdBackground,
	}
	backgroundService = background.Start(processes, nil)

	// register services to api
	api.RegisterBitmarkd(&bitmarkService)

	return nil
}

// finialise - stop all background tasks
func FinaliseBackgroundService() error {

	if err := bitmarkService.Finalise(); nil != err {
		return err
	}

	// stop backgrond services
	background.Stop(backgroundService)
	return nil
}
