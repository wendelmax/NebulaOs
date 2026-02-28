import React from 'react';
import { Activity, Zap, HeartPulse, BarChart3 } from 'lucide-react';

const ObservabilityView: React.FC = () => {
    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
            <header>
                <h1 style={{ fontSize: '1.875rem' }}>System Observability</h1>
                <p style={{ color: 'var(--text-muted)', marginTop: '0.25rem' }}>Real-time telemetry and health diagnostics for your cloud plane.</p>
            </header>

            <div className="resource-grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))' }}>
                <div className="glass p-6" style={{ textAlign: 'center' }}>
                    <HeartPulse color="var(--success)" style={{ margin: '0 auto 1rem' }} />
                    <h2 style={{ margin: '0.5rem 0' }}>100%</h2>
                    <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Platform Uptime</p>
                </div>
                <div className="glass p-6" style={{ textAlign: 'center' }}>
                    <Zap color="#fbbf24" style={{ margin: '0 auto 1rem' }} />
                    <h2 style={{ margin: '0.5rem 0' }}>42ms</h2>
                    <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Avg API Latency</p>
                </div>
                <div className="glass p-6" style={{ textAlign: 'center' }}>
                    <BarChart3 color="var(--primary)" style={{ margin: '0 auto 1rem' }} />
                    <h2 style={{ margin: '0.5rem 0' }}>8.2k</h2>
                    <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Requests / min</p>
                </div>
            </div>

            <div className="glass" style={{ padding: '2rem', minHeight: '400px', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
                <Activity size={48} color="var(--primary)" style={{ marginBottom: '1.5rem', opacity: 0.5 }} />
                <h3>Live Telemetry Feed</h3>
                <p style={{ color: 'var(--text-muted)', maxWidth: '400px', textAlign: 'center', marginTop: '1rem' }}>
                    Establishing secure socket connection to NebulaOS Telemetry Cluster...
                    Real-time metrics visualization will appear here shortly.
                </p>
            </div>

            <div className="glass p-6">
                <h3>Component Health Probes</h3>
                <div style={{ marginTop: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    {[
                        { name: 'Identity Engine (Keycloak)', status: 'Active', latency: '12ms' },
                        { name: 'Secret Vault (HashiCorp)', status: 'Active', latency: '8ms' },
                        { name: 'Audit Broker (NATS)', status: 'Connected', latency: '4ms' },
                        { name: 'Provider Factory', status: 'Healthy', latency: '<1ms' }
                    ].map(comp => (
                        <div key={comp.name} style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: '0.75rem', borderBottom: '1px solid var(--glass-border)' }}>
                            <span>{comp.name}</span>
                            <div style={{ display: 'flex', gap: '1.5rem' }}>
                                <span style={{ color: 'var(--text-muted)' }}>{comp.latency}</span>
                                <span className="badge badge-success">{comp.status}</span>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default ObservabilityView;
