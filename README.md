# Engineering Assistant: Mobile Command Center

A mobile-first, high-density dashboard and remote orchestrator built for modern software engineers. It bridges the gap between high-level project management (GitHub PRs) and low-level execution environments (Remote Raspberry Pi, cloud instances, and local micro-environments).

---

## 🚀 Core Features & Status

### 1. PR Dashboard (Phase 1) — `[COMPLETED]`
* **Direct Integration**: OAuth authentication with GitHub REST/GraphQL APIs.
* **Smart Filter Views**: High-density lists for "Reviewing" (assigned PRs) and "My PRs" (authored PRs).
* **Review Velocity Metrics**: Clear indicator badges showing lead times and status indicators (Draft, Approved, Changes Requested).

### 2. Workspace Manager & RPC Orchestration (Phase 2) — `[COMPLETED]`
* **Node Telemetry**: Live telemetry metrics (CPU load, memory capacity/distribution, goroutines, uptime) via WebSocket custom JSON-RPC 2.0 daemon.
* **Command Shell Streaming**: Execute repository commands directly from the UI and stream stdout/stderr logs in real-time.
* **Auth & Handshake Security**: Subprotocol headers validation, Origin checking, and automatic execution context cancellation.
* **Offline-First Storage**: Local persistent caching of registered nodes utilizing IndexedDB.

### 3. Project Detail & Terminal (Phase 3) — `[UPCOMING]`
* **Integrated TUI Console**: Dense repository tracking displaying active branches, last commits, and live shell logs.
* **Interactive Shortcuts**: Action buttons for instant command invocation (`git pull`, `build`, `deploy`).

### 4. Assistant & Local Storage (Phase 4) — `[UPCOMING]`
* **AI Chat Companion**: An on-device assistant designed to track project status and help troubleshoot failures.
* **Offline-First Storage**: Synchronization of assistant logs and system state to IndexDB.

---

## 🎨 Visual Identity: Terminal Precision

The interface is built strictly around the **Terminal Precision Design System** defined in [DESIGN.md](file:///Users/ngdangdat/voc/dat-personal-stack/DESIGN.md):

* **Surface**: `#051424` (Deep navy slate)
* **Primary Accent**: `#22d3ee` (High-visibility terminal Cyan)
* **Typography**: Monospace JetBrains Mono mixed with clean Geist sans-serif.
* **Vibe**: Flat, high-density, sharp 4px borders, zero soft shadows.

---

## 📂 Repository Structure

* [DESIGN.md](file:///Users/ngdangdat/voc/dat-personal-stack/DESIGN.md): Specifications for typography, colors, layout structures, and components.
* [AGENTS.md](file:///Users/ngdangdat/voc/dat-personal-stack/AGENTS.md): The core application specification and developer agent guidelines.
* [specs/README.md](file:///Users/ngdangdat/voc/dat-personal-stack/specs/README.md): Deep-dive technical specifications including database schemas, JSON-RPC structures, and API patterns.

---

## ⚙️ Local Development & Operations

### Prerequisites
* Docker & Docker Compose
* Go 1.23+ (for local backend testing/running)
* Node.js v20+ (for local frontend testing/running)

### Local Environment Setup

#### 1. Setup Environment Configuration
Before starting, create your local configuration by copying [.env.example](file:///Users/ngdangdat/voc/dat-personal-stack/.env.example):
```bash
cp .env.example .env
```
*(The default values in `.env` are suitable for local development).*

#### 2. Configure GitHub API Access (Personal Access Token)
For security and privacy, GitHub credentials are **not** stored in backend environment variables. Instead, they are configured per-user on the client-side:
1. **Generate a GitHub PAT**:
   * Go to GitHub: **Settings** > **Developer Settings** > **Personal Access Tokens** > **Tokens (classic)**.
   * Click **Generate new token (classic)** and select the `repo` scope.
2. **Save credentials in the App**:
   * Open the Frontend Dashboard (e.g., `http://localhost:3000`).
   * Click the **Settings** tab on the navigation bar.
   * Input your GitHub **Username** and the generated **Personal Access Token (PAT)**, then click **SAVE**.
   * *These are stored securely on your local device inside IndexedDB and automatically attached to backend API requests.*

#### 3. Running via Docker Compose (Recommended)
Build and start the entire stack (Frontend, Backend, and Custom Node Agent) locally:
```bash
docker compose up --build
```
* **Frontend Dashboard**: [http://localhost:3000](http://localhost:3000)
* **Main API Server**: [http://localhost:8080](http://localhost:8080)
* **Telemetry & RPC Agent**: [ws://localhost:8081/rpc](ws://localhost:8081/rpc)

#### 4. Running Separately on Bare Metal (For live reloading / debugging)
* **Frontend (Vite dev server)**:
  ```bash
  cd frontend
  npm install
  npm run dev
  ```
  *(Default server: [http://localhost:5173](http://localhost:5173))*

* **Backend Server**:
  ```bash
  cd backend
  go run cmd/server/main.go
  ```
  *(Default server: [http://localhost:8080](http://localhost:8080))*

* **Node Agent Daemon**:
  ```bash
  cd backend
  PORT=8081 AGENT_SECRET_TOKEN=agent-secret-token go run cmd/node_agent/main.go
  ```
  *(Default WebSocket: `ws://localhost:8081/rpc`)*

---

## 🌐 Production Deployment

For deployment in staging or production environments, follow these security and orchestration practices:

### 1. Configure Production Secrets
Create a secure `.env` file on your server. Do **NOT** use default values for secrets in production:
```env
PORT=8081
AGENT_SECRET_TOKEN=generate-a-long-secure-random-token-here
```

### 2. Build Containerized Images
Dockerfiles are structured with **multi-stage builds** to minimize size and secure production runs:
* **Backend**: Compiles a minimal static binary inside `alpine` running under a non-root `appuser`.
* **Frontend**: Compiles Vite assets, then serves them using Nginx.

Deploy via Docker Compose in daemon mode:
```bash
docker compose -f docker-compose.yml up -d
```

### 3. Deploying the Node Agent on a Remote Host (e.g. Raspberry Pi)
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

### 4. Reverse Proxy & SSL (HTTPS/WSS)
It is highly recommended to place a reverse proxy (such as Nginx, Traefik, or Caddy) in front of the ports to handle SSL/TLS termination:
* Route public web traffic to the frontend Nginx server (port `3000`).
* Proxy `/api` to the backend server (port `8080`).
* Proxy `/rpc` to the node agent WebSocket endpoint (port `8081`), ensuring WebSocket upgrade headers (`Upgrade` and `Connection`) are properly forwarded.

---

## 🧪 Testing & Linting

### Backend Tests:
```bash
cd backend
go test ./...
```

### Frontend Tests & Lint:
```bash
cd frontend
npm install
npm test          # Run Vitest suite
npm run lint      # Run ESLint validation
```

---
*Developed with ❤️ as part of the Google Stitch developer ecosystem.*
