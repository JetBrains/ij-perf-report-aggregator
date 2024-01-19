<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="fullLineDashboardDev"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Full line completion fimModel and rag"
        :measure="['inlineCompletionShow#mean_value']"
        :projects="[
          'kotlin_language_server/full-line-completion/DefinitionTest/empty-fun-full-line',
          'kotlin_language_server_full_line/full-line-completion/DefinitionTest/empty-fun-full-line flcc_fimModel=false flcc_ragEnabled=false',
          'kotlin_language_server_full_line/full-line-completion/DefinitionTest/empty-fun-full-line flcc_fimModel=true flcc_ragEnabled=false',
          'kotlin_language_server_full_line/full-line-completion/DefinitionTest/empty-fun-full-line flcc_fimModel=false flcc_ragEnabled=true',
          'kotlin_language_server_full_line/full-line-completion/DefinitionTest/empty-fun-full-line flcc_fimModel=true flcc_ragEnabled=true',
        ]"
      />
      <GroupProjectsChart
        label="Full line completion maxModelCtx"
        :measure="['inlineCompletionShow#mean_value']"
        :projects="[
          'kotlin_language_server/full-line-completion/DefinitionTest/empty-fun-full-line',
          'kotlin_language_server_full_line/full-line-completion/DefinitionTest/empty-fun-full-line flcc_maxModelCtx=1500',
          'kotlin_language_server_full_line/full-line-completion/DefinitionTest/empty-fun-full-line flcc_maxModelCtx=384',
        ]"
      />
    </section>
  </DashboardPage>
</template>
<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
