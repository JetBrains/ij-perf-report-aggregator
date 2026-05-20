<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="`LLM analysis${analysisId != null ? ` #${analysisId}` : ''}`"
    :style="{ width: '85vw', height: '85vh' }"
    :content-style="{ overflow: 'auto' }"
  >
    <div
      v-if="loading"
      class="flex items-center gap-2 text-sm"
    >
      <i class="pi pi-spin pi-spinner" />
      <span>Loading…</span>
    </div>
    <div
      v-else-if="errorMessage"
      class="text-sm text-red-600"
    >
      {{ errorMessage }}
    </div>
    <dl
      v-else-if="details"
      class="grid grid-cols-[max-content_1fr] gap-x-6 gap-y-2 text-sm"
    >
      <dt class="font-medium text-gray-500">State</dt>
      <dd class="flex items-center gap-2">
        <i :class="stateIconClass(details.state)" />
        <span>{{ stateLabel(details.state) }}</span>
      </dd>

      <dt class="font-medium text-gray-500">Created</dt>
      <dd>{{ formatCreatedAt(details.createdAt) }}</dd>

      <dt class="font-medium text-gray-500">Project</dt>
      <dd>{{ details.project }}</dd>

      <dt class="font-medium text-gray-500">Metric</dt>
      <dd>{{ details.metric }}</dd>

      <dt class="font-medium text-gray-500">Current build</dt>
      <dd>{{ details.currentBuildId }}</dd>

      <dt class="font-medium text-gray-500">Previous build</dt>
      <dd>{{ details.prevBuildId }}</dd>

      <template v-if="details.currentValue != null">
        <dt class="font-medium text-gray-500">Current value</dt>
        <dd>{{ details.currentValue }}</dd>
      </template>

      <template v-if="details.previousValue != null">
        <dt class="font-medium text-gray-500">Previous value</dt>
        <dd>{{ details.previousValue }}</dd>
      </template>

      <template v-if="details.userName">
        <dt class="font-medium text-gray-500">User</dt>
        <dd>{{ details.userName }}</dd>
      </template>

      <template v-if="details.userEmail">
        <dt class="font-medium text-gray-500">User email</dt>
        <dd>{{ details.userEmail }}</dd>
      </template>

      <template v-if="details.firstCommitRevision">
        <dt class="font-medium text-gray-500">First commit</dt>
        <dd class="font-mono">{{ details.firstCommitRevision }}</dd>
      </template>

      <template v-if="details.lastCommitRevision">
        <dt class="font-medium text-gray-500">Last commit</dt>
        <dd class="font-mono">{{ details.lastCommitRevision }}</dd>
      </template>

      <template v-if="details.testMethodName">
        <dt class="font-medium text-gray-500">Test method</dt>
        <dd class="font-mono">{{ details.testMethodName }}</dd>
      </template>

      <template v-if="details.runBuildId">
        <dt class="font-medium text-gray-500">Analyzer build</dt>
        <dd>
          <a
            :href="buildUrl(Number(details.runBuildId))"
            target="_blank"
            class="underline decoration-dotted hover:no-underline"
          >
            TC build {{ details.runBuildId }}
          </a>
        </dd>
      </template>

      <template v-if="details.ytIssueId">
        <dt class="font-medium text-gray-500">YouTrack issue</dt>
        <dd>{{ details.ytIssueId }}</dd>
      </template>

      <template v-if="details.totalCostUsd != null">
        <dt class="font-medium text-gray-500">Total cost</dt>
        <dd>${{ details.totalCostUsd.toFixed(4) }}</dd>
      </template>

      <template v-if="details.llmGuiltyCommits && details.llmGuiltyCommits.length > 0">
        <dt class="font-medium text-gray-500">Guilty commits</dt>
        <dd>
          <ul class="font-mono">
            <li
              v-for="commit in details.llmGuiltyCommits"
              :key="commit"
            >
              {{ commit }}
            </li>
          </ul>
        </dd>
      </template>

      <template v-if="details.llmComment">
        <dt class="col-span-2 mt-4 font-medium text-gray-500">LLM comment</dt>
        <dd class="col-span-2">
          <pre class="whitespace-pre-wrap rounded bg-gray-50 p-3 text-sm">{{ details.llmComment }}</pre>
        </dd>
      </template>
    </dl>
  </Dialog>
</template>

<script setup lang="ts">
import Dialog from "primevue/dialog"
import { ref, watch } from "vue"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { buildUrl } from "../common/sideBar/InfoSidebar"
import { LlmAnalysisClient, LlmAnalysisDetails, LlmAnalysisState } from "../common/llmAnalysis/LlmAnalysisClient"

const visible = defineModel<boolean>("visible", { required: true })
const { analysisId } = defineProps<{ analysisId: number | string | null }>()

const serverConfigurator = new ServerWithCompressConfigurator("", "")
const client = new LlmAnalysisClient(serverConfigurator)

const details = ref<LlmAnalysisDetails | null>(null)
const loading = ref(false)
const errorMessage = ref<string | null>(null)

watch(
  [visible, () => analysisId],
  async ([isVisible, id]) => {
    if (!isVisible) {
      details.value = null
      errorMessage.value = null
      loading.value = false
      return
    }
    if (id == null) {
      details.value = null
      errorMessage.value = null
      loading.value = false
      return
    }
    loading.value = true
    errorMessage.value = null
    details.value = null
    try {
      details.value = await client.getLlmAnalysisById(id)
    } catch (e) {
      errorMessage.value = e instanceof Error ? e.message : String(e)
    } finally {
      loading.value = false
    }
  },
  { immediate: true }
)

function stateIconClass(state: LlmAnalysisState): string {
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

function stateLabel(state: LlmAnalysisState): string {
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

function formatCreatedAt(iso: string): string {
  return new Date(iso).toLocaleString()
}
</script>
