<template>
  <DashboardPage
    db-name="perfint"
    table="rust"
    persistent-id="rust_plugin_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    release-configurator="EAP / Release"
  >
    <section>
      <GroupProjectsChart
        label="Local Inspections (on file open, metric 'firstCodeAnalysis')"
        measure="firstCodeAnalysis"
        :projects="rustLocalInspectionCases"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Local Inspections (on typing top-level, metric 'typingCodeAnalyzing#mean_value')"
        measure="typingCodeAnalyzing#mean_value"
        :projects="rustLocalInspectionCases.map((testCase) => `${testCase}-top-level-typing`)"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Local Inspections (on typing stmt in function, metric 'typingCodeAnalyzing#mean_value')"
        measure="typingCodeAnalyzing#mean_value"
        :projects="rustLocalInspectionCases"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Global Inspection execution time (metric 'globalInspections')"
        measure="globalInspections"
        :projects="rustGlobalInspectionProjects.map((project) => `global-inspection/${project}-inspection`)"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion#mean_value"
        :projects="[
          'completion/arrow-rs',
          'completion/vec',
          'completion/arrow-rs/parse',
          'completion/arrow-rs/cast',
          'completion/arrow-rs/comparison',
          'completion/bat/src/bin/bat/clap_app.rs',
          'completion/bat/src/diff.rs',
          'completion/bevy/crates/bevy_pbr/src/light.rs',
          'completion/bevy/crates/bevy_pbr/src/render/light.rs',
          'completion/bevy/crates/bevy_render/macros/src/as_bind_group.rs',
          'completion/fd/src/cli.rs',
          'completion/fuel/crates/fuel-core/src/executor.rs',
          'completion/cargo/testsuite/build.rs',
          'completion/cargo/testsuite/build_script.rs',
          'completion/cargo/testsuite/metadata.rs',
          'completion/cargo/testsuite/test.rs',
          'completion/cargo/testsuite/git.rs',
          'completion/cargo/config/mod.rs',
          'completion/cargo/fingerprint/mod.rs',
          'completion/cargo/toml/mod.rs',
          'completion/chrono/src/naive/date.rs',
          'completion/chrono/src/duration.rs',
          'completion/clap/toml/mod.rs',
          'completion/clap/derives/subcommand.rs',
          'completion/clap/parser/parser.rs',
          'completion/deno/integration/run_tests.rs',
          'completion/deno/args/flags.rs',
          'completion/diesel/src/table.rs',
          'completion/diesel/connection/mod.rs',
          'completion/diesel/macros/mod.rs',
          'completion/diesel/type_impls/tuples.rs',
          'completion/diesel/src/sql_function.rs',
          'completion/hyperfine/src/command.rs',
          'completion/hyperfine/src/export/markup.rs',
          'completion/lemmy/comment_report_view.rs',
          'completion/nalgebra/src/base/matrix.rs',
          'completion/nalgebra/src/base/alias_slice.rs',
          'completion/nalgebra/src/linalg/cholesky.rs',
          'completion/mySql/conn/mod.rs',
          'completion/mySql/opts/mod.rs',
          'completion/mySql/pool/mod.rs',
          'completion/mySql/io/mod.rs',
          'completion/mySql/src/query.rs',
          'completion/mySql/routines/helpers.rs',
          'completion/rustAnalyzer/crates/ide-db/src/generated/lints.rs',
          'completion/rustAnalyzer/crates/hir-ty/src/chalk_db.rs',
          'completion/solana/progress_map.rs',
          'completion/tokio/tokio/src/net/udp.rs',
          'completion/tokio/tokio/tests/rt_common.rs',
          'completion/tokio/tokio/src/time/interval.rs',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages"
        measure="findUsages"
        :projects="['find-usages/yew', 'find-usages/wasm']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { rustLocalInspectionCases, rustGlobalInspectionProjects } from "./RustTestCases"
</script>
