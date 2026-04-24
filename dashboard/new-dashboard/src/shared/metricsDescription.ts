import { GRADLE_METRICS_NEW_DASHBOARD } from "../components/intelliJ/build-tools/gradle/gradle-metrics"

/**
 * Map of metric names to their descriptions.
 * The syntax is either:
 * ```
 * ["metricName", "MetricDescription"]
 * ["metricNamePrefix#*", "MetricDescription"]
 * ```
 * where `*` is a wildcard that matches any suffix.
 * For example, `processingTime#*` matches `processingTime#java` and `processingTime#cpp`.
 *
 * If you want to provide URL, use `metricInfo` function:
 * ```
 * ["metricName", metricInfo("MetricDescription", "https://example.com")]
 * ```
 */
export const metricsDescription: Map<string, string | MetricInfo> = new Map<string, string | MetricInfo>([
  // FUS events (some of them used for (mega)APDEX calculations)
  ["fus_file_types_usage_duration_ms", 'FUS event with groupID="file.types.usage" eventID="open" eventField="duration_ms", (mega)APDEX: "File Openings: Code Loaded"'],
  ["fus_file_types_usage_time_to_show_ms", 'FUS event with groupID="file.types.usage" eventID="open" eventField="time_to_show", (mega)APDEX: "File Openings: Tab Shown"'],
  [
    "fus_daemon_finished_full_duration_since_started_ms",
    'FUS event with groupID="daemon" eventID="finished" eventField="full_duration_since_started_ms. Full highlighting duration since the file was modified and/or dumb mode status changed. It should be equal to the sum of segments."',
  ],
  ["fus_completion_duration_sum", 'SUM of FUS events with groupID="completion" eventID="finished" eventField="duration"'],
  ["fus_completion_duration_90p", '90 percentile of FUS events with groupID="completion" eventID="finished" eventField="duration"'],
  ["fus_time_to_show_90p", '90 percentile of FUS events with groupID="completion" eventID="finished" eventField="time_to_show"'],
  ["fus_dumb_indexing_time", 'FUS event with groupID="indexing.statistics" eventID="finished" eventField="finished" eventField="indexing_activity_type=dumb_indexing"'],
  ["fus_scanning_time", 'FUS event with groupID="indexing.statistics" eventID="finished" eventField="finished" eventField="indexing_activity_type=scanning"'],
  ["fus_git_branches_checkout_operation", 'FUS event with groupID="git.branches" eventID="checkout.checkout_operation.finished" eventField="duration_ms"'],
  ["fus_git_branches_vfs_refresh", 'FUS event with groupID="git.branches" eventID="checkout.vfs_refresh.finished" eventField="duration_ms"'],
  ["fus_vcs_commit_duration", 'FUS event with groupID="vcs" eventID="commit.finished" eventField="duration_ms"'],
  ["fus_find_usages_all", 'FUS event with groupID="usage.view" eventID="finished" eventField="duration_ms"'],
  ["fus_find_usages_first", 'FUS event with groupID="usage.view" eventID="finished" eventField="duration_first_results_ms. Old startup metric"'],
  ["fus_startup_totalDuration", 'FUS event with groupID="startup" eventID="totalDuration" eventField="duration"'],
  ["fus_reopen_startup_frame_became_interactive", 'FUS event with groupID="reopen.project.startup.performance" eventID="frame.became.interactive" eventField="duration_ms"'],
  ["fus_reopen_startup_first_ui_shown", 'FUS event with groupID="reopen.project.startup.performance" eventID="first.ui.shown" eventField="duration_ms"'],
  ["fus_reopen_startup_frame_became_visible", 'FUS event with groupID="reopen.project.startup.performance" eventID="frame.became.visible" eventField="duration_ms"'],
  [
    "fus_reopen_startup_code_loaded_and_visible_in_editor",
    'FUS event with groupID="reopen.project.startup.performance" eventID="code.loaded.and.visible.in.editor" eventField="duration_ms. New main metric for startup"',
  ],
  [
    "fus_gradle.sync",
    'Difference between durations of FUS events with groupID="build.gradle.import" eventID="gradle.sync.finished" eventField="duration_ms" and groupID="build.gradle.import" eventID="gradle.sync.started" eventField="duration_ms"',
  ],
  ["fus_PROJECT_RESOLVERS", 'FUS event with groupID="build.gradle.import" eventID="phase.finished" phase="PROJECT_RESOLVERS" eventField="duration_ms"'],
  ["fus_GRADLE_CALL", 'FUS event with groupID="build.gradle.import" eventID="phase.finished" phase="GRADLE_CALL" eventField="duration_ms"'],
  ["fus_DATA_SERVICES", 'FUS event with groupID="build.gradle.import" eventID="phase.finished" phase="DATA_SERVICES" eventField="duration_ms"'],

  //completion
  ["completion", "Total time of each completion invocation in test. Completion invocation time is a time that it takes to load all completion variants."],
  ["completion#mean_value", "Mean value of all completion invocation in test. Completion invocation time is a time that it takes to load all completion variants."],
  //find usages
  ["findUsage_popup", "Time to show the find usages popup"],
  ["findUsages", "Time to show all usages in the popup"],
  ["findUsages#number", "Number of found usages"],
  ["findUsages#mean_value", "Mean time to show all usages in the popup"],
  ["findUsages_firstUsage", "Time to show the first usage in the popup"],
  //analysis
  ["firstCodeAnalysis", "Time it takes to perform code analysis on file opening"],
  ["localInspections", "Sum time of all analysis. From Daemon#restart till DaemonListener#daemonFinished."],
  ["localInspections#mean_value", "Code analysis mean time. From Daemon#restart till DaemonListener#daemonFinished."],
  ["runDaemon/executionTime", "Time it takes to complete a first daemon run. It might be restarted so it's not a full time."],
  ["globalInspections", "Time of all inspections runned in batch mode (Inspect Project)."],
  //indexing
  ["indexSize", "Index size in (in kb)"],
  ["numberOfIndexedFiles", "Number of indexed files"],
  ["processingSpeedAvg#*", "Speed of indexing file type (in kb/s)"],
  ["processingTime#*", "CPU time spent on processing file type"],
  ["indexingTimeWithoutPauses", "Indexing time without interruptions"],
  ["scanningTimeWithoutPauses", "Scanning time without interruptions"],
  ["pageLoad", "Number of regular Pages' loads."],
  ["pageMiss", "If the needed page has not existed in the main memory (RAM), it is known as PAGE MISS. The metric displays the number of unsuccessful Pages' obtainment."],
  [
    "pageHit",
    "CPU attempts to obtain a needed page from main memory and the page exists in main memory (RAM), it is referred to as a PAGE HIT. This metric displays the number of successful Pages' obtainment.",
  ],
  //typing
  ["typing", "Typing executing time (usually equal to number of typed characters times delay between key presses)"],
  ["typing#average_awt_delay", "How long on average it takes to process a single empty AWT event in the queue during typing."],
  ["typing#max_awt_delay", "Max value of how long it takes to process a single empty AWT event in the queue during typing."],
  [
    "typing#latency",
    "Average time in ms of inserting a letter in the Editor (approximation of how long does it take from keyboard press till the appearance of the letter); measured during typing",
  ],
  //refactorings
  ["performInlineRename#mean_value", "Rename mean time. Find usages is not included, only actual rename time is counted"],
  ["startInlineRename#mean_value", "Mean time to prepare rename template in the current file"],
  ["prepareForRename#mean_value", "Mean time to perform find usages and other preparations such as conflict detection for write phase of rename"],
  ["fus_refactoring_usages_searched", "Mean time to perform find usages during refactorings"],
  ["moveFiles#mean_value", "Mean time to perform move files refactoring: with find usages, conflict detection and actual move"],
  ["moveFiles_back#mean_value", "Mean time to restore project as it was before move files"],
  ["moveDeclarations#mean_value", "Mean time to perform move files refactoring: with find usages, conflict detection and actual move"],
  ["moveDeclarations_back#mean_value", "Mean time to restore project as it was before move declarations"],

  //editor actions
  ["execute_editor_optimizeimports", "Time to execute optimize imports action in the editor"],
  ["execute_editor_paste", "Time to execute paste action in the editor"],
  ["convertJavaToKotlin", "Time to execute J2K action in the editor"],

  //GC
  ["freedMemoryByGC", metricInfo("Freed memory by GC (in Mb/s)", "https://github.com/chewiebug/GCViewer#readme")],
  ["fullGCPause", metricInfo("Time that full GC was active (IDE is fully paused)", "https://github.com/chewiebug/GCViewer#readme")],
  ["gcPause", metricInfo("Time spent in GC (including minor collections without pausing)", "https://github.com/chewiebug/GCViewer#readme")],
  ["gcPauseCount", metricInfo("Number of minor GCs pauses", "https://github.com/chewiebug/GCViewer#readme")],
  //others
  ["searchEverywhere_*", "Time to fill all search everywhere results"],
  ["FileStructurePopup", "Time needed to display and fill a popup with information about the structure of a given file."],
  //indexing / VFS / VCS indexing
  ["dumbModeWithPauses", "Total time spent in dumb mode (indexes not ready), including pauses like GC and UI freezes"],
  ["vfs_initial_refresh", "Duration of the initial VFS refresh that syncs on-disk files with the VFS cache on project opening"],
  ["vcs-log-indexing", "Duration of VCS Log background indexing of commits for fast search and filtering"],
  //build
  ["build_compilation_duration", "Total elapsed time of a project build (module compile/rebuild/recompile via ProjectTaskManager)"],
  //search everywhere
  ["searchEverywhere", "End-to-end time of a Search Everywhere operation: popup invocation, typing, and optional selection/closing"],
  ["searchEverywhere_dialog_shown", "Time from triggering Search Everywhere until the popup is fully displayed"],
  //vcs
  ["showFileHistory", "Time from invoking Show File History until the initial data pack is loaded and rendered"],
  //menus
  ["%expandMainMenu", "Time to recursively expand all actions in the main menu bar (GROUP_MAIN_MENU)"],
  ["%expandProjectMenu", "Time to expand the Project View context menu popup (GROUP_PROJECT_VIEW_POPUP)"],
  ["%expandEditorMenu", "Time to expand the editor right-click context menu (GROUP_EDITOR_POPUP)"],
  //new file
  ["createKotlinFile", "Time to create a new Kotlin file via the New File template action"],
  ["createJavaFile", "Time to create a new Java class/file via the New File action (JavaDirectoryService.createClass)"],
  //highlighting
  ["highlighting", "Total time of background daemon syntax and semantic highlighting analysis on a file"],
  ["typingCodeAnalyzing", "Time of daemon code analysis triggered after typing, from typing completion until daemon finishes"],
  //rename
  ["startInlineRename", "Sum time of all inline rename invocations, from trigger through template preparation and editor setup"],
  //debug
  ["debugRunConfiguration", "Time from launching a debug run configuration until the debugger first pauses at a breakpoint"],
  ["debugStep_into", "Time from invoking Step Into in the debugger until the next sessionPaused callback fires"],
  //AI completion quality
  ["MatchedRatio", "Average length of accepted completion minus prefix, normalized by expected text (AI completion quality)"],
  ["SyntaxErrorsSessionRatio", "Ratio of completion sessions that left syntax errors in the resulting code"],
  ["EditSimilarity", "Maximum Levenshtein-based similarity between proposed completion suggestions and expected text"],
  //benchmark
  ["attempt.mean.ms", "Mean duration in milliseconds across all benchmark attempt spans"],
  //gc
  ["freedMemory", metricInfo("Amount of memory (bytes) freed by GC during the test, parsed from GCViewer output", "https://github.com/chewiebug/GCViewer#readme")],
  ["test#average_awt_delay", "The average time it takes to process a single empty AWT event in the queue during the whole test."],
  ["showQuickFixes", "Time to show the quick fixes after calling Alt + Enter."],
  [
    "attempt.mad.ms",
    metricInfo(
      "MAD (Median Absolute Deviation) of attempt durations in ms. The MAD is a robust statistic, being more resilient to outliers in a data set than the standard deviation.",
      "https://en.m.wikipedia.org/wiki/Median_absolute_deviation"
    ),
  ],
  ...GRADLE_METRICS_NEW_DASHBOARD,

  // Workspace Model
  [
    "workspaceModel.updates.count",
    "Total number of changes made to the project model, including modifying entities, changing configuration, and updating dependencies, among others.",
  ],
  [
    "workspaceModel.updates.ms",
    "Total time spent on processing modifications to the workspace entities including time required in calling update handlers, collecting changes, initializing bridging operations, and generating snapshots",
  ],
  ["workspaceModel.mutableEntityStorage.to.snapshot.ms", "The time taken to create a snapshot of the mutable entity storage"],
  ["workspaceModel.mutableEntityStorage.replace.by.source.ms", "The time taken to replace entities in the mutable entity storage by source"],
  ["workspaceModel.mutableEntityStorage.add.diff.ms", "The time taken to add differences to the mutable entity storage"],
  ["workspaceModel.loading.from.cache.ms", "The time taken to load the workspace model from cache"],
  ["workspaceModel.do.save.caches.ms", "The time taken to save caches of the workspace model"],
  ["workspaceModel.mutableEntityStorage.collect.changes.ms", "The time taken to collect changes in the mutable entity storage"],
  ["workspaceModel.mutableEntityStorage.add.entity.ms", "The time taken to add an entity to the mutable entity storage"],
  ["jps.load.project.to.empty.storage.ms", "The time taken to load a project into empty storage using the JPS"],
  ["jps.project.serializers.load.ms", "The time taken to load project serializers in the JPS"],
  ["jps.project.serializers.save.ms", "The time taken to save project serializers in the JPS"],
  ["jps.facet.change.listener.process.change.events.ms", "The time taken to process change events by the facet change listener in the JPS"],
])

export interface MetricInfo {
  description: string
  url?: string
}

function metricInfo(description: string, url?: string): MetricInfo {
  return { description, url }
}

function extractMainPrefix(inputString: string): string {
  const regex = /^(\w+#).*/
  const match = regex.exec(inputString)
  return match ? match[1] : "non-matching"
}

export function getMetricDescription(metric: string | undefined): MetricInfo | null {
  if (metric == undefined) return null
  const metricDescription = metricsDescription.get(metric) ?? metricsDescription.get(extractMainPrefix(metric) + "*") ?? null
  return typeof metricDescription == "string" ? metricInfo(metricDescription) : metricDescription
}
