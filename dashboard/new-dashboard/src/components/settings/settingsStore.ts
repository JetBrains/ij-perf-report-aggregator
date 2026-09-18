import { useStorage } from "@vueuse/core"
import { defineStore } from "pinia"
import { computed } from "vue"

export const useSettingsStore = defineStore("settingsStore", () => {
  const storedScaling = useStorage("scalingEnabled", false)
  const storedSmoothing = useStorage("smoothingEnabled", false)

  // scaling and smoothing are mutually exclusive, the rest of the settings are plain flags
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

  return {
    scaling,
    smoothing,
    detectChanges: useStorage("detectChangesEnabled", false),
    flexibleYZero: useStorage("floatingNull", false),
    removeOutliers: useStorage("removeOutliers", false),
    groupBranches: useStorage("groupBranches", true),
    fadeOnHover: useStorage("fadeOnHover", false),
    stdDevInterval: useStorage("stdDevInterval", true),
  }
})
