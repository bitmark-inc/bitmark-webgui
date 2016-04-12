// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file

package main

import(
	"github.com/bitmark-inc/bitmarkd/background"
	"github.com/bitmark-inc/bitmark-mgmt/api"
	"github.com/bitmark-inc/bitmark-mgmt/services"
	"github.com/bitmark-inc/logger"
)

type service struct {

	background *background.T

	log *logger.L

 	service interface{}
}

var bitmarkService services.Bitmarkd
var globalService service

// start service
func Initialise(configFile string) error {

	if err := bitmarkService.Initialise(configFile); nil != err {
		return err
	}

	globalService.service = bitmarkService
	globalService.log = logger.New("service")

	var processes = background.Processes{
		bitmarkService.BitmarkdBackground,
	}
	globalService.background = background.Start(processes, nil)

	api.RegisterBitmarkd(&bitmarkService)

	return nil
}

// finialise - stop all background tasks
func Finalise() error {

	if err := bitmarkService.Finalise(); nil != err {
		return err
	}

	background.Stop(globalService.background)
	return nil
}
