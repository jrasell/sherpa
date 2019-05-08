package client

import (
	nomadAPI "github.com/hashicorp/nomad/api"
)

// NewNomadClient is responsible for generating a reusable Nomad client using the HashiCorp Nomad
// SDK and the default config. This default config pulls Nomad client configuration from env vars
// which can therefore be customized by the user.
func NewNomadClient() (*nomadAPI.Client, error) {
	return nomadAPI.NewClient(nomadAPI.DefaultConfig())
}
