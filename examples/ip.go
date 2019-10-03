package main

import (
	"bufio"
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/gridscale/gsclient-go"
	"github.com/gridscale/gsclient-go"
	"github.com/gridscale/gsclient-go"
)

const locationUUID = "45ed677b-3702-4b36-be2a-a2eab9827950"

var emptyCtx = context.Background()

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
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
	ipc, err := client.CreateIP(emptyCtx, ipRequest)
	if err != nil {
		log.Error("Create IP address has failed with error", err)
		return
	}
	log.WithFields(log.Fields{"ip_uuid": ipc.ObjectUUID}).Info("IP address successfully created")
	defer func() {
		err := client.DeleteIP(emptyCtx, ipc.ObjectUUID)
		if err != nil {
			log.Error("Delete IP address has failed with error", err)
			return
		}
		log.Info("Delete IP address successfully")

		log.Info("Get deleted IP address: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		ips, err := client.GetDeletedIPs(emptyCtx)
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
	ip, err := client.GetIP(emptyCtx, ipc.ObjectUUID)
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
	err = client.UpdateIP(emptyCtx, ip.Properties.ObjectUUID, updateRequest)
	if err != nil {
		log.Error("Update IP address has failed with error", err)
		return
	}
	log.Info("Retrive IP address events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get IP address events
	response, err := client.GetIPEventList(emptyCtx, ip.Properties.ObjectUUID)
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
