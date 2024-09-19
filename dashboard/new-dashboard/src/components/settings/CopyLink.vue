<template>
  <div
    v-if="!timerangeConfigurator.customRange.value"
    class="card flex justify-content-center"
  >
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
    class="flex items-center p-2 bg-white rounded-lg shadow"
    role="alert"
  >
    <div class="font-normal">Copied</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getPersistentLink } from "./CopyLink"

const { timerangeConfigurator } = defineProps<{
  timerangeConfigurator: TimeRangeConfigurator
}>()

async function copyLink() {
  await navigator.clipboard.writeText(getPersistentLink(window.location.href, timerangeConfigurator))
  isToastVisible.value = true
  setTimeout(() => {
    isToastVisible.value = false
  }, 1500)
}

const isToastVisible = ref(false)
</script>
<style scoped></style>
