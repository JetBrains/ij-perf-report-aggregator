<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_performance_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="General" />

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
          label="[Radler] Resolve All References (cmake)"
          :measure="['clangd_light_modules_total_time_s', 'clangd_no_modules_total_time_s', 'nova_resolving_references_s']"
          :projects="['radler/cmakeResolving/measureResolve/cmake']"
          :value-unit="'counter'"
        />
      </div>
    </section>

    <Divider title="Actions" />

    <section class="flex gap-x-6 flex-col md:flex-row">
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
          label="Go to Declaration (ctor)"
          measure="clionGotoDeclaration"
          project="luau/gotoDeclaration/AstStatDeclareFunction.ctor"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import CLionVsRadlerGroupProjectsChart from "./CLionVsRadlerGroupProjectsChart.vue"
</script>
