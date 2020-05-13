package main

import (
	"bufio"
	"context"
	"os"

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

	log.Info("Retrieve labels: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	labels, err := client.GetLabelList(emptyCtx)
	if err != nil {
		log.Error("Retrieve labels has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"labels": labels,
	}).Info("Labels successfully retrieved")
}
