import { useStorage } from "@vueuse/core"
import { defineStore } from "pinia"
import { onMounted, onUnmounted } from "vue"

export const useScrollStore = defineStore("scrollStore", () => {
  const isScrolled = useStorage("scrolled", false)

  const updateIsScrolled = () => {
    isScrolled.value = window.scrollY > 100
  }

  const setupListener = () => {
    window.addEventListener("scroll", updateIsScrolled)
    updateIsScrolled()
  }

  const cleanupListener = () => {
    window.removeEventListener("scroll", updateIsScrolled)
  }

  return { isScrolled, setupListener, cleanupListener }
})

export function useScrollListeners() {
  const scrollStore = useScrollStore()

  onMounted(scrollStore.setupListener)
  onUnmounted(scrollStore.cleanupListener)
}
