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

	log.Info("Create storage and snapshot schedule: Press 'Enter' to continue...")
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
		//Delete all snapshots has been made so far
		snapshots, err := client.GetStorageSnapshotList(context.Background(), cStorage.ObjectUUID)
		if err != nil {
			log.Error("Get storage's snapshots has failed with error", err)
			return
		}
		for _, snapshot := range snapshots {
			err = client.DeleteStorageSnapshot(context.Background(), cStorage.ObjectUUID, snapshot.Properties.ObjectUUID)
			if err != nil {
				log.Error("Delete storage's snapshot has failed with error", err)
				return
			}
		}
		//we have to wait for the snapshot getting deleted firstly
		err = client.DeleteStorage(context.Background(), cStorage.ObjectUUID)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	//Create Snapshot Schedule
	cSnapshotSchedule, err := client.CreateStorageSnapshotSchedule(
		context.Background(),
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
		err := client.DeleteStorageSnapshotSchedule(context.Background(), cStorage.ObjectUUID, cSnapshotSchedule.ObjectUUID)
		if err != nil {
			log.Error("Delete snapshot schedule has failed with error", err)
			return
		}
		log.Info("Snapshot schedule successfully deleted")
	}()

	//Get snapshot schedule to update
	snapshotSchedule, err := client.GetStorageSnapshotSchedule(context.Background(), cStorage.ObjectUUID, cSnapshotSchedule.ObjectUUID)
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
		context.Background(),
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
