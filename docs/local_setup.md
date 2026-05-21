# Local Environment Setup

This document outlines the steps required to set up and run the Engineering Assistant stack locally for development and testing.

## Prerequisites
Before starting, ensure you have the following installed:
* Docker & Docker Compose
* Go 1.23+ (for running the backend manually)
* Node.js v20+ (for running the frontend manually)

---

## 1. Setup Environment Configuration (.env)
Create your local configuration by copying [.env.example](../.env.example) in the project root:
```bash
cp .env.example .env
```
*(The default values in `.env` are suitable for local development).*

---

## 2. Configure GitHub API Access (Personal Access Token)
For security and privacy, GitHub credentials are **not** stored in backend environment variables. Instead, they are configured per-user on the client-side:
1. **Generate a GitHub PAT**:
   * Go to GitHub: **Settings** > **Developer Settings** > **Personal Access Tokens** > **Tokens (classic)**.
   * Click **Generate new token (classic)** and select the `repo` scope.
2. **Save credentials in the App**:
   * Open the Frontend Dashboard (e.g., `http://localhost:3000`).
   * Click the **Settings** tab on the navigation bar.
   * Input your GitHub **Username** and the generated **Personal Access Token (PAT)**, then click **SAVE**.
   * *These are stored securely on your local device inside IndexedDB and automatically attached to backend API requests.*

---

## 3. Running via Docker Compose (Recommended)
Build and start the entire stack (Frontend, Backend, and Custom Node Agent) locally:
```bash
docker compose up --build
```
* **Frontend Dashboard**: [http://localhost:3000](http://localhost:3000)
* **Main API Server**: [http://localhost:8080](http://localhost:8080)
* **Telemetry & RPC Agent**: [ws://localhost:8081/rpc](ws://localhost:8081/rpc)

---

## 4. Running Separately on Bare Metal (For live reloading / debugging)

### Frontend (Vite dev server)
```bash
cd frontend
npm install
npm run dev
```
*(Default dev server runs at [http://localhost:5173](http://localhost:5173))*

### Backend Server
```bash
cd backend
go run cmd/server/main.go
```
*(Default server runs at [http://localhost:8080](http://localhost:8080))*

### Node Agent Daemon
```bash
cd backend
PORT=8081 AGENT_SECRET_TOKEN=agent-secret-token go run cmd/node_agent/main.go
```
*(Default WebSocket runs at `ws://localhost:8081/rpc`)*

---

## 🧪 Testing & Linting

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests & Lint
```bash
cd frontend
npm install
npm test          # Run Vitest suite
npm run lint      # Run ESLint validation
```
