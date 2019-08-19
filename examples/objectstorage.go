package main

import (
	"bufio"
	"os"

	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
)

const LocationUuid = "45ed677b-3702-4b36-be2a-a2eab9827950"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration(
		"https://api.gridscale.io",
		uuid,
		token,
		true,
	)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create object storage access key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	cobj, err := client.CreateObjectStorageAccessKey()
	if err != nil {
		log.Error("Create object storage access key has failed with error", err)
	}
	log.WithFields(log.Fields{
		"accesskey_uuid": cobj.AccessKey,
	}).Info("Create access key successfully")

	log.Info("Get object storage access key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	key, err := client.GetObjectStorageAccessKey(cobj.AccessKey.AccessKey)
	if err != nil {
		log.Error("Retrieve object storage access key has failed with error", err)
	}
	log.WithFields(log.Fields{
		"accesskey": key,
	}).Info("Retrieve access key successfully")

	log.Info("Get buckets: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	buckets, err := client.GetObjectStorageBucketList()
	if err != nil {
		log.Error("Retrieve buckets has failed with error", err)
	}
	log.WithFields(log.Fields{
		"buckets": buckets,
	}).Info("Retrieve buckets successfully")

	log.Info("Delete object storage access key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Delete access key
	err = client.DeleteObjectStorageAccessKey(key.Properties.AccessKey)
	if err != nil {
		log.Error("Delete access key has failed with error", err)
	}
	log.Info("Access key successfully deleted")
}
