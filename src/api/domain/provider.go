package domain

import "time"

type ProviderType string

const (
	ProxmoxProvider   ProviderType = "proxmox"
	OpenStackProvider ProviderType = "openstack"
	AWSProvider       ProviderType = "aws"
	BareMetalProvider ProviderType = "baremetal"
	K8sProvider       ProviderType = "kubernetes"
)

type Provider struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        ProviderType `json:"type"`
	Endpoint    string       `json:"endpoint"`
	Credentials string       `json:"credentials,omitempty"` // Sensitive data, usually encrypted in storage
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
}
