// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package templates

const (
	/**** Configuration template ****/
	ConfigurationTemplate = `
# bitmark-webgui.conf -*- mode: libucl -*-

data_directory = "{{.DataDirectory}}"

port = "{{.Port}}"
password = "{{.Password}}"
enable_https = "{{.EnableHttps}}"
bitmark_config_file = "{{.BitmarkConfigFile}}"

logging {
  size = 1048576
  count = 10

  levels {
    "*" = info
    main = info
    api = info
  }
}
`

	BitmarkConfigTemplate = `
  {{.Field}} = {{.Value}}`

	BitmarkConnectTemplate = `
  connect = {public_key = "{{.PublicKey}}", address = "{{.Address}}"}`

	BitmarkGeneralTemplate = `
{{.Field}} = {{.Value}}
`
)
