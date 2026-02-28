# NebulaOS UI/UX Design System

## 1. Design Philosophy
NebulaOS aims to provide a "Premium Institutional Experience." The interface should feel as modern as DigitalOcean but with the gravitas required for governmental and non-profit use.

### Core Principles:
- **Clarity over Complexity:** Hide technical details by default, but make them accessible via an "Advanced Toggle."
- **Guided Workflows:** Use multi-step wizards for all provisioning tasks.
- **Visual Feedback:** Provide immediate, clear visual confirmation of state changes (e.g., "VM is starting...").
- **Accessibility:** Full compliance with WCAG 2.1 standards for public sector usage.

## 2. Navigation Structure
- **Global Sidebar:**
    - **Dashboard:** Overview of resources and health.
    - **Computing:** VMs, Kubernetes Clusters.
    - **Networking:** VPCs, DNS, Load Balancers.
    - **Storage:** Block Storage, Backups.
    - **Identity:** Users, Groups, Roles (IAM).
    - **Settings:** Organization profile and Quotas.

## 3. The "Wizard" Experience
NebulaOS provides highly opinionated wizards for common tasks:
- **VM Wizard:** 1. Select Blueprint (OS/Apps) -> 2. Select Size -> 3. Select Network -> 4. Review & Launch.
- **K8s Wizard:** Managed cluster provisioning with integrated monitoring and auto-scaling options.
- **VPC visual builder:** Drag-and-drop or checklist-based network topology design.

## 4. Usage & Cost Dashboard
Even in non-commercial environments, transparency of resource usage is critical:
- **Real-time Metrics:** CPU/RAM/Disk consumption.
- **Departmental Chargeback:** Breakdown of resource costs (theoretical) by department or project.
- **Forecasting:** Predictive analysis of quota usage.

## 5. Simple Mode vs. Advanced Mode
- **Simple Mode (Default):** Pre-configured sizes, automated firewalling, and wizard-led navigation.
- **Advanced Mode:** JSON editor for policies, raw YAML for K8s manifests, and direct access to Terraform variables.
