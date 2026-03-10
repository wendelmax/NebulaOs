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

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState<TabType>('overview');
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(!!localStorage.getItem('nebula_token'));
  const [mustChangePassword, setMustChangePassword] = useState<boolean>(false);

  React.useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
  }, [theme]);

  const toggleTheme = () => setTheme(prev => prev === 'dark' ? 'light' : 'dark');

  const handleLogin = (token: string, mustChange: boolean) => {
    localStorage.setItem('nebula_token', token);
    setMustChangePassword(mustChange);
    setIsAuthenticated(true);
  };

  const handleLogout = () => {
    localStorage.removeItem('nebula_token');
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
            <h2 style={{ color: 'var(--text-muted)' }}>Module In Development</h2>
            <p style={{ color: 'rgba(148, 163, 184, 0.6)', marginTop: '1rem' }}>
              The "{activeTab}" capability is being provisioned in the orchestration plane.
            </p>
          </div>
        );
    }
  };

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

export default App;
