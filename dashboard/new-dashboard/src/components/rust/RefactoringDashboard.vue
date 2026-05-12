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
        label="Rename Function with low usage"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameLowUseFn"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename Function with high usage"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameHighUseFn"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename Struct with low usage"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameLowUseStruct"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename Struct with high usage"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameHighUseStruct"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Promote Small File to Directory"
        measure="execute_editor_rust.rspromotemoduletodirectoryaction"
        :projects="rustRefactoringTestCases.moveSmallFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Promote Big File to Directory"
        measure="execute_editor_rust.rspromotemoduletodirectoryaction"
        :projects="rustRefactoringTestCases.moveBigFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Call Optimize Imports in Small File"
        measure="execute_editor_optimizeimports"
        :projects="rustRefactoringTestCases.optimizeImportsSmallFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Call Optimize Imports in Big File"
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
