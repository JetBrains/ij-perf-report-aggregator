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

const measureConfigurator = new SimpleMeasureConfigurator("metrics", persistentStateManager)
measureConfigurator.initData(metricsNames)
const testConfigurator = new SimpleMeasureConfigurator("tests", persistentStateManager)

const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager)

const branchConfigurator1 = createBranchConfigurator(serverConfigurator, persistentStateManager, [], "branch1")
const branchConfigurator2 = createBranchConfigurator(serverConfigurator, persistentStateManager, [], "branch2")

const branch1 = ref<string | null>(null)
const branch2 = ref<string | null>(null)

const triggeredByConfigurator1 = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator1])
const triggeredByConfigurator2 = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator1])

const releaseConfigurator1 = new ReleaseNightlyConfigurator(persistentStateManager)
const releaseConfigurator2 = new ReleaseNightlyConfigurator(persistentStateManager)

const tableData = ref<TableRow[]>()
const fetchedData = ref<TableRow[]>()
combineLatest([branchConfigurator1.createObservable(), branchConfigurator2.createObservable(), serverConfigurator.createObservable(), machineConfigurator.createObservable()])
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

    combineLatest([getAllMetricsFromBranch(machineConfigurator, branch1.value, metricsNames), getAllMetricsFromBranch(machineConfigurator, branch2.value, metricsNames)]).subscribe(
      (data: Result[][]) => {
        const firstBranchResults = data[0]
        const secondBranchResults = data[1]
        const tests = new Set<string>()
        const table: TableRow[] = []
        for (const r1 of firstBranchResults) {
          const r2 = secondBranchResults.find((value) => {
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

function getAllMetricsFromBranch(machineConfigurator: MachineConfigurator, branch: string | null, metricNames: string[]): Observable<Result[]> {
  const params = {
    branch,
    table: dbName + "." + table,
    measure_names: metricNames,
    machine: machineConfigurator.getMergedValue(),
  }
  const compressedParams = serverConfigurator.compressString(JSON.stringify(params))
  return fromFetchWithRetryAndErrorHandling<Result[]>(serverConfigurator.serverUrl + "/api/compareBranches/" + compressedParams)
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

  window.open(router.resolve(testURL + "?" + queryParams.toString() + `&branch=${branch1.value}` + `&branch=${branch2.value}`).href, "_blank")
}
</script>

<style>
.p-button-icon-only.p-button-rounded {
  @apply text-black;
}
.lower {
  font-weight: 700;
  color: #ff5252;
}

.higher {
  font-weight: 700;
  color: #66bb6a;
}
</style>
