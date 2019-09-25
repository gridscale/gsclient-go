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

	log.Info("Create PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get template for creating paas
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		log.Error("Get PaaS templates has failed with error", err)
		return
	}

	//Create security zone
	secZoneRequest := gsclient.PaaSSecurityZoneCreateRequest{
		Name:         "go-client-security-zone",
		LocationUUID: locationUUID,
	}
	cSCZ, err := client.CreatePaaSSecurityZone(context.Background(), secZoneRequest)
	if err != nil {
		log.Error("Create security zone has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"securityzone_uuid": cSCZ.ObjectUUID,
	}).Info("Security zone successfully created")
	defer func() {
		err := client.DeletePaaSSecurityZone(context.Background(), cSCZ.ObjectUUID)
		if err != nil {
			log.Error("Delete security zone has failed with error", err)
			return
		}
		log.Info("Security zone successfully deleted")
	}()

	//Create PaaS service
	paasRequest := gsclient.PaaSServiceCreateRequest{
		Name:                    "go-client-paas",
		PaaSServiceTemplateUUID: paasTemplates[0].Properties.ObjectUUID,
		PaaSSecurityZoneUUID:    cSCZ.ObjectUUID,
	}
	cPaaS, err := client.CreatePaaSService(context.Background(), paasRequest)
	if err != nil {
		log.Error("Create PaaS service has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"paas_uuid": cPaaS.ObjectUUID,
	}).Info("PaaS service create successfully")
	defer func() {
		err := client.DeletePaaSService(context.Background(), cPaaS.ObjectUUID)
		if err != nil {
			log.Error("Delete PaaS service has failed with error", err)
			return
		}
		log.Info("PaaS service successfully deleted")

		log.Info("Get deleted PaaS services: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		paasServices, err := client.GetDeletedPaaSServices(context.Background())
		if err != nil {
			log.Error("Get deleted PaaS services has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"PaaS services": paasServices,
		}).Info("Retrieved deleted PaaS services successfully")
	}()

	log.Info("Update PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get a security zone to update
	secZone, err := client.GetPaaSSecurityZone(context.Background(), cSCZ.ObjectUUID)
	if err != nil {
		log.Error("Get security zone has failed with error", err)
		return
	}
	secZoneUpdateRequest := gsclient.PaaSSecurityZoneUpdateRequest{
		Name:                 "updated security zone",
		LocationUUID:         secZone.Properties.LocationUUID,
		PaaSSecurityZoneUUID: secZone.Properties.ObjectUUID,
	}
	//Update security zone
	err = client.UpdatePaaSSecurityZone(context.Background(), secZone.Properties.ObjectUUID, secZoneUpdateRequest)
	if err != nil {
		log.Error("Update security zone has failed with error", err)
		return
	}
	log.Info("Security Zone successfully updated")

	//Get a PaaS service to update
	paas, err := client.GetPaaSService(context.Background(), cPaaS.ObjectUUID)
	if err != nil {
		log.Error("Get PaaS service has failed with error", err)
		return
	}

	//Update PaaS service
	paasUpdateRequest := gsclient.PaaSServiceUpdateRequest{
		Name:           "updated paas",
		Labels:         paas.Properties.Labels,
		Parameters:     paas.Properties.Parameters,
		ResourceLimits: paas.Properties.ResourceLimits,
	}
	err = client.UpdatePaaSService(context.Background(), paas.Properties.ObjectUUID, paasUpdateRequest)
	if err != nil {
		log.Error("Update PaaS service has failed with error", err)
		return
	}
	log.Info("PaaS service successfully updated")

	//Clean up
	log.Info("Delete PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
