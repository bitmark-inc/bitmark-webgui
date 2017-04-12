// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package main

import (
	"github.com/bitmark-inc/bitmark-webgui/api"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/bitmark-webgui/services"
	"github.com/bitmark-inc/bitmarkd/background"
)

var backgroundService *background.T
var bitcoinService services.Bitcoind
var bitmarkService services.Bitmarkd
var bitmarkConsoleService services.BitmarkConsole

// start service
func InitialiseService(configs *configuration.Configuration) error {

	// initialise all  services
	if err := bitcoinService.Initialise(); nil != err {
		return err
	}
	if err := bitmarkService.Initialise(configs.BitmarkConfigFile); nil != err {
		return err
	}
	if err := bitmarkConsoleService.Initialise(configs.BitmarkConsoleBin); nil != err {
		return err
	}

	// create and start all background service
	var processes = background.Processes{
		&bitcoinService,
		&bitmarkService,
	}
	backgroundService = background.Start(processes, nil)

	// register services to api
	api.Register(&bitcoinService)
	api.Register(&bitmarkService)

	return nil
}

// finialise - stop all background tasks
func FinaliseBackgroundService() error {

	if err := bitcoinService.Finalise(); nil != err {
		return err
	}

	if err := bitmarkService.Finalise(); nil != err {
		return err
	}

	if err := bitmarkConsoleService.Finalise(); nil != err {
		return err
	}

	// stop backgrond services
	backgroundService.Stop()
	return nil
}
