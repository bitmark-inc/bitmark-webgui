// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/bitmark-inc/bitmark-mgmt/fault"
	"github.com/bitmark-inc/certgen"
	"io/ioutil"
	"os"
	"time"
)

const (
	ServerTlsDir   = "./tls"
	ServerCertFile = "bitmark-mgmt.crt"
	ServerKeyFile  = "bitmark-mgmt.key"
)

// config is required
func CheckConfigDir(path string) (string, error) {
	if "" == path {
		return "", fault.ErrRequiredConfigDir
	}

	path = os.ExpandEnv(path)
	return path, nil
}

// check if file exists
func EnsureFileExists(name string) bool {
	_, err := os.Stat(name)
	return nil == err
}

func GetTLSCertFile() (string, string, bool, error) {
	newCreate := false

	// Set up folder
	if !EnsureFileExists(ServerTlsDir) {
		if err := os.MkdirAll(ServerTlsDir, 0755); nil != err {
			return "", "", newCreate, err
		}
		newCreate = true
	}

	//Set up certification files
	certFile := ServerTlsDir + "/" + ServerCertFile
	if !EnsureFileExists(certFile) {
		newCreate = true
	}

	keyFile := ServerTlsDir + "/" + ServerKeyFile
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
