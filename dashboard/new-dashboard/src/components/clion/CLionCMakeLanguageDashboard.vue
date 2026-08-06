<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_find_usage_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="Find Usages" />
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Total Execution Time"
          :measure="['findUsagesInToolWindow']"
          :projects="['radler/llvm/findUsages/function (add_llvm_component_library)', 'radler/llvm/findUsages/function (set)']"
        />
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Number of Usages"
          measure="findUsagesInToolWindow#number"
          :projects="['radler/llvm/findUsages/function (add_llvm_component_library)', 'radler/llvm/findUsages/function (set)']"
        />
      </div>
    </section>
    <Divider title="Go to Declaration" />
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Go to: Target Declaration"
          :measure="['clionGotoDeclaration']"
          :projects="['radler/luau/gotoDeclaration/Luau.Compile.CLI']"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
