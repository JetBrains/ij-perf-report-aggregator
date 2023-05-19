<template>
  <DashboardPage
    v-slot="{averagesConfigurators}"
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_refactoring_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'performInlineRename\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all rename K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'performInlineRename\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all rename K2'"
        />
      </div>

      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'prepareForRename\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all prepare-rename K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'performInlineRename\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all prepare-rename K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'startInlineRename\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all start-rename K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'startInlineRename\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all start-rename K2'"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'PerformInlineRename mean value on  K1"
          measure="performInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k1', 'kotlin_language_server/insertCode/Rename_renameSymbol_k1',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k1']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="PerformInlineRename mean value on  K2"
          measure="performInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k2', 'kotlin_language_server/insertCode/Rename_renameSymbol_k2',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k2']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'StartInlineRename mean value on  K1"
          measure="startInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k1', 'kotlin_language_server/insertCode/Rename_renameSymbol_k1',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k1']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="StartInlineRename mean value on  K2"
          measure="startInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k2', 'kotlin_language_server/insertCode/Rename_renameSymbol_k2',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k2']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'PrepareForRename mean value on  K1"
          measure="prepareForRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k1', 'kotlin_language_server/insertCode/Rename_renameSymbol_k1',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k1']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="PrepareForRename mean value on  K2"
          measure="prepareForRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k2', 'kotlin_language_server/insertCode/Rename_renameSymbol_k2',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k2']"
        />
      </div>
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"

</script>