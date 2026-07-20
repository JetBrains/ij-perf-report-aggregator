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
  "bsp.used.at.exit.mb": "bsp.used.after.sync.mb",
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
