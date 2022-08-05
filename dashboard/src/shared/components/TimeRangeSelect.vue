<template>
  <Dropdown
    v-model="value"
    title="Time Range"
    :options="timeRanges"
    option-label="label"
    option-value="value"
  />
</template>
<script setup lang="ts">
import { computed, shallowRef } from "vue"
import { TimeRangeConfigurator } from "../configurators/TimeRangeConfigurator"

const props =  defineProps<{
  configurator: TimeRangeConfigurator
}>()

const value = computed({
  get() {
    return props.configurator.value.value
  },
  set(value: string) {
    // eslint-disable-next-line vue/no-mutating-props
    props.configurator.value.value = value
  },
})
const timeRanges = shallowRef(TimeRangeConfigurator.timeRanges)
</script>