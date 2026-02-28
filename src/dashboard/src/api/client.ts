import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';

const client = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor for Auth
client.interceptors.request.use((config) => {
    const token = localStorage.getItem('nebula_token');
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
    getResources: () => client.get('/resources'),
    createResource: (data: any) => client.post('/resources', data),
    getNetworkStatus: (domain: string) => client.get(`/network/certificate?domain=${domain}`),
};

export default client;
