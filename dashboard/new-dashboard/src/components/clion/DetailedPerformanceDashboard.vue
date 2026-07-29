<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_detailed_performance_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="Find Usages" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (enumerable)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          project="luau/findUsages/enumerable (LuauOpcode)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (enumerator)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          project="luau/findUsages/enumerator (LOP_NOP)"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (class template)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          project="luau/findUsages/class template (DenseHashTable)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (macro)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          project="luau/findUsages/macro (LUAU_ASSERT)"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (cmake, class)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          project="cmake/findUsages/class (cmCTestResourceAllocator)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (cmake, macro)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          project="cmake/findUsages/macro (SAFEDIV)"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Find Usages (cmake, member)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
          project="cmake/findUsages/member (SlotsNeeded)"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="[CLion] Find Usages (all tests)"
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
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
          :measure="['%syncAction FindUsages', 'findUsagesInToolWindow']"
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

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to Show Context menu in Project View for a single file"
          measure="fus_popup_latency_ms"
          project="llvm/showPopupMenu/SingleSelection"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to Show Context menu in Project View for multiple files"
          measure="fus_popup_latency_ms"
          project="llvm/showPopupMenu/MultipleSelection"
        />
      </div>
      <div class="flex-1 min-w-0">
        <CLionVsRadlerGroupProjectsChart
          label="Time to Show Context menu in Project View for multiple files during indexing"
          measure="fus_popup_latency_ms"
          project="llvm/showPopupMenu/MultipleSelectionDuringIndexing"
        />
      </div>
    </section>

    <!-- END -->
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import CLionVsRadlerGroupProjectsChart from "./CLionVsRadlerGroupProjectsChart.vue"
</script>
