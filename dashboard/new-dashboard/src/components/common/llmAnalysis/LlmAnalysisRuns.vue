<template>
  <Accordion
    v-if="runs.length > 0"
    value="0"
    class="llm-analysis-runs"
  >
    <AccordionPanel value="0">
      <AccordionHeader>LLM analyses</AccordionHeader>
      <AccordionContent>
        <ul class="gap-1.5 break-all">
          <li
            v-for="run in runs"
            :key="run.id"
            class="flex items-center justify-between gap-1.5"
          >
            <span class="flex items-center gap-1.5">
              <i
                v-tooltip.left="stateTooltip(run.state)"
                :class="stateIconClass(run.state)"
              />
              <span v-tooltip.top="run.createdAt">
                {{ formatCreatedAt(run.createdAt) }}
              </span>
            </span>
            <button
              type="button"
              class="underline decoration-dotted hover:no-underline bg-transparent border-0 p-0 cursor-pointer"
              @click="openAnalysis(run.id)"
            >
              Details
            </button>
          </li>
        </ul>
      </AccordionContent>
    </AccordionPanel>
  </Accordion>
  <AnalysisDetailsDialog
    v-model:visible="dialogVisible"
    :analysis-id="selectedAnalysisId"
    :data="data"
  />
</template>
<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { injectOrError } from "../../../shared/injectionKeys"
import { llmAnalysesConfiguratorKey } from "../../../shared/keys"
import { useSelectedPointStore } from "../../../shared/selectedPointStore"
import type { InfoData } from "../sideBar/InfoSidebar"
import AnalysisDetailsDialog from "./AnalysisDetailsDialog.vue"
import { LlmAnalysisState } from "./LlmAnalysisClient"

defineProps<{ data?: InfoData | null }>()

const llmAnalysesConfigurator = injectOrError(llmAnalysesConfiguratorKey)
const selectedPointStore = useSelectedPointStore()

const runs = computed(() => llmAnalysesConfigurator.value.value)

const dialogVisible = ref(false)
const selectedAnalysisId = ref<number | null>(null)

function openAnalysis(id: number): void {
  selectedAnalysisId.value = id
  dialogVisible.value = true
}

let hasAttemptedAutoOpen = false
watch(runs, (loaded) => {
  if (hasAttemptedAutoOpen || dialogVisible.value) return
  const rawId = selectedPointStore.selectedAnalysisId
  const target = Array.isArray(rawId) ? rawId[0] : rawId
  if (target == null || target === "") {
    hasAttemptedAutoOpen = true
    return
  }
  const numericTarget = Number(target)
  if (!Number.isFinite(numericTarget)) {
    hasAttemptedAutoOpen = true
    return
  }
  if (loaded.length === 0) return
  hasAttemptedAutoOpen = true
  const match = loaded.find((r) => r.id === numericTarget)
  if (match != null) openAnalysis(match.id)
})

function stateIconClass(state: LlmAnalysisState): string {
  switch (state) {
    case LlmAnalysisState.InProgress:
      return "pi pi-spin pi-spinner"
    case LlmAnalysisState.Success:
      return "pi pi-verified"
    case LlmAnalysisState.Failed:
      return "pi pi-times-circle"
    default:
      return ""
  }
}

function stateTooltip(state: LlmAnalysisState): string {
  switch (state) {
    case LlmAnalysisState.InProgress:
      return "In progress"
    case LlmAnalysisState.Success:
      return "Success"
    case LlmAnalysisState.Failed:
      return "Failed"
    default:
      return ""
  }
}

function formatCreatedAt(iso: string): string {
  return new Date(iso).toLocaleString()
}
</script>
<style scoped>
.llm-analysis-runs :deep(.p-accordionheader) {
  padding: 0 0 1rem 0;
}

.llm-analysis-runs :deep(.p-accordioncontent) {
  padding: 0;
}

.llm-analysis-runs .pi-spinner {
  color: dodgerblue;
}

.llm-analysis-runs .pi-verified {
  color: #22c55e;
}

.llm-analysis-runs .pi-times-circle {
  color: #ef4444;
}
</style>
