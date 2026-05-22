import { defineStore } from "pinia"
import { UrlParams, useUrlSearchParams } from "@vueuse/core"
import { computed } from "vue"
import type { DefaultLabelFormatterCallbackParams as CallbackDataParams } from "echarts"

export const pointParamName = "point"
export const analysisParamName = "analysis"

export const useSelectedPointStore = defineStore("selectedPointStore", () => {
  const storedSelectedPoint: UrlParams = useUrlSearchParams("history")

  const selectedPoint = computed({
    get: (): string | string[] | undefined => storedSelectedPoint[pointParamName],
    set(value: string | string[] | undefined) {
      if (value != undefined) {
        storedSelectedPoint[pointParamName] = value
      }
    },
  })

  const selectedAnalysisId = computed({
    get: (): string | string[] | undefined => storedSelectedPoint[analysisParamName],
    set(value: string | string[] | undefined) {
      if (value != undefined) {
        storedSelectedPoint[analysisParamName] = value
      }
    },
  })

  return { selectedPoint, selectedAnalysisId }
})

// Render-time bridge: MeasureConfigurator's itemStyle.color callback already
// detects the URL-selected point per data item and has the real ECharts params
// in scope. It pushes those here so the sidebar auto-open consumer can reuse
// them instead of rescanning the dataset and synthesizing fake params.
const matchedSelectedPointParams: CallbackDataParams[] = []

export function captureMatchedSelectedPoint(params: CallbackDataParams): void {
  matchedSelectedPointParams.push(params)
}

export function consumeMatchedSelectedPoints(): CallbackDataParams[] {
  const result = [...matchedSelectedPointParams]
  matchedSelectedPointParams.length = 0
  return result
}
