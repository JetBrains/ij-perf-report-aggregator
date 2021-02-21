<template>
  <el-form
    :inline="true"
    size="small"
  >
    <ServerSelect v-model="serverUrl" />

    <DimensionSelect
      label="Product"
      :dimension="productConfigurator"
    />
    <DimensionSelect
      label="Project"
      :value-label="getProjectName"
      :dimension="projectConfigurator"
    />

    <DimensionHierarchicalSelect
      label="Machine"
      :dimension="machineConfigurator"
    />

    <el-form-item>
      <el-button
        title="Updated automatically, but you can force data reloading"
        type="primary"
        icon="el-icon-refresh"
        @click="loadData"
      />
    </el-form-item>
  </el-form>

  <el-form
    :inline="true"
    size="small"
  >
    <MeasureSelect :configurator="measureConfigurator" />
    <TimeRangeSelect :configurator="timeRangeConfigurator" />
  </el-form>

  <el-row>
    <el-col :span="12">
      <ChartCard :provider="dataQueryExecutor" />
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { ServerConfigurator } from "./dataQuery"
import { defineComponent, watch } from "vue"
import { PersistentStateManager } from "./PersistentStateManager"
import { DimensionConfigurator, SubDimensionConfigurator } from "./configurators/DimensionConfigurator"
import { MachineConfigurator } from "./configurators/MachineConfigurator"
import DimensionSelect from "./components/DimensionSelect.vue"
import TimeRangeSelect from "./components/TimeRangeSelect.vue"
import MeasureSelect from "./components/MeasureSelect.vue"
import ServerSelect from "./components/ServerSelect.vue"
import DimensionHierarchicalSelect from "./components/DimensionHierarchicalSelect.vue"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { MeasureConfigurator } from "./configurators/MeasureConfigurator"
import { TimeRangeConfigurator } from "./configurators/TimeRangeConfigurator"
import { debounceSync } from "./util/debounce"
import ChartCard from "./components/ChartCard.vue"

export default defineComponent({
  name: "IntelliJDashboard",
  components: {ChartCard, ServerSelect, DimensionHierarchicalSelect, DimensionSelect, MeasureSelect, TimeRangeSelect},
  setup() {
    const persistentStateManager = new PersistentStateManager("ij")
    const serverConfigurator = new ServerConfigurator("ij", persistentStateManager)
    const productConfigurator = new DimensionConfigurator("product", serverConfigurator, persistentStateManager)

    const projectConfigurator = createProjectConfigurator(productConfigurator, persistentStateManager)

    const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager)

    const machineConfigurator = new MachineConfigurator(
      new SubDimensionConfigurator("machine", productConfigurator, persistentStateManager),
      persistentStateManager,
    )

    const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

    const dataQueryExecutor = new DataQueryExecutor([
      serverConfigurator,
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
      measureConfigurator,
    ])

    watch(serverConfigurator.server, debounceSync(() => {
      for (const configurator of dataQueryExecutor.configurators) {
        if (configurator.scheduleLoad != null) {
          configurator.scheduleLoad(false)
        }
      }
    }, 1_000))

    persistentStateManager.init()
    dataQueryExecutor.scheduleLoadIncludingConfigurators(true)

    return {
      serverUrl: serverConfigurator.server,
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      measureConfigurator,
      timeRangeConfigurator,
      dataQueryExecutor,
      getProjectName,
      loadData: function () {
        dataQueryExecutor.scheduleLoadIncludingConfigurators()
      },
    }
  },
})

const projectNameToTitle = new Map<string, string>()
// noinspection SpellCheckingInspection
projectNameToTitle.set("/q9N7EHxr8F1NHjbNQnpqb0Q0fs", "joda-time")
// noinspection SpellCheckingInspection
projectNameToTitle.set("73YWaW9bytiPDGuKvwNIYMK5CKI", "simple for IJ")
// noinspection SpellCheckingInspection
projectNameToTitle.set("1PbxeQ044EEghMOG9hNEFee05kM", "light edit (IJ)")

// noinspection SpellCheckingInspection
projectNameToTitle.set("j1a8nhKJexyL/zyuOXJ5CFOHYzU", "simple for PS")
// noinspection SpellCheckingInspection
projectNameToTitle.set("JeNLJFVa04IA+Wasc+Hjj3z64R0", "simple for WS")
// noinspection SpellCheckingInspection
projectNameToTitle.set("nC4MRRFMVYUSQLNIvPgDt+B3JqA", "Idea")
Object.seal(projectNameToTitle)

function getProjectName(value: string) {
  return projectNameToTitle.get(value) || value
}

function createProjectConfigurator(productConfigurator: DimensionConfigurator, persistentStateManager: PersistentStateManager) {
  return new SubDimensionConfigurator("project", productConfigurator, persistentStateManager, (a, b) => {
    const t1 = getProjectName(a)
    const t2 = getProjectName(b)
    if (t1.startsWith("simple ") && !t2.startsWith("simple ")) {
      return -1
    }
    if (t2.startsWith("simple ") && !t1.startsWith("simple ")) {
      return 1
    }
    return t1.localeCompare(t2)
  })
}
</script>
