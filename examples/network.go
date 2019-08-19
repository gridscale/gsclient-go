package main

import (
	"bufio"
	"os"

	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
)

const LocationUuid = "45ed677b-3702-4b36-be2a-a2eab9827950"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration(
		"https://api.gridscale.io",
		uuid,
		token,
		true,
	)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create network: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	networkRequest := gsclient.NetworkCreateRequest{
		Name:         "go-client-network",
		LocationUuid: LocationUuid,
	}
	cnetwork, err := client.CreateNetwork(networkRequest)
	if err != nil {
		log.Fatal("Create network has failed with error", err)
	}
	log.WithFields(log.Fields{
		"network_uuid": cnetwork.ObjectUuid,
	}).Info("Network successfully created")

	//Get network to update
	net, err := client.GetNetwork(cnetwork.ObjectUuid)
	if err != nil {
		log.Fatal("Create network has failed ")
	}

	log.Info("Update network: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	netUpdateRequest := gsclient.NetworkUpdateRequest{
		Name: "Updated network",
	}
	err = client.UpdateNetwork(net.Properties.ObjectUuid, netUpdateRequest)
	if err != nil {
		log.Fatal("Update network has failed with error", err)
	}
	log.WithFields(log.Fields{
		"network_uuid": net.Properties.ObjectUuid,
	}).Info("Network successfully updated")

	log.Info("Retrieve network's events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//get network's events
	events, err := client.GetNetworkEventList(net.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Get network's events has failed with error", err)
	}
	log.WithFields(log.Fields{
		"network_uuid": net.Properties.ObjectUuid,
		"events":       events,
	}).Info("Events successfully retrieved")

	log.Info("Delete network: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//delete network
	err = client.DeleteNetwork(net.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Delete network has failed with error", err)
	}
	log.Info("Network successfully deleted")
}
