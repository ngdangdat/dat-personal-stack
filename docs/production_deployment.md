# Production Deployment Guide

This document outlines the security, building, and orchestration guidelines for deploying the Engineering Assistant stack in a production environment.

## 1. Configure Production Secrets
Create a secure `.env` file on your target server. Do **NOT** use default values for secrets in production:
```env
PORT=8081
AGENT_SECRET_TOKEN=generate-a-long-secure-random-token-here
```

---

## 2. Build Containerized Images
The Dockerfiles are structured with **multi-stage builds** to minimize size and secure production runs:
* **Backend**: Compiles a minimal static binary inside `alpine` running under a non-root `appuser`.
* **Frontend**: Compiles Vite assets, then serves them using Nginx.

Deploy the container stack in detached daemon mode:
```bash
docker compose -f docker-compose.yml up -d
```

---

## 3. Deploying the Node Agent on a Remote Host (e.g. Raspberry Pi)
If the Node Agent runs on a separate target host to report remote metrics:
1. Compile the binary for the target architecture:
   ```bash
   # Example: Build for linux/arm64 (Raspberry Pi 4/5)
   cd backend
   CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o node_agent cmd/node_agent/main.go
   ```
2. Copy the binary to the remote host.
3. Start the daemon with environment secrets:
   ```bash
   PORT=8081 AGENT_SECRET_TOKEN=your-long-secure-secret-token ./node_agent
   ```

---

## 4. Reverse Proxy & SSL (HTTPS/WSS)
It is highly recommended to place a reverse proxy (such as Nginx, Traefik, or Caddy) in front of the ports to handle SSL/TLS termination:
* Route public web traffic to the frontend Nginx server (port `3000`).
* Proxy `/api` to the backend server (port `8080`).
* Proxy `/rpc` to the node agent WebSocket endpoint (port `8081`), ensuring WebSocket upgrade headers (`Upgrade` and `Connection`) are properly forwarded.
