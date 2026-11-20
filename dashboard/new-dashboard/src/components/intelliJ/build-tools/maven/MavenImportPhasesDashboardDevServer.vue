<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_maven_phases_dashboard_dev_server"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="testConfigurator"
        title="Test"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
    <section>
      <GroupProjectsChart
        v-for="project in testConfigurator.selected.value ?? []"
        :key="project"
        :label="project"
        :measure="MAVEN_SYNC_PHASES_WITH_TOTAL"
        :projects="[project]"
        stack="sync"
        :stacked-measures="MAVEN_SYNC_PHASES"
        :legend-formatter="(name: string) => (name === MAVEN_SYNC_TOTAL ? `${name} (= Total)` : name)"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { SimpleMeasureConfigurator } from "../../../../configurators/SimpleMeasureConfigurator"
import GroupProjectsChart from "../../../charts/GroupProjectsChart.vue"
import MeasureSelect from "../../../charts/MeasureSelect.vue"
import DashboardPage from "../../../common/DashboardPage.vue"
import { MAVEN_SYNC_PHASES, MAVEN_SYNC_PHASES_WITH_TOTAL, MAVEN_SYNC_TOTAL } from "./maven-metrics"
import { MAVEN_PROJECTS_FAST_INSTALLERS } from "./maven-projects"

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(MAVEN_PROJECTS_FAST_INSTALLERS)
</script>
