# NebulaOS — Integração com Ecossistemas Governamentais

## 1. GOV.BR (Login Único)

### O que é
Plataforma de identidade digital do governo federal — permite login com CPF, certificado digital (ICP-Brasil), ou contas gov.br (níveis prata/ouro).

### Integração Técnica

```
Usuário                    NebulaOS                  GOV.BR
  │                          │                         │
  │  "Entrar com gov.br"     │                         │
  │ ───────────────────────► │                         │
  │                          │  Redireciona para       │
  │                          │  autorização OIDC       │
  │                          │ ───────────────────────► │
  │                          │                         │
  │                          │  Login + consentimento  │
  │ ◄─────────────────────── │                         │
  │                          │                         │
  │                          │  Callback com           │
  │                          │  authorization code     │
  │                          │ ◄─────────────────────── │
  │                          │                         │
  │                          │  Troca code por token   │
  │                          │ ───────────────────────► │
  │                          │                         │
  │                          │  ID Token + Access Token│
  │                          │ ◄─────────────────────── │
  │                          │                         │
  │  Autenticado via gov.br  │                         │
  │ ◄─────────────────────── │                         │
```

### Fluxo de Implementação

| Etapa | Descrição | Esforço |
|---|---|---|
| 1 | Criar aplicação no portal de desenvolvedor gov.br | 1 dia (burocracia) |
| 2 | Implementar OIDC RP no Keycloak (realm dedicado) | 2 dias |
| 3 | Mapear atributos gov.br → usuário NebulaOS (CPF, nome, e-mail) | 1 dia |
| 4 | Adicionar botão "Entrar com gov.br" na tela de login | 1 dia |
| 5 | Testar com contas nível prata/ouro | 2 dias |
| **Total** | | **~7 dias** |

### Pré-requisitos

- Cadastro como "Software Público" ou convênio com órgão público parceiro
- Política de privacidade e termos de uso (LGPD)
- HTTPS obrigatório (gov.br não redireciona para HTTP)

---

## 2. SEI (Sistema Eletrônico de Informações)

### O que é
Sistema de tramitação de documentos usado por mais de 1.000 órgãos públicos brasileiros.

### Integração (para casos de uso de governança)

NebulaOS não precisa se integrar ao SEI diretamente, mas pode:

| Funcionalidade | Como fazer |
|---|---|
| Notificar conclusão de implantação via SEI | Webhook → API REST do SEI (via SEI Plug-In) |
| Anexar relatório de compliance no SEI | Exportar PDF do relatório + API de anexo |
| Registrar solicitação de recurso como processo SEI | Formulário gera número de protocolo |

### Esforço estimado: 5 dias por integração pontual

---

## 3. SICAF (Sistema de Cadastramento Unificado de Fornecedores)

### O que é
Cadastro obrigatório para fornecer para o governo federal.

### Ação necessária

- [ ] Cadastrar CNPJ no SICAF (categoria: "Serviços de TI — Desenvolvimento de Software")
- [ ] Manter documentação fiscal e trabalhista atualizada
- [ ] Obter certidões: FGTS, INSS, tributos federais, trabalhista
- [ ] Cadastro no SICAF é **pré-requisito para qualquer licitação federal**

---

## 4. Assinatura Eletrônica (ICP-Brasil)

### O que é
Certificado digital padrão do governo brasileiro (A1 ou A3) para assinar documentos.

### Uso no NebulaOS

| Funcionalidade | Implementação |
|---|---|
| Assinar relatórios de compliance | Integração com APIs de assinatura (ex: ZapSign, Clicksign, DocuSign BR) |
| Autorizar operações críticas | Challenge com certificado digital antes de ações destrutivas |
| Autenticação MFA | Certificado A3 como segundo fator |

---

## 5. eMAG (Modelo de Acessibilidade em Governo Eletrônico)

### O que é
Padrão obrigatório para sites e portais do governo brasileiro (baseado em WCAG 2.1).

### Checklist para NebulaOS

| Requisito eMAG | Status | Prioridade |
|---|---|---|
| Contraste mínimo 4.5:1 | ⚠️ Parcial (modo escuro tem contraste baixo) | Alta |
| Navegação por teclado | ❌ Sidebar não é totalmente navegável por Tab | Alta |
| Labels em formulários | ❌ Login e wizards sem aria-labels | Média |
| Texto alternativo | ❌ Ícones sem aria-label | Baixa |
| Pular para conteúdo | ❌ Falta skip-to-content link | Baixa |
| Idioma da página | ✅ Já tem locale system | Feito |

---

## 6. SISP (Sistema de Administração de Recursos de Tecnologia da Informação)

### O que é
Sistema que regula compras de TI no governo federal.

### Requisitos para Software

- [ ] Alinhamento com PDTIC (Plano Diretor de TI) do órgão
- [ ] Conformidade com padrões de interoperabilidade (e-PING)
- [ ] Adoção de software livre é preferencial (Lei 12.527)
- [ ] Análise de viabilidade técnica e econômica

---

## Prioridade de Implementação

| Integração | Impacto | Esforço | Prioridade |
|---|---|---|---|
| **GOV.BR (login)** | 🔑 Desbloqueia adoção em governo federal | 5-7 dias | 🔴 Alta |
| **SICAF** | 📋 Permite participar de licitações | 2 dias (burocrático) | 🔴 Alta |
| **eMAG (acessibilidade)** | ♿ Obrigatório para sites gov | 5 dias | 🟡 Média |
| **SEI** | 📄 Integração secundária | 5 dias por ponto | 🟢 Baixa |
| **ICP-Brasil** | 🔐 Funcionalidade avançada | 10 dias | 🟢 Baixa |
