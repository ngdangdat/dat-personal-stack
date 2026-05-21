<template>
  <div class="app-container">
    <!-- Top App Bar -->
    <TopAppBar :isConnected="state.isConnected" />

    <!-- Main View Area -->
    <main class="main-content">
      <div v-if="currentTab === 'prs'">
        <PRDashboard :state="state" />
      </div>

      <div v-else-if="currentTab === 'workspaces'">
        <WorkspaceManager />
      </div>

      <div v-else-if="currentTab === 'chat'" class="panel placeholder-panel">
        <div class="panel-header">
          <span class="panel-title">// CO-PILOT_AGENT.log</span>
          <span class="status-indicator warning">PHASE_4</span>
        </div>
        <div class="panel-body monospaced">
          <p class="accent-text">> Activating AI model pipeline...</p>
          <p class="dimmed-text">Persistent troubleshooting chat companion and log database synchronization are scheduled for Phase 4 development.</p>
          <div class="terminal-box">
            <span class="prompt">$</span> antigravity --assist "Check server build logs"<br>
            <span class="prompt">Error:</span> Agent engine offline (Status 404).
          </div>
        </div>
      </div>

      <div v-else-if="currentTab === 'settings'">
        <Settings :state="state" @save="saveSettings" />
      </div>
    </main>

    <!-- Bottom Navigation Bar -->
    <BottomNavBar :activeTab="currentTab" @changeTab="setTab" />
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue';
import TopAppBar from './components/TopAppBar.vue';
import BottomNavBar from './components/BottomNavBar.vue';
import PRDashboard from './components/PRDashboard.vue';
import Settings from './components/Settings.vue';
import WorkspaceManager from './components/WorkspaceManager.vue';
import { db } from './services/db';

export default {
  name: 'App',
  components: {
    TopAppBar,
    BottomNavBar,
    PRDashboard,
    Settings,
    WorkspaceManager
  },
  setup() {
    const currentTab = ref('prs');
    const state = reactive({
      token: '',
      username: '',
      isConnected: false,
      isChecking: false
    });

    const setTab = (tab) => {
      currentTab.value = tab;
    };

    // Load configuration on mount (sync backend PostgreSQL & fallback to IndexedDB)
    onMounted(async () => {
      let loadedFromBackend = false;
      try {
        const response = await fetch('/api/settings');
        if (response.ok) {
          const config = await response.json();
          if (config && config.username && config.token) {
            state.token = config.token;
            state.username = config.username;
            loadedFromBackend = true;
            checkConnection();
          }
        }
      } catch (err) {
        console.error("Failed to load config from backend settings:", err);
      }

      if (!loadedFromBackend) {
        try {
          const config = await db.getConfig();
          if (config) {
            state.token = config.token || '';
            state.username = config.username || '';
            if (state.token && state.username) {
              // Sync local IndexedDB config to PostgreSQL database
              try {
                await fetch('/api/settings/save', {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ username: state.username, token: state.token })
                });
              } catch (e) {
                console.error("Auto-sync of local config to backend failed:", e);
              }
              checkConnection();
            }
          }
        } catch (err) {
          console.error("Failed to load config from db:", err);
        }
      }
    });

    // Check backend health/GitHub API connectivity
    const checkConnection = async () => {
      if (!state.token || !state.username) {
        state.isConnected = false;
        return;
      }

      state.isChecking = true;
      try {
        const response = await fetch(`/api/github/prs?type=reviewing&username=${encodeURIComponent(state.username)}`, {
          headers: {
            'Authorization': `Bearer ${state.token}`
          }
        });
        
        state.isConnected = response.ok;
      } catch (err) {
        console.error("Connection check failed:", err);
        state.isConnected = false;
      } finally {
        state.isChecking = false;
      }
    };

    // Save configurations
    const saveSettings = async ({ token, username }) => {
      state.token = token;
      state.username = username;
      
      // Save locally to IndexedDB first
      await db.saveConfig({ token, username });
      
      // Sync to PostgreSQL backend
      try {
        await fetch('/api/settings/save', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ token, username })
        });
      } catch (err) {
        console.error("Failed to save credentials to backend PostgreSQL:", err);
      }

      await checkConnection();
    };

    return {
      currentTab,
      state,
      setTab,
      saveSettings
    };
  }
};
</script>

<style>
/* Base Reset and Layout */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  background-color: var(--color-surface);
  color: var(--color-text-primary);
  font-family: var(--font-sans);
  overflow-x: hidden;
  height: 100vh;
  width: 100vw;
  -webkit-font-smoothing: antialiased;
}

.app-container {
  display: grid;
  grid-template-rows: auto 1fr auto;
  height: 100vh;
  width: 100%;
}

.main-content {
  overflow-y: auto;
  padding: 16px;
  width: 100%;
  max-width: 600px; /* Mobile first lock width */
  margin: 0 auto;
  padding-bottom: 30px;
}

/* Common panels and text utility classes */
.panel {
  background-color: var(--color-surface-container);
  border: 1px solid var(--color-outline);
  border-radius: var(--roundness-border);
  margin-bottom: 16px;
  overflow: hidden;
}

.panel-header {
  background-color: var(--color-surface-header);
  border-bottom: 1px solid var(--color-outline);
  padding: 10px 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-title {
  font-family: var(--font-mono);
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-primary);
  letter-spacing: 0.5px;
}

.panel-body {
  padding: 14px;
}

.placeholder-panel .panel-body p {
  margin-bottom: 12px;
  font-size: 0.9rem;
  line-height: 1.4;
}

.terminal-box {
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-dimmed);
  padding: 10px;
  border-radius: var(--roundness-border);
  font-family: var(--font-mono);
  font-size: 0.8rem;
  line-height: 1.5;
  color: var(--color-text-secondary);
  margin-top: 14px;
}

.terminal-box .prompt {
  color: var(--color-primary);
  user-select: none;
}

.terminal-box .prompt.warning {
  color: var(--color-warning);
}

.monospaced {
  font-family: var(--font-mono);
}

.accent-text {
  color: var(--color-primary) !important;
}

.dimmed-text {
  color: var(--color-text-secondary) !important;
}

.status-indicator {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  padding: 2px 6px;
  border-radius: 2px;
  font-weight: 700;
}

.status-indicator.warning {
  background-color: rgba(234, 179, 8, 0.1);
  color: var(--color-warning);
  border: 1px solid var(--color-warning);
}
</style>
