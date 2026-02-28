import { LayoutDashboard, Server, Globe, Settings, LogOut, ChevronRight, HardDrive, Activity, Landmark, Receipt, Sun, Moon, ShoppingBag, Map, Brain } from 'lucide-react';

export type TabType = 'overview' | 'resources' | 'storage' | 'governance' | 'observability' | 'networking' | 'billing' | 'settings' | 'marketplace' | 'global' | 'advisor';

interface SidebarItemProps {
    icon: any;
    label: string;
    active?: boolean;
    onClick?: () => void;
}

const SidebarItem: React.FC<SidebarItemProps> = ({ icon: Icon, label, active = false, onClick }) => (
    <div className={`sidebar-item ${active ? 'active' : ''}`} onClick={onClick} style={{ cursor: 'pointer' }}>
        <Icon size={20} />
        <span>{label}</span>
        {active && <ChevronRight size={16} style={{ marginLeft: 'auto' }} />}
    </div>
);

interface DashboardShellProps {
    children: React.ReactNode;
    activeTab: TabType;
    onTabChange: (tab: TabType) => void;
    theme: 'dark' | 'light';
    onToggleTheme: () => void;
}

const DashboardShell: React.FC<DashboardShellProps> = ({ children, activeTab, onTabChange, theme, onToggleTheme }) => {
    const logoSrc = theme === 'dark' ? '/logo-dark.png' : '/logo-light.png';
    return (
        <div className="dashboard-container">
            {/* Sidebar */}
            <aside className="sidebar glass">
                <header style={{ padding: '2rem', display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', justifyContent: 'space-between' }}>
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <img src={logoSrc} alt="NebulaOS" style={{ width: '100px', height: '100px', objectFit: 'contain' }} />
                        </div>
                        <button
                            className="glass-hover"
                            onClick={onToggleTheme}
                            style={{
                                padding: '0.75rem',
                                borderRadius: '12px',
                                border: '1px solid var(--glass-border)',
                                background: 'var(--bg-accent)',
                                color: 'var(--text-main)',
                                cursor: 'pointer'
                            }}
                        >
                            {theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
                        </button>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column' }}>
                        <span style={{ fontWeight: 800, fontSize: '1.5rem', letterSpacing: '-0.025em', background: 'var(--primary-gradient)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                            NebulaOS
                        </span>
                        <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                            Enterprise Cloud
                        </span>
                    </div>
                </header>

                <nav className="sidebar-nav">
                    <SidebarItem icon={LayoutDashboard} label="Overview" active={activeTab === 'overview'} onClick={() => onTabChange('overview')} />
                    <SidebarItem icon={Server} label="Resources" active={activeTab === 'resources'} onClick={() => onTabChange('resources')} />
                    <SidebarItem icon={HardDrive} label="Storage" active={activeTab === 'storage'} onClick={() => onTabChange('storage')} />
                    <SidebarItem icon={Landmark} label="Governance" active={activeTab === 'governance'} onClick={() => onTabChange('governance')} />
                    <SidebarItem icon={Activity} label="Observability" active={activeTab === 'observability'} onClick={() => onTabChange('observability')} />
                    <SidebarItem icon={Receipt} label="Billing & usage" active={activeTab === 'billing'} onClick={() => onTabChange('billing')} />
                    <SidebarItem icon={Globe} label="Networking" active={activeTab === 'networking'} onClick={() => onTabChange('networking')} />
                    <SidebarItem icon={Map} label="Global Map" active={activeTab === 'global'} onClick={() => onTabChange('global')} />
                    <SidebarItem icon={ShoppingBag} label="Marketplace" active={activeTab === 'marketplace'} onClick={() => onTabChange('marketplace')} />
                    <SidebarItem icon={Brain} label="AI Advisor" active={activeTab === 'advisor'} onClick={() => onTabChange('advisor')} />
                    <SidebarItem icon={Settings} label="Settings" active={activeTab === 'settings'} onClick={() => onTabChange('settings')} />
                </nav>

                <div style={{ marginTop: 'auto', paddingTop: '1.5rem', borderTop: '1px solid var(--glass-border)' }}>
                    <SidebarItem icon={LogOut} label="Log Out" />
                </div>
            </aside>

            {/* Main Content */}
            <main className="main-content">
                {children}
            </main>
        </div>
    );
};

export default DashboardShell;
