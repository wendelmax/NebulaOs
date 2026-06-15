import React, { useState, useEffect } from 'react';
import { Box, Database, ShieldCheck, Zap, RefreshCw, Rocket, ShoppingBag } from 'lucide-react';
import { api } from '../api/client';

const Marketplace: React.FC = () => {
    const [blueprints, setBlueprints] = useState<any[]>([]);
    const [presets, setPresets] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [deploying, setDeploying] = useState<string | null>(null);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [bpResp, presetResp] = await Promise.all([
                api.getBlueprints(),
                api.getPresets()
            ]);
            
            if (bpResp.status === 200) {
                setBlueprints(bpResp.data || []);
            }
            if (presetResp.status === 200) {
                setPresets(presetResp.data || []);
            }
        } catch (err) {
            console.error("Failed to fetch marketplace data", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const handleDeployBlueprint = async (bpId: string) => {
        setDeploying(bpId);
        try {
            const projectId = import.meta.env.VITE_DEFAULT_PROJECT_ID || 'v-p1';
            const resp = await api.deployBlueprint({ blueprint_id: bpId, project_id: projectId });
            if (resp.status === 200) alert('Blueprint deployment started!');
        } catch (err) {
            console.error(err);
        } finally {
            setDeploying(null);
        }
    };

    const handleProvisionPreset = async (presetId: string) => {
        setDeploying(presetId);
        try {
            const projectId = import.meta.env.VITE_DEFAULT_PROJECT_ID || 'v-p1';
            const resp = await api.provisionPreset({ preset_id: presetId, project_id: projectId });
            if (resp.status === 202) alert('Automated provisioning started!');
        } catch (err) {
            console.error(err);
        } finally {
            setDeploying(null);
        }
    };

    const getIcon = (category: string) => {
        switch (category?.toLowerCase()) {
            case 'infrastructure': return Box;
            case 'databases': return Database;
            case 'security': return ShieldCheck;
            default: return Zap;
        }
    };

    return (
        <div className="flex flex-col gap-12 max-w-[1400px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <ShoppingBag size={14} className="text-secondary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">Digital Marketplace</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">Cloud Catalog</h1>
                    <p className="text-muted mt-2 font-medium">Launch production-ready blueprints or one-click automated stacks.</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchData}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        Refresh
                    </button>
                </div>
            </header>

            {/* One-Click Presets Section */}
            <section className="animate-fade-in" style={{ animationDelay: '0.1s' }}>
                <div className="flex items-center gap-3 mb-8">
                    <div className="p-2 bg-primary/10 rounded-lg text-primary">
                        <Rocket size={20} />
                    </div>
                    <h2 className="text-2xl font-bold tracking-tight">Automated Solutions</h2>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {presets.map(preset => (
                        <div key={preset.id} className="glass p-8 flex flex-col gap-6 group hover:border-primary-glow/30">
                            <div className="flex justify-between items-start">
                                <div className="w-12 h-12 rounded-xl bg-primary-gradient flex items-center justify-center text-white shadow-lg shadow-primary-glow/20">
                                    <Rocket size={24} />
                                </div>
                                <span className="text-[10px] font-black tracking-widest text-primary-light uppercase bg-primary/10 px-3 py-1 rounded-full border border-primary/20">
                                    Stack
                                </span>
                            </div>
                            <div className="flex-1">
                                <h3 className="text-xl font-bold text-main mb-2 tracking-tight group-hover:text-primary-light transition-colors">{preset.name}</h3>
                                <p className="text-sm text-muted leading-relaxed line-clamp-3">{preset.description}</p>
                            </div>
                            <button className="btn-primary w-full py-4 text-sm font-bold" onClick={() => handleProvisionPreset(preset.id)} disabled={deploying === preset.id}>
                                {deploying === preset.id ? 'Initializing...' : 'Provision Cluster'}
                            </button>
                        </div>
                    ))}
                </div>
            </section>

            {/* Blueprints Section */}
            <section className="animate-fade-in" style={{ animationDelay: '0.2s' }}>
                <div className="flex items-center gap-3 mb-8">
                    <div className="p-2 bg-secondary/10 rounded-lg text-secondary">
                        <Box size={20} />
                    </div>
                    <h2 className="text-2xl font-bold tracking-tight">System Blueprints</h2>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {blueprints.map(bp => {
                        const Icon = getIcon(bp.category);
                        return (
                            <div key={bp.id} className="glass p-8 flex flex-col gap-6 group">
                                <div className="flex justify-between items-start">
                                    <div className="w-12 h-12 rounded-xl bg-white/5 flex items-center justify-center text-secondary border border-white/5 shadow-inner">
                                        <Icon size={24} />
                                    </div>
                                    <span className="badge badge-success">{bp.category}</span>
                                </div>
                                <div className="flex-1">
                                    <h3 className="text-xl font-bold text-main mb-2 tracking-tight">{bp.name}</h3>
                                    <p className="text-sm text-muted leading-relaxed line-clamp-3">{bp.description}</p>
                                </div>
                                <button className="btn-secondary w-full py-4 text-sm font-bold group-hover:bg-white/10" onClick={() => handleDeployBlueprint(bp.id)} disabled={deploying === bp.id}>
                                    {deploying === bp.id ? 'Deploying...' : 'Launch Blueprint'}
                                </button>
                            </div>
                        );
                    })}
                </div>
            </section>
        </div>
    );
};

export default Marketplace;
