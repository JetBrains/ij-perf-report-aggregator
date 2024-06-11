import { Observable } from "rxjs"
import { ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../common/dataQuery"
import { useSettingsStore } from "../settingsStore"
import { FilterConfigurator } from "../../../configurators/filter"
import { refToObservable } from "../../../configurators/rxjs"
import { calculatePercentile } from "../../../shared/changeDetector/statistic"

export class RemoveOutliersConfigurator implements DataQueryConfigurator, FilterConfigurator {
  private settingsStore = useSettingsStore()
  readonly value = ref(this.settingsStore.removeOutliers)

  constructor() {
    watch(
      () => this.settingsStore.removeOutliers,
      (newValue) => {
        this.value.value = newValue
      }
    )
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.value)
  }

  configureFilter(_: DataQuery): boolean {
    return true
  }

  configureQuery(_: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    return true
  }
}

export function removeOutliers(data: (string | number)[][]): (string | number)[][] {
  if (data.length === 0) {
    return data
  }
  const values = data[1] as number[]
  const p95th = calculatePercentile(values, 95)
  const p5th = calculatePercentile(values, 5)
  const skippedIndex: number[] = []
  for (const [index, value] of values.entries()) {
    if (value < p5th || value > p95th) {
      skippedIndex.push(index)
    }
  }

  return skippedIndex.length === 0 ? data : data.map((row) => row.filter((_, index) => !skippedIndex.includes(index)))
}
