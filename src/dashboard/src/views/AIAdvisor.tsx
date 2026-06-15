import React, { useState, useEffect } from 'react';
import { Brain, TrendingDown, ShieldAlert, Sparkles, CheckCircle, ArrowRight, RefreshCw, Zap } from 'lucide-react';
import { api } from '../api/client';
import { useLocale } from '../contexts/LocaleContext';

const AIAdvisor: React.FC = () => {
    const { t } = useLocale();
    const [insights, setInsights] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchInsights = async () => {
        setLoading(true);
        try {
            const resp = await api.getAdvisorInsights('v-p1');
            setInsights(resp.data || []);
        } catch (err) {
            console.error("Failed to fetch AI insights", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchInsights();
    }, []);

    const getInsightIcon = (type: string) => {
        switch (type.toLowerCase()) {
            case 'cost': return TrendingDown;
            case 'security': return ShieldAlert;
            case 'performance': return Zap;
            default: return Sparkles;
        }
    };

    return (
        <div className="flex flex-col gap-10 max-w-[1400px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Brain size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">{t.advisor.tag}</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">{t.advisor.title}</h1>
                    <p className="text-muted mt-2 font-medium">{t.advisor.description}</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchInsights}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        {t.advisor.refresh}
                    </button>
                    <button className="btn-secondary">
                        <Sparkles size={18} className="text-primary-light" />
                        {t.advisor.neuralScan}
                    </button>
                </div>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 flex flex-col gap-6">
                    <div className="glass p-8 relative overflow-hidden">
                        <div className="flex justify-between items-center mb-8">
                            <div>
                                <h2 className="text-xl font-bold tracking-tight">{t.advisor.operationalInsights}</h2>
                                <p className="text-xs text-dim font-medium mt-1 uppercase tracking-wider">{t.advisor.insightsSubtitle}</p>
                            </div>
                            <div className="px-4 py-2 rounded-xl bg-white/5 border border-white/5 text-[10px] font-black tracking-widest text-primary-light uppercase">
                                    {insights.length} {t.advisor.recommendations}
                            </div>
                        </div>

                        <div className="flex flex-col gap-4">
                            {insights.length === 0 && !loading && (
                                <div className="flex flex-col items-center justify-center py-20 text-center gap-4 opacity-40">
                                    <div className="w-16 h-16 rounded-full bg-emerald-500/10 flex items-center justify-center text-emerald-400">
                                        <CheckCircle size={32} />
                                    </div>
                                    <p className="text-lg font-bold">{t.advisor.noInsights}</p>
                                    <p className="text-sm text-dim max-w-[280px]">{t.advisor.noInsightsDesc}</p>
                                </div>
                            )}
                            {insights.map((insight: any, idx: number) => {
                                const Icon = getInsightIcon(insight.type);
                                return (
                                    <div key={idx} className="glass p-6 flex items-center gap-6 group hover:bg-white/5 transition-all">
                                        <div className={`p-4 rounded-2xl bg-white/5 ${insight.severity === 'high' ? 'text-red-400 ring-1 ring-red-400/20' : 'text-primary ring-1 ring-primary/20'}`}>
                                            <Icon size={24} />
                                        </div>
                                        <div className="flex-1">
                                            <div className="font-bold text-main leading-tight group-hover:text-primary-light transition-colors">{insight.message}</div>
                                            <div className="flex items-center gap-3 mt-2">
                                                <span className={`text-[9px] font-black tracking-widest px-2 py-0.5 rounded border ${insight.severity === 'high' ? 'border-red-400/30 text-red-400/80' : 'border-primary/30 text-primary/80'} uppercase`}>
                                                    {insight.severity} {t.advisor.impact}
                                                </span>
                                                {insight.actionable && (
                                                    <span className="flex items-center gap-1 text-[9px] font-bold text-dim uppercase tracking-wider">
                                                        <CheckCircle size={10} className="text-emerald-500" />
                                                        {t.advisor.automatedFix}
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                        <button className="btn-primary py-2 px-6 text-xs font-bold whitespace-nowrap">
                                            {t.advisor.execute} <ArrowRight size={14} className="ml-2" />
                                        </button>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                </div>

                <div className="flex flex-col gap-8">
                    <div className="glass p-10 bg-gradient-to-br from-primary/10 to-transparent border-primary/10 flex flex-col items-center text-center">
                        <h3 className="text-xs font-black uppercase tracking-[0.2em] text-dim mb-10">{t.advisor.neuralScore}</h3>
                        <div className="relative">
                            <div className="absolute inset-0 blur-3xl bg-primary/20 rounded-full animate-pulse" />
                            <div className="relative flex flex-col items-center">
                                <span className="text-8xl font-black tracking-tighter bg-clip-text text-transparent bg-gradient-to-b from-white to-white/40 leading-none">
                                    {insights.length === 0 ? '100' : 94 - (insights.length * 2)}
                                </span>
                                <span className="text-xs font-bold text-dim uppercase tracking-widest mt-4">{t.advisor.sectorPerformance}</span>
                            </div>
                        </div>
                        <p className="text-sm text-muted leading-relaxed mt-10 max-w-[200px]">
                            {insights.length === 0 ? t.advisor.perfectScore : t.advisor.acceptableScore}
                        </p>
                    </div>

                    <div className="glass p-8 flex flex-col gap-4">
                        <h3 className="text-sm font-bold uppercase tracking-widest text- dim mb-2">{t.advisor.diagnostics}</h3>
                        <div className="p-4 rounded-xl bg-white/5 border border-white/5 flex items-center gap-4">
                            <Zap size={18} className="text-secondary" />
                            <div>
                                    <div className="text-[10px] font-black text-dim uppercase tracking-widest">{t.advisor.resourceDrift}</div>
                                    <div className="text-sm font-bold text-main mt-0.5">0.42% {t.advisor.variance}</div>
                            </div>
                        </div>
                        <div className="p-4 rounded-xl bg-white/5 border border-white/5 flex items-center gap-4">
                            <ShieldAlert size={18} className="text-red-400" />
                            <div>
                                    <div className="text-[10px] font-black text-dim uppercase tracking-widest">{t.advisor.compliance}</div>
                                    <div className="text-sm font-bold text-main mt-0.5">98.4% {t.advisor.aligned}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AIAdvisor;
