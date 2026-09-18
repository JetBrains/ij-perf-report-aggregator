import { ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../common/dataQuery"
import { useSettingsStore } from "../settingsStore"
import { median } from "../../../shared/changeDetector/statistic"
import { FilterConfigurator } from "../../../configurators/filter"
import { refToObservable } from "../../../configurators/rxjs"

export class ScalingConfigurator implements DataQueryConfigurator, FilterConfigurator {
  private readonly settingsStore = useSettingsStore()
  readonly value = ref(this.settingsStore.scaling)

  constructor() {
    watch(
      () => this.settingsStore.scaling,
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

export function scaleToMedian(arr: number[], scale = getMedianScale(arr)): number[] {
  return scale === 1 ? arr : arr.map((value) => value * scale)
}

/** Shared by the values and their deviation band. */
export function getMedianScale(arr: number[]): number {
  if (arr.length === 0) {
    return 1
  }
  const medianValue = median(arr)
  return medianValue === 0 ? 1 : 50 / medianValue
}
