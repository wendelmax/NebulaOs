import type React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, Server, Globe, Settings, LogOut, ChevronRight, HardDrive, Activity, Landmark, Receipt, Sun, ShoppingBag, Map, Brain, Languages } from 'lucide-react';
import { useLocale } from '../contexts/LocaleContext';

interface SidebarItemProps {
    icon: React.ComponentType<{ size?: number }>;
    label: string;
    to: string;
    active?: boolean;
}

const SidebarItem: React.FC<SidebarItemProps> = ({ icon: Icon, label, to, active = false }) => (
    <Link to={to} className={`sidebar-item ${active ? 'active' : ''}`}>
        <Icon size={20} />
        <span>{label}</span>
        {active && <ChevronRight size={16} style={{ marginLeft: 'auto' }} />}
    </Link>
);

interface DashboardShellProps {
    children: React.ReactNode;
    onLogout: () => void;
}

const navItems: { icon: React.ComponentType<{ size?: number }>; labelKey: string; tab: string }[] = [
    { icon: LayoutDashboard, labelKey: 'overview', tab: '/overview' },
    { icon: Server, labelKey: 'resources', tab: '/resources' },
    { icon: HardDrive, labelKey: 'storage', tab: '/storage' },
    { icon: Landmark, labelKey: 'governance', tab: '/governance' },
    { icon: Activity, labelKey: 'observability', tab: '/observability' },
    { icon: Receipt, labelKey: 'billing', tab: '/billing' },
    { icon: Globe, labelKey: 'networking', tab: '/networking' },
    { icon: Map, labelKey: 'global', tab: '/global' },
    { icon: ShoppingBag, labelKey: 'marketplace', tab: '/marketplace' },
    { icon: Brain, labelKey: 'advisor', tab: '/advisor' },
    { icon: Sun, labelKey: 'hierarchy', tab: '/hierarchy' },
    { icon: Server, labelKey: 'baremetal', tab: '/baremetal' },
    { icon: Settings, labelKey: 'settings', tab: '/settings' },
];

const DashboardShell: React.FC<DashboardShellProps> = ({ children, onLogout }) => {
    const { t, locale, setLocale } = useLocale();
    const location = useLocation();

    return (
        <div className="dashboard-container">
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
                </header>

                <nav className="sidebar-nav">
                    {navItems.map(item => (
                        <SidebarItem
                            key={item.tab}
                            icon={item.icon}
                            label={(t.nav as Record<string, string>)[item.labelKey] || item.labelKey}
                            to={item.tab}
                            active={location.pathname === item.tab}
                        />
                    ))}
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
                    <button className="sidebar-item w-full" onClick={onLogout} style={{ cursor: 'pointer' }}>
                        <LogOut size={20} />
                        <span>{t.nav.signOut}</span>
                    </button>
                </div>
            </aside>

            <main className="main-content">
                <div className="animate-fade-in h-min">
                    {children}
                </div>
            </main>
        </div>
    );
};

export default DashboardShell;
