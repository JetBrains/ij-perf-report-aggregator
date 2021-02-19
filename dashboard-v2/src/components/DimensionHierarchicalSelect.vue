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
import ElCascader from "element-plus/es/el-cascader"

export default defineComponent({
  name: "DimensionHierarchicalSelect",
  components: {
    ElCascader,
  },
  props: {
    label: {
      type: String,
      required: true,
    },
    dimension: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const configurator = props.dimension as MachineConfigurator
    return {
      value: configurator.value,
      values: configurator.values,
    }
  },
})
</script>