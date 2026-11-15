<template>
  <DashboardPage
    :with-installer="false"
    db-name="perfintDev"
    history-start="2026-11-15"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :initial-mode="['default', 'types2']"
    persistent-id="goland_highlighting_dashboard"
    table="goland"
  >
    <template #configurator>
      <DimensionSelect
        label="Metric"
        :dimension="metricConfigurator"
      >
        <template #icon>
          <FunnelIcon class="w-4 h-4" />
        </template>
      </DimensionSelect>
      <DimensionSelect
        label="Quantity"
        :dimension="quantityConfigurator"
        :value-to-label="quantityLabel"
      >
        <template #icon>
          <ScaleIcon class="w-4 h-4" />
        </template>
      </DimensionSelect>
    </template>

    <!-- LEGACY-vs-candidate verdict is the primary content; the per-project line charts are demoted below. -->
    <EngineComparison
      :selected-metric-types="selectedMetricTypes"
      :quantity="selectedQuantity"
    />

    <ChartAccordion :lazy="true">
      <AccordionPanel value="0">
        <AccordionHeader>Per-project history</AccordionHeader>
        <AccordionContent>
          <div class="flex flex-col gap-12">
            <section
              v-for="project in projects"
              :key="project.base"
              class="flex flex-col gap-4"
            >
              <!-- Left-aligned separator: the project title, then a rule filling the remaining width. -->
              <div class="flex items-center gap-3">
                <span class="text-lg font-medium whitespace-nowrap">{{ project.title }}</span>
                <div class="grow border-t border-gray-200 dark:border-gray-700" />
              </div>
              <!-- One chart per (phase, bucket): phases (cold/warm/typing) run down the rows, buckets
                   (fast/medium/slow) across the columns. Deselected metric types leave their cell empty
                   so the columns stay aligned; a fully deselected row or column is dropped. -->
              <div
                class="grid gap-4 overflow-x-auto"
                :style="{ gridTemplateColumns: `repeat(${visibleBuckets.length}, minmax(280px, 1fr))` }"
              >
                <template
                  v-for="phase in visiblePhases"
                  :key="phase"
                >
                  <div
                    v-for="bucket in visibleBuckets"
                    :key="`${phase}_${bucket}`"
                  >
                    <GroupProjectsChart
                      v-if="isTypeSelected(phase, bucket)"
                      :label="metricTypeLabel(`${phase}_${bucket}`)"
                      :measure="buildMeasure(project, `${phase}_${bucket}`, selectedQuantity.suffix)"
                      :projects="[dbProject(project.base)]"
                    />
                  </div>
                </template>
              </div>
            </section>
          </div>
        </AccordionContent>
      </AccordionPanel>
    </ChartAccordion>

    <AdditionalMetrics :projects="allEnabledProjects" />
  </DashboardPage>
</template>

<script lang="ts" setup>
import { computed } from "vue"
import { DimensionConfigurator, selectedToArray } from "../../configurators/DimensionConfigurator"
import DimensionSelect from "../charts/DimensionSelect.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import ChartAccordion from "../charts/ChartAccordion.vue"
import DashboardPage from "../common/DashboardPage.vue"
import AccordionPanel from "openvue/accordionpanel"
import AccordionHeader from "openvue/accordionheader"
import AccordionContent from "openvue/accordioncontent"
import AdditionalMetrics from "./AdditionalMetrics.vue"
import EngineComparison from "./engineComparison/EngineComparison.vue"
import { buildMeasure, metricTypeLabel, PHASES, projects, Quantity } from "./engineComparison/highlightingMetrics"

// Dashboard-only chart settings: which metric types the toolbar offers and the quantity choices. The engine is the
// ij-perf mode of the dashboard mode switch (see engineOfMode). The shared highlighting config (phases, projects,
// measure building, labels) lives in engineComparison/highlightingMetrics.

// Highlighting-cost buckets each project measures one file at (ordered fast -> slow).
const BUCKETS = ["fast", "medium", "slow"]

// The 9 selectable metric types = phase x bucket (e.g. "coldStartHighlighting_fast"). A type maps to a
// project's concrete measure by substituting that project's file for the bucket (see buildMeasure).
const METRIC_TYPES = PHASES.flatMap((phase) => BUCKETS.map((bucket) => `${phase}_${bucket}`))

const QUANTITIES: Quantity[] = [
  { id: "duration", label: "Duration", suffix: "" },
  { id: "allocated", label: "Allocated", suffix: "#jvm.alloc.mb" },
  { id: "cpu", label: "CPU time", suffix: "#jvm.cpu.time.ms" },
  { id: "gcTime", label: "GC time", suffix: "#jvm.gc.time.ms" },
  { id: "gcCount", label: "GC count", suffix: "#jvm.gc.count" },
]

function dbProject(base: string): string {
  return `${base}/highlighting`
}

const allEnabledProjects = projects.map((project) => dbProject(project.base))

// Metric-type picker (toolbar combobox). Client-side: its
// selection chooses which of the 9 phase x bucket metric types the comparison and charts plot.
const metricConfigurator = new DimensionConfigurator("metric", true)
metricConfigurator.values.value = METRIC_TYPES
metricConfigurator.selected.value = [...METRIC_TYPES]
metricConfigurator.state.disabled = false

const selectedMetricTypes = computed<string[]>(() => {
  const ids = selectedToArray(metricConfigurator.selected.value)
  // Keep the canonical METRIC_TYPES order regardless of the order they were clicked.
  return METRIC_TYPES.filter((type) => ids.includes(type))
})

// Quantity picker (single-select): duration or a per-scenario JVM sub-metric. Its suffix is appended to
// each measure name; units and better-direction resolve from the declared metricsDescription entries.
const quantityConfigurator = new DimensionConfigurator("quantity", false)
quantityConfigurator.values.value = QUANTITIES.map((quantity) => quantity.id)
quantityConfigurator.selected.value = "duration"
quantityConfigurator.state.disabled = false

const selectedQuantity = computed<Quantity>(() => {
  const selected = quantityConfigurator.selected.value
  const id = Array.isArray(selected) ? selected[0] : selected
  return QUANTITIES.find((quantity) => quantity.id === id) ?? QUANTITIES[0]
})

function quantityLabel(id: string): string {
  return QUANTITIES.find((quantity) => quantity.id === id)?.label ?? id
}

// Drill-down grid: a (phase, bucket) cell is shown only when its metric type is selected. Rows are the
// phases that still have a selected bucket; columns the buckets that still have a selected phase.
function isTypeSelected(phase: string, bucket: string): boolean {
  return selectedMetricTypes.value.includes(`${phase}_${bucket}`)
}

const visiblePhases = computed(() => PHASES.filter((phase) => BUCKETS.some((bucket) => isTypeSelected(phase, bucket))))
const visibleBuckets = computed(() => BUCKETS.filter((bucket) => PHASES.some((phase) => isTypeSelected(phase, bucket))))
</script>
