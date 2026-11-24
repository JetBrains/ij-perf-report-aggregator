<template>
  <DimensionSelect
    label="Configuration"
    :dimension="configurations"
    :value-to-group="platformOf"
    :value-to-label="configurationOf"
    :selected-label="(items) => `${items.length} configurations`"
  >
    <template #icon>
      <ComputerDesktopIcon class="w-4 h-4" />
    </template>
  </DimensionSelect>
</template>

<script setup lang="ts">
import { watch } from "vue"
import { DataQuery } from "../common/dataQuery"
import DimensionSelect from "../charts/DimensionSelect.vue"
import { dimensionConfigurator, selectedToArray } from "../../configurators/DimensionConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { injectOrError } from "../../shared/injectionKeys"
import { dashboardConfiguratorsKey, serverConfiguratorKey } from "../../shared/keys"
import { compareConfigurations, configurationOf, platformOf } from "./platforms"

/**
 * The configurations of one JBR benchmark (projects "<benchmark>_<configuration>") that the charts show, offered as
 * they appear in the database for the chosen time range and branch, grouped by platform. Configurations that stopped
 * reporting drop out of the list, and new ones join it without a change here.
 */
const selected = defineModel<string[]>({ required: true })
const { benchmark } = defineProps<{ benchmark: string }>()

const benchmarkFilter: FilterConfigurator = {
  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "project", v: benchmark + "_%", o: "like" })
    return true
  },
  createObservable() {
    return null
  },
}

const configurations = dimensionConfigurator(
  "project",
  injectOrError(serverConfiguratorKey),
  null,
  true,
  [...(injectOrError(dashboardConfiguratorsKey) as FilterConfigurator[]), benchmarkFilter],
  compareConfigurations
)
configurations.selected.value = selected.value
watch(configurations.selected, (value) => {
  selected.value = selectedToArray(value)
})
</script>
