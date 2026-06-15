# NebulaOS Marketplace — Modelo de Receita

## Estrutura de Comissionamento

| Tipo de Item | Taxa NebulaOS | Preço ao Cliente | Ganho do Desenvolvedor |
|---|---|---|---|
| Aplicativo pago | 20% | R$ 100/mês | R$ 80/mês |
| Plugin sob medida | 15% | R$ 5.000 (one-time) | R$ 4.250 |
| Template de VM/contêiner | Grátis | — | — (lead gen) |
| Suporte premium add-on | 10% | R$ 500/mês | R$ 450/mês |

## Fluxo Financeiro

```
Cliente paga → Stripe/PagSeguro → NebulaOS retém taxa → Repassa ao dev (Pix mensal)
```

## Requisitos para Publicar

| Requisito | Detalhe |
|---|---|
| Documentação | README + Screenshots + Vídeo (opcional) |
| Segurança | Scanner automático de vulnerabilidades |
| Compatibilidade | Testado contra versão estável do NebulaOS |
| SLA do desenvolvedor | 48h para responder issues críticas |

## Marketplace como Moat

- Plugins de **faturamento** (integração com PagSeguro, Nota Fiscal)
- Plugins de **LGPD** (anonimização, DSR automation)
- Plugins de **monitoração** (Zabbix, Grafana, Prometheus)
- Plugins de **backup** (Veeam, Bacula, restic)

> Quanto mais plugins úteis, maior o lock-in positivo da plataforma.

## Projeção de Receita

| Ano | Plugins Publicados | Receita Bruta | Taxa NebulaOS |
|---|---|---|---|
| 1 | 10 | R$ 60.000 | R$ 12.000 |
| 2 | 30 | R$ 300.000 | R$ 60.000 |
| 3 | 80 | R$ 1.200.000 | R$ 240.000 |

## Regras

1. **NebulaOS não criará plugins concorrentes pagos** — se fizermos um plugin, será gratuito/open-source.
2. Desenvolvedor pode definir preço, licença e suporte do seu plugin.
3. Remoção por violação de segurança ou LGPD.
