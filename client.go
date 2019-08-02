package gsclient

const (
	apiServerBase       = "/objects/servers"
	apiStorageBase      = "/objects/storages"
	apiNetworkBase      = "/objects/networks"
	apiIpBase           = "/objects/ips"
	apiSshkeyBase       = "/objects/sshkeys"
	apiTemplateBase     = "/objects/templates"
	apiLoadBalancerBase = "/objects/loadbalancers"
	apiPaaSBase         = "/objects/paas"
	apiISOBase          = "/objects/isoimages"
)

type Client struct {
	cfg *Config
}

func NewClient(c *Config) *Client {
	client := &Client{
		cfg: c,
	}

	return client
}
