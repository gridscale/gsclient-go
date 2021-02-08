package main

import (
	"bufio"
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gridscale/gsclient-go/v3"
)

var ctx = context.Background()

const backupID = "this-should-be-a-backup-UUID"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create storage from backup: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	cStorage, err := client.CreateStorageFromBackup(ctx, backupID, "My new storage")
	if err != nil {
		log.Error("CreateStorageFromBackup failed with", err)
		return
	}
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUUID,
	}).Info("Storage successfully created")
	defer func() {
		err := client.DeleteStorage(ctx, cStorage.ObjectUUID)
		if err != nil {
			log.Error("Delete storage failed with", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	log.Info("Delete storage: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
