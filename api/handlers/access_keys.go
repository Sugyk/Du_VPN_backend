package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sugyk/rest_vpn/repository"
)

func (s *Service) CreateKey() http.HandlerFunc {
	type request struct {
		Telegram_id   int `json:"telegram_id"`
		Expire_months int `json:"expire_months"`
	}

	type response struct {
		ExpireAt string `json:"expire_at"`
		KeyValue string `json:"access_key"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body request
		if err := getBody(r, &body); err != nil {
			log.Println("error: request have invalid body:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "request have invalid body"}`))
			return
		}

		expire_months := body.Expire_months
		if expire_months == 0 {
			log.Println("error: request have invalid body: expire_months is 0")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "request have invalid body"}`))
			return
		}
		key_response, err := s.OutlineAPI.CreateKey()
		if err != nil {
			log.Println("error: failed to create key:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to create key"}`))
			return
		}

		current_time := time.Now().AddDate(0, 0, 30*body.Expire_months)
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

		if err = s.Repo.CreateKey(repository.CreateKeyParams{
			Telegram_id:   body.Telegram_id,
			AccessUrl:     key_response.AccessUrl,
			Expire_months: expires_at,
		}); err != nil {
			log.Println("error: failed to create key:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to create key"}`))
			// TODO: delete key from outline-server
			return
		}

		resp_bytes, err := json.Marshal(
			response{
				ExpireAt: expires_at.String(),
				KeyValue: key_response.AccessUrl,
			},
		)
		if err != nil {
			log.Println("error: failed to create key:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed generate response"}`))
			// TODO: delete key from outline-server
			// TODO: delete key from repo
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(resp_bytes)
	}
}
