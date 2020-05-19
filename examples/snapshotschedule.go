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

	log.Info("Create storage and snapshot schedule: Press 'Enter' to continue...")
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
		//Delete all snapshots has been made so far
		snapshots, err := client.GetStorageSnapshotList(emptyCtx, cStorage.ObjectUUID)
		if err != nil {
			log.Error("Get storage's snapshots has failed with error", err)
			return
		}
		for _, snapshot := range snapshots {
			err = client.DeleteStorageSnapshot(emptyCtx, cStorage.ObjectUUID, snapshot.Properties.ObjectUUID)
			if err != nil {
				log.Error("Delete storage's snapshot has failed with error", err)
				return
			}
		}
		//we have to wait for the snapshot getting deleted firstly
		err = client.DeleteStorage(emptyCtx, cStorage.ObjectUUID)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	//Create Snapshot Schedule
	cSnapshotSchedule, err := client.CreateStorageSnapshotSchedule(
		emptyCtx,
		cStorage.ObjectUUID,
		gsclient.StorageSnapshotScheduleCreateRequest{
			Name:          "go-client-snapshot-schedule",
			RunInterval:   120,
			KeepSnapshots: 2,
		})
	if err != nil {
		log.Error("Create snapshot schedule has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshotschedule_uuid": cSnapshotSchedule.ObjectUUID,
	}).Info("Snapshot schedule successfully created")
	defer func() {
		err := client.DeleteStorageSnapshotSchedule(emptyCtx, cStorage.ObjectUUID, cSnapshotSchedule.ObjectUUID)
		if err != nil {
			log.Error("Delete snapshot schedule has failed with error", err)
			return
		}
		log.Info("Snapshot schedule successfully deleted")
	}()

	//Get snapshot schedule to update
	snapshotSchedule, err := client.GetStorageSnapshotSchedule(emptyCtx, cStorage.ObjectUUID, cSnapshotSchedule.ObjectUUID)
	if err != nil {
		log.Error("Get snapshot schedule has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshotchedule_uuid": snapshotSchedule.Properties.ObjectUUID,
	}).Info("Snapshot schedule successfully retrieved")

	log.Info("Update snapshot schedule: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	err = client.UpdateStorageSnapshotSchedule(
		emptyCtx,
		cStorage.ObjectUUID,
		snapshotSchedule.Properties.ObjectUUID,
		gsclient.StorageSnapshotScheduleUpdateRequest{
			Name:          "updated snapshot schedule",
			RunInterval:   snapshotSchedule.Properties.RunInterval,
			KeepSnapshots: snapshotSchedule.Properties.KeepSnapshots,
		})
	if err != nil {
		log.Error("Update snapshot schedule has failed with error", err)
		return
	}
	log.Info("Snapshot schedule successfully updated")

	log.Info("Delete snapshot schedule and storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
