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

	locationList, err := client.GetLocationList(emptyCtx)
	if err != nil {
		log.Error("GetLocationList has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"locations": locationList,
	}).Info("Locations successfully retrieved")
	log.Info("Create new location: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	locationCreateRes, err := client.CreateLocation(emptyCtx, gsclient.LocationCreateRequest{
		Name:               "Test location",
		ParentLocationUUID: locationList[0].Properties.ObjectUUID,
		ProductNo:          1500001,
		CPUNodeCount:       1,
	})
	if err != nil {
		log.Error("CreateLocation has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"location uuid": locationCreateRes.ObjectUUID,
	}).Info("A new location is successfully created")
	defer func() {
		err := client.DeleteLocation(emptyCtx, locationCreateRes.ObjectUUID)
		if err != nil {
			log.Error("DeleteLocation has failed with error", err)
			return
		}
		log.WithFields(log.Fields{
			"location uuid": locationCreateRes.ObjectUUID,
		}).Info("Location is successfully deleted")
	}()

	log.Info("Get new location data: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	newLoc, err := client.GetLocation(emptyCtx, locationCreateRes.ObjectUUID)
	if err != nil {
		log.Error("GetLocation has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"new location data": newLoc,
	}).Info("New location is successfully retrieved")
	log.Info("Update location: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	err = client.UpdateLocation(emptyCtx, newLoc.Properties.ObjectUUID, gsclient.LocationUpdateRequest{
		Name: "test updated location",
	})
	if err != nil {
		log.Error("UpdateLocation has failed with error", err)
		return
	}
	updatedLoc, err := client.GetLocation(emptyCtx, newLoc.Properties.ObjectUUID)
	if err != nil {
		log.Error("GetLocation has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"updated location data": updatedLoc,
	}).Info("Updated location is successfully retrieved")
}
