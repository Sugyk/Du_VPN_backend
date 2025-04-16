package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func getBody(r *http.Request, resp interface{}) error {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}
	return nil
}
