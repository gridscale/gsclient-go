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
	config := gsclient.NewConfiguration("https://api.gridscale.io", uuid, token, true)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create storage and snapshot: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Create storage
	cStorage, err := client.CreateStorage(gsclient.StorageCreateRequest{
		Capacity:     1,
		LocationUuid: LocationUuid,
		Name:         "go-client-storage",
	})
	if err != nil {
		log.Error("Create storage has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUuid,
	}).Info("Storage successfully created")
	defer func() {
		//we have to wait for the snapshot getting deleted firstly
		time.Sleep(1 * time.Minute)
		err := client.DeleteStorage(cStorage.ObjectUuid)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	//Create a snapshot
	cSnapshot, err := client.CreateStorageSnapshot(cStorage.ObjectUuid, gsclient.StorageSnapshotCreateRequest{
		Name: "go-client-snapshot",
	})
	if err != nil {
		log.Error("Create snapshot has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshot_uuid": cStorage.ObjectUuid,
	}).Info("Snapshot successfully created")
	defer func() {
		err := client.DeleteStorageSnapshot(cStorage.ObjectUuid, cSnapshot.ObjectUuid)
		if err != nil {
			log.Error("Delete storage snapshot has failed with error", err)
			return
		}
		log.Info("Storage snapshot successfully deleted")
	}()

	//Get a snapshot to update
	snapshot, err := client.GetStorageSnapshot(cStorage.ObjectUuid, cSnapshot.ObjectUuid)
	if err != nil {
		log.Error("Get snapshot has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshot_uuid": snapshot.Properties.ObjectUuid,
	}).Info("Snapshot successfully retrieved")

	log.Info("Update snapshot: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Update a snapshot
	err = client.UpdateStorageSnapshot(cStorage.ObjectUuid, snapshot.Properties.ObjectUuid, gsclient.StorageSnapshotUpdateRequest{
		Name: "updated snapshot",
	})
	if err != nil {
		log.Error("Update snapshot has failed with error", err)
		return
	}
	log.Info("Snapshot successfully updated")

	log.Info("Rollback storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Rollback
	err = client.RollbackStorage(cStorage.ObjectUuid, snapshot.Properties.ObjectUuid, gsclient.StorageRollbackRequest{
		Rollback: true,
	})
	if err != nil {
		log.Error("Rollback storage has failed with error", err)
		return
	}
	log.Info("Storage successfully rollbacked")

	log.Info("Delete snapshot and storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
