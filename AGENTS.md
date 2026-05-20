# Engineering Assistant - Application Specification (AGENTS.md)

## Core Concept
A mobile-first engineering command center designed to bridge the gap between high-level project management (GitHub) and low-level system execution (Remote TUI/RPC). It serves as a unified interface for engineers to monitor contributions and manage remote environments from a single device.

---

## Application Architecture

### 1. PR Dashboard (Phase 1)
- **Source**: GitHub API integration.
- **Views**:
    - **Reviewing**: Pull requests assigned to the user for review.
    - **My PRs**: Pull requests opened by the user.
- **Key Metrics**: Review velocity, average lead time, and PR status (Draft, Approved, Changes Requested).

### 2. Workspace Manager (Phase 2)
- **Infrastructure**: Remote orchestration for connected nodes (Raspberry Pi, Cloud instances).
- **Communication**: Custom RPC protocol for real-time resource monitoring and command execution.
- **Telemetry**: CPU usage, memory distribution, uptime, and thread counts.
- **Workspaces**: Isolated environments per project with language-specific runtimes (Node.js, Python, Rust, Go).

### 3. Project Detail & Terminal
- **Interface**: Integrated TUI (Terminal User Interface) placeholder.
- **Actions**: Direct repository management (Pull, Build, Deploy).
- **Context**: Branch tracking, commit history, and active workspace status.

### 4. Assistant & Storage
- **Chat**: Persistent history for assistant-led project tracking and troubleshooting.
- **Logs**: Centralized data store for project events and engineering notes.
- **Data Persistence**: Local IndexedDB caching for performance with backend database synchronization.

---

## Design Language: Terminal Precision
The visual aesthetics must strictly align with the specifications defined in [DESIGN.md](file:///Users/ngdangdat/voc/dat-personal-stack/DESIGN.md).

- **Theme**: Dark mode primary (`#051424`).
- **Accent**: Cyan interactivity (`#22d3ee`).
- **Typography**: Geist / JetBrains Mono mix for high-density legibility.
- **Philosophy**: Minimalist, flat hierarchy, and status-at-a-glance visibility.

---

## Technical Stack
- **Frontend**: Responsive Mobile Web.
- **Integrations**: GitHub OAuth, Custom RPC.
- **Storage**: Real-time database for chat/logs.

---
*Source: Google Stitch*
