<template>
  <Dropdown
    v-model="model"
    title="Time Range"
    :options="timeRanges"
    option-label="label"
    option-value="value"
  >
    <template #value="slotProps">
      <div class="group flex items-center gap-1">
        <slot name="icon" />

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
  </Dropdown>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { TimeRangeItem } from "shared/src/configurators/TimeRangeConfigurator"
import { computed, shallowRef } from "vue"

const props = defineProps<{
  value: string
  onChange: (value: string) => void
  ranges: TimeRangeItem[]
}>()

const model = computed({
  get() {
    return props.value
  },
  set(value: string) {
    props.onChange(value)
  },
})

const timeRanges = shallowRef(props.ranges)
const currentValue = computed(() => props.ranges.find(item => item.value === model.value))
</script>