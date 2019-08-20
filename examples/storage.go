package main

import (
	"bufio"
	"os"

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
		true,
	)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create storage: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Create a storage
	cStorage, err := client.CreateStorage(gsclient.StorageCreateRequest{
		Capacity:     1,
		LocationUuid: LocationUuid,
		Name:         "go-client-storage",
	})
	if err != nil {
		log.Fatal("Create storage has failed with error", err)
	}
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUuid,
	}).Info("Storage successfully created")
	defer func() {
		err := client.DeleteStorage(cStorage.ObjectUuid)
		if err != nil {
			log.Fatal("Delete storage has failed with error", err)
		}
		log.Info("Storage successfully deleted")
	}()

	//Get storage to update
	storage, err := client.GetStorage(cStorage.ObjectUuid)
	if err != nil {
		log.Fatal("Get storage has failed with error", err)
	}
	log.WithFields(log.Fields{
		"storage_uuid": storage.Properties.ObjectUuid,
	}).Info("Storage successfully retrieved")

	log.Info("Update storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	err = client.UpdateStorage(storage.Properties.ObjectUuid, gsclient.StorageUpdateRequest{
		Name:     "updated storage",
		Labels:   storage.Properties.Labels,
		Capacity: storage.Properties.Capacity,
	})
	if err != nil {
		log.Fatal("Update storage has failed with error", err)
	}
	log.Info("Storage successfully updated")

	log.Info("Get storage's events: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	events, err := client.GetStorageEventList(storage.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Get storage's events has failed with error", err)
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("Storage's events successfully retrieved")

	log.Info("Delete storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
