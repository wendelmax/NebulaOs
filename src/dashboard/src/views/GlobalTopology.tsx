import React from 'react';
import { MapPin } from 'lucide-react';

const GlobalTopology: React.FC = () => {
    const regions = [
        { id: 'us-east-1', name: 'US East (N. Virginia)', status: 'online', lat: '38.8', lng: '-77.1', load: 45 },
        { id: 'eu-west-1', name: 'Europe (Ireland)', status: 'degraded', lat: '53.3', lng: '-6.2', load: 88 },
        { id: 'ap-southeast-1', name: 'Asia Pacific (Singapore)', status: 'online', lat: '1.3', lng: '103.8', load: 12 }
    ];

    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <h1>Global Topology</h1>
                    <p className="text-muted">Live orchestration map across multi-cloud regions.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <div className="glass p-2 px-4" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.8rem' }}>
                        <div style={{ width: '8px', height: '8px', borderRadius: '50%', background: '#10b981' }}></div>
                        Global GSLB: Active
                    </div>
                </div>
            </header>

            <div className="glass p-8" style={{ marginTop: '2rem', height: '500px', display: 'flex', alignItems: 'center', justifyContent: 'center', position: 'relative', overflow: 'hidden' }}>
                {/* Simulated Map Background */}
                <div style={{ position: 'absolute', inset: 0, opacity: 0.1, background: 'radial-gradient(circle at 50% 50%, var(--primary-main) 0%, transparent 70%)' }}></div>

                {/* Region Nodes */}
                <div style={{ position: 'relative', width: '100%', height: '100%' }}>
                    {regions.map((reg, idx) => (
                        <div key={reg.id}
                            className="stat-card glass-hover animate-pulse-subtle"
                            style={{
                                position: 'absolute',
                                top: `${30 + idx * 20}%`,
                                left: `${20 + idx * 25}%`,
                                width: '220px',
                                padding: '1rem',
                                zIndex: 10
                            }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.5rem' }}>
                                <MapPin size={18} className={reg.status === 'online' ? 'text-primary' : 'text-danger'} />
                                <span style={{ fontWeight: 700, fontSize: '0.9rem' }}>{reg.id}</span>
                            </div>
                            <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>{reg.name}</div>
                            <div style={{ marginTop: '0.75rem' }}>
                                <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.7rem', marginBottom: '0.25rem' }}>
                                    <span>Regional Load</span>
                                    <span>{reg.load}%</span>
                                </div>
                                <div style={{ width: '100%', height: '4px', background: 'var(--bg-accent)', borderRadius: '2px' }}>
                                    <div style={{ width: `${reg.load}%`, height: '100%', background: reg.load > 80 ? 'var(--danger-main)' : 'var(--primary-main)', borderRadius: '2px' }}></div>
                                </div>
                            </div>
                        </div>
                    ))}

                    {/* Simulated Connections */}
                    <svg style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', pointerEvents: 'none', opacity: 0.2 }}>
                        <line x1="30%" y1="40%" x2="55%" y2="60%" stroke="var(--primary-main)" strokeWidth="1" strokeDasharray="4 4" />
                        <line x1="55%" y1="60%" x2="80%" y2="80%" stroke="var(--primary-main)" strokeWidth="1" strokeDasharray="4 4" />
                    </svg>
                </div>
            </div>
        </div>
    );
};

export default GlobalTopology;
