package main

import (
	"bufio"
	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
	"os"
)

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

	log.Info("Create label: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	_, err := client.CreateLabel(gsclient.LabelCreateRequest{
		Label: "go-client-label",
	})
	if err != nil {
		log.Error("Create label has failed with error", err)
		return
	}
	log.Info("Label successfully created")
	defer func() {
		err := client.DeleteLabel("go-client-label")
		if err != nil {
			log.Error("Delete label has failed with error", err)
			return
		}
	}()

	log.Info("Retrieve labels: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	labels, err := client.GetLabelList()
	if err != nil {
		log.Error("Retrieve labels has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"labels": labels,
	}).Info("Labels successfully retrieved")

	log.Info("Delete label: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
