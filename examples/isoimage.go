package main

import (
	"bufio"
	"context"
	"github.com/gridscale/gsclient-go"
	"github.com/sirupsen/logrus"
	"os"
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
	logrus.Info("gridscale client configured")

	logrus.Info("Create ISO-image: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	isoRequest := gsclient.ISOImageCreateRequest{
		Name:         "go-client-iso",
		SourceURL:    "http://tinycorelinux.net/10.x/x86/release/TinyCore-current.iso",
		LocationUUID: locationUUID,
	}
	cIso, err := client.CreateISOImage(context.Background(), isoRequest)
	if err != nil {
		logrus.Error("Create ISO-image has failed with error", err)
		return
	}
	logrus.WithFields(logrus.Fields{"isoimage_uuid": cIso.ObjectUUID}).Info("ISO Image successfully created")
	defer func() {
		//Delete ISO-image
		err := client.DeleteISOImage(context.Background(), cIso.ObjectUUID)
		if err != nil {
			logrus.Error("Delete ISO-image has failed with error", err)
			return
		}
		logrus.Info("ISO-image successfully deleted")

		logrus.Info("Get deleted ISO-images: Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		isoImages, err := client.GetDeletedISOImages(context.Background())
		if err != nil {
			logrus.Error("Get deleted ISO-images has failed with error", err)
			return
		}
		logrus.WithFields(logrus.Fields{
			"ISO-images": isoImages,
		}).Info("Retrieved deleted ISO-images successfully")
	}()

	logrus.Info("Update ISO image: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get ISO-image to update
	iso, err := client.GetISOImage(context.Background(), cIso.ObjectUUID)
	if err != nil {
		logrus.Error("Get ISO-image has failed with error", err)
		return
	}

	isoUpdateRequest := gsclient.ISOImageUpdateRequest{
		Name:   "updated ISO",
		Labels: iso.Properties.Labels,
	}
	err = client.UpdateISOImage(context.Background(), iso.Properties.ObjectUUID, isoUpdateRequest)
	if err != nil {
		logrus.Error("Update ISO-image has failed with error", err)
		return
	}
	logrus.WithFields(logrus.Fields{"isoimage_uuid": iso.Properties.ObjectUUID}).Info("ISO image successfully updated")

	//get ISO-image's events
	events, err := client.GetISOImageEventList(context.Background(), iso.Properties.ObjectUUID)
	if err != nil {
		logrus.Error("Get ISO-image's events has failed with error", err)
		return
	}
	logrus.WithFields(logrus.Fields{
		"isoimage_uuid": iso.Properties.ObjectUUID,
		"events":        events,
	}).Info("Events successfully retrieved")

	logrus.Info("Delete ISO-image: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
