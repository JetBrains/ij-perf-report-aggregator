import compareVersions from "compare-versions"
import { numberFormat } from "shared/src/formatter"
import { InputData, InputDataV20, ItemV0, ItemV20, UnitConverter } from "./data"

const markerNames = ["app initialized callback", "module loading"]
export const markerNameToRangeTitle = new Map<string, string>([["app initialized callback", "app initialized"], ["module loading", "project initialized"]])

export type GroupedItems = Array<{ category: string; items: Array<ItemV20> }>

export interface DataDescriptor {
  readonly unitConverter: UnitConverter
  readonly shortenName?: boolean
  readonly threshold?: number

  // used only for time line chart
  rowIndexThreshold?: number
}

export function formatDuration(value: number, dataDescriptor: DataDescriptor): string {
  return numberFormat.format(dataDescriptor.unitConverter.convert(value))
}

export class DataManager {
  private readonly version: string | null

  constructor(readonly data: InputData) {
    this.version = data.version
  }

  private _markerItems: Array<ItemV0 | null> | null = null

  // start, duration in microseconds
  getServiceItems(): GroupedItems {
    const version = this.version
    const isNewCompactFormat = version != null && compareVersions.compare(version, "20", ">=")
    if (isNewCompactFormat) {
      const data = this.data as InputDataV20
      return [
        {category: "app components", items: data.appComponents ?? []},
        {category: "project components", items: data.projectComponents ?? []},
        {category: "module components", items: data.moduleComponents ?? []},

        {category: "app services", items: data.appServices ?? []},
        {category: "project services", items: data.projectServices ?? []},
        {category: "module services", items: data.moduleServices ?? []},
      ]
    }
    else if (version != null && compareVersions.compare(version, "12", ">=")) {
      throw new Error(`Report version ${version} is not supported, ask if needed`)
      // this._serviceEvents = this.data.traceEvents.filter(value => value.cat != null && serviceEventCategorySet.has(value.cat)) as Array<CompleteTraceEvent>
      // return this._serviceEvents.map(it => {
      //   return {
      //     n: it.name,
      //     d: it.dur,
      //     t: it.tid,
      //     s: it.ts - it.dur,
      //     p: "",
      //   }
      // })
    }
    else {
      throw new Error(`Report version ${version} is not supported, ask if needed`)
      // const list: Array<CompleteTraceEvent> = []
      // const data = this.data as InputDataV11AndLess
      //
      // convertV11ToTraceEvent(data.appComponents, "appComponents", list)
      // convertV11ToTraceEvent(data.projectComponents, "projectComponents", list)
      // convertV11ToTraceEvent(data.moduleComponents, "moduleComponents", list)
      //
      // convertV11ToTraceEvent(data.appServices, "appServices", list)
      // convertV11ToTraceEvent(data.projectServices, "projectServices", list)
      // convertV11ToTraceEvent(data.moduleServices, "moduleServices", list)
      //
      // this._serviceEvents = list
      // return list
    }
  }

  get markerItems(): Array<ItemV0 | null> {
    if (this._markerItems != null) {
      return this._markerItems
    }

    const items = this.data == null ? null : this.data.items
    if (items == null || items.length === 0) {
      return []
    }

    const result = new Array<ItemV0 | null>(markerNames.length)
    // JS array is sparse and setting length doesn't pre-fill array
    result.fill(null)
    itemLoop: for (const item of items) {
      for (let i = 0; i < markerNames.length; i++) {
        if (result[i] == null && item.name === markerNames[i]) {
          result[i] = item

          // stop if all items are found
          if (result.every(it => it != null)) {
            break itemLoop
          }
        }
      }
    }

    for (let i = 0; i < markerNames.length; i++) {
      if (result[i] == null) {
        console.warn(`Cannot find item for phase "${markerNames[i]}"`)
      }
    }

    this._markerItems = result
    return result
  }

  // noinspection JSUnusedGlobalSymbols
  toJSON(_key: string): InputData {
    return this.data
  }
}

export function getShortName(qualifiedName: string): string {
  const lastDotIndex = qualifiedName.lastIndexOf(".")
  return lastDotIndex < 0 ? qualifiedName : qualifiedName.substring(lastDotIndex + 1)
}