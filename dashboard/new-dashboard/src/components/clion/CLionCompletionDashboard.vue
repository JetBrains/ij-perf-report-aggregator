<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_completion_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'radler/%hot%'"
          :aggregated-measure="'fus_time_to_show_90p'"
          :is-like="true"
          :title="'Time to show completion list'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'radler/%/typing/%'"
          :aggregated-measure="'typing#latency#max'"
          :is-like="true"
          :title="'Typing latency(max)'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>

    <Divider title="Basic Completion" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Basic: Time to show (90p) (cold)`"
            :measure="['fus_time_to_show_90p']"
            :projects="['radler/fmtlib/completion/std.string (cold)']"
          />
        </section>
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Basic: Time to show (90p)`"
            :measure="['fus_time_to_show_90p']"
            :projects="['radler/fmtlib/completion/std.string (hot)', 'radler/fmtlib/completion/std.shared_ptr (dep) (hot)', 'radler/fmtlib/completion/fmt.join_view (dep) (hot)']"
          />
        </section>
      </div>
    </section>

    <Divider title="Inline Completion" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Inline: Time to show`"
            :measure="['callInlineCompletionOnCompletion']"
            :projects="['radler/fmtlib/completion/if(spec)']"
          />
        </section>
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
