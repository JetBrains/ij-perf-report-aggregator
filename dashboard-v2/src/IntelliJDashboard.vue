<template>
  <el-form
    :inline="true"
    size="small"
  >
    <el-form-item label="Server">
      <el-input
        v-model="serverUrl"
        data-lpignore="true"
        placeholder="Enter the aggregated stats server URL..."
      />
    </el-form-item>

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
      <el-card
        shadow="never"
        :body-style="{ padding: '0px' }"
      >
        <div
          ref="chartElement"
          style="width: 100%; height: 300px;"
        />
      </el-card>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { ServerConfigurator } from "@/dataQuery"
import { defineComponent, onMounted, onUnmounted, ref, Ref, watch } from "vue"
import { PersistentStateManager } from "@/PersistentStateManager"
import { DimensionConfigurator, SubDimensionConfigurator } from "@/configurators/DimensionConfigurator"
import { MachineDimension } from "@/MachineDimension"
import DimensionSelect from "@/components/DimensionSelect.vue"
import TimeRangeSelect from "@/components/TimeRangeSelect.vue"
import MeasureSelect from "@/components/MeasureSelect.vue"
import DimensionHierarchicalSelect from "@/components/DimensionHierarchicalSelect.vue"
import { initEcharts } from "@/echarts"
import { DataQueryExecutor } from "@/DataQueryExecutor"
import { MeasureConfigurator } from "@/configurators/MeasureConfigurator"
import { ChartManager } from "@/ChartManager"
import { TimeRangeConfigurator } from "@/configurators/TimeRangeConfigurator"

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

initEcharts()

export default defineComponent({
  name: "IntelliJDashboard",
  components: {DimensionHierarchicalSelect, DimensionSelect, MeasureSelect, TimeRangeSelect},
  setup() {
    const persistentStateManager = new PersistentStateManager("ij");
    const serverConfigurator = new ServerConfigurator("ij", persistentStateManager)

    const productConfigurator = new DimensionConfigurator("product", serverConfigurator, persistentStateManager)

    watch(serverConfigurator.server, () => {
      productConfigurator.scheduleLoad()
      // we don't call dataQueryExecutor.scheduleExecution here, because it will be called on server info change
    })

    function getProjectName(value: string) {
      return projectNameToTitle.get(value) || value
    }

    const projectConfigurator = new SubDimensionConfigurator("project", productConfigurator, persistentStateManager, (a, b) => {
      const t1 = getProjectName(a)
      const t2 = getProjectName(b)
      if (t1.startsWith("simple ") && !t2.startsWith("simple ")) {
        return -1
      }
      if (t2.startsWith("simple ") && !t1.startsWith("simple ")) {
        return 1
      }
      return t1.localeCompare(t2);
    });

    const dataQueryExecutor = new DataQueryExecutor([
      serverConfigurator,
      productConfigurator,
      projectConfigurator,
    ])

    const measureConfigurator = new MeasureConfigurator(serverConfigurator, dataQueryExecutor, persistentStateManager)

    const machineConfigurator = new MachineDimension(
      new SubDimensionConfigurator("machine", productConfigurator, persistentStateManager),
      persistentStateManager,
    )

    persistentStateManager.init()

    function loadData() {
      productConfigurator.scheduleLoad()
      measureConfigurator.scheduleLoad()
      dataQueryExecutor.scheduleLoad()
    }

    productConfigurator.load()
    measureConfigurator.load()

    const chartElement: Ref<HTMLElement | null> = ref(null)

    let chartManager: ChartManager | null = null
    onMounted(() => {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartManager = new ChartManager(chartElement.value!, dataQueryExecutor)
    })
    onUnmounted(() => {
      let it = chartManager
      if (it == null) {
        return
      }

      chartManager = null
      it.dispose()
    })

    return {
      serverUrl: serverConfigurator.server,
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      measureConfigurator,
      timeRangeConfigurator: new TimeRangeConfigurator(dataQueryExecutor, persistentStateManager),
      getProjectName,
      loadData,
      chartElement,
    }
  },
})
</script>
