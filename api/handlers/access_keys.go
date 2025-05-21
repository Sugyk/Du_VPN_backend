package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sugyk/rest_vpn/repository"
)

// this handler create the key and add it to db with expires_at
//
//	201 - if creating is successfuls
//	400 - if request body is invalid
//	500 - if failed to create a key
func (s *Service) CreateKey() http.HandlerFunc {
	type request struct {
		TelegramId   int `json:"telegram_id"`
		LiveTimeDays int `json:"livetime"`
	}

	type response struct {
		ExpiresAt string `json:"expires_at"`
		KeyValue  string `json:"access_key"`
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

		livetimeDays := body.LiveTimeDays

		// 400
		if livetimeDays <= 0 {
			log.Println("error: request have invalid body: expire_months is <= 0")
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

		current_time := time.Now().AddDate(0, 0, livetimeDays)
		expires_at := time.Date(
			current_time.Year(),
			current_time.Month(),
			current_time.Day(),
			23,
			59,
			59,
			0,
			current_time.Location(),
		)

		err = s.Repo.CreateKey(repository.CreateKeyParams{
			TelegramId: body.TelegramId,
			AccessUrl:  key_response.AccessUrl,
			ExpiresAt:  expires_at,
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
				ExpiresAt: expires_at.String(),
				KeyValue:  key_response.AccessUrl,
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
