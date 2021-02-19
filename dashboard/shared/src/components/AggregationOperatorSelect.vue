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
import { defineComponent, computed, inject } from "vue"
import { aggregationOperatorConfiguratorKey } from "../componentKeys"
import { AggregationOperatorConfigurator } from "../configurators/AggregationOperatorConfigurator"

function getConfigurator(props: {configurator?: AggregationOperatorConfigurator}): AggregationOperatorConfigurator {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  return props.configurator ?? inject(aggregationOperatorConfiguratorKey)!
}

export default defineComponent({
  name: "AggregationOperatorSelect",
  props: {
    configurator: {
      type: AggregationOperatorConfigurator,
      default: null,
    },
  },
  setup(props) {
    return {
      value: computed({
        get() {
          return getConfigurator(props).value.value.operator
        },
        set(value: string) {
          // eslint-disable-next-line
          getConfigurator(props).value.value.operator = value
        },
      }),
      quantile: computed({
        get() {
          return getConfigurator(props).value.value.quantile
        },
        set(value: number) {
          // eslint-disable-next-line
          getConfigurator(props).value.value.quantile = value
        },
      }),
    }
  },
})
</script>