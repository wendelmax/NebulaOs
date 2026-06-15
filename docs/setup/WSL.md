# NebulaOS no Windows Subsystem for Linux (WSL2)

## Pré-requisitos

- Windows 11 (recomendado) ou Windows 10 (build 19041+)
- [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/) com integração WSL2 ativada
- WSL2 instalado com uma distribuição Linux (Ubuntu 22.04 LTS recomendada)

## Configuração

### 1. Instalar WSL2

```powershell
# Windows PowerShell (Admin)
wsl --install -d Ubuntu-22.04
wsl --set-default-version 2
```

### 2. Configurar Docker Desktop

1. Abrir Docker Desktop → Settings → Resources → WSL Integration
2. Ativar integração com a distribuição Ubuntu instalada
3. Aplicar e restartar Docker Desktop

### 3. Clonar o projeto no WSL

```bash
# Dentro do WSL (Ubuntu)
cd ~
git clone https://github.com/wendelmax/nebulaos.git
cd nebulaos
```

### 4. Configurar hosts file

O NebulaOS usa os domínios `nebula.local` e `api.nebula.local`. Eles precisam ser resolvidos tanto no Windows quanto no WSL.

**No Windows (PowerShell Admin):**
```powershell
Add-Content C:\Windows\System32\drivers\etc\hosts "`n127.0.0.1 nebula.local api.nebula.local"
```

**No WSL (/etc/hosts):**
```bash
echo "127.0.0.1 nebula.local api.nebula.local" | sudo tee -a /etc/hosts
```

### 5. Configurar variáveis de ambiente

```bash
cp .env.example .env
# Editar .env se necessário
```

### 6. Rodar o NebulaOS

```bash
./nebula.sh
```

O script detecta o Docker, carrega o `.env` e sobe todos os containers.

## Acesso

| Serviço | URL |
|---|---|
| Dashboard | http://nebula.local |
| API Health | http://api.nebula.local/health |
| Traefik Dashboard | http://localhost:8081 |
| Keycloak Admin | http://localhost:8080 (admin/admin) |

## Troubleshooting

### Porta 80 ocupada (IIS ou outro serviço)
No `.env`:
```env
TRAEFIK_ENTRYPOINT_PORT=8080
NEBULA_DOMAIN=localhost
```
Depois acessar http://localhost:8080 para o dashboard.

### Bind mount permission denied
No WSL, arquivos do Windows podem ter permissão restrita. Prefira clonar o projeto **dentro** do WSL (`~/nebulaos`) em vez de em `C:\Users\...`.

### DNS não resolve nebula.local
Verificar `/etc/hosts` dentro do WSL (não adianta configurar só no Windows — o WSL tem seu próprio resolução de DNS).

### Docker não encontrado no WSL
No Docker Desktop: Settings → Resources → WSL Integration → ativar sua distro. Depois restartar o terminal WSL.
