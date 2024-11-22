<template>
  <DashboardPage
    db-name="fleet"
    table="measure_new"
    persistent-id="fleetPerfStartupComparision_dashboard"
    :with-installer="false"
    :initial-machine="null"
  >
    <Divider title="Windows" />
    <section>
      <GroupProjectsChart
        label="Open File"
        :measure="['test.startup.vscode.open-file', 'test.startup.fleet.open-file']"
        :machines="['Windows Fleet AWS UI']"
        :projects="['metrics']"
      />
      <GroupProjectsChart
        label="Open Window"
        :measure="['test.startup.vscode.open-window', 'test.startup.fleet.open-window']"
        :machines="['Windows Fleet AWS UI']"
        :projects="['metrics']"
      />
    </section>

    <Divider title="macOS" />
    <section>
      <GroupProjectsChart
        label="Open File"
        :measure="['test.startup.vscode.open-file', 'test.startup.fleet.open-file', 'test.startup.zed.open-file']"
        :machines="['Mac Fleet AWS UI']"
        :projects="['metrics']"
      />
      <GroupProjectsChart
        label="Open Window"
        :measure="['test.startup.vscode.open-window', 'test.startup.fleet.open-window', 'test.startup.zed.open-window']"
        :machines="['Mac Fleet AWS UI']"
        :projects="['metrics']"
      />
    </section>

    <Divider title="Linux" />
    <section>
      <GroupProjectsChart
        label="Open File"
        :measure="['test.startup.vscode.open-file', 'test.startup.fleet.open-file']"
        :machines="['Linux Fleet AWS UI']"
        :projects="['metrics']"
      />
      <GroupProjectsChart
        label="Open Window"
        :measure="['test.startup.vscode.open-window', 'test.startup.fleet.open-window']"
        :machines="['Linux Fleet AWS UI']"
        :projects="['metrics']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
