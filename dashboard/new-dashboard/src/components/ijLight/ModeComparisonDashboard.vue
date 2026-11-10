<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    :persistent-id="page.persistentId"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="allCharts"
    :initial-mode="MODES"
    :with-installer="false"
  >
    <template
      v-for="section in page.sections"
      :key="section.key"
    >
      <Divider
        :title="section.title"
        :description="section.description"
      />
      <section>
        <!-- scenario labels repeat across sections, so the key carries the section as well -->
        <GroupProjectsChart
          v-for="chart in section.charts"
          :key="`${section.key}-${chart.definition.label}`"
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
          :aliases="chart.aliases"
          :legend-formatter="legendFormatter"
          :y-axis-scale="page.yAxisScale"
        />
      </section>
    </template>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import { legendFormatter, MODE_COMPARISON_PAGES, ModeComparisonFamily, MODES } from "./modeComparisonTests"

// One page per test family; the route supplies the family (see ijLightModeComparison in routes.ts).
const { family } = defineProps<{ family: ModeComparisonFamily }>()

const page = MODE_COMPARISON_PAGES[family]
const allCharts = page.sections.flatMap((section) => section.charts)
</script>
