<template>
  <div class="flex flex-col gap-4">
    <div
      v-if="candidateMode == null"
      class="text-sm text-gray-500 py-3 px-5 border border-solid rounded-md"
    >
      Select the mode <span class="font-semibold">types2</span> (LEGACY_AND_NEW, compared with <span class="font-semibold">default</span>) or
      <span class="font-semibold">new_only</span> (NEW_ONLY, compared with <span class="font-semibold">annotator_only</span>) to see the comparison.
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
        :candidate="candidate"
        :base="base"
      />
      <EngineCompareHeatmap
        :rows="rows"
        :selected-metric-types="selectedMetricTypes"
        :mode="mode"
        :candidate="candidate"
        :base="base"
      />
      <EngineCompareBars
        :rows="rows"
        :candidate="candidate"
        :base="base"
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
import { computed, inject, ref, toRef } from "vue"
import { selectedToArray } from "../../../configurators/DimensionConfigurator"
import { defaultModeName } from "../../../configurators/TestModeConfigurator"
import { testModeConfiguratorKey } from "../../../shared/keys"
import { baseModeOf, engineOfMode, Quantity } from "./highlightingMetrics"
import { EngineComparisonMode, useEngineComparison } from "./useEngineComparison"
import EngineCompareSummary from "./EngineCompareSummary.vue"
import EngineCompareHeatmap from "./EngineCompareHeatmap.vue"
import EngineCompareBars from "./EngineCompareBars.vue"
import RunTimeline from "./RunTimeline.vue"

const { selectedMetricTypes, quantity } = defineProps<{
  selectedMetricTypes: string[]
  quantity: Quantity
}>()

// The first mode of the dashboard mode switch that has a base mode is the candidate; the comparison selects its base mode.
const modeConfigurator = inject(testModeConfiguratorKey, null)
const candidateMode = computed<string | null>(() => {
  const modes = selectedToArray(modeConfigurator?.selected.value)
  return modes.find((mode) => baseModeOf(mode) != null) ?? null
})
const candidate = computed(() => engineOfMode(candidateMode.value ?? ""))
const baseMode = computed(() => baseModeOf(candidateMode.value ?? "") ?? defaultModeName)
const base = computed(() => engineOfMode(baseMode.value))

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
  candidateMode: toRef(() => candidateMode.value ?? defaultModeName),
  baseMode,
  mode,
  selectedRun,
})

const currentRunLabel = computed(() => {
  if (mode.value !== "single") return ""
  if (selectedRun.value == null) return runs.value[0]?.label ?? "latest"
  return runs.value.find((run) => run.t === selectedRun.value)?.label ?? ""
})
</script>
