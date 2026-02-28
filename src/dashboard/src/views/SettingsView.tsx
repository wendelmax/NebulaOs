import React from 'react';
import { User, Bell, Shield, Cloud, CreditCard, ChevronRight } from 'lucide-react';

const SettingsView: React.FC = () => {
    const sections = [
        { id: 'profile', icon: User, label: 'Profile Settings', desc: 'Update your account details and profile picture.' },
        { id: 'notifications', icon: Bell, label: 'Notifications', desc: 'Configure platform alerts and email preferences.' },
        { id: 'security', icon: Shield, label: 'Security & Auth', desc: 'Manage 2FA, API keys and access tokens.' },
        { id: 'providers', icon: Cloud, label: 'Cloud Providers', desc: 'Connect and configure AWS, OpenStack and Proxmox credentials.' },
        { id: 'billing', icon: CreditCard, label: 'Billing & Plans', desc: 'Manage subscriptions, payment methods and invoices.' }
    ];

    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <h1>Platform Settings</h1>
                    <p className="text-muted">Global configuration for your NebulaOS instance.</p>
                </div>
            </header>

            <div className="stats-grid" style={{ marginTop: '2rem' }}>
                <div className="stat-card glass" style={{ gridColumn: 'span 2' }}>
                    <div className="list-container">
                        {sections.map(section => (
                            <div key={section.id} className="list-item glass-hover" style={{ padding: '1.5rem', cursor: 'pointer' }}>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem', flex: 1 }}>
                                    <div className="stat-icon" style={{ background: 'var(--bg-accent)', color: 'var(--primary-light)' }}>
                                        <section.icon size={20} />
                                    </div>
                                    <div>
                                        <div style={{ fontWeight: 700, fontSize: '1.1rem' }}>{section.label}</div>
                                        <div style={{ fontSize: '0.9rem', color: 'var(--text-muted)', marginTop: '0.25rem' }}>{section.desc}</div>
                                    </div>
                                </div>
                                <ChevronRight size={20} className="text-muted" />
                            </div>
                        ))}
                    </div>
                </div>

                <div className="stat-card glass" style={{ background: 'var(--primary-gradient)', color: 'white' }}>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem', height: '100%', justifyContent: 'center', textAlign: 'center' }}>
                        <div style={{ background: 'rgba(255,255,255,0.2)', padding: '1rem', borderRadius: '50%', width: 'fit-content', margin: '0 auto' }}>
                            <Shield size={32} />
                        </div>
                        <h3 style={{ fontSize: '1.25rem', fontWeight: 800 }}>NebulaOS Shield</h3>
                        <p style={{ fontSize: '0.9rem', opacity: 0.9 }}>Your Enterprise instance is protected by active security group monitoring.</p>
                        <button className="btn-primary" style={{ background: 'white', color: 'var(--primary-main)', marginTop: '1rem' }}>
                            View Security Audit
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SettingsView;
