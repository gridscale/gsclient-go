package main

import (
	"bufio"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gridscale/gsclient-go"
)

const locationUUID = "45ed677b-3702-4b36-be2a-a2eab9827950"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration("https://api.gridscale.io", uuid, token, false, 0, 0, 0)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create IPs and loadbalancer: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// required to create IPv6 and IPv4 to create LB
	ipv4, _ := client.CreateIP(gsclient.IPCreateRequest{
		Family:       gsclient.IPv4Type,
		LocationUUID: locationUUID,
	})
	log.Info("IPv4 has been created")

	ipv6, _ := client.CreateIP(gsclient.IPCreateRequest{
		Family:       gsclient.IPv6Type,
		LocationUUID: locationUUID,
	})
	log.Info("[INFO] IPv6 has been created")

	// populate settings into LoadBalancerCreateRequest
	labels := make([]string, 0)
	labels = append(labels, "lb-http")
	lbRequest := gsclient.LoadBalancerCreateRequest{
		Name:                "go-client-lb",
		Algorithm:           "leastconn",
		LocationUUID:        locationUUID,
		ListenIPv6UUID:      ipv6.ObjectUUID,
		ListenIPv4UUID:      ipv4.ObjectUUID,
		RedirectHTTPToHTTPS: false,
		ForwardingRules: []gsclient.ForwardingRule{
			{
				LetsencryptSSL: nil,
				ListenPort:     8080,
				Mode:           "http",
				TargetPort:     8000,
			},
		},
		BackendServers: []gsclient.BackendServer{
			{
				Weight: 100,
				Host:   "185.201.147.176",
			},
		},
		Labels: labels,
	}

	clb, err := client.CreateLoadBalancer(lbRequest)
	if err != nil {
		log.Fatal("Create loadbalancer has failed with error", err)
	}
	log.WithFields(log.Fields{
		"Loadbalancer_uuid": clb.ObjectUUID}).Info("Loadbalancer successfully created")

	// Get the loadbalacer to update some settings
	glb, err := client.GetLoadBalancer(clb.ObjectUUID)
	if err != nil {
		log.Fatal("Get loadbalancer has failed with error", err)
	}

	log.Info("Update loadbalacer: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	lbUpdateRequest := gsclient.LoadBalancerUpdateRequest{
		Name:                "go-client-lb233",
		Algorithm:           glb.Properties.Algorithm,
		LocationUUID:        glb.Properties.LocationUUID,
		ListenIPv6UUID:      glb.Properties.ListenIPv6UUID,
		ListenIPv4UUID:      glb.Properties.ListenIPv4UUID,
		RedirectHTTPToHTTPS: glb.Properties.RedirectHTTPToHTTPS,
		ForwardingRules: []gsclient.ForwardingRule{
			{
				LetsencryptSSL: nil,
				ListenPort:     443,
				Mode:           "http",
				TargetPort:     443,
			},
		},
		BackendServers: glb.Properties.BackendServers,
		Labels:         labels,
	}
	err = client.UpdateLoadBalancer(glb.Properties.ObjectUUID, lbUpdateRequest)

	if err != nil {
		log.Fatal("Update loadbalancer has failed with error", err)
	}
	log.WithFields(log.Fields{
		"Loadbalancer_uuid": glb.Properties.ObjectUUID}).Info("Loadbalancer successfully updated")

	log.Info("Retrive loadbalancer events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get loadbalancer events
	response, err := client.GetLoadBalancerEventList(glb.Properties.ObjectUUID)
	if err != nil {
		log.Fatal("Events loadbalancer has failed with error", err)
	}
	log.WithFields(log.Fields{
		"Loadbalancer_uuid": glb.Properties.ObjectUUID,
		"events":            response,
	}).Info("Loadbalancer successfully events retrieved")

	log.Info("Delete IPs and loadbalancer: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// finallly clean up delete IPs and loadbalancer
	err = client.DeleteLoadBalancer(glb.Properties.ObjectUUID)
	if err != nil {
		log.Fatal("Delete loadbalancer has failed with error", err)
	}
	log.WithFields(log.Fields{
		"Loadbalancer_uuid": glb.Properties.ObjectUUID}).Info("Loadbalancer successfully deleted")

	time.Sleep(10 * time.Second)

	err = client.DeleteIP(ipv4.ObjectUUID)
	if err != nil {
		log.Fatal("Delete ipv4 has failed with error", err)
	}
	log.Info("IPv4 successfully deleted")

	err = client.DeleteIP(ipv6.ObjectUUID)
	if err != nil {
		log.Fatal("Delete ipv6 has failed with error", err)
	}
	log.Info("IPv6 successfully deleted")
}
