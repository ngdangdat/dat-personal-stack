<template>
  <div class="panel settings-panel">
    <div class="panel-header">
      <span class="panel-title">// SETTINGS.conf</span>
    </div>
    
    <div class="panel-body">
      <form @submit.prevent="handleSave">
        <!-- Username -->
        <div class="form-group">
          <label for="username">GitHub Username</label>
          <input 
            type="text" 
            id="username" 
            v-model="form.username" 
            placeholder="e.g. octocat" 
            required 
          />
        </div>

        <!-- Personal Access Token -->
        <div class="form-group">
          <label for="token">GitHub Personal Access Token (PAT)</label>
          <div class="input-password-wrapper">
            <input 
              :type="showToken ? 'text' : 'password'" 
              id="token" 
              v-model="form.token" 
              placeholder="ghp_xxxxxxxxxxxxxxxxxxxx" 
              required 
            />
            <button type="button" class="btn-toggle-visibility" @click="showToken = !showToken">
              {{ showToken ? 'HIDE' : 'SHOW' }}
            </button>
          </div>
          <span class="help-text">
            Required scopes: <code>repo</code> (for private/public PR metrics)
          </span>
        </div>

        <!-- Status & Connection Alert -->
        <div v-if="state.token && state.username" class="connection-status-box">
          <div class="status-row">
            <span class="label">API_STATUS:</span>
            <span :class="['value', state.isConnected ? 'success' : 'error']">
              {{ state.isChecking ? 'TESTING...' : (state.isConnected ? 'CONNECTED' : 'CONNECTION_FAILED') }}
            </span>
          </div>
        </div>

        <!-- Action buttons -->
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="state.isChecking">
            SAVE_CONFIGURATION
          </button>
        </div>
      </form>

      <!-- Help section -->
      <div class="terminal-box settings-help">
        <span class="prompt"># GitHub Auth Instructions</span><br>
        1. Go to github.com -> Settings -> Developer settings.<br>
        2. Personal access tokens -> Tokens (classic).<br>
        3. Generate a classic token with "repo" scope.<br>
        4. Paste the token above and click SAVE.
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, ref, watch } from 'vue';

export default {
  name: 'Settings',
  props: {
    state: {
      type: Object,
      required: true
    }
  },
  emits: ['save'],
  setup(props, { emit }) {
    const showToken = ref(false);
    const form = reactive({
      username: props.state.username,
      token: props.state.token
    });

    // Keep form in sync if loaded asynchronously
    watch(() => props.state.username, (newVal) => {
      form.username = newVal;
    });
    watch(() => props.state.token, (newVal) => {
      form.token = newVal;
    });

    const handleSave = () => {
      emit('save', {
        username: form.username.trim(),
        token: form.token.trim()
      });
    };

    return {
      form,
      showToken,
      handleSave
    };
  }
};
</script>

<style scoped>
.form-group {
  margin-bottom: 16px;
}

.input-password-wrapper {
  position: relative;
  display: flex;
}

.input-password-wrapper input {
  padding-right: 65px;
}

.btn-toggle-visibility {
  position: absolute;
  right: 1px;
  top: 1px;
  bottom: 1px;
  background: transparent;
  border: none;
  border-left: 1px solid var(--color-outline-variant);
  color: var(--color-on-surface-variant);
  font-family: var(--font-mono);
  font-size: 0.65rem;
  font-weight: 700;
  padding: 0 12px;
  cursor: pointer;
  border-radius: 0 var(--roundness-border) var(--roundness-border) 0;
}

.btn-toggle-visibility:hover {
  color: var(--color-primary);
  background-color: rgba(255, 255, 255, 0.02);
}

.help-text {
  display: block;
  font-size: 0.7rem;
  color: var(--color-text-dimmed);
  margin-top: 6px;
  font-family: var(--font-mono);
}

.help-text code {
  color: var(--color-primary);
  background-color: rgba(34, 211, 238, 0.05);
  padding: 1px 4px;
  border-radius: 2px;
}

.connection-status-box {
  background-color: var(--color-surface-terminal);
  border: 1px solid var(--color-outline-variant);
  padding: 10px 14px;
  border-radius: var(--roundness-border);
  margin-bottom: 20px;
}

.status-row {
  display: flex;
  justify-content: space-between;
  font-family: var(--font-mono);
  font-size: 0.8rem;
}

.status-row .label {
  color: var(--color-on-surface-variant);
}

.status-row .value.success {
  color: var(--color-success);
}

.status-row .value.error {
  color: var(--color-error);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.form-actions button {
  width: 100%;
}

.settings-help {
  margin-top: 24px;
  font-size: 0.75rem;
  line-height: 1.6;
}
</style>
