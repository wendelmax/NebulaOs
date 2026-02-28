import React from 'react';
import { Box, Database, ShieldCheck, Zap, RefreshCw } from 'lucide-react';

const Marketplace: React.FC = () => {
    const [blueprints, setBlueprints] = React.useState<any[]>([]);
    const [loading, setLoading] = React.useState(true);
    const [deploying, setDeploying] = React.useState<string | null>(null);

    const fetchBlueprints = async () => {
        setLoading(true);
        try {
            const resp = await fetch('http://api.nebula.local/marketplace/blueprints');
            if (resp.ok) {
                const data = await resp.json();
                setBlueprints(data || []);
            }
        } catch (err) {
            console.error("Failed to fetch blueprints", err);
        } finally {
            setLoading(false);
        }
    };

    React.useEffect(() => {
        fetchBlueprints();
    }, []);

    const handleDeploy = async (bpId: string) => {
        setDeploying(bpId);
        try {
            const resp = await fetch('http://api.nebula.local/marketplace/deploy', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    blueprint_id: bpId,
                    project_id: 'v-p1'
                })
            });
            if (resp.ok) {
                alert('Deployment started successfully!');
            }
        } catch (err) {
            console.error(err);
        } finally {
            setDeploying(null);
        }
    };

    const getBlueprintIcon = (category: string) => {
        switch (category?.toLowerCase()) {
            case 'infrastructure': return Box;
            case 'databases': return Database;
            case 'security': return ShieldCheck;
            default: return Zap;
        }
    };

    return (
        <div className="view-container animate-fade-in">
            <header className="view-header">
                <div>
                    <h1>Cloud Marketplace</h1>
                    <p className="text-muted">Launch production-ready infrastructure blueprints in seconds.</p>
                </div>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <button className="btn-secondary" onClick={fetchBlueprints} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        Refresh
                    </button>
                </div>
            </header>

            <div className="stats-grid" style={{ marginTop: '2rem' }}>
                {blueprints.length === 0 && !loading && (
                    <div className="stat-card glass" style={{ gridColumn: 'span 3', textAlign: 'center', padding: '4rem' }}>
                        <p className="text-muted">No blueprints available in the marketplace.</p>
                    </div>
                )}
                {blueprints.map(bp => {
                    const Icon = getBlueprintIcon(bp.category);
                    return (
                        <div key={bp.id} className="stat-card glass-hover" style={{ cursor: 'pointer', display: 'flex', flexDirection: 'column', gap: '1rem', padding: '2rem' }}>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                <div className="stat-icon" style={{ background: 'var(--primary-gradient)', color: 'white' }}>
                                    <Icon size={24} />
                                </div>
                                <span className="badge badge-success">{bp.category}</span>
                            </div>
                            <div>
                                <h3 style={{ fontSize: '1.25rem' }}>{bp.name}</h3>
                            </div>
                            <p style={{ fontSize: '0.9rem', color: 'var(--text-muted)', lineHeight: '1.6', flex: 1 }}>
                                {bp.description}
                            </p>
                            <button
                                className="btn-primary"
                                style={{ marginTop: 'auto', padding: '0.75rem', fontSize: '0.9rem' }}
                                onClick={() => handleDeploy(bp.id)}
                                disabled={deploying === bp.id}
                            >
                                {deploying === bp.id ? 'Deploying...' : 'Deploy Blueprint'}
                            </button>
                        </div>
                    );
                })}
            </div>
        </div>
    );
};

export default Marketplace;
