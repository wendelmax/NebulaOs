# NebulaOS — Plano de Compliance e Certificações

## 1. LGPD (Lei Geral de Proteção de Dados) — Lei 13.709/2018

### Mapeamento de Controles

| Requisito LGPD | Controle NebulaOS | Status | Prioridade |
|---|---|---|---|
| **Art. 6 — Finalidade** | Registro de finalidade de tratamento por tenant | ❌ Não implementado | Alta |
| **Art. 7 — Base legal** | Consentimento registrado por usuário | ❌ Não implementado | Alta |
| **Art. 18 — Direitos do titular** | API para exportação e exclusão de dados pessoais | ⚠️ Parcial (export existe via nebula-export) | Alta |
| **Art. 37 — Segurança** | Criptografia em repouso (Postgres) + TLS nas APIs | ⚠️ Parcial (sem criptografia em repouso explícita) | Alta |
| **Art. 38 — Notificação de incidentes** | Sistema de alerta e logging de eventos de segurança | ❌ Não implementado | Média |
| **Art. 41 — Encarregado (DPO)** | Campo de contato do DPO por organização | ❌ Não implementado | Baixa |
| **Art. 46 — Princípios de segurança** | Controles de acesso, auditoria, e privacy by design | ⚠️ Parcial | Alta |

### Plano de Ação LGPD

| Mês | Ação | Entregável |
|---|---|---|
| 1 | Adicionar campo "purpose" no registro de tenant | Schema SQL + API |
| 2 | Implementar consentimento no fluxo de cadastro | Tela de cadastro com checkbox LGPD |
| 3 | Adicionar criptografia em repouso (AES-256 para secrets no Vault + TDE no Postgres) | Documento de arquitetura + config |
| 4 | Criar API de exclusão lógica de dados pessoais (right to be forgotten) | Endpoint /admin/data-privacy |
| 5 | Implementar logging de incidentes de segurança | Módulo de auditoria de segurança |
| 6 | Auditoria interna de conformidade LGPD | Relatório de conformidade |

## 2. ISO 27001 — Sistema de Gestão de Segurança da Informação

### Escopo para Certificação

NebulaOS como plataforma de orquestração cloud — aplicável a:
- Desenvolvimento e manutenção do software
- Infraestrutura de suporte (cloud ou on-prem)
- Processos de suporte e implantação

### Controles ABNT NBR ISO 27001:2022

| Anexo A | Controle | Implementação |
|---|---|---|
| **5.1** | Políticas de segurança | Criar política de segurança da informação |
| **5.15** | Segurança em contratos de parceiros | Contrato de suporte com cláusulas de segurança |
| **5.17** | Autenticação | JWT + Keycloak + MFA |
| **5.18** | Controle de acesso | RBAC + ABAC via OPA |
| **5.25** | Segurança em desenvolvimento | CI/CD com verificação de segurança |
| **5.29** | Segurança em aquisição | Due diligence em parceiros |
| **5.33** | Proteção de registros | Audit logging imutável |
| **5.36** | Conformidade com LGPD | Ver plano LGPD acima |
| **6.8** | Gerenciamento de vulnerabilidades | Scanner de dependências (dependabot) |
| **7.10** | Criptografia | TLS para trânsito, AES-256 para repouso |
| **7.14** | Continuidade de negócio | Backup e restore testado trimestralmente |
| **8.8** | Gestão de mudanças | Processo RFC + approval |

### Cronograma ISO 27001

| Fase | Duração | Atividades |
|---|---|---|
| Planejamento | 2 meses | Definição de escopo, política de segurança, análise de riscos |
| Implementação | 4 meses | Implantação dos controles, documentação, treinamento |
| Auditoria interna | 1 mês | Pré-auditoria, correção de não-conformidades |
| Auditoria externa | 2 meses | Contratação de certificadora (ex: BSI, DNV, QMS) |
| Certificação | — | Emissão do certificado ISO 27001 |

### Investimento Estimado

| Item | Custo |
|---|---|
| Consultoria especializada | R$ 50.000-80.000 |
| Ferramentas de segurança | R$ 10.000-20.000/ano |
| Auditoria externa | R$ 30.000-50.000 |
| **Total** | **R$ 90.000-150.000** |

## 3. Outras Certificações Relevantes

| Certificação | Relevância | Esforço |
|---|---|---|
| **MPS-PR (nível G)** | Governo brasileiro exige para softwares críticos | 6-12 meses |
| **SOC 2 Type II** | Requisito internacional para SaaS | 6-9 meses |
| **SISP (MGI)** | Conformidade com padrões de governo federal | 2-3 meses |
| **eMAG (Modelo de Acessibilidade)** | Obrigatório para sites governamentais | 1-2 meses |
