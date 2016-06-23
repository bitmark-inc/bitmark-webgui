package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/configuration"
	"github.com/bitmark-inc/bitmark-webgui/fault"
	"github.com/bitmark-inc/bitmark-webgui/services"
	"github.com/bitmark-inc/logger"
	"net/http"
)

type BitmarkCliGenerateResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type BitmarkCliInfoResponse struct {
	Default_identity string                `json:"default_identity"`
	Network          string                `json:"network"`
	Connect          string                `json:"connect"`
	Identities       []BitmarkIdentityType `json:"identities"`
}

type BitmarkPaymentAddress struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
}

type BitmarkCliIssueResponse struct {
	AssetId        string                  `json:"assetId"`
	IssueIds       []string                `json:"issueIds"`
	PaymentAddress []BitmarkPaymentAddress `json:"paymentAddress"`
}

type BitmarkCliTransferResponse struct {
	TransferId     string                  `json:"transferId"`
	PaymentAddress []BitmarkPaymentAddress `json:"paymentAddress"`
}

//POST /api/bitmarkCli/*
func BitmarkCliExec(w http.ResponseWriter, req *http.Request, log *logger.L, command string, webguiFilePath string, configuration *configuration.Configuration) {
	log.Infof("POST /api/bitmarCli/%s", command)
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	switch command {
	case "generate":
		if output, err := bitmarkCliService.Generate(); nil != err {
			response.Result = err
		} else {
			var jsonKeyPair BitmarkCliGenerateResponse
			if err := json.Unmarshal(output, &jsonKeyPair); nil != err {
				log.Errorf("Error: %v", err)
			} else {
				response.Ok = true
				response.Result = jsonKeyPair
			}
		}
	case "info":
		if output, err := bitmarkCliService.Info(configuration.BitmarkCliConfigFile); nil != err {
			response.Result = "bitmark-cli info error"
		} else {
			var jsonInfo BitmarkCliInfoResponse
			if err := json.Unmarshal(output, &jsonInfo); nil != err {
				log.Errorf("Error: %v", err)
			} else {
				response.Ok = true
				response.Result = jsonInfo
			}
		}
	case "setup":
		var request services.BitmarkCliSetupType
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&request); nil != err {
			log.Errorf("Error: %v", err)
			response.Result = "bitmark-cli request parsing error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}

		if _, err := bitmarkCliService.Setup(request, webguiFilePath, configuration); nil != err {
			response.Result = "bitmark-cli setup error"
		} else {
			response.Ok = true
			response.Result = "success"
		}
	case "issue":
		var request services.BitmarkCliIssueType
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&request); nil != err {
			log.Errorf("Error: %v", err)
			response.Result = "bitmark-cli request parsing error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}

		output, err := bitmarkCliService.Issue(request)
		if nil != err {
			response.Result = "bitmark-cli issue error"
		} else {
			var jsonIssue BitmarkCliIssueResponse
			if err := json.Unmarshal(output, &jsonIssue); nil != err {
				log.Errorf("Error: %v", err)
			} else {
				response.Ok = true
				response.Result = jsonIssue
			}
		}
	case "transfer":
		var request services.BitmarkCliTransferType
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&request); nil != err {
			log.Errorf("Error: %v", err)
			response.Result = "bitmark-cli request parsing error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}

		output, err := bitmarkCliService.Transfer(request)
		if nil != err {
			response.Result = "bitmark-cli transfer error"
		} else {
			var jsonTransfer BitmarkCliTransferResponse
			if err := json.Unmarshal(output, &jsonTransfer); nil != err {
				log.Errorf("Error: %v", err)
			} else {
				response.Ok = true
				response.Result = jsonTransfer
			}
		}
	case "keypair":
		var request services.BitmarkCliKeyPairType
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&request); nil != err {
			log.Errorf("Error: %v", err)
			response.Result = "bitmark-cli request parsing error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}

		output, err := bitmarkCliService.KeyPair(request, configuration.BitmarkCliConfigFile)
		if nil != err {
			response.Result = "bitmark-cli keypair error"
		} else {
			var jsonKeyPair BitmarkCliGenerateResponse
			if err := json.Unmarshal(output, &jsonKeyPair); nil != err {
				log.Errorf("Error: %v", err)
			} else {
				response.Ok = true
				response.Result = jsonKeyPair
			}
		}
	default:
		response.Result = fault.ErrInvalidCommandType
	}

	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}
