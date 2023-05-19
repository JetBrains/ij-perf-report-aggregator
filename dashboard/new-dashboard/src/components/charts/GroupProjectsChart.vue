<template>
  <LineChart
    :title="props.label"
    :value-unit="props.valueUnit"
    :measures="[measure]"
    :configurators="configurators"
    :skip-zero-values="false"
    :accidents="accidents"
  />
</template>

<script setup lang="ts">
import { ValueUnit } from "shared/src/chart"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { FilterConfigurator } from "shared/src/configurators/filter"
import { DataQueryConfigurator } from "shared/src/dataQuery"
import { inject, onMounted } from "vue"
import { accidentsKeys, serverConfiguratorKey } from "../../shared/keys"
import LineChart from "./LineChart.vue"

interface Props {
  label: string
  measure: string
  projects: Array<string>
  configurators: Array<DataQueryConfigurator>
  valueUnit?: ValueUnit
}

const props = withDefaults(defineProps<Props>(), {
  valueUnit: "ms",
})

const serverConfigurator = inject(serverConfiguratorKey) as ServerConfigurator
const accidents = inject(accidentsKeys)

const scenarioConfigurator = dimensionConfigurator(
  "project", 
  serverConfigurator,
  null, 
  true,
  [...props.configurators] as Array<FilterConfigurator>
)
const configurators = [...props.configurators, scenarioConfigurator, serverConfigurator]

onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
