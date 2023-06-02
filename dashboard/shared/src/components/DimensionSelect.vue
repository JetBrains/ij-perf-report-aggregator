<template>
  <MultiSelect
    v-if="valueToGroup == null && multiple"
    v-model="value"
    :title="label"
    :loading="dimension.state.loading"
    :disabled="dimension.state.disabled"
    :options="items"
    :placeholder="placeholder"
    option-label="label"
    option-value="value"
    :max-selected-labels="hasManyElements ? 1 : 2"
    :filter="hasManyElements"
    :show-toggle-all="hasManyElements"
  >
    <template #value="slotProps">
      <div class="group inline-flex justify-center text-sm font-medium text-gray-700 hover:text-gray-900">
        <span
          v-if="!slotProps.value || slotProps.value.length === 0"
          class="flex items-center gap-1"
        >
          <slot name="icon" />
          {{ placeholder }}
        </span>
        <span
          v-if="slotProps.value && slotProps.value.length === 1"
          class="flex items-center gap-1"
        >
          <slot name="icon" />
          {{ slotProps.value[0] }}
        </span>
        <span
          v-if="slotProps.value && slotProps.value.length > 1"
          class="flex items-center gap-1"
        >
          <slot name="icon" />
          {{ props.selectedLabel(slotProps.value) }}
        </span>
        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #dropdownicon>
      <span />
    </template>
  </MultiSelect>
  <Dropdown
    v-else-if="valueToGroup == null && !multiple"
    v-model="value"
    :title="label"
    :loading="dimension.state.loading"
    :disabled="dimension.state.disabled"
    :options="dimension.values.value"
    :placeholder="placeholder"
    :option-label="optionToLabel"
    :filter="hasManyElements"
    :auto-filter-focus="hasManyElements"
  >
    <!-- eslint-disable vue/no-template-shadow -->
    <template #value="{value}">
      <div class="group inline-flex justify-center text-sm font-medium text-gray-700 hover:text-gray-900">
        {{ value ? valueToLabel(value) : value }}
        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #dropdownicon>
      <!-- empty element to avoid ignoring override of slot -->
      <span />
    </template>
  </Dropdown>
  <MultiSelect
    v-else
    v-model="value"
    :title="label"
    :loading="dimension.state.loading"
    :disabled="dimension.state.disabled"
    :options="items"
    :placeholder="placeholder"
    option-label="label"
    option-value="value"
    option-group-children="options"
    option-group-label="label"
    :selection-limit="multiple ? null : 1"
    :max-selected-labels="1"
    :filter="hasManyElements"
    :auto-filter-focus="true"
  >
    <template #value="slotProps">
      <span
        v-if="!slotProps.value || slotProps.value.length === 0"
        class="flex items-center gap-1"
      >
        <slot name="icon" />
        {{ placeholder }}
      </span>
      <span
        v-if="slotProps.value && slotProps.value.length === 1"
        class="flex items-center gap-1 max-w-[200px] truncate"
      >
        <slot name="icon" />
        {{ slotProps.value[0] }}
      </span>
      <span
        v-if="slotProps.value && slotProps.value.length > 1"
        class="flex items-center gap-1"
      >
        <slot name="icon" />
        {{ props.selectedLabel(slotProps.value) }}
      </span>
    </template>
    <template #dropdownicon>
      <span />
    </template>
  </MultiSelect>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { computed } from "vue"
import { DimensionConfigurator } from "../configurators/DimensionConfigurator"
import { usePlaceholder } from "./placeholder"

const props = withDefaults(defineProps<{
  label: string
  dimension: DimensionConfigurator
  valueToLabel?: (v: string) => string
  // todo not working correctly for now (if value is set to not existing value, runtime error on select)
  valueToGroup?: ((v: string) => string) | null
  selectedLabel?: (items: string[]) => string
}>(), {
  selectedLabel: (items: string[]) => `${items.length} items selected`,
  valueToLabel: (v: string) => v,
  valueToGroup: null,
})

const multiple = computed(() => props.dimension.multiple)

const placeholder = usePlaceholder(props, () => props.dimension.values.value, () => props.dimension.selected.value)

function optionToLabel(value: string): string {
  return props.valueToLabel(value)
}

const value = computed<string | string[] | null>({
  get() {
    const values = props.dimension.values.value
    if (values.length === 0) {
      return null
    }

    const value = props.dimension.selected.value
    if (props.dimension.multiple) {
      if (Array.isArray(value)) {
        return value
      }
      else {
        return value == null || value === "" ? [] : [value]
      }
    }
    else {
      if (Array.isArray(value)) {
        return value[0]
      }
      else {
        return value === "" ? null : value
      }
    }
  },
  set(value) {
    // eslint-disable-next-line vue/no-mutating-props
    props.dimension.selected.value = value == null || value.length === 0 ? null : value
  },
})

const hasManyElements = computed(() => {
  return items.value.length > 3
})

const items = computed(() => {
  const valueToLabel = props.valueToLabel
  const values = props.dimension.values.value
  // map Array<string> to Array<Item> to be able to customize how value is displayed in UI
  if (props.valueToGroup == null) {
    const result = values.map(it => {
      return {label: valueToLabel(it.toString()), value: it}
    })
    if (values.length > 20) {
      // put selected values on top
      result.sort((a, b) => {
        if (value.value == null || !Array.isArray(value.value)) {
          return 0
        }
        return value.value.includes(a.label) ? -1 : (value.value.includes(b.label) ? 1 : 0)
      })
    }
    return result
  }
  else {
    return group(values as string[], props.valueToGroup, valueToLabel)
  }
})

interface Item {
  label: string
  value: string
}

interface GroupItem {
  label: string
  options: Item[]
}

function group(values: string[], groupFunction: (v: string) => string, valueToLabel: (v: string) => string): GroupItem[] {
  const groupNameToGroup = new Map<string, GroupItem>()
  const groups: GroupItem[] = []
  for (const value of values) {
    const groupName = groupFunction(value)
    let group = groupNameToGroup.get(groupName)
    if (group === undefined) {
      group = {
        label: groupName,
        options: [],
      }
      groupNameToGroup.set(groupName, group)
      groups.push(group)
    }
    group.options.push({label: valueToLabel(value), value})
  }
  // console.log(JSON.stringify(groups, null, 2))
  return groups
}
</script>