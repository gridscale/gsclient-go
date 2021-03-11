package main

import (
	"bufio"
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gridscale/gsclient-go/v3"
)

var emptyCtx = context.Background()

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create storage: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Create a storage
	cStorage, err := client.CreateStorage(
		emptyCtx,
		gsclient.StorageCreateRequest{
			Capacity:    1,
			Name:        "go-client-storage",
			StorageType: gsclient.InsaneStorageType,
		})
	if err != nil {
		log.Error("Create storage has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUUID,
	}).Info("Storage successfully created")
	defer func() {
		err := client.DeleteStorage(emptyCtx, cStorage.ObjectUUID)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")

		log.Info("Get deleted storages: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		storages, err := client.GetDeletedStorages(emptyCtx)
		if err != nil {
			log.Error("Get deleted storages has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"storages": storages,
		}).Info("Retrieved deleted storages successfully")
	}()

	// Get storage to update
	storage, err := client.GetStorage(emptyCtx, cStorage.ObjectUUID)
	if err != nil {
		log.Error("Get storage has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"storage_uuid": storage.Properties.ObjectUUID,
	}).Info("Storage successfully retrieved")

	log.Info("Update storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	err = client.UpdateStorage(
		emptyCtx,
		storage.Properties.ObjectUUID,
		gsclient.StorageUpdateRequest{
			Name:     "updated storage",
			Labels:   &storage.Properties.Labels,
			Capacity: storage.Properties.Capacity,
		})
	if err != nil {
		log.Error("Update storage has failed with error", err)
		return
	}
	log.Info("Storage successfully updated")

	log.Info("Get storage's events: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	events, err := client.GetStorageEventList(emptyCtx, storage.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get storage's events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("Storage's events successfully retrieved")

	log.Info("Delete storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
