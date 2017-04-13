// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/structs"
	"github.com/bitmark-inc/bitmark-webgui/templates"
	"github.com/bitmark-inc/bitmark-webgui/utils"
	"github.com/bitmark-inc/bitmarkd/configuration"
	"github.com/bitmark-inc/bitmarkd/peer"
	"github.com/bitmark-inc/logger"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

// Get /api/config
func ListConfig(w http.ResponseWriter, req *http.Request, bitmarkConfigFile string, log *logger.L) {
	log.Info("GET /api/config")
	response := &Response{
		Ok:     false,
		Result: bitmarkdConfigGetErr,
	}

	bitmarkConfigs := &structs.BitmarkdConfiguration{}
	if err := configuration.ParseConfigurationFile(bitmarkConfigFile, bitmarkConfigs); nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
	} else {
		bitmarkConfDir := filepath.Dir(bitmarkConfigFile)
		pubKeyFile := filepath.Join(bitmarkConfDir, bitmarkConfigs.Peering.PublicKey)
		if peerPublicKey, err := getPeerPublicKey(pubKeyFile); nil != err {
			response.Result = err
		} else {
			response.Ok = true
			bitmarkConfigs.Peering.PublicKey = *peerPublicKey
			response.Result = bitmarkConfigs
		}
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

// Post /api/config
func UpdateConfig(w http.ResponseWriter, req *http.Request, bitmarkConfigFile string, log *logger.L) {

	log.Info("POST /api/config")
	response := &Response{
		Ok:     false,
		Result: bitmarkdConfigUpdateErr,
	}

	decoder := json.NewDecoder(req.Body)
	var request structs.BitmarkdConfiguration
	if err := decoder.Decode(&request); nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// Prepare bitmarkConfig
	linesPtr, err := prepareBitmarkConfig(request, bitmarkConfigFile)
	if nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// remove redundant lines
	lines := *linesPtr
	outputLines := make([]string, 0)
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			outputLines = append(outputLines, lines[i])
		} else {
			if i+1 < len(lines) && lines[i+1] != "" {
				outputLines = append(outputLines, lines[i])
			}
		}
	}

	// Write result to bitmark config file
	output := strings.Join(outputLines, "\n")
	err = ioutil.WriteFile(bitmarkConfigFile, []byte(output), 0644)
	if nil != err {
		log.Errorf("Error: %s, %v", bitmarkdConfigUpdateErr, err)
		response.Result = bitmarkdConfigUpdateErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	response.Ok = true
	response.Result = nil
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

func getPeerPublicKey(filePath string) (*string, error) {
	if !utils.EnsureFileExists(filePath) {
		return nil, fault.ErrNotFoundPublicKey
	}
	publicKey, err := ioutil.ReadFile(filePath)
	if nil != err {
		return nil, err
	}

	publicKeyStr := string(publicKey)
	return &publicKeyStr, nil
}

// read existed bitmark config file to string, and set it
func prepareBitmarkConfig(request structs.BitmarkdConfiguration, bitmarkConfigFile string) (*[]string, error) {

	input, err := ioutil.ReadFile(bitmarkConfigFile)
	if nil != err {
		return nil, err
	}

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Index(line, "chain") == 0 {
			item := []string{request.Chain}
			ifItem := prepareBitmarkField("chain", item)
			if err := updateConfigString(lines, i, "chain", ifItem); nil != err {
				return nil, err
			}
		} else if strings.Index(line, "client_rpc") == 0 {
			// Set listen
			listens := prepareBitmarkField("listen", request.ClientRPC.Listen)
			if err := updateConfigString(lines, i, "listen", listens); nil != err {
				return nil, err
			}
			// Set announce
			announces := prepareBitmarkField("announce", request.ClientRPC.Announce)
			if err := updateConfigString(lines, i, "announce", announces); nil != err {
				return nil, err
			}
		} else if strings.Index(line, "peering") == 0 {
			// Set listen
			listens := prepareBitmarkField("listen", request.Peering.Listen)
			if err := updateConfigString(lines, i, "listen", listens); nil != err {
				return nil, err
			}
			// // Set announce
			// announces := prepareBitmarkField("announce", []string{"127.0.0.1:2130"})
			// if err := updateConfigString(lines, i, "announce", announces); nil != err {
			// 	return nil, err
			// }
			// Set connect
			connections := prepareConnectField(request.Peering.Connect)
			if err := updateConfigString(lines, i, "connect", connections); nil != err {
				return nil, err
			}
		} else if strings.Index(line, "bitcoin") == 0 {
			//set  username
			item := []string{request.Bitcoin.Username}
			ifItem := prepareBitmarkField("username", item)
			if err := updateConfigString(lines, i, "username", ifItem); nil != err {
				return nil, err
			}
			// password
			if request.Bitcoin.Password != "" {
				item = []string{request.Bitcoin.Password}
				ifItem = prepareBitmarkField("password", item)
				if err := updateConfigString(lines, i, "password", ifItem); nil != err {
					return nil, err
				}
			}
			// url - use spoon api url
			item = []string{request.Bitcoin.URL}
			ifItem = prepareBitmarkField("url", item)
			if err := updateConfigString(lines, i, "url", ifItem); nil != err {
				return nil, err
			}
		}
	}
	return &lines, nil
}

type bitmarkStringArrayType struct {
	Field string
	Value string
}

func prepareBitmarkField(field string, source []string) []interface{} {

	var localSrc = make([]string, 0)
	for _, src := range source {
		if src != "" {
			localSrc = append(localSrc, src)
		}
	}

	bitmarkFields := make([]bitmarkStringArrayType, len(localSrc))
	for i, src := range localSrc {
		bitmarkFields[i].Field = field
		if field == "listen" || field == "announce" {
			if strings.Contains(src, "[") {
				src = `"` + src + `"`
			}
		} else if field == "fee" {
			src = `"` + src + `"`
		}
		bitmarkFields[i].Value = src
	}

	interfaceBitmarkFields := make([]interface{}, len(bitmarkFields))
	for i, bf := range bitmarkFields {
		interfaceBitmarkFields[i] = bf
	}

	return interfaceBitmarkFields
}

func prepareConnectField(source []peer.Connection) []interface{} {
	localSrc := make([]interface{}, 0)
	for _, src := range source {
		if src.PublicKey != "" && src.Address != "" {
			localSrc = append(localSrc, src)
		}
	}
	return localSrc
}

func updateConfigString(lines []string, index int, field string, values []interface{}) error {
	templateStr := templates.BitmarkConfigTemplate
	if field == "connect" {
		templateStr = templates.BitmarkConnectTemplate
	} else if field == "chain" {
		templateStr = templates.BitmarkGeneralTemplate
	}

	// Prepare update string
	fieldTemp := template.Must(template.New("field").Parse(templateStr))
	fieldBuffer := new(bytes.Buffer)
	for _, v := range values {
		err := fieldTemp.Execute(fieldBuffer, v)
		if nil != err {
			return err
		}
	}

	writePoint := 0
	fieldNotFound := false
	// empty all existed field
	for i := index; i < len(lines); i++ {
		line := lines[i]
		line = strings.TrimSpace(line)
		if strings.Contains(line, field) {
			if strings.Index(line, "#") == 0 { // leave comment
				continue
			} else {
				if writePoint == 0 {
					writePoint = i
				}
				lines[i] = ""
			}
		} else if strings.Index(line, "}") == 0 { // item block finish
			if writePoint == 0 { // no field setting
				fieldNotFound = true
				writePoint = i
			}
			break
		}
	}

	if fieldNotFound {
		lines[writePoint] = fieldBuffer.String() + "\n}"
	} else {
		lines[writePoint] = fieldBuffer.String()
	}

	return nil
}
