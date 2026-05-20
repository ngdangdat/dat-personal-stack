# Engineering Assistant: Mobile Command Center

A mobile-first, high-density dashboard and remote orchestrator built for modern software engineers. It bridges the gap between high-level project management (GitHub PRs) and low-level execution environments (Remote Raspberry Pi, cloud instances, and local micro-environments).

---

## 🚀 Core Features

### 1. PR Dashboard (Phase 1)
* **Direct Integration**: OAuth authentication with GitHub REST/GraphQL APIs.
* **Smart Filter Views**: High-density lists for "Reviewing" (assigned PRs) and "My PRs" (authored PRs).
* **Review Velocity Metrics**: Clear indicator badges showing lead times and status indicators (Draft, Approved, Changes Requested).

### 2. Workspace Manager & RPC Orchestration (Phase 2)
* **Node Telemetry**: Live telemetry metrics (CPU load, memory distribution, threads, uptime) via WebSocket.
* **Isolated Environments**: Support for Node.js, Go, Rust, and Python runtime instances on target nodes.
* **Remote Commands**: Initiate repository actions (pull, rebuild, restart service) directly from the mobile UI.

### 3. Project Detail & Terminal (Phase 3)
* **Integrated TUI Console**: Dense repository tracking displaying active branches, last commits, and live shell logs.
* **Interactive Shortcuts**: Action buttons for instant command invocation (`git pull`, `build`, `deploy`).

### 4. Assistant & Local Storage (Phase 4)
* **AI Chat Companion**: An on-device assistant designed to track project status and help troubleshoot failures.
* **Offline-First Storage**: Local persistent caching utilizing IndexedDB for instantaneous performance.

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

## ⚙️ Local Development

### Prerequisites
* Node.js (v18+)
* Git

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/ngdangdat/dat-personal-stack.git
   cd dat-personal-stack
   ```
2. Set up environments and configure dependencies (to be updated as source code templates are added).

---
*Developed with ❤️ as part of the Google Stitch developer ecosystem.*
