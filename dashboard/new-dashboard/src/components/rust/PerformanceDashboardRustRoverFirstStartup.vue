<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="rust"
    persistent-id="rust_rover_first_startup_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    release-configurator="EAP / Release"
  >
    <section class="flex gap-6">
      <div class="w-1/2">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeedAvg#Rust'"
          :title="'Indexing Rust (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          :measure="['indexingTimeWithoutPauses']"
          :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          :measure="['scanningTimeWithoutPauses']"
          :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
        />
      </div>
    </section>

    <section>
      <GroupProjectsChart
        label="Duration from Start to Work (metric 'rust_duration_from_start_to_work')"
        measure="rust_duration_from_start_to_work"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Warm Duration from Start to Work (first run)"
          :measure="['rust_duration_from_start_to_work_warm_run_1']"
          :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Warm Duration from Start to Work (second run)"
          :measure="['rust_duration_from_start_to_work_warm_run_2']"
          :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Duration from Start to Cargo Sync (metric 'rust_duration_from_start_to_cargo_sync')"
        measure="rust_duration_from_start_to_cargo_sync"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Full Cargo Sync (metric 'cargo_sync_execution_time')"
        measure="cargo_sync_execution_time"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Cargo Metadata (metric 'rust_cargo_metadata_time')"
        measure="rust_cargo_metadata_time"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Buildscript evaluation (metric 'rust_buildscript_evaluation_time')"
        measure="rust_buildscript_evaluation_time"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Stdlib fetching (metric 'rust_stdlib_fetching_time')"
        measure="rust_stdlib_fetching_time"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust CrateDefMaps build (metric 'rust_def_maps_execution_time')"
        measure="rust_def_maps_execution_time"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust macro expansions saving to VFS (metric 'rust_macro_expansion_execution_time')"
        measure="rust_macro_expansion_execution_time"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust CrateDefMaps size in MB (metric 'rust_class_instances_tree_size_mb#org.rust.lang.core.resolve2.CrateDefMap')"
        measure="rust_class_instances_tree_size_mb#org.rust.lang.core.resolve2.CrateDefMap"
        :projects="rustGlobalInspectionProjects.map((project) => `${project}/indexing`)"
        value-unit="counter"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { rustGlobalInspectionProjects } from "./RustTestCases"
</script>
