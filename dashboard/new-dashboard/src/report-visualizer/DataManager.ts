import { compare } from "compare-versions"
import { numberFormat } from "../components/common/formatter"
import { InputData, InputDataV20, ItemV20, UnitConverter } from "./data"

const markerNames = ["app initialized callback", "module loading"]
export const markerNameToRangeTitle = new Map<string, string>([
  ["app initialized callback", "app initialized"],
  ["module loading", "project initialized"],
])

export type GroupedItems = { category: string; items: ItemV20[] }[]

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

  private _markerItems: (ItemV20 | null)[] | null = null

  get isUnifiedItems(): boolean {
    const version = this.version
    return version != null && compare(version, "32", ">=")
  }

  get items(): ItemV20[] {
    return this.data.items
  }

  // start, duration in microseconds
  getServiceItems(): GroupedItems {
    const version = this.version
    const isNewCompactFormat = version != null && compare(version, "20", ">=")
    if (isNewCompactFormat) {
      const data = this.data as InputDataV20
      return [
        { category: "app components", items: data.appComponents ?? [] },
        { category: "project components", items: data.projectComponents ?? [] },
        { category: "module components", items: data.moduleComponents ?? [] },

        { category: "app services", items: data.appServices ?? [] },
        { category: "project services", items: data.projectServices ?? [] },
        { category: "module services", items: data.moduleServices ?? [] },
      ]
    } else if (version != null && compare(version, "12", ">=")) {
      throw new Error(`Report version ${version} is not supported, ask if needed`)
    } else {
      throw new Error(`Report version ${version} is not supported, ask if needed`)
    }
  }

  get markerItems(): (ItemV20 | null)[] {
    if (this._markerItems != null) {
      return this._markerItems
    }

    const items = this.data.items
    if (items.length === 0) {
      return []
    }

    const result = new Array<ItemV20 | null>(markerNames.length)
    // JS array is sparse and setting length doesn't pre-fill array
    result.fill(null)
    itemLoop: for (const item of items) {
      for (const [i, markerName] of markerNames.entries()) {
        if (result[i] == null && item.n === markerName) {
          result[i] = item

          // stop if all items are found
          if (result.every((it) => it != null)) {
            break itemLoop
          }
        }
      }
    }

    for (const [i, markerName] of markerNames.entries()) {
      if (result[i] == null) {
        console.warn(`Cannot find item for phase "${markerName}"`)
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
  return lastDotIndex < 0 ? qualifiedName : qualifiedName.slice(lastDotIndex + 1)
}
