import React, { useState, useEffect } from 'react';
import { Shield, Lock, Globe, RefreshCw, GitBranch, ArrowRight } from 'lucide-react';
import { api } from '../api/client';

const Networking: React.FC = () => {
    const [securityGroups, setSecurityGroups] = useState<any[]>([]);
    const [gslbEndpoints, setGslbEndpoints] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    
    // Routing state
    const [routeSource, setRouteSource] = useState('proxmox');
    const [routeTarget, setRouteTarget] = useState('openstack');
    const [routeDest, setRouteDest] = useState('10.20.0.0/24');
    const [creatingRoute, setCreatingRoute] = useState(false);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [sgResp, _] = await Promise.all([
                fetch('http://localhost:8000/security-groups?project_id=v-p1'),
                api.getNetworkStatus('')
            ]);

            if (sgResp.ok) {
                const sgData = await sgResp.json();
                setSecurityGroups(sgData || []);
            }
            setGslbEndpoints([
                { id: 'g-1', dns_record: 'api.nebula.local', state: 'active', policy: { strategy: 'latency' } }
            ]);
        } catch (err) {
            console.error("Failed to fetch networking data", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const handleCreateRoute = async () => {
        setCreatingRoute(true);
        try {
            await api.createInterProviderRoute({
                source_provider: routeSource,
                target_provider: routeTarget,
                route: {
                    destination: routeDest,
                    next_hop: '172.16.0.1' 
                }
            });
            alert('Cross-provider route established!');
        } catch (err) {
            console.error(err);
        } finally {
            setCreatingRoute(false);
        }
    };

    return (
        <div className="flex flex-col gap-10 max-w-[1400px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Globe size={14} className="text-secondary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">Connectivity Fabric</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">SDN & Global Mesh</h1>
                    <p className="text-muted mt-2 font-medium">Manage VPC firewalls, security groups and multi-provider routing.</p>
                </div>
                <button className="btn-secondary" onClick={fetchData}>
                    <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                    Refresh
                </button>
            </header>

            {/* Inter-Provider Routing Card */}
            <div className="glass p-10 bg-gradient-to-br from-primary/5 to-transparent border-primary/10 relative overflow-hidden">
                <div className="absolute top-0 right-0 p-10 opacity-5 pointer-events-none">
                    <GitBranch size={120} className="text-primary" />
                </div>
                
                <div className="flex items-center gap-4 mb-10">
                    <div className="p-3 bg-primary/20 rounded-xl text-primary ring-1 ring-primary/30">
                        <GitBranch size={24} />
                    </div>
                    <div>
                        <h2 className="text-2xl font-bold tracking-tight text-main/90">Inter-Provider Routing</h2>
                        <p className="text-xs font-bold text-dim uppercase tracking-widest mt-1">Establishing sovereign high-speed bridges</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 xl:grid-cols-5 gap-6 items-center bg-black/40 p-6 rounded-3xl border border-white/5">
                    <div className="xl:col-span-1 space-y-2">
                        <label className="text-[10px] font-black text-dim uppercase tracking-widest px-1">Source Sector</label>
                        <select className="premium-select" value={routeSource} onChange={e => setRouteSource(e.target.value)}>
                            <option value="proxmox">Proxmox Cluster</option>
                            <option value="openstack">OpenStack Cloud</option>
                            <option value="aws">AWS Public</option>
                        </select>
                    </div>
                    
                    <div className="flex justify-center xl:col-span-1">
                        <ArrowRight className="text-dim/30" />
                    </div>

                    <div className="xl:col-span-1 space-y-2">
                        <label className="text-[10px] font-black text-dim uppercase tracking-widest px-1">Target Sector</label>
                        <select className="premium-select" value={routeTarget} onChange={e => setRouteTarget(e.target.value)}>
                            <option value="openstack">OpenStack Cloud</option>
                            <option value="proxmox">Proxmox Cluster</option>
                            <option value="aws">AWS Public</option>
                        </select>
                    </div>

                    <div className="xl:col-span-1 space-y-2">
                        <label className="text-[10px] font-black text-dim uppercase tracking-widest px-1">Destination CIRD</label>
                        <input 
                            type="text" 
                            className="premium-input" 
                            value={routeDest} 
                            onChange={e => setRouteDest(e.target.value)} 
                        />
                    </div>

                    <button className="btn-primary w-full py-4 text-xs font-bold h-fit shadow-xl shadow-primary/10 disabled:opacity-50" onClick={handleCreateRoute} disabled={creatingRoute}>
                        {creatingRoute ? 'establishing tunnel...' : 'Construct Mesh Link'}
                    </button>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <div className="glass p-10 flex flex-col gap-8">
                    <div className="flex items-center gap-4 pb-4 border-b border-white/5">
                        <Lock size={20} className="text-primary" />
                        <h2 className="text-xl font-bold tracking-tight">Security Groups</h2>
                    </div>

                    <div className="flex flex-col gap-4">
                        {securityGroups.length === 0 && !loading && (
                            <div className="py-10 text-center text-dim font-medium italic opacity-40">Zero group definitions found.</div>
                        )}
                        {securityGroups.map((sg: any) => (
                            <div key={sg.id} className="glass p-6 flex items-center gap-6 group hover:bg-white/5">
                                <div className="p-4 rounded-2xl bg-white/5 text-primary border border-white/5 group-hover:border-primary/30 transition-all">
                                    <Shield size={24} />
                                </div>
                                <div className="flex-1">
                                    <div className="font-bold text-main/90 group-hover:text-main transition-colors leading-tight">{sg.name}</div>
                                    <div className="text-[10px] font-mono text-dim mt-1.5 uppercase font-bold tracking-wider">{sg.id}</div>
                                </div>
                                <div className="text-right">
                                    <div className="text-xs font-black text-main tracking-tighter">{sg.rules?.length || 0} RULES</div>
                                    <div className="text-[9px] font-black text-emerald-400/70 uppercase tracking-widest mt-1 italic">ENFORCED</div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="glass p-10 flex flex-col gap-8">
                    <div className="flex items-center gap-4 pb-4 border-b border-white/5">
                        <Globe size={20} className="text-secondary" />
                        <h2 className="text-xl font-bold tracking-tight">Global Endpoints (GSLB)</h2>
                    </div>

                    <div className="flex flex-col gap-4">
                        {gslbEndpoints.length === 0 && !loading && (
                            <div className="py-10 text-center text-dim font-medium italic opacity-40">No global endpoints configured.</div>
                        )}
                        {gslbEndpoints.map((ep: any) => (
                            <div key={ep.id} className="glass p-6 flex items-center gap-6 group hover:bg-white/5">
                                <div className="p-4 rounded-2xl bg-secondary-gradient text-white shadow-xl shadow-secondary/10">
                                    <Globe size={24} />
                                </div>
                                <div className="flex-1">
                                    <div className="font-bold text-main/90 group-hover:text-main transition-colors leading-tight">{ep.dns_record}</div>
                                    <div className="text-[10px] font-black text-dim mt-1.5 uppercase tracking-widest">Dynamic Latency Steering</div>
                                </div>
                                <span className="badge badge-success px-4 py-1.5 text-[9px] font-black uppercase tracking-widest">{ep.state}</span>
                            </div>
                        ))}
                    </div>
                    
                    <button className="btn-secondary w-full py-4 text-xs font-bold mt-2 hover:bg-white/5 border-dashed">
                        Declare New Global Endpoint
                    </button>
                </div>
            </div>
        </div>
    );
};

export default Networking;
