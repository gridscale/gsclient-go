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

	log.Info("Create PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Get template for creating paas
	paasTemplates, err := client.GetPaaSTemplateList(emptyCtx)
	if err != nil {
		log.Error("Get PaaS templates has failed with error", err)
		return
	}

	// Create security zone
	secZoneRequest := gsclient.PaaSSecurityZoneCreateRequest{
		Name: "go-client-security-zone",
	}
	cSCZ, err := client.CreatePaaSSecurityZone(emptyCtx, secZoneRequest)
	if err != nil {
		log.Error("Create security zone has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"securityzone_uuid": cSCZ.ObjectUUID,
	}).Info("Security zone successfully created")
	defer func() {
		err := client.DeletePaaSSecurityZone(emptyCtx, cSCZ.ObjectUUID)
		if err != nil {
			log.Error("Delete security zone has failed with error", err)
			return
		}
		log.Info("Security zone successfully deleted")
	}()

	// Create PaaS service
	paasRequest := gsclient.PaaSServiceCreateRequest{
		Name:                    "go-client-paas",
		PaaSServiceTemplateUUID: paasTemplates[0].Properties.ObjectUUID,
		PaaSSecurityZoneUUID:    cSCZ.ObjectUUID,
	}
	cPaaS, err := client.CreatePaaSService(emptyCtx, paasRequest)
	if err != nil {
		log.Error("Create PaaS service has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"paas_uuid": cPaaS.ObjectUUID,
	}).Info("PaaS service create successfully")
	defer func() {
		err := client.DeletePaaSService(emptyCtx, cPaaS.ObjectUUID)
		if err != nil {
			log.Error("Delete PaaS service has failed with error", err)
			return
		}
		log.Info("PaaS service successfully deleted")

		log.Info("Get deleted PaaS services: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		paasServices, err := client.GetDeletedPaaSServices(emptyCtx)
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

	// Get a security zone to update
	secZone, err := client.GetPaaSSecurityZone(emptyCtx, cSCZ.ObjectUUID)
	if err != nil {
		log.Error("Get security zone has failed with error", err)
		return
	}
	secZoneUpdateRequest := gsclient.PaaSSecurityZoneUpdateRequest{
		Name:                 "updated security zone",
		PaaSSecurityZoneUUID: secZone.Properties.ObjectUUID,
	}
	// Update security zone
	err = client.UpdatePaaSSecurityZone(emptyCtx, secZone.Properties.ObjectUUID, secZoneUpdateRequest)
	if err != nil {
		log.Error("Update security zone has failed with error", err)
		return
	}
	log.Info("Security Zone successfully updated")

	// Get a PaaS service to update
	paas, err := client.GetPaaSService(emptyCtx, cPaaS.ObjectUUID)
	if err != nil {
		log.Error("Get PaaS service has failed with error", err)
		return
	}

	// Update PaaS service
	paasUpdateRequest := gsclient.PaaSServiceUpdateRequest{
		Name:           "updated paas",
		Labels:         &paas.Properties.Labels,
		Parameters:     paas.Properties.Parameters,
		ResourceLimits: paas.Properties.ResourceLimits,
	}
	err = client.UpdatePaaSService(emptyCtx, paas.Properties.ObjectUUID, paasUpdateRequest)
	if err != nil {
		log.Error("Update PaaS service has failed with error", err)
		return
	}
	log.Info("PaaS service successfully updated")

	// Clean up
	log.Info("Delete PaaS and Security zone: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
