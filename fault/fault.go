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
	ErrAlreadyInitialised               = ExistsError("already initialised")
	ErrBitmarkConsoleIsNotRunning       = ProcessError("Bitmark console is not running")
	ErrBitcoinAddress                   = InvalidError("Invalid bitcoin address")
	ErrBitcoindIsNotRunning             = InvalidError("Bitcoind is not running")
	ErrBitcoindIsRunning                = InvalidError("Bitcoind not running")
	ErrBitmarkdIsNotRunning             = InvalidError("Bitmarkd is not running")
	ErrBitmarkdIsRunning                = InvalidError("Bitmarkd is running")
	ErrBitmarkPayIsRunning              = InvalidError("Bitmark-pay is running")
	ErrCertificateFileAlreadyExists     = ExistsError("certificate file already exists")
	ErrConfigFileExited                 = ExistsError("Config file is existed")
	ErrExecBitmarkPayJob                = ProcessError("Failed to execute bitmark pay job ")
	ErrInvalidAccessBitmarkPayJobResult = InvalidError("invalid access bitmark-pay job result")
	ErrInvalidCommandType               = InvalidError("invalid command type")
	ErrInvalidCommandParams             = InvalidError("invalid command params")
	ErrInvalidLoggerChannel             = InvalidError("invalid logger channel")
	ErrInvalidStructPointer             = InvalidError("invalid struct pointer")
	ErrJsonParseFail                    = ProcessError("Parse to json failed")
	ErrKeyFileAlreadyExists             = ExistsError("key file already exists")
	ErrNodeInfoRequestFail              = ProcessError("Send info request failed")
	ErrNotFoundBitmarkPayJob            = NotFoundError("bitmark-pay job is not found")
	ErrNotFoundBinFile                  = NotFoundError("Bin file is not found")
	ErrNotFoundConfigFile               = NotFoundError("Config file is not found")
	ErrNotFoundPublicKey                = NotFoundError("PublicKey is not existed")
	ErrNotInitialised                   = NotFoundError("not initialised")
	ErrPasswordLength                   = InvalidError("Password Length is invalid")
	ErrVerifiedPassword                 = InvalidError("Verified password is different")
	ErrWrongPassword                    = InvalidError("Wrong password")

	// For API response
	ApiErrAlreadyLoggedIn = InvalidError("Already logged in")
	ApiErrChekAuthorize   = InvalidError("Failed to check authorization")
	ApiErrSetAuthorize    = InvalidError("Failed to authorize")
	ApiErrUnauthorized    = InvalidError("API Unauthorized")
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
