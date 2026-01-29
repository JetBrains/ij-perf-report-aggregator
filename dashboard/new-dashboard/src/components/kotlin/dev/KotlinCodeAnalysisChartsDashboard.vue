<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_codeAnalysis_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <ConfiguratorRegistration
      :configurator="projectConfigurator"
      :data="Object.values(PROJECT_CATEGORIES).flatMap((c) => c.label)"
    />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import K1K2DashboardGroupCharts from "../K1K2DashboardGroupCharts.vue"
import { createKotlinCharts, PROJECT_CATEGORIES } from "../projects"
import { SimpleMeasureConfigurator } from "../../../configurators/SimpleMeasureConfigurator"
import ConfiguratorRegistration from "../ConfiguratorRegistration.vue"

const projectConfigurator = new SimpleMeasureConfigurator("project", null)
const { codeAnalysisCharts } = createKotlinCharts(projectConfigurator)
</script>
