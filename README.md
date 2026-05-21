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
* Go 1.23+ (for local backend testing)
* Node.js v20+ (for local frontend testing)

### Quick Start with Docker Compose
To build and start the entire workspace center stack (Frontend, Backend, and Custom node agent) locally:
```bash
docker compose up --build
```
* **Frontend Dashboard**: [http://localhost:3000](http://localhost:3000)
* **Main API Server**: [http://localhost:8080](http://localhost:8080)
* **Telemetry & RPC Agent**: [ws://localhost:8081/rpc](ws://localhost:8081/rpc)

### Running Automated Test Suites & Linters

#### Backend Tests:
```bash
cd backend
go test ./...
```

#### Frontend Tests & Lint:
```bash
cd frontend
npm install
npm test          # Run Vitest suite
npm run lint      # Run ESLint validation
```

---
*Developed with ❤️ as part of the Google Stitch developer ecosystem.*
