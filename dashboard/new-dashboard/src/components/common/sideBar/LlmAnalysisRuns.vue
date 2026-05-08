<template>
  <Accordion
    v-if="visibleRuns.length > 0"
    value="0"
    class="llm-analysis-runs"
  >
    <AccordionPanel value="0">
      <AccordionHeader>LLM analyses</AccordionHeader>
      <AccordionContent>
        <ul class="gap-1.5 break-all">
          <li
            v-for="run in visibleRuns"
            :key="run.id"
            class="flex items-center justify-between gap-1.5"
          >
            <span
              v-tooltip.top="run.createdAt"
              class="text-xs"
            >{{ formatCreatedAt(run.createdAt) }}</span>
            <span class="flex gap-1.5 items-center">
              <a
                v-if="run.runBuildId"
                :href="buildUrl(Number(run.runBuildId))"
                target="_blank"
                :class="getURLStyle()"
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
</template>
<script setup lang="ts">
import { computed } from "vue"
import { computedAsync } from "@vueuse/core"
import { injectOrNull } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"
import { LlmAnalysisClient, LlmAnalysisState } from "../llmAnalysis/LlmAnalysisClient"
import { buildUrl, InfoData } from "./InfoSidebar"

const { data, runsRefreshTrigger } = defineProps<{
  data: InfoData | null
  runsRefreshTrigger: number
}>()

const serverConfigurator = injectOrNull(serverConfiguratorKey)

const llmAnalysisRuns = computedAsync(async () => {
  const metric = data?.series[0]?.metricName
  if (data == null || serverConfigurator == null || data.buildIdPrevious == null || !metric) return []
  void runsRefreshTrigger
  return await new LlmAnalysisClient(serverConfigurator).getLlmAnalysisRuns({
    date: data.date,
    project: data.projectName,
    metric,
    currentBuildId: String(data.buildId),
    prevBuildId: String(data.buildIdPrevious),
  })
}, [])

const visibleRuns = computed(() => llmAnalysisRuns.value.filter((r) => r.state !== LlmAnalysisState.NotStarted))

function stateIconClass(state: LlmAnalysisState): string {
  switch (state) {
    case LlmAnalysisState.Queued:
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
    case LlmAnalysisState.Queued:
      return "Queued"
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

function getURLStyle() {
  return "underline decoration-dotted hover:no-underline"
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
