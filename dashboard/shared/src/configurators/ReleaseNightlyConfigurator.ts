import { ref, ShallowRef } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter } from "../dataQuery"
import { DimensionConfigurator } from "./DimensionConfigurator"

const nightly = "Nightly"
const eap = "EAP / Release"

type ReleaseType = typeof eap | typeof nightly

export class ReleaseNightlyConfigurator extends DimensionConfigurator {
  declare readonly selected: ShallowRef<ReleaseType | Array<ReleaseType> | null>
  readonly values: ShallowRef<Array<ReleaseType>> = ref<Array<ReleaseType>>([eap, nightly])

  constructor(persistentStateManager: PersistentStateManager | null) {
    super("releaseConfigurator", true)
    this.state.disabled = false
    this.selected.value = nightly
    persistentStateManager?.add("releaseConfigurator", this.selected)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return false
    }
    let filter = getFilter(value)

    if (Array.isArray(value) && value.length > 1) {
      configuration.queryProducers.push({
          size(): number {
            return value.length
          },
          mutate(index: number) {
            if (filter == null) {
              filter = getFilter(value[index])
              query.addFilter(filter)
            }
            else {
              filter.o = value[index] == eap ? "!=" : "="
            }
          },
          getSeriesName(index: number): string {
            return value[index]
          },
          getMeasureName(_index: number): string {
            return configuration.measures[0]
          },
        },
      )
    }
    if (filter != null) {
      query.addFilter(filter)
    }
    return true
  }
}
function getFilter(value: ReleaseType): DataQueryFilter
function getFilter(value: ReleaseType | Array<ReleaseType> | null): DataQueryFilter | null
function getFilter(value: ReleaseType | Array<ReleaseType> | null): DataQueryFilter | null {
  if (value == null || value.length === 0) {
    return null
  }
  const filter: DataQueryFilter = {f: "build_c3", v: 0}
  if (Array.isArray(value)) {
    if (value.includes(eap) && value.includes(nightly)) {
      return null
    }
    filter.o = value.includes(eap) ? "!=" : "="
  }
  else {
    filter.o = value == eap ? "!=" : "="
  }
  return filter
}
