<template>
  <div class="flex flex-col gap-4 p-6">
    <div class="flex items-center justify-between gap-4">
      <h1 class="text-xl font-semibold">LLM Analyses</h1>
      <Button
        label="Refresh"
        icon="pi pi-refresh"
        size="small"
        severity="secondary"
        :loading="loading"
        @click="load"
      />
    </div>

    <div class="flex flex-wrap items-end gap-3">
      <div class="flex flex-col gap-1">
        <label class="text-xs font-medium text-gray-500">Search</label>
        <InputText
          v-model="filter.search"
          placeholder="project, metric or user"
          class="w-72"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label class="text-xs font-medium text-gray-500">State</label>
        <MultiSelect
          v-model="filter.states"
          :options="stateOptions"
          option-label="label"
          option-value="value"
          placeholder="Any state"
          class="w-56"
          :show-toggle-all="false"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label class="text-xs font-medium text-gray-500">User</label>
        <MultiSelect
          v-model="filter.users"
          :options="userOptions"
          placeholder="Any user"
          filter
          class="w-56"
          :show-toggle-all="false"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label class="text-xs font-medium text-gray-500">From</label>
        <DatePicker
          v-model="filter.dateFrom"
          date-format="yy-mm-dd"
          show-button-bar
          placeholder="Any"
          class="w-40"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label class="text-xs font-medium text-gray-500">To</label>
        <DatePicker
          v-model="filter.dateTo"
          date-format="yy-mm-dd"
          show-button-bar
          placeholder="Any"
          class="w-40"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label class="text-xs font-medium text-gray-500">Commit</label>
        <InputText
          v-model="filter.commit"
          placeholder="commit hash"
          class="w-56 font-mono"
        />
      </div>
      <div class="flex items-center gap-4 pb-2">
        <div class="flex items-center gap-1.5">
          <Checkbox
            v-model="filter.hasTicket"
            input-id="f-ticket"
            binary
          />
          <label
            for="f-ticket"
            class="text-sm text-gray-600"
          >
            Has ticket
          </label>
        </div>
        <div class="flex items-center gap-1.5">
          <Checkbox
            v-model="filter.hasFeedback"
            input-id="f-feedback"
            binary
          />
          <label
            for="f-feedback"
            class="text-sm text-gray-600"
          >
            Has feedback
          </label>
        </div>
        <div class="flex items-center gap-1.5">
          <Checkbox
            v-model="filter.hasCommits"
            input-id="f-commits"
            binary
          />
          <label
            for="f-commits"
            class="text-sm text-gray-600"
          >
            Found commits
          </label>
        </div>
        <Button
          v-if="isFilterActive"
          label="Clear"
          icon="pi pi-filter-slash"
          size="small"
          text
          @click="clearFilters"
        />
      </div>
    </div>

    <div
      v-if="errorMessage"
      class="flex items-center justify-between gap-3 rounded border border-red-200 bg-red-50 p-3 text-sm text-red-700"
    >
      <span>{{ errorMessage }}</span>
      <Button
        label="Retry"
        size="small"
        severity="danger"
        outlined
        @click="load"
      />
    </div>

    <DataTable
      :value="filteredItems"
      :loading="loading"
      paginator
      :rows="50"
      :rows-per-page-options="[50, 100, 200]"
      sort-field="id"
      :sort-order="-1"
      striped-rows
      class="p-datatable-sm"
      row-hover
      data-key="id"
      @row-click="onRowClick"
    >
      <template #empty>
        <span class="text-sm text-gray-500">No analyses found.</span>
      </template>

      <Column
        field="createdAt"
        header="Created"
        :sortable="true"
        :style="{ width: '12rem' }"
      >
        <template #body="{ data }">
          <span v-tooltip.top="data.createdAt">{{ formatDate(data.createdAt) }}</span>
        </template>
      </Column>

      <Column
        field="state"
        header="State"
        :sortable="true"
        :style="{ width: '7rem' }"
      >
        <template #body="{ data }">
          <span class="flex items-center gap-1.5">
            <i
              v-tooltip.top="stateTooltip(data.state)"
              :class="stateIconClass(data.state)"
            />
            <span class="text-xs text-gray-600">{{ stateTooltip(data.state) }}</span>
          </span>
        </template>
      </Column>

      <Column
        field="project"
        header="Project"
        :sortable="true"
      >
        <template #body="{ data }">
          <span class="break-all">{{ data.project }}</span>
        </template>
      </Column>

      <Column
        field="metric"
        header="Metric"
        :sortable="true"
      >
        <template #body="{ data }">
          <span class="break-all">{{ data.metric }}</span>
        </template>
      </Column>

      <Column
        header="Change"
        :style="{ width: '10rem' }"
      >
        <template #body="{ data }">
          <span v-if="data.previousValue != null || data.currentValue != null"> {{ data.previousValue ?? "—" }} → {{ data.currentValue ?? "—" }} </span>
          <span
            v-else
            class="text-gray-400"
          >
            —
          </span>
        </template>
      </Column>

      <Column
        header="User"
        :style="{ width: '10rem' }"
      >
        <template #body="{ data }">
          <span>{{ userLabel(data) || "—" }}</span>
        </template>
      </Column>

      <Column
        header="Commits"
        :style="{ width: '11rem' }"
      >
        <template #body="{ data }">
          <span
            v-if="data.llmGuiltyCommits && data.llmGuiltyCommits.length > 0"
            class="flex flex-wrap gap-1"
          >
            <button
              v-for="commit in data.llmGuiltyCommits"
              :key="commit"
              v-tooltip.top="`Filter by ${commit}`"
              type="button"
              class="cursor-pointer rounded border-0 bg-gray-100 px-2 py-0.5 font-mono text-xs font-medium text-gray-700 hover:bg-primary hover:text-white dark:bg-gray-700 dark:text-gray-200"
              @click.stop="filterByCommit(commit)"
            >
              {{ commit.slice(0, 8) }}
            </button>
          </span>
          <span
            v-else
            class="text-gray-400"
          >
            —
          </span>
        </template>
      </Column>

      <Column
        header="Ticket"
        :style="{ width: '9rem' }"
      >
        <template #body="{ data }">
          <a
            v-if="data.ytIssueId"
            :href="`https://youtrack.jetbrains.com/issue/${data.ytIssueId}`"
            target="_blank"
            rel="noopener noreferrer"
            class="underline decoration-dotted hover:no-underline"
            @click.stop
          >
            {{ data.ytIssueId }}
          </a>
          <span
            v-else
            class="text-gray-400"
          >
            —
          </span>
        </template>
      </Column>

      <Column
        header="Chart"
        :style="{ width: '5rem' }"
      >
        <template #body="{ data }">
          <a
            v-if="safeDashboardLink(data.dashboardLink)"
            :href="safeDashboardLink(data.dashboardLink)!"
            target="_blank"
            rel="noopener noreferrer"
            v-tooltip.top="'Open chart'"
            class="flex items-center gap-1 text-primary hover:text-primary-dark"
            @click.stop
          >
            <i class="pi pi-chart-line" />
          </a>
          <span
            v-else
            class="text-gray-400"
          >
            —
          </span>
        </template>
      </Column>

      <Column
        header="Feedback"
        :style="{ width: '8rem' }"
      >
        <template #body="{ data }">
          <span
            v-if="data.feedbackCount > 0"
            class="flex items-center gap-1"
          >
            <i class="pi pi-star-fill text-amber-400" />
            <span>{{ data.avgRating != null ? data.avgRating.toFixed(1) : "?" }}</span>
            <span class="text-xs text-gray-500">({{ data.feedbackCount }})</span>
          </span>
          <span
            v-else
            class="text-gray-400"
          >
            —
          </span>
        </template>
      </Column>

      <Column
        field="totalCostUsd"
        header="Cost"
        :sortable="true"
        :style="{ width: '7rem' }"
      >
        <template #body="{ data }">
          <span v-if="data.totalCostUsd != null">{{ formatCost(data.totalCostUsd) }}</span>
          <span
            v-else
            class="text-gray-400"
          >
            —
          </span>
        </template>
      </Column>
    </DataTable>

    <div class="flex items-center justify-center gap-3 py-1">
      <span class="text-xs text-gray-500">{{ items.length }} run{{ items.length === 1 ? "" : "s" }} loaded</span>
      <Button
        v-if="hasMore"
        label="Load more"
        icon="pi pi-chevron-down"
        size="small"
        severity="secondary"
        outlined
        :loading="loadingMore"
        @click="loadMore"
      />
    </div>

    <AnalysisDetailsDialog
      v-model:visible="dialogVisible"
      :analysis-id="selectedAnalysisId"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, provide, ref } from "vue"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { serverConfiguratorKey, youtrackClientKey } from "../../shared/keys"
import { YoutrackClient } from "../common/youtrack/YoutrackClient"
import { LlmAnalysisClient, LlmAnalysisListItem, LlmAnalysisState } from "../common/llmAnalysis/LlmAnalysisClient"
import AnalysisDetailsDialog from "../common/llmAnalysis/AnalysisDetailsDialog.vue"
import { AnalysesFilterState, distinctUsers, emptyAnalysesFilterState, filterAnalyses, userLabel } from "./analysesFilter"

const serverConfigurator = new ServerWithCompressConfigurator("perfint", "report")
provide(serverConfiguratorKey, serverConfigurator)
const youtrackClient = new YoutrackClient(serverConfigurator)
provide(youtrackClientKey, youtrackClient)

const client = new LlmAnalysisClient(serverConfigurator)

const PAGE_SIZE = 100

const items = ref<LlmAnalysisListItem[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const hasMore = ref(false)
const errorMessage = ref<string | null>(null)

const filter = ref<AnalysesFilterState>(emptyAnalysesFilterState())

const dialogVisible = ref(false)
const selectedAnalysisId = ref<number | null>(null)

const stateOptions = [
  { label: "In progress", value: LlmAnalysisState.InProgress },
  { label: "Success", value: LlmAnalysisState.Success },
  { label: "Failed", value: LlmAnalysisState.Failed },
]

const filteredItems = computed(() => filterAnalyses(items.value, filter.value))
const userOptions = computed(() => distinctUsers(items.value))
const isFilterActive = computed(() => {
  const f = filter.value
  return f.search !== "" || f.states.length > 0 || f.users.length > 0 || f.dateFrom != null || f.dateTo != null || f.commit !== "" || f.hasTicket || f.hasFeedback || f.hasCommits
})

async function load(): Promise<void> {
  loading.value = true
  errorMessage.value = null
  try {
    const batch = await client.getAnalyses(PAGE_SIZE, 0)
    items.value = batch
    hasMore.value = batch.length === PAGE_SIZE
  } catch (e) {
    errorMessage.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

async function loadMore(): Promise<void> {
  loadingMore.value = true
  errorMessage.value = null
  try {
    const batch = await client.getAnalyses(PAGE_SIZE, items.value.length)
    items.value = [...items.value, ...batch]
    hasMore.value = batch.length === PAGE_SIZE
  } catch (e) {
    errorMessage.value = e instanceof Error ? e.message : String(e)
  } finally {
    loadingMore.value = false
  }
}

onMounted(load)

function clearFilters(): void {
  filter.value = emptyAnalysesFilterState()
}

function onRowClick(event: { data: LlmAnalysisListItem }): void {
  selectedAnalysisId.value = event.data.id
  dialogVisible.value = true
}

function filterByCommit(commit: string): void {
  filter.value.commit = commit
}

function stateIconClass(state: string): string {
  switch (state) {
    case LlmAnalysisState.InProgress:
      return "pi pi-spin pi-spinner text-sky-500"
    case LlmAnalysisState.Success:
      return "pi pi-verified text-green-600"
    case LlmAnalysisState.Failed:
      return "pi pi-times-circle text-red-500"
    default:
      return ""
  }
}

function stateTooltip(state: string): string {
  switch (state) {
    case LlmAnalysisState.InProgress:
      return "In progress"
    case LlmAnalysisState.Success:
      return "Success"
    case LlmAnalysisState.Failed:
      return "Failed"
    default:
      return state
  }
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString()
}

function formatCost(cost: number): string {
  return `$${cost.toFixed(2)}`
}

function safeDashboardLink(url: string | undefined): string | null {
  if (url == null || url === "") return null
  // Block protocol-relative URLs (//host) that could navigate cross-origin.
  if (url.startsWith("//")) return null
  try {
    // Resolve against the current origin so relative paths (the stored shape, e.g. "/intellij/tests?…")
    // are accepted, while javascript:/data: URIs resolve to a non-http(s) protocol and are rejected.
    const resolved = new URL(url, window.location.origin)
    return resolved.protocol === "http:" || resolved.protocol === "https:" ? url : null
  } catch {
    return null
  }
}
</script>

<style scoped>
/* No dropdown chevrons on the filter selects */
:deep(.p-multiselect-dropdown) {
  display: none;
}
</style>
