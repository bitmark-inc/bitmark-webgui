// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package configuration

import (
	"os"
	"github.com/bitmark-inc/bitmark-webgui/templates"
	"text/template"
)

func UpdateConfiguration(configFile string, configuration *Configuration) error {
	if file, err := os.Create(configFile); nil != err {
		return err
	} else {
		confTemp := template.Must(template.New("config").Parse(templates.ConfigurationTemplate))
		if err := confTemp.Execute(file, configuration); nil != err {
			return err
		}
	}

	return nil
}
