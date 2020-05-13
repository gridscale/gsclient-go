package main

import (
	"bufio"
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gridscale/gsclient-go/v3"
)

const webServerFirewallTemplateUUID = "82aa235b-61ba-48ca-8f47-7060a0435de7"

var emptyCtx = context.Background()

type serviceType int

const (
	serverType   serviceType = iota
	storageType              = iota
	networkType              = iota
	ipType                   = iota
	isoImageType             = iota
)

//enhancedClient inherits all methods from gsclient.Client
//We need this to implement a new additional method
type enhancedClient struct {
	*gsclient.Client
}

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := enhancedClient{
		gsclient.NewClient(config),
	}
	log.Info("gridscale client configured")

	log.Info("Create server: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	serverCreateRequest := gsclient.ServerCreateRequest{
		Name:   "go-client-server",
		Memory: 1,
		Cores:  2,
	}
	cServer, err := client.CreateServer(emptyCtx, serverCreateRequest)
	if err != nil {
		log.Fatal("Create server has failed with error", err)
	}
	log.WithFields(log.Fields{
		"server_uuid": cServer.ObjectUUID,
	}).Info("Server successfully created")
	defer client.deleteService(serverType, cServer.ObjectUUID)

	//get a server to interact with
	server, err := client.GetServer(emptyCtx, cServer.ObjectUUID)
	if err != nil {
		log.Error("Get server has failed with error", err)
		return
	}

	log.Info("Start server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Turn on server
	err = client.StartServer(emptyCtx, server.Properties.ObjectUUID)
	if err != nil {
		log.Error("Start server has failed with error", err)
		return
	}
	log.Info("Server successfully started")

	log.Info("Stop server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Try to shutdown Server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	err = client.ShutdownServer(shutdownCtx, server.Properties.ObjectUUID)
	//if error is returned and it is not caused by an expired context, returns error
	if err != nil && err != shutdownCtx.Err() {
		log.Error("ShutdownServer has failed with error", err)
		return
	}
	// if the server cannot be shutdown gracefully, try to turn it off
	if err == shutdownCtx.Err() {
		//force the sever to stop
		err = client.StopServer(emptyCtx, server.Properties.ObjectUUID)
		if err != nil {
			return
		}
	}
	log.Info("Server successfully stop")

	log.Info("Update server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	autoRecovery := false
	err = client.UpdateServer(
		emptyCtx,
		server.Properties.ObjectUUID,
		gsclient.ServerUpdateRequest{
			Name:         "updated server",
			Memory:       1,
			Cores:        1,
			AutoRecovery: &autoRecovery,
		})
	if err != nil {
		log.Error("Update server has failed with error", err)
		return
	}
	log.Info("Server successfully updated")

	//Get events of server
	events, err := client.GetServerEventList(emptyCtx, server.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("Events successfully retrieved")

	//Create storage, network, IP, and ISO-image to attach to the server
	log.Info("Create storage, Network, IP, ISO-image: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cStorage, err := client.CreateStorage(
		emptyCtx,
		gsclient.StorageCreateRequest{
			Capacity: 1,
			Name:     "go-client-storage",
		})
	if err != nil {
		log.Error("Create storage has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUUID,
	}).Info("Storage successfully created")
	defer client.deleteService(storageType, cStorage.ObjectUUID)

	cNetwork, err := client.CreateNetwork(
		emptyCtx,
		gsclient.NetworkCreateRequest{
			Name: "go-client-network",
		})
	if err != nil {
		log.Error("Create network has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"network_uuid": cNetwork.ObjectUUID,
	}).Info("Network successfully created")
	defer client.deleteService(networkType, cNetwork.ObjectUUID)

	cIP, err := client.CreateIP(
		emptyCtx,
		gsclient.IPCreateRequest{
			Name:   "go-client-ip",
			Family: gsclient.IPv4Type,
		})
	if err != nil {
		log.Error("Create IP has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"IP_uuid": cIP.ObjectUUID,
	}).Info("IP successfully created")
	defer client.deleteService(ipType, cIP.ObjectUUID)

	cISOimage, err := client.CreateISOImage(
		emptyCtx,
		gsclient.ISOImageCreateRequest{
			Name:      "go-client-iso",
			SourceURL: "http://tinycorelinux.net/10.x/x86/release/TinyCore-current.iso",
		})
	if err != nil {
		log.Error("Create ISO-image has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"isoimage_uuid": cISOimage.ObjectUUID,
	}).Info("ISO-image successfully created")
	defer client.deleteService(isoImageType, cISOimage.ObjectUUID)

	//Attach storage, network, IP, and ISO-image to a server
	err = client.LinkStorage(emptyCtx, server.Properties.ObjectUUID, cStorage.ObjectUUID, false)
	if err != nil {
		log.Error("Link storage has failed with error", err)
		return
	}
	log.Info("Storage successfully attached")
	defer client.unlinkService(storageType, server.Properties.ObjectUUID, cStorage.ObjectUUID)

	err = client.LinkNetwork(
		emptyCtx,
		server.Properties.ObjectUUID,
		cNetwork.ObjectUUID,
		webServerFirewallTemplateUUID,
		false,
		1,
		nil,
		&gsclient.FirewallRules{},
	)
	if err != nil {
		log.Error("Link network has failed with error", err)
		return
	}
	log.Info("Network successfully linked")
	defer client.unlinkService(networkType, server.Properties.ObjectUUID, cNetwork.ObjectUUID)

	err = client.LinkIP(emptyCtx, server.Properties.ObjectUUID, cIP.ObjectUUID)
	if err != nil {
		log.Error("Link IP has failed with error", err)
		return
	}
	log.Info("IP successfully linked")
	defer client.unlinkService(ipType, server.Properties.ObjectUUID, cIP.ObjectUUID)

	err = client.LinkIsoImage(emptyCtx, server.Properties.ObjectUUID, cISOimage.ObjectUUID)
	if err != nil {
		log.Error("Link ISO-image has failed with error", err)
		return
	}
	log.Info("ISO-image successfully linked")
	defer client.unlinkService(isoImageType, server.Properties.ObjectUUID, cISOimage.ObjectUUID)

	log.Info("Unlink and delete: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *enhancedClient) deleteService(serviceType serviceType, id string) {
	switch serviceType {
	case serverType:
		//turn off server before deleting
		err := c.StopServer(emptyCtx, id)
		if err != nil {
			log.Error("Stop server has failed with error", err)
			return
		}
		err = c.DeleteServer(emptyCtx, id)
		if err != nil {
			log.Error("Delete server has failed with error", err)
			return
		}
		log.Info("Server successfully deleted")

		log.Info("Get deleted servers: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		servers, err := c.GetDeletedServers(emptyCtx)
		if err != nil {
			log.Error("Get deleted servers has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"servers": servers,
		}).Info("Retrieved deleted servers successfully")
	case storageType:
		err := c.DeleteStorage(emptyCtx, id)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	case networkType:
		err := c.DeleteNetwork(emptyCtx, id)
		if err != nil {
			log.Error("Delete network has failed with error", err)
			return
		}
		log.Info("Network successfully deleted")
	case ipType:
		err := c.DeleteIP(emptyCtx, id)
		if err != nil {
			log.Error("Delete IP has failed with error", err)
			return
		}
		log.Info("IP successfully deleted")
	case isoImageType:
		err := c.DeleteISOImage(emptyCtx, id)
		if err != nil {
			log.Error("Delete ISO-image has failed with error", err)
			return
		}
		log.Info("ISO-image successfully deleted")
	default:
		log.Error("Unknown service type")
		return
	}
}

func (c *enhancedClient) unlinkService(serviceType serviceType, serverID, serviceID string) {
	switch serviceType {
	case storageType:
		err := c.UnlinkStorage(emptyCtx, serverID, serviceID)
		if err != nil {
			log.Error("Unlink storage has failed with error", err)
			return
		}
		log.Info("Storage successfully unlinked")
	case networkType:
		err := c.UnlinkNetwork(emptyCtx, serverID, serviceID)
		if err != nil {
			log.Error("Unlink network has failed with error", err)
			return
		}
		log.Info("Network successfully unlinked")
	case ipType:
		err := c.UnlinkIP(emptyCtx, serverID, serviceID)
		if err != nil {
			log.Error("Unlink IP has failed with error", err)
			return
		}
		log.Info("IP successfully unlinked")
	case isoImageType:
		err := c.UnlinkIsoImage(emptyCtx, serverID, serviceID)
		if err != nil {
			log.Error("Unlink ISO-image has failed with error", err)
			return
		}
		log.Info("ISO-image successfully unlinked")
	default:
		log.Error("Unknown service type")
		return
	}
}
