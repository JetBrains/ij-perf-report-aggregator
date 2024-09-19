<template>
  <div class="pt-3 pb-1 px-5 border border-solid rounded-md border-zinc-200">
    <h3 class="m-0 mb-3">
      {{ title }}
    </h3>

    <div class="mb-3 flex flex-col text-center">
      <span class="text-3xl font-bold">
        {{ vm.average }}
      </span>
      <div class="text-sm text-neutral-500 font-normal">
        Avg
        <span v-if="valueUnit !== 'counter'">, ms</span>
      </div>
    </div>

    <div
      ref="element"
      class="bg-white"
      :style="{ height: `${55}px` }"
    />
  </div>
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, useTemplateRef } from "vue"
import { TimeAverageConfigurator } from "../../configurators/TimeAverageConfigurator"
import { containerKey, dashboardConfiguratorsKey, serverConfiguratorKey } from "../../shared/keys"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ValueUnit } from "../common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { AggregationChartVM } from "./AggregationChartVM"
import { injectOrError } from "../../shared/injectionKeys"

interface AggregationChartProps {
  valueUnit?: ValueUnit
  chartColor?: string
  configurators?: DataQueryConfigurator[]
  aggregatedMeasure: string
  aggregatedProject?: string
  isLike?: boolean
  title: string
}

const { valueUnit = "ms", chartColor = "#4B84EE", configurators, aggregatedMeasure, aggregatedProject, isLike, title } = defineProps<AggregationChartProps>()
const timeAverageConfigurator = new TimeAverageConfigurator()
const measuresConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({ f: "measures.name", v: aggregatedMeasure, o: isLike ? "like" : "=" })
    if (aggregatedProject !== undefined) {
      query.addFilter({ f: "project", v: aggregatedProject, o: isLike ? "like" : "=" })
    }
    return true
  },
  createObservable() {
    return null
  },
}
const avConfigurators = configurators ?? injectOrError(dashboardConfiguratorsKey)
const allConfigurators = [...avConfigurators, injectOrError(serverConfiguratorKey), timeAverageConfigurator, measuresConfigurator]
/* eslint-disable-next-line @typescript-eslint/no-unnecessary-condition */
const queryExecutor = new DataQueryExecutor(allConfigurators.filter((item): item is DataQueryConfigurator => item != null))
const element = useTemplateRef<HTMLElement>("element")
const vm = new AggregationChartVM(queryExecutor, chartColor, valueUnit)
const container = inject(containerKey)

let dispose: () => void
onMounted(() => {
  if (element.value != null) {
    dispose = vm.initChart(element.value, container?.value)
  }
})

onUnmounted(() => {
  dispose()
})
</script>
