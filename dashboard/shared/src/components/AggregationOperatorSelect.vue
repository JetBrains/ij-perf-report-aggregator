<template>
  <div class="flex-initial space-x-2 py-4">
    <label>
      <Dropdown
        v-model="value"
        title="Aggregate by"
        :options="operators"
        class="p-inputtext-sm"
        placeholder="Operator"
      />
    </label>
    <InputNumber
      v-if="value === 'quantile'"
      v-model="quantile"
      class="p-inputtext-sm"
      size="3"
      :min="0"
      :max="100"
      :step="10"
      :show-buttons="true"
    />
  </div>
</template>
<script lang="ts">
import { computed, defineComponent, inject, ref } from "vue"
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
      return props.configurator ?? providedConfigurator
    }

    const operators = ref<Array<string>>(["median", "min", "max", "quantile"])
    return {
      operators,
      value: computed({
        get() {
          return getConfigurator().value.value.operator
        },
        set(value: string) {
          getConfigurator().value.value.operator = value
        },
      }),
      quantile: computed({
        get() {
          return getConfigurator().value.value.quantile
        },
        set(value: number) {
          getConfigurator().value.value.quantile = value
        },
      }),
    }
  },
})
</script>