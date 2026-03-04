# NebulaOS

NebulaOS is a sovereign cloud management platform designed for high-performance, multi-region infrastructure orchestration. The dashboard provides a glassmorphic, data-driven interface for managing compute, storage, networking, and intelligence assets.

## Features

- **Live Infrastructure Telemetry**: Real-time stats for CPU, Storage, and Tenant activity.
- **Resource Inventory**: Unified management of cross-provider compute and storage assets.
- **Network Security & GSLB**: Intelligent multi-region traffic orchestration and VPC firewall management.
- **Marketplace**: One-click deployment of production-ready blueprints (K8s, Postgres, etc.).
- **AI Strategy Advisor**: Operational insights and performance optimization recommendations.
- **Sovereign Billing**: Transparent cost auditing with tenant-level granularity.

## Visual Walkthrough

### Dashboard Overview
The command center for your entire cloud infrastructure, featuring real-time telemetry and project status.
![Overview](docs/images/overview.png)

### Infrastructure Resources
A detailed inventory of all compute and storage assets, synchronized in real-time with the backend orchestration engine.
![Resources](docs/images/resources.png)

### Storage Orchestration
Dedicated management of block volumes and object storage buckets with provider-agnostic controls.
![Storage](docs/images/storage.png)

### Network Security & GSLB
Manage Global Traffic Strategy and security group rules for robust perimeter defense.
![Networking](docs/images/networking.png)

### Marketplace Blueprint Engine
Seamlessly deploy complex infrastructure stacks using pre-validated blueprints.
![Marketplace](docs/images/marketplace.png)

### AI Strategy Advisor
Intelligent insights that help optimize resource utilization and security posture.
![AI Advisor](docs/images/ai_advisor.png)

### Sovereign Billing & Analytics
Full transparency into infrastructure costs and consumption patterns.
![Billing](docs/images/billing.png)

## Development Setup

### Backend (Go)
```bash
cd src/api
go run cmd/server/main.go
```

### Frontend (React + Vite)
```bash
cd src/dashboard
npm install
npm run dev
```

The dashboard will be available at `http://localhost:5173`.
