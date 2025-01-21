import { defineStore } from "pinia"
import { watch } from "vue"
import { usePreferredDark, useStorage } from "@vueuse/core"

export const useDarkModeStore = defineStore("darkModeStore", () => {
  // State
  const darkMode = useStorage("darkMode", usePreferredDark().value)

  watch(
    darkMode,
    () => {
      const element = document.querySelector("html") as HTMLElement
      element.classList.toggle("dark-mode", darkMode.value)
    },
    { immediate: true }
  )

  const toggle = function () {
    darkMode.value = !darkMode.value
  }

  return { darkMode, toggle }
})
