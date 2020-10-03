// Copyright 2000-2019 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file.
import { Component, Vue, Watch } from "vue-property-decorator"
import { getOrCreateReportVisualizerSettings } from "@/state/state"
import { DataManager } from "@/state/DataManager"
import { ChartManager } from "@/charts/ChartManager"
import { Notification } from "element-ui"
import { InputData } from "@/state/data"

// @ts-ignore
@Component
export abstract class BaseChartComponent<T extends ChartManager> extends Vue {
  protected chartManager!: T | null

  /** @final */
  get measurementData(): string | null {
    return getOrCreateReportVisualizerSettings(this.$store).data
  }

  created() {
    this.chartManager = null
  }

  mounted() {
    this.renderDataIfAvailable()
  }

  protected abstract createChartManager(): Promise<T>

  @Watch("measurementData")
  /** @final */
  protected renderDataIfAvailable(): void {
    const data = this.measurementData
    if (data == null || data.length === 0) {
      // do not re-render as empty - null value not expected to be set in valid cases
      return
    }

    const dataManager = new DataManager(JSON.parse(data) as InputData)

    let chartManager = this.chartManager
    if (chartManager == null) {
      this.createChartManager()
        .then(chartManager => {
          this.chartManager = chartManager
          chartManager.render(dataManager)
        })
        .catch(e => {
          console.log(e)
          Notification.error(e)
        })
    }
    else {
      chartManager.render(dataManager)
    }
  }

  beforeDestroy() {
    const chartManager = this.chartManager
    if (chartManager != null) {
      this.chartManager = null
      chartManager.dispose()
    }
  }
}
