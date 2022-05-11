<template>
  <Divider align="center">
    {{ label }}
  </Divider>
  <div class="grid grid-cols-12 gap-4">
    <div class="col-span-12">
      <LineChartCard
        :compound-tooltip="true"
        :chart-type="'line'"
        :value-unit="'ms'"
        :measures="measures"
        :configurators="allConfigurators"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue"
import { dimensionConfigurator } from "../configurators/DimensionConfigurator"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { DataQueryConfigurator } from "../dataQuery"
import LineChartCard from "./LineChartCard.vue"


const props = defineProps<{
  label: string
  measures: Array<string>
  projects: Array<string>
  configurators: Array<DataQueryConfigurator>
  serverConfigurator: ServerConfigurator
}>()
const scenarioConfigurator = dimensionConfigurator("project", props.serverConfigurator, null, true)
const allConfigurators = props.configurators.concat(scenarioConfigurator)
onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
