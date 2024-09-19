<template>
  <MultiSelect
    v-model="value"
    :options="items"
    :loading="configurator.state.loading"
    :disabled="configurator.state.disabled"
    :title="title"
    :placeholder="title"
    :filter="true"
    :auto-filter-focus="true"
    :option-label="(it: string) => it"
    :max-selected-labels="1"
  >
    <template #value="slotProps">
      <div class="group flex items-center gap-1">
        <slot name="icon" />

        <span v-if="!slotProps.value || slotProps.value.length === 0">
          {{ title }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length === 1">
          {{ slotProps.value[0] }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length > 1">
          {{ selectedLabel(slotProps.value) }}
        </span>

        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #dropdownicon>
      <span class="hidden" />
    </template>
    <template #header>
      <div
        v-if="configurator instanceof MeasureConfigurator"
        class="bg-gray-100 rounded-md border"
      >
        <div class="flex items-center w-full mt-2 ml-2 mb-2">
          <ToggleSwitch v-model="showAllMetrics" />
          <span
            v-tooltip.left="'Disable filtering of detailed metrics like completion_32.'"
            class="ml-2"
            >Show all metrics</span
          >
        </div>
      </div>
    </template>
  </MultiSelect>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { computed, ref, watch } from "vue"
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"

interface Props {
  configurator: MeasureConfigurator | SimpleMeasureConfigurator
  selectedLabel?: (items: string[]) => string
  title?: string
}

const { configurator, selectedLabel = (items: string[]) => `${items.length} items selected`, title = "Metrics" } = defineProps<Props>()
const showAllMetrics = ref(false)
watch(showAllMetrics, (value) => {
  if (configurator instanceof MeasureConfigurator) {
    configurator.setShowAllMetrics(value)
  }
})

// const items = props.configurator.data
// put selected values on top
const items = computed(() => {
  let result = configurator.data.value

  if (result == null) {
    return []
  }
  if (result.length < 20) {
    return result
  }

  const selectedValue = configurator.selected.value ?? []
  if (selectedValue.length === 0) {
    return result
  }

  // do not modify original array
  result = [...result]
  result.sort((a, b) => {
    if (selectedValue.includes(a)) {
      return selectedValue.includes(b) ? 0 : -1
    } else {
      return selectedValue.includes(b) ? 1 : 0
    }
  })
  return result
})

const value = computed({
  get(): string[] | null {
    const selectedValue = configurator.selected.value
    const allValues = configurator.data.value
    if (allValues == null) return null
    if (selectedValue != null && selectedValue.length > 0 && !allValues.some((it) => selectedValue.includes(it))) {
      return null
    }
    return allValues.length === 0 ? null : selectedValue
  },
  set(value: string[] | null) {
    configurator.setSelected(value)
  },
})
</script>
