package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getBody(r *http.Request, dest interface{}) error {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &dest); err != nil {
		return err
	}
	return nil
}

func writeResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Write(fmt.Append([]byte{}, body))
}
