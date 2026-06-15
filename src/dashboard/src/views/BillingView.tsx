import React from 'react';
import { CreditCard, TrendingUp, RefreshCw, BarChart3, Globe } from 'lucide-react';
import { api } from '../api/client';
import { useLocale } from '../contexts/LocaleContext';

const BillingView: React.FC = () => {
    const { t } = useLocale();
    const [report, setReport] = React.useState<any>(null);
    const [loading, setLoading] = React.useState(true);

    const fetchReport = async () => {
        setLoading(true);
        try {
            const resp = await api.getBillingReport('v-t1');
            setReport(resp.data);
        } catch (err) {
            console.error("Failed to fetch billing report", err);
        } finally {
            setLoading(false);
        }
    };

    React.useEffect(() => {
        fetchReport();
    }, []);

    return (
        <div className="flex flex-col gap-10 max-w-[1400px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <BarChart3 size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">{t.billing.tag}</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">{t.billing.title}</h1>
                    <p className="text-muted mt-2 font-medium">{t.billing.description}</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchReport}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        {t.billing.refresh}
                    </button>
                    <button className="btn-primary">
                        <CreditCard size={18} className="mr-2" />
                        {t.billing.provisionCapital}
                    </button>
                </div>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="glass p-8 relative overflow-hidden group">
                    <div className="absolute top-0 right-0 p-6 opacity-10 group-hover:opacity-20 transition-opacity">
                        <TrendingUp size={64} className="text-primary" />
                    </div>
                    <div className="text-xs font-black text-dim uppercase tracking-widest mb-4">{t.billing.burnRate}</div>
                    <div className="text-5xl font-black text-main tracking-tighter mb-2">
                        ${report ? report.total_cost.toFixed(2) : '0.00'}
                    </div>
                    <div className="flex items-center gap-2 text-[11px] font-bold text-emerald-400 uppercase tracking-tight">
                        <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]" />
                        {t.billing.activeConsumption}
                    </div>
                </div>

                <div className="glass p-8 relative overflow-hidden group">
                    <div className="absolute top-0 right-0 p-6 opacity-10 group-hover:opacity-20 transition-opacity">
                        <Globe size={64} className="text-secondary" />
                    </div>
                    <div className="text-xs font-black text-dim uppercase tracking-widest mb-4">{t.billing.sovereignCompliance}</div>
                    <div className="text-3xl font-black text-main tracking-tight mb-2 uppercase italic text-glow-primary">
                        {t.billing.enforced}
                    </div>
                    <div className="text-[11px] font-bold text-dim uppercase tracking-widest leading-relaxed">
                        {t.billing.sectorBoundary}: <span className="text-secondary-light">nebula-sovereign-node-1</span>
                    </div>
                </div>

                <div className="glass p-8 flex flex-col justify-center bg-gradient-to-br from-primary/5 to-transparent">
                    <div className="text-xs font-black text-dim uppercase tracking-widest mb-4 italic op-50">{t.billing.nextStatement}</div>
                    <div className="text-xl font-bold text-main/80">March 31, 2026</div>
                    <div className="mt-4 flex items-center gap-2 text-[10px] font-black tracking-widest text-primary-light uppercase border-t border-white/5 pt-4">
                        {t.billing.autoReconciliation}
                    </div>
                </div>
            </div>

            <div className="mt-4">
                <div className="flex items-center gap-3 mb-8">
                    <BarChart3 size={20} className="text-dim" />
                    <h2 className="text-2xl font-bold tracking-tight">{t.billing.usageStatement}</h2>
                </div>

                <table>
                    <thead>
                        <tr>
                            <th>{t.billing.colResource}</th>
                            <th>{t.billing.colCategory}</th>
                            <th>{t.billing.colCost}</th>
                            <th>{t.billing.colGeoCompliance}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {report?.items?.map((item: any) => (
                            <tr key={item.resource_id}>
                                <td className="font-mono text-xs font-bold text-dim">{item.resource_id}</td>
                                <td className="font-bold text-main/90 uppercase tracking-widest text-[10px]">{item.type} Services</td>
                                <td className="font-mono font-bold text-main text-lg">${item.cost.toFixed(2)}</td>
                                <td>
                                    <span className="flex items-center gap-2 text-[10px] font-black tracking-widest text-emerald-500 uppercase">
                                        <div className="w-2 h-2 rounded-full border border-emerald-500/50 flex items-center justify-center p-0.5">
                                            <div className="w-full h-full rounded-full bg-emerald-500" />
                                        </div>
                                        {t.billing.verifiedCompliant}
                                    </span>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            <div className="glass p-10 mt-6 border-white/10 relative">
                <div className="absolute inset-0 bg-primary/5 blur-3xl rounded-full opacity-20 -z-10" />
                <h3 className="text-lg font-black uppercase tracking-widest text-dim mb-8">{t.billing.institutionalBreakdown}</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-x-12 gap-y-6">
                    {report?.items?.reduce((acc: any[], item: any) => {
                        const existing = acc.find(a => a.name === item.type);
                        if (existing) {
                            existing.cost += item.cost;
                        } else {
                            acc.push({ name: item.type, cost: item.cost });
                        }
                        return acc;
                    }, []).map((item: any) => (
                        <div key={item.name} className="flex justify-between items-center py-4 border-b border-white/5">
                            <span className="text-sm font-bold text-dim uppercase tracking-widest">{item.name} {t.billing.segment}</span>
                            <span className="font-black text-main text-xl tracking-tighter">${item.cost.toFixed(2)}</span>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default BillingView;
