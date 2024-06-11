import { Observable } from "rxjs"
import { ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../common/dataQuery"
import { useSettingsStore } from "../settingsStore"
import { FilterConfigurator } from "../../../configurators/filter"
import { refToObservable } from "../../../configurators/rxjs"

export class FlexibleZeroOnYAxisConfigurator implements DataQueryConfigurator, FilterConfigurator {
  private settingsStore = useSettingsStore()
  readonly value = ref(this.settingsStore.flexibleZeroOnYAxis)

  constructor() {
    watch(
      () => this.settingsStore.flexibleZeroOnYAxis,
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
