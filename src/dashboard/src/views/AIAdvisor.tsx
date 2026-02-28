import React, { useState, useEffect } from 'react';
import { Brain, TrendingDown, ShieldAlert, Sparkles, CheckCircle, ArrowRight, RefreshCw, Zap } from 'lucide-react';

const AIAdvisor: React.FC = () => {
    const [insights, setInsights] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchInsights = async () => {
        setLoading(true);
        try {
            const resp = await fetch('http://api.nebula.local/intelligence/advisor?project_id=v-p1');
            if (resp.ok) {
                const data = await resp.json();
                setInsights(data || []);
            }
        } catch (err) {
            console.error("Failed to fetch AI insights", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchInsights();
    }, []);

    const getInsightIcon = (type: string) => {
        switch (type.toLowerCase()) {
            case 'cost': return TrendingDown;
            case 'security': return ShieldAlert;
            case 'performance': return Zap;
            default: return Sparkles;
        }
    };

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
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <button className="btn-secondary" onClick={fetchInsights} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        Refresh
                    </button>
                    <button className="btn-secondary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <Sparkles size={18} />
                        Run Deep Scan
                    </button>
                </div>
            </header>

            <div className="stats-grid" style={{ marginTop: '2rem' }}>
                <div className="stat-card glass" style={{ gridColumn: 'span 2' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1.5rem' }}>
                        <div>
                            <h2 style={{ fontSize: '1.5rem' }}>Top Recommendations</h2>
                            <p className="text-muted" style={{ fontSize: '0.9rem' }}>Applied intelligence for your infrastructure.</p>
                        </div>
                        <span className="stat-label" style={{ background: 'var(--bg-accent)', color: 'var(--primary-light)' }}>
                            {insights.length} Insights Found
                        </span>
                    </div>

                    <div className="list-container">
                        {insights.length === 0 && !loading && (
                            <div className="text-center py-12">
                                <CheckCircle className="text-primary mx-auto mb-4" size={48} />
                                <p className="text-muted">Your infrastructure is fully optimized. No insights found.</p>
                            </div>
                        )}
                        {insights.map((insight: any, idx: number) => {
                            const Icon = getInsightIcon(insight.type);
                            return (
                                <div key={idx} className="list-item glass-hover" style={{ padding: '1.5rem', marginBottom: '1rem' }}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem', flex: 1 }}>
                                        <div className="stat-icon" style={{
                                            background: insight.severity === 'high' ? 'rgba(239, 68, 68, 0.1)' : 'rgba(56, 189, 248, 0.1)',
                                            color: insight.severity === 'high' ? '#ef4444' : '#0ea5e9'
                                        }}>
                                            <Icon size={20} />
                                        </div>
                                        <div>
                                            <div style={{ fontWeight: 700, fontSize: '1rem' }}>{insight.message}</div>
                                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.25rem' }}>
                                                {insight.actionable && <CheckCircle size={14} className="text-primary" />}
                                                <span style={{ fontSize: '0.8rem', color: 'var(--primary-light)', fontWeight: 600 }}>{insight.severity.toUpperCase()} IMPACT</span>
                                            </div>
                                        </div>
                                    </div>
                                    <button className="btn-primary" style={{ padding: '0.5rem 1rem', fontSize: '0.8rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                        Apply Fix
                                        <ArrowRight size={14} />
                                    </button>
                                </div>
                            );
                        })}
                    </div>
                </div>

                <div className="stat-card glass" style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <h3 style={{ fontWeight: 700 }}>AI Performance Score</h3>
                    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '150px' }}>
                        <div style={{ fontSize: '3rem', fontWeight: 800, color: 'var(--primary-light)' }}>{insights.length === 0 ? '100' : 94 - (insights.length * 2)}</div>
                        <div style={{ marginLeft: '0.5rem', color: 'var(--text-muted)' }}>/ 100</div>
                    </div>
                    <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', textAlign: 'center' }}>
                        {insights.length === 0 ? "Perfect score! No optimizations required." : "Your infrastructure is well-maintained with minor optimizations suggested."}
                    </p>
                </div>
            </div>
        </div>
    );
};

export default AIAdvisor;
