package api_test

import (
	"fmt"
	"github.com/bitmark-inc/bitmark-webgui/api"
	"github.com/bitmark-inc/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	reader io.Reader //Ignore this for now
	url    string
)

func init() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/api/bitmarkPay/", handleBitmarkPay)
	serveMux.HandleFunc("/api/bitmarkPay/encrypt", handleBitmarkPay)
	serveMux.HandleFunc("/api/bitmarkPay/info", handleBitmarkPay)
	serveMux.HandleFunc("/api/bitmarkPay/pay", handleBitmarkPay)
	serveMux.HandleFunc("/api/bitmarkPay/status", handleBitmarkPay)
	serveMux.HandleFunc("/api/bitmarkPay/result", handleBitmarkPay)

	//Creating new server with the handlers
	// server = httptest.NewTLSServer(serveMux)
	server = httptest.NewServer(serveMux)

	//Grab the address for the API endpoint
	url = fmt.Sprintf("%s/api/bitmarkPay", server.URL)

}

func handleBitmarkPay(w http.ResponseWriter, req *http.Request) {
	log := logger.New("api-bitmarkPay-test")
	api.SetCORSHeader(w, req)

	switch req.Method {
	case `GET`:
		api.BitmarkPayJobHash(w, req, log)
	case `POST`:
		reqUriArr := strings.Split(req.RequestURI, "/")
		api.BitmarkPay(w, req, log, reqUriArr[3])
	case `DELETE`:
		api.BitmarkPayKill(w, req, log)
	case `OPTIONS`:
		return
	default:
		log.Error("Error: Unknow method")
	}
}

// Will be panic because the bitmarkPay service should be initilized
func TestGetBitmarkPayInfo(t *testing.T) {

	reqJson := `{"username": "dennis", "balance": 200}`

	//Create request with JSON body
	reader = strings.NewReader(reqJson) //Convert string to reader

	request, err := http.NewRequest("POST", url+"/info", reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
	}
}
