package main

import (
	"bufio"
	"context"
	"os"

	"github.com/gridscale/gsclient-go/v3"
	log "github.com/sirupsen/logrus"
)

var emptyCtx = context.Background()

// exampleSSHkey is an example of SSH-key, don't use it in production
const exampleSSHkey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDPUCeKyPUNBZOikJKx+Id7udqm/ZKArvCn2AqwwRr02 john@example.com"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create SSH-key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cSSHkey, err := client.CreateSSHKey(
		emptyCtx,
		gsclient.SSHKeyCreateRequest{
			Name:   "go-client-ssh-key",
			SSHKey: exampleSSHkey,
		})
	if err != nil {
		log.Error("Create SSH-key has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"sshkey_uuid": cSSHkey.ObjectUUID,
	}).Info("SSH-key successfully created")
	defer func() {
		err := client.DeleteSSHKey(emptyCtx, cSSHkey.ObjectUUID)
		if err != nil {
			log.Error("Delete SSH-key has failed with error", err)
			return
		}
		log.Info("SSH-key successfully deleted")
	}()

	// Get a SSH-key to update
	sshkey, err := client.GetSSHKey(emptyCtx, cSSHkey.ObjectUUID)
	if err != nil {
		log.Error("Get SSH-key has failed with error", err)
		return
	}

	log.Info("Update SSH-key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	err = client.UpdateSSHKey(
		emptyCtx,
		sshkey.Properties.ObjectUUID,
		gsclient.SSHKeyUpdateRequest{
			Name:   "updated SSH-key",
			SSHKey: sshkey.Properties.SSHKey,
			Labels: &sshkey.Properties.Labels,
		})
	if err != nil {
		log.Error("Update SSH-key has failed with error", err)
		return
	}
	log.Info("SSH-key successfully updated")

	log.Info("Get SSH-key's events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	events, err := client.GetSSHKeyEventList(emptyCtx, sshkey.Properties.ObjectUUID)
	if err != nil {
		log.Error("Get SSH-key's events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("SSH-key's events successfully retrieved")

	log.Info("Delete SSH-key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
