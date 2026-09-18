import { ref, Ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { useSettingsStore } from "./settingsStore"
import { FilterConfigurator } from "../../configurators/filter"
import { refToObservable } from "../../configurators/rxjs"

type SettingsStore = ReturnType<typeof useSettingsStore>

/** Names of the boolean flags of the settings store. `-?` keeps optional store members from widening the union with `undefined`. */
export type BooleanSetting = { [K in keyof SettingsStore]-?: SettingsStore[K] extends boolean ? K : never }[keyof SettingsStore]

/**
 * Republishes a boolean flag of the settings store as an observable, so that flipping the flag re-runs the query.
 * The flag itself is applied when the data is post-processed, hence query and filter are left untouched.
 */
export class SettingConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly value: Ref<boolean>

  constructor(setting: BooleanSetting) {
    const settingsStore = useSettingsStore()
    this.value = ref(settingsStore[setting])
    watch(
      () => settingsStore[setting],
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
