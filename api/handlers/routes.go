package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sugyk/rest_vpn/api/middlewares"
	"github.com/sugyk/rest_vpn/lib/configs"
)

func registerWithMiddleware(r *mux.Router, path string, handlerFunc http.HandlerFunc, method string) {
	r.HandleFunc(path, middlewares.AuthMiddleware(handlerFunc)).Methods(method)
}

func Register(r *mux.Router, db *sqlx.DB, outlineAPIConfig configs.OutlineAPIConfig) {
	handler := newService(db, outlineAPIConfig)
	registerWithMiddleware(r, "/api/v1/create_user", handler.CreateUser(), http.MethodPost)
	registerWithMiddleware(r, "/api/v1/access_key", handler.CreateKey(), http.MethodPost)
}
