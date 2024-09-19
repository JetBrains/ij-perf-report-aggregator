<template>
  <Select
    v-model="value"
    title="Aggregate by"
    :options="operators"
    option-label="label"
    option-value="value"
    placeholder="Operator"
  >
    <!-- eslint-disable vue/no-template-shadow -->
    <template #value="{ value }">
      <div class="group inline-flex justify-center text-sm font-medium">
        {{ value }}
        <ChevronDownIcon
          class="-mr-1 ml-1 h-5 w-5 flex-shrink-0"
          aria-hidden="true"
        />
      </div>
    </template>
    <template #dropdownicon>
      <!-- empty element to avoid ignoring override of slot -->
      <span />
    </template>
  </Select>
  <InputNumber
    v-if="value === 'quantile'"
    v-model="quantile"
    size="3"
    :min="0"
    :max="100"
    :step="10"
    :show-buttons="true"
  />
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { computed, inject } from "vue"
import { AggregationOperatorConfigurator } from "../../configurators/AggregationOperatorConfigurator"
import { aggregationOperatorConfiguratorKey } from "../../shared/injectionKeys"

const { configurator } = defineProps<{
  configurator?: AggregationOperatorConfigurator
}>()

const providedConfigurator = inject(aggregationOperatorConfiguratorKey, null)

function getConfigurator(): AggregationOperatorConfigurator {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  return configurator ?? providedConfigurator!
}

const operators = ["median", "min", "max", "quantile"].map((it) => ({ label: it, value: it }))
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
