# NebulaOS — Ideal Customer Profile (ICP)

## Segmento Primário: Governo Municipal (Prefeituras)

| Atributo | Descrição |
|---|---|
| **Porte** | 50k - 500k habitantes |
| **Orçamento de TI** | R$ 500k - R$ 5M/ano |
| **Infra atual** | Servidores físicos locais ou nuvem pública esporádica |
| **Dores** | Vendor lock-in (nuvens públicas), custo imprevisível, falta de controle sobre dados |
| **Ciclo de compra** | 6-12 meses (licitação) |
| **Decisor** | Secretário de TI + jurídico |
| **Canal** | Licitação via pregão ou integrador parceiro |
| **Restrições** | LGPD obrigatória, dados devem ficar no município ou no Brasil |
| **Tickets de entrada** | R$ 30k-80k (implantação) |
| **ROI esperado** | Redução de 40-60% no custo de infraestrutura vs. nuvem pública |

## Segmento Secundário: Instituições de Ensino (Universidades)

| Atributo | Descrição |
|---|---|
| **Porte** | 5k - 50k alunos |
| **Infra atual** | Laboratórios próprios + nuvem acadêmica |
| **Dores** | Orçamento limitado, necessidade de ambientes isolados por departamento |
| **Ciclo de compra** | 3-6 meses (dispensa de licitação para PE) |
| **Decisor** | Diretor de TI / Centros de pesquisa |
| **Diferencial NebulaOS** | Multi-tenancy nativo (departamentos isolados), billing interno |

## Segmento Terciário: Órgãos Estaduais / Autarquias

| Atributo | Descrição |
|---|---|
| **Porte** | 500+ servidores |
| **Infra atual** | Datacenter próprio ou contrato vigente com nuvem |
| **Dores** | Soberania de dados, conformidade com legislação estadual |
| **Ciclo de compra** | 12-24 meses |
| **Tickets** | R$ 200k+ |

## Anti-ICP (quem NÃO atender)

- Empresas privadas de pequeno porte (< 20 funcionários) — não justificam ciclo de venda
- Startups em nuvem pública (já resolvido por AWS/GCP)
- Organizações sem nenhum requisito de soberania ou conformidade
- Clientes com orçamento de TI < R$ 100k/ano

## Critérios de Qualificação (BANT)

| Critério | Pergunta |
|---|---|
| **Budget** | Orçamento disponível para infraestrutura de cloud? |
| **Authority** | Decisor tem poder de compra ou influencia a licitação? |
| **Need** | Existe dor explícita com vendor lock-in / custo / soberania? |
| **Timeline** | Há um prazo definido para migrar ou implantar? |

## Canais de Aquisição

1. **Integradores parceiros** — empresas de consultoria que vendem implantação
2. **Pregão eletrônico** — licitação para governo municipal/estadual
3. **Comunidade open source** — adoção orgânica por universidades
4. **Indicação** — um cliente governamental indica para outro
