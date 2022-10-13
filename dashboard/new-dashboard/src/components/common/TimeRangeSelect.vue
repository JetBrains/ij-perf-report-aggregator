<template>
  <Dropdown
    v-model="model"
    title="Time Range"
    :options="timeRanges"
    option-label="label"
    option-value="value"
  >
    <template #value="slotProps">
      <span class="flex items-center gap-1" v-if="slotProps.value">
        <slot name="icon"/>
        {{ currentValue?.label }}
      </span>

      <span class="flex items-center gap-1" v-if="!slotProps.value">
        <slot name="icon"/>
        Select range
      </span>
    </template>
  </Dropdown>
</template>
<script setup lang="ts">
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