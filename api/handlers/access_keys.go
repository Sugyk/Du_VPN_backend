package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
		outline_id, _ := strconv.Atoi(key_response.Id)
		err = s.Repo.CreateKey(repository.CreateKeyParams{
			TelegramId: body.TelegramId,
			AccessUrl:  key_response.AccessUrl,
			Outline_id: outline_id,
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

// this handler return rows of access_keys which pass a filter params
//
//	200 - OK, return
//	500 - error on server
//	400 - bad filter params
func (s *Service) GetAccessKeysList() http.HandlerFunc {
	type AccessKey struct {
		Id        int       `json:"id"`
		KeyValue  string    `json:"key_value"`
		UserId    int       `json:"user_id"`
		OutlineId int       `json:"outline_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	type response struct {
		AccessKeys []AccessKey `json:"access_keys"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		urlParams := r.URL.Query()

		// this params can be used in query filters (e.g. /?id=2&user_id=2)
		// it is built per request, so filters of one request do not leak into another
		filters := map[string]any{
			"id":         0,
			"user_id":    0,
			"outline_id": 0,
		}

		for key, keyList := range urlParams {
			if len(keyList) != 1 {
				writeResponse(
					w,
					http.StatusBadRequest,
					fmt.Sprintf(`{"error": "each filter parameter can be represented in query at most one time: %s"}`, key),
				)
				return
			} else if paramType, ok := filters[key]; ok {
				switch paramType.(type) {
				case int:
					converted, err := strconv.Atoi(keyList[0])
					if err != nil {
						writeResponse(
							w,
							http.StatusBadRequest,
							fmt.Sprintf(`{"error": "incorrect param %s: %s. %e"}`, key, keyList[0], err),
						)
						return
					}
					filters[key] = converted
				}
			} else {
				// return error
				writeResponse(
					w,
					http.StatusBadRequest,
					fmt.Sprintf(`{"error": "unexpected filter parameter: %s"}`, key),
				)
				return
			}
		}
		entries, err := s.Repo.ListEntries(filters)

		// 500
		if err != nil {
			// return error of db
			log.Printf(`error: db: %e`, err)
			writeResponse(
				w,
				http.StatusInternalServerError,
				`{"error": "internal error"}`,
			)
			return
		}

		// forming the response body
		responseBody := response{
			AccessKeys: []AccessKey{},
		}
		for _, v := range entries {
			responseBody.AccessKeys = append(responseBody.AccessKeys, AccessKey(v))
		}

		// marshal the struct to []byte
		bytesBody, _ := json.Marshal(responseBody)

		// 200
		// return the list of needed access_keys
		writeResponse(
			w,
			http.StatusOK,
			string(bytesBody),
		)
	}
}
