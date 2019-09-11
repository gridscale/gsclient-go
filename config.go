package gsclient

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"time"
)

const version = "1.0.0"

const (
	defaultCheckRequestTimeoutSecs     = 120
	defaultServerErrorRetryTimeoutSecs = 60
)

//Config config for client
type Config struct {
	APIUrl                      string
	UserUUID                    string
	APIToken                    string
	UserAgent                   string
	HTTPClient                  *http.Client
	RequestCheckTimeoutSecs     time.Duration
	ServerErrorRetryTimeoutSecs time.Duration
	logger                      logrus.Logger
}

//NewConfiguration creates a new config
func NewConfiguration(apiURL string, uuid string, token string, debugMode bool, requestCheckTimeoutSecs,
	serverErrorRetryTimeoutSecs int) *Config {
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

	if requestCheckTimeoutSecs == 0 {
		requestCheckTimeoutSecs = defaultCheckRequestTimeoutSecs
	}
	if serverErrorRetryTimeoutSecs == 0 {
		serverErrorRetryTimeoutSecs = defaultServerErrorRetryTimeoutSecs
	}
	cfg := &Config{
		APIUrl:                      apiURL,
		UserUUID:                    uuid,
		APIToken:                    token,
		UserAgent:                   "gsclient-go/" + version + " (" + runtime.GOOS + ")",
		HTTPClient:                  http.DefaultClient,
		logger:                      logger,
		RequestCheckTimeoutSecs:     time.Duration(requestCheckTimeoutSecs) * time.Second,
		ServerErrorRetryTimeoutSecs: time.Duration(serverErrorRetryTimeoutSecs) * time.Second,
	}
	return cfg
}
