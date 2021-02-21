<template>
  <el-form-item label="Metrics">
    <el-select
      v-model="selectedValues"
      multiple
      filterable
    >
      <el-option
        v-for="value in values"
        :key="value"
        :label="value"
        :value="value"
      />
    </el-select>
  </el-form-item>
</template>
<script lang="ts">
import { defineComponent } from "vue"
import { MeasureConfigurator } from "@/configurators/MeasureConfigurator"

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
      selectedValues: configurator.selectedValues,
      values: configurator.values,
    }
  },
})
</script>