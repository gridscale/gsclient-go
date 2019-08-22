package main

import (
	"bufio"
	"os"
	"time"

	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
)

const LocationUUID = "45ed677b-3702-4b36-be2a-a2eab9827950"

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

	log.Info("Create PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get template for creating paas
	paasTemplates, err := client.GetPaaSTemplateList()
	if err != nil {
		log.Error("Get PaaS templates has failed with error", err)
		return
	}

	//Create security zone
	secZoneRequest := gsclient.PaaSSecurityZoneCreateRequest{
		Name:         "go-client-security-zone",
		LocationUUID: LocationUUID,
	}
	cSCZ, err := client.CreatePaaSSecurityZone(secZoneRequest)
	if err != nil {
		log.Error("Create security zone has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"securityzone_uuid": cSCZ.ObjectUUID,
	}).Info("Security zone successfully created")
	defer func() {
		//Wait until paas deleted successfully
		//it takes around a minute
		time.Sleep(60 * time.Second)
		err := client.DeletePaaSSecurityZone(cSCZ.ObjectUUID)
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
	cPaaS, err := client.CreatePaaSService(paasRequest)
	if err != nil {
		log.Error("Create PaaS service has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"paas_uuid": cPaaS.ObjectUUID,
	}).Info("PaaS service create successfully")
	defer func() {
		err := client.DeletePaaSService(cPaaS.ObjectUUID)
		if err != nil {
			log.Error("Delete PaaS service has failed with error", err)
			return
		}
		log.Info("PaaS service successfully deleted")
	}()

	log.Info("Update PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get a security zone to update
	secZone, err := client.GetPaaSSecurityZone(cSCZ.ObjectUUID)
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
	err = client.UpdatePaaSSecurityZone(secZone.Properties.ObjectUUID, secZoneUpdateRequest)
	if err != nil {
		log.Error("Update security zone has failed with error", err)
		return
	}
	log.Info("Security Zone successfully updated")

	//Get a PaaS service to update
	paas, err := client.GetPaaSService(cPaaS.ObjectUUID)
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
	err = client.UpdatePaaSService(paas.Properties.ObjectUUID, paasUpdateRequest)
	if err != nil {
		log.Error("Update PaaS service has failed with error", err)
		return
	}
	log.Info("PaaS service successfully updated")

	//Clean up
	log.Info("Delete PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
