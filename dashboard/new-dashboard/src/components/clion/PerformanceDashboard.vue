<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="clion"
    persistent-id="clion_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'clion/%'"
          :aggregated-measure="'processingSpeed#%'"
          :is-like="true"
          :title="'[CLion] Indexing speed (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'clion/%hot%'"
          :aggregated-measure="'completion#mean\_value'"
          :is-like="true"
          :title="'[CLion] Completion'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'clion/%'"
          :aggregated-measure="'%#average_awt_delay'"
          :is-like="true"
          :title="'[CLion] UI responsiveness during typing'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'radler/%'"
          :aggregated-measure="'processingSpeed#%'"
          :is-like="true"
          :title="'[Radler] Indexing speed (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'radler/%hot%'"
          :aggregated-measure="'completion#mean\_value'"
          :is-like="true"
          :title="'[Radler] Completion'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'radler/%'"
          :aggregated-measure="'%#average_awt_delay'"
          :is-like="true"
          :title="'[Radler] UI responsiveness during typing'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>
    <!--<section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing (curl)"
          :measure="['indexingTimeWithoutPauses']"
          :projects="['clion/curl/indexing', 'radler/curl/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning (curl)"
          :measure="['scanningTimeWithoutPauses']"
          :projects="['clion/curl/indexing', 'radler/curl/indexing']"
        />
      </div>
    </section>-->
    <section>
      <GroupProjectsChart
        label="Global Inspections (fmtlib)"
        measure="globalInspections"
        :projects="['clion/fmtlib/inspection', 'radler/fmtlib/inspection']"
      />
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion, mean value (std::string, hot)"
          measure="completion#mean_value"
          :projects="['clion/fmtlib/completion/std.string (hot)', 'radler/fmtlib/completion/std.string (hot)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion, mean value (std::shared_ptr<T>, hot)"
          measure="completion#mean_value"
          :projects="['clion/fmtlib/completion/std.shared_ptr (dep) (hot)', 'radler/fmtlib/completion/std.shared_ptr (dep) (hot)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion, mean value (fmt::join<It>, hot)"
          measure="completion#mean_value"
          :projects="['clion/fmtlib/completion/fmt.join_view (dep) (hot)', 'radler/fmtlib/completion/fmt.join_view (dep) (hot)']"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
