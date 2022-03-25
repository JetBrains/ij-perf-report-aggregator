<template>
  <CascadeSelect
    v-model="value"
    :options="values"
    option-value="value"
    :option-group-children="['children']"
    option-label="value"
    option-group-label="value"
    :placeholder="label"
    @group-change="groupSelected"
  />
</template>
<script lang="ts">
import { defineComponent } from "vue"

import { GroupedDimensionValue, MachineConfigurator } from "../configurators/MachineConfigurator"

export interface GroupEvent {
  value: GroupedDimensionValue
}

export default defineComponent({
  name: "DimensionHierarchicalSelect",
  props: {
    label: {
      type: String,
      required: true,
    },
    dimension: {
      type: MachineConfigurator,
      required: true,
    },
  },
  setup(props) {
    const groupSelected = (e: GroupEvent) => {
      // eslint-disable-next-line vue/no-mutating-props
      props.dimension.value.value = e.value.children!.map(it => it.value)
    }
    return {
      value: props.dimension.value,
      values: props.dimension.values,
      groupSelected,
    }
  },
})
</script>