package api

import (
	"encoding/json"
	"github.com/bitmark-inc/bitmark-webgui/services"
	"github.com/bitmark-inc/logger"
	"net/http"
)

type onestepRequest interface {}

type OnestepStatusRequest struct {
	Network string `json:"network"`
	CliConfig string `json:"cli_config"`
	PayConfig string `json:"pay_config"`
}

type BitmarkIdentityType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public_key  string `json:"public_key"`
}

type OnestepStatusResponse struct {
	Network string `json:"network"`
	Connect string `json:"connect"`
	Identities       []BitmarkIdentityType `json:"identities"`
	Address string `json:"address"`
	EstimatedBalance float64 `json:"estimated_balance"`
	AvailableBalance float64 `json:"available_balance"`
}

type OnestepSetupRequest struct {
	Network string `json:"network"`
	CliConfig string `json:"cli_config"`
	PayConfig string `json:"pay_config"`
	Connect string `json:"connect"`
	Identity    string `json:"identity"`
	Description string `json:"description"`
	CliPassword    string `json:"cli_password"`
	PayPassword  string   `json:"pay_password"`
}


// POST /api/onestep/status, setup, issue, transfer
func OnestepExec(w http.ResponseWriter, req *http.Request, log *logger.L, command string){
	log.Infof("POST /api/onestep/%s", command)

	// get diffrent request instance for json decode
	oneStepRequest := map[string] func() onestepRequest {
		"status": func() onestepRequest {return &OnestepStatusRequest{}},
		"setup": func() onestepRequest {return &OnestepSetupRequest{}},
	}
	request := oneStepRequest[command]()
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(request)
	if nil != err {
		log.Errorf("Error: %v", err)
		response := &Response{
			Ok:     false,
			Result: "bitmarkOnestep "+command+"  request parsing error",
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
	}
}

func execOnestepStatus(w http.ResponseWriter, request OnestepStatusRequest, log *logger.L){
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
		Net: request.Network,
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

func execOnestepSetup(w http.ResponseWriter, request OnestepSetupRequest, log *logger.L){
	response := &Response{
		Ok:     false,
		Result: nil,
	}

	//setup bitmark-cli
	cliRequest := services.BitmarkCliSetupType{
		Config: request.CliConfig,
		Identity: request.Identity,
		Password: request.CliPassword,
		Network: request.Network,
		Connect: request.Connect,
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
		Config: request.PayConfig,
		Net:request.Network,
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
