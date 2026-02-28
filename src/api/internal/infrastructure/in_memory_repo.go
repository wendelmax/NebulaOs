package infrastructure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
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

type SovereignBillingManager struct {
	resRepo    domain.ResourceRepository
	volRepo    domain.VolumeRepository
	bucketRepo domain.BucketRepository
	tenantRepo domain.TenantRepository
}

func NewSovereignBillingManager(resRepo domain.ResourceRepository, volRepo domain.VolumeRepository, bucketRepo domain.BucketRepository, tenantRepo domain.TenantRepository) *SovereignBillingManager {
	return &SovereignBillingManager{resRepo: resRepo, volRepo: volRepo, bucketRepo: bucketRepo, tenantRepo: tenantRepo}
}

func (m *SovereignBillingManager) GenerateReport(ctx context.Context, tenantID string) (*domain.BillingReport, error) {
	// Mock logic: Iterate through resources and sum costs
	report := &domain.BillingReport{
		TenantID:    tenantID,
		TotalCost:   0,
		Items:       []domain.BillingItem{},
		GeneratedAt: time.Now(),
	}

	// Calculate Resource/VM costs
	vms, _ := m.resRepo.List(ctx)
	for _, vm := range vms {
		item := domain.BillingItem{
			ResourceID: vm.ID,
			Type:       string(vm.Type),
			Cost:       15.50,
			Currency:   "USD",
		}
		report.Items = append(report.Items, item)
		report.TotalCost += item.Cost
	}

	// Calculate Volume costs
	vols, _ := m.volRepo.ListByProject(ctx, "default-project")
	for _, vol := range vols {
		item := domain.BillingItem{
			ResourceID: vol.ID,
			Type:       "VOLUME",
			Cost:       5.00,
			Currency:   "USD",
		}
		report.Items = append(report.Items, item)
		report.TotalCost += item.Cost
	}

	return report, nil
}

func (m *SovereignBillingManager) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	vms, _ := m.resRepo.List(ctx)
	// Calculate storage across all buckets/volumes
	tenants, _ := m.tenantRepo.List(ctx)

	stats := &domain.GlobalStats{
		TotalCPUs:     0,
		TotalStorage:  0,
		TotalEgress:   892.4, // Simulated egress
		ActiveTenants: len(tenants),
		TrendCPUs:     12.5,
		TrendStorage:  -4.2,
	}

	for _, vm := range vms {
		if vm.Type == domain.ComputeResource {
			stats.TotalCPUs += 2.0 // Simplified calculation
		}
	}

	// We'll iterate all resources of type storage
	for _, res := range vms {
		if res.Type == domain.StorageResource {
			stats.TotalStorage += 0.5 // Simplified 500GB per storage resource
		}
	}

	return stats, nil
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
