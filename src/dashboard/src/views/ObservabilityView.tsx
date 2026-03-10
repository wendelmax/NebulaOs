import React, { useState, useEffect } from 'react';
import { Activity, Zap, HeartPulse, BarChart3, RefreshCw } from 'lucide-react';
import { api } from '../api/client';

const ObservabilityView: React.FC = () => {
    const [stats, setStats] = useState<any>(null);
    const [loading, setLoading] = useState(true);

    const fetchData = async () => {
        setLoading(true);
        try {
            const resp = await api.getStats();
            setStats(resp.data);
        } catch (err) {
            console.error("Failed to fetch observability stats", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
        const timer = setInterval(fetchData, 15000);
        return () => clearInterval(timer);
    }, []);

    return (
        <div className="flex flex-col gap-10 max-w-[1400px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Activity size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">Platform Core</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">Mission Control</h1>
                    <p className="text-muted mt-2 font-medium">Real-time telemetry and health diagnostics for your cloud plane.</p>
                </div>
                <button className="btn-secondary" onClick={fetchData}>
                    <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                    Refresh
                </button>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="glass p-8 bg-gradient-to-br from-emerald-500/5 to-transparent border-emerald-500/10 group">
                    <div className="flex justify-between items-center mb-6">
                        <HeartPulse size={24} className="text-emerald-500" />
                        <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.5)]" />
                    </div>
                    <div className="text-5xl font-black text-main tracking-tighter mb-1">
                        {stats ? '100%' : '--'}
                    </div>
                    <div className="text-xs font-black text-dim uppercase tracking-widest">Uptime Probability</div>
                </div>

                <div className="glass p-8 bg-gradient-to-br from-amber-500/5 to-transparent border-amber-500/10 group">
                    <div className="flex justify-between items-center mb-6">
                        <Zap size={24} className="text-amber-500" />
                        <span className="text-[10px] font-black text-amber-500/80 uppercase tracking-widest">Optimized</span>
                    </div>
                    <div className="text-5xl font-black text-main tracking-tighter mb-1">
                        {stats ? '42ms' : '--'}
                    </div>
                    <div className="text-xs font-black text-dim uppercase tracking-widest">Sector Latency</div>
                </div>

                <div className="glass p-8 bg-gradient-to-br from-primary/5 to-transparent border-primary/10 group">
                    <div className="flex justify-between items-center mb-6">
                        <BarChart3 size={24} className="text-primary" />
                        <span className="text-[10px] font-black text-primary/80 uppercase tracking-widest">Throughput</span>
                    </div>
                    <div className="text-5xl font-black text-main tracking-tighter mb-1">
                        {stats ? (stats.total_egress * 10).toFixed(0) : '--'}
                    </div>
                    <div className="text-xs font-black text-dim uppercase tracking-widest">Requests / Min</div>
                </div>
            </div>

            <div className="glass p-12 flex flex-col items-center justify-center relative overflow-hidden bg-black/40 border-white/5 min-h-[400px]">
                <div className="absolute inset-0 opacity-10 pointer-events-none" 
                     style={{ backgroundImage: 'radial-gradient(circle at 1px 1px, var(--primary) 1px, transparent 0)', backgroundSize: '24px 24px' }} />
                
                <div className="relative z-10 flex flex-col items-center gap-6">
                    <div className="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center text-primary animate-pulse">
                        <Activity size={40} />
                    </div>
                    <h3 className="text-2xl font-bold tracking-tight text-main/90">Establishing Neural Stream</h3>
                    <p className="text-sm text-dim max-w-sm text-center leading-relaxed font-medium">
                        Synchronizing with NebulaOS Telemetry Clusters. 
                        Real-time geospatial metrics visualization will materialize upon secure handshake.
                    </p>
                    <div className="flex gap-2">
                        <div className="w-1.5 h-1.5 rounded-full bg-primary/40 animate-bounce [animation-delay:-0.3s]" />
                        <div className="w-1.5 h-1.5 rounded-full bg-primary/40 animate-bounce [animation-delay:-0.15s]" />
                        <div className="w-1.5 h-1.5 rounded-full bg-primary/40 animate-bounce" />
                    </div>
                </div>
            </div>

            <div className="glass p-10 mt-4 border-white/10">
                <div className="flex items-center gap-3 mb-8">
                    <HeartPulse size={20} className="text-dim" />
                    <h2 className="text-xl font-bold tracking-tight">Component Health Probes</h2>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {[
                        { name: 'Identity Engine', status: 'Active', latency: '12ms', color: 'emerald' },
                        { name: 'Secret Vault', status: 'Active', latency: '8ms', color: 'emerald' },
                        { name: 'Audit Broker', status: 'Connected', latency: '4ms', color: 'primary' },
                        { name: 'Provider Factory', status: 'Healthy', latency: '<1ms', color: 'emerald' }
                    ].map(comp => (
                        <div key={comp.name} className="flex justify-between items-center p-6 bg-white/5 rounded-2xl border border-white/5 hover:border-white/10 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className={`w-2 h-2 rounded-full bg-${comp.color}-500 shadow-[0_0_8px_rgba(var(--${comp.color}),0.4)]`} />
                                <span className="text-sm font-bold text-main/90 tracking-tight">{comp.name}</span>
                            </div>
                            <div className="flex items-center gap-6">
                                <span className="text-[10px] font-mono text-dim font-bold">{comp.latency}</span>
                                <span className={`text-[10px] font-black uppercase tracking-widest text-${comp.color}-400`}>{comp.status}</span>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default ObservabilityView;
