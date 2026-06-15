import React, { useState, useEffect } from 'react';
import { HardDrive, Database, Plus, Search, RefreshCw } from 'lucide-react';
import { api } from '../api/client';
import { useLocale } from '../contexts/LocaleContext';

const StorageView: React.FC = () => {
    const { t } = useLocale();
    const [volumes, setVolumes] = useState<any[]>([]);
    const [buckets, setBuckets] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [volResp, buckResp] = await Promise.all([
                api.getVolumes('v-p1'),
                api.getBuckets('v-p1')
            ]);

            setVolumes(volResp.data || []);
            setBuckets(buckResp.data || []);
        } catch (err) {
            console.error("Failed to fetch storage data", err);
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
                        <HardDrive size={14} className="text-secondary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">{t.storage.tag}</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">{t.storage.title}</h1>
                    <p className="text-muted mt-2 font-medium">{t.storage.description}</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchData}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        {t.storage.refresh}
                    </button>
                    <button className="btn-primary">
                        <Plus size={20} className="mr-2" />
                        {t.storage.allocate}
                    </button>
                </div>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {volumes.slice(0, 4).map(vol => (
                    <div key={vol.id} className="glass p-8 flex flex-col gap-6 group hover:border-primary-glow/20">
                        <div className="flex justify-between items-start">
                            <div className="p-3 bg-primary/10 rounded-xl text-primary">
                                <HardDrive size={24} />
                            </div>
                            <span className="badge badge-success">{vol.state}</span>
                        </div>
                        <div>
                            <h4 className="text-lg font-bold text-main mb-1 tracking-tight">{vol.name}</h4>
                            <p className="text-xs text-dim font-bold uppercase tracking-widest">{vol.size_gb} GB // Block</p>
                        </div>
                        <div className="mt-auto pt-4 border-t border-white/5 flex items-center justify-between">
                            <span className="text-[10px] text-dim font-mono">{vol.id}</span>
                            <Plus size={14} className="text-dim opacity-0 group-hover:opacity-100 cursor-pointer" />
                        </div>
                    </div>
                ))}
                {buckets.slice(0, Math.max(0, 4 - volumes.length)).map(buck => (
                    <div key={buck.id} className="glass p-8 flex flex-col gap-6 group">
                        <div className="flex justify-between items-start">
                            <div className="p-3 bg-secondary/10 rounded-xl text-secondary">
                                <Database size={24} />
                            </div>
                            <span className="badge badge-success">{buck.state}</span>
                        </div>
                        <div>
                            <h4 className="text-lg font-bold text-main mb-1 tracking-tight">{buck.name}</h4>
                            <p className="text-xs text-dim font-bold uppercase tracking-widest">{buck.region} // Object</p>
                        </div>
                        <div className="mt-auto pt-4 border-t border-white/5 flex items-center justify-between">
                            <span className="text-[10px] text-dim font-mono">Bucket-S3</span>
                            <Plus size={14} className="text-dim opacity-0 group-hover:opacity-100 cursor-pointer" />
                        </div>
                    </div>
                ))}
                {volumes.length === 0 && buckets.length === 0 && !loading && (
                    <div className="lg:col-span-4 glass p-12 text-center flex flex-col items-center gap-4 opacity-30">
                        <HardDrive size={48} />
                        <p className="font-medium text-lg italic">{t.storage.noVolumes}</p>
                    </div>
                )}
            </div>

            <div className="mt-4">
                <div className="flex justify-between items-center mb-6">
                    <h3 className="text-xl font-bold tracking-tight">{t.storage.sectorInventory}</h3>
                    <div className="search-box">
                        <Search size={18} className="text-muted" />
                        <input type="text" placeholder={t.storage.searchPlaceholder} />
                    </div>
                </div>

                <table>
                    <thead>
                        <tr>
                            <th>{t.storage.colResource}</th>
                            <th>{t.storage.colType}</th>
                            <th>{t.storage.colHealth}</th>
                            <th>{t.storage.colCapacity}</th>
                            <th>{t.storage.colRegion}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {volumes.map(vol => (
                            <tr key={vol.id}>
                                <td className="font-bold text-main/90">{vol.name}</td>
                                <td className="text-dim text-xs font-bold uppercase tracking-widest">{t.storage.blockStore}</td>
                                <td><span className="badge badge-success">{vol.state}</span></td>
                                <td className="font-mono text-sm">{vol.size_gb} GB</td>
                                <td className="text-[10px] font-black tracking-tighter text-dim uppercase">{vol.provider_id || 'LOCAL-NODE'}</td>
                            </tr>
                        ))}
                        {buckets.map(buck => (
                            <tr key={buck.id}>
                                <td className="font-bold text-main/90">{buck.name}</td>
                                <td className="text-dim text-xs font-bold uppercase tracking-widest">{t.storage.objectS3}</td>
                                <td><span className="badge badge-success">{buck.state}</span></td>
                                <td className="font-mono text-sm opacity-60">DYNAMIC</td>
                                <td className="text-[10px] font-black tracking-tighter text-dim uppercase">{buck.region}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default StorageView;
