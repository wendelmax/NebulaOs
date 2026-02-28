# NebulaOS Recommended Tech Stack

## 1. Core Services
*   **Language:** Go (Golang) for the Cloud API and Provider Plugins (performance, concurrency, and static typing).
*   **Database:** PostgreSQL with Row-Level Security (RLS) for multi-tenant data storage.
*   **Message Broker:** NATS for lightweight, high-performance internal event bus.

## 2. IAM & Security
*   **Identity Provider:** Keycloak (Open Source, OIDC, SAML, LDAP support).
*   **Policy Engine:** Open Policy Agent (OPA) for evaluating JSON-based access policies.
*   **Secrets Management:** HashiCorp Vault.

## 3. Infrastructure & Automation
*   **IaC Engine:** Terraform for initial provisioning and provider management.
*   **Orchestration:** Crossplane (Kubernetes-native infrastructure management) to provide a "Cloud API" feel.
*   **Configuration:** Ansible for VM-level configuration management.

## 4. Frontend & UX
*   **Framework:** React or Next.js with a custom Design System.
*   **Styling:** Vanilla CSS or CSS Modules (maximizing control and performance).
*   **API Client:** TypeScript for type safety across the stack.

## 5. Network & Edge
*   **Reverse Proxy:** Traefik or Nginx Ingress Controller.
*   **Certificates:** Cert-Manager (Let's Encrypt automation).
