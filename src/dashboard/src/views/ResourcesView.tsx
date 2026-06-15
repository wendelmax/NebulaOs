import React, { useState, useEffect } from 'react';
import { Server, Trash2, ExternalLink, RefreshCw, Filter, Search } from 'lucide-react';
import ResourceWizard from '../components/ResourceWizard';
import { api } from '../api/client';
import { useLocale } from '../contexts/LocaleContext';

const ResourcesView: React.FC = () => {
    const { t } = useLocale();
    const [resources, setResources] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [showWizard, setShowWizard] = useState(false);
    
    const fetchResources = async () => {
        setLoading(true);
        try {
            const resp = await api.getResources('v-p1');
            setResources(resp.data || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchResources();
    }, []);

    return (
        <div className="flex flex-col gap-8 max-w-[1400px]">
            {showWizard && (
                <ResourceWizard 
                    isOpen={showWizard} 
                    onClose={() => setShowWizard(false)} 
                    onSuccess={fetchResources}
                />
            )}
            
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Server size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">{t.resources.tag}</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">{t.resources.title}</h1>
                    <p className="text-muted mt-2 font-medium">{t.resources.description}</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchResources}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        {t.resources.refresh}
                    </button>
                    <button className="btn-primary" onClick={() => setShowWizard(true)}>
                        {t.resources.deploy}
                    </button>
                </div>
            </header>

            <div className="flex gap-4 items-center">
                <div className="search-box">
                    <Search size={18} className="text-muted" />
                    <input
                        type="text"
                        placeholder={t.resources.searchPlaceholder}
                    />
                </div>
                <button className="btn-secondary">
                    <Filter size={18} />
                    {t.resources.filters}
                </button>
            </div>

            <div className="mt-4">
                <table>
                    <thead>
                        <tr>
                            <th>{t.resources.colInstance}</th>
                            <th>{t.resources.colIdentity}</th>
                            <th>{t.resources.colProvider}</th>
                            <th>{t.resources.colState}</th>
                            <th style={{ textAlign: 'right' }}>{t.resources.colManagement}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {resources.length === 0 && !loading && (
                            <tr>
                                <td colSpan={5} className="py-24 text-center text-muted font-medium bg-transparent border-none">
                                    <div className="flex flex-col items-center gap-4 opacity-40">
                                        <Server size={48} />
                                        <p>{t.resources.noResources}</p>
                                    </div>
                                </td>
                            </tr>
                        )}
                        {resources.map((res) => (
                            <tr key={res.id}>
                                <td>
                                    <div className="flex items-center gap-4">
                                        <div className="p-3 bg-white/5 rounded-xl text-primary">
                                            <Server size={20} />
                                        </div>
                                        <div>
                                            <div className="font-bold text-main">{res.type}</div>
                                            <div className="text-xs text-dim font-mono mt-0.5">{res.id}</div>
                                        </div>
                                    </div>
                                </td>
                                <td className="font-semibold text-main/90">{res.name}</td>
                                <td>
                                    <span className="text-xs font-bold uppercase tracking-wider px-2 py-1 bg-white/5 rounded-md border border-white/5 text-dim">
                                        {res.provider}
                                    </span>
                                </td>
                                <td>
                                    <div className="flex items-center gap-2">
                                        <span className={`badge ${res.state.toLowerCase() === 'active' ? 'badge-success' : 'badge-warning'}`}>
                                            {res.state}
                                        </span>
                                    </div>
                                </td>
                                <td style={{ textAlign: 'right' }}>
                                    <div className="flex gap-2 justify-end">
                                        <button className="btn-secondary p-2.5 rounded-xl hover:text-primary">
                                            <ExternalLink size={18} />
                                        </button>
                                        <button className="btn-secondary p-2.5 rounded-xl hover:text-red-400">
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default ResourcesView;
