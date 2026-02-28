import React from 'react';
import { CreditCard, GlobeLock, TrendingUp, History, Download } from 'lucide-react';

const BillingView: React.FC = () => {
    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
            <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                    <h1 style={{ fontSize: '1.875rem' }}>Billing & Sovereign Control</h1>
                    <p style={{ color: 'var(--text-muted)', marginTop: '0.25rem' }}>Fiscal accountability and geographical compliance status.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <button className="glass p-2 px-4 flex items-center gap-2 text-sm font-semibold">
                        <Download size={16} />
                        <span>Export CSV</span>
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
                            <h2 style={{ margin: 0 }}>$1,240.50</h2>
                        </div>
                    </div>
                    <div style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>
                        <span style={{ color: '#4ade80', fontWeight: 600 }}>+12.5%</span> vs last month
                    </div>
                </div>

                <div className="glass p-6">
                    <div style={{ display: 'flex', gap: '1rem', alignItems: 'center', marginBottom: '1.5rem' }}>
                        <div style={{ padding: '0.75rem', background: 'rgba(236, 72, 153, 0.1)', borderRadius: '12px' }}>
                            <GlobeLock color="var(--secondary)" />
                        </div>
                        <div>
                            <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Sovereignty Status</p>
                            <h2 style={{ margin: 0 }}>ENFORCED</h2>
                        </div>
                    </div>
                    <div style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>
                        Regional Boundary: <span style={{ color: 'var(--text-main)', fontWeight: 600 }}>BR-SOUTH-01</span>
                    </div>
                </div>
            </div>

            <div className="glass" style={{ padding: '1.5rem' }}>
                <div style={{ display: 'flex', gap: '1rem', alignItems: 'center', marginBottom: '2rem' }}>
                    <History />
                    <h3>Usage Statement</h3>
                </div>

                <table style={{ width: '100%' }}>
                    <thead>
                        <tr>
                            <th>Resource</th>
                            <th>Category</th>
                            <th>Region</th>
                            <th>Cost</th>
                            <th>Sovereignty</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>Compute-VM-Node-Alpha</td>
                            <td>Virtual Machine</td>
                            <td>br-south-1</td>
                            <td>$45.00</td>
                            <td><span className="badge badge-success">Compliant</span></td>
                        </tr>
                        <tr>
                            <td>Object-Store-Legacy</td>
                            <td>Storage (Bucket)</td>
                            <td>br-south-1</td>
                            <td>$12.40</td>
                            <td><span className="badge badge-success">Compliant</span></td>
                        </tr>
                        <tr>
                            <td>Network-Egress-Global</td>
                            <td>Networking</td>
                            <td>global</td>
                            <td>$185.20</td>
                            <td><span className="badge badge-warning">Exemption Active</span></td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <div className="glass p-8" style={{ background: 'linear-gradient(135deg, rgba(99, 102, 241, 0.05) 0%, transparent 100%)' }}>
                <h3>Institutional Cost breakdown</h3>
                <div style={{ marginTop: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    {[
                        { name: 'Sovereign Compute Tier', usage: '85%', cost: '$850.00' },
                        { name: 'Encrypted Storage Tier', usage: '12%', cost: '$148.00' },
                        { name: 'Audit & Compliance Service', usage: '3%', cost: '$242.50' }
                    ].map(item => (
                        <div key={item.name} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <span>{item.name}</span>
                            <div style={{ display: 'flex', gap: '2rem', alignItems: 'center' }}>
                                <span style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>{item.usage}</span>
                                <span style={{ fontWeight: 600 }}>{item.cost}</span>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default BillingView;
