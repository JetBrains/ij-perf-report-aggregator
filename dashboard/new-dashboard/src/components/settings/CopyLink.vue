<template>
  <div class="card flex justify-content-center">
    <a
      title="Copy link to charts"
      class="info"
      @click="copyLink"
    >
      <span v-tooltip.bottom="'Copy link to dashboard with current date timeline.'">
        <ClipboardDocumentIcon class="w-6 h-6" />
      </span>
    </a>
  </div>
  <div
    v-show="isToastVisible"
    id="toast-default"
    class="flex items-center p-2 bg-white rounded-lg shadow-sm"
    role="alert"
  >
    <div class="font-normal">Copied</div>
  </div>
</template>

<script setup lang="ts">
import { inject, ref } from "vue"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getPersistentLink } from "./CopyLink"
import { sidebarVmKey } from "../../shared/keys"
import { pointParamName } from "../../shared/selectedPointStore"

const { timerangeConfigurator } = defineProps<{
  timerangeConfigurator: TimeRangeConfigurator
}>()

// The sidebar is optional: some hosts of this button don't provide one.
const sidebarVm = inject(sidebarVmKey, null)

async function copyLink() {
  await navigator.clipboard.writeText(getPersistentLink(linkWithSelectedPoint(), timerangeConfigurator))
  isToastVisible.value = true
  setTimeout(() => {
    isToastVisible.value = false
  }, 1500)
}

// Keep the `point` param in sync with the sidebar: include the opened build, and drop any
// stale point (e.g. inherited from a deep-link) when the sidebar isn't showing a point.
function linkWithSelectedPoint(): string {
  const url = new URL(window.location.href)
  const selectedBuildId = sidebarVm?.visible.value === true ? sidebarVm.data.value?.buildId : undefined
  if (selectedBuildId == undefined) {
    url.searchParams.delete(pointParamName)
  } else {
    url.searchParams.set(pointParamName, selectedBuildId.toString())
  }
  return url.toString()
}

const isToastVisible = ref(false)
</script>
<style scoped></style>
