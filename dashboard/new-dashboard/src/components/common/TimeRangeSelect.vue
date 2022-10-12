<template>
  <Dropdown
    v-model="model"
    title="Time Range"
    :options="timeRanges"
    option-label="label"
    option-value="value"
  />
</template>
<script setup lang="ts">
import { TimeRangeItem } from "shared/src/configurators/TimeRangeConfigurator"
import { computed, shallowRef } from "vue"

const props =  defineProps<{
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
</script>