<template>
  <DashboardPage
    db-name="perfint"
    table="rust"
    persistent-id="rust_rover_first_startup_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Global Inspections"
        :measure="['globalInspections']"
        :projects="['cargo/global-inspection', 'clap/global-inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Cargo Sync"
        measure="cargo_sync_execution_time"
        :projects="['arrow-rs/local-inspection/cast', 'arrow-rs/local-inspection/parse']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust Def Maps"
        measure="rust_def_maps_execution_time"
        :projects="['arrow-rs/local-inspection/cast', 'arrow-rs/local-inspection/parse']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust Macro Expansion"
        measure="rust_macro_expansion_execution_time"
        :projects="['arrow-rs/local-inspection/cast', 'arrow-rs/local-inspection/parse']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis"
        measure="firstCodeAnalysis"
        :projects="['arrow-rs/local-inspection/cast', 'arrow-rs/local-inspection/parse']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing latency max"
        measure="typing#latency#max"
        :projects="['arrow-rs/local-inspection/cast', 'arrow-rs/local-inspection/parse']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
