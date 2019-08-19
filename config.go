package gsclient

import (
	"net/http"
)

type Config struct {
	APIUrl     string
	UserUUID   string
	APIToken   string
	HTTPClient *http.Client
	DebugMode  bool
}

func NewConfiguration(uuid string, token string, debugMode bool) *Config {
	cfg := &Config{
		APIUrl:     "https://api.gridscale.io",
		UserUUID:   uuid,
		APIToken:   token,
		HTTPClient: http.DefaultClient,
		DebugMode:  debugMode,
	}
	return cfg
}
