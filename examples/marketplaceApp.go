package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gridscale/gsclient-go/v3"
	log "github.com/sirupsen/logrus"
)

var emptyCtx = context.Background()

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create a storage and a snapshot: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// In order to create a marketplace application, we need to create a storage and its snapshot
	// Create storage
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
	defer func() {
		err := client.DeleteStorage(emptyCtx, cStorage.ObjectUUID)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	}()

	// Create storage snapshot
	cSnapshot, err := client.CreateStorageSnapshot(
		emptyCtx,
		cStorage.ObjectUUID,
		gsclient.StorageSnapshotCreateRequest{
			Name: "go-client-snapshot",
		})
	if err != nil {
		log.Error("Create storage snapshot has failed with error", err)
		return
	}
	defer func() {
		err := client.DeleteStorageSnapshot(emptyCtx, cStorage.ObjectUUID, cSnapshot.ObjectUUID)
		if err != nil {
			log.Error("Delete storage snapshot has failed with error", err)
			return
		}
		log.Info("Storage snapshot successfully deleted")
	}()

	// Create new object storage access key
	cobj, err := client.CreateObjectStorageAccessKey(emptyCtx)
	if err != nil {
		log.Error("Create object storage access key has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"accesskey_uuid": cobj.AccessKey,
	}).Info("Create access key successfully")
	defer func() {
		// Delete access key
		err := client.DeleteObjectStorageAccessKey(emptyCtx, cobj.AccessKey.AccessKey)
		if err != nil {
			log.Error("Delete access key has failed with error", err)
			return
		}
		log.Info("Access key successfully deleted")
	}()

	log.Info("Please create a new bucket in your gridscale panel, then type the name of the bucket and press 'Enter' to continue...")
	bucketNameStdin, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Error("Reading bucket name returns error", err)
		return
	}
	bucketName := strings.TrimSuffix(bucketNameStdin, "\n")
	// Export snapshot to S3
	exportReqBody := gsclient.StorageSnapshotExportToS3Request{
		S3auth: gsclient.S3auth{
			Host:      "gos3.io",
			AccessKey: cobj.AccessKey.AccessKey,
			SecretKey: cobj.AccessKey.SecretKey,
		},
		S3data: gsclient.S3data{
			Host:     "https://gos3.io",
			Bucket:   "testT1",
			Filename: "snapshot.gz",
			Private:  true,
		},
	}
	err = client.ExportStorageSnapshotToS3(emptyCtx, cStorage.ObjectUUID, cSnapshot.ObjectUUID, exportReqBody)
	if err != nil {
		log.Error("Exporting snapshot to s3 has failed with error", err)
		return
	}

	// Create marketplace application
	cMartketApp, err := client.CreateMarketplaceApplication(emptyCtx, gsclient.MarketplaceApplicationCreateRequest{
		Name:              "go-client-marketplace-app",
		ObjectStoragePath: fmt.Sprintf("s3://%s/snapshot.gz", bucketName),
		Category:          gsclient.MarketplaceApplicationAdminpanelCategory,
		Setup: gsclient.MarketplaceApplicationSetup{
			Cores:    1,
			Memory:   2,
			Capacity: 5,
		},
	})
	if err != nil {
		log.Error("Create marketplace application has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"marketplace_app_uuid": cMartketApp.ObjectUUID,
	}).Info("Marketplace application successfully created")
	defer func() {
		err := client.DeleteMarketplaceApplication(emptyCtx, cMartketApp.ObjectUUID)
		if err != nil {
			log.Error("Delete marketplace application has failed with error", err)
			return
		}
		log.Info("Marketplace application successfully deleted")
	}()

	// get a marketplace application to update
	marketApp, err := client.GetMarketplaceApplication(emptyCtx, cMartketApp.ObjectUUID)
	if err != nil {
		log.Error("Get marketplace application has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"marketplace_app": marketApp,
	}).Info("marketplace application successfully retrieved")

	log.Info("Update marketplace application: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// Update marketplace application
	err = client.UpdateMarketplaceApplication(emptyCtx, marketApp.Properties.ObjectUUID, gsclient.MarketplaceApplicationUpdateRequest{
		Name: "updated marketplace application",
		Setup: &gsclient.MarketplaceApplicationSetup{
			Cores:    2,
			Memory:   4,
			Capacity: 10,
		},
	})
	if err != nil {
		log.Error("Update marketplace application has failed with error", err)
		return
	}
	log.Info("Marketplace application successfully updated")

	log.Info("Get marketplace application's events: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// Get marketplace application's events
	events, err := client.GetMarketplaceApplicationEventList(emptyCtx, marketApp.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get marketplace application's events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("marketplace application's events successfully retrieved")

	log.Info("Delete marketplace application: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
