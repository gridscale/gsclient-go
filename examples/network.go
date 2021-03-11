package main

import (
	"bufio"
	"context"
	"os"

	"github.com/gridscale/gsclient-go/v3"
	log "github.com/sirupsen/logrus"
)

var emptyCtx = context.Background()

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create network: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	networkRequest := gsclient.NetworkCreateRequest{
		Name: "go-client-network",
	}
	cnetwork, err := client.CreateNetwork(emptyCtx, networkRequest)
	if err != nil {
		log.Error("Create network has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"network_uuid": cnetwork.ObjectUUID,
	}).Info("Network successfully created")
	defer func() {
		// delete network
		err := client.DeleteNetwork(emptyCtx, cnetwork.ObjectUUID)
		if err != nil {
			log.Error("Delete network has failed with error", err)
			return
		}
		log.Info("Network successfully deleted")

		log.Info("Get deleted networks: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		networks, err := client.GetDeletedNetworks(emptyCtx)
		if err != nil {
			log.Error("Get deleted networks has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"networks": networks,
		}).Info("Retrieved deleted networks successfully")
	}()

	// Get network to update
	net, err := client.GetNetwork(emptyCtx, cnetwork.ObjectUUID)
	if err != nil {
		log.Error("Create network has failed ")
		return
	}

	log.Info("Update network: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	netUpdateRequest := gsclient.NetworkUpdateRequest{
		Name: "Updated network",
	}
	err = client.UpdateNetwork(emptyCtx, net.Properties.ObjectUUID, netUpdateRequest)
	if err != nil {
		log.Error("Update network has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"network_uuid": net.Properties.ObjectUUID,
	}).Info("Network successfully updated")

	log.Info("Retrieve network's events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// get network's events
	events, err := client.GetNetworkEventList(emptyCtx, net.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get network's events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"network_uuid": net.Properties.ObjectUUID,
		"events":       events,
	}).Info("Events successfully retrieved")

	log.Info("Delete network: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
