package domain

import (
	"crypto/rand"
	"fmt"
	"time"
)

func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Department struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
}

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID           string    `json:"id"`
	DepartmentID string    `json:"department_id"`
	TenantID     string    `json:"tenant_id"` // Kept for legacy compatibility
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}

type ResourceType string

const (
	ComputeResource   ResourceType = "compute"
	ContainerResource ResourceType = "container"
	NetworkResource   ResourceType = "network"
	StorageResource   ResourceType = "storage"
)

type Resource struct {
	ID             string                 `json:"id"`
	ProjectID      string                 `json:"project_id"`
	Name           string                 `json:"name"`
	Type           ResourceType           `json:"type"`
	Provider       string                 `json:"provider"`
	ProviderID     string                 `json:"provider_id"`
	RegionID       string                 `json:"region_id,omitempty"`
	ZoneID         string                 `json:"zone_id,omitempty"`
	State          string                 `json:"state"`
	Metadata       map[string]interface{} `json:"metadata"`
	BlueprintID    string                 `json:"blueprint_id,omitempty"`
	SecurityGroups []string               `json:"security_groups,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

type ResourceInsight struct {
	ResourceID string    `json:"resource_id"`
	Type       string    `json:"type"` // e.g., "cost", "performance", "security"
	Message    string    `json:"message"`
	Severity   string    `json:"severity"` // "low", "medium", "high"
	Actionable bool      `json:"actionable"`
	CreatedAt  time.Time `json:"created_at"`
}

type GlobalEndpoint struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	DNSRecord string     `json:"dns_record"`
	Policy    GSLBPolicy `json:"policy"`
	Endpoints []string   `json:"endpoints"` // List of regional resource IDs (e.g., LBs)
	State     string     `json:"state"`
}

type GSLBPolicy struct {
	Strategy string `json:"strategy"` // "round-robin", "latency", "failover"
	Region   string `json:"region"`
}

type Domain struct {
	ID        string    `json:"id"`
	FQDN      string    `json:"fqdn"`
	ProjectID string    `json:"project_id"`
	SSLActive bool      `json:"ssl_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Network struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	CIDR      string    `json:"cidr"`
	Gateway   string    `json:"gateway"`
	Provider  string    `json:"provider"`
	RegionID  string    `json:"region_id,omitempty"`
	ZoneID    string    `json:"zone_id,omitempty"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type Route struct {
	ID             string    `json:"id"`
	NetworkID      string    `json:"network_id"`
	Destination    string    `json:"destination"` // CIDR
	NextHop        string    `json:"next_hop"`     // IP or "cross-provider-interface"
	TargetProvider string    `json:"target_provider,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type Container struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CPU       float64   `json:"cpu"`
	MemoryMB  int       `json:"memory_mb"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type Image struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Type     string `json:"type"` // "iso", "template", "docker"
}

type User struct {
	ID                 string    `json:"id"`
	Username           string    `json:"username"`
	PasswordHash       string    `json:"-"`
	Email              string    `json:"email"`
	TenantID           string    `json:"tenant_id"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Policy struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Statements []string `json:"statements"`
}

type Secret struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"` // Should be encrypted in storage
	CreatedAt time.Time `json:"created_at"`
}

type Quota struct {
	TenantID string `json:"tenant_id"`
	MaxCPUs  int    `json:"max_cpus"`
	MaxRAM   int    `json:"max_ram"`  // MB
	MaxDisk  int    `json:"max_disk"` // GB
}

type Volume struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"project_id"`
	Name       string    `json:"name"`
	SizeGB     int       `json:"size_gb"`
	ProviderID string    `json:"provider_id"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
}

type Bucket struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type BillingItem struct {
	ResourceID string  `json:"resource_id"`
	Type       string  `json:"type"`
	Cost       float64 `json:"cost"`
	Currency   string  `json:"currency"`
}

type BillingReport struct {
	TenantID    string        `json:"tenant_id"`
	TotalCost   float64       `json:"total_cost"`
	Items       []BillingItem `json:"items"`
	GeneratedAt time.Time     `json:"generated_at"`
}

type SovereigntyPolicy struct {
	TenantID       string   `json:"tenant_id"`
	AllowedRegions []string `json:"allowed_regions"`
}

type FirewallProtocol string

const (
	TCP  FirewallProtocol = "tcp"
	UDP  FirewallProtocol = "udp"
	ICMP FirewallProtocol = "icmp"
)

type FirewallRule struct {
	ID        string           `json:"id"`
	Protocol  FirewallProtocol `json:"protocol"`
	FromPort  int              `json:"from_port"`
	ToPort    int              `json:"to_port"`
	SourceIP  string           `json:"source_ip"` // CIDR or "any"
	Action    string           `json:"action"`    // "allow" or "deny"
	CreatedAt time.Time        `json:"created_at"`
}

type SecurityGroup struct {
	ID        string         `json:"id"`
	ProjectID string         `json:"project_id"`
	Name      string         `json:"name"`
	Rules     []FirewallRule `json:"rules"`
	CreatedAt time.Time      `json:"created_at"`
}

type GlobalStats struct {
	TotalCPUs     float64 `json:"total_cpus"`
	TotalStorage  float64 `json:"total_storage"`
	TotalEgress   float64 `json:"total_egress"`
	ActiveTenants int     `json:"active_tenants"`
	TrendCPUs     float64 `json:"trend_cpus"`
	TrendStorage  float64 `json:"trend_storage"`
}

type InfraPreset struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BlueprintID string `json:"blueprint_id"`
}

type AutomaticProvisioningRequest struct {
	PresetID  string                 `json:"preset_id"`
	ProjectID string                 `json:"project_id"`
	Name      string                 `json:"name"`
	Variables map[string]interface{} `json:"variables"`
}

type Region struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	IsDefault bool   `json:"is_default"`
}

type AvailabilityZone struct {
	ID       string    `json:"id"`
	RegionID string    `json:"region_id"`
	Name     string    `json:"name"`
	State    string    `json:"state"` // "available", "maintenance"
}

// Bare Metal & Hardware Inventory
type NodeState string

const (
	NodeStateAvailable    NodeState = "available"
	NodeStateProvisioning NodeState = "provisioning"
	NodeStateActive       NodeState = "active"
	NodeStateMaintenance  NodeState = "maintenance"
	NodeStateError        NodeState = "error"
)

type BareMetalNode struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	MAC          string    `json:"mac"`
	IPMIAddress  string    `json:"ipmi_address"`
	IPMIUser     string    `json:"ipmi_user"`
	IPMIPassword string    `json:"-"`
	CPUCores     int       `json:"cpu_cores"`
	MemoryGB     int       `json:"memory_gb"`
	DiskGB       int       `json:"disk_gb"`
	DepartmentID string    `json:"department_id,omitempty"`
	State        NodeState `json:"state"`
	ProviderID   string    `json:"provider_id,omitempty"` // The physical server provider ref
	CreatedAt    time.Time `json:"created_at"`
}

type ProvisioningLog struct {
	ID        string    `json:"id"`
	NodeID    string    `json:"node_id"`
	Message   string    `json:"message"`
	Level     string    `json:"level"` // "info", "error"
	Timestamp time.Time `json:"timestamp"`
}
