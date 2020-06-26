package client

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
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
