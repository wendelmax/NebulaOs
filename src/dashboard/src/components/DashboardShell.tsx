import { LayoutDashboard, Server, Globe, Settings, LogOut, ChevronRight, HardDrive, Activity, Landmark, Receipt, Sun, Moon, ShoppingBag, Map, Brain, Languages } from 'lucide-react';
import { useLocale } from '../contexts/LocaleContext';

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
    const { t, locale, setLocale } = useLocale();

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
                            {t.app.tagline}
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
                    <SidebarItem icon={LayoutDashboard} label={t.nav.overview} active={activeTab === 'overview'} onClick={() => onTabChange('overview')} />
                    <SidebarItem icon={Server} label={t.nav.resources} active={activeTab === 'resources'} onClick={() => onTabChange('resources')} />
                    <SidebarItem icon={HardDrive} label={t.nav.storage} active={activeTab === 'storage'} onClick={() => onTabChange('storage')} />
                    <SidebarItem icon={Landmark} label={t.nav.governance} active={activeTab === 'governance'} onClick={() => onTabChange('governance')} />
                    <SidebarItem icon={Activity} label={t.nav.observability} active={activeTab === 'observability'} onClick={() => onTabChange('observability')} />
                    <SidebarItem icon={Receipt} label={t.nav.billing} active={activeTab === 'billing'} onClick={() => onTabChange('billing')} />
                    <SidebarItem icon={Globe} label={t.nav.networking} active={activeTab === 'networking'} onClick={() => onTabChange('networking')} />
                    <SidebarItem icon={Map} label={t.nav.global} active={activeTab === 'global'} onClick={() => onTabChange('global')} />
                    <SidebarItem icon={ShoppingBag} label={t.nav.marketplace} active={activeTab === 'marketplace'} onClick={() => onTabChange('marketplace')} />
                    <SidebarItem icon={Brain} label={t.nav.advisor} active={activeTab === 'advisor'} onClick={() => onTabChange('advisor')} />
                    <SidebarItem icon={Sun} label={t.nav.hierarchy} active={activeTab === 'hierarchy'} onClick={() => onTabChange('hierarchy')} />
                    <SidebarItem icon={Server} label={t.nav.baremetal} active={activeTab === 'baremetal'} onClick={() => onTabChange('baremetal')} />
                    <SidebarItem icon={Settings} label={t.nav.settings} active={activeTab === 'settings'} onClick={() => onTabChange('settings')} />
                </nav>

                <div className="mt-auto pt-6 border-t border-white/5 space-y-2">
                    <button
                        className="sidebar-item w-full"
                        onClick={() => setLocale(locale === 'en' ? 'pt-BR' : 'en')}
                        style={{ cursor: 'pointer' }}
                    >
                        <Languages size={20} />
                        <span>{locale === 'en' ? 'Português (BR)' : 'English'}</span>
                    </button>
                    <SidebarItem icon={LogOut} label={t.nav.signOut} onClick={onLogout} />
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