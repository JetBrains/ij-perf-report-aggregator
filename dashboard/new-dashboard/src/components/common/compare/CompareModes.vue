<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <BranchSelect
          :branch-configurator="branchConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <MachineSelect :machine-configurator="machineConfigurator" />
        <DimensionSelect
          label="Mode"
          :dimension="testModeConfigurator1"
          :selected-label="modeSelectLabelFormat"
        >
          <template #icon>
            <AdjustmentsVerticalIcon class="w-4 h-4" />
          </template>
        </DimensionSelect>
        <DimensionSelect
          label="Mode"
          :dimension="testModeConfigurator2"
          :selected-label="modeSelectLabelFormat"
        >
          <template #icon>
            <AdjustmentsVerticalIcon class="w-4 h-4" />
          </template>
        </DimensionSelect>
        <MeasureSelect
          :configurator="testConfigurator"
          title="Test"
        >
          <template #icon>
            <ChartBarIcon class="w-4 h-4" />
          </template>
        </MeasureSelect>
        <MeasureSelect :configurator="measureConfigurator">
          <template #icon>
            <BeakerIcon class="w-4 h-4" />
          </template>
        </MeasureSelect>
      </template>
    </StickyToolbar>

    <DataTable
      :value="tableData"
      show-gridlines
      class="p-datatable-sm"
      sort-field="difference"
      :sort-order="-1"
    >
      <Column
        field="test"
        header="Test"
        :sortable="true"
      >
        <template #body="slotProps">
          <div class="flex items-center">
            <div>{{ slotProps.data.test }}</div>
            <div class="ml-2">
              <Button
                icon="pi pi-external-link"
                class="p-button-rounded p-button-text p-button-sm"
                @click="() => handleNavigateToTest(slotProps.data.test, slotProps.data.metric)"
              />
            </div>
          </div>
        </template>
      </Column>
      <Column
        field="metric"
        header="Metric"
        :sortable="true"
      />
      <Column
        field="mode1"
        :header="mode1 ?? ''"
      >
        <template #body="slotProps">
          <div :class="getColorForBuild(slotProps.data.build1, slotProps.data.build2)">
            {{ slotProps.data.build1 }}
          </div>
        </template>
      </Column>
      <Column
        field="mode2"
        :header="mode2 ?? ''"
      >
        <template #body="slotProps">
          <div :class="getColorForBuild(slotProps.data.build2, slotProps.data.build1)">
            {{ slotProps.data.build2 }}
          </div>
        </template>
      </Column>
      <Column
        field="difference"
        header="Difference (%)"
        :sortable="true"
      />
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { combineLatest, filter, Observable } from "rxjs"
import { provide, ref, watch, useTemplateRef } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../../configurators/BranchConfigurator"
import { MachineConfigurator } from "../../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../../configurators/ServerWithCompressConfigurator"
import { SimpleMeasureConfigurator } from "../../../configurators/SimpleMeasureConfigurator"
import { fromFetchWithRetryAndErrorHandling } from "../../../configurators/rxjs"
import { containerKey } from "../../../shared/keys"
import { MAIN_METRICS } from "../../../util/mainMetrics"
import MeasureSelect from "../../charts/MeasureSelect.vue"
import MachineSelect from "../MachineSelect.vue"
import { PersistentStateManager } from "../PersistentStateManager"
import StickyToolbar from "../StickyToolbar.vue"
import { dbTypeStore } from "../../../shared/dbTypes"
import { DBType } from "../sideBar/InfoSidebar"
import DimensionSelect from "../../charts/DimensionSelect.vue"
import { modeSelectLabelFormat } from "../../../shared/labels"
import { createTestModeConfigurator } from "../../../configurators/TestModeConfigurator"
import BranchSelect from "../BranchSelect.vue"
import { TimeRangeConfigurator } from "../../../configurators/TimeRangeConfigurator"

interface CompareBranchesProps {
  dbName: string
  table: string
  metricsNames?: string[]
}

interface TableRow {
  test: string
  metric: string
  build1: number
  build2: number
  difference: number
}

const { dbName, table, metricsNames = MAIN_METRICS } = defineProps<CompareBranchesProps>()

const initialMachine = "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
const container = useTemplateRef<HTMLElement>("container")
const router = useRouter()

provide(containerKey, container)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, table)
const persistentStateManager = new PersistentStateManager(
  `${dbName}-${table}-compare-branches`,
  {
    machine: initialMachine,
    branch: "master",
    project: [],
    measure: [],
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])

const measureConfigurator = new SimpleMeasureConfigurator("metrics", persistentStateManager)
measureConfigurator.initData(metricsNames)
const testConfigurator = new SimpleMeasureConfigurator("tests", persistentStateManager)
testConfigurator.state.loading = false

const testModeConfigurator1 = createTestModeConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator, branchConfigurator, machineConfigurator, triggeredByConfigurator],
  "mode1",
  false
)
const testModeConfigurator2 = createTestModeConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator, branchConfigurator, machineConfigurator, triggeredByConfigurator],
  "mode2",
  false
)

const mode1 = ref<string | null>(null)
const mode2 = ref<string | null>(null)

const tableData = ref<TableRow[]>()
const fetchedData = ref<TableRow[]>()
combineLatest([testModeConfigurator1.createObservable(), testModeConfigurator2.createObservable(), serverConfigurator.createObservable(), machineConfigurator.createObservable()])
  .pipe(
    filter(() => {
      const mode1SelectedValue = testModeConfigurator1.selected.value
      const mode2SelectedValue = testModeConfigurator2.selected.value
      return mode1SelectedValue !== null && mode2SelectedValue !== null
    })
  )
  .subscribe(() => {
    const mode1SelectedValue = testModeConfigurator1.selected.value
    const mode2SelectedValue = testModeConfigurator2.selected.value

    mode1.value = Array.isArray(mode1SelectedValue) ? mode1SelectedValue[0] : mode1SelectedValue
    mode2.value = Array.isArray(mode2SelectedValue) ? mode2SelectedValue[0] : mode2SelectedValue

    combineLatest([getAllMetricsFromMode(machineConfigurator, mode1.value, metricsNames), getAllMetricsFromMode(machineConfigurator, mode2.value, metricsNames)]).subscribe(
      (data: Result[][]) => {
        const firstModeResults = data[0]
        const secondModeResults = data[1]
        const tests = new Set<string>()
        const table: TableRow[] = []
        for (const r1 of firstModeResults) {
          const r2 = secondModeResults.find((value) => {
            return value.Project == r1.Project && value.MeasureName == r1.MeasureName
          })
          if (
            r2 != undefined &&
            (r1.Median != 0 || r2.Median != 0) && //don't add metrics that are zero
            !/.*_\d+(#.*)?$/.test(r1.MeasureName) //don't add metrics like foo_1
          ) {
            const difference = Number((((r2.Median - r1.Median) / r1.Median) * 100).toFixed(1))
            tests.add(r1.Project)
            table.push({ test: r1.Project, metric: r1.MeasureName, build1: r1.Median, build2: r2.Median, difference })
          }
        }
        testConfigurator.initData([...tests])
        fetchedData.value = table
        tableData.value = table
      }
    )
  })

watch(
  [measureConfigurator.selected, testConfigurator.selected],
  ([metrics, tests]) => {
    tableData.value = fetchedData.value?.filter((value) => {
      return metrics?.includes(value.metric)
    })
    if (tests != null) {
      tableData.value = tableData.value?.filter((value) => {
        return tests.includes(value.test)
      })
    }
  },
  { immediate: true }
)

class Result {
  public constructor(
    readonly Project: string,
    readonly MeasureName: string,
    readonly Median: number
  ) {}
}

function getColorForBuild(build1: number, build2: number) {
  return [
    {
      higher: build1 < build2,
      lower: build1 > build2,
    },
  ]
}

function getAllMetricsFromMode(machineConfigurator: MachineConfigurator, mode: string | null, metricNames: string[]): Observable<Result[]> {
  if (mode == "default") {
    mode = ""
  }
  const params = {
    mode,
    branch: branchConfigurator.selected.value,
    table: dbName + "." + table,
    measure_names: metricNames,
    machine: machineConfigurator.getMergedValue(),
  }
  const compressedParams = serverConfigurator.compressString(JSON.stringify(params))
  return fromFetchWithRetryAndErrorHandling<Result[]>(serverConfigurator.serverUrl + "/api/compareModes/" + compressedParams)
}

function handleNavigateToTest(project: string, metric: string) {
  const currentRoute = router.currentRoute.value
  const parts = currentRoute.path.split("/")
  parts[parts.length - 1] = dbTypeStore().dbType == DBType.INTELLIJ_DEV ? "testsDev" : "tests"
  const testURL = parts.join("/")

  const queryParams = new URLSearchParams(currentRoute.query as Record<string, string>)
  queryParams.set("project", project)
  queryParams.set("measure", metric)
  queryParams.delete("metrics")
  queryParams.delete("tests")
  queryParams.delete("mode1")
  queryParams.delete("mode2")

  window.open(router.resolve(`${testURL}?${queryParams.toString()}&mode=${mode1.value}&mode=${mode2.value}`).href, "_blank")
}
</script>

<style>
.lower {
  font-weight: 700;
  color: #ff5252;
}

.higher {
  font-weight: 700;
  color: #66bb6a;
}
</style>
