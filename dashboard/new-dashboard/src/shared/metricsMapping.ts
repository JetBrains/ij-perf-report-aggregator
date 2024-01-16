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

export function measureNameToLabel(key: string): string {
  key = replaceKeys(key)
  if (key.startsWith("metrics.")) {
    key = key.slice(8)
  }
  return key.includes(".") ? key : /* remove _d or _i suffix */ key.replaceAll(/_[a-z]$/g, "")
}
