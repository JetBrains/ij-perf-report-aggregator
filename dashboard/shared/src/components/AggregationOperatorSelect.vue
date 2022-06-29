<template>
  <div class="flex-initial space-x-2">
    <SelectMenu
      v-model="value"
      title="Aggregate by"
      :options="operators"
      option-label="label"
      option-value="value"
      placeholder="Operator"
    />
    <InputNumber
      v-if="value === 'quantile'"
      v-model="quantile"
      size="3"
      :min="0"
      :max="100"
      :step="10"
      :show-buttons="true"
    />
  </div>
</template>
<script setup lang="ts">
import { computed, inject } from "vue"
import { AggregationOperatorConfigurator } from "../configurators/AggregationOperatorConfigurator"
import { aggregationOperatorConfiguratorKey } from "../injectionKeys"

const props =  defineProps<{
  configurator?: AggregationOperatorConfigurator
}>()

const providedConfigurator = inject(aggregationOperatorConfiguratorKey, null)

function getConfigurator(): AggregationOperatorConfigurator {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  return props.configurator ?? providedConfigurator!
}

const operators = ["median", "min", "max", "quantile"].map(it => ({label: it, value: it}))
const value = computed({
  get() {
    return getConfigurator().operator.value
  },
  set(value: string) {
    getConfigurator().operator.value = value
  },
})
const quantile = computed({
  get() {
    return getConfigurator().quantile.value
  },
  set(value: number) {
    getConfigurator().quantile.value = value
  },
})
</script>