<template>
  <OverlayPanel
    v-if="tooltipData != null"
    ref="panel"
    :show-close-icon="true"
    @onclose="panelClosedExplicitly"
  >
    <div
      class="text-sm"
      @mouseenter="cancelHide"
      @mouseleave="hide"
    >
      <div class="flex gap-x-2 justify-end">
        <div class="w-full">
          <span>{{ linkText }}</span>
        </div>

        <a
          v-if="tooltipData.firstSeriesData[5]"
          title="Changes"
          target="_blank"
          class="info cursor-pointer"
          @click="navigateToSpace"
        >
          <SpaceIcon class="w-5 h-5" />
        </a>

        <a
          v-if="tooltipData.firstSeriesData.length >= 4"
          title="Changes"
          :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[3]}&tab=buildChangesDiv`"
          target="_blank"
          class="info"
        >
          <UserIcon class="w-5 h-5" />
        </a>

        <a
          v-if="tooltipData.firstSeriesData.length >= 4"
          title="Test Artifacts"
          :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[3]}&tab=artifacts`"
          target="_blank"
          class="info"
        >
          <ArchiveBoxIcon class="w-5 h-5" />
        </a>

        <a
          v-if="tooltipData.firstSeriesData.length >= 5"
          title="Installer Artifacts"
          :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[5]}&tab=artifacts`"
          target="_blank"
          class="info"
        >
          <ArrowDownTrayIcon class="w-5 h-5" />
        </a>
      </div>
      <div class="grid grid-cols-[repeat(3,_max-content)] whitespace-nowrap gap-x-2 items-baseline leading-loose">
        <template
          v-for="item in tooltipData.items"
          :key="item.name"
        >
          <span
            class="rounded-lg w-2.5 h-2.5"
            :style="{ 'background-color': item.color }"
          />
          <span>{{ item.name }}</span>
          <span class="font-mono place-self-end">{{ getValueFormatterByMeasureName(item.name)(item.value) }}</span>
        </template>
      </div>
      <!-- in our case data is related to each other, reduce height of panel -->
      <div class="relative mt-3 mb-2">
        <div
          class="absolute inset-0 flex items-center"
          aria-hidden="true"
        >
          <div class="w-full border-t border-gray-300" />
        </div>
      </div>
      <div
        v-if="tooltipData.firstSeriesData.length >= 3"
        class="grid grid-cols-[repeat(2,max-content)] items-center gap-x-1"
      >
        <ServerIcon class="w-5 h-5" />
        <span>{{ tooltipData.firstSeriesData[2] }}</span>
      </div>
    </div>
  </OverlayPanel>
</template>
<script setup lang="ts">
import OverlayPanel from "primevue/overlaypanel"
import { computed, nextTick, shallowRef, watch, WatchStopHandle } from "vue"
import { calculateChanges } from "../../util/changes"
import { debounceSync } from "../../util/debounce"
import SpaceIcon from "../common/SpaceIcon.vue"
import { getValueFormatterByMeasureName, timeFormatWithoutSeconds } from "../common/formatter"
import { StartupTooltipManager, TooltipData } from "./StartupTooltipManager"

const tooltipData = shallowRef<TooltipData | null>(null)
const panel = shallowRef<OverlayPanel | null>()

const linkText = computed(() => {
  const data = tooltipData.value?.firstSeriesData
  if (data == null) {
    return ""
  }

  const generatedTime = data[0]
  let result = timeFormatWithoutSeconds.format(generatedTime)
  if (data[6] && data[7]) {
    result += ` (${data[6]}.${data[7]}${data[8] === 0 ? "" : `.${data[8]}`})`
  }
  return result
})

// do not show tooltip if cursor was not pointed again to some item
let lastManager: StartupTooltipManager | null = null

const hide = debounceSync(() => {
  lastManager = null
  panel.value?.hide()
  tooltipData.value = null
}, 1_000)

let panelTargetIsChanged = false
const consumer: (data: TooltipData | null, event: Event | null) => void = (data, event) => {
  hide.clear()

  tooltipData.value = data

  if (data == null || event == null) {
    return
  }

  let stopHandle: WatchStopHandle | null = null
  const panelElement = panel.value
  if (panelElement == null) {
    stopHandle = watch(panel, (value) => {
      value?.show(event)
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      stopHandle!()
    })
  } else {
    if (panelTargetIsChanged) {
      panelTargetIsChanged = false
      // make sure that we show panel near the target - PrimeVue doesn't handle that if panel is shown
      panelElement.hide()
      // show only after hide is performed
      void nextTick(() => {
        panelElement.show(event)
      })
    } else {
      panelElement.show(event)
    }
  }
}

function panelClosedExplicitly() {
  hide.clear()
  tooltipData.value = null
}

function cancelHide() {
  hide.clear()
}

function navigateToSpace() {
  if (tooltipData.value != null) {
    calculateChanges("ij", tooltipData.value.firstSeriesData[5], (decodedChanges: string | null) => {
      if (decodedChanges == null) {
        window.open(`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.value?.firstSeriesData[3]}&tab=buildChangesDiv`)
      } else {
        window.open("https://jetbrains.team/p/ij/repositories/ultimate/commits?query=%22" + decodedChanges + "%22&tab=changes")
      }
    })
  }
}

defineExpose({
  show(manager: StartupTooltipManager) {
    if (lastManager === manager) {
      hide.clear()
      panelTargetIsChanged = false
    } else if (lastManager != null) {
      panelTargetIsChanged = true
    }

    lastManager = manager
    manager.setConsumer(consumer)
  },
  hide() {
    lastManager?.setConsumer(null)
    hide()
  },
})
</script>
