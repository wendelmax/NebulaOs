import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';
const AUTH_STORAGE_KEY = import.meta.env.VITE_AUTH_STORAGE_KEY || 'nebula_token';

export interface LoginRequest {
    username: string;
    password: string;
}

export interface ChangePasswordRequest {
    old_password: string;
    new_password: string;
    email: string;
}

export interface CreateTenantRequest {
    name: string;
    id?: string;
}

export interface CreateProjectRequest {
    name: string;
    tenant_id: string;
    id?: string;
}

export interface CreateResourceRequest {
    name: string;
    type: string;
    provider: string;
    project_id: string;
    metadata?: Record<string, unknown>;
}

export interface InterProviderRouteRequest {
    route: {
        destination: string;
        next_hop: string;
    };
    source_provider: string;
    target_provider: string;
}

export interface OrchestrateStackRequest {
    project_id: string;
    regions: string[];
}

export interface ProvisionPresetRequest {
    preset_id: string;
    project_id: string;
    variables?: Record<string, unknown>;
}

export interface RegisterProviderRequest {
    name: string;
    type: string;
    endpoint: string;
    id?: string;
}

export interface DeployBlueprintRequest {
    blueprint_id: string;
    project_id: string;
    variables?: Record<string, unknown>;
}

const client = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
});

client.interceptors.request.use((config) => {
    const token = localStorage.getItem(AUTH_STORAGE_KEY);
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const api = {
    getTenants: () => client.get('/tenants'),
    createTenant: (data: CreateTenantRequest) => client.post('/tenants', data),
    getProjects: () => client.get('/projects'),
    createProject: (data: CreateProjectRequest) => client.post('/projects', data),

    login: (data: LoginRequest) => client.post('/auth/login', data),
    logout: () => client.post('/auth/logout'),
    checkAuth: () => client.get('/auth/me'),
    changePassword: (data: ChangePasswordRequest) => client.post('/auth/change-password', data),
    getResources: (projectId?: string) => client.get(`/resources${projectId ? `?project_id=${projectId}` : ''}`),
    createResource: (data: CreateResourceRequest) => client.post('/resources', data),
    getNetworkStatus: (domain: string) => client.get(`/network/certificate?domain=${domain}`),
    getStats: () => client.get('/intelligence/stats'),

    getPresets: () => client.get('/infra/automated/provision'),
    provisionPreset: (data: ProvisionPresetRequest) =>
        client.post('/infra/automated/provision', data),
    getRegions: () => client.get('/cloud/regions'),
    getZones: (regionId?: string) => client.get(`/cloud/zones${regionId ? `?region_id=${regionId}` : ''}`),
    createInterProviderRoute: (data: InterProviderRouteRequest) => client.post('/network/routes', data),
    orchestrateStack: (data: OrchestrateStackRequest) => client.post('/cloud/orchestrate/stack', data),

    getAdvisorInsights: (projectId: string) => client.get(`/intelligence/advisor?project_id=${projectId}`),
    getPolicy: (tenantId: string) => client.get(`/governance/policy?tenant_id=${tenantId}`),

    getProviders: () => client.get('/api/providers'),
    registerProvider: (data: RegisterProviderRequest) => client.post('/api/providers', data),
    deleteProvider: (id: string) => client.delete(`/api/providers?id=${id}`),

    getVolumes: (projectId: string) => client.get(`/storage/volumes?project_id=${projectId}`),
    getBuckets: (projectId: string) => client.get(`/storage/buckets?project_id=${projectId}`),

    getBillingReport: (tenantId: string) => client.get(`/billing/report?tenant_id=${tenantId}`),
    getGlobalStats: () => client.get('/api/billing/stats'),

    getOrganizations: () => client.get('/api/organizations'),
    getDepartments: (orgId?: string) => client.get(`/api/departments${orgId ? `?organization_id=${orgId}` : ''}`),

    getBareMetalNodes: () => client.get('/api/baremetal/nodes'),
    provisionBareMetalNode: (id: string) => client.post(`/api/baremetal/provision?id=${id}`, {}),
    getBareMetalLogs: (nodeId: string) => client.get(`/api/baremetal/logs?node_id=${nodeId}`),

    getBlueprints: () => client.get('/marketplace/blueprints'),
    deployBlueprint: (data: DeployBlueprintRequest) => client.post('/marketplace/deploy', data),

    getSecurityGroups: (projectId: string) => client.get(`/security-groups?project_id=${projectId}`),
};

export const AUTH_TOKEN_KEY = AUTH_STORAGE_KEY;

export function getAuthToken(): string | null {
    return localStorage.getItem(AUTH_STORAGE_KEY);
}

export function setAuthToken(token: string): void {
    localStorage.setItem(AUTH_STORAGE_KEY, token);
}

export function clearAuthToken(): void {
    localStorage.removeItem(AUTH_STORAGE_KEY);
}

export default client;