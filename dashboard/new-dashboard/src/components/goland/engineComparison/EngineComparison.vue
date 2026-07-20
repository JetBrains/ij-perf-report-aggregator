<template>
  <div class="flex flex-col gap-4">
    <div
      v-if="!bothEnginesSelected"
      class="text-sm text-gray-500 py-3 px-5 border border-solid rounded-md"
    >
      Select both <span class="font-semibold">LEGACY</span> and <span class="font-semibold">NEW_ONLY</span> engines to see the comparison.
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
          v-model="selectedRunDay"
          :runs="runDays"
        />
      </div>
      <EngineCompareSummary
        :aggregates="aggregates"
        :mode="mode"
        :run-label="currentRunLabel"
      />
      <EngineCompareHeatmap
        :rows="rows"
        :selected-metric-types="selectedMetricTypes"
        :mode="mode"
      />
      <EngineCompareBars :rows="rows" />
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
import { Quantity } from "./highlightingMetrics"
import { EngineComparisonMode, useEngineComparison } from "./useEngineComparison"
import EngineCompareSummary from "./EngineCompareSummary.vue"
import EngineCompareHeatmap from "./EngineCompareHeatmap.vue"
import EngineCompareBars from "./EngineCompareBars.vue"
import RunTimeline from "./RunTimeline.vue"

const { selectedMetricTypes, quantity, bothEnginesSelected } = defineProps<{
  selectedMetricTypes: string[]
  quantity: Quantity
  // The comparison needs both engines; the Engine toolbar selector gates it (and the drill-down charts).
  bothEnginesSelected: boolean
}>()

// Single run is the default view: it shows one build's exact before/after, so a real improvement is
// visible rather than smeared by the across-runs median. Toggle to Aggregated for the SPEC-style verdict.
const mode = ref<EngineComparisonMode>("single")
// null follows the latest run; a day key pins a specific run (set via the timeline).
const selectedRunDay = ref<number | null>(null)

const modeOptions = [
  { label: "Aggregated (all runs)", value: "aggregated" },
  { label: "Single run", value: "single" },
]

const { rows, aggregates, runDays, loading } = useEngineComparison({
  selectedMetricTypes: toRef(() => selectedMetricTypes),
  quantity: toRef(() => quantity),
  mode,
  selectedRunDay,
})

const currentRunLabel = computed(() => {
  if (mode.value !== "single") return ""
  if (selectedRunDay.value == null) return runDays.value[0]?.label ?? "latest"
  return runDays.value.find((run) => run.day === selectedRunDay.value)?.label ?? ""
})
</script>
