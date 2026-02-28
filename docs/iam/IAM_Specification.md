# NebulaOS IAM Design Specification

## 1. Objective
Design a multi-tenant Identity and Access Management (IAM) system for NebulaOS that provides the same granularity and power as AWS IAM, but using open-source foundations (Keycloak + OPA).

## 2. Policy Model (AWS-style JSON)
NebulaOS uses a JSON-based policy engine. Policies can be attached to Users, Groups, or Roles.

### Example Policy:
```json
{
  "Version": "2026-02-27",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["compute:CreateVM", "compute:ListVMs"],
      "Resource": ["arn:nebula:compute:region-1:tenant-id:vm/*"],
      "Condition": {
        "IpAddress": {"nebula:SourceIp": "10.0.0.0/24"}
      }
    }
  ]
}
```

## 3. Tenant Isolation Strategy
Tenant isolation is enforced at multiple levels:
- **Identity Level:** Each Organization is a separate Realm in Keycloak.
- **Data Level:** Every API request must include a `Tenant-ID`. The Cloud API validates that the authenticated subject has access to that Tenant.
- **Resource Level:** Resources (VMs, Networks) are tagged with `nebula:owner` (Tenant ID) and `nebula:project`.

## 4. Security Model
- **RBAC (Role-Based Access Control):** For high-level permissions (e.g., "Administrator", "Viewer").
- **ABAC (Attribute-Based Access Control):** Using Open Policy Agent (OPA) to evaluate complex conditions like time-of-day, IP range, or resource tags.
- **Service Accounts:** For automation and CI/CD, using long-lived API tokens or OIDC client credentials.

## 5. API Design
- `POST /iam/policies`: Create a new policy.
- `POST /iam/roles`: Create a role with attached policies.
- `POST /iam/auth/token`: Exchange credentials or keys for a short-lived bearer token.
- `GET /iam/audit`: Retrieve immutable access logs for compliance.

## 6. Technology Integration
- **Keycloak:** Stores users, groups, and handles authentication (OIDC).
- **Nebula Policy Engine:** A custom Go service that validates IAM JSON policies using OPA for decision making.
- **PostgreSQL:** Stores IAM metadata and audit logs with strict RLS policies.
