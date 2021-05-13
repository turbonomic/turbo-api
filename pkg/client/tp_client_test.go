package client

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/turbonomic/turbo-api/pkg/api"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestGetProbeIDWithRetry(t *testing.T) {
	probeType := "Kubernetes"
	probeCategory := "Cloud Native"
	baseURL, _ := url.Parse("http://localhost")
	tpClient := &TPClient{
		&RESTClient{&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}, baseURL, TopologyProcessorPath, nil}}
	start := time.Now()
	_, err := tpClient.getProbeID(probeType, probeCategory)
	assert.Error(t, err)
	assert.True(t, time.Since(start).Seconds() >= float64(retryAttempts-1)*retryDelay.Seconds())
	fmt.Println(err)
}

func TestExtractCommunicationBindingChannel(t *testing.T) {
	communicationBindingChannel := "xoxo"
	inputField1 := api.InputField{Name: "foo", Value: "123"}
	inputField2 := api.InputField{Name: api.CommunicationBindingChannel, Value: communicationBindingChannel}
	inputField3 := api.InputField{Name: "bar", Value: "456"}
	inputFields := []*api.InputField{&inputField1, &inputField2, &inputField3}
	expectedInputFields := []*api.InputField{&inputField1, &inputField3}

	tpClient := TPClient{}
	// with all 3 fields as input
	extractedInputFields, extractedChannel := tpClient.extractCommunicationBindingChannel(inputFields)
	assert.Equalf(t, communicationBindingChannel, extractedChannel,
		"Actual extracted communication binding channel %v is different than the expected %v", extractedChannel,
		communicationBindingChannel)
	assert.True(t, reflect.DeepEqual(expectedInputFields, extractedInputFields), "Expected input fields: %v,"+
		" are not the same as the actual: %v", expectedInputFields, extractedInputFields)

	// with only field 1 and field 3
	extractedInputFields, extractedChannel = tpClient.extractCommunicationBindingChannel(expectedInputFields)
	assert.Equalf(t, "", extractedChannel,
		"Actual extracted communication binding channel %v should be an empty string but not", extractedChannel)
	assert.True(t, reflect.DeepEqual(expectedInputFields, extractedInputFields), "Expected input fields: %v,"+
		" are not the same as the actual: %v", expectedInputFields, extractedInputFields)
}
