import React, { useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import DashboardShell from './components/DashboardShell';
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
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [mustChangePassword, setMustChangePassword] = useState<boolean>(false);
  const [authLoading, setAuthLoading] = useState<boolean>(true);
  const navigate = useNavigate();

  React.useEffect(() => {
    api.checkAuth()
      .then(() => setIsAuthenticated(true))
      .catch(() => setIsAuthenticated(false))
      .finally(() => setAuthLoading(false));
  }, []);

  const handleLogin = (token: string, mustChange: boolean) => {
    setMustChangePassword(mustChange);
    setIsAuthenticated(true);
    navigate(mustChange ? '/change-password' : '/overview');
  };

  const handleLogout = async () => {
    try { await api.logout(); } catch { /* ignore */ }
    setIsAuthenticated(false);
    navigate('/login');
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

  return (
    <Routes>
      <Route path="/login" element={
        !isAuthenticated
          ? <LoginView onLogin={handleLogin} />
          : <Navigate to="/overview" replace />
      } />
      <Route path="/change-password" element={
        isAuthenticated && mustChangePassword
          ? <ChangePasswordView onComplete={() => { setMustChangePassword(false); navigate('/overview'); }} />
          : !isAuthenticated
            ? <Navigate to="/login" replace />
            : <Navigate to="/overview" replace />
      } />
      <Route path="/*" element={
        isAuthenticated
          ? <DashboardShell onLogout={handleLogout}>
              <Routes>
                <Route path="/overview" element={<Overview />} />
                <Route path="/resources" element={<ResourcesView />} />
                <Route path="/storage" element={<StorageView />} />
                <Route path="/governance" element={<GovernanceView />} />
                <Route path="/observability" element={<ObservabilityView />} />
                <Route path="/billing" element={<BillingView />} />
                <Route path="/networking" element={<Networking />} />
                <Route path="/marketplace" element={<Marketplace />} />
                <Route path="/global" element={<GlobalTopology />} />
                <Route path="/advisor" element={<AIAdvisor />} />
                <Route path="/settings" element={<SettingsView />} />
                <Route path="/hierarchy" element={<HierarchyView />} />
                <Route path="/baremetal" element={<BareMetalView />} />
                <Route path="*" element={<Navigate to="/overview" replace />} />
              </Routes>
            </DashboardShell>
          : <Navigate to="/login" replace />
      } />
    </Routes>
  );
}

const App: React.FC = () => (
  <LocaleProvider>
    <BrowserRouter>
      <AppInner />
    </BrowserRouter>
  </LocaleProvider>
);

export default App;
