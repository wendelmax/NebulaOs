import React from 'react';
import { api } from '../api/client';
import ResourceCard from '../components/ResourceCard';
import { Cpu, HardDrive, Network, Users, LayoutDashboard, Brain, Activity, ChevronRight } from 'lucide-react';
import { useLocale } from '../contexts/LocaleContext';

interface GlobalStats {
    total_cpus: number;
    total_storage: number;
    total_egress: number;
    active_tenants: number;
    trend_cpus: number;
    trend_storage: number;
}

interface OverviewProps {
}

const Overview: React.FC<OverviewProps> = () => {
    const { t } = useLocale();
    const [stats, setStats] = React.useState<GlobalStats>({
        total_cpus: 0,
        total_storage: 0,
        total_egress: 0,
        active_tenants: 0,
        trend_cpus: 0,
        trend_storage: 0
    });

    React.useEffect(() => {
        const fetchStats = async () => {
            try {
                const resp = await api.getStats();
                setStats(resp.data);
            } catch (err) {
                console.error("Failed to fetch stats", err);
            }
        };
        fetchStats();
        const interval = setInterval(fetchStats, 10000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="flex flex-col gap-12 max-w-[1600px] mx-auto">
            <header className="flex items-end gap-8 mb-4">
                <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                        <span className="badge badge-success animate-pulse">{t.overview.tag}</span>
                        <span className="text-dim text-xs font-bold uppercase tracking-widest">{t.overview.version}</span>
                    </div>
                    <h1 className="text-6xl font-black leading-tight tracking-tighter">
                        {t.overview.title}<span className="bg-clip-text text-transparent bg-gradient-to-r from-indigo-500 to-fuchsia-500">{t.overview.titleAccent}</span>
                    </h1>
                    <p className="text-muted text-xl mt-4 font-medium max-w-2xl leading-relaxed">
                        {t.overview.description}
                        <span className="text-main"> Proxmox</span>, <span className="text-main">OpenStack</span>, and <span className="text-main">Bare Metal</span>.
                    </p>
                </div>
            </header>

            <div className="resource-grid">
                <ResourceCard
                    title={t.overview.computeFabric}
                    value={stats.total_cpus.toFixed(1)}
                    unit={t.overview.vCpus}
                    icon={Cpu}
                    trend={stats.trend_cpus}
                    color="var(--primary)"
                />
                <ResourceCard
                    title={t.overview.coldStorage}
                    value={stats.total_storage.toFixed(1)}
                    unit={t.overview.tb}
                    icon={HardDrive}
                    trend={stats.trend_storage}
                    color="var(--secondary)"
                />
                <ResourceCard
                    title={t.overview.dataEgress}
                    value={stats.total_egress.toFixed(0)}
                    unit={t.overview.gb}
                    icon={Network}
                    trend={28}
                    color="#818cf8"
                />
                <ResourceCard
                    title={t.overview.sovereignTenants}
                    value={stats.active_tenants.toString()}
                    unit={t.overview.units}
                    icon={Users}
                    color="#fbbf24"
                />
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 glass p-10 relative overflow-hidden group">
                    <div className="absolute top-0 right-0 p-8 opacity-10 group-hover:opacity-20 transition-opacity">
                        <LayoutDashboard size={120} />
                    </div>
                    <h3 className="text-2xl mb-8 flex items-center gap-3">
                        <Activity size={24} className="text-primary" />
                        {t.overview.livePerformance}
                    </h3>
                    <div className="flex flex-col items-center justify-center py-20 border-2 border-dashed border-white/5 rounded-3xl bg-white/autoborder-white/5">
                        <div className="relative mb-6">
                            <div className="absolute inset-0 blur-2xl bg-primary/20 rounded-full animate-pulse"></div>
                            <LayoutDashboard className="relative text-dim" size={56} />
                        </div>
                        <h4 className="text-muted font-bold tracking-tight">{t.overview.telemetryInitializing}</h4>
                        <p className="text-dim text-sm mt-2">{t.overview.telemetryDesc}</p>
                    </div>
                </div>

                <div className="flex flex-col gap-6">
                    <div className="glass p-8 bg-gradient-to-br from-indigo-500/5 to-transparent border-indigo-500/10">
                        <h3 className="text-xl mb-6 flex items-center gap-3">
                            <Brain size={20} className="text-indigo-400" />
                            {t.overview.nebulaIQ}
                        </h3>
                        <div className="p-4 rounded-xl bg-indigo-500/10 border border-indigo-500/20 mb-4">
                            <p className="text-sm font-medium text-indigo-200">{t.overview.recommendationTitle}</p>
                            <p className="text-xs text-indigo-300/70 mt-1">{t.overview.recommendationDesc}</p>
                        </div>
                        <button className="btn-primary w-full py-4 text-sm mt-2">
                            {t.overview.viewAI}
                        </button>
                    </div>

                    <div className="glass p-8">
                        <h3 className="text-xl mb-6 flex items-center gap-3">
                            <Activity size={20} className="text-muted" />
                            {t.overview.systemAudit}
                        </h3>
                        <div className="flex flex-col gap-4">
                            {[1, 2, 3].map((i) => (
                                <div key={i} className="flex gap-4 items-center p-4 rounded-2xl bg-white/5 hover:bg-white/10 transition-colors cursor-pointer group">
                                    <div className="w-2 h-2 rounded-full bg-primary shadow-[0_0_10px_var(--primary-glow)]" />
                                    <div className="flex-1">
                                        <p className="text-sm font-bold opacity-80 group-hover:opacity-100 transition-opacity">RESOURCE_PROVISIONED</p>
                                        <p className="text-[10px] text-muted uppercase font-bold tracking-wider mt-1">2m ago // AlphaNode</p>
                                    </div>
                                    <ChevronRight size={14} className="text-dim opacity-0 group-hover:opacity-100 transition-all -translate-x-2 group-hover:translate-x-0" />
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Overview;
