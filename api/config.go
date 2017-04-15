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
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Get /api/config
func ListConfig(w http.ResponseWriter, req *http.Request, bitmarkConfigFile, prooferdConfigFile string, log *logger.L) {
	log.Info("GET /api/config")
	results := map[string]interface{}{}
	var err error
	bitmarkConfigs := &structs.BitmarkdConfiguration{}
	prooferdConfig := &structs.ProoferdConfiguration{}
	if err = configuration.ParseConfigurationFile(bitmarkConfigFile, bitmarkConfigs); nil != err {
		log.Errorf("Error: %v", err)
		results["bitmarkd"] = err
	} else {
		pubKeyFile, err := filepath.Abs(bitmarkConfigs.Peering.PublicKey)
		if err != nil {
			log.Errorf("Error: %v", err)
			results["bitmarkd"] = err
		}
		if peerPublicKey, err := getPeerPublicKey(pubKeyFile); nil != err {
			results["bitmarkd"] = err
		} else {
			bitmarkConfigs.Peering.PublicKey = *peerPublicKey
			results["bitmarkd"] = bitmarkConfigs
		}
	}

	if err := configuration.ParseConfigurationFile(prooferdConfigFile, prooferdConfig); nil != err {
		log.Errorf("Error: %v", err)
		results["prooferd"] = err
	} else {
		pubKeyFile, err := filepath.Abs(prooferdConfig.Peering.PublicKey)
		if err != nil {
			log.Errorf("Error: %v", err)
			results["bitmarkd"] = err
		}
		if peerPublicKey, err := getPeerPublicKey(pubKeyFile); nil != err {
			results["prooferd"] = err
		} else {
			prooferdConfig.Peering.PublicKey = *peerPublicKey
			results["prooferd"] = prooferdConfig
		}
	}

	response := &Response{
		Ok:     true,
		Result: results,
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

// Post /api/config
func UpdateConfig(w http.ResponseWriter, req *http.Request, chain, bitmarkConfigFile, prooferdConfigFile string, log *logger.L) {

	log.Info("POST /api/config")
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	decoder := json.NewDecoder(req.Body)
	var request struct {
		BitmarkConfig  structs.BitmarkdConfiguration
		ProoferdConfig structs.ProoferdConfiguration
	}
	if err := decoder.Decode(&request); nil != err {
		log.Errorf("Error: %v", err)
		response.Result = err
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	_, err := os.OpenFile(bitmarkConfigFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Errorf("Error: %s, %v", bitmarkdConfigUpdateErr, err)
		response.Result = bitmarkdConfigUpdateErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	_, err = os.OpenFile(prooferdConfigFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Errorf("Error: %s, %v", bitmarkdConfigUpdateErr, err)
		response.Result = bitmarkdConfigUpdateErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	bitmarkdConfig, err := structs.NewBitmarkdConfiguration(bitmarkConfigFile)
	if nil != err {
		log.Errorf("Error: %s, %v", bitmarkdConfigUpdateErr, err)
		response.Result = bitmarkdConfigUpdateErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}
	bitmarkdConfig.Chain = chain
	bitmarkdConfig.ClientRPC.Listen = request.BitmarkConfig.ClientRPC.Listen
	bitmarkdConfig.ClientRPC.Announce = request.BitmarkConfig.ClientRPC.Announce
	bitmarkdConfig.Peering.Broadcast = request.BitmarkConfig.Peering.Broadcast
	bitmarkdConfig.Peering.Listen = request.BitmarkConfig.Peering.Listen
	bitmarkdConfig.Proofing.Publish = request.BitmarkConfig.Proofing.Publish
	bitmarkdConfig.Proofing.Submit = request.BitmarkConfig.Proofing.Submit
	bitmarkdConfig.SaveToJson(bitmarkConfigFile)

	prooferdConfig, err := structs.NewProoferdConfiguration(prooferdConfigFile)
	if nil != err {
		log.Errorf("Error: %s, %v", prooferdConfigUpdateErr, err)
		response.Result = prooferdConfigUpdateErr

		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}
	prooferdConfig.Chain = chain
	prooferdConfig.Peering.Connect = request.ProoferdConfig.Peering.Connect
	prooferdConfig.SaveToJson(prooferdConfigFile)

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
