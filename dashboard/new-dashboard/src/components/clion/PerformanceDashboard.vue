<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="clion"
    persistent-id="clion_performance_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section class="flex gap-x-6 flex-col md:flex-row">
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
          :aggregated-project="'clion/%/typing/%'"
          :aggregated-measure="'%#average_awt_delay'"
          :is-like="true"
          :title="'[CLion] UI responsiveness during typing'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
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
          :aggregated-project="'radler/%/typing/%'"
          :aggregated-measure="'%#average_awt_delay'"
          :is-like="true"
          :title="'[Radler] UI responsiveness during typing'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>

    <section>
      <GroupProjectsChart
        label="Global Inspections (fmtlib)"
        measure="globalInspections"
        :projects="['clion/fmtlib/inspection', 'radler/fmtlib/inspection']"
      />
    </section>

    <Divider title="Completion" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion, mean value (std::string, hot)"
          measure="completion#mean_value"
          :projects="['clion/fmtlib/completion/std.string (hot)', 'radler/fmtlib/completion/std.string (hot)']"
        />
      </div>
    </section>

    <Divider title="Actions" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Find Usages (macro)"
          measure="%syncAction FindUsages"
          :projects="['clion/luau/findUsages/macro (LUAU_ASSERT)', 'radler/luau/findUsages/macro (LUAU_ASSERT)']"
        />
      </div>
    </section>

    <!-- TODO: get rid of %action -->
    <!--<section class="flex gap-x-6 flex-col md:flex-row">-->
    <!--  <div class="flex-1 min-w-0">-->
    <!--    <GroupProjectsChart-->
    <!--      label="Go to Declaration (ctor)"-->
    <!--      measure="%action GotoDeclaration"-->
    <!--      :projects="['clion/luau/gotoDeclaration/AstStatDeclareFunction.ctor', 'radler/luau/gotoDeclaration/AstStatDeclareFunction.ctor']"-->
    <!--    />-->
    <!--  </div>-->
    <!--</section>-->
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
