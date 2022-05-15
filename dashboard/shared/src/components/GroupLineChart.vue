<template>
  <Divider :label="label" />
  <div class="grid grid-cols-12 gap-4">
    <div class="col-span-12">
      <LineChartCard
        :compound-tooltip="true"
        :chart-type="'line'"
        :value-unit="'ms'"
        :measures="[measure]"
        :configurators="configurators"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import Divider from "tailwind-ui/src/Divider.vue"
import { inject, onMounted } from "vue"
import { dimensionConfigurator } from "../configurators/DimensionConfigurator"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { configuratorListKey } from "../injectionKeys"
import LineChartCard from "./LineChartCard.vue"

const props = defineProps<{
  label: string
  measure: string
  projects: Array<string>
  serverConfigurator: ServerConfigurator
}>()
const providedConfigurators = inject(configuratorListKey, null)
if (providedConfigurators == null) {
  throw new Error("`dataQueryExecutor` is not provided")
}
const scenarioConfigurator = dimensionConfigurator("project", props.serverConfigurator, null, true)
const configurators = providedConfigurators.concat(scenarioConfigurator)
onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
