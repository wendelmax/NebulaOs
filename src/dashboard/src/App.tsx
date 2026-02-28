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

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState<TabType>('overview');
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');

  React.useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
  }, [theme]);

  const toggleTheme = () => setTheme(prev => prev === 'dark' ? 'light' : 'dark');

  const renderContent = () => {
    switch (activeTab) {
      case 'overview':
        return <Overview theme={theme} />;
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

  return (
    <DashboardShell
      activeTab={activeTab}
      onTabChange={setActiveTab}
      theme={theme}
      onToggleTheme={toggleTheme}
    >
      {renderContent()}
    </DashboardShell>
  );
}

export default App;
