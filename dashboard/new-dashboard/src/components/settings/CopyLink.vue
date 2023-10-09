<template>
  <div
    v-if="!props.timerangeConfigurator.customRange.value"
    class="card flex justify-content-center"
  >
    <a
      title="Copy link to charts"
      class="info"
      @click="copyLink"
    >
      <VTooltip theme="info">
        <span>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-6 h-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M8.25 7.5V6.108c0-1.135.845-2.098 1.976-2.192.373-.03.748-.057 1.123-.08M15.75 18H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08M15.75 18.75v-1.875a3.375 3.375 0 00-3.375-3.375h-1.5a1.125 1.125 0 01-1.125-1.125v-1.5A3.375 3.375 0 006.375 7.5H5.25m11.9-3.664A2.251 2.251 0 0015 2.25h-1.5a2.251 2.251 0 00-2.15 1.586m5.8 0c.065.21.1.433.1.664v.75h-6V4.5c0-.231.035-.454.1-.664M6.75 7.5H4.875c-.621 0-1.125.504-1.125 1.125v12c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V16.5a9 9 0 00-9-9z"
            />
          </svg>
        </span>
        <template #popper><span class="text-sm">Copy link to dashboard with current date timeline.</span></template>
      </VTooltip>
    </a>
  </div>
  <div
    v-show="isToastVisible"
    id="toast-default"
    class="flex items-center w-full max-w-xs p-4 text-gray-500 bg-white rounded-lg shadow dark:text-gray-400 dark:bg-gray-800"
    role="alert"
  >
    <div class="ml-3 text-sm font-normal">Copied</div>
    <button
      type="button"
      class="ml-auto -mx-1.5 -my-1.5 bg-white text-gray-400 hover:text-gray-900 rounded-lg focus:ring-2 focus:ring-gray-300 p-1.5 hover:bg-gray-100 inline-flex items-center justify-center h-8 w-8 dark:text-gray-500 dark:hover:text-white dark:bg-gray-800 dark:hover:bg-gray-700"
      data-dismiss-target="#toast-default"
    ></button>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { parseDuration, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"

const props = defineProps<{
  timerangeConfigurator: TimeRangeConfigurator
}>()

async function copyLink() {
  let url = window.location.href

  const now = new Date()
  const ago = getDateAgoByDuration(props.timerangeConfigurator.value.value)
  const filter = `${ago.getFullYear()}-${ago.getMonth() + 1}-${ago.getDate()}:${now.getFullYear()}-${now.getMonth() + 1}-${now.getDate()}`
  url = url.replace(new RegExp("&?customRange=.+&?"), "")
  url = url.replace(new RegExp("&?timeRange=.+&?"), "")
  await navigator.clipboard.writeText(url + "&timeRange=custom&customRange=" + filter)
  isToastVisible.value = true
  setTimeout(() => {
    isToastVisible.value = false
  }, 1500)
}
const isToastVisible = ref(false)

function getDateAgoByDuration(s: string): Date {
  const result = parseDuration(s)
  let days = 0

  if (result.days != null) {
    days += result.days
  }
  if (result.months != null) {
    days += result.months * 31
  }

  if (result.weeks != null) {
    days += result.weeks * 7
  }
  if (result.years != null) {
    days += result.years * 365
  }
  const date = new Date()
  date.setDate(date.getDate() - days)
  return date
}
</script>
<style scoped></style>
