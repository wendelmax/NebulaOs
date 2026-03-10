import React, { useState } from 'react';
import { api } from '../api/client';
import { jwtDecode } from 'jwt-decode';
import { Lock, User, Rocket, ShieldCheck, AlertCircle } from 'lucide-react';

interface LoginViewProps {
    onLogin: (token: string, mustChange: boolean) => void;
}

const LoginView: React.FC<LoginViewProps> = ({ onLogin }) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        
        try {
            const response = await api.login({ username, password });
            const token = response.data.token;
            const decoded: any = jwtDecode(token);
            onLogin(token, decoded.must_change_password);
        } catch (err: any) {
            setError(err.response?.data || 'Access denied. Please check your credentials.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="login-container">
            <div className="login-glass-card glass-morphism">
                <div className="login-header">
                    <div className="login-logo pulsate">
                        <Rocket size={40} className="text-accent-primary" />
                    </div>
                    <h1>NebulaOS</h1>
                    <p className="subtitle">Next-Gen Infrastructure Orchestrator</p>
                </div>

                <form onSubmit={handleSubmit} className="login-form">
                    {error && (
                        <div className="login-error-badge">
                            <AlertCircle size={18} />
                            <span>{error}</span>
                        </div>
                    )}

                    <div className="input-group-premium">
                        <label>Identity</label>
                        <div className="input-with-icon">
                            <User size={18} />
                            <input 
                                type="text" 
                                placeholder="Username" 
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                required
                            />
                        </div>
                    </div>

                    <div className="input-group-premium">
                        <label>Credentials</label>
                        <div className="input-with-icon">
                            <Lock size={18} />
                            <input 
                                type="password" 
                                placeholder="Password" 
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
                        </div>
                    </div>

                    <button 
                        type="submit" 
                        className="btn-premium btn-primary full-width"
                        disabled={loading}
                    >
                        {loading ? 'Authenticating...' : 'Enter Dashboard'}
                    </button>
                </form>

                <div className="login-footer">
                    <div className="security-badge">
                        <ShieldCheck size={14} />
                        <span>Production Hardened v14.2</span>
                    </div>
                </div>
            </div>
            
            {/* Background elements */}
            <div className="bg-glow-1"></div>
            <div className="bg-glow-2"></div>
        </div>
    );
};

export default LoginView;
