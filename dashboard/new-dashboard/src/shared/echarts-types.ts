// Type definitions for echarts that are not exported from the public API
// These types are defined based on the internal echarts types to maintain compatibility
// when using bundler module resolution

export type OptionDataValue = string | number | Date | null | undefined

export type ScaleDataValue = string | number | Date

export type OptionDataItem = OptionDataValue | Record<string, OptionDataValue> | OptionDataValue[]

export type OptionSourceData = OptionDataItem[] | Record<string, OptionDataValue[]> | OptionDataValue[][]

export interface DimensionDefinition {
  name?: string
  type?: "number" | "ordinal" | "float" | "int" | "time"
  displayName?: string
}
