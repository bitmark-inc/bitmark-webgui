package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/services"
	"github.com/bitmark-inc/logger"
	"net/http"
)

type onestepRequest interface{}

// POST /api/onestep/status, setup, issue, transfer
func OnestepExec(w http.ResponseWriter, req *http.Request, log *logger.L, command string) {
	log.Infof("POST /api/onestep/%s", command)

	// get diffrent request instance for json decode
	oneStepRequest := map[string]func() onestepRequest{
		"status":   func() onestepRequest { return &OnestepStatusRequest{} },
		"setup":    func() onestepRequest { return &OnestepSetupRequest{} },
		"issue":    func() onestepRequest { return &OnestepIssueRequest{} },
		"transfer": func() onestepRequest { return &OnestepTransferRequest{} },
	}
	request := oneStepRequest[command]()
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(request)
	if nil != err {
		log.Errorf("Error: %v", err)
		response := &Response{
			Ok:     false,
			Result: "bitmarkOnestep " + command + "  request parsing error",
		}
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	switch request.(type) {
	case *OnestepStatusRequest:
		realRequest := request.(*OnestepStatusRequest)
		execOnestepStatus(w, *realRequest, log)
	case *OnestepSetupRequest:
		realRequest := request.(*OnestepSetupRequest)
		execOnestepSetup(w, *realRequest, log)
	case *OnestepIssueRequest:
		realRequest := request.(*OnestepIssueRequest)
		execOnestepIssue(w, *realRequest, log)
	case *OnestepTransferRequest:
		realRequest := request.(*OnestepTransferRequest)
		execOnestepTransfer(w, *realRequest, log)
	}
}

type OnestepStatusRequest struct {
	Network   string `json:"network"`
	CliConfig string `json:"cli_config"`
	PayConfig string `json:"pay_config"`
}

type BitmarkIdentityType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public_key  string `json:"public_key"`
}

type OnestepStatusResponse struct {
	Network          string                `json:"network"`
	Connect          string                `json:"connect"`
	Identities       []BitmarkIdentityType `json:"identities"`
	Address          string                `json:"address"`
	EstimatedBalance float64               `json:"estimated_balance"`
	AvailableBalance float64               `json:"available_balance"`
}

func execOnestepStatus(w http.ResponseWriter, request OnestepStatusRequest, log *logger.L) {
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	var statusResponse OnestepStatusResponse

	// get bitmark-cli info
	cliRequest := services.BitmarkCliInfoType{
		Config: request.CliConfig,
	}
	cliOutput, err := bitmarkCliService.Info(cliRequest)
	if nil != err {
		response.Result = onestepCliInfoErr
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	} else {
		if err := json.Unmarshal(cliOutput, &statusResponse); nil != err {
			log.Errorf("Error: %v", err)
			response.Result = "bitmarkOnestep status response parsing error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}
	}

	//get bitmark-pay info
	payRequest := services.BitmarkPayType{
		Config: request.PayConfig,
		Net:    request.Network,
	}
	if payRequest.Net == "local" {
		payRequest.Net = "local_bitcoin_reg"
	}

	payOutput, err := bitmarkPayService.Info(payRequest)
	if nil != err {
		response.Result = "bitmark-pay info error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	} else {
		if err := json.Unmarshal(payOutput, &statusResponse); nil != err {
			log.Errorf("Error: %v", err)
			response.Result = "bitmarkOnestep status response parsing error"
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}
	}

	// return success response
	response.Ok = true
	response.Result = statusResponse
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

type OnestepSetupRequest struct {
	Network     string `json:"network"`
	CliConfig   string `json:"cli_config"`
	PayConfig   string `json:"pay_config"`
	Connect     string `json:"connect"`
	Identity    string `json:"identity"`
	Description string `json:"description"`
	CliPassword string `json:"cli_password"`
	PayPassword string `json:"pay_password"`
}

func execOnestepSetup(w http.ResponseWriter, request OnestepSetupRequest, log *logger.L) {
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	//setup bitmark-cli
	cliRequest := services.BitmarkCliSetupType{
		Config:      request.CliConfig,
		Identity:    request.Identity,
		Password:    request.CliPassword,
		Network:     request.Network,
		Connect:     request.Connect,
		Description: request.Description,
	}
	_, err := bitmarkCliService.Setup(cliRequest)
	if nil != err {
		response.Result = "bitmark-cli setup error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	//encrypt bitmark-pay
	payRequest := services.BitmarkPayType{
		Config:   request.PayConfig,
		Net:      request.Network,
		Password: request.PayPassword,
	}
	if payRequest.Net == "local" {
		payRequest.Net = "local_bitcoin_reg"
	}
	_, err = bitmarkPayService.Encrypt(payRequest)
	if nil != err {
		response.Result = "bitmark-pay encrypt error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// return success response
	response.Ok = true
	response.Result = "Success"
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

type OnestepIssueRequest struct {
	Network     string `json:"network"`
	CliConfig   string `json:"cli_config"`
	PayConfig   string `json:"pay_config"`
	Identity    string `json:"identity"`
	Asset       string `json:"asset"`
	Description string `json:"description"`
	Fingerprint string `json:"fingerprint"`
	Quantity    int    `json:"quantity"`
	CliPassword string `json:"cli_password"`
	PayPassword string `json:"pay_password"`
}

type OnestepIssueFailResponse struct {
	CliResult BitmarkCliIssueResponse `json:"cli_result"`
	FailStart int                     `json:"fail_start"`
}

func execOnestepIssue(w http.ResponseWriter, request OnestepIssueRequest, log *logger.L) {
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	// bitmark-cli issue
	cliRequest := services.BitmarkCliIssueType{
		Config:      request.CliConfig,
		Identity:    request.Identity,
		Password:    request.CliPassword,
		Asset:       request.Asset,
		Description: request.Description,
		Fingerprint: request.Fingerprint,
		Quantity:    request.Quantity,
	}

	cliOutput, err := bitmarkCliService.Issue(cliRequest)
	if nil != err {
		response.Result = "bitmark-cli issue error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	var cliIssueResponse BitmarkCliIssueResponse
	if err := json.Unmarshal(cliOutput, &cliIssueResponse); nil != err {
		log.Errorf("Error: %v", err)
		response.Result = "bitmark-cli issue success, but parsing fail."
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// bitmark-pay txid address
	payRequest := services.BitmarkPayType{
		Net:       request.Network,
		Config:    request.PayConfig,
		Password:  request.PayPassword,
		Addresses: []string{cliIssueResponse.PaymentAddress[0].Address},
	}
	if payRequest.Net == "local" {
		payRequest.Net = "local_bitcoin_reg"
	}

	for i, issueId := range cliIssueResponse.IssueIds {
		log.Tracef("pay issueId: %s", issueId)
		payRequest.Txid = issueId
		if _, err := bitmarkPayService.Pay(payRequest); nil != err {
			failResponse := OnestepIssueFailResponse{
				FailStart: i,
				CliResult: cliIssueResponse,
			}
			response.Result = failResponse
			if err := writeApiResponseAndSetCookie(w, response); nil != err {
				log.Errorf("Error: %v", err)
			}
			return
		}
	}

	// return success response
	response.Ok = true
	response.Result = cliIssueResponse
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}

type OnestepTransferRequest struct {
	Network     string `json:"network"`
	CliConfig   string `json:"cli_config"`
	PayConfig   string `json:"pay_config"`
	Identity    string `json:"identity"`
	Txid        string `json:"txid"`
	Receiver    string `json:"receiver"`
	CliPassword string `json:"cli_password"`
	PayPassword string `json:"pay_password"`
}

type OnestepTransferFailResponse struct {
	CliResult BitmarkCliTransferResponse `json:"cli_result"`
}

func execOnestepTransfer(w http.ResponseWriter, request OnestepTransferRequest, log *logger.L) {
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	// bitmark-cli transfer
	cliRequest := services.BitmarkCliTransferType{
		Config:   request.CliConfig,
		Identity: request.Identity,
		Password: request.CliPassword,
		Txid:     request.Txid,
		Receiver: request.Receiver,
	}

	output, err := bitmarkCliService.Transfer(cliRequest)
	if nil != err {
		response.Result = "bitmark-cli transfer error"
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	var cliTransfer BitmarkCliTransferResponse
	if err := json.Unmarshal(output, &cliTransfer); nil != err {
		log.Errorf("Error: %v", err)
		response.Result = "bitmark-cli transfer success, but parsing fail."
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// bitmark-pay
	payRequest := services.BitmarkPayType{
		Net:       request.Network,
		Config:    request.PayConfig,
		Password:  request.PayPassword,
		Addresses: []string{cliTransfer.PaymentAddress[0].Address},
		Txid:      cliTransfer.TransferId,
	}
	if payRequest.Net == "local" {
		payRequest.Net = "local_bitcoin_reg"
	}

	if _, err := bitmarkPayService.Pay(payRequest); nil != err {
		failResponse := OnestepTransferFailResponse{
			CliResult: cliTransfer,
		}
		response.Result = failResponse
		if err := writeApiResponseAndSetCookie(w, response); nil != err {
			log.Errorf("Error: %v", err)
		}
		return
	}

	// return success response
	response.Ok = true
	response.Result = cliTransfer
	if err := writeApiResponseAndSetCookie(w, response); nil != err {
		log.Errorf("Error: %v", err)
	}
}
