import { Observable } from "rxjs"
import { ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { useSettingsStore } from "../components/settings/settingsStore"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

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
