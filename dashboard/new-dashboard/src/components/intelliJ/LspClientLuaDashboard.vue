<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_lsp_client_lua"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"

const declaration: ChartDefinition[] = [
  {
    labels: ["Code Lens on File Open"],
    measures: ["lspCodeLensOnFileOpen"],
    projects: ["lua_lsp_code_lens/code_lens/file_open"],
  },
  {
    labels: ["Completion", "Completion First Element Shown", "Completion Number"],
    measures: ["completion", "completion#firstElementShown", "completion#number"],
    projects: ["lua_lsp_completion/completion/imported_module", "lua_lsp_completion/completion/local_variables", "lua_lsp_completion/completion/table_fields"],
  },
  {
    labels: ["Diagnostics on First File Open", "Diagnostics on File Open"],
    measures: ["lspDiagnosticsOnFirstFileOpen", "lspDiagnosticsOnFileOpen"],
    projects: ["lua_lsp_diagnostics/diagnostics/file_open"],
  },
  {
    labels: ["Diagnostics on Typing", "Diagnostics on Deleting"],
    measures: ["lspDiagnosticsOnTyping", "lspDiagnosticsOnDeleting"],
    projects: ["lua_lsp_diagnostics/diagnostics/typing"],
  },
  {
    labels: ["Document Color on File Open"],
    measures: ["lspDocumentColorOnFileOpen"],
    projects: ["lua_lsp_document_color/document_color/file_open"],
  },
  {
    labels: ["Document Highlight on First Caret Move", "Document Highlight on Caret Move"],
    measures: ["lspDocumentHighlightOnFirstCaretMove", "lspDocumentHighlightOnCaretMove"],
    projects: ["lua_lsp_document_highlights/document_highlights/caret_move"],
  },
  {
    labels: ["Find References"],
    measures: ["lspFindUsages"],
    projects: ["lua_lsp_find_references/find_references/same_file", "lua_lsp_find_references/find_references/project"],
  },
  {
    labels: ["Folding Range on File Open"],
    measures: ["lspFoldingRangeOnFileOpen"],
    projects: ["lua_lsp_folding_range/folding_range/file_open"],
  },
  {
    labels: ["Formatting on Reformat Code"],
    measures: ["lspFormattingOnReformatCode"],
    projects: ["lua_lsp_formatting/formatting/reformat_code"],
  },
  {
    labels: ["Go to Definition"],
    measures: ["execute_editor_gotodeclaration"],
    projects: ["lua_lsp_go_to_definition/go_to_definition/same_file", "lua_lsp_go_to_definition/go_to_definition/other_file"],
  },
  {
    labels: ["Hover"],
    measures: ["lspHover"],
    projects: ["lua_lsp_hover/hover/function_of_other_file"],
  },
  {
    labels: ["Inlay Hint on File Open"],
    measures: ["lspInlayHintOnFileOpen"],
    projects: ["lua_lsp_inlay_hint/inlay_hint/file_open"],
  },
  {
    labels: ["On-Type Formatting on Enter"],
    measures: ["lspOnTypeFormattingOnEnter"],
    projects: ["lua_lsp_on_type_formatting/on_type_formatting/enter"],
  },
  {
    labels: ["Prepare Rename on Start", "Rename on Enter"],
    measures: ["lspPrepareRenameOnStart", "lspRenameOnEnter"],
    projects: ["lua_lsp_rename/rename/same_file", "lua_lsp_rename/rename/project"],
  },
  {
    labels: ["Semantic Tokens on File Open"],
    measures: ["lspSemanticTokensOnFileOpen"],
    projects: ["lua_lsp_semantic_tokens/semantic_tokens/file_open"],
  },
]

const charts = combineCharts(declaration)
</script>
