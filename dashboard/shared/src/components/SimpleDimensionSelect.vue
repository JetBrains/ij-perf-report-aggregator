<template>
  <DropdownFilter
    v-model="value"
    :label="label"
    :options="items"
  />
</template>
<script setup lang="ts">
import DropdownFilter from "tailwind-ui/src/DropdownFilter.vue"
import { Option } from "tailwind-ui/src/tabModel"
import { computed } from "vue"
import { DimensionConfigurator } from "../configurators/DimensionConfigurator"

const props = withDefaults(defineProps<{
  label: string
  dimension: DimensionConfigurator
  valueToLabel?: (v: string) => string
}>(), {
  valueToLabel: (v: string) => v,
  valueToGroup: null,
})

const items = computed<Array<Option>>(() => {
  const valueToLabel = props.valueToLabel
  return props.dimension.values.value.map(it => {
    return {label: valueToLabel(it.toString()), value: it as string}
  })
})

const value = computed<Option | null>({
  get(): Option | null {
    const values = props.dimension.values.value
    if (values == null || values.length === 0) {
      return null
    }

    const value = props.dimension.selected.value
    const normalizedValue = Array.isArray(value) ? value[0] : (value === "" ? null : value)
    return normalizedValue == null ? null : {label: props.valueToLabel(normalizedValue), value: normalizedValue}
  },
  set(value) {
    // eslint-disable-next-line vue/no-mutating-props
    props.dimension.selected.value = value?.value ?? null
  },
})

</script>