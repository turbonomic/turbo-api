package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/turbonomic/turbo-api/pkg/api"
	"github.com/golang/glog"
)

type Client struct {
	*RESTClient
	SessionCookie *http.Cookie
}

const (
	SessionCookie string = "JSESSIONID"
)

// Create a Turbo API Client based on basic authentication.
func NewAPIClientWithBA(c *Config) (*Client, error) {
	if c.basicAuth == nil {
		return nil, errors.New("Basic authentication is not set")
	}
	client := http.DefaultClient
	// If use https, disable the security check.
	if c.serverAddress.Scheme == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	restClient := NewRESTClient(client, c.serverAddress, c.apiPath).BasicAuthentication(c.basicAuth)
	return &Client{restClient}, nil
}

// Login to the Turbo API server
func (c *Client) Login() (*Result, error) {
	var data []byte
	data = []byte(fmt.Sprintf("username=%s&password=%s", c.basicAuth.username, c.basicAuth.password))
	request := c.Post().Resource("login").
		Header("Content-Type", "application/x-www-form-urlencoded").
		Data(data)

	response, err := request.Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to login  %s: %s", c.baseURL, err)
	}
	if response.statusCode != 200 {
		return nil, buildResponseError("Turbo server login", response.status, response.body)
	}

	// Save the session cookie
	sessionCookie, ok := response.cookies[SessionCookie]
	if ok {
		c.SessionCookie = sessionCookie
		glog.V(4).Infof("Session Cookie = %s:%s\n", c.SessionCookie.Name, c.SessionCookie.Value)
	} else {
		return nil, buildResponseError("Invalid session cookie", response.status, fmt.Sprintf("%s", response.cookies))
	}
	return &response, nil
}

// Discover a target using API
func (c *Client) DiscoverTarget(uuid string) (*Result, error) {
	response, err := c.Post().Resource(api.Resource_Type_Target).Name(uuid).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to discover target %s: %s", uuid, err)
	}
	if response.statusCode != 200 {
		return nil, buildResponseError("target discovery", response.status, response.body)
	}
	return &response, nil
}

//Add a target using API
func (c *Client) AddTarget(target *api.Target) (*Result, error) {
	// find if the target exists - this is a workaround since in current XL server,
	// duplicate targets can be added
	targetExists, _ := c.FindTarget(target)
	if targetExists {
		return nil, fmt.Errorf("Target %v exists", target)
	}

	glog.V(2).Infof("***************** [AddTarget] %++v\n", target)

	targetData, err := json.Marshal(target)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshall target instance: %s", err)
	}
	request := c.Post().Resource(api.Resource_Type_Target).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Data(targetData).
		Header("Cookie", fmt.Sprintf("%s=%s", c.SessionCookie.Name, c.SessionCookie.Value))

	glog.V(4).Infof("[AddTarget] Request %++v\n", request)

	response, err := request.Do()

	if err != nil {
		return nil, fmt.Errorf("Failed to add target: %s", err)
	}
	if response.statusCode != 200 {
		return nil, buildResponseError("target addition", response.status, response.body)
	}
	return &response, nil
}

//Add a target using API
func (c *Client) FindTarget(target *api.Target) (bool, error) {
	targetData, err := json.Marshal(target)
	if err != nil {
		return false, fmt.Errorf("Failed to marshall target instance: %s", err)
	}

	request := c.Get().Resource(api.Resource_Type_Target).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Header("Cookie", fmt.Sprintf("%s=%s", c.SessionCookie.Name, c.SessionCookie.Value)).
		Data(targetData)

	glog.V(4).Infof("[FindTarget] Request %++v\n", request)

	response, err := request.Do()

	if err != nil {
		fmt.Printf("Failed to execute find target request: %s", err)
		return false, fmt.Errorf("Failed to execute find target request: %s", err)
	}

	if response.statusCode != 200 {
		return false, buildResponseError("find target", response.status, response.body)
	}
	// Target identifier for the given target
	targetId := getTargetId(target)

	// Parse the response
	var targetList []interface{}
	json.Unmarshal([]byte(response.body), &targetList)

	for _, tgt := range targetList {
		m := tgt.(map[string]interface{})
		var category, targetType, tgtId string
		for _, v := range m {
			category, _ = m["category"].(string)
			if target.Category != category {
				continue
			}
			targetType, _ = m["type"].(string)
			if target.Type != targetType {
				continue
			}
			switch val := v.(type) {
			case []interface{}:
				for _, value := range val {
					inputFieldMap, ok := value.(map[string]interface{})
					if ok {
						field, ok := inputFieldMap["name"].(string)
						if ok && field == "targetIdentifier" {
							tgtId, _ = inputFieldMap["value"].(string)
							if tgtId == targetId {
								return true, nil
							}
						}
					}
				}
			default:
			}
		}
		glog.V(4).Infof("%s::%s::%s\n", category, targetType, tgtId)
	}
	return false, nil
}

// Get the target identifier for the given target
func getTargetId(target *api.Target) string {
	for _, inputField := range target.InputFields {
		field := inputField.Name
		if field == "targetIdentifier" {
			tgtId := inputField.Value
			return tgtId
		}
	}
	return ""
}

func buildResponseError(requestDesc string, status string, content string) error {
	errorMsg := fmt.Sprintf("unsuccessful %s response: %s.", requestDesc, status)
	errorDTO, err := parseAPIErrorDTO(content)
	if err == nil && errorDTO.Message != "" {
		// Add error message only if we can parse result content to errorDTO.
		errorMsg = errorMsg + fmt.Sprintf(" %s.", errorDTO.Message)
	}
	return fmt.Errorf("%s", errorMsg)
}
