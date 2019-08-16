package gsclient

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

const (
	apiServerBase        = "/objects/servers"
	apiStorageBase       = "/objects/storages"
	apiNetworkBase       = "/objects/networks"
	apiIpBase            = "/objects/ips"
	apiSshkeyBase        = "/objects/sshkeys"
	apiTemplateBase      = "/objects/templates"
	apiLoadBalancerBase  = "/objects/loadbalancers"
	apiPaaSBase          = "/objects/paas"
	apiISOBase           = "/objects/isoimages"
	apiObjectStorageBase = "/objects/objectstorages"
	apiFirewallBase      = "/objects/firewalls"
)

type Client struct {
	cfg    *Config
	logger *logrus.Logger
}

func NewClient(c *Config) *Client {
	var logLevel logrus.Level
	if c.LogLevel == logLevel {
		logLevel = logrus.ErrorLevel
	} else {
		logLevel = c.LogLevel
	}
	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
	client := &Client{
		cfg:    c,
		logger: logger,
	}

	return client
}
