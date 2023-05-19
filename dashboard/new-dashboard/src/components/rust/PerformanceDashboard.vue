<template>
  <DashboardPage
    v-slot="{serverConfigurator, dashboardConfigurators, averagesConfigurators, warnings}"
    db-name="perfint"
    table="rust"
    persistent-id="rust_dashboard"
    initial-machine="Linux EC2 m5d.xlarge or 5d.xlarge or m5ad.xlarge"
  >
    <section class="flex gap-6">
      <div class="w-1/2">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#Rust'"
          :title="'Indexing Rust (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        measure="indexing"
        :projects="['intelli-j-with-rust-test/test-rustling-cargo-sync', 'intelli-j-with-rust-test/run-ide-with-rust-plugin',
                    'intelli-j-with-rust-test/test-drogue-cloud-c-p-u-usage']"
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inspection execution time"
        measure="globalInspections"
        :projects="['intelli-j-with-rust-test/test-cargo-inspection']"
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inspection execution time"
        measure="globalInspections"
        :projects="['intelli-j-with-rust-test/test-cargo-inspection']"
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion"
        :projects="['intelli-j-with-rust-test/test-arrow-rs-completion', 'intelli-j-with-rust-test/test-vec-completion']"
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages"
        measure="findUsages"
        :projects="['intelli-j-with-rust-test/run-ide-with-rust-plugin-find-usages', 'intelli-j-with-rust-test/run-ide-with-rust-plugin-wasm-find-usages']"
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local Inspections (on file open)"
        measure="firstCodeAnalysis"
        :projects="[ 'intelli-j-with-rust-test/test-arrow-rs-highlighting', 'intelli-j-with-rust-test/test-cargo-highlighting', 
                     'intelli-j-with-rust-test/test-my-sql-async-highlighting']"
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Cargo Sync"
        measure="cargo_sync_execution_time"
        :projects="[ 'intelli-j-with-rust-test/test-rustling-cargo-sync']"
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

</script>