import React from 'react';
import { CreditCard, TrendingUp, Download, RefreshCw, BarChart3, Globe } from 'lucide-react';

const BillingView: React.FC = () => {
    const [report, setReport] = React.useState<any>(null);
    const [loading, setLoading] = React.useState(true);

    const fetchReport = async () => {
        setLoading(true);
        try {
            const resp = await fetch('http://api.nebula.local/billing/report?tenant_id=v-t1');
            if (resp.ok) {
                const data = await resp.json();
                setReport(data);
            }
        } catch (err) {
            console.error("Failed to fetch billing report", err);
        } finally {
            setLoading(false);
        }
    };

    React.useEffect(() => {
        fetchReport();
    }, []);

    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
            <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <h1 style={{ fontSize: '1.875rem' }}>Billing & Sovereign Control</h1>
                    <p style={{ color: 'var(--text-muted)', marginTop: '0.25rem' }}>Fiscal accountability and geographical compliance status.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <button className="btn-secondary" onClick={fetchReport} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        Refresh
                    </button>
                    <button className="button-primary">
                        <CreditCard size={20} />
                        <span>Add Credits</span>
                    </button>
                </div>
            </header>

            <div className="resource-grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))' }}>
                <div className="glass p-6">
                    <div style={{ display: 'flex', gap: '1rem', alignItems: 'center', marginBottom: '1.5rem' }}>
                        <div style={{ padding: '0.75rem', background: 'rgba(99, 102, 241, 0.1)', borderRadius: '12px' }}>
                            <TrendingUp color="var(--primary)" />
                        </div>
                        <div>
                            <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Monthly Estimated Cost</p>
                            <h2 style={{ margin: 0 }}>${report ? report.total_cost.toFixed(2) : '0.00'}</h2>
                        </div>
                    </div>
                    <div style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>
                        <span style={{ color: '#4ade80', fontWeight: 600 }}>Active</span> consumption billing
                    </div>
                </div>

                <div className="glass p-6">
                    <div style={{ display: 'flex', gap: '1rem', alignItems: 'center', marginBottom: '1.5rem' }}>
                        <div style={{ padding: '0.75rem', background: 'rgba(236, 72, 153, 0.1)', borderRadius: '12px' }}>
                            <Globe color="var(--secondary)" />
                        </div>
                        <div>
                            <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Sovereignty Status</p>
                            <h2 style={{ margin: 0 }}>ENFORCED</h2>
                        </div>
                    </div>
                    <div style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>
                        Regional Boundary: <span style={{ color: 'var(--text-main)', fontWeight: 600 }}>nebula-local</span>
                    </div>
                </div>
            </div>

            <div className="glass" style={{ padding: '1.5rem' }}>
                <div style={{ display: 'flex', gap: '1rem', alignItems: 'center', marginBottom: '2rem' }}>
                    <BarChart3 />
                    <h3>Usage Statement</h3>
                </div>

                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                    <thead>
                        <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--glass-border)' }}>
                            <th style={{ padding: '1rem 0' }}>Resource ID</th>
                            <th>Category</th>
                            <th>Cost</th>
                            <th>Sovereignty</th>
                        </tr>
                    </thead>
                    <tbody>
                        {report?.items?.map((item: any) => (
                            <tr key={item.resource_id} style={{ borderBottom: '1px solid var(--glass-border)' }}>
                                <td style={{ padding: '1rem 0' }}>{item.resource_id}</td>
                                <td>{item.type}</td>
                                <td>${item.cost.toFixed(2)}</td>
                                <td><span className="badge badge-success">Compliant</span></td>
                            </tr>
                        ))}
                        {(!report || !report.items || report.items.length === 0) && !loading && (
                            <tr>
                                <td colSpan={4} style={{ padding: '4rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                                    No billing records found for this period.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>

            <div className="glass p-8" style={{ background: 'linear-gradient(135deg, rgba(99, 102, 241, 0.05) 0%, transparent 100%)' }}>
                <h3>Institutional Cost breakdown</h3>
                <div style={{ marginTop: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    {report?.items?.reduce((acc: any[], item: any) => {
                        const existing = acc.find(a => a.name === item.type);
                        if (existing) {
                            existing.cost += item.cost;
                        } else {
                            acc.push({ name: item.type, cost: item.cost });
                        }
                        return acc;
                    }, []).map((item: any) => (
                        <div key={item.name} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <span style={{ textTransform: 'capitalize' }}>{item.name} Services</span>
                            <div style={{ display: 'flex', gap: '2rem', alignItems: 'center' }}>
                                <span style={{ fontWeight: 600 }}>${item.cost.toFixed(2)}</span>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default BillingView;
