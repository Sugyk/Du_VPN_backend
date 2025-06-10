package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sugyk/rest_vpn/repository"
)

// this handler create the key and add it to db with expires_at
//
//	201 - if creating is successfull
//	400 - if request body is invalid
//	500 - if failed to create a key
func (s *Service) CreateKey() http.HandlerFunc {
	type request struct {
		TelegramId int `json:"telegram_id"`
	}

	type response struct {
		KeyValue string `json:"access_key"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body request

		// 400
		if err := getBody(r, &body); err != nil {
			log.Println("error: request have invalid body:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "request have invalid body"}`))
			return
		}

		key_response, err := s.OutlineAPI.CreateKey()

		// 500
		if err != nil {
			log.Println("error: failed to create key:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to create key"}`))
			return
		}

		err = s.Repo.CreateKey(repository.CreateKeyParams{
			TelegramId: body.TelegramId,
			AccessUrl:  key_response.AccessUrl,
		})

		// 500
		if err != nil {
			log.Println("error: failed to create key:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to create key"}`))
			// TODO: delete key from outline-server
			return
		}

		resp_bytes, err := json.Marshal(
			response{
				KeyValue: key_response.AccessUrl,
			},
		)

		// 500
		if err != nil {
			log.Println("error: failed to create key:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed generate response"}`))
			// TODO: delete key from outline-server
			// TODO: delete key from repo
			return
		}

		// 201
		w.WriteHeader(http.StatusCreated)
		w.Write(resp_bytes)
	}
}
