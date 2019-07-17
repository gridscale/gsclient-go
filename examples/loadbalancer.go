package main

import (
	"bufio"
	"net/http"
	"os"
	"time"

	"bitbucket.org/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
)

const LocationUuid = "45ed677b-3702-4b36-be2a-a2eab9827950"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.Config{
		APIUrl:     "https://api.gridscale.io",
		UserUUID:   uuid,
		APIToken:   token,
		HTTPClient: http.DefaultClient,
	}
	client := gsclient.NewClient(&config)
	log.Info("gridscale client configured")

	log.Info("Create IPs and loadbalancer: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// required to create IPv6 and IPv4 to create LB
	ipv4, _ := client.CreateIp(gsclient.IpCreateRequest{
		Family:       4,
		LocationUuid: LocationUuid,
	})
	log.Info("IPv4 has been created")

	ipv6, _ := client.CreateIp(gsclient.IpCreateRequest{
		Family:       6,
		LocationUuid: LocationUuid,
	})
	log.Info("[INFO] IPv6 has been created")

	// populate settings into LoadBalancerCreateRequest
	labels := make([]interface{}, 0)
	labels = append(labels, "lb-http")
	lbRequest := gsclient.LoadBalancerCreateRequest{
		Name:                "go-client-lb",
		Algorithm:           "leastconn",
		LocationUuid:        LocationUuid,
		ListenIPv6Uuid:      ipv6.ObjectUuid,
		ListenIPv4Uuid:      ipv4.ObjectUuid,
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
		"Loadbalancer_uuid": clb.ObjectUuid}).Info("Loadbalancer successfully created")

	log.Info("Update loadbalacer: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Get the loadbalacer to update some settings
	glb, err := client.GetLoadBalancer(clb.ObjectUuid)
	if err != nil {
		log.Fatal("Get loadbalancer has failed with error", err)
	}

	lbUpdateRequest := gsclient.LoadBalancerUpdateRequest{
		Name:                "go-client-lb233",
		Algorithm:           glb.Properties.Algorithm,
		LocationUuid:        glb.Properties.LocationUuid,
		ListenIPv6Uuid:      glb.Properties.ListenIPv6Uuid,
		ListenIPv4Uuid:      glb.Properties.ListenIPv4Uuid,
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
	err = client.UpdateLoadBalancer(glb.Properties.ObjectUuid, lbUpdateRequest)

	if err != nil {
		log.Fatal("Update loadbalancer has failed with error", err)
	}
	log.WithFields(log.Fields{
		"Loadbalancer_uuid": glb.Properties.ObjectUuid}).Info("Loadbalancer successfully updated")

	log.Info("Retrive loadbalancer events: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//Get loadbalancer events
	response, err := client.GetLoadBalancerEventList(glb.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Events loadbalancer has failed with error", err)
	}
	log.WithFields(log.Fields{
		"Loadbalancer_uuid": glb.Properties.ObjectUuid,
		"events":            response.Events,
	}).Info("Loadbalancer successfully events retrived")

	log.Info("Delete IPs and loadbalancer: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// finallly clean up delete IPs and loadbalancer
	err = client.DeleteLoadBalancer(glb.Properties.ObjectUuid)
	if err != nil {
		log.Fatal("Delete loadbalancer has failed with error", err)
	}
	log.WithFields(log.Fields{
		"Loadbalancer_uuid": glb.Properties.ObjectUuid}).Info("Loadbalancer successfully deleted")

	time.Sleep(10 * time.Second)

	err = client.DeleteIp(ipv4.ObjectUuid)
	if err != nil {
		log.Fatal("Delete ipv4 has failed with error", err)
	}
	log.Info("IPv4 successfully deleted")

	err = client.DeleteIp(ipv6.ObjectUuid)
	if err != nil {
		log.Fatal("Delete ipv6 has failed with error", err)
	}
	log.Info("IPv6 successfully deleted")
}
