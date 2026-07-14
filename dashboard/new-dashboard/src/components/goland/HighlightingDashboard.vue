<template>
  <DashboardPage
    :with-installer="false"
    :with-mode="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_highlighting_dashboard"
    table="goland"
    branch="vietage/new-type-system"
  >
    <template #configurator>
      <DimensionSelect
        label="Engine"
        :dimension="engineConfigurator"
      >
        <template #icon>
          <AdjustmentsVerticalIcon class="w-4 h-4" />
        </template>
      </DimensionSelect>
      <DimensionSelect
        label="Metric"
        :dimension="metricConfigurator"
      >
        <template #icon>
          <FunnelIcon class="w-4 h-4" />
        </template>
      </DimensionSelect>
    </template>

    <template
      v-for="project in projects"
      :key="project.base"
    >
      <Divider :label="project.title" />
      <section>
        <GroupProjectsChart
          :label="project.title"
          :measure="measuresFor(project)"
          :projects="projectVariants(project.base)"
          :aliases="engineAliases"
        />
      </section>
    </template>

    <AdditionalMetrics :projects="allEnabledProjects" />
  </DashboardPage>
</template>

<script lang="ts" setup>
import { computed } from "vue"
import { DimensionConfigurator } from "../../configurators/DimensionConfigurator"
import DimensionSelect from "../charts/DimensionSelect.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import AdditionalMetrics from "./AdditionalMetrics.vue"

// Highlighting phases published per file by GoHighlightingTest (see runGoHighlightingScenario). Each is
// a duration in ms; the file label ("<speed>_<name>") is appended to build the full metric name.
const PHASES = ["coldStartHighlighting", "warmStartHighlighting", "typingHighlighting"]

// Highlighting-cost tiers each project measures one file at (ordered fast -> slow).
const SPEEDS = ["fast", "medium", "slow"]

// The 9 selectable metric types = phase x speed (e.g. "coldStartHighlighting_fast"). A type maps to a
// project's concrete measure by substituting that project's file for the tier (see measuresFor).
const METRIC_TYPES = PHASES.flatMap((phase) => SPEEDS.map((speed) => `${phase}_${speed}`))

// The two type-system engines the test compares. The engine is not a `mode` column value — it is encoded
// in the project name via GoHighlightingTest.Engine.suffix: LEGACY has no suffix, NEW_ONLY appends
// "-types2" (e.g. "caddy/highlighting" vs "caddy-types2/highlighting").
interface Engine {
  id: string
  suffix: string
}
const ENGINES: Engine[] = [
  { id: "LEGACY", suffix: "" },
  { id: "NEW_ONLY", suffix: "-types2" },
]

interface HighlightingProject {
  // Chart title.
  title: string
  // DB project base, before the engine suffix and the "/highlighting" scenario segment.
  base: string
  // The three scenario files (fast -> medium -> slow), as their metric-name labels ("<speed>_<name>").
  files: string[]
}

// Base project names and per-project file labels mirror the ScenarioFile lists in GoHighlightingTest.kt.
const projects: HighlightingProject[] = [
  { title: "kubernetes", base: "kubernetes", files: ["fast_clientTest", "medium_validation", "slow_generatedPb"] },
  { title: "mattermost", base: "mattermost-server", files: ["fast_webContext", "medium_userStore", "slow_opentracingLayer"] },
  { title: "cockroach", base: "cockroach", files: ["fast_storageParam", "medium_geoBuiltins", "slow_projNonConstOps"] },
  { title: "milvus", base: "milvus", files: ["fast_milvusClient", "medium_proxyImpl", "slow_dataCoordPb"] },
  { title: "rclone", base: "rclone", files: ["fast_httpBackend", "medium_azureblob", "slow_encoderCasesTest"] },
  { title: "volcano", base: "volcano", files: ["fast_queueCli", "medium_eventHandlers", "slow_allocateTest"] },
  { title: "caddy", base: "caddy", files: ["fast_importGraph", "medium_reverseProxyCaddyfile", "slow_httpType"] },
  { title: "k8sDevice", base: "k8sDevice", files: ["fast_resourceFactory", "medium_pciUtil", "slow_serverTest"] },
  { title: "fakeKub", base: "fake_kub", files: ["fast_clientTest", "medium_validation", "slow_generatedPb"] },
]

// Engine picker rendered in the toolbar next to branch/machine (same DimensionSelect combobox). It is
// driven entirely client-side: its selection only chooses which project variants each chart plots, so it
// is never added to a query's configurators.
const engineConfigurator = new DimensionConfigurator("engine", true)
engineConfigurator.values.value = ENGINES.map((engine) => engine.id)
engineConfigurator.selected.value = ENGINES.map((engine) => engine.id)
// createComponentState() starts disabled (it normally waits for a server load); enable it right away.
engineConfigurator.state.disabled = false

const selectedEngines = computed<Engine[]>(() => {
  const selected = engineConfigurator.selected.value
  const ids = Array.isArray(selected) ? selected : selected == null ? [] : [selected]
  // Keep the canonical LEGACY-before-NEW_ONLY order regardless of the order they were clicked.
  return ENGINES.filter((engine) => ids.includes(engine.id))
})

// The selected engine variants of a project, e.g. ["caddy/highlighting", "caddy-types2/highlighting"].
function projectVariants(base: string): string[] {
  return selectedEngines.value.map((engine) => `${base}${engine.suffix}/highlighting`)
}

// Legend label for each engine series, parallel (same order) to projectVariants().
const engineAliases = computed(() => selectedEngines.value.map((engine) => engine.id))

const allEnabledProjects = computed(() => projects.flatMap((project) => projectVariants(project.base)))

// Metric-type picker (second toolbar combobox, same DimensionSelect as Engine). Also client-side: its
// selection chooses which of the 9 phase x speed metric types each chart plots.
const metricConfigurator = new DimensionConfigurator("metric", true)
metricConfigurator.values.value = METRIC_TYPES
metricConfigurator.selected.value = [...METRIC_TYPES]
metricConfigurator.state.disabled = false

const selectedMetricTypes = computed<string[]>(() => {
  const selected = metricConfigurator.selected.value
  const ids = Array.isArray(selected) ? selected : selected == null ? [] : [selected]
  // Keep the canonical METRIC_TYPES order regardless of the order they were clicked.
  return METRIC_TYPES.filter((type) => ids.includes(type))
})

// Concrete measures for a project: each selected "<phase>_<speed>" type becomes "<phase>_<projectFile>",
// substituting the project's file for that tier (e.g. "coldStartHighlighting_fast" -> ..._fast_importGraph).
function measuresFor(project: HighlightingProject): string[] {
  return selectedMetricTypes.value.map((type) => {
    const phase = PHASES.find((p) => type.startsWith(`${p}_`)) ?? ""
    const speed = phase === "" ? "" : type.slice(phase.length + 1)
    const file = project.files.find((f) => f.startsWith(`${speed}_`)) ?? ""
    return `${phase}_${file}`
  })
}
</script>
