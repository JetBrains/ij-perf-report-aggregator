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
        {{ formatMeasure(slotProps.data.baselineValue) }}
      </template>
    </Column>
    <Column
      field="currentValue"
      :header="currentColumnLabel"
      :sortable="true"
    >
      <template #body="slotProps">
        {{ formatMeasure(slotProps.data.currentValue) }}
      </template>
    </Column>
    <Column
      field="difference"
      :header="differenceColumnLabel"
      :sortable="true"
    >
      <template #body="slotProps">
        {{ formatDifference(slotProps.data.difference) }}
      </template>
    </Column>
  </DataTable>
</template>

<script setup lang="ts">
import { FilterMatchMode } from "primevue/api"
import { ColumnFilterModelType } from "primevue/column"

import { Observable } from "rxjs"
import { onMounted, onUnmounted, ref } from "vue"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { injectOrError } from "../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../shared/keys"
import { DataQueryExecutor } from "./DataQueryExecutor"
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

export interface TestComparisonTableEntry {
  test: string
  baselineValue: number
  currentValue: number
  difference: number
}

const resultData = ref<TestComparisonTableEntry[]>([])

const filters = ref({
  test: { value: null, matchMode: FilterMatchMode.CONTAINS },
})

const formatMeasure = getValueFormatterByMeasureName(props.measure)

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
      query.addField({ n: "measures", subName: "name" })
      query.addField({ n: "measures", subName: "value" })

      configuration.measures = [props.measure]
      query.addFilter({ f: "measures.name", v: props.measure })

      query.order = "project"

      return true
    }

    createObservable(): Observable<unknown> | null {
      return null
    }
  })(),
] as DataQueryConfigurator[])

function applyData(data: (string | number)[][][]) {
  const rawMeasuresByTestName = new Map<string, number>()

  // The `data` array consists of one result for each configured "project", i.e. one result for each test name. We can then take the last entry
  // from the value arrays of that result to get the most up-to-date measure value.
  for (const resultForSingleProject of data) {
    const testNames = resultForSingleProject[0]
    if (testNames.length === 0) continue

    const measureValues = resultForSingleProject[2]
    rawMeasuresByTestName.set(testNames.at(-1), measureValues.at(-1))
  }

  const tableData: TestComparisonTableEntry[] = []

  for (const testComparison of props.comparisons) {
    const baselineValue = rawMeasuresByTestName.get(testComparison.baselineTestName) as number
    const currentValue = rawMeasuresByTestName.get(testComparison.currentTestName) as number
    const difference = Number.isFinite(baselineValue) && Number.isFinite(currentValue) ? (baselineValue - currentValue) / currentValue : 0

    tableData.push({
      test: testComparison.label,
      baselineValue,
      currentValue,
      difference,
    })
  }

  resultData.value = tableData

  emit("update:resultData", resultData.value)
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -- Initialization and teardown
// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

let unsubscribe: (() => void) | null = null

onMounted(() => {
  // Ensure that the API query only requests results for projects/tests which should be displayed by this comparison table.
  projectConfigurator.selected.value = props.comparisons.flatMap((testComparison) => [testComparison.baselineTestName, testComparison.currentTestName])

  initializeTable()
})

onUnmounted(() => {
  unsubscribe?.()
})

function initializeTable() {
  unsubscribe = dataQueryExecutor.subscribe((data, _configuration, isLoading) => {
    if (isLoading || data == null) {
      return
    }
    applyData(data)
  })
}
</script>
