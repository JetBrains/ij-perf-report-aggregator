<template>
  <DashboardPage
    db-name="perfintDev"
    table="rust"
    persistent-id="rust_refactoring_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Rename Low Use Function"
        measure="performInlineRename#mean_value"
        :projects="rustRefactoringTestCases.renameLowUseFn"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename High Use Function"
        measure="performInlineRename#mean_value"
        :projects="rustRefactoringTestCases.renameHighUseFn"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename Low Use Struct"
        measure="performInlineRename#mean_value"
        :projects="rustRefactoringTestCases.renameLowUseStruct"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename High Use Struct"
        measure="performInlineRename#mean_value"
        :projects="rustRefactoringTestCases.renameHighUseStruct"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Move Small File"
        measure="execute_editor_rust.rspromotemoduletodirectoryaction"
        :projects="rustRefactoringTestCases.moveSmallFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Move Big File"
        measure="execute_editor_rust.rspromotemoduletodirectoryaction"
        :projects="rustRefactoringTestCases.moveBigFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Optimize Imports Small File"
        measure="execute_editor_optimizeimports"
        :projects="rustRefactoringTestCases.optimizeImportsSmallFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Optimize Imports Big File"
        measure="execute_editor_optimizeimports"
        :projects="rustRefactoringTestCases.optimizeImportsBigFile"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { rustRefactoringTestCases } from "./RustTestCases"
</script>
