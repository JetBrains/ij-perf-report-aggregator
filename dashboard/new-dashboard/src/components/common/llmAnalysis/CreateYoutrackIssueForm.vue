<template>
  <div class="flex flex-col gap-3 rounded bg-gray-50 p-3 dark:bg-gray-800">
    <SelectButton
      v-model="mode"
      :options="modeOptions"
      option-label="label"
      option-value="value"
      :allow-empty="false"
      :disabled="isSubmitting"
    />

    <template v-if="mode === 'new'">
      <div class="flex flex-col gap-1">
        <label
          for="yt-project"
          class="font-medium text-gray-500"
        >
          Project
        </label>
        <Select
          id="yt-project"
          v-model="selectedProject"
          :options="ytProjects"
          option-label="name"
          placeholder="Project"
          :disabled="isSubmitting"
          class="w-80"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label
          for="yt-title"
          class="font-medium text-gray-500"
        >
          Title
        </label>
        <InputText
          id="yt-title"
          v-model="effectiveTitle"
          :disabled="isSubmitting"
          placeholder="Issue title"
          class="w-full"
        />
      </div>
    </template>

    <div
      v-else
      class="flex flex-col gap-1"
    >
      <label
        for="yt-issue-id"
        class="font-medium text-gray-500"
      >
        Issue ID or URL
      </label>
      <InputText
        id="yt-issue-id"
        v-model="issueIdInput"
        :disabled="isSubmitting"
        placeholder="IJPL-1234 or https://youtrack.jetbrains.com/issue/IJPL-1234"
        class="w-full"
      />
      <small class="text-gray-500"> The LLM analysis will be posted as a comment on this issue, which will be linked and tagged. </small>
    </div>

    <div class="flex items-center gap-2">
      <Button
        :label="submitLabel"
        icon="pi pi-check"
        :disabled="isSubmitting || submitDisabled"
        :loading="isSubmitting"
        @click="submit"
      />
      <Button
        label="Cancel"
        severity="secondary"
        :disabled="isSubmitting"
        @click="emit('cancel')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "primevue/usetoast"
import Button from "primevue/button"
import InputText from "primevue/inputtext"
import Select from "primevue/select"
import SelectButton from "primevue/selectbutton"
import { computed, ref } from "vue"
import { getSpaceUrl, InfoData } from "../sideBar/InfoSidebar"
import { generateDefaultReason, inferKindFromData } from "../sideBar/AccidentUtils"
import type { Project } from "../youtrack/YoutrackClient"
import { injectOrError, injectOrNull } from "../../../shared/injectionKeys"
import { serverConfiguratorKey, youtrackClientKey } from "../../../shared/keys"
import { fetchChartPngAsBase64 } from "../uploadAttachments/uploadAttachmentsUtils"

type Mode = "new" | "link"

const { analysisId, data } = defineProps<{
  analysisId: number | string
  data?: InfoData | null
}>()

const emit = defineEmits<{
  cancel: []
  created: [issue: { id: string; idReadable: string }, action: "created" | "linked"]
}>()

const youtrackClient = injectOrError(youtrackClientKey)
const serverConfigurator = injectOrNull(serverConfiguratorKey)
const toast = useToast()

const modeOptions: { label: string; value: Mode }[] = [
  { label: "New issue", value: "new" },
  { label: "Link existing", value: "link" },
]
const mode = ref<Mode>("new")

const ytProjects = ref<Project[]>(youtrackClient.getProjects())
const selectedProject = ref<Project | null>(ytProjects.value[0] ?? null)
const titleOverride = ref<string | null>(null)
const issueIdInput = ref("")
const isSubmitting = ref(false)

const ISSUE_ID_RE = /^[A-Z][A-Z0-9]*-\d+$/

// Accepts a bare readable id (IJPL-1234) or a YouTrack issue URL and extracts the readable id.
function parseIssueId(input: string): string | null {
  const trimmed = input.trim()
  if (trimmed === "") return null
  const fromUrl = /\/issue\/([A-Z][A-Z0-9]*-\d+)/.exec(trimmed)
  const candidate = (fromUrl?.[1] ?? trimmed).toUpperCase()
  return ISSUE_ID_RE.test(candidate) ? candidate : null
}

const defaultTitle = computed<string>(() => {
  if (data == null) return ""
  return `${inferKindFromData(data)} ${generateDefaultReason(data)} ${data.mode ? `on ${data.mode} mode` : ""}`
})

const effectiveTitle = computed<string>({
  get: () => titleOverride.value ?? defaultTitle.value,
  set: (v) => {
    titleOverride.value = v
  },
})

const submitLabel = computed<string>(() => (mode.value === "new" ? "Create" : "Link"))

const submitDisabled = computed<boolean>(() => (mode.value === "new" ? selectedProject.value == null : issueIdInput.value.trim() === ""))

async function submit() {
  await (mode.value === "new" ? submitNew() : submitLink())
}

async function submitNew() {
  if (selectedProject.value == null) return
  const title = effectiveTitle.value.trim()
  if (title.length < 5) {
    toast.add({ severity: "error", summary: "Validation Error", detail: "Title must be at least 5 characters long", life: 5000 })
    return
  }
  isSubmitting.value = true
  try {
    const { changesLink, delta, chartPng } = await collectAnalysisContext()
    const resp = await youtrackClient.createIssueByAnalysis(Number(analysisId), {
      projectId: selectedProject.value.id,
      ticketLabel: title,
      delta,
      changesLink,
      chartPng,
    })
    toast.add({ severity: "success", summary: "Issue created", detail: resp.issue.idReadable, life: 4000 })
    emit("created", { id: resp.issue.id, idReadable: resp.issue.idReadable }, "created")
  } catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    toast.add({ severity: "error", summary: "Issue creation failed", detail: msg, life: 8000 })
  } finally {
    isSubmitting.value = false
  }
}

async function submitLink() {
  const issueId = parseIssueId(issueIdInput.value)
  if (issueId == null) {
    toast.add({ severity: "error", summary: "Validation Error", detail: "Enter a valid issue ID (e.g. IJPL-1234) or YouTrack issue URL", life: 5000 })
    return
  }
  isSubmitting.value = true
  try {
    const resp = await youtrackClient.linkIssueByAnalysis(Number(analysisId), { issueId })
    toast.add({ severity: "success", summary: "Issue linked", detail: resp.issue.idReadable, life: 4000 })
    emit("created", { id: resp.issue.id, idReadable: resp.issue.idReadable }, "linked")
  } catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    toast.add({ severity: "error", summary: "Issue link failed", detail: msg, life: 8000 })
  } finally {
    isSubmitting.value = false
  }
}

async function collectAnalysisContext(): Promise<{ changesLink: string; delta: string; chartPng: string | undefined }> {
  let changesLink = ""
  let delta = ""
  let chartPng: string | undefined
  if (data != null) {
    const spaceUrls = await getSpaceUrl(data, serverConfigurator)
    changesLink = spaceUrls.length > 0 ? spaceUrls.join(",") : data.changesUrl
    delta = data.deltaPrevious?.replaceAll(/[+-]/g, (match) => (match === "+" ? "-" : "+")) ?? ""
    if (data.chartDataUrl) {
      try {
        chartPng = await fetchChartPngAsBase64(data.chartDataUrl)
      } catch (e) {
        console.error("Failed to prepare chart for upload", e)
        toast.add({ severity: "warn", summary: "Chart not attached", detail: "Failed to prepare chart for upload; the issue will be created without it.", life: 5000 })
      }
    }
  }
  return { changesLink, delta, chartPng }
}
</script>
