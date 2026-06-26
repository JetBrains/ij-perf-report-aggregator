<template>
  <LineChart
    v-if="renderMode === 'charts'"
    :title="label"
    :value-unit="valueUnit"
    :measures="measureArray"
    :configurators="configurators"
    :skip-zero-values="false"
    :legend-formatter="legendFormatter"
    :can-be-closed="canBeClosed"
    :description="description"
    :better-direction="betterDirection"
    @chart-closed="onChartClosed"
  />
</template>

<script setup lang="ts">
import { computed, inject, onUnmounted, Ref, ref, watch } from "vue"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { injectOrError } from "../../shared/injectionKeys"
import { compareSectionsRegistryKey, dashboardConfiguratorsKey, renderModeKey, serverConfiguratorKey } from "../../shared/keys"
import { removeCommonSegments } from "../../util/removeCommonPrefixes"
import { ValueUnit } from "../common/chart"
import { BetterDirection } from "../../shared/changeDetector/algorithm"
import LineChart from "./LineChart.vue"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { CompareSectionConfig, RenderMode } from "./compareMode"

interface Props {
  label: string
  measure: string | string[]
  projects: string[]
  machines?: string[] | null
  valueUnit?: ValueUnit
  legendFormatter?: (name: string) => string
  aliases?: string[] | null
  canBeClosed?: boolean
  description?: string
  betterDirection?: BetterDirection
}

const {
  label,
  measure,
  projects,
  machines = null,
  valueUnit = "auto",
  legendFormatter = (name: string) => name,
  aliases = null,
  canBeClosed = false,
  description,
  betterDirection = "lower",
} = defineProps<Props>()

const serverConfigurator = injectOrError(serverConfiguratorKey)
const dashboardConfigurators = injectOrError(dashboardConfiguratorsKey)
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, null, true)
const configurators = [...dashboardConfigurators, scenarioConfigurator, serverConfigurator]

if (machines != null) {
  const machineConfigurator = new MachineConfigurator(serverConfigurator, undefined, [], true, machines)
  configurators.push(machineConfigurator)
}

const measureArray: Ref<string[]> = computed(() => {
  return Array.isArray(measure) ? measure : [measure]
})

const emit = defineEmits<{
  chartClosed: [projects: string[]]
}>()

function onChartClosed(): void {
  emit("chartClosed", projects)
}

watch(
  () => [projects, aliases],
  ([projects, aliases]) => {
    if (projects != null) {
      const aliasesToUse = aliases ?? removeCommonSegments(projects)
      scenarioConfigurator.aliases = new Map(projects.map((key, index) => [key, aliasesToUse[index]]))
    }
    scenarioConfigurator.selected.value = projects
  },
  { immediate: true }
)

const renderMode = inject(renderModeKey, ref<RenderMode>("charts"))
const compareRegistry = inject(compareSectionsRegistryKey, null)

const sectionId = compareRegistry?.nextSectionId() ?? null

function buildSectionConfig(): CompareSectionConfig | null {
  if (sectionId == null) return null
  return {
    id: sectionId,
    label,
    measure,
    projects,
    aliases,
    machines,
    valueUnit,
  }
}

if (compareRegistry != null && sectionId != null) {
  const initial = buildSectionConfig()
  if (initial != null) compareRegistry.register(initial)

  watch(
    () => [label, measure, projects, aliases, machines, valueUnit],
    () => {
      const config = buildSectionConfig()
      if (config != null) compareRegistry.register(config)
    },
    { deep: true }
  )

  onUnmounted(() => {
    compareRegistry.unregister(sectionId)
  })
}
</script>
