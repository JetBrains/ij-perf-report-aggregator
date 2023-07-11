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
          :aggregated-measure="'processingSpeed#C++'"
          :title="'Indexing C++ (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#C++ Header'"
          :title="'Indexing C++ Header (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :is-like="true"
          :title="'Completion'"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        :measure="['indexing', 'indexingTimeWithoutPauses']"
        :projects="['curl/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files"
        measure="numberOfIndexedFiles"
        :projects="['curl/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning"
        :measure="['scanning', 'scanningTimeWithoutPauses']"
        :projects="['curl/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion"
        :projects="['fmtlib/completion']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
