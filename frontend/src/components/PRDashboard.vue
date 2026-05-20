<template>
  <div class="pr-dashboard">
    <!-- Unconfigured State Hint -->
    <div v-if="!state.token || !state.username" class="panel setup-hint-panel">
      <div class="panel-header">
        <span class="panel-title">// GIT_INITIALIZE.err</span>
      </div>
      <div class="panel-body monospaced">
        <p class="error-text">No active session found.</p>
        <p class="dimmed-text">Please set up your GitHub Personal Access Token in Settings to display your PR workspace telemetry.</p>
        <div class="terminal-box">
          <span class="prompt">$</span> git config --global assistant.token "ghp_..."<br>
          <span class="prompt">Error:</span> Missing authentication credential key.
        </div>
      </div>
    </div>

    <!-- Main Dashboard -->
    <div v-else>
      <!-- Metrics Grid -->
      <div class="metrics-grid">
        <!-- Lead Time Metric -->
        <div class="metric-card">
          <span class="metric-label">LEAD_TIME</span>
          <span class="metric-value">
            {{ formatHours(dashboardData.metrics?.avg_lead_time_hours) }}
          </span>
          <span class="metric-sub text-dimmed">commit to merge</span>
        </div>

        <!-- Review Velocity Metric -->
        <div class="metric-card">
          <span class="metric-label">REV_VELOCITY</span>
          <span class="metric-value">
            {{ formatHours(dashboardData.metrics?.avg_review_velocity_hours) }}
          </span>
          <span class="metric-sub text-dimmed">creation to review</span>
        </div>

        <!-- PR Count Metric -->
        <div class="metric-card">
          <span class="metric-label">TOTAL_PRs</span>
          <div class="metric-value pr-counts">
            <span>{{ dashboardData.metrics?.total_count || 0 }}</span>
            <div class="pr-breakdown">
              <span class="count-item approved" title="Approved">{{ dashboardData.metrics?.approved_count || 0 }}A</span>
              <span class="count-item pending" title="Changes Requested">{{ dashboardData.metrics?.changes_requested_count || 0 }}C</span>
            </div>
          </div>
          <span class="metric-sub text-dimmed">active pull requests</span>
        </div>
      </div>

      <!-- Filter Sub-Tabs -->
      <div class="tab-container">
        <button 
          :class="['tab-btn', { active: activeQueryType === 'reviewing' }]" 
          @click="setQueryType('reviewing')"
        >
          REVIEWING ({{ activeQueryType === 'reviewing' ? dashboardData.prs?.length || 0 : '?' }})
        </button>
        <button 
          :class="['tab-btn', { active: activeQueryType === 'mine' }]" 
          @click="setQueryType('mine')"
        >
          MY_PRs ({{ activeQueryType === 'mine' ? dashboardData.prs?.length || 0 : '?' }})
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="terminal-box loading-box">
        <span class="prompt">></span> FETCHING FROM GITHUB API...<br>
        <span class="loading-indicator">|||||||||||||||||||||| 100%</span>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="terminal-box error-box">
        <span class="prompt warning">Error:</span> Failed to fetch pull request telemetry.<br>
        <span class="dimmed-text">{{ error }}</span>
      </div>

      <!-- Empty State -->
      <div v-else-if="!dashboardData.prs || dashboardData.prs.length === 0" class="panel empty-panel">
        <div class="panel-body monospaced text-center">
          <p class="dimmed-text">> No pull requests found matching the current workspace filters.</p>
        </div>
      </div>

      <!-- PRs List -->
      <div v-else class="prs-list">
        <div 
          v-for="pr in dashboardData.prs" 
          :key="pr.id" 
          class="pr-card"
        >
          <div class="pr-card-header">
            <div class="pr-meta-info">
              <img :src="pr.author_avatar" class="author-avatar" alt="Avatar" />
              <span class="pr-repo monospaced">{{ pr.repo }}</span>
              <span class="pr-number monospaced">#{{ pr.number }}</span>
            </div>
            
            <!-- PR Status Badge -->
            <span :class="['badge', `badge-${pr.status}`]">
              {{ formatStatus(pr.status) }}
            </span>
          </div>

          <div class="pr-card-body">
            <h3 class="pr-title">{{ pr.title }}</h3>
          </div>

          <div class="pr-card-footer monospaced">
            <div class="pr-metrics-detail">
              <span v-if="pr.review_velocity_hours !== null" class="metric-pill">
                VEL: {{ pr.review_velocity_hours.toFixed(1) }}h
              </span>
              <span v-if="pr.lead_time_hours !== null" class="metric-pill">
                LT: {{ pr.lead_time_hours.toFixed(1) }}h
              </span>
              <span v-if="pr.is_draft" class="draft-pill">DRAFT</span>
            </div>
            <a :href="pr.url" target="_blank" rel="noopener noreferrer" class="btn btn-dimmed btn-open">
              OPEN_PR
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, watch, onMounted } from 'vue';

export default {
  name: 'PRDashboard',
  props: {
    state: {
      type: Object,
      required: true
    }
  },
  setup(props) {
    const activeQueryType = ref('reviewing');
    const loading = ref(false);
    const error = ref(null);
    const dashboardData = reactive({
      prs: [],
      metrics: null
    });

    const setQueryType = (type) => {
      activeQueryType.value = type;
      fetchDashboardData();
    };

    // Load data
    const fetchDashboardData = async () => {
      if (!props.state.token || !props.state.username) return;

      loading.value = true;
      error.value = null;

      try {
        const response = await fetch(
          `/api/github/prs?type=${activeQueryType.value}&username=${encodeURIComponent(props.state.username)}`,
          {
            headers: {
              'Authorization': `Bearer ${props.state.token}`
            }
          }
        );

        if (!response.ok) {
          const errData = await response.json();
          throw new Error(errData.error || `HTTP error ${response.status}`);
        }

        const data = await response.json();
        dashboardData.prs = data.prs || [];
        dashboardData.metrics = data.metrics || null;
      } catch (err) {
        console.error("Fetch dashboard data error:", err);
        error.value = err.message || "Failed to load PR dashboard.";
        dashboardData.prs = [];
        dashboardData.metrics = null;
      } finally {
        loading.value = false;
      }
    };

    // Format metrics time
    const formatHours = (val) => {
      if (val === undefined || val === null || isNaN(val)) return 'N/A';
      if (val < 1) return `${(val * 60).toFixed(0)}m`;
      return `${val.toFixed(1)}h`;
    };

    // Clean status name
    const formatStatus = (status) => {
      return status ? status.replace('_', ' ') : 'reviewing';
    };

    // Refetch when credentials change
    watch(() => props.state.isConnected, (connected) => {
      if (connected) {
        fetchDashboardData();
      } else {
        dashboardData.prs = [];
        dashboardData.metrics = null;
      }
    });

    onMounted(() => {
      if (props.state.isConnected) {
        fetchDashboardData();
      }
    });

    return {
      activeQueryType,
      loading,
      error,
      dashboardData,
      setQueryType,
      formatHours,
      formatStatus
    };
  }
};
</script>

<style scoped>
.error-text {
  color: var(--color-error);
  font-weight: 600;
  margin-bottom: 8px;
}

/* Metrics Grid Style */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--grid-gap);
  margin-bottom: 20px;
}

.metric-card {
  background-color: var(--color-surface-container);
  border: 1px solid var(--color-outline);
  border-radius: var(--roundness-border);
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  height: 80px;
  text-align: center;
}

.metric-label {
  font-family: var(--font-mono);
  font-size: 0.6rem;
  color: var(--color-text-secondary);
  font-weight: 700;
  letter-spacing: 0.5px;
}

.metric-value {
  font-family: var(--font-mono);
  font-size: 1.15rem;
  color: var(--color-primary);
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 4px;
}

.pr-counts {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.pr-breakdown {
  display: flex;
  flex-direction: column;
  font-size: 0.5rem;
  line-height: 1.1;
}

.count-item.approved {
  color: var(--color-approved);
}

.count-item.pending {
  color: var(--color-changes-req);
}

.metric-sub {
  font-size: 0.55rem;
}

/* Loading Box */
.loading-box {
  margin-bottom: 20px;
}

.loading-indicator {
  color: var(--color-primary);
  letter-spacing: 2px;
}

.error-box {
  border-color: var(--color-error);
  margin-bottom: 20px;
}

/* PR Cards List */
.prs-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.pr-card {
  background-color: var(--color-surface-container);
  border: 1px solid var(--color-outline);
  border-radius: var(--roundness-border);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.pr-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.pr-meta-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-avatar {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1px solid var(--color-outline);
}

.pr-repo {
  font-size: 0.7rem;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.pr-number {
  font-size: 0.7rem;
  color: var(--color-text-dimmed);
}

.pr-card-body {
  padding: 2px 0;
}

.pr-title {
  font-size: 0.85rem;
  font-weight: 500;
  line-height: 1.3;
  color: var(--color-text-primary);
}

.pr-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 4px;
}

.pr-metrics-detail {
  display: flex;
  gap: 6px;
  font-size: 0.65rem;
}

.metric-pill {
  background-color: rgba(34, 211, 238, 0.05);
  border: 1px solid var(--color-outline-dimmed);
  color: var(--color-primary);
  padding: 1px 6px;
  border-radius: 2px;
}

.draft-pill {
  background-color: rgba(100, 116, 139, 0.1);
  border: 1px solid var(--color-outline-dimmed);
  color: var(--color-draft);
  padding: 1px 6px;
  border-radius: 2px;
}

.btn-open {
  font-size: 0.65rem;
  padding: 4px 10px;
}

.text-center {
  text-align: center;
}
</style>
