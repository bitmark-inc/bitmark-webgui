// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package templates

const (
	/**** Configuration template ****/
	ConfigurationTemplate = `
# bitmark-mgmt.conf -*- mode: libucl -*-

port = "{{.Port}}"
password = "{{.Password}}"
enable_https = "{{.EnableHttps}}"
bitmark_config_file = "{{.BitmarkConfigFile}}"
`

	BitmarkConfigTemplate = `
  {{.Field}} = {{.Value}}`

	BitmarkConnectTemplate = `
  connect = {public_key = "{{.PublicKey}}", address = "{{.Address}}"}`
)
