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
            <button
              v-tooltip.top="run.createdAt"
              type="button"
              class="text-xs underline decoration-dotted hover:no-underline bg-transparent border-0 p-0 cursor-pointer text-left"
              @click="openAnalysis(run.id)"
            >
              {{ formatCreatedAt(run.createdAt) }}
            </button>
            <span class="flex gap-1.5 items-center">
              <a
                v-if="run.runBuildId"
                :href="buildUrl(Number(run.runBuildId))"
                target="_blank"
                class="underline decoration-dotted hover:no-underline"
              >
                TC build
              </a>
              <i
                v-tooltip.left="stateTooltip(run.state)"
                :class="stateIconClass(run.state)"
              />
            </span>
          </li>
        </ul>
      </AccordionContent>
    </AccordionPanel>
  </Accordion>
  <AnalysisDetailsDialog
    v-model:visible="dialogVisible"
    :analysis-id="selectedAnalysisId"
  />
</template>
<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { injectOrError } from "../../../shared/injectionKeys"
import { llmAnalysesConfiguratorKey } from "../../../shared/keys"
import { useSelectedPointStore } from "../../../shared/selectedPointStore"
import AnalysisDetailsDialog from "../../llmAnalysis/AnalysisDetailsDialog.vue"
import { LlmAnalysisState } from "../llmAnalysis/LlmAnalysisClient"
import { buildUrl } from "./InfoSidebar"

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
  if (hasAttemptedAutoOpen) return
  hasAttemptedAutoOpen = true
  if (dialogVisible.value) return
  const rawId = selectedPointStore.selectedAnalysisId
  const target = Array.isArray(rawId) ? rawId[0] : rawId
  if (target == null || target === "") return
  const numericTarget = Number(target)
  if (!Number.isFinite(numericTarget)) return
  const match = loaded.find((r) => r.id === numericTarget)
  if (match == null) return
  openAnalysis(match.id)
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
