<template>
  <Cog8ToothIcon
    ref="settingsIcon"
    class="w-6 h-6 text-primary"
    @click="showSettings"
  />
  <OverlayPanel
    ref="settingsPanel"
    class="flex flex-col"
    append-to="body"
  >
    <SmoothingSwitch class="mb-2" />
    <ScalingSwitch class="mb-2" />
    <DetectChangesSwitch />
  </OverlayPanel>
</template>

<script setup lang="ts">
import OverlayPanel from "primevue/overlaypanel"
import { onBeforeUnmount, shallowRef } from "vue"
import { DetectChangesConfigurator } from "../../configurators/DetectChangesConfigurator"
import { ScalingConfigurator } from "../../configurators/ScalingConfigurator"
import { SmoothingConfigurator } from "../../configurators/SmoothingConfigurator"
import DetectChangesSwitch from "./DetectChangesSwitch.vue"
import ScalingSwitch from "./ScalingSwitch.vue"
import SmoothingSwitch from "./SmoothingSwitch.vue"

const settingsPanel = shallowRef<OverlayPanel | null>(null)
const settingsIcon = shallowRef<HTMLElement>()

const emit = defineEmits(["update:configurators"])
emit("update:configurators", new ScalingConfigurator())
emit("update:configurators", new SmoothingConfigurator())
emit("update:configurators", new DetectChangesConfigurator())

const showSettings = function (event: Event) {
  settingsPanel.value?.toggle(event, settingsIcon.value) // Toggle the panel first
  setTimeout(() => {
    adjustPosition()
    window.addEventListener("scroll", adjustPosition)
  }, 0)
}

// this is a hack since appendTo doesn't work with Icon for some reason
function adjustPosition() {
  const iconRect = settingsIcon.value?.getBoundingClientRect()

  // Query for the OverlayPanel's DOM element.
  const overlayElement = document.querySelector(".p-overlaypanel")

  if (iconRect && overlayElement != null) {
    let leftPosition = iconRect.left
    const overlayHTMLElement = overlayElement as HTMLElement
    const overlayWidth = overlayHTMLElement.offsetWidth

    // Screen margin to prevent the overlay from sticking to the edge.
    const screenMargin = 20

    const topPosition = iconRect.bottom + window.scrollY

    // If the OverlayPanel would overflow the right edge of the screen
    if (leftPosition + overlayWidth + screenMargin > window.innerWidth) {
      leftPosition = window.innerWidth - overlayWidth - screenMargin
    }

    // If the adjusted position would still overflow the left edge (i.e., it's wider than the screen), just set it to the margin value.
    if (leftPosition < screenMargin) {
      leftPosition = screenMargin
    }

    const verticalMargin = 10 // Margin between the icon and the OverlayPanel
    overlayHTMLElement.style.top = `${topPosition + verticalMargin}px`
    overlayHTMLElement.style.left = `${leftPosition}px`
  }
}

onBeforeUnmount(() => {
  window.removeEventListener("scroll", adjustPosition)
})
</script>

<style scoped></style>
