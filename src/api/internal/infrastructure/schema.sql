-- NebulaOS Production Schema

CREATE TABLE IF NOT EXISTS tenants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    tenant_id TEXT REFERENCES tenants(id),
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS resources (
    id TEXT PRIMARY KEY,
    project_id TEXT REFERENCES projects(id),
    type TEXT NOT NULL,
    provider TEXT NOT NULL,
    state TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    blueprint_id TEXT,
    security_groups TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS quotas (
    tenant_id TEXT PRIMARY KEY REFERENCES tenants(id),
    max_cpus INTEGER NOT NULL,
    max_ram INTEGER NOT NULL,
    max_disk INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS volumes (
    id TEXT PRIMARY KEY,
    project_id TEXT REFERENCES projects(id),
    name TEXT NOT NULL,
    size INTEGER NOT NULL,
    state TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS buckets (
    id TEXT PRIMARY KEY,
    project_id TEXT REFERENCES projects(id),
    name TEXT NOT NULL,
    state TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sovereignty_policies (
    tenant_id VARCHAR(255) PRIMARY KEY REFERENCES tenants(id),
    allowed_regions TEXT[] NOT NULL
);

CREATE TABLE IF NOT EXISTS security_groups (
    id VARCHAR(255) PRIMARY KEY,
    project_id VARCHAR(255) NOT NULL REFERENCES projects(id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS firewall_rules (
    id VARCHAR(255) PRIMARY KEY,
    security_group_id VARCHAR(255) NOT NULL REFERENCES security_groups(id),
    protocol VARCHAR(50) NOT NULL,
    from_port INTEGER NOT NULL,
    to_port INTEGER NOT NULL,
    source_ip VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS terraform_states (
    id VARCHAR(255) PRIMARY KEY,
    project_id VARCHAR(255) NOT NULL REFERENCES projects(id),
    version INTEGER NOT NULL DEFAULT 1,
    state BYTEA NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS blueprints (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    resources JSONB NOT NULL,
    variables JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT,
    tenant_id TEXT REFERENCES tenants(id),
    must_change_password BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS departments (
    id TEXT PRIMARY KEY,
    organization_id TEXT REFERENCES organizations(id),
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Update projects to support department hierarchy
ALTER TABLE projects ADD COLUMN IF NOT EXISTS department_id TEXT;

CREATE TABLE IF NOT EXISTS bare_metal_nodes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    mac TEXT UNIQUE NOT NULL,
    ipmi_address TEXT,
    ipmi_user TEXT,
    ipmi_password TEXT,
    cpu_cores INTEGER,
    memory_gb INTEGER,
    disk_gb INTEGER,
    department_id TEXT,
    state TEXT NOT NULL,
    provider_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS provisioning_logs (
    id TEXT PRIMARY KEY,
    node_id TEXT REFERENCES bare_metal_nodes(id),
    message TEXT NOT NULL,
    level TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
