<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <div class="flex items-center">
          <MachineSelect :machine-configurator="machineConfigurator" />
        </div>
        <BranchSelect
          :branch-configurator="branchConfigurator1"
          :release-configurator="releaseConfigurator1"
          :triggered-by-configurator="triggeredByConfigurator1"
          :selection-limit="1"
        />
        <BranchSelect
          :branch-configurator="branchConfigurator2"
          :release-configurator="releaseConfigurator2"
          :triggered-by-configurator="triggeredByConfigurator2"
          :selection-limit="1"
        />
        <DimensionSelect
          label="Mode"
          :dimension="testModeConfigurator"
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
        field="branch1"
        :header="branch1 ?? ''"
      >
        <template #body="slotProps">
          <div :class="getColorForBuild(slotProps.data.build1, slotProps.data.build2)">
            {{ slotProps.data.build1 }}
          </div>
        </template>
      </Column>
      <Column
        field="branch2"
        :header="branch2 ?? ''"
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
import { ReleaseNightlyConfigurator } from "../../../configurators/ReleaseNightlyConfigurator"
import { ServerWithCompressConfigurator } from "../../../configurators/ServerWithCompressConfigurator"
import { SimpleMeasureConfigurator } from "../../../configurators/SimpleMeasureConfigurator"
import { fromFetchWithRetryAndErrorHandling } from "../../../configurators/rxjs"
import { containerKey } from "../../../shared/keys"
import { MAIN_METRICS } from "../../../util/mainMetrics"
import MeasureSelect from "../../charts/MeasureSelect.vue"
import BranchSelect from "../BranchSelect.vue"
import MachineSelect from "../MachineSelect.vue"
import { PersistentStateManager } from "../PersistentStateManager"
import StickyToolbar from "../StickyToolbar.vue"
import { dbTypeStore } from "../../../shared/dbTypes"
import { DBType } from "../sideBar/InfoSidebar"
import { modeSelectLabelFormat } from "../../../shared/labels"
import DimensionSelect from "../../charts/DimensionSelect.vue"
import { createTestModeConfigurator } from "../../../configurators/TestModeConfigurator"
import { DimensionConfigurator } from "../../../configurators/DimensionConfigurator"
import { DataQuery } from "../dataQuery"
import type { FilterConfigurator } from "../../../configurators/filter"

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

const recentTimeFilter: FilterConfigurator = {
  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "generated_time", q: ">subtractMonths(now(),1)" })
    return true
  },
  createObservable() {
    return null
  },
}

const measureConfigurator = new SimpleMeasureConfigurator("metrics", persistentStateManager)
measureConfigurator.initData(metricsNames)
const testConfigurator = new SimpleMeasureConfigurator("tests", persistentStateManager)
testConfigurator.state.loading = false

const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [recentTimeFilter])

const branchConfigurator1 = createBranchConfigurator(serverConfigurator, persistentStateManager, [recentTimeFilter], "branch1")
const branchConfigurator2 = createBranchConfigurator(serverConfigurator, persistentStateManager, [recentTimeFilter], "branch2")

const branch1 = ref<string | null>(null)
const branch2 = ref<string | null>(null)

const triggeredByConfigurator1 = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator1, recentTimeFilter])
const triggeredByConfigurator2 = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator1, recentTimeFilter])

const releaseConfigurator1 = new ReleaseNightlyConfigurator(persistentStateManager)
const releaseConfigurator2 = new ReleaseNightlyConfigurator(persistentStateManager)

const testModeConfigurator = createTestModeConfigurator(
  serverConfigurator,
  persistentStateManager,
  [branchConfigurator1, machineConfigurator, triggeredByConfigurator1, triggeredByConfigurator2, recentTimeFilter],
  "mode",
  false
)

const tableData = ref<TableRow[]>()
const fetchedData = ref<TableRow[]>()
combineLatest([
  branchConfigurator1.createObservable(),
  branchConfigurator2.createObservable(),
  serverConfigurator.createObservable(),
  machineConfigurator.createObservable(),
  testModeConfigurator.createObservable(),
])
  .pipe(
    filter(() => {
      const branch1SelectedValue = branchConfigurator1.selected.value
      const branch2SelectedValue = branchConfigurator2.selected.value
      return branch1SelectedValue !== null && branch2SelectedValue !== null
    })
  )
  .subscribe(() => {
    const branch1SelectedValue = branchConfigurator1.selected.value
    const branch2SelectedValue = branchConfigurator2.selected.value

    branch1.value = Array.isArray(branch1SelectedValue) ? branch1SelectedValue[0] : branch1SelectedValue
    branch2.value = Array.isArray(branch2SelectedValue) ? branch2SelectedValue[0] : branch2SelectedValue

    compareBranches(machineConfigurator, branch1.value, branch2.value, metricsNames, testModeConfigurator).subscribe((results: ComparisonResult[]) => {
      const tests = new Set<string>()
      const table: TableRow[] = []
      for (const r of results) {
        if (
          (r.Median1 != 0 || r.Median2 != 0) && //don't add metrics that are zero
          !/.*_\d+(#.*)?$/.test(r.MeasureName) //don't add metrics like foo_1
        ) {
          tests.add(r.Project)
          table.push({ test: r.Project, metric: r.MeasureName, build1: r.Median1, build2: r.Median2, difference: r.Diff })
        }
      }
      testConfigurator.initData([...tests])
      fetchedData.value = table
      tableData.value = table
    })
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

class ComparisonResult {
  public constructor(
    readonly Project: string,
    readonly MeasureName: string,
    readonly Median1: number,
    readonly Median2: number,
    readonly Diff: number
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

function compareBranches(
  machineConfigurator: MachineConfigurator,
  branch1Value: string | null,
  branch2Value: string | null,
  metricNames: string[],
  testModeConfigurator: DimensionConfigurator
): Observable<ComparisonResult[]> {
  const params = {
    branch1: branch1Value,
    branch2: branch2Value,
    table: dbName + "." + table,
    measure_names: metricNames,
    machine: machineConfigurator.getMergedValue(),
    mode: testModeConfigurator.selected.value,
  }
  const compressedParams = serverConfigurator.compressString(JSON.stringify(params))
  return fromFetchWithRetryAndErrorHandling<ComparisonResult[]>(serverConfigurator.serverUrl + "/api/compareBranches/" + compressedParams)
}

function handleNavigateToTest(project: string, metric: string) {
  const currentRoute = router.currentRoute.value
  const parts = currentRoute.path.split("/")
  parts[parts.length - 1] = dbTypeStore().dbType == DBType.INTELLIJ_DEV ? "testsDev" : "tests"
  const testURL = parts.join("/")

  const queryParams = new URLSearchParams(currentRoute.query as Record<string, string>)
  queryParams.delete("branch")
  queryParams.set("project", project)
  queryParams.set("measure", metric)
  queryParams.delete("metrics")
  queryParams.delete("tests")

  window.open(router.resolve(`${testURL}?${queryParams.toString()}&branch=${branch1.value}&branch=${branch2.value}`).href, "_blank")
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
