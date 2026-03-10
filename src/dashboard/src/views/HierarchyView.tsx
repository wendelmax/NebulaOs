import React, { useState, useEffect } from 'react';
import { Landmark, Users, Plus, RefreshCw, ChevronRight } from 'lucide-react';
import { api } from '../api/client';

const HierarchyView: React.FC = () => {
    const [organizations, setOrganizations] = useState<any[]>([]);
    const [departments, setDepartments] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedOrg, setSelectedOrg] = useState<string | null>(null);

    const fetchData = async () => {
        setLoading(true);
        try {
            const orgResp = await api.getOrganizations();
            setOrganizations(orgResp.data || []);
            
            if (orgResp.data && orgResp.data.length > 0) {
                const firstOrg = orgResp.data[0].id;
                setSelectedOrg(firstOrg);
                const deptResp = await api.getDepartments(firstOrg);
                setDepartments(deptResp.data || []);
            }
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleOrgSelect = async (orgId: string) => {
        setSelectedOrg(orgId);
        try {
            const resp = await api.getDepartments(orgId);
            setDepartments(resp.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    return (
        <div className="flex flex-col gap-8 max-w-[1400px]">
             <header className="flex justify-between items-end pb-4 border-b border-white/5">
                <div>
                    <div className="flex items-center gap-2 mb-2">
                        <Landmark size={14} className="text-primary" />
                        <span className="text-dim text-[10px] font-bold uppercase tracking-widest">Enterprise Architecture</span>
                    </div>
                    <h1 className="text-4xl font-extrabold tracking-tight">Organization Hierarchy</h1>
                    <p className="text-muted mt-2 font-medium">Manage multi-tenant organizations and departments.</p>
                </div>
                <div className="flex gap-4">
                    <button className="btn-secondary" onClick={fetchData}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        Sync Registry
                    </button>
                    <button className="btn-primary">
                        <Plus size={18} />
                        New Organization
                    </button>
                </div>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Organizations List */}
                <div className="lg:col-span-1 flex flex-col gap-4">
                    <div className="flex items-center justify-between">
                        <h3 className="text-lg font-bold">Organizations</h3>
                        <span className="badge badge-secondary">{organizations.length}</span>
                    </div>
                    <div className="flex flex-col gap-2">
                        {organizations.map(org => (
                            <div 
                                key={org.id} 
                                onClick={() => handleOrgSelect(org.id)}
                                className={`p-4 rounded-2xl border cursor-pointer transition-all ${
                                    selectedOrg === org.id 
                                    ? 'bg-primary/10 border-primary/20 shadow-lg shadow-primary/5' 
                                    : 'bg-white/5 border-white/5 hover:border-white/10'
                                }`}
                            >
                                <div className="flex items-center justify-between">
                                    <div className="font-bold text-main">{org.name}</div>
                                    <ChevronRight size={16} className={selectedOrg === org.id ? 'text-primary' : 'text-muted'} />
                                </div>
                                <div className="text-[10px] font-mono text-muted mt-1 uppercase tracking-wider">{org.id}</div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Departments List */}
                <div className="lg:col-span-2 flex flex-col gap-4">
                    <div className="flex items-center justify-between">
                        <h3 className="text-lg font-bold">Departments</h3>
                        <button className="btn-secondary text-xs py-1 px-3">Add Department</button>
                    </div>
                    <div className="glass p-6 min-h-[400px]">
                        {departments.length === 0 ? (
                            <div className="flex flex-col items-center justify-center h-full text-muted opacity-40 gap-4">
                                <Users size={48} />
                                <p>No departments found for the selected organization.</p>
                            </div>
                        ) : (
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                {departments.map(dept => (
                                    <div key={dept.id} className="p-6 rounded-2xl bg-white/5 border border-white/5 hover:border-white/10 transition-all">
                                        <div className="flex items-start justify-between mb-4">
                                            <div className="p-3 bg-primary/10 rounded-xl text-primary">
                                                <Users size={20} />
                                            </div>
                                            <span className="text-[10px] font-mono text-muted uppercase tracking-wider">{dept.id}</span>
                                        </div>
                                        <h4 className="text-xl font-bold mb-1">{dept.name}</h4>
                                        <p className="text-sm text-dim leading-relaxed">System-assigned department for infrastructure partitioning.</p>
                                        <div className="mt-6 pt-4 border-t border-white/5 flex justify-between items-center text-xs">
                                            <span className="text-muted">Projects: 0</span>
                                            <button className="text-primary font-bold hover:underline">View Details</button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default HierarchyView;
