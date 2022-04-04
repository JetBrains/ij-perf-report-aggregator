<template>
  <OverlayPanel
    v-if="tooltipData != null"
    ref="panel"
    :show-close-icon="true"
  >
    <div
      @mouseenter="cancelHide"
      @mouseleave="hide"
    >
      <a
        v-if="tooltipData.linkUrl != null"
        :href="tooltipData.linkUrl"
        target="_blank"
      >
        {{ tooltipData.linkText }}
      </a>
      <span v-else>{{ tooltipData.linkText }}</span>

      <a
        v-if="tooltipData.firstSeriesData.length >= 3"
        title="Changes"
        :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[2]}&tab=buildChangesDiv`"
        target="_blank"
        class="info"
      >
        changes
      </a>

      <a
        v-if="tooltipData.firstSeriesData.length >= 4"
        title="Test Artifacts"
        :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[3]}&tab=artifacts`"
        target="_blank"
        class="info"
      >
        artifacts
      </a>

      <div
        v-for="item in tooltipData.items"
        :key="item.name"
        style="margin: 10px 0 0;white-space: nowrap"
      >
        <span
          class="tooltipNameMarker"
          :style='{"background-color": item.color}'
        />
        <span style="margin-left:2px;">{{ item.name }}</span>
        <span class="tooltipValue">{{ item.value }}</span>
      </div>
      <div
        v-if="tooltipData.firstSeriesData.length >= 8"
        style="margin: 10px 0 0;white-space: nowrap"
      >
        <span style="margin-left:2px;">machine</span>
        <span class="tooltipValue">{{ tooltipData.firstSeriesData[7] }}</span>
      </div>
    </div>
  </OverlayPanel>
</template>
<script setup lang="ts">
import OverlayPanel from "primevue/overlaypanel"
import { onBeforeUnmount, onMounted, ref } from "vue"
import { debounceSync } from "../util/debounce"
import { ChartToolTipManager, TooltipData } from "./ChartToolTipManager"

const tooltipData = ref<TooltipData | null>(null)
const panel = ref<OverlayPanel | null>()

let currentTarget: EventTarget | null

const hide = debounceSync(() => {
  currentTarget = null
  debouncedShow.clear()
  panel.value?.hide()
}, 2_000)

const debouncedShow = debounceSync(() => {
  if (currentTarget != null) {
    panel.value?.show({currentTarget} as Event)
    currentTarget = null
  }
}, 1)

function cancelHide() {
  hide.clear()
}

let isMetaPressed = false
let metaKey = ""
let lastManager: ChartToolTipManager | null = null
const keyDown = (event: KeyboardEvent) => {
  if (event.metaKey) {
    isMetaPressed = true
    metaKey = event.code
    if (lastManager != null) {
      lastManager.paused = true
    }
  }
}
const keyUp = (event: KeyboardEvent) => {
  if (event.code == metaKey) {
    isMetaPressed = false
    if (lastManager != null) {
      lastManager.paused = false
    }
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
<style scoped>
.tooltipNameMarker {
  display: inline-block;
  margin-right: 4px;
  border-radius: 10px;
  width: 10px;
  height: 10px;
}

.tooltipValue {
  @apply font-mono;
  float: right;
  margin-left: 20px;
}

a {
  text-decoration: none;
}

a.info {
  @apply text-gray-600;
}
</style>