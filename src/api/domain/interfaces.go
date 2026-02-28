package domain

import "context"

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id string) (*Tenant, error)
	List(ctx context.Context) ([]*Tenant, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id string) (*Project, error)
	GetByTenant(ctx context.Context, tenantID string) ([]*Project, error)
	List(ctx context.Context) ([]*Project, error)
}

type ResourceRepository interface {
	Create(ctx context.Context, resource *Resource) error
	GetByID(ctx context.Context, id string) (*Resource, error)
	GetByProject(ctx context.Context, projectID string) ([]*Resource, error)
	UpdateState(ctx context.Context, id string, state string) error
	List(ctx context.Context) ([]*Resource, error)
}

type CloudProvider interface {
	Provision(ctx context.Context, resource *Resource) error
	Decommission(ctx context.Context, resource *Resource) error
	GetStatus(ctx context.Context, resourceID string) (string, error)
}

type NetworkProvider interface {
	ConfigureIngress(ctx context.Context, domain string, targetService string) error
	RequestCertificate(ctx context.Context, domain string) error
	GetNetworkStatus(ctx context.Context, domain string) (string, error)
}

type IdentityManager interface {
	Authenticate(ctx context.Context, username, password string) (string, error) // Returns JWT
	ValidateToken(ctx context.Context, token string) (*User, error)
}

type SecretManager interface {
	StoreSecret(ctx context.Context, key string, value string) error
	GetSecret(ctx context.Context, key string) (string, error)
}

type QuotaRepository interface {
	GetByTenant(ctx context.Context, tenantID string) (*Quota, error)
	Update(ctx context.Context, quota *Quota) error
}

type StorageProvider interface {
	CreateVolume(ctx context.Context, volume *Volume) error
	DeleteVolume(ctx context.Context, volume *Volume) error
	CreateBucket(ctx context.Context, bucket *Bucket) error
	DeleteBucket(ctx context.Context, bucket *Bucket) error
}

type VolumeRepository interface {
	Create(ctx context.Context, vol *Volume) error
	GetByID(ctx context.Context, id string) (*Volume, error)
	ListByProject(ctx context.Context, projectID string) ([]*Volume, error)
}

type BucketRepository interface {
	Create(ctx context.Context, b *Bucket) error
	GetByID(ctx context.Context, id string) (*Bucket, error)
	ListByProject(ctx context.Context, projectID string) ([]*Bucket, error)
}

type BillingManager interface {
	GenerateReport(ctx context.Context, tenantID string) (*BillingReport, error)
}

type PolicyService interface {
	ValidateRegion(ctx context.Context, tenantID string, region string) error
	GetPolicy(ctx context.Context, tenantID string) (*SovereigntyPolicy, error)
	UpdatePolicy(ctx context.Context, policy *SovereigntyPolicy) error
}

type SovereigntyPolicyRepository interface {
	Save(ctx context.Context, p *SovereigntyPolicy) error
	GetByTenantID(ctx context.Context, tenantID string) (*SovereigntyPolicy, error)
}

type SecurityGroupRepository interface {
	Create(ctx context.Context, sg *SecurityGroup) error
	GetByID(ctx context.Context, id string) (*SecurityGroup, error)
	ListByProject(ctx context.Context, projectID string) ([]*SecurityGroup, error)
	AddRule(ctx context.Context, sgID string, rule FirewallRule) error
	RemoveRule(ctx context.Context, sgID string, ruleID string) error
}
