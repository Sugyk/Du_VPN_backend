package handlers

import (
	"github.com/jmoiron/sqlx"
	"github.com/sugyk/rest_vpn/lib/configs"
	"github.com/sugyk/rest_vpn/outline_api"
	"github.com/sugyk/rest_vpn/repository"
)

type Service struct {
	Repo       *repository.Repository
	OutlineAPI *outline_api.OutlineAPI
}

func newService(db *sqlx.DB, outlineAPI_url configs.OutlineAPIConfig) *Service {
	return &Service{
		Repo:       repository.NewRepo(db),
		OutlineAPI: outline_api.NewOutlineAPI(outlineAPI_url),
	}
}
