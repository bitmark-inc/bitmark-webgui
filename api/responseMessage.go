// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

const (
	// general error
	invalidValueErr = `Invalid value`

	// bitmarkConfig error
	bitmarkdConfigGetErr    = `Failed to get bitmarkd config`
	bitmarkdConfigUpdateErr = `Failed to update bitmarkd config`

	// bitmarkMgmt error
	setPasswordErr = `Failed to set up bitmark-mgmt password`

	// bitmarkd api
	bitmarkdStartSuccess = `start running bitmarkd`
	bitmarkdStopSuccess  = `stop running bitmarkd`
	bitmarkdStarted      = `started`
	bitmarkdStopped      = `stopped`

	// bitmarkd error
	bitmarkdStartErr        = `Failed to start bitmarkd`
	bitmarkdStopErr         = `Failed to stop bitmarkd`
	bitmarkdAlreadyStartErr = `Already started bitmarkd`
	bitmarkdAlreadyStopErr  = `Already stoped bitmarkd`
	bitmarkdConnectErr      = `Failed to connect to bitmarkd`
	bitmarkdGetInfoErr      = `Failed to get bitmarkd info`
	bitmarkdGetConfigErr    = `Failed to get bitamrkd configuration`
	// login
	loginErr = `Failed to log in`
)
