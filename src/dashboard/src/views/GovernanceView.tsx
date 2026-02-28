import React from 'react';
import { ShieldCheck, FileText } from 'lucide-react';

const GovernanceView: React.FC = () => {
    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
            <header>
                <h1 style={{ fontSize: '1.875rem' }}>Governance & Compliance</h1>
                <p style={{ color: 'var(--text-muted)', marginTop: '0.25rem' }}>Enforcing policy, quotas, and sovereign audit standards.</p>
            </header>

            <div className="resource-grid">
                <div className="glass p-6">
                    <h3>CPU Quota</h3>
                    <div style={{ marginTop: '1rem' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                            <span style={{ fontSize: '0.875rem' }}>Usage: 14 / 20 vCPUs</span>
                            <span style={{ fontSize: '0.875rem', fontWeight: 600 }}>70%</span>
                        </div>
                        <div style={{ height: '8px', background: 'rgba(255,255,255,0.1)', borderRadius: '4px' }}>
                            <div style={{ height: '100%', width: '70%', background: 'var(--primary)', borderRadius: '4px' }} />
                        </div>
                    </div>
                </div>
                <div className="glass p-6">
                    <h3>RAM Quota</h3>
                    <div style={{ marginTop: '1rem' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                            <span style={{ fontSize: '0.875rem' }}>Usage: 32 / 64 GB</span>
                            <span style={{ fontSize: '0.875rem', fontWeight: 600 }}>50%</span>
                        </div>
                        <div style={{ height: '8px', background: 'rgba(255,255,255,0.1)', borderRadius: '4px' }}>
                            <div style={{ height: '100%', width: '50%', background: 'var(--secondary)', borderRadius: '4px' }} />
                        </div>
                    </div>
                </div>
                <div className="glass p-6" style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <ShieldCheck color="var(--success)" size={32} />
                    <div>
                        <h4 style={{ margin: 0 }}>Policy Compliance</h4>
                        <p style={{ margin: 0, fontSize: '0.875rem', color: 'var(--text-muted)' }}>Status: Fully Compliant</p>
                    </div>
                </div>
            </div>

            <div className="glass" style={{ padding: '1.5rem' }}>
                <div style={{ display: 'flex', gap: '1rem', alignItems: 'center', marginBottom: '1.5rem' }}>
                    <FileText />
                    <h3>Detailed Audit Logs</h3>
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                    {[1, 2, 3, 4, 5].map(i => (
                        <div key={i} className="glass-hover" style={{ padding: '1rem', display: 'grid', gridTemplateColumns: '150px 1fr 100px', alignItems: 'center', fontSize: '0.875rem' }}>
                            <span style={{ color: 'var(--text-muted)' }}>2026-02-28 14:0{i}</span>
                            <span style={{ fontWeight: 500 }}>API_REQUEST: POST /resources (TenantID: G-NGO-DE)</span>
                            <span style={{ textAlign: 'right', color: 'var(--primary)' }}>SUCCESS</span>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default GovernanceView;
