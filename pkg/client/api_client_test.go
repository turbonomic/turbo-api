package client

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"github.com/dongyiyang/turbo-api/pkg/api"
)

func TestNewAPIClientWithBA(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	apiPath := "path/to/api"
	table := []struct {
		config         *Config
		expectedClient *Client
		expectsError    bool
	}{
		{
			config:&Config{baseURL, apiPath, nil},
			expectsError:true,
		},
		{
			config:&Config{baseURL, apiPath, &BasicAuthentication{"foo", "bar"}},
			expectedClient:&Client{
				&RESTClient{http.DefaultClient, baseURL, apiPath, &BasicAuthentication{"foo", "bar"}},
			},
			expectsError:false,
		},
	}
	for _, item := range table {
		client, err := NewAPIClientWithBA(item.config)
		if item.expectsError && err == nil {
			t.Error("Expects error, got no error")
		}
		if !reflect.DeepEqual(client, item.expectedClient) {
			t.Errorf("Expected client %++v, got %++v", item.expectedClient, client)
		}
	}
}

func TestClient_DiscoverTarget_WithError(t *testing.T) {
	address := ""
	baseURL, _ := url.Parse("http://localhost")
	apiPath := "path/to/api"
	config:=&Config{baseURL, apiPath, &BasicAuthentication{"foo", "bar"}}
	client, _ := NewAPIClientWithBA(config)
	err := client.DiscoverTarget(address)
	if err == nil {
		t.Error("Expected error, but got no error.")
	}
}

func TestClient_AddExternalTarget_WithError(t *testing.T) {
	target := &api.Target{}
	baseURL, _ := url.Parse("http://localhost")
	apiPath := "path/to/api"
	config:=&Config{baseURL, apiPath, &BasicAuthentication{"foo", "bar"}}
	client, _ := NewAPIClientWithBA(config)
	err := client.AddExternalTarget(target)
	if err == nil {
		t.Error("Expected error, but got no error.")
	}
}
