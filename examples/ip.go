package main

import (
	"bufio"
	"context"
	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
	"os"
)

const locationUUID = "45ed677b-3702-4b36-be2a-a2eab9827950"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration(
		"https://api.gridscale.io",
		uuid,
		token,
		true,
		true,
		0,
		0,
		0,
	)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	ipRequest := gsclient.IPCreateRequest{
		Name:         "go-client-ip",
		Family:       gsclient.IPv4Type,
		LocationUUID: locationUUID,
	}
	//Create new IP
	ipc, err := client.CreateIP(context.Background(), ipRequest)
	if err != nil {
		log.Error("Create IP address has failed with error", err)
		return
	}
	log.WithFields(log.Fields{"ip_uuid": ipc.ObjectUUID}).Info("IP address successfully created")
	defer func() {
		err := client.DeleteIP(context.Background(), ipc.ObjectUUID)
		if err != nil {
			log.Error("Delete IP address has failed with error", err)
			return
		}
		log.Info("Delete IP address successfully")

		log.Info("Get deleted IP address: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		ips, err := client.GetDeletedIPs(context.Background())
		if err != nil {
			log.Error("Get delete IP address has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"ips": ips,
		}).Info("Retrieved deleted IP successfully")
	}()

	log.Info("Update IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get IP to update
	ip, err := client.GetIP(context.Background(), ipc.ObjectUUID)
	if err != nil {
		log.Error("Get IP address has failed with error", err)
		return
	}
	updateRequest := gsclient.IPUpdateRequest{
		Name:       "Updated IP address",
		Failover:   ip.Properties.Failover,
		ReverseDNS: ip.Properties.ReverseDNS,
		Labels:     ip.Properties.Labels,
	}
	err = client.UpdateIP(context.Background(), ip.Properties.ObjectUUID, updateRequest)
	if err != nil {
		log.Error("Update IP address has failed with error", err)
		return
	}
	log.Info("Retrive IP address events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get IP address events
	response, err := client.GetIPEventList(context.Background(), ip.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get IP address events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"ip_uuid": ip.Properties.ObjectUUID,
		"events":  response,
	}).Info("Events successfully events retrieved")
	log.Info("Delete IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
