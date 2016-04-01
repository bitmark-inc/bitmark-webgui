// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// error instances
//
// Provides a single instance of errors to allow easy comparison
package fault

// error base
type GenericError string

// to allow for different classes of errors
type ExistsError GenericError
type InvalidError GenericError
type NotFoundError GenericError
type ProcessError GenericError

// common errors - keep in alphabetic order
var (
	ErrCertificateFileAlreadyExists = ExistsError("certificate file already exists")
	ErrKeyFileAlreadyExists         = ExistsError("key file already exists")
	ErrPasswordLength               = InvalidError("Password Length is invalid")
	ErrVerifiedPassword             = InvalidError("Verified password is different")
	ErrRequiredConfigDir            = InvalidError("Config folder is required")
	ErrConfigDirPath                = InvalidError("Config is not a folder")
	ErrWrongPassword                = InvalidError("Wrong password")
	ErrNotFoundPublicKey            = NotFoundError("PublicKey is not existed")
	ErrNotFoundConfigFile           = NotFoundError("Config file is not found")
	ErrJsonParseFail                = ProcessError("Parse to json failed")
	ErrInvalidStructPointer         = InvalidError("invalid struct pointer")

	// For API response
	ApiErrAlreadyLoggedIn      = InvalidError("Already logged in")
	ApiErrUnauthorized         = InvalidError("API Unauthorized")
	ApiErrChekAuthorize        = InvalidError("Failed to check authorization")
	ApiErrSetAuthorize         = InvalidError("Failed to authorize")
	ApiErrStatusBitmarkd       = InvalidError("Failed to get bitmarkd status")
	ApiErrStartBitmarkd        = InvalidError("Failed to start bitmarkd")
	ApiErrStopBitmarkd         = InvalidError("Failed to stop bitmarkd")
	ApiErrAlreadyStartBitmarkd = InvalidError("Already started bitmarkd")
	ApiErrInvalidValue         = InvalidError("Invalid value")
	ApiErrGetBitmarkConfig     = ProcessError("Failed to get bitmarkd config")
	ApiErrUpdateBitmarkdConfig = ProcessError("Failed to update bitmarkd config")
	ApiErrSetPassword          = ProcessError("Failed to set up bitmark-mgmt password")
	ApiErrLogin                = ProcessError("Failed to log in")
)

// the error interface base method
func (e GenericError) Error() string { return string(e) }

// the error interface methods
func (e ExistsError) Error() string   { return string(e) }
func (e InvalidError) Error() string  { return string(e) }
func (e NotFoundError) Error() string { return string(e) }
func (e ProcessError) Error() string  { return string(e) }

// determine the class of an error
func IsErrExists(e error) bool   { _, ok := e.(ExistsError); return ok }
func IsErrInvalid(e error) bool  { _, ok := e.(InvalidError); return ok }
func IsErrNotFound(e error) bool { _, ok := e.(NotFoundError); return ok }
func IsErrProcess(e error) bool  { _, ok := e.(ProcessError); return ok }
