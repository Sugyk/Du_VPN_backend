package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sugyk/rest_vpn/repository"
)

// this handler create a new user, who is admin or not
//
//	200 - if user created or already in db
//	500 - if something is wrong
func (s *Service) CreateUser() http.HandlerFunc {
	type request struct {
		Telegram_id int  `json:"telegram_id"`
		Is_admin    bool `json:"is_admin"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Handle request
		req_body := &request{}
		err := getBody(r, &req_body)
		if err != nil {
			log.Println("Error parsing request body: ", err)
		}

		// Create user in repo
		err = s.Repo.CreateUser(repository.CreateUserParams{
			Telegram_id: req_body.Telegram_id,
			Is_admin:    req_body.Is_admin,
		})

		// Form response
		if err != nil {
			// 500
			log.Println("Error of creating of user: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// 200
			log.Println("User successfully created")
			w.WriteHeader(http.StatusOK)
			w.Write(fmt.Append([]byte{}, `{"message": "OK"}`))
		}
	}
}
