<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="clion"
    persistent-id="clion_detailed_performance_dashboard"
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

    <Divider title="Completion" />

    <section class="flex gap-x-6 flex-col md:flex-row">
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
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Completion, mean value"
          measure="completion#mean_value"
          :projects="[
            'clion/fmtlib/completion/std.string (cold)',
            'clion/fmtlib/completion/std.string (hot)',
            'clion/fmtlib/completion/std.shared_ptr (dep) (hot)',
            'clion/fmtlib/completion/fmt.join_view (dep) (hot)',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler] Completion, mean value"
          measure="completion#mean_value"
          :projects="[
            'radler/fmtlib/completion/std.string (cold)',
            'radler/fmtlib/completion/std.string (hot)',
            'radler/fmtlib/completion/std.shared_ptr (dep) (hot)',
            'radler/fmtlib/completion/fmt.join_view (dep) (hot)',
          ]"
        />
      </div>
    </section>

    <Divider title="Find Usages" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Find Usages (enumerable)"
          measure="%syncAction FindUsages"
          :projects="['clion/luau/findUsages/enumerable (LuauOpcode)', 'radler/luau/findUsages/enumerable (LuauOpcode)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Find Usages (enumerator)"
          measure="%syncAction FindUsages"
          :projects="['clion/luau/findUsages/enumerator (LOP_NOP)', 'radler/luau/findUsages/enumerator (LOP_NOP)']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Find Usages (class template)"
          measure="%syncAction FindUsages"
          :projects="['clion/luau/findUsages/class template (DenseHashTable)', 'radler/luau/findUsages/class template (DenseHashTable)']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Find Usages (macro)"
          measure="%syncAction FindUsages"
          :projects="['clion/luau/findUsages/macro (LUAU_ASSERT)', 'radler/luau/findUsages/macro (LUAU_ASSERT)']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Find Usages"
          measure="%syncAction FindUsages"
          :projects="[
            'clion/luau/findUsages/enumerable (LuauOpcode)',
            'clion/luau/findUsages/enumerator (LOP_NOP)',
            'clion/luau/findUsages/class template (DenseHashTable)',
            'clion/luau/findUsages/macro (LUAU_ASSERT)',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler] Find Usages"
          measure="%syncAction FindUsages"
          :projects="[
            'radler/luau/findUsages/enumerable (LuauOpcode)',
            'radler/luau/findUsages/enumerator (LOP_NOP)',
            'radler/luau/findUsages/class template (DenseHashTable)',
            'radler/luau/findUsages/macro (LUAU_ASSERT)',
          ]"
        />
      </div>
    </section>

    <Divider title="Go to Declaration" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Go to Declaration (ctor)"
          measure="%syncAction GotoDeclaration"
          :projects="['clion/luau/gotoDeclaration/AstStatDeclareFunction.ctor', 'radler/luau/gotoDeclaration/AstStatDeclareFunction.ctor']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Go to Declaration (method)"
          measure="%syncAction GotoDeclaration"
          :projects="['clion/luau/gotoDeclaration/TypeChecker.getScopes', 'radler/luau/gotoDeclaration/TypeChecker.getScopes']"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Go to Declaration (std::string - alias)"
          measure="%syncAction GotoDeclaration"
          :projects="['clion/luau/gotoDeclaration/std.string', 'radler/luau/gotoDeclaration/std.string']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Go to Declaration (macro)"
          measure="%syncAction GotoDeclaration"
          :projects="['clion/luau/gotoDeclaration/LUAU_ASSERT', 'radler/luau/gotoDeclaration/LUAU_ASSERT']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Go to Declaration (time.h - header)"
          measure="%syncAction GotoDeclaration"
          :projects="['clion/luau/gotoDeclaration/time.h', 'radler/luau/gotoDeclaration/time.h']"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Go to Declaration"
          measure="%syncAction GotoDeclaration"
          :projects="[
            'clion/luau/gotoDeclaration/AstStatDeclareFunction.ctor',
            'clion/luau/gotoDeclaration/TypeChecker.getScopes',
            'clion/luau/gotoDeclaration/std.string',
            'clion/luau/gotoDeclaration/LUAU_ASSERT',
            'clion/luau/gotoDeclaration/time.h',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler] Go to Declaration"
          measure="%syncAction GotoDeclaration"
          :projects="[
            'radler/luau/gotoDeclaration/AstStatDeclareFunction.ctor',
            'radler/luau/gotoDeclaration/TypeChecker.getScopes',
            'radler/luau/gotoDeclaration/std.string',
            'radler/luau/gotoDeclaration/LUAU_ASSERT',
            'radler/luau/gotoDeclaration/time.h',
          ]"
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
