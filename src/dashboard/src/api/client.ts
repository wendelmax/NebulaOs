import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';
const AUTH_STORAGE_KEY = import.meta.env.VITE_AUTH_STORAGE_KEY || 'nebula_token';

const client = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
});

// Deprecated: kept for backward compat with old tokens
client.interceptors.request.use((config) => {
    const token = localStorage.getItem(AUTH_STORAGE_KEY);
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const api = {
    getTenants: () => client.get('/tenants'),
    createTenant: (data: any) => client.post('/tenants', data),
    getProjects: () => client.get('/projects'),
    createProject: (data: any) => client.post('/projects', data),

    // Auth (cookie-based)
    login: (data: any) => client.post('/auth/login', data),
    logout: () => client.post('/auth/logout'),
    checkAuth: () => client.get('/auth/me'),
    changePassword: (data: any) => client.post('/auth/change-password', data),
    getResources: (projectId?: string) => client.get(`/resources${projectId ? `?project_id=${projectId}` : ''}`),
    createResource: (data: any) => client.post('/resources', data),
    getNetworkStatus: (domain: string) => client.get(`/network/certificate?domain=${domain}`),
    getStats: () => client.get('/intelligence/stats'),
    
    // Open Cloud Extensions
    getPresets: () => client.get('/infra/automated/provision'),
    provisionPreset: (data: { preset_id: string, project_id: string, variables?: any }) => 
        client.post('/infra/automated/provision', data),
    getRegions: () => client.get('/cloud/regions'),
    getZones: (regionId?: string) => client.get(`/cloud/zones${regionId ? `?region_id=${regionId}` : ''}`),
    createInterProviderRoute: (data: any) => client.post('/network/routes', data),
    orchestrateStack: (data: any) => client.post('/cloud/orchestrate/stack', data),

    // Intelligence & Governance
    getAdvisorInsights: (projectId: string) => client.get(`/intelligence/advisor?project_id=${projectId}`),
    getPolicy: (tenantId: string) => client.get(`/governance/policy?tenant_id=${tenantId}`),

    // Provider Management
    getProviders: () => client.get('/api/providers'),
    registerProvider: (data: any) => client.post('/api/providers', data),
    deleteProvider: (id: string) => client.delete(`/api/providers?id=${id}`),

    // Storage
    getVolumes: (projectId: string) => client.get(`/storage/volumes?project_id=${projectId}`),
    getBuckets: (projectId: string) => client.get(`/storage/buckets?project_id=${projectId}`),

    // Billing
    getBillingReport: (tenantId: string) => client.get(`/billing/report?tenant_id=${tenantId}`),
    getGlobalStats: () => client.get('/api/billing/stats'),

    // Enterprise Hierarchy
    getOrganizations: () => client.get('/api/organizations'),
    getDepartments: (orgId?: string) => client.get(`/api/departments${orgId ? `?organization_id=${orgId}` : ''}`),

    // Bare Metal Orchestration
    getBareMetalNodes: () => client.get('/api/baremetal/nodes'),
    provisionBareMetalNode: (id: string) => client.post(`/api/baremetal/provision?id=${id}`, {}),
    getBareMetalLogs: (nodeId: string) => client.get(`/api/baremetal/logs?node_id=${nodeId}`),

    // Marketplace
    getBlueprints: () => client.get('/marketplace/blueprints'),
    deployBlueprint: (data: any) => client.post('/marketplace/deploy', data),

    // Security Groups
    getSecurityGroups: (projectId: string) => client.get(`/security-groups?project_id=${projectId}`),
};

export const AUTH_TOKEN_KEY = AUTH_STORAGE_KEY;

// Deprecated: prefer cookie-based auth
export function getAuthToken(): string | null {
    return localStorage.getItem(AUTH_STORAGE_KEY);
}

// Deprecated: prefer cookie-based auth
export function setAuthToken(token: string): void {
    localStorage.setItem(AUTH_STORAGE_KEY, token);
}

// Deprecated: prefer api.logout()
export function clearAuthToken(): void {
    localStorage.removeItem(AUTH_STORAGE_KEY);
}

export default client;