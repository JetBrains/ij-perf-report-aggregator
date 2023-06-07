<template>
  <!-- https://github.com/primefaces/primevue/issues/1725 loading is not supported -->
  <TreeSelect
    v-model="value"
    :disabled="dimension.state.disabled"
    title="Machine"
    :selection-mode="dimension.multiple ? 'multiple' : 'single'"
    :options="values"
    :placeholder="placeholder"
    class="max-w-lg"
  >
    <!-- eslint-disable vue/no-template-shadow -->
    <template #value="{ value }">
      <div class="group inline-flex justify-center text-sm font-medium text-gray-700 hover:text-gray-900">
        <template v-if="value && value.length > 0">
          <span class="flex items-center gap-1">
            <slot name="icon" />
            {{ value[0].label }}
          </span>
        </template>
        <template v-if="!value || value.length === 0">
          <span class="flex items-center gap-1">
            <slot name="icon" />
            {{ placeholder }}
          </span>
        </template>
        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #triggericon>
      <span />
    </template>
  </TreeSelect>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid"
import { TreeNode } from "primevue/tree"
import { computed } from "vue"
import { GroupedDimensionValue, MachineConfigurator } from "../../configurators/MachineConfigurator"
import { usePlaceholder } from "./placeholder"

function convertItemToTreeSelectModel(item: GroupedDimensionValue): TreeNode {
  return {
    key: item.value,
    label: item.value,
    children: item.children?.map(convertItemToTreeSelectModel),
    leaf: item.children == null,
  }
}

const props = defineProps<{
  label: string
  dimension: MachineConfigurator
}>()

type SelectedValue = Record<string, boolean>

const placeholder = usePlaceholder(
  props,
  () => props.dimension.values.value,
  () => props.dimension.selected.value
)

const value = computed<SelectedValue>({
  get(): SelectedValue {
    const result: SelectedValue = {}
    for (const k of props.dimension.selected.value) {
      result[k] = true
    }
    return result
  },
  set(value: SelectedValue) {
    // eslint-disable-next-line vue/no-mutating-props
    props.dimension.selected.value = Object.entries(value)
      .filter((it) => it[1])
      .map((it) => it[0])
  },
})
const values = computed(() => {
  return props.dimension.values.value.map((element) => convertItemToTreeSelectModel(element))
})
</script>
