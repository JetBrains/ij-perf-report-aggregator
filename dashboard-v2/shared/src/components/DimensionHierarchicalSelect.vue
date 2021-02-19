<template>
  <el-form-item :label="label">
    <el-cascader
      v-model="value"
      :show-all-levels="false"
      :props="{label: 'value', checkStrictly: true, emitPath: false}"
      :options="values"
    />
  </el-form-item>
</template>
<script lang="ts">
import { defineComponent } from "vue"

import { MachineConfigurator } from "../configurators/MachineConfigurator"

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
    return {
      value: props.dimension.value,
      values: props.dimension.values,
    }
  },
})
</script>