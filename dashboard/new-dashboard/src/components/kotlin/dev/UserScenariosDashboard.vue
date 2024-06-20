<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlin_scenarios_dashboard_dev"
    :with-installer="false"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="KOTLIN_SCENARIO_CONFIGURATOR"
        title="Project"
        :selected-label="scenarioSelectedLabel"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </MeasureSelect>
    </template>
    <SlackLink></SlackLink>
    <div
      v-for="(label, index) in KOTLIN_SCENARIO_CONFIGURATOR.selected.value"
      :key="index"
    >
      <Divider :title="label" />
      <K1K2DashboardGroupCharts :definitions="Object.values(USER_SCENARIOS).find((d) => d.label == label)?.charts?.value" />
    </div>
  </DashboardPage>
</template>

<script setup lang="ts">
import K1K2DashboardGroupCharts from "../K1K2DashboardGroupCharts.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import Divider from "../../common/Divider.vue"
import { KOTLIN_SCENARIO_CONFIGURATOR, USER_SCENARIOS } from "../projects"
import SlackLink from "../SlackLink.vue"
import MeasureSelect from "../../charts/MeasureSelect.vue"
import { scenarioSelectedLabel } from "../label-formatter"
</script>
