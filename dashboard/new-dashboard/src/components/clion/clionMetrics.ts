import { MAIN_METRICS } from "../../util/mainMetrics"

const CLION_SPECIFIC_METRICS = [
  // indexing
  "backendIndexingTimeMs",
  "nova_collecting_files_ms",
  "nova_indexing_time_ms",
  "nova_total_time_ms",
  "nova_resolving_references_ms",
  "workspaceModel.updates.ms",
  "workspaceModel.replace.project.model.ms",
  // typing and highlighting
  "typing#latency#mean_value",
  "typing#latency#max",
  "fus_file_types_usage_duration_ms",
  "fus_file_types_usage_time_to_show_ms",
  // completion
  "fus_time_to_show_90p",
  "callInlineCompletionOnCompletion",
  // search everywhere
  "searchEverywhere_first_elements_added",
  // debugger
  "fus_debug_session_initialized_ms",
  "fus_frame_variables_computed_ms",
  "debugStep_out",
  "debugStep_over",
  "evaluateExpression#mean_value",
  // go to declaration
  "clionGotoDeclaration",
  // find usages
  "%syncAction FindUsages",
  "findUsagesInToolWindow",
  "findUsagesInToolWindow#number",
  // run gutter
  "waitFirstTestGutter",
  "fus_popup_latency_ms",
  // lagging/latency
  "ui.lagging#average",
  "ui.lagging#max",
  "ui.lagging#percentage_share",
  // memory
  "nova_total_memory_mb",
  "rd.memory.allocatedManagedMemoryMb/afterIndexing",
  "JVM.heapUsageMb/afterIndexing",
]

export const CLION_MAIN_METRICS: string[] = [...MAIN_METRICS, ...CLION_SPECIFIC_METRICS]
