# NebulaOS — Análise Competitiva

## Matriz Comparativa

| Funcionalidade | NebulaOS | OpenStack | CloudStack | Proxmox VE | Mist.io | VMware vCloud |
|---|---|---|---|---|---|---|
| **Licença** | Apache 2.0 | Apache 2.0 | Apache 2.0 | AGPLv3 | AGPLv3 | Proprietária |
| **Multi-cloud** | ✅ Nativo | ❌ Apenas OpenStack | ❌ Apenas CloudStack | ❌ Apenas Proxmox | ✅ Multi-cloud | ❌ Apenas VMware |
| **Multi-tenancy** | ✅ Organização > Depto > Projeto | ✅ Projeto | ✅ Conta | ⚠️ Via PVE | ⚠️ Limitado | ✅ Organização |
| **IAM estilo AWS** | ✅ (Keycloak + OPA) | ✅ (Keystone) | ❌ | ❌ | ❌ | ✅ (vIDM) |
| **Marketplace** | ✅ Blueprints | ✅ (Glance/Horizon) | ✅ (Templates) | ✅ (CT/VM templates) | ❌ | ✅ (VMware Marketplace) |
| **Billing** | ✅ Soberano (chargeback) | ⚠️ (Ceilometer) | ❌ | ❌ | ⚠️ | ✅ (vRealize) |
| **Soberania** | ✅ Design sovereign-first | ❌ Neutro | ❌ Neutro | ❌ Neutro | ❌ Neutro | ❌ Vendor lock-in |
| **Bare metal** | ✅ iPXE + Redfish | ✅ (Ironic) | ❌ | ❌ | ❌ | ❌ |
| **GSLB multi-região** | ✅ Nativo | ❌ | ❌ | ❌ | ❌ | ❌ |
| **AI Advisor** | ✅ Nativo | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Setup** | Docker Compose (5 min) | 8 serviços + rabbit + mysql (dias) | 4 serviços (horas) | Instalação apt (30 min) | Docker Compose | Enterprise (semanas) |

## Posicionamento

### Vantagens competitivas do NebulaOS

1. **Soberania como feature, não como add-on** — diferenciação única contra OpenStack e CloudStack
2. **Menor complexidade operacional** — Docker Compose vs. OpenStack (8+ serviços)
3. **Multi-provedor nativo** — ao contrário de OpenStack e Proxmox que só gerenciam o próprio hypervisor
4. **Billing inclusivo** — chargeback por departamento sem custo adicional
5. **AI Advisor + GSLB** — features que concorrentes open source não têm

### Fraquezas

1. **Ecossistema imaturo** — OpenStack tem dezenas de integradores, comunidade massiva
2. **Menos providers** — AWS/Azure/GCP stubs, não funcionais
3. **Sem track record** — nenhuma implantação de produção conhecida
4. **Time pequeno** — sustentabilidade do projeto é risco

## Estratégia Recomendada

| Concorrente | Estratégia |
|---|---|
| **OpenStack** | Não competir em escala; focar em simplicidade e setup 5 min |
| **Proxmox VE** | Complementar: NebulaOS orquestra múltiplos clusters Proxmox |
| **VMware** | Migração de clientes insatisfeitos com custo de licenciamento |
| **CloudStack** | Comparação direta: NebulaOS tem mais funcionalidades nativas |
| **Mist.io** | Foco em soberania como diferencial (Mist.io não tem) |
