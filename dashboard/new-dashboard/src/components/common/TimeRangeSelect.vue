<template>
  <Select
    v-model="model"
    title="Time Range"
    :options="timerangeConfigurator.timeRanges.value.filter((element) => element.label != '')"
    option-label="label"
    option-value="value"
  >
    <template #value="slotProps">
      <div class="group flex items-center gap-1">
        <slot name="icon">
          <CalendarIcon class="w-4 h-4 text-gray-500" />
        </slot>

        <span v-if="!slotProps.value">Select range</span>
        <span>{{ currentValue?.label }}</span>

        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #dropdownicon>
      <span class="hidden" />
    </template>
    <template #footer>
      <ul class="p-select-list">
        <div
          class="p-select-option"
          @click="showCalendar"
        >
          Range
        </div>
      </ul>
      <DatePicker
        v-if="isShowCalendar"
        v-model="date"
        class="text-sm"
        date-format="dd/mm/yy"
        selection-mode="range"
        :inline="true"
        @click.stop
      />
    </template>
  </Select>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { computed, ref, watch } from "vue"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"

const { timerangeConfigurator } = defineProps<{
  timerangeConfigurator: TimeRangeConfigurator
}>()

const model = timerangeConfigurator.value
const currentValue = computed(() => timerangeConfigurator.timeRanges.value.find((item) => item.value === model.value))

const date = ref()
watch(date, (value: (Date | null)[] | null) => {
  if (value != null) {
    const start = value[0]
    const end = value[1]
    if (start != null && end != null) {
      timerangeConfigurator.setCustomRange(start, end)
    }
  }
})

const isShowCalendar = ref(false)

function showCalendar() {
  isShowCalendar.value = !isShowCalendar.value
}
</script>
<style #scoped>
.p-datepicker table {
  @apply text-sm;
}
.p-link {
  @apply text-sm;
}
</style>
