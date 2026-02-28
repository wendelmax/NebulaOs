import React, { useState, useEffect } from 'react';
import { Shield, Plus, Lock, Globe, RefreshCw } from 'lucide-react';

const Networking: React.FC = () => {
    const [securityGroups, setSecurityGroups] = useState<any[]>([]);
    const [gslbEndpoints, setGslbEndpoints] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchData = async () => {
        setLoading(true);
        try {
            const [sgResp, gslbResp] = await Promise.all([
                fetch('http://api.nebula.local/security-groups?project_id=v-p1'),
                fetch('http://api.nebula.local/network/gslb')
            ]);

            if (sgResp.ok) {
                const sgData = await sgResp.json();
                setSecurityGroups(sgData || []);
            }
            if (gslbResp.ok) {
                const gslbData = await gslbResp.json();
                setGslbEndpoints(gslbData || []);
            }
        } catch (err) {
            console.error("Failed to fetch networking data", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <h1>Network Security & Global Traffic</h1>
                    <p className="text-muted">Manage VPC firewalls, security groups and GSLB endpoints.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <button className="btn-secondary" onClick={fetchData} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        Refresh
                    </button>
                    <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <Plus size={18} />
                        New Security Group
                    </button>
                </div>
            </header>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem', marginTop: '2rem' }}>
                <div className="glass p-8">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '2rem' }}>
                        <Lock className="text-primary" />
                        <h2 style={{ fontSize: '1.5rem' }}>Security Groups</h2>
                    </div>

                    <div className="list-container">
                        {securityGroups.length === 0 && !loading && (
                            <p className="text-muted text-center py-8">No security groups found.</p>
                        )}
                        {securityGroups.map((sg: any) => (
                            <div key={sg.id} className="list-item glass-hover" style={{ padding: '1.5rem', marginBottom: '1rem' }}>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                                    <div className="stat-icon" style={{ width: '40px', height: '40px' }}><Shield size={20} /></div>
                                    <div>
                                        <h4 style={{ fontWeight: 700 }}>{sg.name}</h4>
                                        <span className="text-muted" style={{ fontSize: '0.8rem' }}>{sg.id}</span>
                                    </div>
                                </div>
                                <div style={{ display: 'flex', gap: '2rem' }}>
                                    <div className="text-center">
                                        <span className="stat-label">Rules</span>
                                        <div style={{ fontWeight: 700 }}>{sg.rules?.length || 0} Active</div>
                                    </div>
                                </div>
                                <button className="btn-secondary" style={{ padding: '0.5rem 1rem' }}>Edit Rules</button>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="glass p-8">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '2rem' }}>
                        <Globe className="text-secondary" />
                        <h2 style={{ fontSize: '1.5rem' }}>Global Endpoints (GSLB)</h2>
                    </div>

                    <div className="list-container">
                        {gslbEndpoints.length === 0 && !loading && (
                            <p className="text-muted text-center py-8">No GSLB endpoints configured.</p>
                        )}
                        {gslbEndpoints.map((ep: any) => (
                            <div key={ep.id} className="list-item glass-hover" style={{ padding: '1.5rem', marginBottom: '1rem' }}>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', flex: 1 }}>
                                    <div className="stat-icon" style={{ width: '40px', height: '40px', background: 'var(--secondary-gradient)' }}><Globe size={20} /></div>
                                    <div>
                                        <h4 style={{ fontWeight: 700 }}>{ep.dns_record}</h4>
                                        <span className="text-muted" style={{ fontSize: '0.8rem' }}>Strategy: {ep.policy?.strategy}</span>
                                    </div>
                                </div>
                                <div className="badge badge-success">{ep.state}</div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Networking;
