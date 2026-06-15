# NebulaOS — Governança Comunitária

## Modelo: Open Core com Fundação

NebulaOS é open-core: funcionalidades core são abertas, complementos são comerciais.

Uma **Fundação NebulaOS** será criada quando houver 3+ mantenedores de empresas diferentes.

## Estrutura

```
┌─────────────────────────────────────┐
│         TSC (Technical Steering     │
│          Committee)                  │
│   Define roadmap, aprova RFCs       │
├─────────────────────────────────────┤
│         Mantenedores                 │
│   Review PRs, triage issues         │
├─────────────────────────────────────┤
│         Colaboradores                │
│   Contribuem com PRs, docs, bugs    │
├─────────────────────────────────────┤
│         Comunidade                   │
│   Usuários, parceiros,            │
│   entusiastas                       │
└─────────────────────────────────────┘
```

## TSC — Technical Steering Committee

### Composição
- **5 membros** (ideal: 1 do maintainer original, 2 de empresas parceiras, 2 da comunidade)
- Mandato de **1 ano**, reeleição ilimitada
- Decisões por consenso; se não houver, votação simples (maioria)

### Responsabilidades
- Aprovar ou rejeitar RFCs
- Definir prioridades do roadmap
- Nomear novos mantenedores
- Resolver disputas técnicas
- Definir política de segurança

## Processo de RFC

Toda mudança significativa no NebulaOS deve passar por RFC.

```
Idea → GitHub Discussion (1 semana) → Draft RFC (PR no repositório)
→ Período de comentários (2 semanas) → Votação TSC → Aprovado/Rejeitado
```

### O que precisa de RFC

- Mudanças que quebram compatibilidade
- Novos provedores core (não plugin)
- Mudanças na API pública
- Alterações na segurança/autenticação
- Mudanças na licença

### O que NÃO precisa de RFC

- Bug fixes
- Melhorias de performance sem mudança de API
- Testes e documentação
- Plugins do marketplace

## Diretrizes de Contribuição

1. **Fork + Branch**: feature/nome-da-feature ou fix/issue-123
2. **Commits**: mensagens em inglês, descritivas
3. **PRs**: referenciar issue, descrever mudança, screenshots se UI
4. **Code review**: 2 approvals de mantenedores (1 para docs/tests)
5. **CI**: todos os checks passando
6. **Testes**: novos testes obrigatórios

## Código de Conduta

Contribuidores devem seguir o [Contributor Covenant 2.1](https://www.contributor-covenant.org/).

## Marcas e Logotipos

- Nome "NebulaOS" e logotipo são marcas registradas do mantenedor original
- Uso permitido para referenciar o projeto (fair use)
- Parceiros certificados podem usar selos oficiais

## Releases

| Versão | Frequência | Suporte |
|---|---|---|
| **Edge** | Semanal (main) | ❌ |
| **Stable** | Trimestral | 6 meses de patches |
| **LTS** | Anual | 2 anos de patches de segurança |

## Repositórios

| Repositório | Descrição |
|---|---|
| nebulaos/core | API + provedores core |
| nebulaos/dashboard | Frontend React |
| nebulaos/plugins | Plugins marketplace |
| nebulaos/docs | Documentação geral |
| nebulaos/rfcs | Propostas de mudança |

## Financiamento da Comunidade

- **Open Collective** (ou similar): doações voluntárias para custear infraestrutura
- **Sponsors**: empresas que bancam features específicas (mencionadas no README)
- **Contribuições em espécie**: AWS/GCP credits, servidores para testes
