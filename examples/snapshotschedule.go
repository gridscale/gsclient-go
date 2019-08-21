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

	log.Info("Create storage and snapshot schedule: Press 'Enter' to continue...")
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
		time.Sleep(30 * time.Second)
		//Delete all snapshots has been made so far
		snapshots, err := client.GetStorageSnapshotList(cStorage.ObjectUuid)
		if err != nil {
			log.Error("Get storage's snapshots has failed with error", err)
			return
		}
		for _, snapshot := range snapshots {
			err = client.DeleteStorageSnapshot(cStorage.ObjectUuid, snapshot.Properties.ObjectUuid)
			if err != nil {
				log.Error("Delete storage's snapshot has failed with error", err)
				return
			}
		}
		//we have to wait for the snapshot getting deleted firstly
		time.Sleep(30 * time.Second)
		err = client.DeleteStorage(cStorage.ObjectUuid)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	//Create Snapshot Schedule
	cSnapshotSchedule, err := client.CreateStorageSnapshotSchedule(cStorage.ObjectUuid, gsclient.StorageSnapshotScheduleCreateRequest{
		Name:          "go-client-snapshot-schedule",
		RunInterval:   120,
		KeepSnapshots: 2,
	})
	if err != nil {
		log.Error("Create snapshot schedule has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshotschedule_uuid": cSnapshotSchedule.ObjectUuid,
	}).Info("Snapshot schedule successfully created")
	defer func() {
		err := client.DeleteStorageSnapshotSchedule(cStorage.ObjectUuid, cSnapshotSchedule.ObjectUuid)
		if err != nil {
			log.Error("Delete snapshot schedule has failed with error", err)
			return
		}
		log.Info("Snapshot schedule successfully deleted")
	}()

	//Get snapshot schedule to update
	snapshotSchedule, err := client.GetStorageSnapshotSchedule(cStorage.ObjectUuid, cSnapshotSchedule.ObjectUuid)
	if err != nil {
		log.Error("Get snapshot schedule has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"snapshotchedule_uuid": snapshotSchedule.Properties.ObjectUuid,
	}).Info("Snapshot schedule successfully retrieved")

	log.Info("Update snapshot schedule: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	err = client.UpdateStorageSnapshotSchedule(cStorage.ObjectUuid, snapshotSchedule.Properties.ObjectUuid, gsclient.StorageSnapshotScheduleUpdateRequest{
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
