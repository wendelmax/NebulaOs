import React from 'react';
import { Box, Database, ShieldCheck, Zap } from 'lucide-react';

const Marketplace: React.FC = () => {
    const blueprints = [
        { id: 'bp-1', name: 'K8s Cluster', desc: 'Secure, multi-zone Kubernetes control plane.', cat: 'Infrastructure', icon: Box },
        { id: 'bp-2', name: 'Postgres High-Availability', desc: 'Managed DB with auto-failover and backups.', cat: 'Databases', icon: Database },
        { id: 'bp-3', name: 'Nebula Firewall Edge', desc: 'DDoS protection and advanced WAF.', cat: 'Security', icon: ShieldCheck },
        { id: 'bp-4', name: 'Redis Cache', desc: 'Low-latency in-memory data store.', cat: 'Cache', icon: Zap }
    ];

    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <h1>Cloud Marketplace</h1>
                    <p className="text-muted">Launch production-ready infrastructure blueprints in seconds.</p>
                </div>
            </header>

            <div className="stats-grid" style={{ marginTop: '2rem' }}>
                {blueprints.map(bp => {
                    const Icon = bp.icon;
                    return (
                        <div key={bp.id} className="stat-card glass-hover" style={{ cursor: 'pointer', display: 'flex', flexDirection: 'column', gap: '1rem', padding: '2rem' }}>
                            <div className="stat-icon" style={{ background: 'var(--primary-gradient)', color: 'white' }}>
                                <Icon size={24} />
                            </div>
                            <div>
                                <h3 style={{ fontSize: '1.25rem' }}>{bp.name}</h3>
                                <span className="stat-label" style={{ background: 'var(--bg-accent)', color: 'var(--primary-light)', padding: '0.25rem 0.5rem', borderRadius: '6px', fontSize: '0.7rem' }}>
                                    {bp.cat}
                                </span>
                            </div>
                            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: '1.6' }}>
                                {bp.desc}
                            </p>
                            <button className="btn-primary" style={{ marginTop: 'auto', padding: '0.75rem', fontSize: '0.9rem' }}>
                                Deploy Blueprint
                            </button>
                        </div>
                    );
                })}
            </div>
        </div>
    );
};

export default Marketplace;
