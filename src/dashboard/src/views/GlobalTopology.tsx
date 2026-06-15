import React, { useState, useEffect } from 'react';
import { MapPin, RefreshCw, Layers, ChevronRight } from 'lucide-react';
import { api } from '../api/client';
import { useLocale } from '../contexts/LocaleContext';

const GlobalTopology: React.FC = () => {
    const { t } = useLocale();
    const [regions, setRegions] = useState<any[]>([]);
    const [zones, setZones] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedRegion, setSelectedRegion] = useState<string | null>(null);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [regResp, zoneResp] = await Promise.all([
                api.getRegions(),
                api.getZones()
            ]);
            setRegions(regResp.data || []);
            setZones(zoneResp.data || []);
        } catch (err) {
            console.error("Failed to fetch topology data", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const filteredZones = selectedRegion 
        ? zones.filter(z => z.region_id === selectedRegion)
        : [];

    return (
        <div className="flex flex-col gap-10 max-w-[1500px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <MapPin size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">{t.globalTopology.tag}</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">{t.globalTopology.title}</h1>
                    <p className="text-muted mt-2 font-medium">{t.globalTopology.description}</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchData}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        {t.globalTopology.syncNodes}
                    </button>
                    <div className="glass px-5 py-3 rounded-2xl flex items-center gap-3 border-emerald-500/10">
                        <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.4)]" />
                        <span className="text-xs font-bold text-main opacity-80 uppercase tracking-tighter">{t.globalTopology.meshOptimal}</span>
                    </div>
                </div>
            </header>

            <div className="grid grid-cols-1 xl:grid-cols-3 gap-8">
                {/* Map View */}
                <div className="xl:col-span-2 glass p-10 min-h-[600px] relative overflow-hidden group">
                    {/* Decorative background grid/dots */}
                    <div className="absolute inset-0 opacity-20 pointer-events-none" 
                         style={{ backgroundImage: 'radial-gradient(circle at 1px 1px, var(--text-dim) 1px, transparent 0)', backgroundSize: '40px 40px' }} />
                    <div className="absolute inset-0 bg-radial-at-c from-primary/10 via-transparent to-transparent opacity-50 pointer-events-none" />
                    
                    <div className="relative w-full h-full flex items-center justify-center">
                        {regions.map((reg, idx) => (
                            <div key={reg.id}
                                className={`glass p-6 w-72 cursor-pointer transition-all duration-500 hover:shadow-2xl ${selectedRegion === reg.id ? 'border-primary shadow-primary-glow/20 -translate-y-2' : 'hover:border-white/20'}`}
                                onClick={() => setSelectedRegion(reg.id)}
                                style={{
                                    position: 'absolute',
                                    top: `${25 + idx * 25}%`,
                                    left: `${10 + idx * 30}%`,
                                    zIndex: selectedRegion === reg.id ? 20 : 10
                                }}>
                                <div className="flex items-center gap-4 mb-3">
                                    <div className="w-10 h-10 rounded-xl bg-primary-gradient flex items-center justify-center text-white">
                                        <MapPin size={20} />
                                    </div>
                                    <div className="flex-1">
                                        <div className="text-xs font-black text-dim tracking-widest uppercase mb-0.5">{reg.id}</div>
                                        <h4 className="font-bold text-main">{reg.name}</h4>
                                    </div>
                                </div>
                                <div className="flex justify-between items-center text-[11px]">
                                    <span className="text-dim font-medium uppercase">{reg.location}</span>
                                    {reg.is_default && <span className="text-primary-light font-black tracking-tighter uppercase opacity-80">Default Home</span>}
                                </div>
                            </div>
                        ))}

                        <svg className="absolute inset-0 w-full h-full pointer-events-none opacity-10">
                            {regions.length > 1 && regions.map((_, i) => i > 0 && (
                                <line 
                                    key={i}
                                    x1={`${10 + (i-1) * 30 + 15}%`} y1={`${25 + (i-1) * 25 + 10}%`} 
                                    x2={`${10 + i * 30 + 15}%`} y2={`${25 + i * 25 + 10}%`} 
                                    stroke="var(--primary)" strokeWidth="1" strokeDasharray="8 8" 
                                />
                            ))}
                        </svg>
                    </div>
                </div>

                {/* Details Sidebar */}
                <div className="flex flex-col gap-6">
                    <div className="glass p-8 bg-gradient-to-br from-indigo-500/5 to-transparent flex-1 border-indigo-500/10">
                        <div className="flex items-center gap-3 mb-8 pb-4 border-b border-white/5">
                            <Layers size={20} className="text-indigo-400" />
                            <h2 className="text-xl font-bold tracking-tight">{t.globalTopology.regionalFabric}</h2>
                        </div>

                        {!selectedRegion ? (
                            <div className="flex flex-col items-center justify-center py-20 text-center gap-4 opacity-30">
                                <div className="w-16 h-16 rounded-full border-2 border-dashed border-white/20 flex items-center justify-center">
                                    <MapPin size={28} />
                                </div>
                                <p className="text-sm font-medium leading-relaxed max-w-[200px]">{t.globalTopology.selectRegion}</p>
                            </div>
                        ) : (
                            <div className="animate-fade-in flex flex-col gap-6">
                                <div className="p-4 rounded-2xl bg-primary/5 border border-primary/10">
                                    <div className="text-[10px] uppercase font-black text-primary-light tracking-widest mb-1">{t.globalTopology.activeSector}</div>
                                    <div className="text-lg font-bold text-main">{selectedRegion} {t.globalTopology.statusPanel}</div>
                                </div>

                                <div className="flex flex-col gap-3">
                                    <h3 className="text-xs font-bold text-dim uppercase tracking-widest mb-2 px-1">{t.globalTopology.availabilityClusters}</h3>
                                    {filteredZones.length === 0 && <p className="text-sm text-dim italic p-4">{t.globalTopology.noClusters}</p>}
                                    {filteredZones.map(zone => (
                                        <div key={zone.id} className="glass p-5 flex items-center gap-4 group cursor-pointer hover:bg-white/5 transition-colors">
                                            <div className="relative">
                                                <div className={`w-3 h-3 rounded-full ${zone.state === 'available' ? 'bg-emerald-500 shadow-[0_0_10px_rgba(16,185,129,0.3)]' : 'bg-amber-500 shadow-[0_0_10px_rgba(245,158,11,0.3)]'}`} />
                                            </div>
                                            <div className="flex-1">
                                                <div className="font-bold text-main/90 group-hover:text-main transition-colors">{zone.name}</div>
                                                <div className="text-[10px] text-dim font-bold uppercase tracking-wider mt-1">Cluster-{zone.id.split('-').pop()} // {zone.state}</div>
                                            </div>
                                            <ChevronRight size={14} className="text-dim opacity-0 group-hover:opacity-100 transition-all -translate-x-2 group-hover:translate-x-0" />
                                        </div>
                                    ))}
                                </div>
                                
                                <button className="btn-primary w-full py-4 text-sm font-bold mt-4 shadow-xl shadow-primary/10">
                                    {t.globalTopology.analyze}
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default GlobalTopology;
