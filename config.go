package gsclient

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type Config struct {
	APIUrl     string
	UserUUID   string
	APIToken   string
	HTTPClient *http.Client
	LogLevel   logrus.Level
}

func NewConfiguration(uuid string, token string, logLevel logrus.Level) *Config {
	cfg := &Config{
		APIUrl:     "https://api.gridscale.io",
		UserUUID:   uuid,
		APIToken:   token,
		HTTPClient: http.DefaultClient,
		LogLevel:   logLevel,
	}
	return cfg
}
