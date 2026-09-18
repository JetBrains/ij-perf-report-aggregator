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
    <SettingSwitch
      v-for="{ setting, label, tooltip } in settings"
      :key="setting"
      :setting="setting"
      :label="label"
      :tooltip="tooltip"
      class="mb-2 last:mb-0"
    />
  </Popover>
</template>

<script setup lang="ts">
import { onBeforeUnmount, useTemplateRef } from "vue"
import SettingSwitch from "./SettingSwitch.vue"
import { BooleanSetting, SettingConfigurator } from "./SettingConfigurator"
import { useSettingsStore } from "./settingsStore"
import { storeToRefs } from "pinia"
import { PopoverMethods } from "openvue/popover"
import { DataQueryConfigurator } from "../common/dataQuery"
import { FilterConfigurator } from "../../configurators/filter"

/**
 * `rerunsQuery` marks a setting that post-processes the loaded data, so the query has to be re-run when it is flipped.
 * The other settings are observed where they are applied - `stdDevInterval` by the measure configurator, `fadeOnHover` by the chart itself.
 */
const settings: { setting: BooleanSetting; label: string; tooltip: string; rerunsQuery?: boolean }[] = [
  {
    setting: "smoothing",
    label: "Smoothing",
    tooltip: "Applies exponential smoothing to the dataset.",
    rerunsQuery: true,
  },
  {
    setting: "scaling",
    label: "Scaling",
    tooltip: "Scales each value in the dataset on its median. Each element is divided by the median and multiplied by 50.",
    rerunsQuery: true,
  },
  {
    setting: "detectChanges",
    label: "Detect Changes",
    tooltip: "Apply change detector algorithm. ",
    rerunsQuery: true,
  },
  {
    setting: "flexibleYZero",
    label: "Flexible zero on Y axis",
    tooltip: "Make zero on Y axis flexible and not always equals to 0.",
    rerunsQuery: true,
  },
  {
    setting: "removeOutliers",
    label: "Remove outliers (Alpha)",
    tooltip: "Remove outliers based on rolling MAD (Median Value Deviation) score",
    rerunsQuery: true,
  },
  {
    setting: "stdDevInterval",
    label: "Std dev interval",
    tooltip: "Wrap metrics reported as a mean in a ±1 standard deviation band, so the spread between the runs behind the mean is visible",
  },
  {
    setting: "fadeOnHover",
    label: "Fade others on hover",
    tooltip: "Fade out other series when hovering over one",
  },
]

const settingsPanel = useTemplateRef<PopoverMethods>("settingsPanel")
const settingsIcon = useTemplateRef<HTMLElement>("settingsIcon")

const emit = defineEmits<{
  "update:configurators": [configurator: DataQueryConfigurator & FilterConfigurator]
}>()
for (const { setting, rerunsQuery } of settings) {
  if (rerunsQuery) {
    emit("update:configurators", new SettingConfigurator(setting))
  }
}

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
