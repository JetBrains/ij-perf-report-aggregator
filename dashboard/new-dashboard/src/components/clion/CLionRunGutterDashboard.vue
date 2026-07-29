<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_run_gutter_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
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

    <Divider title="Context menu in Project View" />

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
