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

	log.Info("Create storage and storage backup schedule: Press 'Enter' to continue...")
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

	//Create a storage backup schedule
	cBackupSchedule, err := client.CreateStorageBackupSchedule(
		emptyCtx,
		cStorage.ObjectUUID,
		gsclient.StorageBackupScheduleCreateRequest{
			Name:        "go-client-backup-schedule",
			RunInterval: 120,
			KeepBackups: 2,
		})
	if err != nil {
		log.Error("Create storage backup schedule has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"backup_schedule_uuid": cStorage.ObjectUUID,
	}).Info("Backup schedule successfully created")
	defer func() {
		err := client.DeleteStorageBackupSchedule(emptyCtx, cStorage.ObjectUUID, cBackupSchedule.ObjectUUID)
		if err != nil {
			log.Error("Delete backup schedule has failed with error", err)
			return
		}
		log.Info("Backup schedule successfully deleted")
	}()

	//Get a backup schedule to update
	backupSchedule, err := client.GetStorageBackupSchedule(emptyCtx, cStorage.ObjectUUID, cBackupSchedule.ObjectUUID)
	if err != nil {
		log.Error("Get backup schedule has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"backup_schedule_uuid": backupSchedule.Properties.ObjectUUID,
	}).Info("Backup schedule successfully retrieved")

	log.Info("Update backup schedule: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Update a backup schedule
	err = client.UpdateStorageBackupSchedule(
		emptyCtx,
		cStorage.ObjectUUID,
		backupSchedule.Properties.ObjectUUID,
		gsclient.StorageBackupScheduleUpdateRequest{
			Name:        "updated backup schedule",
			RunInterval: 60,
			KeepBackups: backupSchedule.Properties.KeepBackups,
		})
	if err != nil {
		log.Error("Update backup schedule has failed with error", err)
		return
	}
	log.Info("Backup schedule successfully updated")

	log.Info("Get a list of backups: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	backupList, err := client.GetStorageBackupList(emptyCtx, cStorage.ObjectUUID)
	if err != nil {
		log.Error("Get backups has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"backups": backupList,
	}).Info("Backups successfully retrieved")

	log.Info("Delete  storage backup schedules and storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
