package outline_api

type AccessKey struct {
	// This struct will represent the access key

	Id        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	AccessUrl string `json:"accessUrl"`
}

type CreateKeyResponse struct {
	AccessKey
}
