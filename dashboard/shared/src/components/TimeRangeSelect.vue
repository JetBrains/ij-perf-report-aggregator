<template>
  <SelectMenu
    v-model="selected"
    title="Time Range"
    :items="items"
    class="w-30"
  />
</template>
<script setup lang="ts">
import SelectMenu from "tailwind-ui/src/SelectMenu.vue"
import { computed } from "vue"
import { TimeRangeConfigurator, TimeRangeItem } from "../configurators/TimeRangeConfigurator"

const props =  defineProps<{
  configurator: TimeRangeConfigurator
}>()

const items = TimeRangeConfigurator.timeRanges
const selected = computed<TimeRangeItem>({
  get() {
    return TimeRangeConfigurator.findItemByValue(props.configurator.value.value) ?? TimeRangeConfigurator.timeRanges[0]
  },
  set(value) {
    // eslint-disable-next-line vue/no-mutating-props
    props.configurator.value.value = value?.value
  },
})
</script>