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
        label="Rename src/foo.rs/fn add (low usage)"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameLowUseFn"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename src/small/utils.rs/fn format_event (high usage)"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameHighUseFn"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename src/small/model.rs/struct Time (low usage)"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameLowUseStruct"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rename src/big/big_file.rs/struct BigRegistry (high usage)"
        measure="performInlineRename"
        :projects="rustRefactoringTestCases.renameHighUseStruct"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Promote src/small/utils.rs to Directory"
        measure="execute_editor_rust.rspromotemoduletodirectoryaction"
        :projects="rustRefactoringTestCases.moveSmallFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Promote src/big/big_file.rs to Directory"
        measure="execute_editor_rust.rspromotemoduletodirectoryaction"
        :projects="rustRefactoringTestCases.moveBigFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Optimize Imports: src/small/model.rs"
        measure="execute_editor_optimizeimports"
        :projects="rustRefactoringTestCases.optimizeImportsSmallFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Optimize Imports: src/big/big_file.rs"
        measure="execute_editor_optimizeimports"
        :projects="rustRefactoringTestCases.optimizeImportsBigFile"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Optimize Imports in Directory: src/small"
        measure="optimizeImportsOnDirectory"
        :projects="rustRefactoringTestCases.optimizeImportsSmallDirectory"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Optimize Imports in Directory: src/big"
        measure="optimizeImportsOnDirectory"
        :projects="rustRefactoringTestCases.optimizeImportsBigDirectory"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { rustRefactoringTestCases } from "./RustTestCases"
</script>
