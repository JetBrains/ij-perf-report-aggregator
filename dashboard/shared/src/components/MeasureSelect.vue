<template>
  <el-form-item label="Metrics">
    <el-select
      v-model="value"
      multiple
      filterable
    >
      <el-option
        v-for="value in data"
        :key="value"
        :label="value"
        :value="value"
      />
    </el-select>
  </el-form-item>
</template>
<script lang="ts">
import { defineComponent } from "vue"

import { MeasureConfigurator } from "../configurators/MeasureConfigurator"

export default defineComponent({
  name: "MeasureSelect",
  props: {
    configurator: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const configurator = props.configurator as MeasureConfigurator
    return {
      value: configurator.value,
      data: configurator.data,
    }
  },
})
</script>
<style>
.el-form-item__label {
  display: flex;
  align-items: center;
}
</style>