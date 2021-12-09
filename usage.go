package gsclient

import (
	"context"
	"net/http"
	"path"
	"strconv"
)

// Usage represents usage of a product.
type Usage struct {
	// Number of a product.
	ProductNumber int `json:"product_number"`

	// Total usage of a product.
	Value int `json:"value"`
}

// UsagePerInterval represents usage of active product within a specific interval.
type UsagePerInterval struct {
	// Start accumulation period.
	IntervalStart GSTime `json:"interval_start"`

	// interval_end
	IntervalEnd GSTime `json:"interval_end"`

	// Accumulation of product's usage in given period
	AccumulatedUsage []Usage `json:"accumulated_usage"`
}

// ResourceUsageInfo represents usage of a specific resource (e.g. server, storage, etc.).
type ResourceUsageInfo struct {
	CurrentUsagePerMinute []Usage            `json:"current_usage_per_minute"`
	UsagePerInterval      []UsagePerInterval `json:"usage_per_interval"`
}

// GeneralUsage represents general usage.
type GeneralUsage struct {
	Servers             ResourceUsageInfo `json:"servers"`
	RocketStorages      ResourceUsageInfo `json:"rocket_storages"`
	DistributedStorages ResourceUsageInfo `json:"distributed_storages"`
	StorageBackups      ResourceUsageInfo `json:"storage_backups"`
	Snapshots           ResourceUsageInfo `json:"snapshots"`
	Templates           ResourceUsageInfo `json:"templates"`
	IsoImages           ResourceUsageInfo `json:"iso_images"`
	IPAddresses         ResourceUsageInfo `json:"ip_addresses"`
	LoadBalancers       ResourceUsageInfo `json:"load_balancers"`
	PaaSServices        ResourceUsageInfo `json:"paas_services"`
}

// ServersUsage represents usage of servers.
type ServersUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// Indicates the amount of memory in GB.
	Memory int `json:"memory"`

	// Number of server cores.
	Cores int `json:"cores"`

	// The power status of the server.
	Power bool `json:"power"`

	// List of labels.
	Labels []string `json:"labels"`

	// True if the object is deleted.
	Deleted bool `json:"deleted"`

	// Status indicates the status of the object. it could be in-provisioning or active.
	Status string `json:"status"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// DistributedStoragesUsage represents usage of distributed storages.
type DistributedStoragesUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The UUID of the Storage used to create this Snapshot.
	ParentUUID string `json:"parent_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// List of labels.
	Labels []string `json:"labels"`

	// True if the object is deleted.
	Deleted bool `json:"deleted"`

	// Status indicates the status of the object. it could be in-provisioning or active.
	Status string `json:"status"`

	// Storage type.
	// (one of storage, storage_high, storage_insane).
	StorageType string `json:"storage_type"`

	// Indicates the UUID of the last used template on this storage.
	LastUsedTemplate string `json:"last_used_template"`

	// The capacity of a storage/ISO image/template/snapshot in GB.
	Capacity int `json:"capacity"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// RocketStoragesUsage represents usage of rocket storages.
type RocketStoragesUsage struct {
	DistributedStoragesUsage
}

// StorageBackupsUsage represents usage of storage backups.
type StorageBackupsUsage struct {
	// The UUID of a backup is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The name of the backup equals schedule name plus backup UUID.
	Name string `json:"name"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	// The size of a backup in GB.
	Capacity int `json:"capacity"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// SnapshotsUsage represents usage of snapshots.
type SnapshotsUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// Uuid of the storage used to create this snapshot.
	ParentUUID string `json:"parent_uuid"`

	// Name of the storage used to create this snapshot.
	ParentName string `json:"parent_name"`

	// Uuid of the project used to create this snapshot.
	ProjectUUID string `json:"project_uuid"`

	// List of labels.
	Labels []string `json:"labels"`

	// Status indicates the status of the object.
	Status string `json:"status"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	// The capacity of a storage/ISO image/template/snapshot in GB.
	Capacity int `json:"capacity"`

	// True if the object is deleted.
	Deleted bool `json:"deleted"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// TemplatesUsage represents usage of templates.
type TemplatesUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// Status indicates the status of the object.
	Status string `json:"status"`

	// Status indicates the status of the object.
	Ostype string `json:"ostype"`

	// A version string for this template.
	Version string `json:"version"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last change.
	ChangeTime GSTime `json:"change_time"`

	// Whether the object is private, the value will be true. Otherwise the value will be false.
	Private bool `json:"private"`

	// If a template has been used that requires a license key (e.g. Windows Servers)
	// this shows the product_no of the license (see the /prices endpoint for more details).
	LicenseProductNo int `json:"license_product_no"`

	// The capacity of a storage/ISO image/template/snapshot in GiB.
	Capacity int `json:"capacity"`

	// The OS distribution of this template.
	Distro string `json:"distro"`

	// Description of the template.
	Description string `json:"description"`

	// List of labels.
	Labels []string `json:"labels"`

	// Uuid of the project used to create this template.
	ProjectUUID string `json:"project_uuid"`

	// True if the object is deleted.
	Deleted bool `json:"deleted"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// ISOImagesUsage represents usage of ISO images.
type ISOImagesUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// Description of the ISO image release.
	Description string `json:"description"`

	// Contains the source URL of the ISO image that it was originally fetched from.
	SourceURL string `json:"source_url"`

	// List of labels.
	Labels []string `json:"labels"`

	// Status indicates the status of the object.
	Status string `json:"status"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	// Upstream version of the ISO image release.
	Version string `json:"version"`

	// Whether the ISO image is private or not.
	Private bool `json:"private"`

	// The capacity of an ISO image in GB.
	Capacity int `json:"capacity"`

	// Uuid of the project used to create this ISO image.
	ProjectUUID string `json:"project_uuid"`

	// True if the object is deleted.
	Deleted bool `json:"deleted"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// IPsUsage represents usage of IP addresses.
type IPsUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// Defines the IP Address (v4 or v6).
	IP string `json:"ip"`

	// Enum:4 6. The IP Address family (v4 or v6).
	Family int `json:"family"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	// Status indicates the status of the object.
	Status string `json:"status"`

	// The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.
	LocationCountry string `json:"location_country"`

	// The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.
	LocationName string `json:"location_name"`

	// Uses IATA airport code, which works as a location identifier.
	LocationIata string `json:"location_iata"`

	// Helps to identify which data center an object belongs to.
	LocationUUID string `json:"location_uuid"`

	// The IP prefix.
	Prefix string `json:"prefix"`

	// Defines if the object is administratively blocked. If true, it can not be deleted by the user.
	DeleteBlock bool `json:"delete_block"`

	// Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.
	Failover bool `json:"failover"`

	// List of labels.
	Labels []string `json:"labels"`

	// Defines the reverse DNS entry for the IP Address (PTR Resource Record).
	ReverseDNS string `json:"reverse_dns"`

	// Uuid of the partner used to create this IP address.
	PartnerUUID string `json:"partner_uuid"`

	// Uuid of the project used to create this IP address.
	ProjectUUID string `json:"project_uuid"`

	// True if the object is deleted.
	Deleted bool `json:"deleted"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// LoadBalancersUsage represents usage of storage backups.
type LoadBalancersUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// Forwarding rules of a load balancer.
	ForwardingRules []ForwardingRule `json:"forwarding_rules"`

	// The servers that this Load balancer can communicate with.
	BackendServers []BackendServer `json:"backend_servers"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	// Status indicates the status of the object.
	Status string `json:"status"`

	// Whether the Load balancer is forced to redirect requests from HTTP to HTTPS.
	RedirectHTTPToHTTPS bool `json:"redirect_http_to_https"`

	// The algorithm used to process requests. Accepted values: roundrobin / leastconn.
	Algorithm string `json:"algorithm"`

	// The UUID of the IPv6 address the Load balancer will listen to for incoming requests.
	ListenIPv6UUID string `json:"listen_ipv6_uuid"`

	// The UUID of the IPv4 address the Load balancer will listen to for incoming requests.
	ListenIPv4UUID string `json:"listen_ipv4_uuid"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// PaaSUsage represents usage of PaaS services.
type PaaSUsage struct {
	// The UUID of an object is always unique, and refers to a specific object.
	ObjectUUID string `json:"object_uuid"`

	// The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
	Name string `json:"name"`

	// Status indicates the status of the object.
	Status string `json:"status"`

	// Contains the initial setup credentials for Service.
	Credentials []Credential `json:"credentials"`

	// Defines the date and time the object was initially created.
	CreateTime GSTime `json:"create_time"`

	// Defines the date and time of the last object change.
	ChangeTime GSTime `json:"change_time"`

	// The template used to create the service, you can find an available list at the /service_templates endpoint.
	ServiceTemplateUUID string `json:"service_template_uuid"`

	// Contains the service parameters for the service.
	Parameters map[string]interface{} `json:"parameters"`

	// A list of service resource limits.
	ResourceLimits []ResourceLimit `json:"resource_limits"`

	// Uuid of the project used to create this PaaS.
	ProjectUUID string `json:"project_uuid"`

	// True if the object is deleted.
	Deleted bool `json:"deleted"`

	// Current usage of active product.
	CurrentUsagePerMinute []Usage `json:"current_usage_per_minute"`

	// Usage of active product within a specific interval.
	UsagePerInterval []UsagePerInterval `json:"usage_per_interval"`
}

// IntervalVariable represents interval variable when querying usage.
// IntervalVariable is used to get usage of resources within time intervals.
type IntervalVariable string

// All allowed interval variable's values
const (
	HourIntervalVariable  IntervalVariable = "H"
	DayIntervalVariable   IntervalVariable = "D"
	WeekIntervalVariable  IntervalVariable = "W"
	MonthIntervalVariable IntervalVariable = "M"
)

type usageQueryLevel int

const (
	// ProjectLevelUsage is used to query resources' usage in project level.
	ProjectLevelUsage usageQueryLevel = iota

	// ContractLevelUsage is used to query resources' usage in contract level.
	ContractLevelUsage = iota
)

// GetGeneralUsage returns general usage of all resources in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelUsageGet
func (c *Client) GetGeneralUsage(ctx context.Context, queryLevel usageQueryLevel, fromTime GSTime, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (GeneralUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 uri,
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response GeneralUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetServersUsage returns usage of all servers in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelServerUsageGet
func (c *Client) GetServersUsage(ctx context.Context, queryLevel usageQueryLevel, fromTime GSTime, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (ServersUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "servers"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response ServersUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetDistributedStoragesUsage returns usage of all distributed storages in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelDistributedStorageUsageGet
func (c *Client) GetDistributedStoragesUsage(ctx context.Context, queryLevel usageQueryLevel, fromTime GSTime, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (DistributedStoragesUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "distributed_storages"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response DistributedStoragesUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetRocketStoragesUsage returns usage of all servers in project/contract level.
// Args:
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelRocketStorageUsageGet
func (c *Client) GetRocketStoragesUsage(ctx context.Context, queryLevel usageQueryLevel, fromTime GSTime, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (RocketStoragesUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "rocket_storages"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response RocketStoragesUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetStorageBackupsUsage returns usage of all storage backups in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelStorageBackupUsageGet
func (c *Client) GetStorageBackupsUsage(ctx context.Context, fromTime GSTime, queryLevel usageQueryLevel, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (StorageBackupsUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "rocket_storages"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response StorageBackupsUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetSnapshotsUsage returns usage of all snapshots in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelSnapshotUsageGet
func (c *Client) GetSnapshotsUsage(ctx context.Context, fromTime GSTime, queryLevel usageQueryLevel, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (SnapshotsUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "snapshots"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response SnapshotsUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetTemplatesUsage returns usage of all templates in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelTemplateUsageGet
func (c *Client) GetTemplatesUsage(ctx context.Context, fromTime GSTime, queryLevel usageQueryLevel, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (TemplatesUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "templates"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response TemplatesUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetISOImagesUsage returns usage of all ISO images in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelIsoimageUsageGet
func (c *Client) GetISOImagesUsage(ctx context.Context, fromTime GSTime, queryLevel usageQueryLevel, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (ISOImagesUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "iso_images"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response ISOImagesUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetIPsUsage returns usage of all IP addresses' usage in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelIpUsageGet
func (c *Client) GetIPsUsage(ctx context.Context, fromTime GSTime, queryLevel usageQueryLevel, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (IPsUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "ip_addresses"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response IPsUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetLoadBalancersUsage returns usage of all Load Balancers' usage in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelLoadbalancerUsageGet
func (c *Client) GetLoadBalancersUsage(ctx context.Context, fromTime GSTime, queryLevel usageQueryLevel, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (LoadBalancersUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "load_balancers"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response LoadBalancersUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}

// GetPaaSServicesUsage returns usage of all PaaS services' usage in project/contract level.
// Args:
//		- queryLevel (Required): resources' usage query level. Either ProjectLevelUsage or ContractLevelUsage.
// 		- fromTime (Required): Starting time when the usage should be calculated.
//		- toTime (Optional, can be nil): End time when the usage should be calculated.
//		- withoutDeleted (Optional, can be nil): To calculate the usage with or without deleted resources.
//		- intervalVariable (Optional, can be nil): HourIntervalVariable, DayIntervalVariable, WeekIntervalVariable, MonthIntervalVariable.
//
// See: https://gridscale.io/en/api-documentation/index.html#operation/ProjectLevelPaasServiceUsageGet
func (c *Client) GetPaaSServicesUsage(ctx context.Context, fromTime GSTime, queryLevel usageQueryLevel, toTime *GSTime, withoutDeleted *bool, intervalVariable *IntervalVariable) (PaaSUsage, error) {
	queryParam := map[string]string{
		"from_time": fromTime.String(),
	}
	if toTime != nil {
		queryParam["to_time"] = toTime.String()
	}
	if withoutDeleted != nil {
		queryParam["without_deleted"] = strconv.FormatBool(*withoutDeleted)
	}
	if intervalVariable != nil {
		queryParam["interval_variable"] = string(*intervalVariable)
	}
	var uri string
	switch queryLevel {
	case ProjectLevelUsage:
		uri = apiProjectLevelUsage
	case ContractLevelUsage:
		uri = apiContractLevelUsage
	}
	r := gsRequest{
		uri:                 path.Join(uri, "paas_services"),
		method:              http.MethodGet,
		skipCheckingRequest: true,
		queryParameters:     queryParam,
	}
	var response PaaSUsage
	err := r.execute(ctx, *c, &response)
	return response, err
}
