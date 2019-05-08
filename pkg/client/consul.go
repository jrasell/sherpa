package client

import consulAPI "github.com/hashicorp/consul/api"

// NewConsulClient is responsible for generating a reusable Consul client using the HashiCorp
// Consul SDK and the default config. This default config pulls Nomad client configuration from env
// vars which can therefore be customized by the user.
func NewConsulClient() (*consulAPI.Client, error) {
	return consulAPI.NewClient(consulAPI.DefaultConfig())
}
