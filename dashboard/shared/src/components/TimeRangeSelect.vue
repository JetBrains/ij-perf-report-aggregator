<template>
  <Dropdown
    v-model="value"
    title="Time Range"
    :options="timeRanges"
    option-label="label"
  />
</template>
<script setup lang="ts">
import { computed, shallowRef } from "vue"
import { TimeRangeConfigurator, TimeRangeItem } from "../configurators/TimeRangeConfigurator"

const props = defineProps<{
  configurator: TimeRangeConfigurator
}>()

const value = computed<TimeRangeItem>({
  get(): TimeRangeItem {
    return TimeRangeConfigurator.timeRangeValueToItem.get(props.configurator.value.value) ?? TimeRangeConfigurator.timeRanges[0]
  },
  set(item: TimeRangeItem) {
    // eslint-disable-next-line vue/no-mutating-props
    props.configurator.value.value = item.value
  },
})
const timeRanges = shallowRef(TimeRangeConfigurator.timeRanges)
</script>