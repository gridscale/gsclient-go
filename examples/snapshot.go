package main

import (
	"bufio"
	"context"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/nvthongswansea/gsclient-go"
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

	log.Info("Create storage and snapshot: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
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
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUUID,
	}).Info("Storage successfully created")
	defer func() {
		err := client.DeleteStorage(context.Background(), cStorage.ObjectUUID)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	//Create a snapshot
	cSnapshot, err := client.CreateStorageSnapshot(
		context.Background(),
		cStorage.ObjectUUID,
		gsclient.StorageSnapshotCreateRequest{
			Name: "go-client-snapshot",
		})
	if err != nil {
		log.Error("Create snapshot has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshot_uuid": cStorage.ObjectUUID,
	}).Info("Snapshot successfully created")
	defer func() {
		err := client.DeleteStorageSnapshot(context.Background(), cStorage.ObjectUUID, cSnapshot.ObjectUUID)
		if err != nil {
			log.Error("Delete storage snapshot has failed with error", err)
			return
		}
		log.Info("Storage snapshot successfully deleted")

		log.Info("Get deleted snapshots: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		snapshots, err := client.GetDeletedSnapshots(context.Background())
		if err != nil {
			log.Error("Get deleted snapshots has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"snapshots": snapshots,
		}).Info("Retrieved deleted snapshots successfully")
	}()

	//Get a snapshot to update
	snapshot, err := client.GetStorageSnapshot(context.Background(), cStorage.ObjectUUID, cSnapshot.ObjectUUID)
	if err != nil {
		log.Error("Get snapshot has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshot_uuid": snapshot.Properties.ObjectUUID,
	}).Info("Snapshot successfully retrieved")

	log.Info("Update snapshot: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Update a snapshot
	err = client.UpdateStorageSnapshot(
		context.Background(),
		cStorage.ObjectUUID,
		snapshot.Properties.ObjectUUID,
		gsclient.StorageSnapshotUpdateRequest{
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
	err = client.RollbackStorage(
		context.Background(),
		cStorage.ObjectUUID,
		snapshot.Properties.ObjectUUID,
		gsclient.StorageRollbackRequest{
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
