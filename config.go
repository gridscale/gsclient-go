package gsclient

import (
	"net/http"
	"os"

	version                        = "1.0.0"
	defaultAPIURL                  = "https://api.gridscale.io"
	version                        = "1.0.0"
	version                        = "1.0.0"
)

//Config config for client
type Config struct {
	APIUrl     string
	UserUUID   string
	APIToken   string
	HTTPClient *http.Client
	logger     logrus.Logger
}

//NewConfiguration creates a new config
func NewConfiguration(apiURL string, uuid string, token string, debugMode bool) *Config {
	logLevel := logrus.InfoLevel
	if debugMode {
		logLevel = logrus.DebugLevel
	}

	logger := logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: false,
		},
	}

	cfg := &Config{
		APIUrl:     apiURL,
		UserUUID:   uuid,
		APIToken:   token,
		HTTPClient: http.DefaultClient,
		logger:     logger,
	}
	return cfg
}

//DefaultConfiguration creates a default configuration
func DefaultConfiguration(uuid string, token string) *Config {
	logger := logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: false,
		},
	}
	cfg := &Config{
		apiURL:                  defaultAPIURL,
		userUUID:                uuid,
		apiToken:                token,
		userAgent:               "gsclient-go/" + version + " (" + runtime.GOOS + ")",
		sync:                    true,
		httpClient:              http.DefaultClient,
		logger:                  logger,
		requestCheckTimeoutSecs: time.Duration(defaultCheckRequestTimeoutSecs) * time.Second,
		delayInterval:           time.Duration(defaultDelayIntervalMilliSecs) * time.Millisecond,
		maxNumberOfRetries:      defaultMaxNumberOfRetries,
	}
	return cfg
}
