package outline_api

import "github.com/sugyk/rest_vpn/lib/configs"

type OutlineAPI struct {
	// The layers of the Outline API
	API_URL string
}

func NewOutlineAPI(outline_cnf configs.OutlineAPIConfig) *OutlineAPI {
	return &OutlineAPI{
		API_URL: outline_cnf.API_URL,
	}
}
