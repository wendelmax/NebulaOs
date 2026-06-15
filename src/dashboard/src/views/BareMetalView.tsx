import React, { useState, useEffect } from 'react';
import { Server, Activity, Terminal, Play, RefreshCw, Cpu, Database } from 'lucide-react';
import { api } from '../api/client';
import { useLocale } from '../contexts/LocaleContext';

const BareMetalView: React.FC = () => {
    const { t } = useLocale();
    const [nodes, setNodes] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedNode, setSelectedNode] = useState<any | null>(null);
    const [logs, setLogs] = useState<any[]>([]);

    const fetchNodes = async () => {
        setLoading(true);
        try {
            const resp = await api.getBareMetalNodes();
            setNodes(resp.data || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const fetchLogs = async (nodeId: string) => {
        try {
            const resp = await api.getBareMetalLogs(nodeId);
            setLogs(resp.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    const handleProvision = async (id: string) => {
        try {
            await api.provisionBareMetalNode(id);
            fetchNodes();
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchNodes();
    }, []);

    useEffect(() => {
        if (selectedNode) {
            fetchLogs(selectedNode.id);
            const interval = setInterval(() => fetchLogs(selectedNode.id), 5000);
            return () => clearInterval(interval);
        }
    }, [selectedNode]);

    return (
        <div className="flex flex-col gap-8 max-w-[1400px]">
            <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Server size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">{t.baremetal.tag}</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">{t.baremetal.title}</h1>
                    <p className="text-muted mt-2 font-medium">{t.baremetal.description}</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchNodes}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        {t.baremetal.scanNetwork}
                    </button>
                    <button className="btn-primary">{t.baremetal.addNode}</button>
                </div>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
                {/* Node List */}
                <div className="lg:col-span-2 flex flex-col gap-4">
                    <div className="flex items-center justify-between">
                        <h3 className="text-lg font-bold">{t.baremetal.physicalCapacity}</h3>
                        <span className="badge badge-secondary">{nodes.length} {t.baremetal.nodes}</span>
                    </div>
                    
                    <div className="flex flex-col gap-3">
                        {nodes.map(node => (
                            <div 
                                key={node.id} 
                                onClick={() => setSelectedNode(node)}
                                className={`p-5 rounded-2xl border cursor-pointer transition-all ${
                                    selectedNode?.id === node.id 
                                    ? 'bg-primary/10 border-primary/20 shadow-lg shadow-primary/5' 
                                    : 'bg-white/5 border-white/5 hover:border-white/10'
                                }`}
                            >
                                <div className="flex items-start justify-between">
                                    <div className="flex gap-4">
                                        <div className={`p-3 bg-white/5 rounded-xl ${node.state === 'active' ? 'text-primary' : 'text-dim'}`}>
                                            <Server size={24} />
                                        </div>
                                        <div>
                                            <div className="font-bold text-main text-lg">{node.name}</div>
                                            <div className="text-[10px] font-mono text-muted uppercase tracking-wider">{node.mac}</div>
                                        </div>
                                    </div>
                                    <span className={`badge ${
                                        node.state === 'active' ? 'badge-success' : 
                                        node.state === 'provisioning' ? 'badge-warning' : 'badge-secondary'
                                    }`}>
                                        {node.state}
                                    </span>
                                </div>
                                <div className="grid grid-cols-3 gap-2 mt-6 pt-4 border-t border-white/5">
                                    <div className="flex flex-col">
                                        <span className="text-[9px] uppercase tracking-tighter text-muted">{t.baremetal.cpu}</span>
                                        <span className="text-sm font-bold text-main">{node.cpu_cores} {t.baremetal.cores}</span>
                                    </div>
                                    <div className="flex flex-col">
                                        <span className="text-[9px] uppercase tracking-tighter text-muted">{t.baremetal.ram}</span>
                                        <span className="text-sm font-bold text-main">{node.memory_gb} GB</span>
                                    </div>
                                    <div className="flex flex-col">
                                        <span className="text-[9px] uppercase tracking-tighter text-muted">{t.baremetal.storage}</span>
                                        <span className="text-sm font-bold text-main">{node.disk_gb} GB</span>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Node Details & Logs */}
                <div className="lg:col-span-2 flex flex-col gap-6">
                    {selectedNode ? (
                        <>
                            <div className="glass p-6">
                                <div className="flex items-center justify-between mb-8">
                                    <div className="flex items-center gap-3">
                                        <div className="p-2 bg-primary/20 rounded-lg text-primary">
                                            <Activity size={18} />
                                        </div>
                                        <h3 className="text-xl font-bold">{t.baremetal.hardwareOps}</h3>
                                    </div>
                                    <button 
                                        className="btn-primary py-2 px-4 text-sm"
                                        onClick={() => handleProvision(selectedNode.id)}
                                        disabled={selectedNode.state === 'provisioning'}
                                    >
                                        <Play size={14} />
                                        {t.baremetal.provisionIPXE}
                                    </button>
                                </div>

                                <div className="grid grid-cols-2 gap-4 mb-8">
                                    <div className="p-4 rounded-xl bg-white/3 border border-white/5">
                                        <div className="text-xs text-muted mb-1 flex items-center gap-1">
                                            <Activity size={12} /> {t.baremetal.ipmiAddr}
                                        </div>
                                        <div className="font-mono text-sm">{selectedNode.ipmi_address || t.baremetal.notConfigured}</div>
                                    </div>
                                    <div className="p-4 rounded-xl bg-white/3 border border-white/5">
                                        <div className="text-xs text-muted mb-1 flex items-center gap-1">
                                            <Database size={12} /> {t.baremetal.deptId}
                                        </div>
                                        <div className="font-mono text-sm">{selectedNode.department_id || t.baremetal.global}</div>
                                    </div>
                                </div>

                                <div className="flex flex-col gap-3">
                                    <div className="flex items-center justify-between text-xs font-bold uppercase tracking-widest text-dim">
                                        <div className="flex items-center gap-2">                                        <Terminal size={14} /> {t.baremetal.provisioningLogs}</div>
                                        <span className="animate-pulse flex items-center gap-1">
                                            <div className="w-1 h-1 bg-primary rounded-full" />                                             {t.baremetal.live}
                                        </span>
                                    </div>
                                    <div className="console shadow-inner">
                                        {logs.length === 0 ? (
                                            <div className="text-muted italic opacity-50">{t.baremetal.waitingEvents}</div>
                                        ) : (
                                            logs.map((log, idx) => (
                                                <div key={idx} className="flex gap-4 mb-2">
                                                    <span className="text-muted opacity-30 pointer-events-none">{new Date(log.timestamp).toLocaleTimeString()}</span>
                                                    <span className={log.level === 'error' ? 'text-red-400' : 'text-primary'}>[{log.level.toUpperCase()}]</span>
                                                    <span className="text-dim">{log.message}</span>
                                                </div>
                                            ))
                                        )}
                                    </div>
                                </div>
                            </div>
                        </>
                    ) : (
                        <div className="glass p-12 flex flex-col items-center justify-center text-center opacity-40 min-h-[400px]">
                            <Cpu size={64} className="mb-4" />
                            <h3 className="text-xl font-bold">{t.baremetal.noNodeSelected}</h3>
                            <p className="text-muted mt-2 max-w-[300px]">{t.baremetal.noNodeDesc}</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default BareMetalView;
