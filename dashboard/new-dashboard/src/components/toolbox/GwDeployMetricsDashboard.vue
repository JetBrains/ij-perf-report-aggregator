<template>
  <DashboardPage
    db-name="toolbox"
    table="report"
    persistent-id="gw_deploy_metrics_dashboard"
  >
    <section>
      <GroupProjectsChart
        label="GW: open recent with backend stopped"
        :measure="backendStoppedMetrics"
        :projects="['GW: open recent with backend stopped']"
      />
      <GroupProjectsChart
        label="GW: open recent with backend running"
        :measure="backendRunningMetrics"
        :projects="['GW: open recent with backend running']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const backendStoppedMetrics = [
  "gw.cold.start.from.gw.link.clicked.to.editors.opened",
  "gw.cold.start.from.gw.link.clicked.to.host.first.metric",
  "gw.cold.start.from.host.first.metric.to.joinLink.ready",
  "gw.cold.start.from.joinLink.ready.to.client.first.metric",
  "gw.cold.start.from.client.first.metric.to.connection.established",
  "gw.cold.start.from.connection.established.to.editors.opened",
]

const backendRunningMetrics = [
  "gw.warm.start.from.gw.link.clicked.to.editors.opened",
  "gw.warm.start.from.gw.link.clicked.to.client.first.metric",
  "gw.warm.start.from.client.first.metric.to.connection.established",
  "gw.warm.start.from.connection.established.to.editors.opened",
]
</script>
