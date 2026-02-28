# NebulaOS Infrastructure Automation Model

## 1. Objective
Define the automation backbone that allows NebulaOS to manage infrastructure as code (IaC) with high reliability and visibility.

## 2. Infrastructure as Code (IaC) Stack
- **Terraform:** For managing low-level provider resources (VMs on Proxmox, VPCs on AWS).
- **Crossplane / Juju:** For higher-level orchestration and providing a Kubernetes-native infrastructure API.
- **Ansible / AWX:** For late-stage configuration of VM-based workloads (installing software, hardening OS).

## 3. GitOps Workflow (ArgoCD)
NebulaOS internal components and customer "Infrastructure Blueprints" are managed via GitOps:
1. **User Action:** User selects a "Blueprint" in the Console.
2. **API Action:** The Cloud API generates a Git repository or directory containing the specific IaC manifest.
3. **ArgoCD Action:** ArgoCD detects the change and synchronizes the state with the underlying cluster/provider.
4. **Drift Detection:** Any manual changes to the infrastructure are automatically detected and reverted (if policy allows).

## 4. Execution Modes
- **Wizard Mode (UI):** Simplified, opinionated forms for non-technical users.
- **Advanced Mode (CLI/IaC):** Direct access to manifest files and Terraform state for engineers.

## 5. State Management & Rollbacks
- **State Storage:** Encrypted Terraform state stored in the NebulaOS metadata DB.
- **Rollback Strategy:** Atomic updates using the Git commit history as the source of truth. If a sync fails, ArgoCD can roll back to the last known-good commit.

## 6. Multi-Cloud Abstraction Method
NebulaOS uses "Resource Classes". A single "Nebula:Compute" resource can be scheduled on a "Proxmox:Node" or "AWS:EC2" based on the tenant's policy and regional locking.
