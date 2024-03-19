import { ToastServiceMethods } from "primevue/toastservice"
import { useToast } from "primevue/usetoast"
import { combineLatest, debounceTime, Subject } from "rxjs"
import { onBeforeUnmount, onMounted, Ref, watch } from "vue"
import { refToObservable } from "../../configurators/rxjs"
import { DataManager } from "../DataManager"
import { InputData } from "../data"
import { reportData } from "../state"

export interface ChartManager {
  render(data: DataManager): void

  dispose(): void
}

export class ChartComponent {
  chartManager: ChartManager | null = null
  private readonly toast: ToastServiceMethods

  private readonly subject = new Subject<void>()

  constructor(
    containerRef: Ref<HTMLElement | null>,
    private readonly createChartManager: (container: HTMLElement) => Promise<ChartManager>
  ) {
    const containerObservable = refToObservable(containerRef)

    onBeforeUnmount(() => {
      const chartManager = this.chartManager
      if (chartManager != null) {
        this.chartManager = null
        chartManager.dispose()
        console.log("DISPOSED!!")
      }
    })

    watch(reportData, () => {
      this.requestRender()
    })

    onMounted(() => {
      this.requestRender()
    })

    combineLatest([this.subject, containerObservable])
      .pipe(debounceTime(10))
      .subscribe((value) => {
        const container = value[1]
        if (container != null) {
          this.renderDataIfAvailable(container)
        }
      })

    this.toast = useToast()
  }

  requestRender() {
    this.subject.next()
  }

  private renderDataIfAvailable(container: HTMLElement): void {
    const data = reportData.value
    if (data.length === 0) {
      // do not re-render as empty - null value not expected to be set in valid cases
      return
    }

    const dataManager = new DataManager(JSON.parse(data) as InputData)
    const chartManager = this.chartManager
    if (chartManager == null) {
      this.createChartManager(container)
        .then((chartManager) => {
          this.chartManager = chartManager
          chartManager.render(dataManager)
        })
        .catch((error: unknown) => {
          console.error("Cannot create chart", error)
          this.toast.add({ severity: "error", summary: "Cannot create chart", detail: (error as Error).toString() })
        })
    } else {
      chartManager.render(dataManager)
    }
  }
}
