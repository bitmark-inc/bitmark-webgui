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
var bitmarkService services.Bitmarkd
var bitmarkPayService services.BitmarkPay

// start service
func InitialiseService(configs *configuration.Configuration) error {

	// initialise all  services
	if err := bitmarkService.Initialise(configs.BitmarkConfigFile); nil != err {
		return err
	}
	if err := bitmarkPayService.Initialise(configs.BitmarkPayServiceBin); nil != err {
		return err
	}

	// create and start all background service
	var processes = background.Processes{
		bitmarkService.BitmarkdBackground,
	}
	backgroundService = background.Start(processes, nil)

	// register services to api
	api.Register(&bitmarkService)
	api.Register(&bitmarkPayService)

	return nil
}

// finialise - stop all background tasks
func FinaliseBackgroundService() error {

	if err := bitmarkService.Finalise(); nil != err {
		return err
	}

	if err := bitmarkPayService.Finalise(); nil != err {
		return err
	}

	// stop backgrond services
	background.Stop(backgroundService)
	return nil
}
