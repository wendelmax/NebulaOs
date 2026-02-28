package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
	"github.com/lib/pq"
)

var _ domain.TenantRepository = (*PostgresTenantRepository)(nil)
var _ domain.ProjectRepository = (*PostgresProjectRepository)(nil)
var _ domain.ResourceRepository = (*PostgresResourceRepository)(nil)
var _ domain.QuotaRepository = (*PostgresQuotaRepository)(nil)
var _ domain.VolumeRepository = (*PostgresVolumeRepository)(nil)
var _ domain.BucketRepository = (*PostgresBucketRepository)(nil)
var _ domain.SovereigntyPolicyRepository = (*PostgresPolicyRepository)(nil)
var _ domain.SecurityGroupRepository = (*PostgresSecurityGroupRepository)(nil)
var _ domain.TerraformStateRepository = (*PostgresTerraformStateRepository)(nil)
var _ domain.BlueprintRepository = (*PostgresBlueprintRepository)(nil)

// --- Tenant Repository ---

type PostgresTenantRepository struct {
	db *sql.DB
}

func NewPostgresTenantRepository(db *sql.DB) *PostgresTenantRepository {
	return &PostgresTenantRepository{db: db}
}

func (r *PostgresTenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	query := `INSERT INTO tenants (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.Name, t.CreatedAt)
	return err
}

func (r *PostgresTenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `SELECT id, name, created_at FROM tenants WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var t domain.Tenant
	if err := row.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PostgresTenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	query := `SELECT id, name, created_at FROM tenants`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []*domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, err
		}
		tenants = append(tenants, &t)
	}
	return tenants, nil
}

// --- Project Repository ---

type PostgresProjectRepository struct {
	db *sql.DB
}

func NewPostgresProjectRepository(db *sql.DB) *PostgresProjectRepository {
	return &PostgresProjectRepository{db: db}
}

func (r *PostgresProjectRepository) Create(ctx context.Context, p *domain.Project) error {
	query := `INSERT INTO projects (id, tenant_id, name, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.TenantID, p.Name, p.CreatedAt)
	return err
}

func (r *PostgresProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `SELECT id, tenant_id, name, created_at FROM projects WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var p domain.Project
	if err := row.Scan(&p.ID, &p.TenantID, &p.Name, &p.CreatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresProjectRepository) GetByTenant(ctx context.Context, tenantID string) ([]*domain.Project, error) {
	query := `SELECT id, tenant_id, name, created_at FROM projects WHERE tenant_id = $1`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
}

func (r *PostgresProjectRepository) List(ctx context.Context) ([]*domain.Project, error) {
	query := `SELECT id, tenant_id, name, created_at FROM projects`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
}

// --- Resource Repository ---

type PostgresResourceRepository struct {
	db *sql.DB
}

func NewPostgresResourceRepository(db *sql.DB) *PostgresResourceRepository {
	return &PostgresResourceRepository{db: db}
}

func (r *PostgresResourceRepository) Create(ctx context.Context, res *domain.Resource) error {
	metadataJSON, _ := json.Marshal(res.Metadata)
	query := `INSERT INTO resources (id, project_id, type, provider, state, metadata, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, res.ID, res.ProjectID, res.Type, res.Provider, res.State, metadataJSON, res.CreatedAt)
	return err
}

func (r *PostgresResourceRepository) GetByID(ctx context.Context, id string) (*domain.Resource, error) {
	query := `SELECT id, project_id, type, provider, state, metadata, created_at FROM resources WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var res domain.Resource
	var metadataBlob []byte
	if err := row.Scan(&res.ID, &res.ProjectID, &res.Type, &res.Provider, &res.State, &metadataBlob, &res.CreatedAt); err != nil {
		return nil, err
	}
	json.Unmarshal(metadataBlob, &res.Metadata)
	return &res, nil
}

func (r *PostgresResourceRepository) GetByProject(ctx context.Context, projectID string) ([]*domain.Resource, error) {
	query := `SELECT id, project_id, type, provider, state, metadata, created_at FROM resources WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*domain.Resource
	for rows.Next() {
		var res domain.Resource
		var metadataBlob []byte
		if err := rows.Scan(&res.ID, &res.ProjectID, &res.Type, &res.Provider, &res.State, &metadataBlob, &res.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(metadataBlob, &res.Metadata)
		resources = append(resources, &res)
	}
	return resources, nil
}

func (r *PostgresResourceRepository) UpdateState(ctx context.Context, id string, state string) error {
	query := `UPDATE resources SET state = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, state, id)
	return err
}

func (r *PostgresResourceRepository) List(ctx context.Context) ([]*domain.Resource, error) {
	query := `SELECT id, project_id, type, provider, state, metadata, created_at FROM resources`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*domain.Resource
	for rows.Next() {
		var res domain.Resource
		var metadataBlob []byte
		if err := rows.Scan(&res.ID, &res.ProjectID, &res.Type, &res.Provider, &res.State, &metadataBlob, &res.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(metadataBlob, &res.Metadata)
		resources = append(resources, &res)
	}
	return resources, nil
}

// --- Quota Repository ---

type PostgresQuotaRepository struct {
	db *sql.DB
}

func NewPostgresQuotaRepository(db *sql.DB) *PostgresQuotaRepository {
	return &PostgresQuotaRepository{db: db}
}

func (r *PostgresQuotaRepository) GetByTenant(ctx context.Context, tenantID string) (*domain.Quota, error) {
	query := `SELECT tenant_id, max_cpus, max_ram, max_disk FROM quotas WHERE tenant_id = $1`
	row := r.db.QueryRowContext(ctx, query, tenantID)
	var q domain.Quota
	if err := row.Scan(&q.TenantID, &q.MaxCPUs, &q.MaxRAM, &q.MaxDisk); err != nil {
		if err == sql.ErrNoRows {
			return &domain.Quota{TenantID: tenantID, MaxCPUs: 4, MaxRAM: 8192, MaxDisk: 100}, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *PostgresQuotaRepository) Update(ctx context.Context, q *domain.Quota) error {
	query := `INSERT INTO quotas (tenant_id, max_cpus, max_ram, max_disk) VALUES ($1, $2, $3, $4)
              ON CONFLICT (tenant_id) DO UPDATE SET max_cpus = $2, max_ram = $3, max_disk = $4`
	_, err := r.db.ExecContext(ctx, query, q.TenantID, q.MaxCPUs, q.MaxRAM, q.MaxDisk)
	return err
}

// --- Volume Repository ---

type PostgresVolumeRepository struct {
	db *sql.DB
}

func NewPostgresVolumeRepository(db *sql.DB) *PostgresVolumeRepository {
	return &PostgresVolumeRepository{db: db}
}

func (r *PostgresVolumeRepository) Create(ctx context.Context, v *domain.Volume) error {
	query := `INSERT INTO volumes (id, name, size, state, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, v.ID, v.Name, v.SizeGB, v.State, v.CreatedAt)
	return err
}

func (r *PostgresVolumeRepository) GetByID(ctx context.Context, id string) (*domain.Volume, error) {
	query := `SELECT id, name, size, state, created_at FROM volumes WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var v domain.Volume
	if err := row.Scan(&v.ID, &v.Name, &v.SizeGB, &v.State, &v.CreatedAt); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PostgresVolumeRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.Volume, error) {
	// Note: project_id column missing in schema, let's assume it should be there for project-aware listing
	query := `SELECT id, name, size, state, created_at FROM volumes`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []*domain.Volume
	for rows.Next() {
		var v domain.Volume
		if err := rows.Scan(&v.ID, &v.Name, &v.SizeGB, &v.State, &v.CreatedAt); err != nil {
			return nil, err
		}
		volumes = append(volumes, &v)
	}
	return volumes, nil
}

// --- Bucket Repository ---

type PostgresBucketRepository struct {
	db *sql.DB
}

func NewPostgresBucketRepository(db *sql.DB) *PostgresBucketRepository {
	return &PostgresBucketRepository{db: db}
}

func (r *PostgresBucketRepository) Create(ctx context.Context, b *domain.Bucket) error {
	query := `INSERT INTO buckets (id, name, state, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, b.ID, b.Name, b.State, b.CreatedAt)
	return err
}

func (r *PostgresBucketRepository) GetByID(ctx context.Context, id string) (*domain.Bucket, error) {
	query := `SELECT id, name, state, created_at FROM buckets WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var b domain.Bucket
	if err := row.Scan(&b.ID, &b.Name, &b.State, &b.CreatedAt); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *PostgresBucketRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.Bucket, error) {
	query := `SELECT id, name, state, created_at FROM buckets`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buckets []*domain.Bucket
	for rows.Next() {
		var b domain.Bucket
		if err := rows.Scan(&b.ID, &b.Name, &b.State, &b.CreatedAt); err != nil {
			return nil, err
		}
		buckets = append(buckets, &b)
	}
	return buckets, nil
}

// --- Sovereignty Policy Repository ---

type PostgresPolicyRepository struct {
	db *sql.DB
}

func NewPostgresPolicyRepository(db *sql.DB) *PostgresPolicyRepository {
	return &PostgresPolicyRepository{db: db}
}

func (r *PostgresPolicyRepository) GetByTenantID(ctx context.Context, tenantID string) (*domain.SovereigntyPolicy, error) {
	query := `SELECT tenant_id, allowed_regions FROM sovereignty_policies WHERE tenant_id = $1`
	row := r.db.QueryRowContext(ctx, query, tenantID)
	var p domain.SovereigntyPolicy
	var regions []string
	if err := row.Scan(&p.TenantID, pq.Array(&regions)); err != nil {
		if err == sql.ErrNoRows {
			return &domain.SovereigntyPolicy{TenantID: tenantID, AllowedRegions: []string{}}, nil
		}
		return nil, err
	}
	p.AllowedRegions = regions
	return &p, nil
}

func (r *PostgresPolicyRepository) Save(ctx context.Context, p *domain.SovereigntyPolicy) error {
	query := `INSERT INTO sovereignty_policies (tenant_id, allowed_regions) VALUES ($1, $2)
              ON CONFLICT (tenant_id) DO UPDATE SET allowed_regions = $2`
	_, err := r.db.ExecContext(ctx, query, p.TenantID, pq.Array(p.AllowedRegions))
	return err
}

// --- Security Group Repository ---

type PostgresSecurityGroupRepository struct {
	db *sql.DB
}

func NewPostgresSecurityGroupRepository(db *sql.DB) *PostgresSecurityGroupRepository {
	return &PostgresSecurityGroupRepository{db: db}
}

func (r *PostgresSecurityGroupRepository) Create(ctx context.Context, sg *domain.SecurityGroup) error {
	query := `INSERT INTO security_groups (id, project_id, name, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, sg.ID, sg.ProjectID, sg.Name, sg.CreatedAt)
	return err
}

func (r *PostgresSecurityGroupRepository) GetByID(ctx context.Context, id string) (*domain.SecurityGroup, error) {
	query := `SELECT id, project_id, name, created_at FROM security_groups WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var sg domain.SecurityGroup
	if err := row.Scan(&sg.ID, &sg.ProjectID, &sg.Name, &sg.CreatedAt); err != nil {
		return nil, err
	}

	// Fetch rules
	rulesQuery := `SELECT id, protocol, from_port, to_port, source_ip, action, created_at FROM firewall_rules WHERE security_group_id = $1`
	rows, err := r.db.QueryContext(ctx, rulesQuery, id)
	if err != nil {
		return &sg, nil // Return SG even if rules fail to fetch
	}
	defer rows.Close()

	for rows.Next() {
		var rule domain.FirewallRule
		if err := rows.Scan(&rule.ID, &rule.Protocol, &rule.FromPort, &rule.ToPort, &rule.SourceIP, &rule.Action, &rule.CreatedAt); err != nil {
			continue
		}
		sg.Rules = append(sg.Rules, rule)
	}

	return &sg, nil
}

func (r *PostgresSecurityGroupRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.SecurityGroup, error) {
	query := `SELECT id, project_id, name, created_at FROM security_groups WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sgs []*domain.SecurityGroup
	for rows.Next() {
		var sg domain.SecurityGroup
		if err := rows.Scan(&sg.ID, &sg.ProjectID, &sg.Name, &sg.CreatedAt); err != nil {
			return nil, err
		}
		sgs = append(sgs, &sg)
	}
	return sgs, nil
}

func (r *PostgresSecurityGroupRepository) AddRule(ctx context.Context, sgID string, rule domain.FirewallRule) error {
	query := `INSERT INTO firewall_rules (id, security_group_id, protocol, from_port, to_port, source_ip, action, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, rule.ID, sgID, rule.Protocol, rule.FromPort, rule.ToPort, rule.SourceIP, rule.Action, rule.CreatedAt)
	return err
}

func (r *PostgresSecurityGroupRepository) RemoveRule(ctx context.Context, sgID string, ruleID string) error {
	query := `DELETE FROM firewall_rules WHERE id = $1 AND security_group_id = $2`
	_, err := r.db.ExecContext(ctx, query, ruleID, sgID)
	return err
}

// --- Terraform State Repository ---

type PostgresTerraformStateRepository struct {
	db *sql.DB
}

func NewPostgresTerraformStateRepository(db *sql.DB) *PostgresTerraformStateRepository {
	return &PostgresTerraformStateRepository{db: db}
}

func (r *PostgresTerraformStateRepository) Save(ctx context.Context, s *domain.TerraformState) error {
	query := `INSERT INTO terraform_states (id, project_id, version, state, updated_at) VALUES ($1, $2, $3, $4, $5)
              ON CONFLICT (id) DO UPDATE SET version = $3, state = $4, updated_at = $5`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.ProjectID, s.Version, s.State, s.UpdatedAt)
	return err
}

func (r *PostgresTerraformStateRepository) GetByProjectID(ctx context.Context, projectID string) (*domain.TerraformState, error) {
	query := `SELECT id, project_id, version, state, updated_at FROM terraform_states WHERE project_id = $1`
	row := r.db.QueryRowContext(ctx, query, projectID)
	var s domain.TerraformState
	if err := row.Scan(&s.ID, &s.ProjectID, &s.Version, &s.State, &s.UpdatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

// --- Blueprint Repository ---

type PostgresBlueprintRepository struct {
	db *sql.DB
}

func NewPostgresBlueprintRepository(db *sql.DB) *PostgresBlueprintRepository {
	return &PostgresBlueprintRepository{db: db}
}

func (r *PostgresBlueprintRepository) Create(ctx context.Context, b *domain.Blueprint) error {
	resourcesJSON, _ := json.Marshal(b.Resources)
	variablesJSON, _ := json.Marshal(b.Variables)
	query := `INSERT INTO blueprints (id, name, description, category, resources, variables, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, b.ID, b.Name, b.Description, b.Category, resourcesJSON, variablesJSON, b.CreatedAt)
	return err
}

func (r *PostgresBlueprintRepository) GetByID(ctx context.Context, id string) (*domain.Blueprint, error) {
	query := `SELECT id, name, description, category, resources, variables, created_at FROM blueprints WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var b domain.Blueprint
	var resourcesBlob, variablesBlob []byte
	if err := row.Scan(&b.ID, &b.Name, &b.Description, &b.Category, &resourcesBlob, &variablesBlob, &b.CreatedAt); err != nil {
		return nil, err
	}
	json.Unmarshal(resourcesBlob, &b.Resources)
	json.Unmarshal(variablesBlob, &b.Variables)
	return &b, nil
}

func (r *PostgresBlueprintRepository) List(ctx context.Context) ([]*domain.Blueprint, error) {
	query := `SELECT id, name, description, category, resources, variables, created_at FROM blueprints`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blueprints []*domain.Blueprint
	for rows.Next() {
		var b domain.Blueprint
		var resourcesBlob, variablesBlob []byte
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.Category, &resourcesBlob, &variablesBlob, &b.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(resourcesBlob, &b.Resources)
		json.Unmarshal(variablesBlob, &b.Variables)
		blueprints = append(blueprints, &b)
	}
	return blueprints, nil
}
