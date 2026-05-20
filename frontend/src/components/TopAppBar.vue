<template>
  <header class="top-app-bar">
    <div class="header-left">
      <svg class="terminal-icon" viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="4 17 10 11 4 5"></polyline>
        <line x1="12" y1="19" x2="20" y2="19"></line>
      </svg>
      <span class="app-title">ENGINEERING_ASSISTANT</span>
    </div>
    <div class="header-right">
      <span :class="['status-dot', isConnected ? 'connected' : 'offline']"></span>
      <span class="status-text">{{ isConnected ? 'CONNECTED' : 'OFFLINE' }}</span>
    </div>
  </header>
</template>

<script>
export default {
  name: 'TopAppBar',
  props: {
    isConnected: {
      type: Boolean,
      default: false
    }
  }
};
</script>

<style scoped>
.top-app-bar {
  background-color: #131313;
  border-bottom: 1px solid var(--color-outline);
  height: 48px;
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-primary);
}

.terminal-icon {
  stroke: var(--color-primary);
}

.app-title {
  font-family: var(--font-mono);
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 1px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-dot.connected {
  background-color: var(--color-approved);
  box-shadow: 0 0 6px var(--color-approved);
}

.status-dot.offline {
  background-color: var(--color-error);
}

.status-text {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  color: var(--color-text-secondary);
  font-weight: 600;
}
</style>
