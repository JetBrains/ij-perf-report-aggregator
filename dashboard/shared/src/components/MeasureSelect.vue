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

        <span
          v-if="!slotProps.value || slotProps.value.length === 0"
        >
          {{ title }}
        </span>

        <span v-if="slotProps.value && slotProps.value.length === 1">
          {{ slotProps.value[0] }}
        </span>

        <span
          v-if="slotProps.value && slotProps.value.length > 1"
        >
          {{ props.selectedLabel(slotProps.value) }}
        </span>

        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #dropdownicon>
      <span class="hidden" />
    </template>
  </MultiSelect>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { computed } from "vue"
import { MeasureConfigurator } from "../configurators/MeasureConfigurator"

interface Props {
  configurator: MeasureConfigurator
  selectedLabel?: (items: string[]) => string
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: "Metrics",
  selectedLabel: (items: string[]) => `${items.length} items selected`,
})

// const items = props.configurator.data
// put selected values on top
const items = computed(() => {
  const configurator = props.configurator
  let result = configurator.data.value

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
    }
    else {
      return selectedValue.includes(b) ? 1 : 0
    }
  })
  return result
})

const value = computed({
  get(): string[] | null {
    const configurator = props.configurator
    const selectedValue = configurator.selected.value
    const allValues = configurator.data.value
    if (selectedValue != null && selectedValue.length > 0 && !allValues.some(it => selectedValue != null && selectedValue.includes(it))) {
      return null
    }
    return allValues.length === 0 ? null : selectedValue
  },
  set(value: string[] | null) {
    props.configurator.setSelected(value)
  },
})
</script>