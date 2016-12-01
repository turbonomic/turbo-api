package client

import (
	"fmt"
	"net/http"

	"github.com/dongyiyang/turbo-api/pkg/api"
)

type Client struct {
	*RESTClient
}

// Create a Turbo API Client based on basic authentication.
func NewAPIClientWithBA(c *Config) (*Client, error) {
	if c.BasicAuth == nil {
		return nil, fmt.Errorf("Basic authentication is not set")
	}
	restClient := NewRESTClient(http.DefaultClient, c.ServerAddress, c.APIPath).BasicAuthentication(c.BasicAuth)
	return &Client{restClient}, nil
}

// Discover a target using api
// <turbo_server_address>/vmturbo/api/targets/<name_or_address>
func (c *Client) DiscoverTarget(nameOrAddress string) error {
	_, err := c.Post().Resource(api.Resource_Type_Target).Name(nameOrAddress).Do()
	return fmt.Errorf("Failed to discover target %s: %s", nameOrAddress, err)
}

// Add a ExampleProbe target to server
// example : <turbo_server_address>/vmturbo/api/externaltargets?
//                     type=<target_type>&nameOrAddress=<host_address>&username=<username>&targetIdentifier=<target_identifier>&password=<password>
func (c *Client) AddExternalTarget(target *api.Target) error {
	_, err := c.Post().Resource(api.Resource_Type_External_Target).
		Param("type", target.TargetType).
		Param("nameOrAddress", target.NameOrAddress).
		Param("targetIdentifier", target.TargetIdentifier).
		Param("username", target.Username).
		Param("password", target.Password).
		Do()

	return fmt.Errorf("Failed to add external target: %s", err)
}
