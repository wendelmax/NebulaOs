import React, { useState } from 'react';
import DashboardShell from './components/DashboardShell';
import type { TabType } from './components/DashboardShell';
import Overview from './views/Overview';
import ResourcesView from './views/ResourcesView';
import StorageView from './views/StorageView';
import GovernanceView from './views/GovernanceView';
import ObservabilityView from './views/ObservabilityView';
import BillingView from './views/BillingView';
import Marketplace from './views/Marketplace';
import Networking from './views/Networking';
import GlobalTopology from './views/GlobalTopology';
import AIAdvisor from './views/AIAdvisor';
import SettingsView from './views/SettingsView';
import LoginView from './views/LoginView';
import ChangePasswordView from './views/ChangePasswordView';
import HierarchyView from './views/HierarchyView';
import BareMetalView from './views/BareMetalView';
import { api } from './api/client';
import { LocaleProvider, useLocale } from './contexts/LocaleContext';

const AppInner: React.FC = () => {
  const { t } = useLocale();
  const [activeTab, setActiveTab] = useState<TabType>('overview');
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [mustChangePassword, setMustChangePassword] = useState<boolean>(false);
  const [authLoading, setAuthLoading] = useState<boolean>(true);

  React.useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
  }, [theme]);

  React.useEffect(() => {
    api.checkAuth()
      .then(() => setIsAuthenticated(true))
      .catch(() => setIsAuthenticated(false))
      .finally(() => setAuthLoading(false));
  }, []);

  const toggleTheme = () => setTheme(prev => prev === 'dark' ? 'light' : 'dark');

  const handleLogin = (token: string, mustChange: boolean) => {
    setMustChangePassword(mustChange);
    setIsAuthenticated(true);
  };

  const handleLogout = async () => {
    try { await api.logout(); } catch { /* ignore */ }
    setIsAuthenticated(false);
  };

  const renderContent = () => {
    switch (activeTab) {
      case 'overview':
        return <Overview />;
      case 'resources':
        return <ResourcesView />;
      case 'storage':
        return <StorageView />;
      case 'governance':
        return <GovernanceView />;
      case 'observability':
        return <ObservabilityView />;
      case 'billing':
        return <BillingView />;
      case 'networking':
        return <Networking />;
      case 'marketplace':
        return <Marketplace />;
      case 'global':
        return <GlobalTopology />;
      case 'advisor':
        return <AIAdvisor />;
      case 'settings':
        return <SettingsView />;
      case 'hierarchy':
        return <HierarchyView />;
      case 'baremetal':
        return <BareMetalView />;
      default:
        return (
          <div className="glass p-12 text-center">
            <h2 style={{ color: 'var(--text-muted)' }}>{t.common.moduleInDevelopment}</h2>
            <p style={{ color: 'rgba(148, 163, 184, 0.6)', marginTop: '1rem' }}>
              {t.common.moduleInDevelopmentDesc}
            </p>
          </div>
        );
    }
  };

  if (authLoading) {
    return (
      <div className="flex items-center justify-center h-screen" style={{ background: 'var(--bg-primary)' }}>
        <div className="text-center">
          <div className="text-4xl font-black tracking-tight text-main/50 animate-pulse">NebulaOS</div>
          <div className="text-dim text-sm mt-4">{t.app.loading}</div>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <LoginView onLogin={handleLogin} />;
  }

  if (mustChangePassword) {
    return <ChangePasswordView onComplete={() => setMustChangePassword(false)} />;
  }

  return (
    <DashboardShell
      activeTab={activeTab}
      onTabChange={setActiveTab}
      theme={theme}
      onToggleTheme={toggleTheme}
      onLogout={handleLogout}
    >
      {renderContent()}
    </DashboardShell>
  );
}

const App: React.FC = () => (
  <LocaleProvider>
    <AppInner />
  </LocaleProvider>
);

export default App;