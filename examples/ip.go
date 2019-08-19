package main

import (
	"bufio"
	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
	"os"
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

	log.Info("Create IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	ipRequest := gsclient.IpCreateRequest{
		Name:         "go-client-ip",
		Family:       4,
		LocationUuid: LocationUuid,
	}
	//Create new IP
	ipc, err := client.CreateIp(ipRequest)
	if err != nil {
		log.Fatal("Create IP address has failed with error", err)
	}
	log.WithFields(log.Fields{"ip_uuid": ipc.ObjectUuid}).Info("IP address successfully created")

	log.Info("Update IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get IP to update
	ip, err := client.GetIp(ipc.ObjectUuid)
	if err != nil {
		log.Fatal("Get IP address has failed with error", err)
	}
	updateRequest := gsclient.IpUpdateRequest{
		Name:       "Updated IP address",
		Failover:   ip.Properties.Failover,
		ReverseDns: ip.Properties.ReverseDns,
		Labels:     ip.Properties.Labels,
	}
	err = client.UpdateIp(ip.Properties.ObjectUuid, updateRequest)
	if err != nil {
		log.Fatal("Update IP address has failed with error", err)
	}
	log.Info("Retrive IP address events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get IP address events
	response, err := client.GetIpEventList(ip.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Get IP address events has failed with error", err)
	}
	log.WithFields(log.Fields{
		"ip_uuid": ip.Properties.ObjectUuid,
		"events":  response,
	}).Info("Events successfully events retrieved")
	log.Info("Delete IP address: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	err = client.DeleteIp(ip.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Delete IP address has failed with error", err)
	}
	log.Info("Delete IP address successfully")
}
