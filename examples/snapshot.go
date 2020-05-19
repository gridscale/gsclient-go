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

	log.Info("Create storage and snapshot: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Create storage
	cStorage, err := client.CreateStorage(
		emptyCtx,
		gsclient.StorageCreateRequest{
			Capacity: 1,
			Name:     "go-client-storage",
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
	}()

	//Create a snapshot
	cSnapshot, err := client.CreateStorageSnapshot(
		emptyCtx,
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
		err := client.DeleteStorageSnapshot(emptyCtx, cStorage.ObjectUUID, cSnapshot.ObjectUUID)
		if err != nil {
			log.Error("Delete storage snapshot has failed with error", err)
			return
		}
		log.Info("Storage snapshot successfully deleted")

		log.Info("Get deleted snapshots: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		snapshots, err := client.GetDeletedSnapshots(emptyCtx)
		if err != nil {
			log.Error("Get deleted snapshots has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"snapshots": snapshots,
		}).Info("Retrieved deleted snapshots successfully")
	}()

	//Get a snapshot to update
	snapshot, err := client.GetStorageSnapshot(emptyCtx, cStorage.ObjectUUID, cSnapshot.ObjectUUID)
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
		emptyCtx,
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
		emptyCtx,
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
