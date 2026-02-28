import React from 'react';
import { Shield, Plus, Lock } from 'lucide-react';

const Networking: React.FC = () => {
    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <h1>Network Security</h1>
                    <p className="text-muted">Manage VPC firewalls and security groups for your projects.</p>
                </div>
                <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                    <Plus size={18} />
                    New Security Group
                </button>
            </header>

            <div className="glass p-8" style={{ marginTop: '2rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '2rem' }}>
                    <Lock className="text-primary" />
                    <h2 style={{ fontSize: '1.5rem' }}>Security Groups</h2>
                </div>

                <div className="list-container">
                    <div className="list-item glass-hover" style={{ padding: '1.5rem' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                            <div className="stat-icon" style={{ width: '40px', height: '40px' }}><Shield size={20} /></div>
                            <div>
                                <h4 style={{ fontWeight: 700 }}>default-vpc-sg</h4>
                                <span className="text-muted" style={{ fontSize: '0.8rem' }}>Allows HTTP/HTTPS and SSH ingress</span>
                            </div>
                        </div>
                        <div style={{ display: 'flex', gap: '2rem' }}>
                            <div className="text-center">
                                <span className="stat-label">Rules</span>
                                <div style={{ fontWeight: 700 }}>5 Active</div>
                            </div>
                            <div className="text-center">
                                <span className="stat-label">Applied to</span>
                                <div style={{ fontWeight: 700 }}>12 Resources</div>
                            </div>
                        </div>
                        <button className="btn-secondary" style={{ padding: '0.5rem 1rem' }}>Edit Rules</button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Networking;
