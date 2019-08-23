package main

import (
	"bufio"
	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
	"os"
)

const locationUUID = "45ed677b-3702-4b36-be2a-a2eab9827950"

//exampleSSHkey is an example of SSH-key, don't use it in production
const exampleSSHkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC9BlRsUvqRNKi59UkQmmztP5g+1jX5Ettr9C0+udwu9ATOukoM3rr0dXGVEVOJKQO1QCoEvMxn5HhZO2+klTVC1inapOrFrlUveqhcXvx6Fr1l3AmBsgY7loa5ELgi0qcKNcM/c9J7gB3EadKei/kfo5EXLDchn8SGHEq9Rhi8n8RcpGCEFnuvbao7uRsSj1QxTBaZgl5FL+W7wq2/dtwNhUk/KVA+ZKkMd4EnVlkF2ngQ02WQsu+0TN1gusMhBfph5sqtFT0twoOvYE3ejVaCc5LwT+5oxZulQ4TvggbJjzGD618q0QFkJ0CUtuh2s0otJkx1RqABX3TjfgmDjA8L example@gridscales.local"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration("https://api.gridscale.io", uuid, token, true)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create SSH-key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cSSHkey, err := client.CreateSshkey(gsclient.SshkeyCreateRequest{
		Name:   "go-client-ssh-key",
		Sshkey: exampleSSHkey,
	})
	if err != nil {
		log.Error("Create SSH-key has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"sshkey_uuid": cSSHkey.ObjectUUID,
	}).Info("SSH-key successfully created")
	defer func() {
		err := client.DeleteSshkey(cSSHkey.ObjectUUID)
		if err != nil {
			log.Error("Delete SSH-key has failed with error", err)
			return
		}
		log.Info("SSH-key successfully deleted")
	}()

	//Get a SSH-key to update
	sshkey, err := client.GetSshkey(cSSHkey.ObjectUUID)
	if err != nil {
		log.Error("Get SSH-key has failed with error", err)
		return
	}

	log.Info("Update SSH-key: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	err = client.UpdateSshkey(sshkey.Properties.ObjectUUID, gsclient.SshkeyUpdateRequest{
		Name:   "updated SSH-key",
		Sshkey: sshkey.Properties.Sshkey,
		Labels: sshkey.Properties.Labels,
	})
	if err != nil {
		log.Error("Update SSH-key has failed with error", err)
		return
	}
	log.Info("SSH-key successfully updated")

	log.Info("Get SSH-key's events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	events, err := client.GetSshkeyEventList(sshkey.Properties.ObjectUUID)
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
