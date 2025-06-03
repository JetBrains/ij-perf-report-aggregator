import { defineStore } from "pinia"
import { UrlParams, useUrlSearchParams } from "@vueuse/core"
import { computed } from "vue"

export const pointParamName = "point"

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

  return { selectedPoint }
})
