import { useStorage } from "@vueuse/core"
import { defineStore } from "pinia"
import { computed } from "vue"

export const useSettingsStore = defineStore("settingsStore", () => {
  const storedScaling = useStorage("scalingEnabled", false)
  const storedSmoothing = useStorage("smoothingEnabled", false)

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

  return { scaling, smoothing }
})
