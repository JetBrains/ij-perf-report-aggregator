<template>
  <el-form
    :inline="true"
    size="small"
  >
    <el-form-item label="Operator">
      <el-select
        v-model="value"
        filterable
      >
        <el-option
          v-for='name in ["median", "min", "max", "quantile"]'
          :key="name"
          :label="name"
          :value="name"
        />
      </el-select>
    </el-form-item>
    <el-form-item v-if="value === 'quantile'">
      <el-input-number
        v-model="quantile"
        :min="0"
        :max="100"
        :step="10"
      />
    </el-form-item>
  </el-form>
</template>
<script lang="ts">
import { defineComponent, computed } from "vue"
import { AggregationOperatorConfigurator } from "../configurators/AggregationOperatorConfigurator"

export default defineComponent({
  name: "AggregationOperatorSelect",
  props: {
    configurator: {
      type: AggregationOperatorConfigurator,
      required: true,
    },
  },
  setup(props) {
    return {
      value: computed({
        get() {
          return props.configurator.value.value.operator
        },
        set(value: string) {
          // eslint-disable-next-line
          props.configurator.value.value.operator = value
        },
      }),
      quantile: computed({
        get() {
          return props.configurator.value.value.quantile
        },
        set(value: number) {
          // eslint-disable-next-line
          props.configurator.value.value.quantile = value
        },
      }),
    }
  },
})
</script>