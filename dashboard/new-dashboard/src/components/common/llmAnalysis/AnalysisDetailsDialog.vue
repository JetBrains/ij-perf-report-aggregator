<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: 'min(85vw, 1100px)', height: '85vh' }"
    :content-style="{ overflow: 'auto' }"
  >
    <template #header>
      <div class="flex items-center gap-2 text-base font-semibold">
        <i
          v-if="details"
          :class="stateIconClass(details.state)"
        />
        <span>{{ details ? `Launched ${formatCreatedAt(details.createdAt)}` : "LLM analysis" }}</span>
        <span
          v-if="details?.userEmail"
          class="font-normal text-gray-500"
        >
          by {{ userLocalPart(details.userEmail) }}
        </span>
      </div>
    </template>
    <Tabs
      v-model:value="activeTab"
      class="flex flex-col flex-1 min-h-0"
    >
      <TabList class="flex-shrink-0">
        <Tab :value="0">Analysis</Tab>
        <Tab :value="1">Feedback</Tab>
      </TabList>
      <TabPanels class="flex-1 min-h-0 overflow-auto">
        <TabPanel :value="0">
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
            <dt class="font-medium text-gray-500">Project</dt>
            <dd>{{ details.project }}</dd>

            <dt class="font-medium text-gray-500">Metric</dt>
            <dd>{{ details.metric }}</dd>

            <dt class="font-medium text-gray-500">Current build</dt>
            <dd>
              <a
                :href="buildUrl(Number(details.currentBuildId))"
                target="_blank"
                class="underline decoration-dotted hover:no-underline"
              >
                {{ details.currentBuildId }}
              </a>
            </dd>

            <dt class="font-medium text-gray-500">Previous build</dt>
            <dd>
              <a
                :href="buildUrl(Number(details.prevBuildId))"
                target="_blank"
                class="underline decoration-dotted hover:no-underline"
              >
                {{ details.prevBuildId }}
              </a>
            </dd>

            <template v-if="details.currentValue != null || details.previousValue != null">
              <dt class="font-medium text-gray-500">Metric change</dt>
              <dd>{{ details.previousValue ?? "—" }} → {{ details.currentValue ?? "—" }}</dd>
            </template>

            <template v-if="details.firstCommitRevision && details.lastCommitRevision">
              <dt class="font-medium text-gray-500">Commits range</dt>
              <dd class="font-mono">{{ details.firstCommitRevision }}^..{{ details.lastCommitRevision }}</dd>
            </template>
            <template v-else-if="details.firstCommitRevision || details.lastCommitRevision">
              <dt class="font-medium text-gray-500">Commit</dt>
              <dd class="font-mono">{{ details.firstCommitRevision || details.lastCommitRevision }}</dd>
            </template>
            <template v-else-if="details.firstCommitRevision === '' && details.lastCommitRevision === ''">
              <dt class="font-medium text-gray-500">Commits range</dt>
              <dd class="text-gray-500 italic">empty</dd>
            </template>

            <template v-if="details.runBuildId">
              <dt class="font-medium text-gray-500">Analysis artifacts</dt>
              <dd>
                <a
                  :href="`${buildUrl(Number(details.runBuildId))}&buildTab=artifacts`"
                  target="_blank"
                  class="underline decoration-dotted hover:no-underline"
                >
                  {{ details.runBuildId }}
                </a>
              </dd>
            </template>

            <template v-if="details.llmComment">
              <dt class="col-span-2 mt-4 border-t border-gray-200 pt-3 text-base font-semibold text-gray-900">
                <div class="flex items-center justify-between gap-3">
                  <span>{{ showCreateForm ? "Create YouTrack issue" : "Result" }}</span>
                  <Button
                    v-if="!showCreateForm && createdIssue == null && isTerminalState"
                    label="Create YouTrack issue"
                    icon="pi pi-plus"
                    size="small"
                    @click="showCreateForm = true"
                  />
                  <a
                    v-else-if="!showCreateForm && createdIssue"
                    :href="`https://youtrack.jetbrains.com/issue/${createdIssue.idReadable}`"
                    target="_blank"
                    class="flex items-center gap-1 text-sm font-normal underline decoration-dotted hover:no-underline"
                  >
                    <i class="pi pi-verified text-green-600" />
                    Created {{ createdIssue.idReadable }} ↗
                  </a>
                </div>
              </dt>
              <dd class="col-span-2">
                <CreateYoutrackIssueForm
                  v-if="showCreateForm && analysisId != null"
                  :analysis-id="analysisId"
                  :data="data"
                  @cancel="showCreateForm = false"
                  @created="onIssueCreated"
                />
                <div
                  v-else
                  class="markdown-body rounded bg-gray-50 p-3 text-base dark:bg-gray-800 dark:text-gray-100"
                  v-html="renderedComment"
                />
              </dd>
            </template>
          </dl>
        </TabPanel>
        <TabPanel :value="1">
          <AnalysisFeedbackTab
            v-if="visible"
            :analysis-id="analysisId"
            :is-terminal-state="isTerminalState"
            :active="activeTab === 1"
          />
        </TabPanel>
      </TabPanels>
    </Tabs>
  </Dialog>
</template>

<script setup lang="ts">
import { Marked } from "marked"
import Dialog from "primevue/dialog"
import { computed, ref, watch } from "vue"
import { buildUrl, InfoData } from "../sideBar/InfoSidebar"
import { LlmAnalysisClient, LlmAnalysisDetails, LlmAnalysisState } from "./LlmAnalysisClient"
import AnalysisFeedbackTab from "./AnalysisFeedbackTab.vue"
import CreateYoutrackIssueForm from "./CreateYoutrackIssueForm.vue"
import { injectOrNull } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"

const escapeHtml = (s: string): string => s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;").replaceAll('"', "&quot;")

const SAFE_HREF = /^(https?:|\/|\.\.?\/)/i

const md = new Marked({
  gfm: true,
  breaks: false,
  async: false,
  walkTokens(token) {
    if (token.type === "html") {
      token.text = escapeHtml(token.text)
    } else if ((token.type === "link" || token.type === "image") && !SAFE_HREF.test(token.href)) {
      token.href = ""
    }
  },
})

const visible = defineModel<boolean>("visible", { required: true })
const { analysisId } = defineProps<{ analysisId: number | string | null; data?: InfoData | null }>()

const serverConfigurator = injectOrNull(serverConfiguratorKey)
const client = new LlmAnalysisClient(serverConfigurator)

const details = ref<LlmAnalysisDetails | null>(null)
const loading = ref(false)
const errorMessage = ref<string | null>(null)

const renderedComment = computed(() => {
  const text = details.value?.llmComment
  return text == null || text === "" ? "" : (md.parse(text) as string)
})

const activeTab = ref<0 | 1>(0)

const createdIssue = ref<{ id: string; idReadable: string } | null>(null)
const showCreateForm = ref(false)

const isTerminalState = computed(() => details.value != null && details.value.state !== LlmAnalysisState.InProgress)

function onIssueCreated(issue: { id: string; idReadable: string }) {
  createdIssue.value = issue
  if (details.value != null) details.value.ytIssueId = issue.idReadable
  showCreateForm.value = false
}

watch(
  [visible, () => analysisId],
  async ([isVisible, id]) => {
    details.value = null
    errorMessage.value = null
    loading.value = false
    activeTab.value = 0
    createdIssue.value = null
    showCreateForm.value = false
    if (!isVisible || id == null) return
    loading.value = true
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

function formatCreatedAt(iso: string): string {
  return new Date(iso).toLocaleString()
}

function userLocalPart(email: string): string {
  const at = email.indexOf("@")
  return at === -1 ? email : email.slice(0, at)
}
</script>

<style scoped>
.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4) {
  font-weight: 600;
  margin: 0.75rem 0 0.25rem;
}
.markdown-body :deep(h1) {
  font-size: 1.125rem;
}
.markdown-body :deep(h2) {
  font-size: 1rem;
}
.markdown-body :deep(h3),
.markdown-body :deep(h4) {
  font-size: 0.9375rem;
}

.markdown-body :deep(p) {
  margin: 0.5rem 0;
}
.markdown-body :deep(p:first-child) {
  margin-top: 0;
}
.markdown-body :deep(p:last-child) {
  margin-bottom: 0;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  margin: 0.5rem 0;
  padding-left: 1.5rem;
}
.markdown-body :deep(ul) {
  list-style: disc;
}
.markdown-body :deep(ol) {
  list-style: decimal;
}
.markdown-body :deep(li) {
  margin: 0.125rem 0;
}

.markdown-body :deep(a) {
  color: #2563eb;
  text-decoration: underline;
  text-decoration-style: dotted;
}
.markdown-body :deep(a:hover) {
  text-decoration: none;
}

.markdown-body :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 0.85em;
  background: #e5e7eb;
  padding: 0.05rem 0.25rem;
  border-radius: 0.25rem;
}
.markdown-body :deep(pre) {
  background: #1f2937;
  color: #f9fafb;
  padding: 0.75rem;
  border-radius: 0.375rem;
  overflow-x: auto;
  margin: 0.5rem 0;
}
.markdown-body :deep(pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
  font-size: 0.85em;
}

.markdown-body :deep(blockquote) {
  border-left: 3px solid #d1d5db;
  margin: 0.5rem 0;
  padding-left: 0.75rem;
  color: #4b5563;
}

.markdown-body :deep(table) {
  border-collapse: collapse;
  margin: 0.5rem 0;
}
.markdown-body :deep(th),
.markdown-body :deep(td) {
  border: 1px solid #d1d5db;
  padding: 0.25rem 0.5rem;
}
.markdown-body :deep(th) {
  background: #e5e7eb;
  font-weight: 600;
}

.dark-mode .markdown-body :deep(code) {
  background: #374151;
  color: #f3f4f6;
}
.dark-mode .markdown-body :deep(pre) {
  background: #0f172a;
}
.dark-mode .markdown-body :deep(a) {
  color: #60a5fa;
}
.dark-mode .markdown-body :deep(blockquote) {
  border-left-color: #4b5563;
  color: #d1d5db;
}
.dark-mode .markdown-body :deep(th),
.dark-mode .markdown-body :deep(td) {
  border-color: #4b5563;
}
.dark-mode .markdown-body :deep(th) {
  background: #374151;
}
</style>
