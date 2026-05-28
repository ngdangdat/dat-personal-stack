<template>
  <nav class="bottom-nav-bar">
    <button 
      v-for="item in navItems" 
      :key="item.id" 
      :id="`nav-btn-${item.id}`"
      :class="['nav-item', { active: activeTab === item.id }]" 
      @click="$emit('changeTab', item.id)"
    >
      <div class="icon-wrapper" v-html="item.icon"></div>
      <span class="nav-label">{{ item.label }}</span>
    </button>
  </nav>
</template>

<script>
export default {
  name: 'BottomNavBar',
  props: {
    activeTab: {
      type: String,
      required: true
    }
  },
  setup() {
    const navItems = [
      {
        id: 'prs',
        label: 'PRs',
        // Git Pull Request/Merge icon
        icon: `<svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="18" cy="18" r="3"></circle>
          <circle cx="6" cy="6" r="3"></circle>
          <circle cx="6" cy="18" r="3"></circle>
          <path d="M18 15V9a4 4 0 0 0-4-4H9"></path>
          <line x1="6" y1="9" x2="6" y2="15"></line>
        </svg>`
      },
      {
        id: 'workspaces',
        label: 'Nodes',
        // Grid/Workspace icon
        icon: `<svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="7" height="9"></rect>
          <rect x="14" y="3" width="7" height="5"></rect>
          <rect x="14" y="12" width="7" height="9"></rect>
          <rect x="3" y="16" width="7" height="5"></rect>
        </svg>`
      },
      {
        id: 'chat',
        label: 'Chat',
        // Message/Forum icon
        icon: `<svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>`
      },
      {
        id: 'settings',
        label: 'Settings',
        // Settings gear icon
        icon: `<svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>`
      }
    ];

    return {
      navItems
    };
  }
};
</script>

<style scoped>
.bottom-nav-bar {
  background-color: var(--color-surface-container);
  border-top: 1px solid var(--color-outline-variant);
  height: 56px;
  width: 100%;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  position: sticky;
  bottom: 0;
  z-index: 100;
}

.nav-item {
  background: transparent;
  border: none;
  color: var(--color-on-surface-variant);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  outline: none;
  transition: transform 0.1s ease, background-color 0.15s ease;
  height: 100%;
}

.nav-item:active {
  transform: scale(0.95); /* Scale-95 transition on active tap */
}

.nav-item.active {
  background-color: var(--color-secondary-container);
  color: var(--color-primary);
}

.icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-label {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  font-weight: 500;
}
</style>
