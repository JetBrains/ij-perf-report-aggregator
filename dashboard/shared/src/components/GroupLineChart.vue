<template>
  <div class="w-full">
    <LineChartCard
      :label="label"
      :compound-tooltip="true"
      :chart-type="'line'"
      :value-unit="props.valueUnit"
      :measures="[measure]"
      :configurators="configurators"
      :skip-zero-values="false"
      trigger="item"
    />
  </div>
</template>

<script setup lang="ts">
import { Observable } from "rxjs"
import { inject, onMounted } from "vue"
import { dimensionConfigurator } from "../configurators/DimensionConfigurator"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { FilterConfigurator } from "../configurators/filter"
import { DataQuery } from "../dataQuery"
import { configuratorListKey } from "../injectionKeys"
import LineChartCard from "./LineChartCard.vue"

const props = withDefaults(defineProps<{
  label: string
  measure: string
  projects: Array<string>
  serverConfigurator: ServerConfigurator
  valueUnit?: "ns"|"ms"
}>(), {
  valueUnit: "ms"
})
const providedConfigurators = inject(configuratorListKey, null)
if (providedConfigurators == null) {
  throw new Error("`dataQueryExecutor` is not provided")
}
class ProjectFilter implements FilterConfigurator{
  createObservable(): Observable<unknown> | null {
    return null
  }

  configureFilter(query: DataQuery): boolean {
    return query.addFilter({f: "project", v: props.projects})
  }
}
const scenarioConfigurator = dimensionConfigurator("project", props.serverConfigurator, null, true, [...providedConfigurators, new ProjectFilter()])
const configurators = [...providedConfigurators, scenarioConfigurator]
onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
