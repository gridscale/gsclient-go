package main

import (
	"bufio"
	"os"

	"github.com/nvthongswansea/gsclient-go"
	log "github.com/sirupsen/logrus"
)

const LocationUuid = "45ed677b-3702-4b36-be2a-a2eab9827950"
const webServerFirewallTemplateUuid = "82aa235b-61ba-48ca-8f47-7060a0435de7"

type ServiceType string

const (
	Server   ServiceType = "server"
	Storage  ServiceType = "storage"
	Network  ServiceType = "network"
	IP       ServiceType = "ip"
	ISOImage ServiceType = "isoimage"
)

//enhancedClient inherits all methods from gsclient.Client
//We need this to implement a new additional method
type enhancedClient struct {
	*gsclient.Client
}

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration(
		"https://api.gridscale.io",
		uuid,
		token,
		true,
	)
	client := enhancedClient{
		gsclient.NewClient(config),
	}
	log.Info("gridscale client configured")

	log.Info("Create server: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	serverCreateRequest := gsclient.ServerCreateRequest{
		Name:         "go-client-server",
		Memory:       1,
		Cores:        1,
		LocationUuid: LocationUuid,
	}
	cServer, err := client.CreateServer(serverCreateRequest)
	if err != nil {
		log.Fatal("Create server has failed with error", err)
	}
	log.WithFields(log.Fields{
		"server_uuid": cServer.ObjectUuid,
	}).Info("Server successfully created")
	defer client.deleteService(Server, cServer.ObjectUuid)

	//get a server to interact with
	server, err := client.GetServer(cServer.ObjectUuid)
	if err != nil {
		log.Fatal("Get server has failed with error", err)
	}

	log.Info("Start server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Turn on server
	err = client.StartServer(server.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Start server has failed with error", err)
	}
	log.Info("Server successfully started")

	log.Info("Stop server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Turn off server
	err = client.StopServer(server.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Stop server has failed with error", err)
	}
	log.Info("Server successfully stop")

	log.Info("Update server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	err = client.UpdateServer(server.Properties.ObjectUuid, gsclient.ServerUpdateRequest{
		Name:   "updated server",
		Memory: 2,
	})
	if err != nil {
		log.Fatal("Update server has failed with error", err)
	}
	log.Info("Server successfully updated")

	//Get events of server
	events, err := client.GetServerEventList(server.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Get events has failed with error", err)
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("Events successfully retrieved")

	//Create storage, network, IP, and ISO-image to attach to the server
	log.Info("Create storage, Network, IP, ISO-image: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cStorage, err := client.CreateStorage(gsclient.StorageCreateRequest{
		Capacity:     1,
		LocationUuid: LocationUuid,
		Name:         "go-client-storage",
	})
	if err != nil {
		log.Fatal("Create storage has failed with error", err)
	}
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUuid,
	}).Info("Storage successfully created")
	defer client.deleteService(Storage, cStorage.ObjectUuid)

	cNetwork, err := client.CreateNetwork(gsclient.NetworkCreateRequest{
		Name:         "go-client-network",
		LocationUuid: LocationUuid,
	})
	if err != nil {
		log.Fatal("Create network has failed with error", err)
	}
	log.WithFields(log.Fields{
		"network_uuid": cNetwork.ObjectUuid,
	}).Info("Network successfully created")
	defer client.deleteService(Network, cNetwork.ObjectUuid)

	cIp, err := client.CreateIp(gsclient.IpCreateRequest{
		Name:         "go-client-ip",
		Family:       4,
		LocationUuid: LocationUuid,
	})
	if err != nil {
		log.Fatal("Create IP has failed with error", err)
	}
	log.WithFields(log.Fields{
		"IP_uuid": cIp.ObjectUuid,
	}).Info("IP successfully created")
	defer client.deleteService(IP, cIp.ObjectUuid)

	cISOimage, err := client.CreateISOImage(gsclient.ISOImageCreateRequest{
		Name:         "go-client-iso",
		SourceUrl:    "http://releases.ubuntu.com/16.04.4/ubuntu-16.04.4-server-amd64.iso?_ga=2.188975915.108704605.1521033305-403279979.1521033305",
		LocationUuid: LocationUuid,
	})
	if err != nil {
		log.Fatal("Create ISO-image has failed with error", err)
	}
	log.WithFields(log.Fields{
		"isoimage_uuid": cISOimage.ObjectUuid,
	}).Info("ISO-image successfully created")
	defer client.deleteService(ISOImage, cISOimage.ObjectUuid)

	//Attach storage, network, IP, and ISO-image to a server
	err = client.LinkStorage(server.Properties.ObjectUuid, cStorage.ObjectUuid, false)
	if err != nil {
		log.Fatal("Link storage has failed with error", err)
	}
	log.Info("Storage successfully attached")
	defer client.unlinkService(Storage, server.Properties.ObjectUuid, cStorage.ObjectUuid)

	err = client.LinkNetwork(
		server.Properties.ObjectUuid,
		cNetwork.ObjectUuid,
		webServerFirewallTemplateUuid,
		false,
		1,
		nil,
		gsclient.FirewallRules{},
	)
	if err != nil {
		log.Fatal("Link network has failed with error", err)
	}
	log.Info("Network successfully linked")
	defer client.unlinkService(Network, server.Properties.ObjectUuid, cNetwork.ObjectUuid)

	err = client.LinkIp(server.Properties.ObjectUuid, cIp.ObjectUuid)
	if err != nil {
		log.Fatal("Link IP has failed with error", err)
	}
	log.Info("IP successfully linked")
	defer client.unlinkService(IP, server.Properties.ObjectUuid, cIp.ObjectUuid)

	err = client.LinkIsoImage(server.Properties.ObjectUuid, cISOimage.ObjectUuid)
	if err != nil {
		log.Fatal("Link ISO-image has failed with error", err)
	}
	log.Info("ISO-image successfully linked")
	defer client.unlinkService(ISOImage, server.Properties.ObjectUuid, cISOimage.ObjectUuid)

	log.Info("Unlink and delete: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *enhancedClient) deleteService(serviceType ServiceType, id string) {
	switch serviceType {
	case Server:
		err := c.DeleteServer(id)
		if err != nil {
			log.Fatal("Delete server has failed with error", err)
		}
		log.Info("Server successfully deleted")
	case Storage:
		err := c.DeleteStorage(id)
		if err != nil {
			log.Fatal("Delete storage has failed with error", err)
		}
		log.Info("Storage successfully deleted")
	case Network:
		err := c.DeleteNetwork(id)
		if err != nil {
			log.Fatal("Delete network has failed with error", err)
		}
		log.Info("Network successfully deleted")
	case IP:
		err := c.DeleteIp(id)
		if err != nil {
			log.Fatal("Delete IP has failed with error", err)
		}
		log.Info("IP successfully deleted")
	case ISOImage:
		err := c.DeleteISOImage(id)
		if err != nil {
			log.Fatal("Delete ISO-image has failed with error", err)
		}
		log.Info("ISO-image successfully deleted")
	default:
		log.Fatal("Unknown service type")
	}
}

func (c *enhancedClient) unlinkService(serviceType ServiceType, serverId, serviceId string) {
	switch serviceType {
	case Storage:
		err := c.UnlinkStorage(serverId, serviceId)
		if err != nil {
			log.Fatal("Unlink storage has failed with error", err)
		}
		log.Info("Storage successfully unlinked")
	case Network:
		err := c.UnlinkNetwork(serverId, serviceId)
		if err != nil {
			log.Fatal("Unlink network has failed with error", err)
		}
		log.Info("Network successfully unlinked")
	case IP:
		err := c.UnlinkIp(serverId, serviceId)
		if err != nil {
			log.Fatal("Unlink IP has failed with error", err)
		}
		log.Info("IP successfully unlinked")
	case ISOImage:
		err := c.UnlinkIsoImage(serverId, serviceId)
		if err != nil {
			log.Fatal("Unlink ISO-image has failed with error", err)
		}
		log.Info("ISO-image successfully unlinked")
	default:
		log.Fatal("Unknown service type")
	}
}
