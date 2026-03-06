import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import ResourceCard from './ResourceCard';
import { Cpu } from 'lucide-react';

describe('ResourceCard', () => {
    const defaultProps = {
        title: 'CPU Usage',
        value: '45',
        unit: 'vCPUs',
        icon: Cpu,
        trend: 10,
        color: 'blue'
    };

    it('renders the title and value correctly', () => {
        render(<ResourceCard {...defaultProps} />);

        expect(screen.getByText('CPU Usage')).toBeInTheDocument();
        expect(screen.getByText('45')).toBeInTheDocument();
        expect(screen.getByText('vCPUs')).toBeInTheDocument();
    });

    it('displays the trend correctly when positive', () => {
        render(<ResourceCard {...defaultProps} />);
        expect(screen.getByText('+10%')).toBeInTheDocument();
    });

    it('displays the trend correctly when negative', () => {
        render(<ResourceCard {...defaultProps} trend={-5} />);
        expect(screen.getByText('-5%')).toBeInTheDocument();
    });
});
