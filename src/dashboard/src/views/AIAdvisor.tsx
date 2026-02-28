import React from 'react';
import { Brain, TrendingDown, ShieldAlert, Sparkles, CheckCircle, ArrowRight } from 'lucide-react';

const AIAdvisor: React.FC = () => {
    const insights = [
        { id: 1, type: 'cost', msg: "Optimization: 4 'Zombie' volumes detected in dev-project-01.", impact: 'Save $124/mo', severity: 'medium', icon: TrendingDown },
        { id: 2, type: 'performance', msg: "latency peak in ap-southeast-1. Suggesting GSLB failover to eu-west-1.", impact: '35% Faster UX', severity: 'high', icon: Sparkles },
        { id: 3, type: 'security', msg: "Insecure SSH rule detected in 'default-vpc-sg'.", impact: 'Risk Mitigation', severity: 'high', icon: ShieldAlert }
    ];

    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.5rem' }}>
                        <Brain className="text-primary" size={28} />
                        <h1>AI Resource Advisor</h1>
                    </div>
                    <p className="text-muted">Proactive operational intelligence powered by NebulaAI.</p>
                </div>
                <button className="btn-secondary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                    <Sparkles size={18} />
                    Run Deep Scan
                </button>
            </header>

            <div className="stats-grid" style={{ marginTop: '2rem' }}>
                <div className="stat-card glass" style={{ gridColumn: 'span 2' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1.5rem' }}>
                        <div>
                            <h2 style={{ fontSize: '1.5rem' }}>Top Recommendations</h2>
                            <p className="text-muted" style={{ fontSize: '0.9rem' }}>Applied intelligence for your infrastructure.</p>
                        </div>
                        <span className="stat-label" style={{ background: 'var(--bg-accent)', color: 'var(--primary-light)' }}>
                            3 Insights Found
                        </span>
                    </div>

                    <div className="list-container">
                        {insights.map(insight => (
                            <div key={insight.id} className="list-item glass-hover" style={{ padding: '1.5rem' }}>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem', flex: 1 }}>
                                    <div className="stat-icon" style={{
                                        background: insight.severity === 'high' ? 'rgba(239, 68, 68, 0.1)' : 'rgba(56, 189, 248, 0.1)',
                                        color: insight.severity === 'high' ? '#ef4444' : '#0ea5e9'
                                    }}>
                                        <insight.icon size={20} />
                                    </div>
                                    <div>
                                        <div style={{ fontWeight: 700, fontSize: '1rem' }}>{insight.msg}</div>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.25rem' }}>
                                            <CheckCircle size={14} className="text-primary" />
                                            <span style={{ fontSize: '0.8rem', color: 'var(--primary-light)', fontWeight: 600 }}>{insight.impact}</span>
                                        </div>
                                    </div>
                                </div>
                                <button className="btn-primary" style={{ padding: '0.5rem 1rem', fontSize: '0.8rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                    Apply Fix
                                    <ArrowRight size={14} />
                                </button>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="stat-card glass" style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <h3 style={{ fontWeight: 700 }}>AI Performance Score</h3>
                    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '150px' }}>
                        <div style={{ fontSize: '3rem', fontWeight: 800, color: 'var(--primary-light)' }}>94</div>
                        <div style={{ marginLeft: '0.5rem', color: 'var(--text-muted)' }}>/ 100</div>
                    </div>
                    <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', textAlign: 'center' }}>
                        Your infrastructure is optimized. Minor cost issues detected.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default AIAdvisor;
