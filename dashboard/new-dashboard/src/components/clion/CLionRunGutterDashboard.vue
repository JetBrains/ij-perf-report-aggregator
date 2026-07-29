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
        <GroupProjectsChart
          label="Time to show test gutter (all tests)"
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
        <GroupProjectsChart
          label="Time to Show Context menu in Project View for a single file"
          measure="fus_popup_latency_ms"
          :projects="['radler/llvm/showPopupMenu/SingleSelection']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Time to Show Context menu in Project View for multiple files"
          measure="fus_popup_latency_ms"
          :projects="['radler/llvm/showPopupMenu/MultipleSelection']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Time to Show Context menu in Project View for multiple files during indexing"
          measure="fus_popup_latency_ms"
          :projects="['radler/llvm/showPopupMenu/MultipleSelectionDuringIndexing']"
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
</script>
