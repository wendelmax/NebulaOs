import { LayoutDashboard, Server, Globe, Settings, LogOut, ChevronRight, HardDrive, Activity, Landmark, Receipt, Sun, Moon, ShoppingBag, Map, Brain } from 'lucide-react';

export type TabType = 'overview' | 'resources' | 'storage' | 'governance' | 'observability' | 'networking' | 'billing' | 'settings' | 'marketplace' | 'global' | 'advisor' | 'hierarchy' | 'baremetal';

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
    onLogout: () => void;
}

const DashboardShell: React.FC<DashboardShellProps> = ({ children, activeTab, onTabChange, theme, onToggleTheme, onLogout }) => {
    return (
        <div className="dashboard-container">
            {/* Sidebar */}
            <aside className="sidebar">
                <header className="sidebar-brand">
                    <div className="brand-icon">
                        <Brain size={28} />
                    </div>
                    <div className="flex flex-col">
                        <span className="brand-text">NebulaOS</span>
                        <span className="text-[10px] text-muted font-bold uppercase tracking-[0.2em] -mt-1 opacity-60">
                            Enterprise Cloud
                        </span>
                    </div>
                    <button
                        className="btn-secondary p-2 ml-auto"
                        onClick={onToggleTheme}
                        style={{ padding: '0.6rem', borderRadius: '12px' }}
                    >
                        {theme === 'dark' ? <Sun size={18} /> : <Moon size={18} />}
                    </button>
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
                    <SidebarItem icon={Sun} label="Enterprise" active={activeTab === 'hierarchy'} onClick={() => onTabChange('hierarchy')} />
                    <SidebarItem icon={Server} label="Bare Metal" active={activeTab === 'baremetal'} onClick={() => onTabChange('baremetal')} />
                    <SidebarItem icon={Settings} label="Settings" active={activeTab === 'settings'} onClick={() => onTabChange('settings')} />
                </nav>

                <div className="mt-auto pt-6 border-t border-white/5">
                    <SidebarItem icon={LogOut} label="Sign Out" onClick={onLogout} />
                </div>
            </aside>

            {/* Main Content */}
            <main className="main-content">
                <div className="animate-fade-in h-min">
                    {children}
                </div>
            </main>
        </div>
    );
};

export default DashboardShell;
