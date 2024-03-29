<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="rust"
    persistent-id="rust_plugin_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section class="flex gap-6">
      <div class="w-1/2">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeedAvg#Rust'"
          :title="'Indexing Rust (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          :measure="['indexingTimeWithoutPauses']"
          :projects="['indexing/rustling', 'indexing/yew']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          :measure="['scanningTimeWithoutPauses']"
          :projects="['indexing/rustling', 'indexing/yew']"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Global Inspection execution time"
        measure="globalInspections"
        :projects="inspectionProjects.map((project) => `global-inspection/${project}-inspection`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Global Inspection (Rust-Only))"
        measure="globalInspections"
        :projects="inspectionProjects.map((project) => `rust-only-inspection/${project}-inspection`)"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion#mean_value"
        :projects="['completion/arrow-rs', 'completion/vec']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages"
        measure="findUsages"
        :projects="['find-usages/yew', 'find-usages/wasm']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local Inspections (on file open)"
        measure="firstCodeAnalysis"
        :projects="[
          'arrow-rs/local-inspection/parse',
          'arrow-rs/local-inspection/cast',
          'arrow-rs/local-inspection/comparison',
          'bat/local-inspection/src/bin/bat/clap_app.rs',
          'bat/local-inspection/src/diff.rs',
          'bevy/local-inspection/crates/bevy_pbr/src/light.rs',
          'bevy/local-inspection/crates/bevy_pbr/src/render/light.rs',
          'bevy/local-inspection/crates/bevy_render/macros/src/as_bind_group.rs',
          'fd/local-inspection/src/cli.rs',
          'fuel/local-inspection/crates/fuel-core/src/executor.rs',
          'cargo/local-inspection/testsuite/build.rs',
          'cargo/local-inspection/testsuite/build_script.rs',
          'cargo/local-inspection/testsuite/metadata.rs',
          'cargo/local-inspection/testsuite/test.rs',
          'cargo/local-inspection/testsuite/git.rs',
          'cargo/local-inspection/config/mod.rs',
          'cargo/local-inspection/fingerprint/mod.rs',
          'cargo/local-inspection/toml/mod.rs',
          'chrono/local-inspection/src/naive/date.rs',
          'chrono/local-inspection/src/duration.rs',
          'clap/local-inspection/toml/mod.rs',
          'clap/local-inspection/derives/subcommand.rs',
          'clap/local-inspection/parser/parser.rs',
          'deno/local-inspection/integration/run_tests.rs',
          'deno/local-inspection/args/flags.rs',
          'diesel/local-inspection/src/table.rs',
          'diesel/local-inspection/connection/mod.rs',
          'diesel/local-inspection/macros/mod.rs',
          'diesel/local-inspection/type_impls/tuples.rs',
          'diesel/local-inspection/src/sql_function.rs',
          'hyperfine/local-inspection/src/command.rs',
          'hyperfine/local-inspection/src/export/markup.rs',
          'lemmy/local-inspection/comment_report_view.rs',
          'nalgebra/local-inspection/src/base/matrix.rs',
          'nalgebra/local-inspection/src/base/alias_slice.rs',
          'nalgebra/local-inspection/src/linalg/cholesky.rs',
          'mySql/local-inspection/conn/mod.rs',
          'mySql/local-inspection/opts/mod.rs',
          'mySql/local-inspection/pool/mod.rs',
          'mySql/local-inspection/io/mod.rs',
          'mySql/local-inspection/src/query.rs',
          'mySql/local-inspection/routines/helpers.rs',
          'rustAnalyzer/local-inspection/crates/ide-db/src/generated/lints.rs',
          'rustAnalyzer/local-inspection/crates/hir-ty/src/chalk_db.rs',
          'solana/local-inspection/progress_map.rs',
          'tokio/local-inspection/tokio/src/net/udp.rs',
          'tokio/local-inspection/tokio/tests/rt_common.rs',
          'tokio/local-inspection/tokio/src/time/interval.rs',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local Inspections (on typing top-level)"
        measure="typingCodeAnalyzing#mean_value"
        :projects="[
          'arrow-rs/local-inspection/parse-top-level-typing',
          'arrow-rs/local-inspection/cast-top-level-typing',
          'arrow-rs/local-inspection/comparison-top-level-typing',
          'bat/local-inspection/src/bin/bat/clap_app.rs-top-level-typing',
          'bat/local-inspection/src/diff.rs-top-level-typing',
          'bevy/local-inspection/crates/bevy_pbr/src/light.rs-top-level-typing',
          'bevy/local-inspection/crates/bevy_pbr/src/render/light.rs-top-level-typing',
          'bevy/local-inspection/crates/bevy_render/macros/src/as_bind_group.rs-top-level-typing',
          'fd/local-inspection/src/cli.rs-top-level-typing',
          'fuel/local-inspection/crates/fuel-core/src/executor.rs-top-level-typing',
          'cargo/local-inspection/testsuite/build.rs-top-level-typing',
          'cargo/local-inspection/testsuite/build_script.rs-top-level-typing',
          'cargo/local-inspection/testsuite/metadata.rs-top-level-typing',
          'cargo/local-inspection/testsuite/test.rs-top-level-typing',
          'cargo/local-inspection/testsuite/git.rs-top-level-typing',
          'cargo/local-inspection/config/mod.rs-top-level-typing',
          'cargo/local-inspection/fingerprint/mod.rs-top-level-typing',
          'cargo/local-inspection/toml/mod.rs-top-level-typing',
          'chrono/local-inspection/src/naive/date.rs-top-level-typing',
          'chrono/local-inspection/src/duration.rs-top-level-typing',
          'clap/local-inspection/toml/mod.rs-top-level-typing',
          'clap/local-inspection/derives/subcommand.rs-top-level-typing',
          'clap/local-inspection/parser/parser.rs-top-level-typing',
          'deno/local-inspection/integration/run_tests.rs-top-level-typing',
          'deno/local-inspection/args/flags.rs-top-level-typing',
          'diesel/local-inspection/src/table.rs-top-level-typing',
          'diesel/local-inspection/connection/mod.rs-top-level-typing',
          'diesel/local-inspection/macros/mod.rs-top-level-typing',
          'diesel/local-inspection/type_impls/tuples.rs-top-level-typing',
          'diesel/local-inspection/src/sql_function.rs-top-level-typing',
          'hyperfine/local-inspection/src/command.rs-top-level-typing',
          'hyperfine/local-inspection/src/export/markup.rs-top-level-typing',
          'lemmy/local-inspection/comment_report_view.rs-top-level-typing',
          'nalgebra/local-inspection/src/base/matrix.rs-top-level-typing',
          'nalgebra/local-inspection/src/base/alias_slice.rs-top-level-typing',
          'nalgebra/local-inspection/src/linalg/cholesky.rs-top-level-typing',
          'mySql/local-inspection/conn/mod.rs-top-level-typing',
          'mySql/local-inspection/opts/mod.rs-top-level-typing',
          'mySql/local-inspection/pool/mod.rs-top-level-typing',
          'mySql/local-inspection/io/mod.rs-top-level-typing',
          'mySql/local-inspection/src/query.rs-top-level-typing',
          'mySql/local-inspection/routines/helpers.rs-top-level-typing',
          'rustAnalyzer/local-inspection/crates/ide-db/src/generated/lints.rs-top-level-typing',
          'rustAnalyzer/local-inspection/crates/hir-ty/src/chalk_db.rs-top-level-typing',
          'solana/local-inspection/progress_map.rs-top-level-typing',
          'tokio/local-inspection/tokio/src/net/udp.rs-top-level-typing',
          'tokio/local-inspection/tokio/tests/rt_common.rs-top-level-typing',
          'tokio/local-inspection/tokio/src/time/interval.rs-top-level-typing',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local Inspections (on typing stmt)"
        measure="typingCodeAnalyzing#mean_value"
        :projects="[
          'arrow-rs/local-inspection/parse',
          'arrow-rs/local-inspection/cast',
          'arrow-rs/local-inspection/comparison',
          'bat/local-inspection/src/bin/bat/clap_app.rs',
          'bat/local-inspection/src/diff.rs',
          'bevy/local-inspection/crates/bevy_pbr/src/light.rs',
          'bevy/local-inspection/crates/bevy_pbr/src/render/light.rs',
          'bevy/local-inspection/crates/bevy_render/macros/src/as_bind_group.rs',
          'fd/local-inspection/src/cli.rs',
          'fuel/local-inspection/crates/fuel-core/src/executor.rs',
          'cargo/local-inspection/testsuite/build.rs',
          'cargo/local-inspection/testsuite/build_script.rs',
          'cargo/local-inspection/testsuite/metadata.rs',
          'cargo/local-inspection/testsuite/test.rs',
          'cargo/local-inspection/testsuite/git.rs',
          'cargo/local-inspection/config/mod.rs',
          'cargo/local-inspection/fingerprint/mod.rs',
          'cargo/local-inspection/toml/mod.rs',
          'chrono/local-inspection/src/naive/date.rs',
          'chrono/local-inspection/src/duration.rs',
          'clap/local-inspection/toml/mod.rs',
          'clap/local-inspection/derives/subcommand.rs',
          'clap/local-inspection/parser/parser.rs',
          'deno/local-inspection/integration/run_tests.rs',
          'deno/local-inspection/args/flags.rs',
          'diesel/local-inspection/src/table.rs',
          'diesel/local-inspection/connection/mod.rs',
          'diesel/local-inspection/macros/mod.rs',
          'diesel/local-inspection/type_impls/tuples.rs',
          'diesel/local-inspection/src/sql_function.rs',
          'hyperfine/local-inspection/src/command.rs',
          'hyperfine/local-inspection/src/export/markup.rs',
          'lemmy/local-inspection/comment_report_view.rs',
          'nalgebra/local-inspection/src/base/matrix.rs',
          'nalgebra/local-inspection/src/base/alias_slice.rs',
          'nalgebra/local-inspection/src/linalg/cholesky.rs',
          'mySql/local-inspection/conn/mod.rs',
          'mySql/local-inspection/opts/mod.rs',
          'mySql/local-inspection/pool/mod.rs',
          'mySql/local-inspection/io/mod.rs',
          'mySql/local-inspection/src/query.rs',
          'mySql/local-inspection/routines/helpers.rs',
          'rustAnalyzer/local-inspection/crates/ide-db/src/generated/lints.rs',
          'rustAnalyzer/local-inspection/crates/hir-ty/src/chalk_db.rs',
          'solana/local-inspection/progress_map.rs',
          'tokio/local-inspection/tokio/src/net/udp.rs',
          'tokio/local-inspection/tokio/tests/rt_common.rs',
          'tokio/local-inspection/tokio/src/time/interval.rs',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Cargo Sync"
        measure="cargo_sync_execution_time"
        :projects="[
          'arrow-rs/local-inspection/cast',
          'arrow-rs/local-inspection/comparison',
          'arrow-rs/local-inspection/parse',
          'cargo/local-inspection/config/build_script.rs',
          'cargo/local-inspection/fingerprint/mod.rs',
          'cargo/local-inspection/testsuite/build.rs',
          'cargo/local-inspection/testsuite/build_script.rs',
          'cargo/local-inspection/testsuite/git.rs',
          'cargo/local-inspection/testsuite/metadata.rs',
          'cargo/local-inspection/testsuite/test.rs',
          'cargo/local-inspection/toml/mod.rs',
          'clap/local-inspection/derives/subcommand.rs',
          'clap/local-inspection/parser/parser.rs',
          'clap/local-inspection/toml/mod.rs',
          'completion/arrow-rs',
          'completion/vec',
          'deno/local-inspection/integration/run_tests.rs',
          'diesel/local-inspection/connection/mod.rs',
          'diesel/local-inspection/macros/mod.rs',
          'diesel/local-inspection/src/sql_function.rs',
          'diesel/local-inspection/src/table.rs',
          'diesel/local-inspection/type_impls/tuples.rs',
          'find-usages/wasm',
          'find-usages/yew',
          'global-inspection/arrowRs-inspection',
          'global-inspection/arti-inspection',
          'global-inspection/bat-inspection',
          'global-inspection/bevy-inspection',
          'global-inspection/cargo-inspection',
          'global-inspection/chrono-inspection',
          'global-inspection/clap-inspection',
          'global-inspection/deno-inspection',
          'global-inspection/diem-inspection',
          'global-inspection/fd-inspection',
          'global-inspection/fuel-inspection',
          'global-inspection/hyperfine-inspection',
          'global-inspection/meiliSearch-inspection',
          'global-inspection/nalgebra-inspection',
          'global-inspection/rand-inspection',
          'global-inspection/ripgrep-inspection',
          'global-inspection/ruffle-inspection',
          'global-inspection/solana-inspection',
          'global-inspection/spotify-inspection',
          'global-inspection/tokio-inspection',
          'global-inspection/turbo-inspection',
          'global-inspection/yew-inspection',
          'indexing/rustling',
          'indexing/yew',
          'mySql/local-inspection/conn/mod.rs',
          'mySql/local-inspection/io/mod.rs',
          'mySql/local-inspection/opts/mod.rs',
          'mySql/local-inspection/pool/mod.rs',
          'mySql/local-inspection/src/query.rs',
          'rust-only-inspection/arrowRs-inspection',
          'rust-only-inspection/arti-inspection',
          'rust-only-inspection/bat-inspection',
          'rust-only-inspection/bevy-inspection',
          'rust-only-inspection/cargo-inspection',
          'rust-only-inspection/chrono-inspection',
          'rust-only-inspection/clap-inspection',
          'rust-only-inspection/deno-inspection',
          'rust-only-inspection/diem-inspection',
          'rust-only-inspection/fd-inspection',
          'rust-only-inspection/fuel-inspection',
          'rust-only-inspection/hyperfine-inspection',
          'rust-only-inspection/lemmy-inspection',
          'rust-only-inspection/meiliSearch-inspection',
          'rust-only-inspection/nalgebra-inspection',
          'rust-only-inspection/rand-inspection',
          'rust-only-inspection/ripgrep-inspection',
          'rust-only-inspection/ruffle-inspection',
          'rust-only-inspection/solana-inspection',
          'rust-only-inspection/spotify-inspection',
          'rust-only-inspection/tokio-inspection',
          'rust-only-inspection/turbo-inspection',
          'rust-only-inspection/veloren-inspection',
          'rust-only-inspection/yew-inspection',
          'solana/progress_map.rs',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust macro expansion time"
        measure="rust_macro_expansion_execution_time"
        :projects="[
          'arrow-rs/local-inspection/cast',
          'arrow-rs/local-inspection/comparison',
          'arrow-rs/local-inspection/parse',
          'cargo/local-inspection/config/build_script.rs',
          'cargo/local-inspection/fingerprint/mod.rs',
          'cargo/local-inspection/testsuite/build.rs',
          'cargo/local-inspection/testsuite/build_script.rs',
          'cargo/local-inspection/testsuite/git.rs',
          'cargo/local-inspection/testsuite/metadata.rs',
          'cargo/local-inspection/testsuite/test.rs',
          'cargo/local-inspection/toml/mod.rs',
          'clap/local-inspection/derives/subcommand.rs',
          'clap/local-inspection/parser/parser.rs',
          'clap/local-inspection/toml/mod.rs',
          'completion/arrow-rs',
          'completion/vec',
          'deno/local-inspection/integration/run_tests.rs',
          'diesel/local-inspection/connection/mod.rs',
          'diesel/local-inspection/macros/mod.rs',
          'diesel/local-inspection/src/sql_function.rs',
          'diesel/local-inspection/src/table.rs',
          'diesel/local-inspection/type_impls/tuples.rs',
          'find-usages/wasm',
          'find-usages/yew',
          'indexing/rustling',
          'indexing/yew',
          'mySql/local-inspection/conn/mod.rs',
          'mySql/local-inspection/io/mod.rs',
          'mySql/local-inspection/opts/mod.rs',
          'mySql/local-inspection/pool/mod.rs',
          'mySql/local-inspection/src/query.rs',
          'solana/progress_map.rs',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Rust DefMaps execution time"
        measure="rust_def_maps_execution_time"
        :projects="[
          'arrow-rs/local-inspection/cast',
          'arrow-rs/local-inspection/comparison',
          'arrow-rs/local-inspection/parse',
          'cargo/local-inspection/config/build_script.rs',
          'cargo/local-inspection/fingerprint/mod.rs',
          'cargo/local-inspection/testsuite/build.rs',
          'cargo/local-inspection/testsuite/build_script.rs',
          'cargo/local-inspection/testsuite/git.rs',
          'cargo/local-inspection/testsuite/metadata.rs',
          'cargo/local-inspection/testsuite/test.rs',
          'cargo/local-inspection/toml/mod.rs',
          'clap/local-inspection/derives/subcommand.rs',
          'clap/local-inspection/parser/parser.rs',
          'clap/local-inspection/toml/mod.rs',
          'completion/arrow-rs',
          'completion/vec',
          'deno/local-inspection/integration/run_tests.rs',
          'diesel/local-inspection/connection/mod.rs',
          'diesel/local-inspection/macros/mod.rs',
          'diesel/local-inspection/src/sql_function.rs',
          'diesel/local-inspection/src/table.rs',
          'diesel/local-inspection/type_impls/tuples.rs',
          'find-usages/wasm',
          'find-usages/yew',
          'global-inspection/arrowRs-inspection',
          'global-inspection/arti-inspection',
          'global-inspection/bat-inspection',
          'global-inspection/bevy-inspection',
          'global-inspection/cargo-inspection',
          'global-inspection/chrono-inspection',
          'global-inspection/clap-inspection',
          'global-inspection/deno-inspection',
          'global-inspection/diem-inspection',
          'global-inspection/fd-inspection',
          'global-inspection/fuel-inspection',
          'global-inspection/hyperfine-inspection',
          'global-inspection/meiliSearch-inspection',
          'global-inspection/nalgebra-inspection',
          'global-inspection/rand-inspection',
          'global-inspection/ripgrep-inspection',
          'global-inspection/ruffle-inspection',
          'global-inspection/solana-inspection',
          'global-inspection/spotify-inspection',
          'global-inspection/tokio-inspection',
          'global-inspection/turbo-inspection',
          'global-inspection/yew-inspection',
          'indexing/rustling',
          'indexing/yew',
          'mySql/local-inspection/conn/mod.rs',
          'mySql/local-inspection/io/mod.rs',
          'mySql/local-inspection/opts/mod.rs',
          'mySql/local-inspection/pool/mod.rs',
          'mySql/local-inspection/src/query.rs',
          'rust-only-inspection/arrowRs-inspection',
          'rust-only-inspection/arti-inspection',
          'rust-only-inspection/bat-inspection',
          'rust-only-inspection/bevy-inspection',
          'rust-only-inspection/cargo-inspection',
          'rust-only-inspection/chrono-inspection',
          'rust-only-inspection/clap-inspection',
          'rust-only-inspection/deno-inspection',
          'rust-only-inspection/diem-inspection',
          'rust-only-inspection/fd-inspection',
          'rust-only-inspection/fuel-inspection',
          'rust-only-inspection/hyperfine-inspection',
          'rust-only-inspection/lemmy-inspection',
          'rust-only-inspection/meiliSearch-inspection',
          'rust-only-inspection/nalgebra-inspection',
          'rust-only-inspection/rand-inspection',
          'rust-only-inspection/ripgrep-inspection',
          'rust-only-inspection/ruffle-inspection',
          'rust-only-inspection/solana-inspection',
          'rust-only-inspection/spotify-inspection',
          'rust-only-inspection/tokio-inspection',
          'rust-only-inspection/turbo-inspection',
          'rust-only-inspection/veloren-inspection',
          'rust-only-inspection/yew-inspection',
          'solana/progress_map.rs',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const inspectionProjects = [
  "arrowRs",
  "arti",
  "bat",
  "bevy",
  "cargo",
  "chrono",
  "clap",
  "deno",
  "diem",
  "fd",
  "fuel",
  "hyperfine",
  "meiliSearch",
  "nalgebra",
  "rand",
  "ripgrep",
  "ruffle",
  "solana",
  "spotify",
  "tokio",
  "turbo",
  "yew",
]
</script>
