package main

import (
	"bufio"
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
	)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	ipRequest := gsclient.IPCreateRequest{
		Name:         "go-client-ip",
		Family:       4,
		LocationUUID: locationUUID,
	}
	//Create new IP
	ipc, err := client.CreateIP(ipRequest)
	if err != nil {
		log.Error("Create IP address has failed with error", err)
		return
	}
	log.WithFields(log.Fields{"ip_uuid": ipc.ObjectUUID}).Info("IP address successfully created")
	defer func() {
		err := client.DeleteIP(ipc.ObjectUUID)
		if err != nil {
			log.Error("Delete IP address has failed with error", err)
			return
		}
		log.Info("Delete IP address successfully")
	}()

	log.Info("Update IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get IP to update
	ip, err := client.GetIP(ipc.ObjectUUID)
	if err != nil {
		log.Error("Get IP address has failed with error", err)
		return
	}
	updateRequest := gsclient.IPUpdateRequest{
		Name:       "Updated IP address",
		Failover:   ip.Properties.Failover,
		ReverseDns: ip.Properties.ReverseDns,
		Labels:     ip.Properties.Labels,
	}
	err = client.UpdateIP(ip.Properties.ObjectUUID, updateRequest)
	if err != nil {
		log.Error("Update IP address has failed with error", err)
		return
	}
	log.Info("Retrive IP address events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get IP address events
	response, err := client.GetIPEventList(ip.Properties.ObjectUUID)
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
