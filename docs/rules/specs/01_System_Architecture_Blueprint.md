# NebulaOS -- System Architecture Design Prompt

## Objective

Design a modular architecture for NebulaOS as a Cloud Abstraction Layer
over:

-   Proxmox
-   OpenStack
-   Kubernetes
-   Bare Metal
-   AWS / Azure / GCP

## Must Include

-   Cloud Abstraction API Layer
-   IAM Layer (Open Source equivalent to AWS IAM)
-   Multi-Tenant Core
-   Network Abstraction (VPC, Subnet, Firewall)
-   Storage Abstraction (Block, Object, File)
-   Compute Abstraction (VM, K8s, Baremetal)
-   Certificate Management Layer
-   Domain & Reverse Proxy Orchestration
-   Observability Stack

## Output Required

-   Layered Architecture Diagram (textual)
-   API Boundaries
-   Service Responsibilities
-   Internal Communication Patterns
-   Recommended Tech Stack
-   Scalability Strategy
