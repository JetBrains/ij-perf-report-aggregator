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
import { AggregationOperatorConfigurator } from "../configurators/AggregationOperatorConfigurator"
import { aggregationOperatorConfiguratorKey } from "../injectionKeys"

export default defineComponent({
  name: "AggregationOperatorSelect",
  props: {
    configurator: {
      type: AggregationOperatorConfigurator,
      default: null,
    },
  },
  setup(props) {
    const providedConfigurator = inject(aggregationOperatorConfiguratorKey, null)

    function getConfigurator(): AggregationOperatorConfigurator {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      return props.configurator ?? providedConfigurator
    }

    return {
      value: computed({
        get() {
          return getConfigurator().value.value.operator
        },
        set(value: string) {
          // eslint-disable-next-line
          getConfigurator().value.value.operator = value
        },
      }),
      quantile: computed({
        get() {
          return getConfigurator().value.value.quantile
        },
        set(value: number) {
          // eslint-disable-next-line
          getConfigurator().value.value.quantile = value
        },
      }),
    }
  },
})
</script>