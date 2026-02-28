import React from 'react';
import { motion } from 'framer-motion';
import { LayoutDashboard } from 'lucide-react';

interface ResourceCardProps {
    title: string;
    value: string | number;
    unit: string;
    icon: typeof LayoutDashboard;
    trend?: number;
    color?: string;
}

const ResourceCard: React.FC<ResourceCardProps> = ({ title, value, unit, icon: Icon, trend, color = 'var(--primary)' }) => {
    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            whileHover={{ scale: 1.02 }}
            className="resource-card glass glass-hover"
        >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                <div className="card-icon-wrapper" style={{ backgroundColor: `rgba(99, 102, 241, 0.1)`, color: color }}>
                    <Icon size={24} />
                </div>
                {trend !== undefined && (
                    <span style={{
                        fontSize: '0.75rem',
                        padding: '0.25rem 0.5rem',
                        borderRadius: '999px',
                        backgroundColor: trend > 0 ? 'rgba(16, 185, 129, 0.1)' : 'rgba(244, 63, 94, 0.1)',
                        color: trend > 0 ? '#34d399' : '#fb7185'
                    }}>
                        {trend > 0 ? '+' : ''}{trend}%
                    </span>
                )}
            </div>

            <div>
                <h3 className="card-title">{title}</h3>
                <div className="card-value-group">
                    <span className="card-value">{value}</span>
                    <span style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>{unit}</span>
                </div>
            </div>

            <div className="progress-bar-bg">
                <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: '65%' }}
                    transition={{ duration: 1, delay: 0.5 }}
                    className="progress-bar-fill"
                    style={{ backgroundColor: color }}
                />
            </div>
        </motion.div>
    );
};

export default ResourceCard;
