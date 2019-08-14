# The gridscale Go Client
[![Build Status](https://travis-ci.com/gridscale/gsclient-go.svg?branch=master)](https://travis-ci.com/gridscale/gsclient-go) [![Go Report Card](https://goreportcard.com/badge/github.com/gridscale/gsclient-go)](https://goreportcard.com/report/github.com/gridscale/gsclient-go)

This is a client for the gridscale API. It can be used to make an application interact with the gridscale cloud platform to create and manage resources.

## Preparation for Use

To be able to use this client a number of steps need to be taken. First a gridscale account will be required, which can be created [here](https://my.gridscale.io/signup/). Then an API-token [should be created](https://my.gridscale.io/APIs/).

## Installation

First the Go programming language will need to be installed. This can be done by using [the official Go installation guide](https://golang.org/doc/install) or by using the packages provided by your distribution.
 
Downloading the gridscale Go client can be done with the following go command:

```
go get "github.com/gridscale/gsclient-go"
```

## Using the gridscale Client

To be able to use the gridscale Go client in an application it can be imported in a go file. This can be done with the following code:

```go
import (
	"github.com/gridscale/gsclient-go"
)
```

To get access to the functions of the Go client, a Client type needs to be created. This requires a Config type. Both of these can be created with the following code: 

```go
config := gsclient.NewConfiguration("User-UUID", "API-token")
client := gsclient.NewClient(config)
```

Make sure to replace the user-UUID and API-token strings with valid credentials or variables containing valid credentials. It is recommended to use environment variables for them.

## Using API endpoints

After having created a Client type, as shown above, it will be possible to interact with the API. An example would be the [Servers Get endpoint](https://gridscale.io/en/api-documentation/index.html#servers-get):

```go
servers := client.GetServerList()
```

For creating and updating/patching objects in gridscale, it will be required to use the respective CreateRequest and UpdateRequest types. For creating an SSH-key that would be SshkeyCreateRequest and SshkeyUpdateRequest. Here an example:

```go
requestBody := gsclient.IpCreateRequest {
	Family: 6,
	Name:   "IpTest",
}

client.CreateIp(requestBody)
```

What options are available for each create and update request can be found in the source code. After installing it should be located in: 
```
~/go/src/github.com/gridscale/gsclient-go
```

## Implemented API Endpoints

Not all endpoints have been implemented in this client, but new ones will be added in the future. Here is the current list of implemented endpoints and their respective function written like endpoint (function):

* Servers
    * Servers Get (GetServerList)
    * Server Get (GetServer)
    * Server Create (CreateServer)
    * Server Patch (UpdateServer)
    * Server Delete (DeleteServer)
    * ACPI Shutdown (ShutdownServer)
    * Server On/Off (StartServer, StopServer)
    * Link Storage (LinkStorage)
    * Unlink Storage (UnlinkStorage)
    * Link Network (LinkNetwork)
    * Unlink Network (UnlinkNetwork)
    * Link IP (LinkIp)
    * Unlink IP (UnlinkIp)
    * Link ISO-Image (LinkIsoimage)
    * Unlink ISO-Image (UnlinkIsoimage)
* Storages
    * Storages Get (GetStorageList)
    * Storage Get (GetStorage)
    * Storage Create (CreateStorage)
    * Storage Patch (UpdateStorage)
    * Storage Delete (DeleteStorage)
* Networks
    * Networks Get (GetNetworkList)
    * Network Get (GetNetwork)
    * Network Create (CreateNetwork)
    * Network Patch (UpdateNetwork)
    * Network Delete (DeleteNetwork)
    * Network Events Get (GetNetworkEventList)
    * (GetNetworkPublic) No official endpoint, but gives the Public Network
* Loadbalancers
    * LoadBalancers Get (GetLoadBalancerList)
    * LoadBalancer Get (GetLoadBalancer)
    * LoadBalancer Create (CreateLoadBalancer)
    * LoadBalancer Patch (UpdateLoadBalancer)
    * LoadBalancer Delete (DeleteLoadBalancer)
    * LoadBalancerEvents Get (GetLoadBalancerEventList)
* IPs
    * IPs Get (GetIpList)
    * IP Get (GetIp)
    * IP Create (CreateIp)
    * IP Patch (UpdateIp)
    * IP Delete (DeleteIp)
    * IP Events Get (GetIpEventList)
    * IP Version Get (GetIpVersion)
* SSH-Keys
    * SSH-Keys Get (GetSshkeyList)
    * SSH-Key Get (GetSshkey)
    * SSH-Key Create (CreateSshkey)
    * SSH-Key Patch (UpdateSshkey)
    * SSH-Key Delete (DeleteSshkey)
* Template
    * Templates Get (GetTemplateList)
    * Template Get (GetTemplate)
    * (GetTemplateByName) No official endpoint, but gives a template which matches the exact name given.
* PaaS
    * PaaS services Get (GetPaaSServiceList)
    * PaaS service Get (GetPaaSService)
    * PaaS service Create (CreatePaaSService)
    * PaaS service Update (UpdatePaaSService)
    * PaaS service Delete (DeletePaaSService)
    * PaaS service metrics Get (GetPaaSServiceMetrics)
    * PaaS service templates Get (GetPaaSTemplateList)
    * PaaS service security zones Get (GetPaaSSecurityZoneList)
    * Paas service security zone Get (GetPaaSSecurityZone)
    * PaaS service security zone Create (CreatePaaSSecurityZone)
    * PaaS service security zone Update (UpdatePaaSSecurityZone)
    * PaaS service security zone Delete (DeletePaaSSecurityZone)
* ISO Image
    * ISO Images Get (GetISOImageList)
    * ISO Image Get (GetISOImage)
    * ISO Image Create (CreateISOImage)
    * ISO Image Update (UpdateISOImage)
    * ISO Image Delete (DeleteISOImage)
    * ISO Image Events Get (GetISOImageEventList)
* Object Storage
    * Object Storage's Access Keys Get (GetObjectStorageAccessKeyList)
    * Object Storage's Access Key Get (GetObjectStorageAccessKey)
    * Object Storage's Access Key Create (CreateObjectStorageAccessKey)
    * Object Storage's Access Key Delete (DeleteObjectStorageAccessKey)
    * Object Storage's Buckets Get (GetObjectStorageBucketList)
* Storage Snapshot Scheduler
    * Storage Snapshot Schedules Get (GetStorageSnapshotScheduleList)
    * Storage Snapshot Schedule Get (GetStorageSnapshotSchedule)
    * Storage Snapshot Schedule Create (CreateStorageSnapshotSchedule)
    * Storage Snapshot Schedule Update (UpdateStorageSnapshotSchedule)
    * Storage Snapshot Schedule Delete (DeleteStorageSnapshotSchedule)
* Storage Snapshot
    * Storage Snapshots Get (GetStorageSnapshotList)
    * Storage Snapshot Get (GetStorageSnapshot)
    * Storage Snapshot Create (CreateStorageSnapshot)
    * Storage Snapshot Update (UpdateStorageSnapshot)
    * Storage Snapshot Delete (DeleteStorageSnapshot)
    * Storage Rollback (RollbackStorage)
    * Storage Snapshot Export to S3 (ExportStorageSnapshotToS3)
* Firewall
    * Firewalls Get (GetFirewallList)
    * Firewall Get (GetFirewall)
    * Firewall Create (CreateFirewall)
    * Firewall Update (UpdateFirewall)
    * Firewall Delete (DeleteFirewall)
    * Firewall Events Get (GetFirewallEventList) 
    
Note: The functions in this list can be called with a Client type.

## Known Issues
The following issues are known to us:

* L3security isn't read in the network relation of a server.
* The autorecovery attribute of a server can't be changed with this client.
