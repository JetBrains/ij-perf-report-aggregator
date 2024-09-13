import { useStorage } from "@vueuse/core"
import { defineStore } from "pinia"
import { computed } from "vue"

export const useSettingsStore = defineStore("settingsStore", () => {
  const storedScaling = useStorage("scalingEnabled", false)
  const storedSmoothing = useStorage("smoothingEnabled", false)
  const storedDetectChanges = useStorage("detectChangesEnabled", false)
  const storedFlexibleYZero = useStorage("floatingNull", false)
  const storedRemoveOutliers = useStorage("removeOutliers", false)
  const storedShowAllMetrics = useStorage("showAllMetrics", false)

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

  const flexibleYZero = computed({
    get: () => storedFlexibleYZero.value,
    set(value) {
      storedFlexibleYZero.value = value
    },
  })

  const removeOutliers = computed({
    get: () => storedRemoveOutliers.value,
    set(value) {
      storedRemoveOutliers.value = value
    },
  })

  const showAllMetrics = computed({
    get: () => storedShowAllMetrics.value,
    set(value) {
      storedShowAllMetrics.value = value
    },
  })

  return { scaling, smoothing, detectChanges, flexibleYZero, removeOutliers, showAllMetrics }
})
