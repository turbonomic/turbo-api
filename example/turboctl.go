package main

import (
	"fmt"
	"net/url"

	"github.com/golang/glog"
	"github.com/turbonomic/turbo-api/pkg/api"
	"github.com/turbonomic/turbo-api/pkg/client"
)

func main() {
	discoverTargetExample()
}

func addTarget() {
	serverAddress, err := url.Parse("<Server_Address>")
	if err != nil {
		glog.Errorf("Incorrect URL: %s", err)
	}
	config := client.NewConfigBuilder(serverAddress).
		BasicAuthentication("<UI-username>", "UI-password").
		Create()
	turboClient, err := client.NewTurboClient(config)
	if err != nil {
		glog.Errorf("Error creating client: %s", err)
	}

	target := &api.Target{
		Category: "Hypervisor",
		Type:     "vCenter",
		InputFields: []*api.InputField{
			{
				Value:           "<VC_Address>",
				Name:            "nameOrAddress",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<VC_Username>",
				Name:            "username",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<VC_Password>",
				Name:            "password",
				GroupProperties: []*api.List{},
			},
		},
	}
	if err = turboClient.AddTarget(target, client.API); err != nil {
		glog.Errorf("Error adding target: %s", err)
		return
	}
}

// Add an external target. This type of type is registered through SDK.
// Here we use Kubernetes target for example.
func addExternalTarget() {
	// Get Turbonomic server address.
	serverAddress, err := url.Parse("<SERVER_ADDRESS>")
	if err != nil {
		glog.Errorf("Incorrect URL: %s", err)
	}

	// Create API client config.
	config := client.NewConfigBuilder(serverAddress).
		BasicAuthentication("<UI_USERNAME>", "<UI_PASSWORD>").
		Create()
	turboClient, err := client.NewTurboClient(config)
	if err != nil {
		glog.Errorf("Error creating client: %s", err)
	}

	// Configure target data.
	target := &api.Target{
		Category: "Custom",
		Type:     "Kubernetes",
		InputFields: []*api.InputField{
			{
				Value:           "<Kubernetes_TargetID>",
				Name:            "targetIdentifier",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<Kubernetes_Target_Username>",
				Name:            "username",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "<Kubernetes_Target_Password>",
				Name:            "password",
				GroupProperties: []*api.List{},
			},
		},
	}

	// Make API calls.
	if err = turboClient.AddTarget(target, client.API); err != nil {
		fmt.Printf("Error adding target: %s\n", err)
		return
	}
}

func discoverTargetExample() {
	serverAddress, err := url.Parse("<SERVER_ADDRESS>")
	if err != nil {
		glog.Errorf("Incorrect URL: %s", err)
	}
	config := client.NewConfigBuilder(serverAddress).
		BasicAuthentication("<UI_USERNAME>", "<UI_PASSWORD>").
		Create()
	turboClient, err := client.NewTurboClient(config)
	if err != nil {
		glog.Errorf("Error creating client: %s", err)
	}
	uuid := "<TARGET_UUID>"
	resp, err := turboClient.DiscoverTarget(uuid, client.API)
	if err != nil {
		glog.Errorf("Error adding target: %s", err)
		return
	}
	glog.Infof("Response is %+v", resp)
}
