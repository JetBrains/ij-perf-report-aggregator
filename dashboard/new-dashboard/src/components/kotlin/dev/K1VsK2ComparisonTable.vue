<template>
  <section class="flex flex-col w-full mt-8">
    <div class="flex flex-row gap-6 mb-8">
      <div>
        <h3 class="text-2xl mb-3">{{ name }}</h3>
        <p class="text-sm text-gray-600">Measure: {{ measure }}</p>
      </div>
      <div class="flex-grow"></div>
      <div
        v-for="topStat in topStats"
        :key="topStat.label"
      >
        <div class="py-3 px-5 border border-solid rounded-md border-zinc-200">
          <h3 class="m-0 text-sm mb-2">{{ topStat.label }}</h3>
          <span class="text-2xl font-bold">{{ topStat.value }}</span>
        </div>
      </div>
    </div>

    <TestComparisonTable
      :measure="measure"
      :comparisons="projects.map(transformToTestComparison)"
      :configurators="configurators"
      baseline-column-label="K1"
      current-column-label="K2"
      difference-column-label="Improvement (%)"
      @update:result-data="(newValue) => (resultData = newValue)"
    />

    <p class="text-sm text-gray-500 text-right mt-4">The table only displays the results of the last build from the selected branch.</p>
  </section>
</template>

<script setup lang="ts">
import { computed, Ref, ref } from "vue"
import { FilterConfigurator } from "../../../configurators/filter"
import TestComparisonTable, { TestComparison } from "../../common/TestComparisonTable.vue"
import { DataQueryConfigurator } from "../../common/dataQuery"
import { formatPercentage, getValueFormatterByMeasureName } from "../../common/formatter"
import { TestComparisonTableEntry } from "../../common/TestComparisonTableEntry"

interface Props {
  name: string
  measure: string
  projects: string[]

  configurators: (DataQueryConfigurator | FilterConfigurator)[]
}

const props = defineProps<Props>()

const resultData = ref<TestComparisonTableEntry[]>([]) as Ref<TestComparisonTableEntry[]>

// eslint-disable-next-line @typescript-eslint/no-unsafe-return,@typescript-eslint/no-unsafe-member-access
const totalTimeK1 = computed(() => calculateTotalTime((entry) => entry.baselineValue))

// eslint-disable-next-line @typescript-eslint/no-unsafe-return,@typescript-eslint/no-unsafe-member-access
const totalTimeK2 = computed(() => calculateTotalTime((entry) => entry.currentValue))

const totalImprovement = computed(() => (totalTimeK2.value == 0 ? 0 : (totalTimeK1.value - totalTimeK2.value) / totalTimeK2.value))

function calculateTotalTime(getTime: (entry: TestComparisonTableEntry) => number): number {
  return resultData.value
    .map((entry) => {
      const value = getTime(entry)
      return Number.isFinite(value) ? value : 0
    })
    .reduce((a, b) => a + b, 0)
}

const topStats = computed(() => [
  {
    label: "K1 / Total Time",
    value: formatTime(totalTimeK1.value),
  },
  {
    label: "K2 / Total Time",
    value: formatTime(totalTimeK2.value),
  },
  {
    label: "Total Improvement",
    value: formatPercentage(totalImprovement.value),
  },
])

const formatTime = getValueFormatterByMeasureName(props.measure)

function transformToTestComparison(projectName: string): TestComparison {
  // We want to compare K1 and K2 tests against each other, and they are respectively suffixed with "_k1" and "_k2".
  return {
    label: projectName,
    baselineTestName: `${projectName}_k1`,
    currentTestName: `${projectName}_k2`,
  }
}
</script>
