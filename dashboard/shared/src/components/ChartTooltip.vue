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
      <div class="text-xs pt-3">
        Hold {{ metaKeySymbol }} to prevent popup refresh
      </div>
    </div>
  </OverlayPanel>
</template>
<script setup lang="ts">

import OverlayPanel from "primevue/overlaypanel"
import Divider from "tailwind-ui/src/Divider.vue"
import { computed, onBeforeUnmount, onMounted, shallowRef } from "vue"
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
  const data = tooltipData?.value
  const query = data?.query
  if (query == null) {
    return null
  }
  return data?.reportInfoProvider?.createReportUrl(data.firstSeriesData[0], query)
})

let currentTarget: EventTarget | null

// noinspection JSDeprecatedSymbols
const metaKeySymbol = window.navigator.platform.toUpperCase().includes("MAC") ? "⌘" : "⊞"

// do not show tooltip if cursor was not pointed again to some item
function resetState() {
  const value = tooltipData.value
  if (value != null) {
    value.items.length = 0
  }

  if (lastManager != null) {
    lastManager.paused = false
  }
}

const hide = debounceSync(() => {
  if (lastManager?.paused) {
    return
  }
  currentTarget = null
  debouncedShow.clear()
  panel.value?.hide()
  resetState()
}, 2_000)

function panelClosedExplicitly() {
  hide.clear()
  currentTarget = null
  resetState()
}

const debouncedShow = debounceSync(() => {
  if (currentTarget != null) {
    panel.value?.show({currentTarget} as Event)
    currentTarget = null
  }
}, 300)

function cancelHide() {
  hide.clear()
}

let metaKey = ""
let lastManager: ChartToolTipManager | null = null
const keyDown = (event: KeyboardEvent) => {
  if (event.metaKey) {
    metaKey = event.code
    if (lastManager != null) {
      lastManager.paused = true
    }
  }
}
const keyUp = (event: KeyboardEvent) => {
  if (event.code === metaKey && lastManager != null) {
      lastManager.paused = false
    }
}

onMounted(() => {
  window.addEventListener("keyup", keyUp)
  window.addEventListener("keydown", keyDown)
})

onBeforeUnmount(() => {
  window.removeEventListener("keydown", keyDown)
  window.removeEventListener("keyup", keyUp)
})

defineExpose({
  show(event: MouseEvent, manager: ChartToolTipManager) {
    if(lastManager?.paused) return
    if (lastManager !== manager && lastManager !== null) {
      currentTarget = null
      debouncedShow.clear()
      panel.value?.hide()
    }

    lastManager = manager
    hide.clear()
    if (!lastManager.paused) {
      tooltipData.value = manager.reportTooltipData
    }
    currentTarget = event.currentTarget
    debouncedShow()
  },
  hide,
})

</script>