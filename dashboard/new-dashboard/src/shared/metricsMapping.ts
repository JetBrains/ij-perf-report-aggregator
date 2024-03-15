import { CallbackDataParams } from "echarts/types/src/util/types"

export const METRICS_MAPPING: Record<string, string> = {
  "processingTime#": "",
  "processingSpeedAvg#": "",
  "lexingSpeed#": "",
  "parsingSpeed#": "",
  "lexingTime#": "",
  "parsingTime#": "",
  "lexingSize#": "",
  "parsingSize#": "",
  indexingTimeWithoutPauses: "indexing",
  scanningTimeWithoutPauses: "scanning",
}

function replaceKeys(originalKey: string): string {
  let modifiedKey = originalKey
  for (const [searchValue, replaceValue] of Object.entries(METRICS_MAPPING)) {
    modifiedKey = modifiedKey.replaceAll(searchValue, replaceValue)
  }
  return modifiedKey
}

export function measureNameToLabel(key: string | undefined): string {
  if (key == undefined) return ""
  key = replaceKeys(key)
  if (key.startsWith("metrics.")) {
    key = key.slice(8)
  }
  return key.includes(".") ? key : /* remove _d or _i suffix */ key.replaceAll(/_[a-z]$/g, "")
}

function getPrefix(name: string): string {
  const lastDot = name.lastIndexOf(".")
  const lastSlash = name.lastIndexOf("/")
  const lastIndex = Math.max(lastDot, lastSlash)
  return name.slice(0, Math.max(0, lastIndex))
}

export function getCommonPrefix(params: CallbackDataParams[]): string {
  const names = params.map((param) => getPrefix(param.seriesName as string))
  if (names.length === 0) {
    return ""
  }
  let commonPrefix = names[0]
  for (let i = 1; i < names.length; i++) {
    while (!names[i].startsWith(commonPrefix) && commonPrefix) {
      const lastIndex = Math.max(commonPrefix.lastIndexOf("/"), commonPrefix.lastIndexOf("."))
      commonPrefix = commonPrefix.slice(0, Math.max(0, lastIndex))
    }
  }
  if (commonPrefix.endsWith("/") || commonPrefix.endsWith(".")) {
    commonPrefix = commonPrefix.slice(0, Math.max(0, commonPrefix.length - 1))
  }
  return commonPrefix
}

export function removePrefix(name: string, prefix: string): string {
  if (name.startsWith(prefix + ".") || name.startsWith(prefix + "/")) {
    return name.slice(Math.max(0, prefix.length + 1))
  }
  return name
}
