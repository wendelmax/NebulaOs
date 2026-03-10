import React, { useState, useEffect } from 'react';
import { ShieldCheck, FileText, RefreshCw, Landmark } from 'lucide-react';
import { api } from '../api/client';

const GovernanceView: React.FC = () => {
    const [policy, setPolicy] = useState<any>(null);
    const [loading, setLoading] = useState(true);

    const fetchData = async () => {
        setLoading(true);
        try {
            const resp = await api.getPolicy('t1');
            setPolicy(resp.data);
        } catch (err) {
            console.error("Failed to fetch governance data", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    return (
        <div className="flex flex-col gap-10 max-w-[1400px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Landmark size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">Sovereign Authority</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">Compliance Fabric</h1>
                    <p className="text-muted mt-2 font-medium">Enforcing policy, quotas, and sovereign audit standards across global clusters.</p>
                </div>
                <button className="btn-secondary" onClick={fetchData}>
                    <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                    Audit Sync
                </button>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <div className="glass p-8 flex flex-col gap-8 group">
                    <div className="flex justify-between items-center">
                        <h3 className="text-xs font-black uppercase tracking-widest text- dim italic">Resource Quota // 01</h3>
                        <span className="text-lg font-black text-main leading-none">CPU</span>
                    </div>
                    <div className="flex flex-col gap-3">
                        <div className="flex justify-between items-end">
                            <span className="text-xs font-bold text-dim uppercase tracking-tighter">14 / 20 vCPUs</span>
                            <span className="text-2xl font-black text-main tracking-tighter">70%</span>
                        </div>
                        <div className="h-1 w-full bg-white/5 rounded-full overflow-hidden">
                            <div className="h-full bg-primary-gradient rounded-full shadow-[0_0_8px_rgba(var(--primary-rgb),0.3)] transition-all duration-1000" style={{ width: '70%' }} />
                        </div>
                    </div>
                </div>

                <div className="glass p-8 flex flex-col gap-8 group">
                    <div className="flex justify-between items-center">
                        <h3 className="text-xs font-black uppercase tracking-widest text- dim italic">Resource Quota // 02</h3>
                        <span className="text-lg font-black text-main leading-none">RAM</span>
                    </div>
                    <div className="flex flex-col gap-3">
                        <div className="flex justify-between items-end">
                            <span className="text-xs font-bold text-dim uppercase tracking-tighter">32 / 64 GB</span>
                            <span className="text-2xl font-black text-main tracking-tighter">50%</span>
                        </div>
                        <div className="h-1 w-full bg-white/5 rounded-full overflow-hidden">
                            <div className="h-full bg-secondary-gradient rounded-full shadow-[0_0_8px_rgba(var(--secondary-rgb),0.3)] transition-all duration-1000" style={{ width: '50%' }} />
                        </div>
                    </div>
                </div>

                <div className="glass p-8 flex items-center gap-6 group hover:border-emerald-500/20">
                    <div className={`p-4 rounded-2xl ${policy?.enforce_sovereignty ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 shadow-[0_0_20px_rgba(16,185,129,0.1)]' : 'bg-primary/10 text-primary border border-primary/20'}`}>
                        <ShieldCheck size={32} />
                    </div>
                    <div>
                        <div className="text-[10px] font-black text-dim uppercase tracking-widest mb-1">Audit Posture</div>
                        <h4 className="text-lg font-black text-main tracking-tight leading-none uppercase italic">
                            {policy?.enforce_sovereignty ? 'Fully Compliant' : 'Monitoring Active'}
                        </h4>
                    </div>
                </div>
            </div>

            <div className="glass p-10 mt-4 border-white/10">
                <div className="flex items-center gap-3 mb-10 pb-4 border-b border-white/5">
                    <FileText size={20} className="text-dim" />
                    <h2 className="text-2xl font-bold tracking-tight">Geospatial Audit Stream</h2>
                </div>
                <div className="flex flex-col gap-2">
                    {[1, 2, 3, 4, 5].map(i => (
                        <div key={i} className="flex grid grid-cols-12 gap-6 items-center p-5 rounded-2xl hover:bg-white/5 transition-colors group cursor-pointer border border-transparent hover:border-white/5">
                            <div className="col-span-2 text-[10px] font-mono font-bold text-dim uppercase tracking-widest">2026-03-10 14:0{i}</div>
                            <div className="col-span-8 flex items-center gap-3">
                                <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_4px_rgba(16,185,129,0.5)]" />
                                <span className="text-sm font-medium text-main/80 group-hover:text-main transition-colors">API_REQUEST: SOVEREIGNTY_CHECK_PASSED (TenantID: <span className="text-primary-light font-bold">t1</span>)</span>
                            </div>
                            <div className="col-span-2 text-right">
                                <span className="text-[10px] font-black tracking-widest text-primary-light uppercase bg-primary/10 px-3 py-1 rounded-md border border-primary/20">SUCCESS</span>
                            </div>
                        </div>
                    ))}
                </div>
                
                <button className="w-full py-4 mt-8 text-[10px] font-black uppercase tracking-[0.2em] text-dim hover:text-main transition-colors border border-dashed border-white/10 rounded-2xl hover:bg-white/5">
                    Load Archive logs
                </button>
            </div>
        </div>
    );
};

export default GovernanceView;
