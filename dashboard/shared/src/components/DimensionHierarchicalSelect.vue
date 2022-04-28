<template>
  <!-- https://github.com/primefaces/primevue/issues/1725 loading is not supported -->
  <TreeSelect
    v-model="value"
    title="Machine"
    :selection-mode="props.dimension.multiple ? 'multiple' : 'single'"
    :options="values"
    :placeholder="placeholder"
    class="max-w-lg"
  />
</template>
<script setup lang="ts">
import { computed } from "vue"
import { GroupedDimensionValue, MachineConfigurator } from "../configurators/MachineConfigurator"
import { usePlaceholder } from "./placeholder"

function convertItemToTreeSelectModel(item: GroupedDimensionValue): unknown {
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

interface SelectedValue {
  [key: string]: boolean
}

const placeholder = usePlaceholder(props, () => props.dimension.values.value, () => props.dimension.value.value)

const value = computed<SelectedValue>({
  get(): SelectedValue {
    let v = props.dimension.value.value
    if (!Array.isArray(v)) {
      v = [v]
    }

    const result: SelectedValue = {}
    for (const k of v) {
      result[k] = true
    }
    return result
  },
  set(value: SelectedValue) {
    // eslint-disable-next-line vue/no-mutating-props
    props.dimension.value.value = Object.entries(value).filter(it => it[1]).map(it => it[0])
  },
})
const values = computed(() => {
  return props.dimension.values.value.map(convertItemToTreeSelectModel)
})
</script>