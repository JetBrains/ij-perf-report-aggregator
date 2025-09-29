<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_performance_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
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
          :aggregated-measure="'typing#latency#max'"
          :is-like="true"
          :title="'[CLion] Typing latency(max)'"
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
          :aggregated-measure="'typing#latency#max'"
          :is-like="true"
          :title="'[Radler] Typing latency(max)'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>

    <Divider title="General" />

    <section>
      <CLionVsRadlerIndexingChart
        label="Index project (LLVM)"
        project="llvm/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerGroupProjectsChart
        label="Workspace Model Update (LLVM)"
        measure="cidr.workspace.metrics#duration_total_ms"
        project="llvm/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerGroupProjectsChart
        label="Inspect project (not only C/C++) (fmtlib)"
        measure="globalInspections"
        project="fmtlib/globalInspection"
      />
    </section>

    <section>
      <CLionVsRadlerGroupProjectsChart
        label="Time to show test gutter (luau, Linter.test.cpp)"
        measure="waitFirstTestGutter"
        project="luau/checkLocalTestConfig/Linter.test.cpp.marks"
      />
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler vs clangd] Resolve All References (cmake)"
          :measure="['clangd_light_modules_total_time_s', 'clangd_no_modules_total_time_s', 'nova_resolving_references_s']"
          :projects="['radler/cmakeResolving/measureResolve/cmake']"
          :value-unit="'counter'"
        />
      </div>
    </section>

    <Divider title="Completion" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show completion list, 90th percentile (std::string, hot)"
          measure="fus_time_to_show_90p"
          project="fmtlib/completion/std.string (hot)"
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
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (macro)"
          measure="%syncAction FindUsages"
          project="luau/findUsages/macro (LUAU_ASSERT)"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Go to Declaration (ctor)"
          measure="clionGotoDeclaration"
          project="luau/gotoDeclaration/AstStatDeclareFunction.ctor"
        />
      </div>
    </section>

    <Divider title="Debugger" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Stepping (fmtlib)"
          :measure="['debugStep_into', 'debugStep_out', 'debugStep_over']"
          project="fmtlib/debug/args-test/basic"
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
import CLionVsRadlerGroupProjectsChart from "./CLionVsRadlerGroupProjectsChart.vue"
import CLionVsRadlerIndexingChart from "./CLionVsRadlerIndexingChart.vue"
</script>
