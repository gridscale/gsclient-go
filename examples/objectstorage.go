package main

import (
	"bufio"
	"context"
	"os"

	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
)

var emptyCtx = context.Background()

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create object storage access key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	cobj, err := client.CreateObjectStorageAccessKey(emptyCtx)
	if err != nil {
		log.Error("Create object storage access key has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"accesskey_uuid": cobj.AccessKey,
	}).Info("Create access key successfully")
	defer func() {
		//Delete access key
		err := client.DeleteObjectStorageAccessKey(emptyCtx, cobj.AccessKey.AccessKey)
		if err != nil {
			log.Error("Delete access key has failed with error", err)
			return
		}
		log.Info("Access key successfully deleted")
	}()

	log.Info("Get object storage access key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	key, err := client.GetObjectStorageAccessKey(emptyCtx, cobj.AccessKey.AccessKey)
	if err != nil {
		log.Error("Retrieve object storage access key has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"accesskey": key,
	}).Info("Retrieve access key successfully")

	log.Info("Get buckets: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	buckets, err := client.GetObjectStorageBucketList(emptyCtx)
	if err != nil {
		log.Error("Retrieve buckets has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"buckets": buckets,
	}).Info("Retrieve buckets successfully")

	log.Info("Delete object storage access key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
