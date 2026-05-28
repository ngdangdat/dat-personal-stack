<template>
  <div class="project-terminal">
    <!-- Grid Summary Panel -->
    <div class="panel status-panel">
      <div class="panel-header">
        <span class="panel-title">// REPOSITORY_STATUS.cfg</span>
        <span v-if="gitStatus" :class="['status-indicator', gitStatus.is_clean ? 'clean' : 'dirty']">
          {{ gitStatus.is_clean ? 'CLEAN' : 'DIRTY' }}
        </span>
      </div>
      <div class="panel-body monospaced">
        <div class="status-grid">
          <div class="grid-item">
            <span class="grid-label">ACTIVE_BRANCH:</span>
            <span class="grid-value accent-text">{{ gitStatus?.branch || 'UNKNOWN' }}</span>
          </div>
          <div class="grid-item">
            <span class="grid-label">LAST_BUILD:</span>
            <span :class="['grid-value', buildStatusClass]">
              {{ buildStatus.toUpperCase() }}
            </span>
          </div>
          <div class="grid-item">
            <span class="grid-label">LAST_DEPLOY:</span>
            <span :class="['grid-value', deployStatusClass]">
              {{ deployStatus.toUpperCase() }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Commits List Panel -->
    <div class="panel commits-panel">
      <div class="panel-header">
        <span class="panel-title">// RECENT_COMMITS.log</span>
      </div>
      <div class="panel-body monospaced no-pad-y">
        <div v-if="!gitStatus || !gitStatus.commits || gitStatus.commits.length === 0" class="dimmed-text pad-y">
          > No commit logs available.
        </div>
        <div v-else class="commits-list">
          <div class="commit-header">
            <span>HASH</span>
            <span>AUTHOR</span>
            <span>DATE</span>
            <span>MESSAGE</span>
          </div>
          <div v-for="commit in gitStatus.commits" :key="commit.hash" class="commit-item">
            <span class="commit-hash accent-text">{{ commit.hash }}</span>
            <span class="commit-author dimmed-text truncate">{{ commit.author }}</span>
            <span class="commit-date dimmed-text truncate">{{ commit.date }}</span>
            <span class="commit-message truncate">{{ commit.message }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Direct Actions Panel -->
    <div class="panel actions-panel">
      <div class="panel-header">
        <span class="panel-title">// REPOSITORY_ACTIONS.sh</span>
      </div>
      <div class="panel-body monospaced">
        <div class="actions-grid">
          <button 
            @click="$emit('action', 'pull')" 
            :disabled="actionRunning" 
            class="btn btn-primary"
          >
            GIT_PULL
          </button>
          <button 
            @click="$emit('action', 'build')" 
            :disabled="actionRunning" 
            class="btn btn-primary"
          >
            RUN_BUILD
          </button>
          <button 
            @click="$emit('action', 'deploy')" 
            :disabled="actionRunning" 
            class="btn btn-warning-outline"
          >
            RUN_DEPLOY
          </button>
        </div>
      </div>
    </div>

    <!-- Terminal Log Console -->
    <div class="panel console-panel">
      <div class="panel-header">
        <span class="panel-title">// CONSOLE_LOGS: STREAM_STDOUT</span>
        <button @click="$emit('clear-console')" class="btn btn-dimmed btn-sm">CLEAR</button>
      </div>
      <div class="panel-body monospaced">
        <div ref="terminalLog" class="terminal-log-box">
          <div v-for="(line, idx) in consoleLines" :key="idx" :class="['log-line', line.type]">
            <span class="prompt" v-if="line.type === 'input'">$ </span>
            <span class="data">
              <span v-for="(chunk, cIdx) in parseAnsi(line.text)" :key="cIdx" :class="chunk.class">
                {{ chunk.text }}
              </span>
            </span>
          </div>
          <div v-if="actionRunning" class="log-line system blink-cursor">
            <span class="prompt">></span> EXECUTING PIPELINE...
          </div>
        </div>

        <!-- Custom Command Run Form -->
        <form @submit.prevent="submitCustomCommand" class="terminal-input-form">
          <span class="prompt">$</span>
          <input 
            v-model="customCommand" 
            type="text" 
            placeholder="Run custom command..." 
            :disabled="actionRunning"
            class="terminal-input"
          />
          <button 
            type="submit" 
            :disabled="actionRunning || !customCommand" 
            class="btn btn-primary btn-sm"
          >
            RUN
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, nextTick, computed } from 'vue';

export default {
  name: 'ProjectTerminal',
  props: {
    gitStatus: {
      type: Object,
      default: () => null
    },
    consoleLines: {
      type: Array,
      default: () => []
    },
    actionRunning: {
      type: Boolean,
      default: false
    },
    buildStatus: {
      type: String,
      default: 'unknown' // 'unknown' | 'building' | 'success' | 'failed'
    },
    deployStatus: {
      type: String,
      default: 'unknown' // 'unknown' | 'deploying' | 'success' | 'failed'
    }
  },
  emits: ['action', 'clear-console', 'run-custom'],
  setup(props, { emit }) {
    const customCommand = ref('');
    const terminalLog = ref(null);

    // Auto-scroll console
    const scrollConsole = () => {
      nextTick(() => {
        if (terminalLog.value) {
          terminalLog.value.scrollTop = terminalLog.value.scrollHeight;
        }
      });
    };

    watch(() => props.consoleLines.length, scrollConsole);
    watch(() => props.actionRunning, scrollConsole);

    const submitCustomCommand = () => {
      if (customCommand.value.trim() && !props.actionRunning) {
        emit('run-custom', customCommand.value);
        customCommand.value = '';
      }
    };

    // ANSI escape sequence color parser
    const parseAnsi = (text) => {
      if (!text) return [];
      
      const ansiRegex = /\u001b\[([0-9;]+)m/g;
      const parts = [];
      let lastIndex = 0;
      let currentClass = '';
      let match;
      
      const colorMap = {
        '0': '',
        '30': 'ansi-black',
        '31': 'ansi-red',
        '32': 'ansi-green',
        '33': 'ansi-yellow',
        '34': 'ansi-blue',
        '35': 'ansi-magenta',
        '36': 'ansi-cyan',
        '37': 'ansi-white',
        '90': 'ansi-gray',
      };

      while ((match = ansiRegex.exec(text)) !== null) {
        if (match.index > lastIndex) {
          parts.push({
            text: text.substring(lastIndex, match.index),
            class: currentClass
          });
        }
        
        const codes = match[1].split(';');
        for (const code of codes) {
          if (code === '0') {
            currentClass = '';
          } else if (colorMap[code]) {
            currentClass = colorMap[code];
          }
        }
        
        lastIndex = ansiRegex.lastIndex;
      }
      
      if (lastIndex < text.length) {
        parts.push({
          text: text.substring(lastIndex),
          class: currentClass
        });
      }
      
      return parts;
    };

    const buildStatusClass = computed(() => {
      if (props.buildStatus === 'success') return 'text-success';
      if (props.buildStatus === 'failed') return 'text-error';
      if (props.buildStatus === 'building') return 'text-warning';
      return 'dimmed-text';
    });

    const deployStatusClass = computed(() => {
      if (props.deployStatus === 'success') return 'text-success';
      if (props.deployStatus === 'failed') return 'text-error';
      if (props.deployStatus === 'deploying') return 'text-warning';
      return 'dimmed-text';
    });

    return {
      customCommand,
      terminalLog,
      submitCustomCommand,
      parseAnsi,
      buildStatusClass,
      deployStatusClass
    };
  }
};
</script>

<style scoped>
.project-terminal {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.grid-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.grid-label {
  font-size: 0.65rem;
  color: var(--color-on-surface-variant);
}

.grid-value {
  font-size: 0.85rem;
  font-weight: 700;
}

.status-indicator.clean {
  background-color: rgba(34, 211, 238, 0.1);
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
}

.status-indicator.dirty {
  background-color: rgba(251, 191, 36, 0.1);
  color: var(--color-warning);
  border: 1px solid var(--color-warning);
}

/* Commits List style */
.commits-list {
  display: flex;
  flex-direction: column;
  font-size: 0.75rem;
}

.commit-header {
  display: grid;
  grid-template-columns: 80px 100px 100px 1fr;
  padding: 6px 0;
  border-bottom: 1px solid var(--color-outline-variant);
  color: var(--color-on-surface-variant);
  font-weight: 700;
}

.commit-item {
  display: grid;
  grid-template-columns: 80px 100px 100px 1fr;
  padding: 8px 0;
  border-bottom: 1px solid var(--color-outline-variant);
  align-items: center;
}

.commit-item:last-child {
  border-bottom: none;
}

.truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.no-pad-y {
  padding-top: 0 !important;
  padding-bottom: 0 !important;
}

.pad-y {
  padding: 16px 0;
}

/* Actions Grid */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
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
  text-align: center;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn:active:not(:disabled) {
  transform: scale(0.95);
}

.btn-primary {
  background-color: var(--color-primary-bg);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--color-primary);
  color: var(--color-surface);
}

.btn-warning-outline {
  background-color: rgba(251, 191, 36, 0.05);
  border-color: var(--color-warning);
  color: var(--color-warning);
}

.btn-warning-outline:hover:not(:disabled) {
  background-color: var(--color-warning);
  color: var(--color-surface);
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

/* Terminal Feed */
.terminal-log-box {
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-variant);
  padding: 12px;
  border-radius: var(--roundness-border);
  height: 200px;
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

/* Status Classes */
.text-success {
  color: var(--color-success);
}

.text-error {
  color: var(--color-error);
}

.text-warning {
  color: var(--color-warning);
}

/* ANSI styled text colors */
:deep(.ansi-red) {
  color: var(--color-error);
}
:deep(.ansi-green) {
  color: var(--color-success);
}
:deep(.ansi-yellow) {
  color: var(--color-warning);
}
:deep(.ansi-blue) {
  color: #3b82f6;
}
:deep(.ansi-magenta) {
  color: #a855f7;
}
:deep(.ansi-cyan) {
  color: var(--color-primary);
}
:deep(.ansi-white) {
  color: var(--color-on-surface);
}
:deep(.ansi-gray) {
  color: var(--color-on-surface-variant);
}
</style>
