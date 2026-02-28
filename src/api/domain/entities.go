package domain

import "time"

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ResourceType string

const (
	ComputeResource ResourceType = "compute"
	NetworkResource ResourceType = "network"
	StorageResource ResourceType = "storage"
)

type Resource struct {
	ID         string                 `json:"id"`
	ProjectID  string                 `json:"project_id"`
	Name       string                 `json:"name"`
	Type       ResourceType           `json:"type"`
	Provider   string                 `json:"provider"`
	ProviderID string                 `json:"provider_id"`
	State      string                 `json:"state"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  time.Time              `json:"created_at"`
}

type Domain struct {
	ID        string    `json:"id"`
	FQDN      string    `json:"fqdn"`
	ProjectID string    `json:"project_id"`
	SSLActive bool      `json:"ssl_active"`
	CreatedAt time.Time `json:"created_at"`
}

type NetworkConfig struct {
	VPCID      string `json:"vpc_id"`
	SubnetCIDR string `json:"subnet_cidr"`
	Gateway    string `json:"gateway"`
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	TenantID  string    `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
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
