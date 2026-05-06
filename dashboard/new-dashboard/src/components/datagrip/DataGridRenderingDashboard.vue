<template>
  <DashboardPage
    db-name="perfintDev"
    table="datagrip"
    persistent-id="datagrip_data_grid_rendering_dashboard"
    initial-machine="Linux EC2 m5d.xlarge (4 vCPU Xeon, 16 GB)"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Base Path (selection off)"
        :measure="[
          'grid.rendering.base.path.old.selection.off.attempt.mean.ms',
          'grid.rendering.base.path.old.selection.off.attempt.median.ms',
          'grid.rendering.base.path.old.selection.off.attempt.min.ms',
          'grid.rendering.base.path.old.selection.off.attempt.range.ms',
        ]"
        :projects="['data-grid-rendering-performance-test/render-cells-with-long-strings-local']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Expansion Hint Path (selection off)"
        :measure="[
          'grid.rendering.expansion.hint.path.old.selection.off.attempt.mean.ms',
          'grid.rendering.expansion.hint.path.old.selection.off.attempt.median.ms',
          'grid.rendering.expansion.hint.path.old.selection.off.attempt.min.ms',
          'grid.rendering.expansion.hint.path.old.selection.off.attempt.range.ms',
        ]"
        :projects="['data-grid-rendering-performance-test/render-cells-with-long-strings-local']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
