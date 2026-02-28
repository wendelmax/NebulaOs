import React, { useState, useEffect } from 'react';
import { Server, Trash2, ExternalLink, RefreshCw, Filter, Search } from 'lucide-react';

const ResourcesView: React.FC = () => {
    const [resources, setResources] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchResources = async () => {
        setLoading(true);
        try {
            const resp = await fetch('http://api.nebula.local/resources?project_id=v-p1');
            if (resp.ok) {
                const data = await resp.json();
                setResources(data || []);
            }
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchResources();
    }, []);

    const getStatusColor = (state: string) => {
        switch (state.toLowerCase()) {
            case 'active': return '#10b981';
            case 'provisioning': return '#3b82f6';
            case 'deleted': return '#f43f5e';
            default: return 'var(--text-muted)';
        }
    };

    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <h1>Infrastructure Resources</h1>
                    <p className="text-muted">Manage and monitor your multi-cloud compute assets.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <button className="btn-secondary" onClick={fetchResources} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <RefreshCw size={18} className={loading === true ? 'animate-spin' : ''} />
                        Refresh
                    </button>
                    <button className="btn-primary">Provision Resource</button>
                </div>
            </header>

            <div className="glass p-4 mb-8 mt-8" style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                <div style={{ position: 'relative', flex: 1 }}>
                    <Search size={18} style={{ position: 'absolute', left: '1rem', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }} />
                    <input
                        type="text"
                        placeholder="Search resources by ID, name or provider..."
                        style={{ width: '100%', padding: '0.75rem 1rem 0.75rem 3rem', background: 'var(--bg-accent)', border: '1px solid var(--glass-border)', borderRadius: '12px', color: 'var(--text-main)' }}
                    />
                </div>
                <button className="btn-secondary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                    <Filter size={18} />
                    Filter
                </button>
            </div>

            <div className="glass" style={{ overflow: 'hidden' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                    <thead style={{ background: 'var(--bg-accent)' }}>
                        <tr>
                            <th style={{ textAlign: 'left', padding: '1.25rem', fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-muted)', textTransform: 'uppercase' }}>Type / ID</th>
                            <th style={{ textAlign: 'left', padding: '1.25rem', fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-muted)', textTransform: 'uppercase' }}>Resource Name</th>
                            <th style={{ textAlign: 'left', padding: '1.25rem', fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-muted)', textTransform: 'uppercase' }}>Provider</th>
                            <th style={{ textAlign: 'left', padding: '1.25rem', fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-muted)', textTransform: 'uppercase' }}>Status</th>
                            <th style={{ textAlign: 'right', padding: '1.25rem', fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-muted)', textTransform: 'uppercase' }}>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {resources.length === 0 && !loading && (
                            <tr>
                                <td colSpan={5} style={{ padding: '4rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                                    No resources found for this project.
                                </td>
                            </tr>
                        )}
                        {resources.map((res) => (
                            <tr key={res.id} className="glass-hover" style={{ borderBottom: '1px solid var(--glass-border)' }}>
                                <td style={{ padding: '1.25rem' }}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                                        <div className="stat-icon" style={{ padding: '0.5rem', background: 'var(--bg-accent)', color: 'var(--primary-light)' }}>
                                            <Server size={18} />
                                        </div>
                                        <div>
                                            <div style={{ fontWeight: 600, fontSize: '0.9rem' }}>{res.type}</div>
                                            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{res.id}</div>
                                        </div>
                                    </div>
                                </td>
                                <td style={{ padding: '1.25rem', fontWeight: 500, fontSize: '0.9rem' }}>{res.name}</td>
                                <td style={{ padding: '1.25rem' }}>
                                    <span style={{ fontSize: '0.8rem', background: 'rgba(255,255,255,0.05)', padding: '0.25rem 0.6rem', borderRadius: '4px' }}>
                                        {res.provider}
                                    </span>
                                </td>
                                <td style={{ padding: '1.25rem' }}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                        <div style={{ width: '8px', height: '8px', borderRadius: '50%', background: getStatusColor(res.state) }}></div>
                                        <span style={{ fontSize: '0.9rem', fontWeight: 500 }}>{res.state}</span>
                                    </div>
                                </td>
                                <td style={{ padding: '1.25rem', textAlign: 'right' }}>
                                    <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'flex-end' }}>
                                        <button className="glass-hover p-2 rounded-lg text-muted">
                                            <ExternalLink size={18} />
                                        </button>
                                        <button className="glass-hover p-2 rounded-lg text-danger">
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default ResourcesView;
