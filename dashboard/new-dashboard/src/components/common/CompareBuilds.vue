<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <div class="flex items-center">
          <DimensionHierarchicalSelect
            label="Machine"
            :dimension="machineConfigurator"
          >
            <template #icon>
              <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
            </template>
          </DimensionHierarchicalSelect>
          <span class="ml-5">
            Build 1
          </span>
        </div>
        <BranchSelect
          :branch-configurator="branchConfigurator1"
          :release-configurator="releaseConfigurator1"
          :triggered-by-configurator="triggeredByConfigurator1"
        />
        <DimensionSelect
          label="Build N1"
          :dimension="firstBuildConfigurator"
        />
        Build 2
        <BranchSelect
          :branch-configurator="branchConfigurator2"
          :release-configurator="releaseConfigurator2"
          :triggered-by-configurator="triggeredByConfigurator2"
        />
        <DimensionSelect
          label="Build N2"
          :dimension="secondBuildConfigurator"
        />
      </template>
    </Toolbar>

    <DataTable
      v-model:filters="filters"
      :value="metricData"
      responsive-layout="scroll"
      show-gridlines
      filter-display="row"
      class="p-datatable-sm"
    >
      <Column
        field="test"
        header="Test"
        :sortable="true"
      >
        <template #filter="{filterModel}">
          <InputText
            v-model="filterModel.value"
            type="text"
            class="p-column-filter"
            placeholder="Search by name"
          />
        </template>
      </Column>
      <Column
        field="metric"
        header="Metric"
        :sortable="true"
        :filter-match-mode-options="metricsMatchModes"
      >
        <template #filter="{filterModel}">
          <InputText
            v-if="filterModel.matchMode == FilterMatchMode.CONTAINS"
            v-model="filterModel.value"
            type="text"
            class="p-column-filter"
            placeholder="Search by name"
          />
          <div
            v-else
            class="p-column-filter"
          >
            {{ metricsMatchModes.find(e=>e.value == filterModel.matchMode)?.label }}
          </div>
        </template>
      </Column>
      <Column
        field="build1"
        header="Build 1"
      >
        <template #body="slotProps">
          <div :class="getColorForBuild(slotProps.data.build1, slotProps.data.build2)">
            {{ slotProps.data.build1 }}
          </div>
        </template>
      </Column>
      <Column
        field="build2"
        header="Build 2"
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
        :filter-match-mode-options="differenceMatchModes"
      >
        <template #filter="{filterModel}">
          <Slider
            v-model="filterModel.value"
            class="m-3"
          />
          <div class="flex px-2">
            <span class="text-sm">Difference ≥ {{ filterModel.value ? filterModel.value : 0 }}%</span>
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { FilterMatchMode, FilterService } from "primevue/api"
import { combineLatest, Observable } from "rxjs"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { buildConfigurator } from "shared/src/configurators/BuildConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "shared/src/dataQuery"
import {MAIN_METRICS} from "shared/src/util/mainMetrics"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey } from "../../shared/keys"
import BranchSelect from "./BranchSelect.vue"

const props = defineProps<{
  dbName: string
  table: string
}>()

const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()

provide(containerKey, container)

const serverConfigurator = new ServerConfigurator(props.dbName, props.table)
const persistentStateManager = new PersistentStateManager(
  `${props.dbName}-${props.table}-dashboard`,
  {
    machine: initialMachine,
    branch: "master",
    project: [],
    measure: [],
  }, router)

const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistentStateManager,
)

const branchConfigurator1 = createBranchConfigurator(serverConfigurator, persistentStateManager)
const branchConfigurator2 = createBranchConfigurator(serverConfigurator, persistentStateManager)

const triggeredByConfigurator1 = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator1])
const triggeredByConfigurator2 = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator1])

const releaseConfigurator1 = new ReleaseNightlyConfigurator(persistentStateManager)
const releaseConfigurator2 = new ReleaseNightlyConfigurator(persistentStateManager)


const metricData = ref()
const firstBuildConfigurator = buildConfigurator("firstBuild", serverConfigurator, persistentStateManager, [branchConfigurator1, machineConfigurator])
const secondBuildConfigurator = buildConfigurator("secondBuild", serverConfigurator, persistentStateManager, [branchConfigurator2, machineConfigurator])
combineLatest([refToObservable(firstBuildConfigurator.selected), refToObservable(secondBuildConfigurator.selected)]).subscribe(data => {
  combineLatest([getAllMetricsFromBuild(data[0] as string),
    getAllMetricsFromBuild(data[1] as string)]).subscribe(data => {
      const firstBuildResults = data[0]
      const secondBuildResults = data[1]
      const table = []
      for (const r1 of firstBuildResults) {
        const r2 = secondBuildResults.find(value => {
          return value.test == r1.test && value.metric == r1.metric
        })
        if (r2 != undefined
          && (r1.value != 0 || r2.value != 0) //don't add metrics that are zero
          && !/.*_\d+(#.*)?$/.test(r1.metric) //don't add metrics like foo_1
          && (r1.value != r2.value) //don't add equal metrics
        ) {
          const difference = (((r2.value - r1.value) / r1.value) * 100).toFixed(1)
          table.push({test: r1.test, metric: r1.metric, build1: r1.value, build2: r2.value, difference})
        }
      }
      metricData.value = table
    },
  )
})

FilterService.register("metricsFilter", value => {
  return MAIN_METRICS.has(value as string)
})
const indexingMetrics = new Set(["indexing", "scanning", "numberOfIndexingRuns"])
FilterService.register("indexingFilter", value => {
  return indexingMetrics.has(value as string)
})
const memoryMetrics = new Set(["freedMemoryByGC", "fullGCPause", "gcPause", "gcPauseCount", "totalHeapUsedMax", "Memory | IDE | RESIDENT SIZE (MB) 95th pctl"])
FilterService.register("memoryFilter", value => {
  return memoryMetrics.has(value as string)
})

const metricsMatchModes = [
  {label: "Main metrics", value: "metricsFilter"},
  {label: "Indexing metrics", value: "indexingFilter"},
  {label: "Memory metrics", value: "memoryFilter"},
  {label: "Contains", value: FilterMatchMode.CONTAINS}
]

FilterService.register("differenceFilter", (a, b) => {
  return a > b || a < -b
})
const differenceMatchModes = [
  {label: "Set difference", value: "differenceFilter"},
]

const filters = ref({
  "test": {value: null, matchMode: FilterMatchMode.CONTAINS},
  "metric": {value: null, matchMode: "metricsFilter"},
  "difference": {value: 30, matchMode: "differenceFilter"},
})

class Result {
  public constructor(readonly test: string, readonly metric: string, readonly value: number) {}
}

function getColorForBuild(build1: number, build2: number) {
  return [
    {
      "higher": build1 < build2,
      "lower": build1 > build2,
    },
  ]
}

function getAllMetricsFromBuild(build: string): Observable<Array<Result>> {
  return new Observable<Array<Result>>(subscriber => {
    new DataQueryExecutor([serverConfigurator, new class implements DataQueryConfigurator {
      configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
        configuration.queryProducers.push(new SimpleQueryProducer())
        query.addField("project")
        query.addField({n: "measures", subName: "name"})
        query.addField({n: "measures", subName: "value"})
        const buildParts = build.split(".")
        query.addFilter({f: "build_c1", v: buildParts[0]})
        if (buildParts[1] != undefined) {
          query.addFilter({f: "build_c2", v: buildParts[1]})
        }
        if (buildParts[2] != undefined) {
          query.addFilter({f: "build_c3", v: buildParts[2]})
        }
        query.order = "project"
        return true
      }

      createObservable(): Observable<unknown> | null {
        return null
      }
    }]).subscribe((data, _configuration, isLoading) => {
      if (isLoading || data == null){
        return
      }
      const result: Array<Result> = new Array<Result>()
      const datum = data[0]
      for (let i = 0; i < datum[0].length; i++) {
        result.push(new Result(datum[0][i] as string, datum[1][i] as string, datum[2][i] as number))
      }
      subscriber.next(result)
    })
  })
}

</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}

.lower {
  font-weight: 700;
  color: #FF5252;
}

.higher {
  font-weight: 700;
  color: #66BB6A;
}
</style>