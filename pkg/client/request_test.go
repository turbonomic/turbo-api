package client

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/dongyiyang/turbo-api/pkg/api"
)

func TestNewRequest(t *testing.T) {
	// NewRequest(client HTTPClient, verb string, baseURL *url.URL, apiPath string) *Request
	client := http.DefaultClient
	verb := "GET"
	baseURL, _ := url.Parse("http://localhost")
	table := []struct {
		apiPath       string
		expectRequest *Request
	}{
		{
			"",
			&Request{
				client:     client,
				verb:       verb,
				baseURL:    baseURL,
				pathPrefix: "",
			},
		},
		{
			"foo",
			&Request{
				client:     client,
				verb:       verb,
				baseURL:    baseURL,
				pathPrefix: "/foo",
			},
		},
		{
			"/bar",
			&Request{
				client:     client,
				verb:       verb,
				baseURL:    baseURL,
				pathPrefix: "/bar",
			},
		},
	}

	for _, item := range table {
		request := NewRequest(client, verb, baseURL, item.apiPath)
		if !reflect.DeepEqual(request, item.expectRequest) {
			t.Errorf("expected %++v, got %++v", item.expectRequest, request)
		}
	}
}

func TestBasicAuthentication(t *testing.T) {
	table := []struct {
		username   string
		password   string
		expectAuth *BasicAuthentication
	}{
		{"foo", "31415", &BasicAuthentication{"foo", "31415"}},
		{"bar", "", &BasicAuthentication{"bar", ""}},
	}

	u, _ := url.Parse("http://localhost")
	for _, item := range table {
		basicAuth := &BasicAuthentication{item.username, item.password}
		request := NewRequest(http.DefaultClient, "GET", u, "").BasicAuthentication(basicAuth)
		if !reflect.DeepEqual(request.basicAuth, item.expectAuth) {
			t.Errorf("expected %++v, got %++v", item.expectAuth, request)
		}
	}
}

func TestParam(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	table := []struct {
		name      string
		testVal   string
		expectStr string
	}{
		{"foo", "31415", "http://localhost?foo=31415"},
		{"bar", "42", "http://localhost?bar=42"},
		{"baz", "0", "http://localhost?baz=0"},
	}

	for _, item := range table {
		r := NewRequest(http.DefaultClient, "GET", u, "").Param(item.name, item.testVal)
		if e, a := item.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %v, got %v", e, a)
		}
	}
}

func TestName(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	tests := []struct {
		name      string
		expectStr string
	}{
		{"bar", "http://localhost/bar"},
		{"foo", "http://localhost/foo"},
	}
	for _, test := range tests {
		r := NewRequest(http.DefaultClient, "GET", u, "").Name(test.name)
		if e, a := test.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}

func TestResource(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	tests := []struct {
		resource  api.ResourceType
		expectStr string
	}{
		{api.Resource_Type_Target, u.String() + "/targets"},
		{api.Resource_Type_External_Target, u.String() + "/externaltargets"},
	}
	for _, test := range tests {
		r := NewRequest(http.DefaultClient, "GET", u, "").Resource(test.resource)
		if e, a := test.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}

func TestURLInOrder(t *testing.T) {
	u, _ := url.Parse("http://localhost")
	tests := []struct {
		resource     api.ResourceType
		resourceName string
		parameters   map[string]string
		expectStr    string
	}{
		{
			resource:     api.Resource_Type_Target,
			resourceName: "foo",
			expectStr:    "http://localhost/targets/foo",
		},
		{
			resource: api.Resource_Type_External_Target,
			parameters: map[string]string{
				"foo": "12",
			},
			expectStr: "http://localhost/externaltargets?foo=12"},
	}
	for _, test := range tests {
		r := NewRequest(http.DefaultClient, "GET", u, "").Resource(test.resource).Name(test.resourceName)
		for key, val := range test.parameters {
			r.Param(key, val)
		}
		if e, a := test.expectStr, r.URL().String(); e != a {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}
