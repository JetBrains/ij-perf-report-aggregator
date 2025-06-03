import { defineStore } from "pinia"
import { useUrlSearchParams } from "@vueuse/core"
import { computed } from "vue"

export const pointParamName = "point"

export const useSelectedPointStore = defineStore("selectedPointStore", () => {
  const storedSelectedPoint = useUrlSearchParams("history")

  const selectedPoint = computed({
    get: (): string|undefined => storedSelectedPoint[pointParamName],
    set(value) {
      storedSelectedPoint[pointParamName] = value
    },
  })

  return { selectedPoint }
})
