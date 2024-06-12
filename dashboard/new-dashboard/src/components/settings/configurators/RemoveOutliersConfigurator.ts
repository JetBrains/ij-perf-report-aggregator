import { Observable } from "rxjs"
import { ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../common/dataQuery"
import { useSettingsStore } from "../settingsStore"
import { FilterConfigurator } from "../../../configurators/filter"
import { refToObservable } from "../../../configurators/rxjs"
import { rollingMad } from "../../../shared/changeDetector/statistic"

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

export function removeOutliers(data: (string | number)[][], windowSize: number = 5, threshold: number = 3): (string | number)[][] {
  if (data.length === 0) {
    return data
  }

  const values = data[1] as number[]

  const { medians, mads } = rollingMad(values, windowSize)
  const skippedIndex: number[] = []

  for (const [i, value] of values.entries()) {
    const madScore = Math.abs(value - medians[i]) / mads[i]
    if (madScore > threshold) {
      skippedIndex.push(i)
    }
  }

  return skippedIndex.length === 0 ? data : data.map((row) => row.filter((_, index) => !skippedIndex.includes(index)))
}
