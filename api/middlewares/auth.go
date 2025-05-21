package middlewares

import (
	"log"
	"net/http"
	"os"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Checking auth token")
		if r.Header.Get("X-Auth-Token") == "" {
			http.Error(w, "Auth token is required", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("X-Auth-Token") != os.Getenv("AUTH_TOKEN") {
			log.Println("middleware: token is invalid", r.Header.Get("X-Auth-Token"))
			http.Error(w, "Auth token is incorrect", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
