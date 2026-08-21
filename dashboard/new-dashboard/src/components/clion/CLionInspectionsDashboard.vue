<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_inspections_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="First Code Analysis" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`First Code Analysis (fmt)`"
            :measure="['firstCodeAnalysis']"
            :projects="['radler/fmtlib/typing/simple (4 lines)', 'radler/fmtlib/typing/simple (4 lines), Grazie disabled']"
          />
        </section>
      </div>
    </section>

    <Divider title="Inspection" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Local Inspection (luau)`"
            :measure="['localInspections']"
            :projects="['radler/luau/localInspection/Lexer']"
          />
        </section>
      </div>
    </section>
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <section>
          <GroupProjectsChart
            :label="`Global Inspection (not only C/C++) (fmt)`"
            :measure="['globalInspections']"
            :projects="['radler/fmtlib/globalInspection', 'radler/fmtlib/inspection']"
          />
        </section>
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
