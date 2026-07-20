<template>
  <div class="flex flex-col gap-3 py-3 px-5 border border-solid rounded-md">
    <div class="text-sm font-semibold">Movers, ranked by Δ% <span class="font-normal text-gray-500">(worst regression at top; red = slower, green = faster; hover for values)</span></div>
    <div
      v-if="phaseGroups.length === 0"
      class="text-sm text-gray-500 py-2"
    >
      No comparable scenarios to rank yet.
    </div>
    <div
      v-else
      class="flex flex-col gap-6"
    >
      <EngineCompareBarsChart
        v-for="group in phaseGroups"
        :key="group.phase"
        :title="group.label"
        :rows="group.rows"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { EngineCompareRow } from "./useEngineComparison"
import { PHASES, phaseLabel } from "./highlightingMetrics"
import EngineCompareBarsChart from "./EngineCompareBarsChart.vue"

const { rows } = defineProps<{ rows: EngineCompareRow[] }>()

// One ranked chart per phase (Cold / Warm / Typing), in the canonical phase order, dropping phases with
// no comparable rows.
const phaseGroups = computed(() =>
  PHASES.map((phase) => ({
    phase,
    label: phaseLabel(phase),
    rows: rows.filter((row) => row.phase === phase && Number.isFinite(row.diffPercent)),
  })).filter((group) => group.rows.length > 0)
)
</script>
