<template>
  <Dropdown
    v-model="model"
    title="Time Range"
    :options="props.timerangeConfigurator.timeRanges.value.filter((element) => element.label != '')"
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
  </Dropdown>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { computed } from "vue"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"

const props = defineProps<{
  timerangeConfigurator: TimeRangeConfigurator
}>()

const model = props.timerangeConfigurator.value

const currentValue = computed(() => props.timerangeConfigurator.timeRanges.value.find((item) => item.value === model.value))
</script>
