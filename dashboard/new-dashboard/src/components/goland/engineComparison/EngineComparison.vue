<template>
  <div class="flex flex-col gap-4">
    <div
      v-if="!bothEnginesSelected"
      class="text-sm text-gray-500 py-3 px-5 border border-solid rounded-md"
    >
      Select both <span class="font-semibold">LEGACY</span> and <span class="font-semibold">{{ candidate.id }}</span> engines to see the comparison.
    </div>
    <div
      v-else-if="selectedMetricTypes.length === 0"
      class="text-sm text-gray-500 py-3 px-5 border border-solid rounded-md"
    >
      Select at least one metric type to compare.
    </div>
    <template v-else>
      <div class="flex flex-col gap-2">
        <SelectButton
          v-model="mode"
          :options="modeOptions"
          option-label="label"
          option-value="value"
          :allow-empty="false"
          class="self-start"
        />
        <RunTimeline
          v-if="mode === 'single'"
          v-model="selectedRun"
          :runs="runs"
        />
      </div>
      <EngineCompareSummary
        :aggregates="aggregates"
        :mode="mode"
        :run-label="currentRunLabel"
        :quantity="quantity"
        :candidate="candidate.id"
      />
      <EngineCompareHeatmap
        :rows="rows"
        :selected-metric-types="selectedMetricTypes"
        :mode="mode"
        :candidate="candidate.id"
      />
      <EngineCompareBars
        :rows="rows"
        :candidate="candidate.id"
      />
      <div
        v-if="loading && rows.length === 0"
        class="text-sm text-gray-500"
      >
        Loading comparison…
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRef } from "vue"
import { Engine, Quantity } from "./highlightingMetrics"
import { EngineComparisonMode, useEngineComparison } from "./useEngineComparison"
import EngineCompareSummary from "./EngineCompareSummary.vue"
import EngineCompareHeatmap from "./EngineCompareHeatmap.vue"
import EngineCompareBars from "./EngineCompareBars.vue"
import RunTimeline from "./RunTimeline.vue"

const { selectedMetricTypes, quantity, candidate, bothEnginesSelected } = defineProps<{
  selectedMetricTypes: string[]
  quantity: Quantity
  // The engine compared with LEGACY; the `new` runs of the comparison are its runs.
  candidate: Engine
  // The comparison needs both engines; the Engine toolbar selector gates it (and the drill-down charts).
  bothEnginesSelected: boolean
}>()

// Single run is the default view: it shows one build's exact before/after, so a real improvement is
// visible rather than smeared by the across-runs median. Toggle to Aggregated for the SPEC-style verdict.
const mode = ref<EngineComparisonMode>("single")
// null follows the latest run; a build timestamp pins a specific run (set via the timeline).
const selectedRun = ref<number | null>(null)

const modeOptions = [
  { label: "Aggregated (all runs)", value: "aggregated" },
  { label: "Single run", value: "single" },
]

const { rows, aggregates, runs, loading } = useEngineComparison({
  selectedMetricTypes: toRef(() => selectedMetricTypes),
  quantity: toRef(() => quantity),
  candidateSuffix: toRef(() => candidate.suffix),
  mode,
  selectedRun,
})

const currentRunLabel = computed(() => {
  if (mode.value !== "single") return ""
  if (selectedRun.value == null) return runs.value[0]?.label ?? "latest"
  return runs.value.find((run) => run.t === selectedRun.value)?.label ?? ""
})
</script>
