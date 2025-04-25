package outline_api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (o *OutlineAPI) CreateKey() (CreateKeyResponse, error) {
	// This function will create a key for user
	// from the outline server

	// Allow self signed certificate on transport level
	transport_level := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Create the client with custom transport level
	httpClient := &http.Client{
		Transport: &transport_level,
	}

	// Send the request of creating key
	resp, err := httpClient.Post(o.API_URL+"/access-keys", "application/json", bytes.NewReader([]byte{}))

	if err != nil {
		return CreateKeyResponse{}, err
	}

	if resp.StatusCode == http.StatusCreated {
		response := CreateKeyResponse{}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return CreateKeyResponse{}, fmt.Errorf("Error while reading response body: %v", err)
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return CreateKeyResponse{}, fmt.Errorf("Error while parsing response body: %v", err)
		}
		log.Println("Access key created successfully")
		return response, nil
	}
	return CreateKeyResponse{}, fmt.Errorf("Wrong status when created a key: %d", resp.StatusCode)
}
