// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/certgen"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

const (
	ServerTlsDir   = "tls"
	ServerCertFile = "bitmark-mgmt.crt"
	ServerKeyFile  = "bitmark-mgmt.key"
)

// check if file exists
func EnsureFileExists(name string) bool {
	_, err := os.Stat(name)
	return nil == err
}

func GetTLSCertFile(baseDir string) (string, string, bool, error) {
	newCreate := false
	serverTlsDir := baseDir + "/" + ServerTlsDir

	// Set up folder
	if !EnsureFileExists(serverTlsDir) {
		if err := os.MkdirAll(serverTlsDir, 0755); nil != err {
			return "", "", newCreate, err
		}
		newCreate = true
	}

	//Set up certification files
	certFile := serverTlsDir + "/" + ServerCertFile
	if !EnsureFileExists(certFile) {
		newCreate = true
	}

	keyFile := serverTlsDir + "/" + ServerKeyFile
	if !EnsureFileExists(keyFile) {
		newCreate = true
	}

	return certFile, keyFile, newCreate, nil

}

// create a self-signed certificate
func MakeSelfSignedCertificate(name string, certificateFileName string, privateKeyFileName string, override bool, extraHosts []string) error {
	if EnsureFileExists(certificateFileName) {
		return fault.ErrCertificateFileAlreadyExists
	}

	if EnsureFileExists(privateKeyFileName) {
		return fault.ErrKeyFileAlreadyExists
	}

	org := "bitmark self signed cert for: " + name
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, err := certgen.NewTLSCertPair(org, validUntil, override, extraHosts)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(certificateFileName, cert, 0666); err != nil {
		return err
	}

	if err = ioutil.WriteFile(privateKeyFileName, key, 0600); err != nil {
		os.Remove(certificateFileName)
		return err
	}

	return nil
}

var validBitcoinAddress = regexp.MustCompile(`^([a-z]*[A-Z]*[0-9]*)+$`)

func CheckBitcoinAddress(address string) error {
	addrLen := len(address)
	if addrLen < 26 || addrLen > 35 || !validBitcoinAddress.MatchString(address) {
		return fault.ErrBitcoinAddress
	}
	return nil
}
