import React from 'react';
import { HardDrive, Database, Plus, Search } from 'lucide-react';

const StorageView: React.FC = () => {
    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
            <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <h1 style={{ fontSize: '1.875rem' }}>Storage Orchestration</h1>
                    <p style={{ color: 'var(--text-muted)', marginTop: '0.25rem' }}>Manage block volumes and object storage buckets across providers.</p>
                </div>
                <button className="button-primary">
                    <Plus size={20} />
                    <span>Create Resource</span>
                </button>
            </header>

            <div className="resource-grid">
                <div className="glass p-6">
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
                        <HardDrive className="text-primary" size={24} />
                        <span className="badge badge-success">Active</span>
                    </div>
                    <h4>Mock-Volume-01</h4>
                    <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Size: 100GB | Type: Block</p>
                </div>
                <div className="glass p-6">
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
                        <Database className="text-secondary" size={24} />
                        <span className="badge badge-success">Active</span>
                    </div>
                    <h4>Mock-Bucket-Alpha</h4>
                    <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Region: us-east-1 | Type: Object</p>
                </div>
            </div>

            <div className="glass" style={{ padding: '1.5rem' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '2rem' }}>
                    <h3>Storage Inventory</h3>
                    <div className="search-box">
                        <Search size={18} />
                        <input type="text" placeholder="Search volumes or buckets..." />
                    </div>
                </div>

                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                    <thead>
                        <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--glass-border)' }}>
                            <th style={{ padding: '1rem 0' }}>Name</th>
                            <th>Type</th>
                            <th>Status</th>
                            <th>Capacity</th>
                            <th>Provider</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr style={{ borderBottom: '1px solid var(--glass-border)' }}>
                            <td style={{ padding: '1rem 0' }}>block-vol-001</td>
                            <td>Block Storage</td>
                            <td><span className="badge badge-success">online</span></td>
                            <td>250 GB</td>
                            <td>Mock-Storage</td>
                        </tr>
                        <tr>
                            <td style={{ padding: '1rem 0' }}>object-store-beta</td>
                            <td>Object (S3)</td>
                            <td><span className="badge badge-success">ready</span></td>
                            <td>Unlimited</td>
                            <td>Mock-Storage</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default StorageView;
