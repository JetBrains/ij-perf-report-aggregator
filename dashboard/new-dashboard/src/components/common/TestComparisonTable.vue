<!--
-- A table component which allows comparing pairs of tests *from the same build* based on some measurement. The table will always base its
-- comparisons on the latest available data point.
--
-- For example, the Kotlin K1 vs. K2 comparison dashboard uses this component to compare the time it takes K1 and K2 to complete the same test.
-->

<template>
  <DataTable
    v-model:filters="filters"
    :value="resultData"
    show-gridlines
    class="p-datatable-sm"
  >
    <Column
      field="test"
      header="Test"
      :sortable="true"
    >
      <template #filter="{ filterModel }">
        <InputText
          v-model="(filterModel as ColumnFilterModelType).value"
          type="text"
          class="p-column-filter"
          placeholder="Search by name"
        />
      </template>
    </Column>
    <Column
      field="baselineValue"
      :header="baselineColumnLabel"
      :sortable="true"
    >
      <template #body="slotProps">
        {{ formatMeasureOrFallback(slotProps.data.baselineValue) }}
      </template>
    </Column>
    <Column
      field="currentValue"
      :header="currentColumnLabel"
      :sortable="true"
    >
      <template #body="slotProps">
        {{ formatMeasureOrFallback(slotProps.data.currentValue) }}
      </template>
    </Column>
    <Column
      field="difference"
      :header="differenceColumnLabel"
      :sortable="true"
    >
      <template #body="slotProps">
        {{ formatDifferenceOrFallback(slotProps.data.difference) }}
      </template>
    </Column>
  </DataTable>
</template>

<script setup lang="ts">
import { FilterMatchMode } from "primevue/api"
import { ColumnFilterModelType } from "primevue/column"

import { Observable } from "rxjs"
import { onMounted, onUnmounted, ref, watch } from "vue"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { injectOrError } from "../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../shared/keys"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { TestComparisonTableEntry } from "./TestComparisonTableEntry"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "./dataQuery"
import { formatPercentage, getValueFormatterByMeasureName } from "./formatter"

/**
 * Defines that a `baseline` test should be compared against a `current` test. This represents a single row in the comparison table.
 */
export interface TestComparison {
  /**
   * Some label shown in the table as the test name.
   */
  label: string

  /**
   * The full test name of the baseline test, e.g. `kotlin_empty/completion/empty_place_with_library_cache_k1`.
   */
  baselineTestName: string

  /**
   * The full test name of the current test, e.g. `kotlin_empty/completion/empty_place_with_library_cache_k2`.
   */
  currentTestName: string
}

interface TestComparisonTableProps {
  measure: string
  comparisons: TestComparison[]

  configurators: (DataQueryConfigurator | FilterConfigurator)[]

  baselineColumnLabel?: string
  currentColumnLabel?: string
  differenceColumnLabel?: string

  formatDifference?: (difference: number) => string
}

const props = withDefaults(defineProps<TestComparisonTableProps>(), {
  baselineColumnLabel: "Baseline",
  currentColumnLabel: "Current",
  differenceColumnLabel: "Difference (%)",
  formatDifference: formatPercentage,
})

const emit = defineEmits<(e: "update:resultData", resultData: TestComparisonTableEntry[]) => void>()

const resultData = ref<TestComparisonTableEntry[]>([])

watch(resultData, () => {
  emit("update:resultData", resultData.value)
})

const filters = ref({
  test: { value: null, matchMode: FilterMatchMode.CONTAINS },
})

const formatMeasure = getValueFormatterByMeasureName(props.measure)

function formatMeasureOrFallback(value: number | null) {
  if (value === null) return "N/A"
  return formatMeasure(value)
}

function formatDifferenceOrFallback(value: number | null) {
  if (value === null) return "N/A"
  return props.formatDifference(value)
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -- Query configuration and data evaluation
// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

const serverConfigurator = injectOrError(serverConfiguratorKey)

const projectConfigurator = dimensionConfigurator("project", serverConfigurator, null, true, [...(props.configurators as FilterConfigurator[])])

const dataQueryExecutor = new DataQueryExecutor([
  serverConfigurator,
  ...props.configurators,
  projectConfigurator,
  new (class implements DataQueryConfigurator {
    configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
      query.addField("project")
      query.addField("generated_time")
      query.addField({ n: "measures", subName: "name" })
      query.addField({ n: "measures", subName: "value" })

      configuration.measures = [props.measure]
      query.addFilter({ f: "measures.name", v: props.measure })

      query.order = ["project", "generated_time"]

      return true
    }

    createObservable(): Observable<unknown> | null {
      return null
    }
  })(),
] as DataQueryConfigurator[])

function applyData(data: (string | number)[][][]) {
  const rawMeasuresByTestName = new Map<string, number | null>()

  // The `data` array consists of one result for each configured "project", i.e. one result for each test name. We can then take the last entry
  // from the value arrays of that result to get the most up-to-date measure value.
  for (const resultForSingleProject of data) {
    const testNames = resultForSingleProject[0] as string[]
    if (testNames.length === 0) continue

    const measureValues = resultForSingleProject[3] as number[]
    rawMeasuresByTestName.set(testNames.at(-1) ?? "", measureValues.at(-1) ?? null)
  }

  const tableData: TestComparisonTableEntry[] = []

  for (const testComparison of props.comparisons) {
    const baselineValue = (rawMeasuresByTestName.get(testComparison.baselineTestName) as number | null) ?? null // Replace `undefined` with `null`.
    const currentValue = (rawMeasuresByTestName.get(testComparison.currentTestName) as number | null) ?? null

    let difference: number | null = null
    if (baselineValue !== null && currentValue !== null) {
      difference = Number.isFinite(baselineValue) && Number.isFinite(currentValue) ? (baselineValue - currentValue) / currentValue : 0
    }

    tableData.push({
      test: testComparison.label,
      baselineValue,
      currentValue,
      difference,
    })
  }

  resultData.value = tableData
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -- Initialization, updating and teardown
// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

let unsubscribe: (() => void) | null = null

onMounted(() => {
  updateProjectConfigurator(props.comparisons)
  initializeTable()
})

watch(
  () => props.comparisons,
  (newValue) => {
    updateProjectConfigurator(newValue)
  }
)

onUnmounted(() => {
  unsubscribe?.()
})

function updateProjectConfigurator(comparisons: TestComparison[]) {
  // Ensure that the API query only requests results for projects/tests which should be displayed by this comparison table.
  projectConfigurator.selected.value = comparisons.flatMap((testComparison) => [testComparison.baselineTestName, testComparison.currentTestName])

  // If there are no comparisons and thus no projects, the query will not be performed, so we need to clear out the old result data manually.
  if (comparisons.length === 0) {
    resultData.value = []
  }
}

function initializeTable() {
  unsubscribe = dataQueryExecutor.subscribe((data, _configuration, isLoading) => {
    if (isLoading || data == null) {
      return
    }
    applyData(data)
  })
}
</script>
