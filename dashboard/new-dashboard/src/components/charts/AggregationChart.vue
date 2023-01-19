<template>
  <div class="pt-3 pb-1 px-5 border border-solid rounded-md border-zinc-200">
    <h3 class="m-0 text-sm mb-3">
      {{ props.title }}
    </h3>

    <div class="mb-3 flex flex-col text-center">
      <span class="text-3xl text-black font-bold">
        {{ vm.average }}
      </span>
      <div class="text-sm text-neutral-500 font-normal">
        Avg<div v-if="props.valueUnit !== 'counter'">, ms</div>
      </div>
    </div>

    <div
      ref="element"
      class="bg-white"
      :style="{height: `${55}px`}"
    />
  </div>
</template>
<script setup lang="ts">
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ValueUnit } from "shared/src/chart"
import { TimeAverageConfigurator } from "shared/src/configurators/TimeAverageConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import { inject, onMounted, onUnmounted, ref, shallowRef, withDefaults } from "vue"
import { containerKey } from "../../shared/keys"
import { AggregationChartVM } from "./AggregationChartVM"

interface AggregationChartProps {
  valueUnit?: ValueUnit
  chartColor?: string
  configurators: Array<DataQueryConfigurator>
  aggregatedMeasure: string
  isLike?: boolean
  title: string
}

const props = withDefaults(defineProps<AggregationChartProps>(), {
  valueUnit: "ms",
  chartColor: "#4B84EE",
})
const timeAverageConfigurator = new TimeAverageConfigurator()
const measuresConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({f: "measures.name", v: props.aggregatedMeasure, o: props.isLike ? "like" : "="})
    return true
  },
  createObservable() {
    return null
  },
}
const configurators = [...props.configurators, timeAverageConfigurator, measuresConfigurator]
const queryExecutor = new DataQueryExecutor(configurators)
const element = shallowRef<HTMLElement>()
const vm = new AggregationChartVM(queryExecutor, props.chartColor)
const container = inject(containerKey)

let dispose: () => void
onMounted(() => {
  dispose = vm.initChart(element.value!, container?.value)
})

onUnmounted(() => {
  dispose?.()
})
</script>