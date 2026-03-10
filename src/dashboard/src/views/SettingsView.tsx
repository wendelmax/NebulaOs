import React, { useState, useEffect } from 'react';
import { User, Bell, Shield, Cloud, CreditCard, ChevronRight, Trash2, Globe, RefreshCw, Key, Settings } from 'lucide-react';
import { api } from '../api/client';

const SettingsView: React.FC = () => {
    const [activeSection, setActiveSection] = useState('providers');
    const [providers, setProviders] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [showAddProvider, setShowAddProvider] = useState(false);
    
    // New provider state
    const [newProvider, setNewProvider] = useState({
        name: '',
        type: 'proxmox',
        endpoint: '',
        credentials: ''
    });

    const sections = [
        { id: 'profile', icon: User, label: 'Neural Identity', desc: 'Secure account details and profile verification.' },
        { id: 'notifications', icon: Bell, label: 'Pulse Alerts', desc: 'Configure platform telemetry and system signals.' },
        { id: 'security', icon: Shield, label: 'Shield & Core', desc: 'Manage cryptographic keys and access tokens.' },
        { id: 'providers', icon: Cloud, label: 'Cloud Sectors', desc: 'Connect and bridge distributed cloud clusters.' },
        { id: 'billing', icon: CreditCard, label: 'Fiscal Core', desc: 'Manage capital provision and consumption nodes.' }
    ];

    useEffect(() => {
        if (activeSection === 'providers') {
            fetchProviders();
        }
    }, [activeSection]);

    const fetchProviders = async () => {
        setLoading(true);
        try {
            const resp = await api.getProviders();
            setProviders(resp.data || []);
        } catch (err) {
            console.error("Failed to fetch providers", err);
        } finally {
            setLoading(false);
        }
    };

    const handleAddProvider = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await api.registerProvider(newProvider);
            setNewProvider({ name: '', type: 'proxmox', endpoint: '', credentials: '' });
            setShowAddProvider(false);
            fetchProviders();
        } catch (err) {
            console.error("Failed to register provider", err);
        }
    };

    const handleDeleteProvider = async (id: string) => {
        if (!confirm("Are you sure you want to decommission this provider?")) return;
        try {
            await api.deleteProvider(id);
            fetchProviders();
        } catch (err) {
            console.error("Failed to delete provider", err);
        }
    };

    const renderProviders = () => (
        <div className="flex flex-col gap-10">
            <header className="flex justify-between items-end">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight text-main/90">Infrastructure Sectors</h2>
                    <p className="text-sm text-dim mt-2 font-medium">Manage primary and secondary cloud provider clusters.</p>
                </div>
                <button 
                    onClick={() => setShowAddProvider(!showAddProvider)}
                    className={showAddProvider ? 'btn-secondary' : 'btn-primary'}
                >
                    {showAddProvider ? 'Abandon Declaration' : 'Declare New Sector'}
                </button>
            </header>

            {showAddProvider && (
                <form onSubmit={handleAddProvider} className="glass p-10 bg-gradient-to-br from-primary/5 to-transparent border-primary/20 relative">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
                        <div className="space-y-2">
                            <label className="text-[10px] font-black text-dim uppercase tracking-widest px-1">Sector Label</label>
                            <input 
                                type="text" 
                                required
                                value={newProvider.name}
                                onChange={e => setNewProvider({...newProvider, name: e.target.value})}
                                placeholder="e.g. Proxmox-Alpha-Node"
                                className="premium-input" 
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-[10px] font-black text-dim uppercase tracking-widest px-1">Infrastructure Type</label>
                            <select 
                                value={newProvider.type}
                                onChange={e => setNewProvider({...newProvider, type: e.target.value})}
                                className="premium-select"
                            >
                                <option value="proxmox">Proxmox Virtual Environment</option>
                                <option value="openstack">OpenStack Cloud</option>
                                <option value="aws">AWS (Geospatial Compatible)</option>
                                <option value="baremetal">Bare Metal / IPMI</option>
                            </select>
                        </div>
                    </div>
                    <div className="space-y-2 mb-8">
                        <label className="text-[10px] font-black text-dim uppercase tracking-widest px-1">Neural API Endpoint</label>
                        <input 
                            type="url" 
                            required
                            value={newProvider.endpoint}
                            onChange={e => setNewProvider({...newProvider, endpoint: e.target.value})}
                            placeholder="https://nebula-node-1.sovereign.local/api"
                            className="premium-input" 
                        />
                    </div>
                    <div className="space-y-2 mb-10">
                        <label className="text-[10px] font-black text-dim uppercase tracking-widest px-1">Secure Handshake Credentials</label>
                        <div className="relative group">
                            <Key size={16} className="absolute left-4 top-1/2 -translate-y-1/2 text-dim group-focus-within:text-primary transition-colors" />
                            <input 
                                type="password" 
                                required
                                value={newProvider.credentials}
                                onChange={e => setNewProvider({...newProvider, credentials: e.target.value})}
                                placeholder="****************"
                                className="premium-input pl-12" 
                            />
                        </div>
                    </div>
                    <button type="submit" className="btn-primary w-full py-4 text-xs font-bold shadow-xl shadow-primary/20">Establish Secure Connection</button>
                </form>
            )}

            <div className="flex flex-col gap-4">
                {loading && <div className="text-center py-20"><RefreshCw className="animate-spin mx-auto text-primary opacity-50" size={32} /></div>}
                {!loading && providers.length === 0 && !showAddProvider && (
                    <div className="glass p-20 text-center border-dashed border-white/10 bg-white/2">
                        <Cloud size={64} className="mx-auto mb-6 text-dim opacity-20" />
                        <h3 className="text-lg font-bold text-dim/60 mb-2">Zero Sectors Connected</h3>
                        <p className="text-sm text-dim max-w-xs mx-auto">Connect your first private or public cloud cluster to begin global orchestration.</p>
                    </div>
                )}
                {providers.map(p => (
                    <div key={p.id} className="glass p-8 flex items-center gap-8 group hover:bg-white/5 border-white/5 hover:border-white/10 transition-all">
                        <div className="p-5 rounded-2xl bg-white/5 text-primary border border-white/5 group-hover:border-primary/30 transition-all shadow-inner">
                            <Globe size={28} />
                        </div>
                        <div className="flex-1">
                            <div className="flex items-center gap-3 mb-1.5">
                                <h4 className="text-lg font-black text-main/90 tracking-tight group-hover:text-main">{p.name}</h4>
                                <span className="text-[9px] font-black uppercase bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 px-3 py-1 rounded-full">{p.status}</span>
                            </div>
                            <div className="text-xs font-mono text-dim font-bold tracking-tight bg-black/20 w-fit px-2 py-0.5 rounded-md border border-white/5">{p.endpoint}</div>
                        </div>
                        <div className="flex items-center gap-10">
                            <div className="text-right">
                                <div className="text-[10px] font-black uppercase text-primary tracking-widest mb-1 italic op-80">{p.type} Engine</div>
                                <div className="text-[10px] font-bold text-dim uppercase tracking-widest op-50">Handshake established {new Date(p.created_at).toLocaleDateString()}</div>
                            </div>
                            <button 
                                onClick={() => handleDeleteProvider(p.id)}
                                className="p-3 text-red-500/40 hover:text-red-400 hover:bg-red-500/10 rounded-xl border border-transparent hover:border-red-500/20 transition-all"
                            >
                                <Trash2 size={20} />
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );

    return (
        <div className="flex flex-col gap-10">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Settings size={14} className="text-dim" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">Platform Parameters</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">Control Plane</h1>
                    <p className="text-muted mt-2 font-medium">Global configuration for your autonomous NebulaOS instance.</p>
                </div>
            </header>

            <div className="flex flex-col lg:flex-row gap-10">
                {/* Section Selector */}
                <div className="w-full lg:w-96 flex flex-col gap-3">
                    {sections.map(section => (
                        <div 
                            key={section.id} 
                            onClick={() => setActiveSection(section.id)}
                            className={`p-6 rounded-3xl cursor-pointer transition-all flex items-center gap-6 border relative group overflow-hidden ${
                                activeSection === section.id 
                                ? 'glass border-primary/50 shadow-2xl shadow-primary/20' 
                                : 'hover:bg-white/5 border-white/5'
                            }`}
                        >
                            {activeSection === section.id && (
                                <div className="absolute inset-0 bg-primary/5 -z-10" />
                            )}
                            <div className={`p-4 rounded-2xl transition-all ${activeSection === section.id ? 'bg-primary text-white shadow-lg shadow-primary/30' : 'bg-white/5 text-dim border border-white/5 group-hover:border-white/10 group-hover:text-main'}`}>
                                <section.icon size={24} />
                            </div>
                            <div className="flex-1">
                                <div className={`text-lg font-black tracking-tight ${activeSection === section.id ? 'text-main' : 'text-dim group-hover:text-main'}`}>{section.label}</div>
                                <div className="text-[10px] text-dim/60 font-medium uppercase tracking-widest leading-relaxed">{section.desc}</div>
                            </div>
                            <ChevronRight size={18} className={`transition-transform duration-300 ${activeSection === section.id ? 'text-primary translate-x-1' : 'text-dim/20'}`} />
                        </div>
                    ))}
                </div>

                {/* Section Content */}
                <div className="flex-1">
                    <div className="glass p-12 min-h-[600px] border-white/10 relative overflow-hidden">
                        <div className="absolute -top-24 -right-24 w-96 h-96 bg-primary/2 blur-[150px] rounded-full pointer-events-none" />
                        
                        {activeSection === 'providers' ? renderProviders() : (
                            <div className="flex flex-col items-center justify-center h-full text-center py-24 relative z-10">
                                <div className="w-24 h-24 rounded-full bg-white/5 flex items-center justify-center mb-8 border border-white/5">
                                    {React.createElement(sections.find(s => s.id === activeSection)?.icon || User, { size: 48, className: 'text-dim opacity-50' })}
                                </div>
                                <h2 className="text-3xl font-black text-main tracking-tight mb-4">{sections.find(s => s.id === activeSection)?.label}</h2>
                                <p className="text-dim font-medium max-w-sm leading-relaxed mb-10">
                                    The <span className="text-primary-light">{activeSection}</span> configuration array is currently being synchronized with the neural core. Detailed sub-parameters will stabilize shortly.
                                </p>
                                <button className="btn-secondary px-10 border-white/10 hover:border-white/20">Acknowledge Status</button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SettingsView;
