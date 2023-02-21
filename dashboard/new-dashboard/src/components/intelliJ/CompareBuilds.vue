<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect
          :ranges="TimeRangeConfigurator.timeRanges"
          :value="timeRangeConfigurator.value.value"
          :on-change="onChangeRange"
        >
          <template #icon>
            <CalendarIcon class="w-4 h-4 text-gray-500" />
          </template>
        </TimeRangeSelect>
        <BranchSelect
          :branch-configurator="branchConfigurator"
          :release-configurator="releaseConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        >
          <template #icon>
            <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionHierarchicalSelect>
        <DimensionSelect
          label="Build N1"
          :dimension="firstBuildConfigurator"
        />
        <DimensionSelect
          label="Build N2"
          :dimension="secondBuildConfigurator"
        />
      </template>
    </Toolbar>

    <DataTable
      :value="metricData"
      responsive-layout="scroll"
      show-gridlines
    >
      <Column
        field="test"
        header="Test"
      />
      <Column
        field="metric"
        header="Metric"
      />
      <Column
        field="build1"
        header="Value"
      >
        <template #body="slotProps">
          <div :class="getColorForBuild(slotProps.data.build1, slotProps.data.build2)">
            {{ slotProps.data.build1 }}
          </div>
        </template>
      </Column>
      <Column
        field="build2"
        header="Value"
      >
        <template #body="slotProps">
          <div :class="getColorForBuild(slotProps.data.build2, slotProps.data.build1)">
            {{ slotProps.data.build2 }}
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
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
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "shared/src/dataQuery"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey } from "../../shared/keys"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"


const dbName = "perfint"
const dbTable = "idea"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()

provide(containerKey, container)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager(
  `${dbName}-${dbTable}-dashboard`,
  {
    machine: initialMachine,
    branch: "master",
    project: [],
    measure: [],
  }, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])

const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator, branchConfigurator],
)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistentStateManager,
  [branchConfigurator, timeRangeConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistentStateManager)


function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}

const metricData = ref()
const firstBuildConfigurator = buildConfigurator("firstBuild", serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator, machineConfigurator])
const secondBuildConfigurator = buildConfigurator("secondBuild", serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator, machineConfigurator])
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
          && (r1.value/r2.value > 1.2 || r1.value/r2.value < 0.8) //don't add metrics that differ less than 20%
        ) {
          table.push({test: r1.test, metric: r1.metric, build1: r1.value, build2: r2.value})
        }
      }
      metricData.value = table
    },
  )
})

class Result {
  public constructor(readonly test: string, readonly metric: string, readonly value: number) {}
}

function getColorForBuild(build1: number, build2: number) {
  return [
    {
      "higher": build1 > build2,
      "lower": build1 < build2,
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
    }]).subscribe((data, _configuration) => {
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