<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_distribution_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Distributive Size (in Kb, unpacked)"
        measure="unpackedDistributiveSize"
        :projects="['distributive-size-GO']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
