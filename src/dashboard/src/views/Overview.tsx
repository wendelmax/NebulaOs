import React from 'react';
import ResourceCard from '../components/ResourceCard';
import { Cpu, HardDrive, Network, Users, LayoutDashboard } from 'lucide-react';

interface OverviewProps {
    theme?: 'dark' | 'light';
}

const Overview: React.FC<OverviewProps> = ({ theme = 'dark' }) => {
    const [stats, setStats] = React.useState<any>({
        total_cpus: 42.8,
        total_storage: 1.2,
        total_egress: 892,
        active_tenants: 14,
        trend_cpus: 12,
        trend_storage: -4
    });

    React.useEffect(() => {
        const fetchStats = async () => {
            try {
                const resp = await fetch(`http://api.nebula.local/intelligence/stats?t=${Date.now()}`);
                if (resp.ok) {
                    const data = await resp.json();
                    setStats(data);
                }
            } catch (err) {
                console.error("Failed to fetch stats", err);
            }
        };
        fetchStats();
        const interval = setInterval(fetchStats, 10000);
        return () => clearInterval(interval);
    }, []);

    const logoSrc = theme === 'dark' ? '/logo-dark.png' : '/logo-light.png';
    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
            <header style={{ display: 'flex', alignItems: 'center', gap: '2.5rem', marginBottom: '2rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                    <img src={logoSrc} alt="NebulaOS" style={{ width: '160px', height: '160px', objectFit: 'contain' }} />
                </div>
                <div>
                    <h1 style={{ fontSize: '3.5rem', margin: 0, fontWeight: 900, background: 'var(--primary-gradient)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>Overview</h1>
                    <p style={{ color: 'var(--text-muted)', marginTop: '0.75rem', fontSize: '1.5rem', fontWeight: 500 }}>Enterprise Control Plane | Sovereign Status: Active</p>
                </div>
            </header>

            <div className="resource-grid">
                <ResourceCard
                    title="Compute Consumption"
                    value={stats.total_cpus.toFixed(1)}
                    unit="vCPUs"
                    icon={Cpu}
                    trend={stats.trend_cpus}
                    color="var(--primary)"
                />
                <ResourceCard
                    title="Storage Tier 1"
                    value={stats.total_storage.toFixed(1)}
                    unit="TB"
                    icon={HardDrive}
                    trend={stats.trend_storage}
                    color="var(--secondary)"
                />
                <ResourceCard
                    title="Egress Traffic"
                    value={stats.total_egress.toFixed(0)}
                    unit="GB/mo"
                    icon={Network}
                    trend={28}
                    color="#818cf8"
                />
                <ResourceCard
                    title="Active Tenants"
                    value={stats.active_tenants.toString()}
                    unit="units"
                    icon={Users}
                    color="#fbbf24"
                />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '1.5rem' }}>
                <div style={{ gridColumn: 'span 2' }} className="glass p-8 min-h-[400px]">
                    <div style={{ textAlign: 'center', padding: '4rem 0' }}>
                        <LayoutDashboard style={{ color: 'var(--text-muted)', marginBottom: '1rem' }} size={48} />
                        <h3 style={{ color: 'var(--text-muted)' }}>Real-time Performance Metrics</h3>
                        <p style={{ color: 'rgba(148, 163, 184, 0.6)', fontSize: '0.875rem' }}>Infrastructure synchronization in progress...</p>
                    </div>
                </div>

                <div className="glass" style={{ padding: '1.5rem' }}>
                    <h3 style={{ marginBottom: '1.5rem' }}>Recent Audit Events</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                        {[1, 2, 3, 4, 5].map((i) => (
                            <div key={i} style={{ display: 'flex', gap: '1rem', alignItems: 'center', padding: '0.75rem', borderRadius: '8px' }} className="glass-hover">
                                <div style={{ width: '8px', height: '8px', borderRadius: '50%', backgroundColor: 'var(--primary)' }} />
                                <div>
                                    <p style={{ fontSize: '0.875rem', fontWeight: 500 }}>Audit Event: RESOURCE_PROVISIONED</p>
                                    <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>2 mins ago | Project: ALPHA-NODE</p>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Overview;
