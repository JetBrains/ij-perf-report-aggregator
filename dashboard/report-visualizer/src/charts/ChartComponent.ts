import { debounceSync } from "shared/src/util/debounce"
import { watch, onBeforeUnmount, onMounted } from "vue"
import { DataManager } from "../DataManager"
import { InputData } from "../data"
import { reportData } from "../state"

export interface ChartManager {
  render(data: DataManager): void

  dispose(): void
}

export class ChartComponent {
  chartManager: ChartManager | null = null

  private readonly renderDataIfAvailableDebounced = debounceSync(() => this.renderDataIfAvailable(), 10)

  constructor(private readonly createChartManager: () => Promise<ChartManager>) {
    onBeforeUnmount(() => {
      const chartManager = this.chartManager
      if (chartManager != null) {
        this.chartManager = null
        chartManager.dispose()
      }
    })

    onMounted(() => {
      this.renderDataIfAvailableDebounced()
    })

    watch(reportData, () => {
      this.renderDataIfAvailableDebounced()
    })
  }

  renderDataIfAvailable(): void {
    const data = reportData.value
    if (data == null || data.length === 0) {
      // do not re-render as empty - null value not expected to be set in valid cases
      return
    }

    const dataManager = new DataManager(JSON.parse(data) as InputData)
    const chartManager = this.chartManager
    if (chartManager == null) {
      this.createChartManager()
        .then(chartManager => {
          this.chartManager = chartManager
          chartManager.render(dataManager)
        })
        .catch(e => {
          console.error("Cannot create chart", e)
          ElNotification({
            type: "error",
            message: (e as Error).toString(),
          })
        })
    }
    else {
      chartManager.render(dataManager)
    }
  }
}
