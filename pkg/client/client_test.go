package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/turbonomic/turbo-api/pkg/api"
)

func TestNewTurboClient(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost")
	secureURL, _ := url.Parse("https://localhost")
	table := []struct {
		config         *Config
		service        string
		expectedClient Client
	}{
		{
			config:  &Config{baseURL, &BasicAuthentication{"foo", "bar"}, "", "", ""},
			service: API,
			expectedClient: &APIClient{
				&RESTClient{http.DefaultClient, baseURL, APIPath, &BasicAuthentication{"foo", "bar"}},
				nil, "", "",
			},
		},
		{
			config:  &Config{secureURL, &BasicAuthentication{"foo", "bar"}, "", "", ""},
			service: API,
			expectedClient: &APIClient{
				&RESTClient{&http.Client{Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}}, secureURL, APIPath, &BasicAuthentication{"foo", "bar"}},
				nil, "", "",
			},
		},
		{
			config:  &Config{secureURL, nil, "", "", ""},
			service: TopologyProcessor,
			expectedClient: &TPClient{
				&RESTClient{&http.Client{Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}}, secureURL, TopologyProcessorPath, nil},
			},
		},
	}
	for _, item := range table {
		turboClient, err := NewTurboClient(item.config)
		if err != nil {
			t.Error(err)
		}
		client, found := turboClient.clients[item.service]
		if !found {
			t.Errorf("Client for service %v is not found", item.service)
		}
		if !reflect.DeepEqual(client, item.expectedClient) {
			t.Errorf("Expected client %+v, got %+v", item.expectedClient, client)
		}
	}
}

// Error is expected because of empty address
func TestClient_DiscoverTarget_WithError(t *testing.T) {
	uuid := ""
	baseURL, _ := url.Parse("http://localhost")
	config := &Config{baseURL, &BasicAuthentication{"foo", "bar"}, "", "", ""}
	turboClient, _ := NewTurboClient(config)
	_, err := turboClient.DiscoverTarget(uuid, API)
	if err == nil {
		t.Error("Expected error, but got no error.")
	}
}

func TestClient_AddTarget_WithError(t *testing.T) {
	target := &api.Target{}
	baseURL, _ := url.Parse("http://localhost")
	config := &Config{baseURL, &BasicAuthentication{"foo", "bar"}, "", "", ""}
	turboClient, _ := NewTurboClient(config)
	if err := turboClient.AddTarget(target, API); err == nil {
		t.Error("Expected error, but got no error.")
	}
}

func TestBuildErrorAPIDTO(t *testing.T) {
	table := []struct {
		requestDesc    string
		status         string
		contentMessage string
	}{
		{
			requestDesc:    "target addition",
			status:         "400 Bad Request",
			contentMessage: "some message",
		},
		{
			requestDesc:    "target addition",
			status:         "400 Bad Request",
			contentMessage: "",
		},
	}
	for _, item := range table {
		content := fmt.Sprintf("{\"message\":\"%s\"}", item.contentMessage)
		err := buildResponseError(item.requestDesc, item.status, content)
		expectedErrString := fmt.Sprintf("unsuccessful %s response: %s.", item.requestDesc, item.status)
		if item.contentMessage != "" {
			expectedErrString = fmt.Sprintf("%s %s.", expectedErrString, item.contentMessage)
		}
		expectedErr := errors.New(expectedErrString)
		if !reflect.DeepEqual(err, expectedErr) {
			t.Errorf("Expected error %s, got %s", expectedErrString, err)
		}
	}
}
