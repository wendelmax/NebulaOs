# NebulaOS — Calculadora de TCO (Total Cost of Ownership)

## Cenário Base: Prefeitura de 200k habitantes

- 50 VMs (média: 4 vCPU, 8GB RAM, 100GB cada)
- 10TB de armazenamento block
- 5TB de object storage
- 3 ambientes (dev, staging, production)
- 1 administrador de infraestrutura (R$ 12k/mês)

## Comparativo Anual

| Item | Nuvem Pública (AWS) | OpenStack (DIY) | NebulaOS |
|---|---|---|---|
| **Compute** (50 VMs) | R$ 240.000 | R$ 30.000 (servidores) | R$ 30.000 (servidores) |
| **Storage** (15TB) | R$ 36.000 | R$ 12.000 (discos) | R$ 12.000 (discos) |
| **Licenciamento** | R$ 0 | R$ 0 | R$ 30.000 (Suporte Enterprise) |
| **Engenharia** (setup) | R$ 0 | R$ 120.000 (3 meses de especialista OpenStack) | R$ 30.000 (implantação assistida) |
| **Operação** (admin) | R$ 60.000 (meio período) | R$ 144.000 (dedicado OpenStack) | R$ 60.000 (meio período) |
| **Network egress** | R$ 24.000 | R$ 0 | R$ 0 |
| **Total Ano 1** | **R$ 360.000** | **R$ 306.000** | **R$ 162.000** |
| **Total Ano 2+** | **R$ 324.000** | **R$ 186.000** | **R$ 102.000** |

## Economia NebulaOS vs. Nuvem Pública

| Horizonte | Economia |
|---|---|
| Ano 1 | **55%** (R$ 198.000) |
| Ano 2+ | **69%** (R$ 222.000/ano) |

## Premissas

- Servidores físicos: 3 units de R$ 10.000 cada (amortizados em 36 meses)
- Custo nuvem pública baseado em AWS us-east-1 (on-demand, sem reserva)
- Custo OpenStack inclui apenas hardware, não inclui complexidade operacional
- NebulaOS considera Suporte Enterprise + implantação assistida

## Planilha (formato TSV para copiar para Excel)

```
Categoria	AWS	OpenStack	NebulaOS
Compute	240000	30000	30000
Storage	36000	12000	12000
Licenciamento	0	0	30000
Setup	0	120000	30000
Operacao	60000	144000	60000
Egress	24000	0	0
Total_Ano1	360000	306000	162000
Total_Ano2+	324000	186000	102000
```
