package outline_api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sugyk/rest_vpn/lib/configs"
	"github.com/sugyk/rest_vpn/outline_api"
)

func TestGetKeys(t *testing.T) {

	// Prepare the mock response
	mock_response := outline_api.CreateKeyResponse{
		AccessKey: outline_api.AccessKey{
			Id:        "0",
			Name:      "my_key",
			Password:  "6Tb3uWK3bKcFMiAMgwOeB1",
			Port:      31766,
			Method:    "chacha20-ietf-poly1305",
			AccessUrl: "ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTo2VGIzdVdLM2JLY0ZNaUFNZ3dPZUIx@85.192.25.246:31766/?outline=1",
		},
	}

	// convert mock response to json
	mock_response_json, err := json.Marshal(&mock_response)
	if err != nil {
		t.Fatal(err)
	}

	// Create the mock server
	mockserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(mock_response_json)
	}))
	defer mockserver.Close()

	// Create layer instance
	outlineAPI := outline_api.NewOutlineAPI(
		configs.OutlineAPIConfig{
			API_URL: mockserver.URL,
		},
	)

	// Send tested function of layer
	response, err := outlineAPI.CreateKey()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(response, mock_response) {
		t.Error("GetKeys returned unexpected result")
	}
}
