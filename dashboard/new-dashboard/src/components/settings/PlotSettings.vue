<template>
  <Cog8ToothIcon
    ref="settingsIcon"
    :class="'w-6 h-6 ' + (removeOutliers ? 'text-red-500' : 'text-primary dark:text-primary-dark')"
    @click="showSettings"
  />
  <Popover
    ref="settingsPanel"
    class="flex flex-col"
    append-to="body"
  >
    <SmoothingSwitch class="mb-2" />
    <ScalingSwitch class="mb-2" />
    <DetectChangesSwitch class="mb-2" />
    <FlexibleZeroOnYAxis class="mb-2" />
    <RemoveOutliersSwitch />
  </Popover>
</template>

<script setup lang="ts">
import { onBeforeUnmount, useTemplateRef } from "vue"
import { DetectChangesConfigurator } from "./configurators/DetectChangesConfigurator"
import { ScalingConfigurator } from "./configurators/ScalingConfigurator"
import { SmoothingConfigurator } from "./configurators/SmoothingConfigurator"
import DetectChangesSwitch from "./DetectChangesSwitch.vue"
import ScalingSwitch from "./ScalingSwitch.vue"
import SmoothingSwitch from "./SmoothingSwitch.vue"
import FlexibleZeroOnYAxis from "./FlexibleYZeroSwitch.vue"
import { FlexibleZeroOnYAxisConfigurator } from "./configurators/FlexibleZeroOnYAxisConfigurator"
import RemoveOutliersSwitch from "./RemoveOutliersSwitch.vue"
import { RemoveOutliersConfigurator } from "./configurators/RemoveOutliersConfigurator"
import { useSettingsStore } from "./settingsStore"
import { storeToRefs } from "pinia"
import { PopoverMethods } from "primevue/popover"
import { DataQueryConfigurator } from "../common/dataQuery"
import { FilterConfigurator } from "../../configurators/filter"

const settingsPanel = useTemplateRef<PopoverMethods>("settingsPanel")
const settingsIcon = useTemplateRef<HTMLElement>("settingsIcon")

const emit = defineEmits<{
  "update:configurators": [configurator: DataQueryConfigurator & FilterConfigurator]
}>()
emit("update:configurators", new ScalingConfigurator())
emit("update:configurators", new SmoothingConfigurator())
emit("update:configurators", new DetectChangesConfigurator())
emit("update:configurators", new FlexibleZeroOnYAxisConfigurator())
emit("update:configurators", new RemoveOutliersConfigurator())

const settingsStore = useSettingsStore()
const { removeOutliers } = storeToRefs(settingsStore)

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
  const overlayElement = document.querySelector(".p-popover")

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
