import React, { useState } from 'react';
import { api } from '../api/client';
import { KeyRound, Mail, AlertCircle, CheckCircle2, ShieldAlert } from 'lucide-react';
import { useLocale } from '../contexts/LocaleContext';

interface ChangePasswordViewProps {
    onComplete: () => void;
}

const ChangePasswordView: React.FC<ChangePasswordViewProps> = ({ onComplete }) => {
    const { t } = useLocale();
    const [oldPassword, setOldPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [email, setEmail] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        
        if (newPassword !== confirmPassword) {
            setError(t.changePassword.mismatch);
            return;
        }

        if (newPassword.length < 6) {
            setError(t.changePassword.tooShort);
            return;
        }

        setLoading(true);
        try {
            await api.changePassword({ 
                old_password: oldPassword, 
                new_password: newPassword, 
                email 
            });
            setSuccess(true);
            setTimeout(() => {
                onComplete();
            }, 2000);
        } catch (err: any) {
            setError(err.response?.data || t.changePassword.failed);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="login-container">
            <div className="login-glass-card glass-morphism max-w-md">
                <div className="login-header">
                    <div className="login-logo text-accent-primary animate-bounce-slow">
                        <ShieldAlert size={40} />
                    </div>
                    <h1>{t.changePassword.title}</h1>
                    <p className="subtitle text-warning">{t.changePassword.description}</p>
                </div>

                {success ? (
                    <div className="security-success-overlay animate-fade-in">
                        <CheckCircle2 size={64} className="text-accent-secondary mb-4" />
                        <h2>{t.changePassword.success}</h2>
                        <p>{t.changePassword.redirecting}</p>
                    </div>
                ) : (
                    <form onSubmit={handleSubmit} className="login-form">
                        {error && (
                            <div className="login-error-badge">
                                <AlertCircle size={18} />
                                <span>{error}</span>
                            </div>
                        )}

                        <div className="input-group-premium">
                        <label>{t.changePassword.currentLabel}</label>
                        <div className="input-with-icon">
                            <KeyRound size={18} />
                            <input 
                                type="password" 
                                placeholder={t.changePassword.currentPlaceholder} 
                                    value={oldPassword}
                                    onChange={(e) => setOldPassword(e.target.value)}
                                    required
                                />
                            </div>
                        </div>

                        <div className="divider-premium"></div>

                        <div className="input-group-premium">
                            <label>{t.changePassword.newLabel}</label>
                            <div className="input-with-icon">
                                <KeyRound size={18} />
                                <input 
                                    type="password" 
                                    placeholder={t.changePassword.newPlaceholder} 
                                    value={newPassword}
                                    onChange={(e) => setNewPassword(e.target.value)}
                                    required
                                />
                            </div>
                        </div>

                        <div className="input-group-premium">
                            <label>{t.changePassword.confirmLabel}</label>
                            <div className="input-with-icon">
                                <KeyRound size={18} />
                                <input 
                                    type="password" 
                                    placeholder={t.changePassword.confirmPlaceholder} 
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    required
                                />
                            </div>
                        </div>

                        <div className="input-group-premium">
                            <label>{t.changePassword.emailLabel}</label>
                            <div className="input-with-icon">
                                <Mail size={18} />
                                <input 
                                    type="email" 
                                    placeholder={t.changePassword.emailPlaceholder} 
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    required
                                />
                            </div>
                        </div>

                        <button 
                            type="submit" 
                            className="btn-premium btn-primary full-width mt-4"
                            disabled={loading}
                        >
                            {loading ? t.changePassword.submitting : t.changePassword.submit}
                        </button>
                    </form>
                )}
            </div>
            
            <div className="bg-glow-1 opacity-50"></div>
            <div className="bg-glow-secret"></div>
        </div>
    );
};

export default ChangePasswordView;
