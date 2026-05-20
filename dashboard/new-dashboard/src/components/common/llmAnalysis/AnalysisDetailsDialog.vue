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
              <dd class="font-mono">{{ details.firstCommitRevision ?? details.lastCommitRevision }}</dd>
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
              <dt class="col-span-2 mt-4 border-t border-gray-200 pt-3 text-base font-semibold text-gray-900">Result</dt>
              <dd class="col-span-2">
                <div
                  class="markdown-body rounded bg-gray-50 p-3 text-base"
                  v-html="renderedComment"
                />
              </dd>
            </template>
          </dl>
        </TabPanel>
        <TabPanel :value="1">
          <div class="flex flex-col gap-3 text-sm">
            <div
              v-if="feedbacksLoading"
              class="text-xs text-gray-500"
            >
              Loading feedback…
            </div>
            <div
              v-else-if="feedbacksError"
              class="text-xs text-red-600"
            >
              {{ feedbacksError }}
            </div>
            <ul
              v-else-if="feedbacks.length > 0"
              class="flex flex-col gap-1 border-b border-gray-200 pb-2 text-xs"
            >
              <li
                v-for="fb in feedbacks"
                :key="fb.id"
                class="flex items-center gap-2"
              >
                <span class="font-medium whitespace-nowrap">★ {{ fb.rate }}/5</span>
                <span
                  v-if="fb.userEmail"
                  class="text-gray-500 whitespace-nowrap"
                >
                  {{ userLocalPart(fb.userEmail) }}
                </span>
                <span class="text-gray-400 whitespace-nowrap">{{ formatCreatedAt(fb.createdAt) }}</span>
                <span
                  v-if="fb.feedback"
                  v-tooltip.top="fb.feedback"
                  class="truncate text-gray-700"
                >
                  {{ fb.feedback }}
                </span>
              </li>
            </ul>
            <div class="flex flex-col gap-1">
              <label
                for="analysis-feedback-rating"
                class="font-medium text-gray-500"
              >
                Rating
              </label>
              <Select
                id="analysis-feedback-rating"
                v-model="rating"
                :options="ratingOptions"
                option-label="label"
                option-value="value"
                placeholder="Select a rating"
                class="w-80"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label
                for="analysis-feedback-comments"
                class="font-medium text-gray-500"
              >
                Feedback
              </label>
              <Textarea
                id="analysis-feedback-comments"
                v-model="feedbackText"
                rows="5"
                class="w-full"
                placeholder="What could work better, ideally the actually guilty commit(s)"
              />
            </div>
            <div
              v-tooltip.top="submitDisabledReason"
              class="inline-block w-fit"
            >
              <Button
                label="Send"
                :disabled="submitDisabledReason != null || submitting"
                :loading="submitting"
                @click="submitFeedback"
              />
            </div>
          </div>
        </TabPanel>
      </TabPanels>
    </Tabs>
  </Dialog>
</template>

<script setup lang="ts">
import { Marked } from "marked"
import Dialog from "primevue/dialog"
import { computed, ref, watch } from "vue"
import { buildUrl } from "../sideBar/InfoSidebar"
import { AnalysisFeedback, LlmAnalysisClient, LlmAnalysisDetails, LlmAnalysisState } from "./LlmAnalysisClient"
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
const { analysisId } = defineProps<{ analysisId: number | string | null }>()

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
const rating = ref<number | null>(null)
const feedbackText = ref("")

const feedbacks = ref<AnalysisFeedback[]>([])
const feedbacksLoaded = ref(false)
const feedbacksLoading = ref(false)
const feedbacksError = ref<string | null>(null)
const submitting = ref(false)

const ratingOptions = [
  { value: 1, label: "1 — Poor" },
  { value: 2, label: "2 — Fair" },
  { value: 3, label: "3 — Good" },
  { value: 4, label: "4 — Very good" },
  { value: 5, label: "5 — Excellent" },
]

const isTerminalState = computed(() => details.value != null && details.value.state !== LlmAnalysisState.InProgress)

const submitDisabledReason = computed<string | null>(() => {
  if (!isTerminalState.value) return "Please wait till analysis is finished"
  if (rating.value == null) return "Select a rating"
  return null
})

function resetFeedback() {
  activeTab.value = 0
  rating.value = null
  feedbackText.value = ""
  feedbacks.value = []
  feedbacksLoaded.value = false
  feedbacksError.value = null
}

async function loadFeedbacks() {
  if (analysisId == null || feedbacksLoaded.value || feedbacksLoading.value) return
  feedbacksLoading.value = true
  feedbacksError.value = null
  try {
    feedbacks.value = await client.getAnalysisFeedback(analysisId)
    feedbacksLoaded.value = true
  } catch (e) {
    feedbacksError.value = e instanceof Error ? e.message : String(e)
  } finally {
    feedbacksLoading.value = false
  }
}

async function submitFeedback() {
  if (rating.value == null || analysisId == null || submitting.value) return
  submitting.value = true
  try {
    await client.submitAnalysisFeedback(analysisId, rating.value, feedbackText.value)
    resetFeedback()
  } catch (e) {
    feedbacksError.value = e instanceof Error ? e.message : String(e)
  } finally {
    submitting.value = false
  }
}

watch(activeTab, (t) => {
  if (t === 1) void loadFeedbacks()
})

watch(
  [visible, () => analysisId],
  async ([isVisible, id]) => {
    details.value = null
    errorMessage.value = null
    loading.value = false
    resetFeedback()
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
</style>
