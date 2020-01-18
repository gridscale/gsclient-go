package main

import (
	"bufio"
	"context"
	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
	"os"
)

var emptyCtx = context.Background()

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create firewall: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	labels := []string{"fw-http"}
	fwRequest := gsclient.FirewallCreateRequest{
		Name:   "go-client-firewall",
		Labels: labels,
		Rules: gsclient.FirewallRules{
			RulesV4In: []gsclient.FirewallRuleProperties{
				{
					Protocol: gsclient.TCPTransport,
					Action:   "accept",
					Order:    1,
				},
			},
		},
	}
	//Create a new firewall
	cfw, err := client.CreateFirewall(emptyCtx, fwRequest)
	if err != nil {
		log.Error("Create firewall has failed with error", err)
		return
	}
	log.WithFields(log.Fields{"Firewall_uuid": cfw.ObjectUUID}).Info("Firewall successfully created")
	log.Info("Update firewall: Press 'Enter' to continue...")
	defer func() {
		err := client.DeleteFirewall(emptyCtx, cfw.ObjectUUID)
		if err != nil {
			log.Error("Delete firewall has failed with error", err)
			return
		}
		log.Info("Firewall has successfully deleted")
	}()
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get a firewall to update
	fw, err := client.GetFirewall(emptyCtx, cfw.ObjectUUID)
	if err != nil {
		log.Errorf("Get firewall %s has failed with error %v", cfw.ObjectUUID, err)
		return
	}
	fwUpdateRequest := gsclient.FirewallUpdateRequest{
		Name:   "Updated name",
		Labels: &fw.Properties.Labels,
		Rules:  &fw.Properties.Rules,
	}
	err = client.UpdateFirewall(emptyCtx, fw.Properties.ObjectUUID, fwUpdateRequest)
	if err != nil {
		log.Error("Update firewall has failed with error", err)
		return
	}

	//Get firewall events
	events, err := client.GetFirewallEventList(emptyCtx, fw.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get firewall's events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"firewall_uuid": fw.Properties.ObjectUUID,
		"events":        events}).Info("Firewall's events successfully retrieved")
	log.Info("Delete firewall: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
