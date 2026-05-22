<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '85vw', height: '85vh' }"
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
        <dt class="col-span-2 mt-4 font-bold text-gray-900">Result</dt>
        <dd class="col-span-2">
          <div
            class="markdown-body rounded bg-gray-50 p-3 text-sm"
            v-html="renderedComment"
          />
        </dd>
      </template>
    </dl>
  </Dialog>
</template>

<script setup lang="ts">
import { marked } from "marked"
import Dialog from "primevue/dialog"
import { computed, ref, watch } from "vue"
import { ServerWithCompressConfigurator } from "../../../configurators/ServerWithCompressConfigurator"
import { buildUrl } from "../sideBar/InfoSidebar"
import { LlmAnalysisClient, LlmAnalysisDetails, LlmAnalysisState } from "./LlmAnalysisClient"

const escapeHtml = (s: string): string => s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;").replaceAll('"', "&quot;")

marked.use({
  gfm: true,
  breaks: false,
  walkTokens(token) {
    if (token.type === "html") {
      token.text = escapeHtml(token.text)
    }
  },
})

const visible = defineModel<boolean>("visible", { required: true })
const { analysisId } = defineProps<{ analysisId: number | string | null }>()

const serverConfigurator = new ServerWithCompressConfigurator("", "")
const client = new LlmAnalysisClient(serverConfigurator)

const details = ref<LlmAnalysisDetails | null>(null)
const loading = ref(false)
const errorMessage = ref<string | null>(null)

const renderedComment = computed(() => {
  const text = details.value?.llmComment
  return text == null || text === "" ? "" : (marked.parse(text) as string)
})

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
