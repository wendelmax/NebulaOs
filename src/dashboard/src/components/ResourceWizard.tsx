import React, { useState, useEffect } from 'react';
import { X, ChevronRight, ChevronLeft, Globe, Cpu, Layers, CheckCircle, Loader2 } from 'lucide-react';
import { api } from '../api/client';

interface ResourceWizardProps {
    isOpen: boolean;
    onClose: () => void;
    onSuccess: () => void;
}

const ResourceWizard: React.FC<ResourceWizardProps> = ({ isOpen, onClose, onSuccess }) => {
    const [step, setStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [regions, setRegions] = useState<any[]>([]);
    
    // Selection state
    const [config, setConfig] = useState({
        region: '',
        zone: '',
        type: 'COMPUTE',
        template: 'ubuntu-22.04',
        name: '',
        cpus: 2,
        ram: 4
    });

    useEffect(() => {
        if (isOpen) {
            fetchData();
        }
    }, [isOpen]);

    const fetchData = async () => {
        try {
            const resp = await api.getRegions();
            setRegions(resp.data || []);
            if (resp.data?.length > 0) {
                setConfig(prev => ({ ...prev, region: resp.data[0].id }));
            }
        } catch (err) {
            console.error("Failed to fetch regions", err);
        }
    };

    if (!isOpen) return null;

    const nextStep = () => setStep(step + 1);
    const prevStep = () => setStep(step - 1);

    const handleProvision = async () => {
        setLoading(true);
        try {
            await api.createResource({
                project_id: 'v-p1', // Default project
                name: config.name || `nebula-${Math.random().toString(36).substring(7)}`,
                type: config.type,
                provider: 'mock', // In a real scenario, this would be based on region/zone mapping
            });
            onSuccess();
            onClose();
        } catch (err) {
            console.error("Provisioning failed", err);
        } finally {
            setLoading(false);
        }
    };

    const renderStep = () => {
        switch (step) {
            case 1:
                return (
                    <div className="animate-fade-in">
                        <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
                            <Globe size={20} className="text-primary" />
                            Target Destination
                        </h3>
                        <p className="text-muted mb-6">Where should this resource be physically located?</p>
                        
                        <div className="grid grid-cols-1 gap-4">
                            {regions.map(r => (
                                <div 
                                    key={r.id} 
                                    onClick={() => setConfig({ ...config, region: r.id })}
                                    className={`p-4 rounded-xl border cursor-pointer transition-all ${config.region === r.id ? 'bg-primary-dark/20 border-primary shadow-lg shadow-primary/10' : 'bg-bg-accent/50 border-white/5 hover:border-white/20'}`}
                                >
                                    <div className="font-bold">{r.name}</div>
                                    <div className="text-sm text-muted">{r.location}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                );
            case 2:
                return (
                    <div className="animate-fade-in">
                        <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
                            <Layers size={20} className="text-primary" />
                            Software Template
                        </h3>
                        <p className="text-muted mb-6">Select the base image or orchestrated stack.</p>
                        
                        <div className="grid grid-cols-2 gap-4">
                            {['ubuntu-22.04', 'debian-12', 'k8s-worker', 'docker-host'].map(t => (
                                <div 
                                    key={t}
                                    onClick={() => setConfig({ ...config, template: t })}
                                    className={`p-4 rounded-xl border cursor-pointer transition-all ${config.template === t ? 'bg-primary-dark/20 border-primary shadow-lg' : 'bg-bg-accent/50 border-white/5 hover:border-white/20'}`}
                                >
                                    <div className="font-mono text-sm uppercase tracking-wider">{t}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                );
            case 3:
                return (
                    <div className="animate-fade-in">
                        <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
                            <Cpu size={20} className="text-primary" />
                            Specifications
                        </h3>
                        <p className="text-muted mb-6">Define the resource capacity and identification.</p>
                        
                        <div className="space-y-6">
                            <div>
                                <label className="block text-sm font-medium text-muted mb-2">Resource Name</label>
                                <input 
                                    type="text" 
                                    value={config.name}
                                    onChange={e => setConfig({ ...config, name: e.target.value })}
                                    placeholder="e.g. prod-web-01"
                                    className="w-full bg-bg-accent/50 border border-white/10 rounded-xl p-3 focus:border-primary outline-none transition-all"
                                />
                            </div>
                            
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-muted mb-2">vCPUs</label>
                                    <select 
                                        value={config.cpus}
                                        onChange={e => setConfig({ ...config, cpus: parseInt(e.target.value) })}
                                        className="w-full bg-bg-accent/50 border border-white/10 rounded-xl p-3 outline-none"
                                    >
                                        {[1, 2, 4, 8, 16].map(v => <option key={v} value={v}>{v} Blocks</option>)}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-muted mb-2">RAM (GB)</label>
                                    <select 
                                        value={config.ram}
                                        onChange={e => setConfig({ ...config, ram: parseInt(e.target.value) })}
                                        className="w-full bg-bg-accent/50 border border-white/10 rounded-xl p-3 outline-none"
                                    >
                                        {[2, 4, 8, 16, 32].map(v => <option key={v} value={v}>{v} GB</option>)}
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>
                );
            default:
                return null;
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
            <div className="glass w-full max-w-lg overflow-hidden flex flex-col shadow-2xl animate-scale-in">
                {/* Header */}
                <div className="p-6 border-b border-white/5 flex justify-between items-center">
                    <div>
                        <span className="text-xs font-bold text-primary uppercase tracking-widest">Step {step} of 3</span>
                        <h2 className="text-2xl font-black">Provision Infrastructure</h2>
                    </div>
                    <button onClick={onClose} className="p-2 hover:bg-white/5 rounded-full transition-colors">
                        <X size={24} />
                    </button>
                </div>

                {/* Progress Bar */}
                <div className="h-1 w-full bg-white/5">
                    <div 
                        className="h-full bg-primary-gradient transition-all duration-500 shadow-[0_0_10px_rgba(56,189,248,0.5)]" 
                        style={{ width: `${(step / 3) * 100}%` }} 
                    />
                </div>

                {/* Content */}
                <div className="p-8 flex-1 min-h-[400px]">
                    {renderStep()}
                </div>

                {/* Footer */}
                <div className="p-6 border-t border-white/5 flex justify-between bg-bg-accent/30">
                    <button 
                        onClick={prevStep} 
                        disabled={step === 1 || loading}
                        className="btn-secondary flex items-center gap-2 disabled:opacity-30"
                    >
                        <ChevronLeft size={18} />
                        Back
                    </button>
                    
                    {step < 3 ? (
                        <button 
                            onClick={nextStep}
                            className="btn-primary flex items-center gap-2"
                        >
                            Continue
                            <ChevronRight size={18} />
                        </button>
                    ) : (
                        <button 
                            onClick={handleProvision}
                            disabled={loading}
                            className="btn-primary flex items-center gap-2 bg-primary-gradient border-none"
                        >
                            {loading ? <Loader2 className="animate-spin" size={18} /> : <CheckCircle size={18} />}
                            {loading ? 'Orchestrating...' : 'Confirm & Launch'}
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ResourceWizard;
