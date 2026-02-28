# NebulaOS Network & Edge Architecture

## 1. Objective
Provide a simplified, institutional-grade networking layer that automates the complexity of VPCs, Firewalls, and SSL/TLS management.

## 2. VPC & Subnet Abstraction
NebulaOS abstracts underlying networking (Proxmox SDN, AWS VPC) into a unified interface:
- **Private VPCs:** Isolation at the project level.
- **Visual Subnetting:** Guided allocation of CIDR blocks.
- **Firewall Builder:** A visual tool to define ingress/egress rules (Security Groups) that apply across VMs and Kubernetes workloads.

## 3. Edge Services: Reverse Proxy Engine
The "Nebula Edge" handles incoming traffic:
- **Routing:** Path-based or Host-based routing to VMs or K8s services.
- **Load Balancing:** Automated distribution of traffic across multiple instances.
- **Web Application Firewall (WAF):** Integrated basic protection against common attacks.

## 4. DNS & Domain Orchestration
- **Binding:** Simple flow to bind a custom domain to a NebulaOS resource.
- **DNS Integration:** Automated creation of A/CNAME records if managed by NebulaOS, or guidance for external management.

## 5. SSL/TLS & Certificate Automation
- **Let's Encrypt:** Native integration for zero-touch SSL certificate generation and auto-renewal.
- **Custom Certificates:** Support for uploading proprietary or institutional certificates.
- **Automated Renewal:** Periodic checks and background renewal of all certificates before expiration.

## 6. Integration Boundaries
- **Kubernetes:** Automated creation of Ingress resources and Cert-Manager synchronization.
- **VM Workloads:** Automated configuration of the host reverse proxy to point to private VM IPs.
