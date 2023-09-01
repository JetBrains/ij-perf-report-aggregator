import { defineStore } from "pinia"
import { ref } from "vue"

export const useSmoothingStore = defineStore("smoothing", () => {
  const isSmoothingEnabled = ref(false)

  return { isSmoothingEnabled }
})
