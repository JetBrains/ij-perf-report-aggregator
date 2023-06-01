<template>
  <DashboardPage
    v-slot="{averagesConfigurators}"
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
        :projects="['rustling/cargo-sync', 'yew/indexing',
                    'drogue-cloud/cpu']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Global Inspection execution time"
        measure="globalInspections"
        :projects="['cargo/global-inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion#mean_value"
        :projects="['arrow-rs/completion', 'vec/completion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages"
        measure="findUsages"
        :projects="['yew/find-usages', 'wasm/find-usages']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local Inspections (on file open)"
        measure="firstCodeAnalysis"
        :projects="[ 'arrow-rs/local-inspection', 'cargo/local-inspection',
                     'my-sql/local-inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local Inspections (on typing)"
        measure="typingCodeAnalyzing#mean_value"
        :projects="[ 'arrow-rs/local-inspection', 'cargo/local-inspection',
                     'my-sql/local-inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Cargo Sync"
        measure="cargo_sync_execution_time"
        :projects="[ 'cargo/local-inspection', 'my-sql/local-inspection',
                     'arrow-rs/local-inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust macro expansion time"
        measure="rust_macro_expansion_execution_time"
        :projects="[ 'cargo/local-inspection', 'my-sql/local-inspection',
                     'arrow-rs/local-inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust DefMaps execution time"
        measure="rust_def_maps_execution_time"
        :projects="[ 'cargo/local-inspection', 'my-sql/local-inspection',
                     'arrow-rs/local-inspection']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

</script>