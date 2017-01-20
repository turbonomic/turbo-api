package main

import (
	"fmt"
	"net/url"

	"github.com/turbonomic/turbo-api/pkg/api"
	"github.com/turbonomic/turbo-api/pkg/client"
)

func main() {
	discoverTargetExample()
}

func addTarget() {
	serverAddress, err := url.Parse("Server-address")
	if err != nil {
		fmt.Errorf("Incorrect URL: %s", err)
	}
	config := client.NewConfigBuilder(serverAddress).
		APIPath("/vmturbo/rest").
		BasicAuthentication("<UI-username>", "UI-password").
		Create()
	client, err := client.NewAPIClientWithBA(config)
	if err != nil {
		fmt.Errorf("Error creating client: %s", err)
	}

	target := &api.Target{
		Category: "Hypervisor",
		Type:     "vCenter",
		InputFields: []*api.InputField{
			{
				Value:           "VC-address",
				Name:            "nameOrAddress",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "username",
				Name:            "VC-username",
				GroupProperties: []*api.List{},
			},
			{
				Value:           "password",
				Name:            "VC-password",
				GroupProperties: []*api.List{},
			},
		},
	}
	resp, err := client.AddTarget(target)
	if err != nil {
		fmt.Errorf("Error adding target: %s", err)
		return
	}
	fmt.Printf("Response is %++v", resp)
}

func discoverTargetExample() {
	serverAddress, err := url.Parse("Server-address")
	if err != nil {
		fmt.Errorf("Incorrect URL: %s", err)
	}
	config := client.NewConfigBuilder(serverAddress).
		APIPath("/vmturbo/rest").
		BasicAuthentication("<UI-username>", "UI-password").
		Create()
	client, err := client.NewAPIClientWithBA(config)
	if err != nil {
		fmt.Errorf("Error creating client: %s", err)
	}
	uuid := "TARGET_UUID"
	resp, err := client.DiscoverTarget(uuid)
	if err != nil {
		fmt.Errorf("Error adding target: %s", err)
		return
	}
	fmt.Printf("Response is %++v", resp)
}
