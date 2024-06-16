package models

type Config struct {
	ListenPort       string `json:"listenPort"`
	ConnectionString string `json:"ConnectionString"`
	APIURL           string `json:"apiURL"`
	IsProd           bool   `json:"isProd"`
}
