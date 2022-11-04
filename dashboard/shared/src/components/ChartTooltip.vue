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
          <a
            v-if="linkUrl != null"
            :href="linkUrl"
            target="_blank"
          >
            {{ linkText }}
          </a>
          <span v-else>{{ linkText }}</span>
        </div>

        <a
          v-if="tooltipData.firstSeriesData.length >= 4"
          title="Changes"
          :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[3]}&tab=buildChangesDiv`"
          target="_blank"
          class="info"
        >
          <UsersIcon class="w-5 h-5" />
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
          :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[4]}&tab=artifacts`"
          target="_blank"
          class="info"
        >
          <ArrowDownTrayIcon class="w-5 h-5" />
        </a>
      </div>
      <div
        class="grid grid-cols-[repeat(3,_max-content)] whitespace-nowrap gap-x-2 items-baseline leading-loose"
      >
        <template
          v-for="item in tooltipData.items"
          :key="item.name"
        >
          <span
            class="rounded-lg w-2.5 h-2.5"
            :style='{"background-color": item.color}'
          />
          <span>{{ item.name }}</span>
          <span class="font-mono place-self-end">{{ getValueFormatterByMeasureName(item.name)(item.value) }}</span>
        </template>
      </div>
      <!-- in our case data is related to each other, reduce height of panel -->
      <Divider class="my-2" />
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

import { ArchiveBoxIcon, ArrowDownTrayIcon, ServerIcon, UsersIcon } from "@heroicons/vue/24/outline"
import OverlayPanel from "primevue/overlaypanel"
import Divider from "tailwind-ui/src/Divider.vue"
import { computed, nextTick, shallowRef, watch, WatchStopHandle } from "vue"
import { getValueFormatterByMeasureName, timeFormatWithoutSeconds } from "../formatter"
import { debounceSync } from "../util/debounce"
import { ChartToolTipManager, TooltipData } from "./ChartToolTipManager"

const tooltipData = shallowRef<TooltipData | null>(null)
const panel = shallowRef<OverlayPanel | null>()

const linkText = computed(() => {
  const data = tooltipData.value?.firstSeriesData
  if (data == null) {
    return ""
  }

  const generatedTime = data[0]
  let result = timeFormatWithoutSeconds.format(generatedTime)
  if (data[5] && data[6]) {
    result += ` (${data[5]}.${data[6]}${data[7] === 0 ? "" : `.${data[7]}`})`
  }
  return result
})

const linkUrl = computed(() => {
  const data = tooltipData.value
  const query = data?.query
  if (query == null) {
    return null
  }
  return data?.reportInfoProvider?.createReportUrl(data.firstSeriesData[0], query)
})

// do not show tooltip if cursor was not pointed again to some item
let lastManager: ChartToolTipManager | null = null

const hide = debounceSync(() => {
  console.log("hide")
  lastManager = null
  panel.value?.hide()
  tooltipData.value = null
}, 1_000)

let panelTargetIsChanged = false
const consumer: (data: TooltipData | null, event: Event | null) => void = (data, event) => {
  console.log("new data", data == null, event == null, panelTargetIsChanged)
  // console.log(event)
  hide.clear()

  tooltipData.value = data

  if (data == null || event == null) {
    return
  }

  let stopHandle: WatchStopHandle | null = null
  const panelElement = panel.value
  if (panelElement == null) {
    stopHandle = watch(panel, value => {
      value?.show(event)
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      stopHandle!()
    })
  }
  else {
    if (panelTargetIsChanged) {
      panelTargetIsChanged = false
      // make sure that we show panel near the target - PrimeVue doesn't handle that if panel is shown
      panelElement.hide()
      // show only after hide is performed
      void nextTick(() => {
        panelElement.show(event)
      })
    }
    else {
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

defineExpose({
  show(manager: ChartToolTipManager) {
    if (lastManager === manager) {
      hide.clear()
      panelTargetIsChanged = false
    }
    else if (lastManager != null) {
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