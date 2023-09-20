<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="clion"
    persistent-id="clion_performance_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section class="flex gap-x-6 flex-col md:flex-row">
      <!--<div class="flex-1 min-w-0">-->
      <!--  <AggregationChart-->
      <!--    :configurators="averagesConfigurators"-->
      <!--    :aggregated-project="'clion/%'"-->
      <!--    :aggregated-measure="'processingSpeed#%'"-->
      <!--    :is-like="true"-->
      <!--    :title="'[CLion] Indexing speed (kB/s)'"-->
      <!--    :chart-color="'#219653'"-->
      <!--    :value-unit="'counter'"-->
      <!--  />-->
      <!--</div>-->
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'clion/%hot%'"
          :aggregated-measure="'fus_time_to_show_90p'"
          :is-like="true"
          :title="'[CLion] Time to show completion list'"
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
      <!--<div class="flex-1 min-w-0">-->
      <!--  <AggregationChart-->
      <!--    :configurators="averagesConfigurators"-->
      <!--    :aggregated-project="'radler/%'"-->
      <!--    :aggregated-measure="'processingSpeed#%'"-->
      <!--    :is-like="true"-->
      <!--    :title="'[Radler] Indexing speed (kB/s)'"-->
      <!--    :chart-color="'#219653'"-->
      <!--    :value-unit="'counter'"-->
      <!--  />-->
      <!--</div>-->
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-project="'radler/%hot%'"
          :aggregated-measure="'fus_time_to_show_90p'"
          :is-like="true"
          :title="'[Radler] Time to show completion list'"
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

    <Divider title="General" />

    <section>
      <GroupProjectsChart
        label="[CLion vs Radler] Inspect project (not only C/C++) (fmtlib)"
        measure="globalInspections"
        :projects="['clion/fmtlib/inspection', 'radler/fmtlib/inspection']"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="[CLion vs Radler] Time to show test gutter (luau, Linter.test.cpp)"
        measure="%waitFirstTestGutter"
        :projects="['clion/luau/checkLocalTestConfig/Linter.test.cpp.marks', 'radler/luau/checkLocalTestConfig/Linter.test.cpp.marks']"
      />
    </section>

    <Divider title="Completion" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion vs Radler] Time to show completion list, 90th percentile (std::string, hot)"
          measure="fus_time_to_show_90p"
          :projects="['clion/fmtlib/completion/std.string (hot)', 'radler/fmtlib/completion/std.string (hot)']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler, clangd vs R#] First element calculated, 90th percentile (std::string, hot)"
          :measure="['fus_clangd_time_ms_90p', 'fus_rider_time_ms_90p']"
          :projects="['radler/fmtlib/completion/std.string (hot)']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler, clangd vs R#] Total items, mean (std::string, hot)"
          :measure="['fus_clangd_items_count_mean', 'fus_rider_items_count_mean']"
          :projects="['radler/fmtlib/completion/std.string (hot)']"
        />
      </div>
    </section>

    <Divider title="Actions" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion vs Radler] Find Usages (macro)"
          measure="%syncAction FindUsages"
          :projects="['clion/luau/findUsages/macro (LUAU_ASSERT)', 'radler/luau/findUsages/macro (LUAU_ASSERT)']"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Go to Declaration (ctor)"
          measure="clionGotoDeclaration"
          :projects="['clion/luau/gotoDeclaration/AstStatDeclareFunction.ctor', 'radler/luau/gotoDeclaration/AstStatDeclareFunction.ctor']"
        />
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
