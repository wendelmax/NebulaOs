package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type InMemoryTenantRepository struct {
	mu      sync.RWMutex
	tenants map[string]*domain.Tenant
}

func NewInMemoryTenantRepository() *InMemoryTenantRepository {
	return &InMemoryTenantRepository{
		tenants: make(map[string]*domain.Tenant),
	}
}

func (r *InMemoryTenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tenants[tenant.ID]; ok {
		return fmt.Errorf("tenant already exists")
	}
	r.tenants[tenant.ID] = tenant
	return nil
}

func (r *InMemoryTenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tenant, ok := r.tenants[id]
	if !ok {
		return nil, fmt.Errorf("tenant not found")
	}
	return tenant, nil
}

func (r *InMemoryTenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*domain.Tenant, 0, len(r.tenants))
	for _, t := range r.tenants {
		list = append(list, t)
	}
	return list, nil
}

type InMemoryProjectRepository struct {
	mu       sync.RWMutex
	projects map[string]*domain.Project
}

func NewInMemoryProjectRepository() *InMemoryProjectRepository {
	return &InMemoryProjectRepository{
		projects: make(map[string]*domain.Project),
	}
}

func (r *InMemoryProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[project.ID] = project
	return nil
}

func (r *InMemoryProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.projects[id]
	if !ok {
		return nil, fmt.Errorf("project not found")
	}
	return p, nil
}

func (r *InMemoryProjectRepository) GetByTenant(ctx context.Context, tenantID string) ([]*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Project
	for _, p := range r.projects {
		if p.TenantID == tenantID {
			list = append(list, p)
		}
	}
	return list, nil
}

func (r *InMemoryProjectRepository) GetByDepartment(ctx context.Context, deptID string) ([]*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Project
	for _, p := range r.projects {
		if p.DepartmentID == deptID {
			list = append(list, p)
		}
	}
	return list, nil
}

func (r *InMemoryProjectRepository) List(ctx context.Context) ([]*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*domain.Project, 0, len(r.projects))
	for _, p := range r.projects {
		list = append(list, p)
	}
	return list, nil
}

type InMemoryResourceRepository struct {
	mu        sync.RWMutex
	resources map[string]*domain.Resource
}

func NewInMemoryResourceRepository() *InMemoryResourceRepository {
	return &InMemoryResourceRepository{
		resources: make(map[string]*domain.Resource),
	}
}

func (r *InMemoryResourceRepository) Create(ctx context.Context, res *domain.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resources[res.ID] = res
	return nil
}

func (r *InMemoryResourceRepository) GetByID(ctx context.Context, id string) (*domain.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res, ok := r.resources[id]
	if !ok {
		return nil, fmt.Errorf("resource not found")
	}
	return res, nil
}

func (r *InMemoryResourceRepository) GetByProject(ctx context.Context, projectID string) ([]*domain.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Resource
	for _, res := range r.resources {
		if res.ProjectID == projectID {
			list = append(list, res)
		}
	}
	return list, nil
}

func (r *InMemoryResourceRepository) UpdateState(ctx context.Context, id string, state string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, ok := r.resources[id]
	if !ok {
		return fmt.Errorf("resource not found")
	}
	res.State = state
	return nil
}

func (r *InMemoryResourceRepository) List(ctx context.Context) ([]*domain.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*domain.Resource, 0, len(r.resources))
	for _, res := range r.resources {
		list = append(list, res)
	}
	return list, nil
}

type InMemoryQuotaRepository struct {
	mu     sync.RWMutex
	quotas map[string]*domain.Quota
}

func NewInMemoryQuotaRepository() *InMemoryQuotaRepository {
	return &InMemoryQuotaRepository{
		quotas: make(map[string]*domain.Quota),
	}
}

func (r *InMemoryQuotaRepository) GetByTenant(ctx context.Context, tenantID string) (*domain.Quota, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q, ok := r.quotas[tenantID]
	if !ok {
		// Return a default quota if none exists
		return &domain.Quota{
			TenantID: tenantID,
			MaxCPUs:  4,
			MaxRAM:   8192,
			MaxDisk:  100,
		}, nil
	}
	return q, nil
}

func (r *InMemoryQuotaRepository) Update(ctx context.Context, quota *domain.Quota) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.quotas[quota.TenantID] = quota
	return nil
}

type InMemoryVolumeRepository struct {
	mu      sync.RWMutex
	volumes map[string]*domain.Volume
}

func NewInMemoryVolumeRepository() *InMemoryVolumeRepository {
	return &InMemoryVolumeRepository{
		volumes: make(map[string]*domain.Volume),
	}
}

func (r *InMemoryVolumeRepository) Create(ctx context.Context, vol *domain.Volume) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.volumes[vol.ID] = vol
	return nil
}

func (r *InMemoryVolumeRepository) GetByID(ctx context.Context, id string) (*domain.Volume, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.volumes[id]
	if !ok {
		return nil, fmt.Errorf("volume not found")
	}
	return v, nil
}

func (r *InMemoryVolumeRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.Volume, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Volume
	for _, v := range r.volumes {
		if v.ProjectID == projectID {
			list = append(list, v)
		}
	}
	return list, nil
}

type InMemoryBucketRepository struct {
	mu      sync.RWMutex
	buckets map[string]*domain.Bucket
}

func NewInMemoryBucketRepository() *InMemoryBucketRepository {
	return &InMemoryBucketRepository{
		buckets: make(map[string]*domain.Bucket),
	}
}

func (r *InMemoryBucketRepository) Create(ctx context.Context, b *domain.Bucket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buckets[b.ID] = b
	return nil
}

func (r *InMemoryBucketRepository) GetByID(ctx context.Context, id string) (*domain.Bucket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.buckets[id]
	if !ok {
		return nil, fmt.Errorf("bucket not found")
	}
	return b, nil
}

func (r *InMemoryBucketRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.Bucket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Bucket
	for _, b := range r.buckets {
		if b.ProjectID == projectID {
			list = append(list, b)
		}
	}
	return list, nil
}

type InMemorySovereigntyPolicyRepository struct {
	policies map[string]*domain.SovereigntyPolicy
	mu       sync.RWMutex
}

func NewInMemorySovereigntyPolicyRepository() *InMemorySovereigntyPolicyRepository {
	return &InMemorySovereigntyPolicyRepository{policies: make(map[string]*domain.SovereigntyPolicy)}
}

func (r *InMemorySovereigntyPolicyRepository) Save(ctx context.Context, p *domain.SovereigntyPolicy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[p.TenantID] = p
	return nil
}

func (r *InMemorySovereigntyPolicyRepository) GetByTenantID(ctx context.Context, tenantID string) (*domain.SovereigntyPolicy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.policies[tenantID]
	if !ok {
		return nil, fmt.Errorf("policy not found")
	}
	return p, nil
}



type InMemorySecurityGroupRepository struct {
	mu             sync.RWMutex
	securityGroups map[string]*domain.SecurityGroup
}

func NewInMemorySecurityGroupRepository() *InMemorySecurityGroupRepository {
	return &InMemorySecurityGroupRepository{
		securityGroups: make(map[string]*domain.SecurityGroup),
	}
}

func (r *InMemorySecurityGroupRepository) Create(ctx context.Context, sg *domain.SecurityGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.securityGroups[sg.ID] = sg
	return nil
}

func (r *InMemorySecurityGroupRepository) GetByID(ctx context.Context, id string) (*domain.SecurityGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sg, ok := r.securityGroups[id]
	if !ok {
		return nil, fmt.Errorf("security group not found")
	}
	return sg, nil
}

func (r *InMemorySecurityGroupRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.SecurityGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.SecurityGroup
	for _, sg := range r.securityGroups {
		if sg.ProjectID == projectID {
			list = append(list, sg)
		}
	}
	return list, nil
}

func (r *InMemorySecurityGroupRepository) AddRule(ctx context.Context, sgID string, rule domain.FirewallRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sg, ok := r.securityGroups[sgID]
	if !ok {
		return fmt.Errorf("security group not found")
	}
	sg.Rules = append(sg.Rules, rule)
	return nil
}

func (r *InMemorySecurityGroupRepository) RemoveRule(ctx context.Context, sgID string, ruleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	sg, ok := r.securityGroups[sgID]
	if !ok {
		return fmt.Errorf("security group not found")
	}
	for i, rule := range sg.Rules {
		if rule.ID == ruleID {
			sg.Rules = append(sg.Rules[:i], sg.Rules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("rule not found")
}

type InMemoryTerraformStateRepository struct {
	mu     sync.RWMutex
	states map[string]*domain.TerraformState
}

func NewInMemoryTerraformStateRepository() *InMemoryTerraformStateRepository {
	return &InMemoryTerraformStateRepository{states: make(map[string]*domain.TerraformState)}
}

func (r *InMemoryTerraformStateRepository) Save(ctx context.Context, s *domain.TerraformState) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.states[s.ProjectID] = s
	return nil
}

func (r *InMemoryTerraformStateRepository) GetByProjectID(ctx context.Context, projectID string) (*domain.TerraformState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.states[projectID]
	if !ok {
		return nil, fmt.Errorf("state not found")
	}
	return s, nil
}

type InMemoryBlueprintRepository struct {
	mu         sync.RWMutex
	blueprints map[string]*domain.Blueprint
}

func NewInMemoryBlueprintRepository() *InMemoryBlueprintRepository {
	return &InMemoryBlueprintRepository{blueprints: make(map[string]*domain.Blueprint)}
}

func (r *InMemoryBlueprintRepository) List(ctx context.Context) ([]*domain.Blueprint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Blueprint
	for _, b := range r.blueprints {
		list = append(list, b)
	}
	return list, nil
}

func (r *InMemoryBlueprintRepository) GetByID(ctx context.Context, id string) (*domain.Blueprint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.blueprints[id]
	if !ok {
		return nil, fmt.Errorf("blueprint not found")
	}
	return b, nil
}

func (r *InMemoryBlueprintRepository) Create(ctx context.Context, b *domain.Blueprint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blueprints[b.ID] = b
	return nil
}

// --- GSLB Repository ---

type InMemoryGSLBRepository struct {
	mu        sync.RWMutex
	endpoints map[string]*domain.GlobalEndpoint
}

func NewInMemoryGSLBRepository() *InMemoryGSLBRepository {
	return &InMemoryGSLBRepository{
		endpoints: make(map[string]*domain.GlobalEndpoint),
	}
}

func (r *InMemoryGSLBRepository) Save(ctx context.Context, ep *domain.GlobalEndpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.endpoints[ep.ID] = ep
	return nil
}

func (r *InMemoryGSLBRepository) GetByID(ctx context.Context, id string) (*domain.GlobalEndpoint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ep, ok := r.endpoints[id]
	if !ok {
		return nil, fmt.Errorf("endpoint not found")
	}
	return ep, nil
}

func (r *InMemoryGSLBRepository) List(ctx context.Context) ([]*domain.GlobalEndpoint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.GlobalEndpoint
	for _, ep := range r.endpoints {
		list = append(list, ep)
	}
	return list, nil
}

// --- Region & Zone Repositories ---

type InMemoryRegionRepository struct {
	mu      sync.RWMutex
	regions map[string]*domain.Region
}

func NewInMemoryRegionRepository() *InMemoryRegionRepository {
	return &InMemoryRegionRepository{regions: make(map[string]*domain.Region)}
}

func (r *InMemoryRegionRepository) Create(ctx context.Context, reg *domain.Region) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.regions[reg.ID] = reg
	return nil
}

func (r *InMemoryRegionRepository) GetByID(ctx context.Context, id string) (*domain.Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	reg, ok := r.regions[id]
	if !ok {
		return nil, fmt.Errorf("region not found")
	}
	return reg, nil
}

func (r *InMemoryRegionRepository) List(ctx context.Context) ([]*domain.Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Region
	for _, reg := range r.regions {
		list = append(list, reg)
	}
	return list, nil
}

type InMemoryAvailabilityZoneRepository struct {
	mu    sync.RWMutex
	zones map[string]*domain.AvailabilityZone
}

func NewInMemoryAvailabilityZoneRepository() *InMemoryAvailabilityZoneRepository {
	return &InMemoryAvailabilityZoneRepository{zones: make(map[string]*domain.AvailabilityZone)}
}

func (r *InMemoryAvailabilityZoneRepository) Create(ctx context.Context, az *domain.AvailabilityZone) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.zones[az.ID] = az
	return nil
}

func (r *InMemoryAvailabilityZoneRepository) GetByRegion(ctx context.Context, regionID string) ([]*domain.AvailabilityZone, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.AvailabilityZone
	for _, az := range r.zones {
		if az.RegionID == regionID {
			list = append(list, az)
		}
	}
	return list, nil
}

func (r *InMemoryAvailabilityZoneRepository) List(ctx context.Context) ([]*domain.AvailabilityZone, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.AvailabilityZone
	for _, az := range r.zones {
		list = append(list, az)
	}
	return list, nil
}

type InMemoryProviderRepository struct {
	mu        sync.RWMutex
	providers map[string]*domain.Provider
}

func NewInMemoryProviderRepository() *InMemoryProviderRepository {
	return &InMemoryProviderRepository{providers: make(map[string]*domain.Provider)}
}

func (r *InMemoryProviderRepository) Create(ctx context.Context, p *domain.Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.ID] = p
	return nil
}

func (r *InMemoryProviderRepository) GetByID(ctx context.Context, id string) (*domain.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[id]
	if !ok {
		return nil, fmt.Errorf("provider not found")
	}
	return p, nil
}

func (r *InMemoryProviderRepository) List(ctx context.Context) ([]*domain.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Provider
	for _, p := range r.providers {
		list = append(list, p)
	}
	return list, nil
}

func (r *InMemoryProviderRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.providers[id]; !ok {
		return fmt.Errorf("provider not found")
	}
	delete(r.providers, id)
	return nil
}

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*domain.User)}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

type InMemoryOrganizationRepository struct {
	mu   sync.RWMutex
	orgs map[string]*domain.Organization
}

func NewInMemoryOrganizationRepository() *InMemoryOrganizationRepository {
	return &InMemoryOrganizationRepository{orgs: make(map[string]*domain.Organization)}
}

func (r *InMemoryOrganizationRepository) Create(ctx context.Context, o *domain.Organization) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orgs[o.ID] = o
	return nil
}

func (r *InMemoryOrganizationRepository) GetByID(ctx context.Context, id string) (*domain.Organization, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orgs[id]
	if !ok {
		return nil, fmt.Errorf("org not found")
	}
	return o, nil
}

func (r *InMemoryOrganizationRepository) List(ctx context.Context) ([]*domain.Organization, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Organization
	for _, o := range r.orgs {
		list = append(list, o)
	}
	return list, nil
}

type InMemoryDepartmentRepository struct {
	mu    sync.RWMutex
	depts map[string]*domain.Department
}

func NewInMemoryDepartmentRepository() *InMemoryDepartmentRepository {
	return &InMemoryDepartmentRepository{depts: make(map[string]*domain.Department)}
}

func (r *InMemoryDepartmentRepository) Create(ctx context.Context, d *domain.Department) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.depts[d.ID] = d
	return nil
}

func (r *InMemoryDepartmentRepository) GetByID(ctx context.Context, id string) (*domain.Department, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.depts[id]
	if !ok {
		return nil, fmt.Errorf("dept not found")
	}
	return d, nil
}

func (r *InMemoryDepartmentRepository) GetByOrganization(ctx context.Context, orgID string) ([]*domain.Department, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Department
	for _, d := range r.depts {
		if d.OrganizationID == orgID {
			list = append(list, d)
		}
	}
	return list, nil
}

type InMemoryBareMetalRepository struct {
	mu    sync.RWMutex
	nodes map[string]*domain.BareMetalNode
	logs  map[string][]*domain.ProvisioningLog
}

func NewInMemoryBareMetalRepository() *InMemoryBareMetalRepository {
	return &InMemoryBareMetalRepository{
		nodes: make(map[string]*domain.BareMetalNode),
		logs:  make(map[string][]*domain.ProvisioningLog),
	}
}

func (r *InMemoryBareMetalRepository) Create(ctx context.Context, n *domain.BareMetalNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[n.ID] = n
	return nil
}

func (r *InMemoryBareMetalRepository) GetByID(ctx context.Context, id string) (*domain.BareMetalNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.nodes[id]
	if !ok {
		return nil, fmt.Errorf("node not found")
	}
	return n, nil
}

func (r *InMemoryBareMetalRepository) List(ctx context.Context) ([]*domain.BareMetalNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.BareMetalNode
	for _, n := range r.nodes {
		list = append(list, n)
	}
	return list, nil
}

func (r *InMemoryBareMetalRepository) Update(ctx context.Context, n *domain.BareMetalNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[n.ID] = n
	return nil
}

func (r *InMemoryBareMetalRepository) AddLog(ctx context.Context, l *domain.ProvisioningLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs[l.NodeID] = append(r.logs[l.NodeID], l)
	return nil
}

func (r *InMemoryBareMetalRepository) GetLogs(ctx context.Context, nodeID string) ([]*domain.ProvisioningLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.logs[nodeID], nil
}
