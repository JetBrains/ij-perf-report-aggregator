<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_detailed_performance_dashboard"
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
    <!--<section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Indexing (curl)"
          :measure="['indexingTimeWithoutPauses']"
          project="curl/indexing"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Scanning (curl)"
          :measure="['scanningTimeWithoutPauses']"
          project="curl/indexing"
        />
      </div>
    </section>-->

    <Divider title="Indexing" />

    <section>
      <CLionVsRadlerIndexingChart
        label="Index project (LLVM)"
        project="llvm/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerIndexingChart
        label="Index project (50k sources, 10k headers)"
        project="big_project_50k_10k/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerIndexingChart
        label="Index project (OpenCV)"
        project="opencv/indexing"
      />
    </section>

    <section>
      <CLionVsRadlerIndexingChart
        label="Index project (curl)"
        project="curl/indexing"
      />
    </section>

    <Divider title="Inspection" />

    <section>
      <CLionVsRadlerGroupProjectsChart
        label="Inspect project (not only C/C++) (fmtlib)"
        measure="globalInspections"
        project="fmtlib/globalInspection"
      />
    </section>

    <Divider title="Completion" />

    <!-- Completion: std::string (cold) -->
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show completion list, 90th percentile (std::string, cold)"
          measure="fus_time_to_show_90p"
          project="fmtlib/completion/std.string (cold)"
        />
      </div>
    </section>

    <!-- Completion: std::string (hot) -->
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show completion list, 90th percentile (std::string, hot)"
          measure="fus_time_to_show_90p"
          project="fmtlib/completion/std.string (hot)"
        />
      </div>
    </section>

    <!-- Completion: std::shared_ptr<T> (hot) -->
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show completion list, 90th percentile (std::shared_ptr<T>, hot)"
          measure="fus_time_to_show_90p"
          project="fmtlib/completion/std.shared_ptr (dep) (hot)"
        />
      </div>
    </section>

    <!-- Completion: fmt::join<It> -->
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show completion list, 90th percentile (fmt::join<It>, hot)"
          measure="fus_time_to_show_90p"
          project="fmtlib/completion/fmt.join_view (dep) (hot)"
        />
      </div>
    </section>

    <!-- Completion: overall -->
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Time to show completion list, 90th percentile (all tests)"
          measure="fus_time_to_show_90p"
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
          label="[Radler] Time to show completion list, 90th percentile (all tests)"
          measure="fus_time_to_show_90p"
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
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (enumerable)"
          measure="%syncAction FindUsages"
          project="luau/findUsages/enumerable (LuauOpcode)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (enumerator)"
          measure="%syncAction FindUsages"
          project="luau/findUsages/enumerator (LOP_NOP)"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (class template)"
          measure="%syncAction FindUsages"
          project="luau/findUsages/class template (DenseHashTable)"
        />
      </div>
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
          label="Find Usages (cmake, class)"
          measure="%syncAction FindUsages"
          project="cmake/findUsages/class (cmCTestResourceAllocator)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (cmake, macro)"
          measure="%syncAction FindUsages"
          project="cmake/findUsages/macro (SAFEDIV)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (cmake, member)"
          measure="%syncAction FindUsages"
          project="cmake/findUsages/member (SlotsNeeded)"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Find Usages (all tests)"
          measure="%syncAction FindUsages"
          :projects="[
            'clion/luau/findUsages/enumerable (LuauOpcode)',
            'clion/luau/findUsages/enumerator (LOP_NOP)',
            'clion/luau/findUsages/class template (DenseHashTable)',
            'clion/luau/findUsages/macro (LUAU_ASSERT)',
            'clion/cmake/findUsages/class (cmCTestResourceAllocator)',
            'clion/cmake/findUsages/macro (SAFEDIV)',
            'clion/cmake/findUsages/member (SlotsNeeded)',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler] Find Usages (all tests)"
          measure="%syncAction FindUsages"
          :projects="[
            'radler/luau/findUsages/enumerable (LuauOpcode)',
            'radler/luau/findUsages/enumerator (LOP_NOP)',
            'radler/luau/findUsages/class template (DenseHashTable)',
            'radler/luau/findUsages/macro (LUAU_ASSERT)',
            'radler/cmake/findUsages/class (cmCTestResourceAllocator)',
            'radler/cmake/findUsages/macro (SAFEDIV)',
            'radler/cmake/findUsages/member (SlotsNeeded)',
          ]"
        />
      </div>
    </section>

    <Divider title="Go to Declaration" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Go to Declaration (ctor)"
          measure="clionGotoDeclaration"
          project="luau/gotoDeclaration/AstStatDeclareFunction.ctor"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Go to Declaration (method)"
          measure="clionGotoDeclaration"
          project="luau/gotoDeclaration/TypeChecker.getScopes"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Go to Declaration (std::string - alias)"
          measure="clionGotoDeclaration"
          project="luau/gotoDeclaration/std.string"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Go to Declaration (macro)"
          measure="clionGotoDeclaration"
          project="luau/gotoDeclaration/LUAU_ASSERT"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Go to Declaration (time.h - header)"
          measure="clionGotoDeclaration"
          project="luau/gotoDeclaration/time.h"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Go to Declaration (all tests)"
          measure="clionGotoDeclaration"
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
          label="[Radler] Go to Declaration (all tests)"
          measure="clionGotoDeclaration"
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

    <Divider title="Test gutters" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show test gutter (luau, AstQuery.test.cpp)"
          measure="waitFirstTestGutter"
          project="luau/checkLocalTestConfig/AstQuery.test.cpp.marks"
        />
      </div>

      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show test gutter (luau, Linter.test.cpp)"
          measure="waitFirstTestGutter"
          project="luau/checkLocalTestConfig/Linter.test.cpp.marks"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show test gutter (luau, Repl.test.cpp)"
          measure="waitFirstTestGutter"
          project="luau/checkLocalTestConfig/Repl.test.cpp.marks"
        />
      </div>

      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show test gutter (luau, TypeInfer.unionTypes.test.cpp)"
          measure="waitFirstTestGutter"
          project="luau/checkLocalTestConfig/TypeInfer.unionTypes.test.cpp.marks"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show test gutter (openCV, test_houghlines.cpp)"
          measure="waitFirstTestGutter"
          project="opencv/checkLocalTestConfig/test.houghlines.cpp.marks"
        />
      </div>

      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to show test gutter (openCV, test_kalman.cpp)"
          measure="waitFirstTestGutter"
          project="opencv/checkLocalTestConfig/test.kalman.cpp.marks"
        />
      </div>
    </section>

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Time to show test gutter (all tests)"
          measure="waitFirstTestGutter"
          :projects="[
            'clion/luau/checkLocalTestConfig/AstQuery.test.cpp.marks',
            'clion/luau/checkLocalTestConfig/Linter.test.cpp.marks',
            'clion/luau/checkLocalTestConfig/Repl.test.cpp.marks',
            'clion/luau/checkLocalTestConfig/TypeInfer.unionTypes.test.cpp.marks',
            'clion/opencv/checkLocalTestConfig/test.houghlines.cpp.marks',
            'clion/opencv/checkLocalTestConfig/test.kalman.cpp.marks',
          ]"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[Radler] Time to show test gutter (all tests)"
          measure="waitFirstTestGutter"
          :projects="[
            'radler/luau/checkLocalTestConfig/AstQuery.test.cpp.marks',
            'radler/luau/checkLocalTestConfig/Linter.test.cpp.marks',
            'radler/luau/checkLocalTestConfig/Repl.test.cpp.marks',
            'radler/luau/checkLocalTestConfig/TypeInfer.unionTypes.test.cpp.marks',
            'radler/opencv/checkLocalTestConfig/test.houghlines.cpp.marks',
            'radler/opencv/checkLocalTestConfig/test.kalman.cpp.marks',
          ]"
        />
      </div>
    </section>

    <Divider title="Debugger" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Stepping (fmtlib)"
          :measure="['debugStep_into#mean_value', 'debugStep_out#mean_value', 'debugStep_out#mean_value']"
          project="fmtlib/debug/args-test/basic"
        />
      </div>
    </section>
    <!-- END -->
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
