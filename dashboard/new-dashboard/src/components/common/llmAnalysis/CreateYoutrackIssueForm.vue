<template>
  <div class="flex flex-col gap-3 rounded bg-gray-50 p-3 dark:bg-gray-800">
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
    <div class="flex items-center gap-2">
      <Button
        label="Create"
        icon="pi pi-check"
        :disabled="isSubmitting || selectedProject == null"
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
import { computed, ref } from "vue"
import { getSpaceUrl, InfoData } from "../sideBar/InfoSidebar"
import { generateDefaultReason, inferKindFromData } from "../sideBar/AccidentUtils"
import type { Project } from "../youtrack/YoutrackClient"
import { injectOrError, injectOrNull } from "../../../shared/injectionKeys"
import { serverConfiguratorKey, youtrackClientKey } from "../../../shared/keys"

const { analysisId, data } = defineProps<{
  analysisId: number | string
  data?: InfoData | null
}>()

const emit = defineEmits<{
  cancel: []
  created: [issue: { id: string; idReadable: string }]
}>()

const youtrackClient = injectOrError(youtrackClientKey)
const serverConfigurator = injectOrNull(serverConfiguratorKey)
const toast = useToast()

const ytProjects = ref<Project[]>(youtrackClient.getProjects())
const selectedProject = ref<Project | null>(ytProjects.value[0] ?? null)
const titleOverride = ref<string | null>(null)
const isSubmitting = ref(false)

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

async function submit() {
  if (selectedProject.value == null) return
  const title = effectiveTitle.value.trim()
  if (title.length < 5) {
    toast.add({ severity: "error", summary: "Validation Error", detail: "Title must be at least 5 characters long", life: 5000 })
    return
  }
  isSubmitting.value = true
  try {
    let changesLink = ""
    let delta = ""
    if (data != null) {
      const spaceUrls = await getSpaceUrl(data, serverConfigurator)
      changesLink = spaceUrls.length > 0 ? spaceUrls.join(",") : data.changesUrl
      delta = data.deltaPrevious?.replaceAll(/[+-]/g, (match) => (match === "+" ? "-" : "+")) ?? ""
    }
    const resp = await youtrackClient.createIssueByAnalysis(Number(analysisId), {
      projectId: selectedProject.value.id,
      ticketLabel: title,
      delta,
      changesLink,
    })
    toast.add({ severity: "success", summary: "Issue created", detail: resp.issue.idReadable, life: 4000 })
    emit("created", { id: resp.issue.id, idReadable: resp.issue.idReadable })
  } catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    toast.add({ severity: "error", summary: "Issue creation failed", detail: msg, life: 8000 })
  } finally {
    isSubmitting.value = false
  }
}
</script>
