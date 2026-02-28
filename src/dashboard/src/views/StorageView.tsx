import React, { useState, useEffect } from 'react';
import { HardDrive, Database, Plus, Search, RefreshCw } from 'lucide-react';

const StorageView: React.FC = () => {
    const [volumes, setVolumes] = useState<any[]>([]);
    const [buckets, setBuckets] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [volResp, buckResp] = await Promise.all([
                fetch('http://api.nebula.local/storage/volumes?project_id=v-p1'),
                fetch('http://api.nebula.local/storage/buckets?project_id=v-p1')
            ]);

            if (volResp.ok) {
                const volData = await volResp.json();
                setVolumes(volData || []);
            }
            if (buckResp.ok) {
                const buckData = await buckResp.json();
                setBuckets(buckData || []);
            }
        } catch (err) {
            console.error("Failed to fetch storage data", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
            <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <h1 style={{ fontSize: '1.875rem' }}>Storage Orchestration</h1>
                    <p style={{ color: 'var(--text-muted)', marginTop: '0.25rem' }}>Manage block volumes and object storage buckets across providers.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <button className="btn-secondary" onClick={fetchData} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        Refresh
                    </button>
                    <button className="button-primary">
                        <Plus size={20} />
                        <span>Create Resource</span>
                    </button>
                </div>
            </header>

            <div className="resource-grid">
                {volumes.slice(0, 2).map(vol => (
                    <div key={vol.id} className="glass p-6">
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
                            <HardDrive className="text-primary" size={24} />
                            <span className="badge badge-success">{vol.state}</span>
                        </div>
                        <h4>{vol.name}</h4>
                        <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Size: {vol.size_gb}GB | Type: Block</p>
                    </div>
                ))}
                {buckets.slice(0, Math.max(0, 2 - volumes.length)).map(buck => (
                    <div key={buck.id} className="glass p-6">
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
                            <Database className="text-secondary" size={24} />
                            <span className="badge badge-success">{buck.state}</span>
                        </div>
                        <h4>{buck.name}</h4>
                        <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Region: {buck.region} | Type: Object</p>
                    </div>
                ))}
                {volumes.length === 0 && buckets.length === 0 && !loading && (
                    <div className="glass p-6 text-center" style={{ gridColumn: 'span 2' }}>
                        <p style={{ color: 'var(--text-muted)' }}>No high-priority storage resources detected.</p>
                    </div>
                )}
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
                            <th>Provider ID</th>
                        </tr>
                    </thead>
                    <tbody>
                        {volumes.map(vol => (
                            <tr key={vol.id} style={{ borderBottom: '1px solid var(--glass-border)' }}>
                                <td style={{ padding: '1rem 0' }}>{vol.name}</td>
                                <td>Block Storage</td>
                                <td><span className="badge badge-success">{vol.state}</span></td>
                                <td>{vol.size_gb} GB</td>
                                <td>{vol.provider_id || 'manual'}</td>
                            </tr>
                        ))}
                        {buckets.map(buck => (
                            <tr key={buck.id} style={{ borderBottom: '1px solid var(--glass-border)' }}>
                                <td style={{ padding: '1rem 0' }}>{buck.name}</td>
                                <td>Object (S3)</td>
                                <td><span className="badge badge-success">{buck.state}</span></td>
                                <td>Unlimited</td>
                                <td>{buck.region}</td>
                            </tr>
                        ))}
                        {volumes.length === 0 && buckets.length === 0 && !loading && (
                            <tr>
                                <td colSpan={5} style={{ padding: '4rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                                    No storage resources found for this project.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default StorageView;
