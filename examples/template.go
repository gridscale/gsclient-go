package main

import (
	"bufio"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gridscale/gsclient-go"
)

const LocationUuid = "45ed677b-3702-4b36-be2a-a2eab9827950"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration(
		"https://api.gridscale.io",
		uuid,
		token,
		false,
	)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create template: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//In order to create a template, we need to create a storage and its snapshot
	//Create storage
	cStorage, err := client.CreateStorage(gsclient.StorageCreateRequest{
		Capacity:     0,
		LocationUuid: LocationUuid,
		Name:         "go-client-storage",
	})
	if err != nil {
		log.Fatal("Create storage has failed with error", err)
	}
	defer func() {
		time.Sleep(30 * time.Second)
		err := client.DeleteStorage(cStorage.ObjectUuid)
		if err != nil {
			log.Fatal("Delete storage has failed with error", err)
		}
		log.Info("Storage successfully deleted")
	}()

	//Create storage snapshot
	cSnapshot, err := client.CreateStorageSnapshot(cStorage.ObjectUuid, gsclient.StorageSnapshotCreateRequest{
		Name: "go-client-snapshot",
	})
	if err != nil {
		log.Fatal("Create storage snapshot has failed with error", err)
	}
	defer func() {
		time.Sleep(40 * time.Second)
		err := client.DeleteStorageSnapshot(cStorage.ObjectUuid, cSnapshot.ObjectUuid)
		if err != nil {
			log.Fatal("Delete storage snapshot has failed with error", err)
		}
		log.Info("Storage snapshot successfully deleted")
	}()

	//Create template
	cTemplate, err := client.CreateTemplate(gsclient.TemplateCreateRequest{
		Name:         "go-client-template",
		SnapshotUuid: cSnapshot.ObjectUuid,
	})
	if err != nil {
		log.Fatal("Create template has failed with error", err)
	}
	log.WithFields(log.Fields{
		"template_uuid": cTemplate.ObjectUuid,
	}).Info("Template successfully created")
	defer func() {
		err := client.DeleteTemplate(cTemplate.ObjectUuid)
		if err != nil {
			log.Fatal("Delete template has failed with error", err)
		}
		log.Info("Template successfully deleted")
	}()

	//get a template to update
	template, err := client.GetTemplate(cTemplate.ObjectUuid)
	if err != nil {
		log.Fatal("Get template has failed with error", err)
	}
	log.WithFields(log.Fields{
		"template_uuid": template.Properties.ObjectUuid,
	}).Info("Template successfully retrieved")

	log.Info("Update template: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Update template
	err = client.UpdateTemplate(template.Properties.ObjectUuid, gsclient.TemplateUpdateRequest{
		Name:   "updated template",
		Labels: template.Properties.Labels,
	})
	if err != nil {
		log.Fatal("Update template has failed with error", err)
	}
	log.Info("Template successfully updated")

	log.Info("Get template's events: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Get template's events
	events, err := client.GetTemplateEventList(template.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Get template's events has failed with error", err)
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("Template's events successfully retrieved")

	log.Info("Delete template: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
