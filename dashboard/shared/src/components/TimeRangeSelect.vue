<template>
  <DropdownFilter
    v-model="value"
    label="Time Range"
    :options="TimeRangeConfigurator.timeRanges"
  />
</template>
<script setup lang="ts">
import DropdownFilter from "tailwind-ui/src/DropdownFilter.vue"
import { computed } from "vue"
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
</script>