package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/turbonomic/turbo-api/pkg/api"
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

	//var result map[string]interface{}
	//json.Unmarshal([]byte(response.body), &result)
	//
	//for key, value := range result {
	//	// Each value is an interface{} type, that is type asserted as a string
	//	fmt.Println(key, value)
	//}
	//fmt.Printf("AuthToken = %s\n", result["authToken"])

	//authToken, ok := result["authToken"].(string)
	//if ok {
	//	c.AuthToken = authToken
	//}

	sessionCookie, ok := response.cookies[SessionCookie]
	fmt.Printf("%s", response.cookies)
	if ok {
		c.SessionCookie = sessionCookie
		fmt.Printf("Session Cookie = %s:%s\n", c.SessionCookie.Name, c.SessionCookie.Value)
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
	// find if the target exists
	targetExists, _ := c.FindTarget(target)
	if targetExists {
		return nil, fmt.Errorf("Target %v exists", target)
	}

	fmt.Printf("***************** [AddTarget] %++v\n", target)

	targetData, err := json.Marshal(target)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshall target instance: %s", err)
	}
	request := c.Post().Resource(api.Resource_Type_Target).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Data(targetData).
		Header("Cookie", fmt.Sprintf("%s=%s", c.SessionCookie.Name, c.SessionCookie.Value))

	fmt.Printf("Request %++v\n", request)

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
	fmt.Printf("***************** [FindTarget] %++v\n", target)
	targetData, err := json.Marshal(target)
	if err != nil {
		return false, fmt.Errorf("Failed to marshall target instance: %s", err)
	}

	request := c.Get().Resource(api.Resource_Type_Target).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Header("Cookie", fmt.Sprintf("%s=%s", c.SessionCookie.Name, c.SessionCookie.Value)).
		Data(targetData)

	fmt.Printf("Request %++v\n", request)

	response, err := request.Do()

	if err != nil {
		fmt.Printf("Failed to execute find target request: %s", err)
		return false, fmt.Errorf("Failed to execute find target request: %s", err)
	}
	fmt.Printf("[FindTarget] Response %++v\n", response)
	if response.statusCode != 200 {
		return false, buildResponseError("find target", response.status, response.body)
	}

	// Parse the response
	targetId := getTargetId(target)
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
		fmt.Printf("%s::%s::%s\n", category, targetType, tgtId)
	}
	return false, nil
}

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
