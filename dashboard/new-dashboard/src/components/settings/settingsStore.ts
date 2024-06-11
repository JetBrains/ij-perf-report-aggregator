import { useStorage } from "@vueuse/core"
import { defineStore } from "pinia"
import { computed } from "vue"

export const useSettingsStore = defineStore("settingsStore", () => {
  const storedScaling = useStorage("scalingEnabled", false)
  const storedSmoothing = useStorage("smoothingEnabled", false)
  const storedDetectChanges = useStorage("detectChangesEnabled", false)
  const storedFlexibleZeronOnYAxis = useStorage("floatingNull", false)

  const scaling = computed({
    get: () => storedScaling.value,
    set(value) {
      storedScaling.value = value
      storedSmoothing.value = false
    },
  })

  const smoothing = computed({
    get: () => storedSmoothing.value,
    set(value) {
      storedSmoothing.value = value
      storedScaling.value = false
    },
  })

  const detectChanges = computed({
    get: () => storedDetectChanges.value,
    set(value) {
      storedDetectChanges.value = value
    },
  })

  const flexibleZeroOnYAxis = computed({
    get: () => storedFlexibleZeronOnYAxis.value,
    set(value) {
      storedFlexibleZeronOnYAxis.value = value
    },
  })

  return { scaling, smoothing, detectChanges, flexibleZeroOnYAxis }
})
