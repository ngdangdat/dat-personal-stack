<template>
  <div class="workspace-manager">
    <!-- Add Node Panel -->
    <div class="panel form-panel">
      <div class="panel-header">
        <span class="panel-title">// ADD_NODE.cfg</span>
      </div>
      <div class="panel-body monospaced">
        <form @submit.prevent="addNode" class="node-form">
          <div class="form-group">
            <label for="node-name">NODE_NAME:</label>
            <input 
              id="node-name" 
              v-model="newNode.name" 
              type="text" 
              placeholder="e.g. Local-Pi" 
              required 
              class="form-input" 
            />
          </div>
          <div class="form-group">
            <label for="node-address">HOST_ADDRESS:</label>
            <input 
              id="node-address" 
              v-model="newNode.address" 
              type="text" 
              placeholder="e.g. localhost:8081" 
              required 
              class="form-input" 
            />
          </div>
          <div class="form-group">
            <label for="node-token">AUTH_TOKEN:</label>
            <input 
              id="node-token" 
              v-model="newNode.token" 
              type="password" 
              placeholder="e.g. agent-secret-token" 
              required 
              class="form-input" 
            />
          </div>
          <div class="form-group-row">
            <div class="form-group">
              <label for="node-build-cmd">BUILD_COMMAND:</label>
              <input 
                id="node-build-cmd" 
                v-model="newNode.build_cmd" 
                type="text" 
                placeholder="e.g. npm run build" 
                class="form-input" 
              />
            </div>
            <div class="form-group">
              <label for="node-deploy-cmd">DEPLOY_COMMAND:</label>
              <input 
                id="node-deploy-cmd" 
                v-model="newNode.deploy_cmd" 
                type="text" 
                placeholder="e.g. npm run preview" 
                class="form-input" 
              />
            </div>
          </div>
          <button type="submit" class="btn btn-primary">ADD_NODE</button>
        </form>
      </div>
    </div>

    <!-- Nodes List Panel -->
    <div class="panel nodes-panel">
      <div class="panel-header">
        <span class="panel-title">// ACTIVE_NODES.list</span>
      </div>
      <div class="panel-body monospaced">
        <div v-if="nodes.length === 0" class="dimmed-text text-center pad-y">
          > No active nodes configured. Please add one above.
        </div>
        <div v-else class="nodes-list">
          <div 
            v-for="node in nodes" 
            :key="node.id" 
            :class="['node-item', { active: activeNode && activeNode.id === node.id }]"
          >
            <div class="node-info">
              <span class="node-label">{{ node.name }}</span>
              <span class="node-addr dimmed-text">{{ node.address }}</span>
            </div>
            
            <div class="node-status-area">
              <span :class="['status-badge', getNodeStatus(node.id)]">
                {{ getNodeStatus(node.id).toUpperCase() }}
              </span>
            </div>

            <div class="node-actions">
              <button 
                v-if="getNodeStatus(node.id) === 'connected'" 
                @click="disconnectNode()" 
                class="btn btn-error btn-sm"
              >
                DISCONNECT
              </button>
              <button 
                v-else-if="getNodeStatus(node.id) === 'connecting'" 
                disabled 
                class="btn btn-dimmed btn-sm"
              >
                CONNECTING...
              </button>
              <button 
                v-else 
                @click="connectNode(node)" 
                class="btn btn-primary btn-sm"
              >
                CONNECT
              </button>
              <button @click="deleteNode(node.id)" class="btn btn-dimmed btn-sm btn-delete">
                DELETE
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Connected Node Tabs -->
    <div v-if="activeNode && getNodeStatus(activeNode.id) === 'connected'" class="panel tabs-panel">
      <div class="panel-body monospaced no-pad-y">
        <div class="sub-tabs">
          <button 
            :class="['sub-tab-btn', { active: activeSubTab === 'telemetry' }]" 
            @click="activeSubTab = 'telemetry'"
          >
            [ TELEMETRY_INFO ]
          </button>
          <button 
            :class="['sub-tab-btn', { active: activeSubTab === 'project_terminal' }]" 
            @click="activeSubTab = 'project_terminal'"
          >
            [ PROJECT_TERMINAL ]
          </button>
        </div>
      </div>
    </div>

    <!-- Telemetry Panel -->
    <div v-if="activeNode && telemetry && activeSubTab === 'telemetry'" class="panel telemetry-panel">
      <div class="panel-header">
        <span class="panel-title">// TELEMETRY: {{ activeNode.name.toUpperCase() }}</span>
        <span class="status-indicator online">MONITORING</span>
      </div>
      <div class="panel-body monospaced">
        <div class="telemetry-grid">
          <!-- CPU -->
          <div class="telemetry-card">
            <span class="card-label">CPU_USAGE</span>
            <div class="progress-bar-container">
              <div 
                class="progress-bar" 
                :style="{ width: `${telemetry.cpu?.usage_percent || 0}%` }"
              ></div>
            </div>
            <div class="telemetry-details">
              <span>{{ telemetry.cpu?.usage_percent?.toFixed(1) || '0.0' }}%</span>
              <span class="dimmed-text">{{ telemetry.cpu?.temp_celsius?.toFixed(1) || '0.0' }}°C</span>
            </div>
          </div>

          <!-- Memory -->
          <div class="telemetry-card">
            <span class="card-label">MEM_DISTRIBUTION</span>
            <div class="progress-bar-container">
              <div 
                class="progress-bar" 
                :style="{ width: `${telemetry.memory ? (telemetry.memory.used_bytes / telemetry.memory.total_bytes) * 100 : 0}%` }"
              ></div>
            </div>
            <div class="telemetry-details">
              <span>{{ formatBytes(telemetry.memory?.used_bytes) }}</span>
              <span class="dimmed-text">/ {{ formatBytes(telemetry.memory?.total_bytes) }}</span>
            </div>
          </div>

          <!-- Threads -->
          <div class="telemetry-card">
            <span class="card-label">ACTIVE_THREADS</span>
            <span class="large-val">{{ telemetry.threads || 0 }}</span>
            <span class="metric-sub text-dimmed">goroutines</span>
          </div>

          <!-- Uptime -->
          <div class="telemetry-card">
            <span class="card-label">SYSTEM_UPTIME</span>
            <span class="large-val">{{ formatUptime(telemetry.uptime_seconds) }}</span>
            <span class="metric-sub text-dimmed">since agent start</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Terminal log panel -->
    <div v-if="activeNode && getNodeStatus(activeNode.id) === 'connected' && activeSubTab === 'telemetry'" class="panel terminal-panel">
      <div class="panel-header">
        <span class="panel-title">// TERMINAL_CONSOLE: {{ activeNode.name.toUpperCase() }}</span>
        <button @click="clearConsole" class="btn btn-dimmed btn-sm">CLEAR</button>
      </div>
      <div class="panel-body monospaced">
        <div ref="terminalLog" class="terminal-log-box">
          <div v-for="(line, idx) in consoleLines" :key="idx" :class="['log-line', line.type]">
            <span class="prompt" v-if="line.type === 'input'">$ </span>
            <span class="data">{{ line.text }}</span>
          </div>
          <div v-if="commandRunning" class="log-line system blink-cursor">
            <span class="prompt">></span> RUNNING COMMAND...
          </div>
        </div>
        
        <form @submit.prevent="runCommand" class="terminal-input-form">
          <span class="prompt">$</span>
          <input 
            v-model="currentCommand" 
            type="text" 
            placeholder="Type command e.g. ls -la..." 
            :disabled="commandRunning"
            class="terminal-input"
          />
          <button type="submit" :disabled="commandRunning || !currentCommand" class="btn btn-primary btn-sm">RUN</button>
        </form>
      </div>
    </div>

    <!-- Project Terminal Panel (Phase 3) -->
    <div v-if="activeNode && getNodeStatus(activeNode.id) === 'connected' && activeSubTab === 'project_terminal'">
      <ProjectTerminal 
        :gitStatus="gitStatus"
        :consoleLines="consoleLines"
        :actionRunning="commandRunning"
        :buildStatus="buildStatus"
        :deployStatus="deployStatus"
        @action="handleProjectAction"
        @clear-console="clearConsole"
        @run-custom="runCustomProjectCommand"
      />
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from 'vue';
import { db } from '../services/db';
import ProjectTerminal from './ProjectTerminal.vue';

export default {
  name: 'WorkspaceManager',
  components: {
    ProjectTerminal
  },
  setup() {
    const nodes = ref([]);
    const newNode = reactive({
      name: '',
      address: '',
      token: 'agent-secret-token',
      build_cmd: 'npm run build',
      deploy_cmd: 'npm run preview'
    });

    const activeNode = ref(null);
    const telemetry = ref(null);
    const connectionStates = reactive({}); // node.id -> 'disconnected' | 'connecting' | 'connected'
    
    // Sub-tab state
    const activeSubTab = ref('telemetry'); // 'telemetry' | 'project_terminal'
    const gitStatus = ref(null);
    const buildStatus = ref('unknown');
    const deployStatus = ref('unknown');
    const runningActionType = ref(null);

    // Command runner states
    const currentCommand = ref('');
    const commandRunning = ref(false);
    const consoleLines = ref([
      { type: 'system', text: 'Console initialized. Add a node and connect to execute commands.' }
    ]);
    const terminalLog = ref(null);
    
    let activeSocket = null;
    let commandIdCounter = 1;
    let pendingCommandId = null;

    const fetchGitStatus = () => {
      if (!activeSocket) return;
      const payload = {
        jsonrpc: '2.0',
        method: 'workspace.get_git_status',
        id: 'git_status'
      };
      try {
        activeSocket.send(JSON.stringify(payload));
      } catch (err) {
        console.error('Failed to fetch git status:', err);
      }
    };

    // Load nodes from DB
    const loadNodes = async () => {
      try {
        nodes.value = await db.getNodes();
      } catch (err) {
        console.error('Failed to load nodes:', err);
      }
    };

    const addNode = async () => {
      const id = 'node-' + Date.now();
      const nodeToSave = {
        id,
        name: newNode.name,
        address: newNode.address,
        token: newNode.token,
        build_cmd: newNode.build_cmd || 'npm run build',
        deploy_cmd: newNode.deploy_cmd || 'npm run preview'
      };

      try {
        await db.saveNode(nodeToSave);
        newNode.name = '';
        newNode.address = '';
        newNode.token = 'agent-secret-token';
        newNode.build_cmd = 'npm run build';
        newNode.deploy_cmd = 'npm run preview';
        await loadNodes();
      } catch (err) {
        console.error('Failed to save node:', err);
      }
    };

    const deleteNode = async (id) => {
      if (activeNode.value && activeNode.value.id === id) {
        disconnectNode();
      }
      try {
        await db.deleteNode(id);
        await loadNodes();
      } catch (err) {
        console.error('Failed to delete node:', err);
      }
    };

    const getNodeStatus = (id) => {
      return connectionStates[id] || 'disconnected';
    };

    const connectNode = (node) => {
      disconnectNode(); // ensure any existing socket is closed
      
      activeNode.value = node;
      connectionStates[node.id] = 'connecting';
      
      let wsUrl = node.address;
      if (!/^wss?:\/\//i.test(wsUrl)) {
        const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
        wsUrl = protocol + wsUrl;
      }
      
      if (!wsUrl.includes('/rpc')) {
        wsUrl = wsUrl.replace(/\/$/, '') + '/rpc';
      }

      appendSystemLog(`Connecting to ${node.name} (${node.address})...`);

      try {
        // Securely pass token via WebSocket subprotocol rather than URL parameters
        const ws = new WebSocket(wsUrl, node.token);
        activeSocket = ws;

        ws.onopen = () => {
          if (ws !== activeSocket) return;
          connectionStates[node.id] = 'connected';
          appendSystemLog(`Successfully connected to ${node.name}. Streaming live telemetry...`);
          fetchGitStatus();
        };

        ws.onmessage = (event) => {
          if (ws !== activeSocket) return;
          try {
            const data = JSON.parse(event.data);
            
            // Check for notifications
            if (data.method) {
              if (data.method === 'system.telemetry_stream') {
                telemetry.value = data.params;
              } else if (data.method === 'workspace.stdout') {
                appendConsoleLine('stdout', data.params.data);
              } else if (data.method === 'workspace.stderr') {
                appendConsoleLine('stderr', data.params.data);
              }
            } 
            
            // Check for git status responses
            if (data.id === 'git_status') {
              if (data.result) {
                gitStatus.value = data.result;
              }
              return;
            }

            // Check for responses
            if (data.id !== undefined && data.id === pendingCommandId) {
              if (data.error) {
                appendConsoleLine('stderr', `RPC Error: ${data.error.message}\n`);
                if (runningActionType.value === 'build') buildStatus.value = 'failed';
                if (runningActionType.value === 'deploy') deployStatus.value = 'failed';
              } else if (data.result) {
                const exitCode = data.result.exit_code;
                appendSystemLog(`Process finished with exit code ${exitCode}`);
                
                if (runningActionType.value === 'build') {
                  buildStatus.value = exitCode === 0 ? 'success' : 'failed';
                } else if (runningActionType.value === 'deploy') {
                  deployStatus.value = exitCode === 0 ? 'success' : 'failed';
                } else if (runningActionType.value === 'pull') {
                  fetchGitStatus();
                }
              }
              commandRunning.value = false;
              pendingCommandId = null;
              runningActionType.value = null;
            }
          } catch (e) {
            console.error('Failed to parse WebSocket message:', e);
          }
        };

        ws.onclose = () => {
          if (ws !== activeSocket) return;
          connectionStates[node.id] = 'disconnected';
          if (activeNode.value && activeNode.value.id === node.id) {
            appendSystemLog(`Connection to ${node.name} closed.`);
            activeNode.value = null;
            telemetry.value = null;
            commandRunning.value = false;
            gitStatus.value = null;
            buildStatus.value = 'unknown';
            deployStatus.value = 'unknown';
            runningActionType.value = null;
          }
        };

        ws.onerror = (err) => {
          if (ws !== activeSocket) return;
          console.error('WebSocket error:', err);
          appendSystemLog(`Connection error on ${node.name}.`);
        };
      } catch (err) {
        console.error('WebSocket connection initialization failed:', err);
        connectionStates[node.id] = 'disconnected';
        appendSystemLog(`Failed to initialize connection.`);
      }
    };

    const disconnectNode = () => {
      if (activeSocket) {
        activeSocket.close();
        activeSocket = null;
      }
      if (activeNode.value) {
        connectionStates[activeNode.value.id] = 'disconnected';
        activeNode.value = null;
        telemetry.value = null;
        gitStatus.value = null;
        buildStatus.value = 'unknown';
        deployStatus.value = 'unknown';
        runningActionType.value = null;
      }
    };

    const runCommand = () => {
      if (!activeSocket || !currentCommand.value || commandRunning.value) return;

      const cmd = currentCommand.value;
      const cmdId = commandIdCounter++;
      pendingCommandId = cmdId;
      commandRunning.value = true;

      // Add to console log
      consoleLines.value.push({ type: 'input', text: cmd });
      scrollConsole();

      const payload = {
        jsonrpc: '2.0',
        method: 'workspace.run_command',
        params: {
          workspace_id: 'local',
          command: cmd
        },
        id: cmdId
      };

      try {
        activeSocket.send(JSON.stringify(payload));
        currentCommand.value = '';
      } catch (err) {
        console.error('Failed to send command:', err);
        appendConsoleLine('stderr', `Failed to send command: ${err.message}\n`);
        commandRunning.value = false;
        pendingCommandId = null;
      }
    };

    const appendConsoleLine = (type, text) => {
      // Chunk-based streams might not end with newline, clean up representation
      consoleLines.value.push({ type, text });
      if (consoleLines.value.length > 1000) {
        consoleLines.value.shift();
      }
      scrollConsole();
    };

    const appendSystemLog = (text) => {
      consoleLines.value.push({ type: 'system', text: `> ${text}` });
      if (consoleLines.value.length > 1000) {
        consoleLines.value.shift();
      }
      scrollConsole();
    };

    const clearConsole = () => {
      consoleLines.value = [{ type: 'system', text: 'Console logs cleared.' }];
    };

    const scrollConsole = () => {
      nextTick(() => {
        if (terminalLog.value) {
          terminalLog.value.scrollTop = terminalLog.value.scrollHeight;
        }
      });
    };

    // Helper formatting functions
    const formatBytes = (bytes) => {
      if (!bytes) return '0 B';
      const k = 1024;
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    const formatUptime = (seconds) => {
      if (seconds === undefined || seconds === null) return 'N/A';
      const h = Math.floor(seconds / 3600);
      const m = Math.floor((seconds % 3600) / 60);
      const s = Math.floor(seconds % 60);
      return `${h}h ${m}m ${s}s`;
    };

    const handleProjectAction = (action) => {
      if (!activeSocket || commandRunning.value) return;

      let cmd = '';
      if (action === 'pull') {
        const branch = gitStatus.value?.branch;
        if (!branch || branch === 'unknown') {
          appendConsoleLine('stderr', 'Error: Active branch is unknown or cannot be determined. Pull aborted.\n');
          return;
        }
        cmd = `git pull origin ${branch}`;
        runningActionType.value = 'pull';
      } else if (action === 'build') {
        cmd = activeNode.value?.build_cmd || 'npm run build';
        buildStatus.value = 'building';
        runningActionType.value = 'build';
      } else if (action === 'deploy') {
        cmd = activeNode.value?.deploy_cmd || 'npm run preview';
        deployStatus.value = 'deploying';
        runningActionType.value = 'deploy';
      }

      executeProjectCommand(cmd);
    };

    const runCustomProjectCommand = (cmd) => {
      runningActionType.value = 'custom';
      executeProjectCommand(cmd);
    };

    const executeProjectCommand = (cmd) => {
      const cmdId = commandIdCounter++;
      pendingCommandId = cmdId;
      commandRunning.value = true;

      // Add to console log
      consoleLines.value.push({ type: 'input', text: cmd });
      if (consoleLines.value.length > 1000) {
        consoleLines.value.shift();
      }
      scrollConsole();

      const payload = {
        jsonrpc: '2.0',
        method: 'workspace.run_command',
        params: {
          workspace_id: 'local',
          command: cmd
        },
        id: cmdId
      };

      try {
        activeSocket.send(JSON.stringify(payload));
      } catch (err) {
        console.error('Failed to send project command:', err);
        appendConsoleLine('stderr', `Failed to send command: ${err.message}\n`);
        commandRunning.value = false;
        pendingCommandId = null;
        runningActionType.value = null;
        if (buildStatus.value === 'building') buildStatus.value = 'failed';
        if (deployStatus.value === 'deploying') deployStatus.value = 'failed';
      }
    };

    onMounted(() => {
      loadNodes();
    });

    onBeforeUnmount(() => {
      disconnectNode();
    });

    return {
      nodes,
      newNode,
      activeNode,
      telemetry,
      connectionStates,
      currentCommand,
      commandRunning,
      consoleLines,
      terminalLog,
      addNode,
      deleteNode,
      getNodeStatus,
      connectNode,
      disconnectNode,
      runCommand,
      clearConsole,
      formatBytes,
      formatUptime,
      activeSubTab,
      gitStatus,
      buildStatus,
      deployStatus,
      handleProjectAction,
      runCustomProjectCommand
    };
  }
};
</script>

<style scoped>
.workspace-manager {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Form Styles */
.node-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group label {
  color: var(--color-primary);
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.form-input {
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-variant);
  color: var(--color-on-surface);
  padding: 8px 10px;
  border-radius: var(--roundness-border);
  font-family: var(--font-mono);
  font-size: 0.85rem;
  outline: none;
  transition: border-color 0.15s ease;
}

.form-input:focus {
  border-color: var(--color-primary);
}

.btn {
  background-color: var(--color-outline);
  border: 1px solid var(--color-outline);
  color: var(--color-on-surface);
  padding: 8px 16px;
  border-radius: var(--roundness-border);
  font-family: var(--font-mono);
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  outline: none;
}

.btn:active {
  transform: scale(0.95);
}

.btn-primary {
  background-color: var(--color-primary-bg);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.btn-primary:hover {
  background-color: var(--color-primary);
  color: var(--color-surface);
}

.btn-error {
  background-color: rgba(248, 113, 113, 0.1);
  border-color: var(--color-error);
  color: var(--color-error);
}

.btn-error:hover {
  background-color: var(--color-error);
  color: var(--color-on-surface);
}

.btn-dimmed {
  background-color: transparent;
  border-color: var(--color-outline-variant);
  color: var(--color-on-surface-variant);
}

.btn-dimmed:hover {
  border-color: var(--color-outline);
  color: var(--color-on-surface);
}

.btn-sm {
  padding: 4px 8px;
  font-size: 0.7rem;
}

/* Active Nodes List */
.nodes-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.node-item {
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-variant);
  padding: 10px;
  border-radius: var(--roundness-border);
  display: grid;
  grid-template-columns: 2fr 1fr auto;
  align-items: center;
  gap: 12px;
}

.node-item.active {
  border-color: var(--color-primary);
  box-shadow: 0 0 4px rgba(34, 211, 238, 0.1);
}

.node-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.node-label {
  font-weight: 700;
  font-size: 0.85rem;
  color: var(--color-on-surface);
}

.node-addr {
  font-size: 0.7rem;
}

.status-badge {
  font-size: 0.65rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 2px;
  border: 1px solid transparent;
  display: inline-block;
}

.status-badge.disconnected {
  background-color: rgba(100, 116, 139, 0.1);
  border-color: var(--color-draft);
  color: var(--color-draft);
}

.status-badge.connecting {
  background-color: rgba(251, 191, 36, 0.1);
  border-color: var(--color-warning);
  color: var(--color-warning);
}

.status-badge.connected {
  background-color: rgba(34, 211, 238, 0.1);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.node-actions {
  display: flex;
  gap: 6px;
}

/* Telemetry styles */
.telemetry-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.telemetry-card {
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-variant);
  border-radius: var(--roundness-border);
  padding: 12px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 8px;
  min-height: 80px;
}

.card-label {
  font-size: 0.65rem;
  color: var(--color-on-surface-variant);
  font-weight: 700;
}

.large-val {
  font-size: 1.25rem;
  color: var(--color-primary);
  font-weight: 700;
}

.progress-bar-container {
  background-color: var(--color-outline-variant);
  height: 6px;
  width: 100%;
  border-radius: 3px;
  overflow: hidden;
}

.progress-bar {
  background-color: var(--color-primary);
  height: 100%;
  border-radius: 3px;
  transition: width 0.3s ease;
}

.telemetry-details {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--color-on-surface);
}

.metric-sub {
  font-size: 0.6rem;
}

.status-indicator.online {
  background-color: rgba(34, 211, 238, 0.1);
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
}

/* Terminal Styles */
.terminal-log-box {
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-variant);
  padding: 12px;
  border-radius: var(--roundness-border);
  height: 180px;
  overflow-y: auto;
  font-family: var(--font-mono);
  font-size: 0.8rem;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.log-line {
  line-height: 1.4;
  white-space: pre-wrap;
}

.log-line.system {
  color: var(--color-on-surface-variant);
}

.log-line.input {
  color: var(--color-primary);
}

.log-line.stdout {
  color: var(--color-on-surface);
}

.log-line.stderr {
  color: var(--color-error);
}

.terminal-input-form {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-variant);
  padding: 6px 10px;
  border-radius: var(--roundness-border);
  margin-top: 8px;
}

.terminal-input-form .prompt {
  color: var(--color-primary);
  font-weight: 700;
  user-select: none;
}

.terminal-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--color-on-surface);
  font-family: var(--font-mono);
  font-size: 0.8rem;
}

.blink-cursor::after {
  content: '▋';
  color: var(--color-primary);
  animation: blink 1s step-start infinite;
}

@keyframes blink {
  50% { opacity: 0; }
}

.pad-y {
  padding: 16px 0;
}

.text-center {
  text-align: center;
}

/* Connected Node Tabs */
.tabs-panel {
  margin-bottom: 12px;
}

.sub-tabs {
  display: flex;
  gap: 16px;
  padding: 8px 0;
}

.sub-tab-btn {
  background: transparent;
  border: none;
  color: var(--color-on-surface-variant);
  font-family: var(--font-mono);
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  padding: 4px 8px;
  transition: color 0.15s ease;
  outline: none;
}

.sub-tab-btn.active {
  color: var(--color-primary);
}

.sub-tab-btn:hover:not(.active) {
  color: var(--color-on-surface);
}

.form-group-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

@media (max-width: 480px) {
  .form-group-row {
    grid-template-columns: 1fr;
  }
}
</style>
