<template>
  <CascadeSelect
    v-model="value"
    :options="values"
    option-value="value"
    :option-group-children="['children']"
    option-label="value"
    option-group-label="value"
    :placeholder="cLabel"
    @group-change="groupSelect"
  />
</template>
<script lang="ts">
import { computed, defineComponent, ref } from "vue"

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
    const selectedGroup = ref<string>()
    const groupSelect = (e: GroupEvent) => {
      selectedGroup.value = e.value.value
      if(e.value.children != null) {
        // eslint-disable-next-line vue/no-mutating-props
        props.dimension.value.value = e.value.children.map(it => it.value)
      }
    }
    return {
      cLabel: computed({
        get(): string {
          if (typeof props.dimension.value.value == "string") {
            return props.dimension.value.value
          }
          else {
            return selectedGroup.value ?? ""
          }
        },
        set(value: string) {
          console.log(value)
        },
      }),
      value: props.dimension.value,
      values: props.dimension.values,
      groupSelect,
    }
  },
})
</script>