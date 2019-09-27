package main

import (
	"bufio"
	"context"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/gridscale/gsclient-go"
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

	log.Info("Create template: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//In order to create a template, we need to create a storage and its snapshot
	//Create storage
	cStorage, err := client.CreateStorage(
		context.Background(),
		gsclient.StorageCreateRequest{
			Capacity:     1,
			LocationUUID: locationUUID,
			Name:         "go-client-storage",
		})
	if err != nil {
		log.Error("Create storage has failed with error", err)
		return
	}
	defer func() {
		err := client.DeleteStorage(context.Background(), cStorage.ObjectUUID)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	//Create storage snapshot
	cSnapshot, err := client.CreateStorageSnapshot(
		context.Background(),
		cStorage.ObjectUUID,
		gsclient.StorageSnapshotCreateRequest{
			Name: "go-client-snapshot",
		})
	if err != nil {
		log.Error("Create storage snapshot has failed with error", err)
		return
	}
	defer func() {
		err := client.DeleteStorageSnapshot(context.Background(), cStorage.ObjectUUID, cSnapshot.ObjectUUID)
		if err != nil {
			log.Error("Delete storage snapshot has failed with error", err)
			return
		}
		log.Info("Storage snapshot successfully deleted")
	}()

	//Create template
	cTemplate, err := client.CreateTemplate(
		context.Background(),
		gsclient.TemplateCreateRequest{
			Name:         "go-client-template",
			SnapshotUUID: cSnapshot.ObjectUUID,
		})
	if err != nil {
		log.Error("Create template has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"template_uuid": cTemplate.ObjectUUID,
	}).Info("Template successfully created")
	defer func() {
		err := client.DeleteTemplate(context.Background(), cTemplate.ObjectUUID)
		if err != nil {
			log.Error("Delete template has failed with error", err)
			return
		}
		log.Info("Template successfully deleted")

		log.Info("Get deleted templates: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		templates, err := client.GetDeletedTemplates(context.Background())
		if err != nil {
			log.Error("Get deleted templates has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"templates": templates,
		}).Info("Retrieved deleted templates successfully")
	}()

	//get a template to update
	template, err := client.GetTemplate(context.Background(), cTemplate.ObjectUUID)
	if err != nil {
		log.Error("Get template has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"template_uuid": template.Properties.ObjectUUID,
	}).Info("Template successfully retrieved")

	log.Info("Update template: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Update template
	err = client.UpdateTemplate(
		context.Background(),
		template.Properties.ObjectUUID,
		gsclient.TemplateUpdateRequest{
			Name:   "updated template",
			Labels: template.Properties.Labels,
		})
	if err != nil {
		log.Error("Update template has failed with error", err)
		return
	}
	log.Info("Template successfully updated")

	log.Info("Get template's events: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Get template's events
	events, err := client.GetTemplateEventList(context.Background(), template.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get template's events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("Template's events successfully retrieved")

	log.Info("Delete template: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
