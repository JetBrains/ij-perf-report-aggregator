import { BetterDirection } from "./changeDetector/algorithm"
import { GRADLE_METRICS_NEW_DASHBOARD } from "../components/intelliJ/build-tools/gradle/gradle-metrics"
import type { MeasureUnit } from "../components/common/formatter"
import { SERIES_NAME_SEPARATOR } from "../components/common/DataQueryExecutor"

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
 * To declare a unit, documentation URL, or better-direction, use the `metricInfo` function
 * (`metricInfo(description, unit?, url?, betterDirection?)`):
 * ```
 * ["metricName", metricInfo("MetricDescription", "milliseconds", "https://example.com")]
 * ```
 */
export const metricsDescription: Map<string, string | MetricInfo> = new Map<string, string | MetricInfo>([
  // FUS events (some of them used for (mega)APDEX calculations)
  [
    "fus_file_types_usage_duration_ms",
    metricInfo('Time from file open until code is loaded. Used for (mega)APDEX "File Openings: Code Loaded". FUS: file.types.usage/open/duration_ms.', "milliseconds"),
  ],
  [
    "fus_file_types_usage_time_to_show_ms",
    metricInfo('Time from file open until the editor tab is shown. Used for (mega)APDEX "File Openings: Tab Shown". FUS: file.types.usage/open/time_to_show.', "milliseconds"),
  ],
  [
    "fus_daemon_finished_full_duration_since_started_ms",
    metricInfo(
      "Full highlighting duration since the file was modified or dumb mode status changed. Equal to the sum of highlighting segments. FUS: daemon/finished/full_duration_since_started_ms.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-295/Highlighting-Metric"
    ),
  ],
  [
    "fus_completion_duration_sum",
    metricInfo(
      "Sum of all code completion durations. Reported by FUS telemetry. FUS: completion/finished/duration.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-304/Code-Completion-Metric"
    ),
  ],
  [
    "fus_completion_duration_90p",
    metricInfo(
      "90th percentile of code completion duration. Reported by FUS telemetry. FUS: completion/finished/duration.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-304/Code-Completion-Metric"
    ),
  ],
  [
    "fus_time_to_show_90p",
    metricInfo(
      "90th percentile time until the completion popup is shown. Reported by FUS telemetry. FUS: completion/finished/time_to_show.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-304/Code-Completion-Metric"
    ),
  ],
  [
    "fus_dumb_indexing_time",
    metricInfo(
      "Time spent indexing in dumb mode. Reported by FUS telemetry. FUS: indexing.statistics/finished/indexing_activity_type=dumb_indexing.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-296/Indexing-Metric"
    ),
  ],
  [
    "fus_scanning_time",
    metricInfo(
      "Time spent scanning the file system for changes. Reported by FUS telemetry. FUS: indexing.statistics/finished/indexing_activity_type=scanning.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-296/Indexing-Metric"
    ),
  ],
  [
    "fus_git_branches_checkout_operation",
    metricInfo("Time spent on the Git checkout operation. FUS: git.branches/checkout.checkout_operation.finished/duration_ms.", "milliseconds"),
  ],
  [
    "fus_git_branches_vfs_refresh",
    metricInfo("Time spent refreshing the virtual file system after Git checkout. FUS: git.branches/checkout.vfs_refresh.finished/duration_ms.", "milliseconds"),
  ],
  ["fus_vcs_commit_duration", metricInfo("Time to complete a VCS commit operation. FUS: vcs/commit.finished/duration_ms.", "milliseconds")],
  ["fus_find_usages_all", metricInfo("Time to find all usages. FUS: usage.view/finished/duration_ms.", "milliseconds")],
  ["fus_find_usages_first", metricInfo("Time to find the first usage result. FUS: usage.view/finished/duration_first_results_ms.", "milliseconds")],
  [
    "fus_startup_totalDuration",
    metricInfo("Total IDE startup duration. FUS: startup/totalDuration/duration.", "milliseconds", "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric"),
  ],
  [
    "fus_reopen_startup_frame_became_interactive",
    metricInfo(
      "Time from project reopen until the IDE frame became interactive. FUS: reopen.project.startup.performance/frame.became.interactive/duration_ms.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric"
    ),
  ],
  [
    "fus_reopen_startup_first_ui_shown",
    metricInfo(
      "Time from project reopen until the first UI frame was shown. FUS: reopen.project.startup.performance/first.ui.shown/duration_ms.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric"
    ),
  ],
  [
    "fus_reopen_startup_frame_became_visible",
    metricInfo(
      "Time from project reopen until the IDE frame became visible. FUS: reopen.project.startup.performance/frame.became.visible/duration_ms.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric"
    ),
  ],
  [
    "fus_reopen_startup_code_loaded_and_visible_in_editor",
    metricInfo(
      "Time from project reopen until code is loaded and visible in the editor. Primary startup metric. FUS: reopen.project.startup.performance/code.loaded.and.visible.in.editor/duration_ms.",
      "milliseconds",
      "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric"
    ),
  ],
  [
    "fus_gradle.sync",
    metricInfo(
      "Duration of Gradle sync. Difference between sync finished and sync started FUS events. FUS: build.gradle.import/gradle.sync.finished minus gradle.sync.started/duration_ms.",
      "milliseconds"
    ),
  ],
  [
    "fus_PROJECT_RESOLVERS",
    metricInfo("Time spent resolving Gradle projects. Phase of Gradle sync. FUS: build.gradle.import/phase.finished[PROJECT_RESOLVERS]/duration_ms.", "milliseconds"),
  ],
  ["fus_GRADLE_CALL", metricInfo("Time spent calling Gradle daemon. Phase of Gradle sync. FUS: build.gradle.import/phase.finished[GRADLE_CALL]/duration_ms.", "milliseconds")],
  [
    "fus_DATA_SERVICES",
    metricInfo("Time spent loading Gradle data services. Phase of Gradle sync. FUS: build.gradle.import/phase.finished[DATA_SERVICES]/duration_ms.", "milliseconds"),
  ],

  //completion
  ["completion", metricInfo("Total time to load all completion variants for each invocation.", "milliseconds")],
  ["completion#mean_value", metricInfo("Mean time to load all completion variants across invocations.", "milliseconds")],
  //code typing
  [
    "codeTyping#mean_value",
    metricInfo(
      "Mean time of a code typing iteration. A code typing test types a piece of code with medium complexity from beginning to end, invoking completion and waiting for code analysis at predetermined points.",
      "milliseconds"
    ),
  ],
  //find usages
  ["findUsage_popup", metricInfo("Time to show the find usages popup.", "milliseconds")],
  ["findUsages", metricInfo("Time to find and display all usages in the popup.", "milliseconds")],
  ["findUsages#number", metricInfo("Number of usages found. Stable for a fixed query; a drop signals a resolve regression.", "counter", undefined, "stable")],
  ["findUsages#mean_value", metricInfo("Mean time to show all usages in the popup.", "milliseconds")],
  ["findUsages_firstUsage", metricInfo("Time until the first usage appears in the popup.", "milliseconds")],
  //analysis
  ["firstCodeAnalysis", metricInfo("Time to perform syntax and semantic highlighting when a file opens for the first time (cold daemon pass).", "milliseconds")],
  ["localInspections", metricInfo("Total time of on-the-fly code analysis. Measured from daemon restart to daemon finish.", "milliseconds")],
  ["localInspections#mean_value", metricInfo("Mean time of on-the-fly code analysis. Measured from daemon restart to daemon finish.", "milliseconds")],
  ["runDaemon/executionTime", metricInfo("Time to complete a first daemon run. It might be restarted so it is not a full time.", "milliseconds")],
  ["globalInspections", metricInfo("Time to run all inspections in batch mode (Inspect Code / Inspect Project).", "milliseconds")],
  //indexing
  ["indexSize", metricInfo("Total size of indexes written to disk after indexing.", "kilobytes")],
  ["lexingSize#*", metricInfo("Size of the lexed content of a file type (in bytes).", "bytes")],
  ["parsingSize#*", metricInfo("Size of the parsed content of a file type (in bytes).", "bytes")],
  ["numberOfIndexedFiles", metricInfo("Total files that went through the indexing pipeline.", "counter", undefined, "stable")],
  ["processingSpeedAvg#*", metricInfo("Average indexing throughput for a file type (kB/s).", "kilobytesPerSecond", undefined, "higher")],
  ["lexingSpeed#*", metricInfo("Lexing speed of a file type (in kB/s).", "kilobytesPerSecond", undefined, "higher")],
  ["parsingSpeed#*", metricInfo("Parsing speed of a file type (in kB/s).", "kilobytesPerSecond", undefined, "higher")],
  ["processingTime#*", metricInfo("CPU time spent on processing files of this type.", "milliseconds")],
  ["lexingTime#*", metricInfo("Time the lexer spent tokenizing files of this type.", "milliseconds")],
  ["parsingTime#*", metricInfo("Time the parser spent building PSI for files of this type.", "milliseconds")],
  [
    "processingSpeedWorst#*",
    metricInfo("Worst-case indexing throughput for a file type (kB/s). Low values indicate pathological files.", "kilobytesPerSecond", undefined, "higher"),
  ],
  [
    "processingSpeedOfBaseLanguageWorst#*",
    metricInfo("Worst-case indexing throughput mapped to the base language (kB/s). Low values indicate pathological files.", "kilobytesPerSecond", undefined, "higher"),
  ],
  ["numberOfIndexedFiles#*", metricInfo("Number of files of this type that were indexed.", "counter", undefined, "stable")],
  ["indexingTimeWithoutPauses", metricInfo("Time to build indexes, excluding paused intervals.", "milliseconds")],
  ["scanningTimeWithoutPauses", metricInfo("Time to scan the file system for changes before indexing, excluding pauses.", "milliseconds")],
  ["pageLoad", metricInfo("Number of regular Pages' loads.", "counter", undefined, "stable")],
  ["pageMiss", metricInfo("Number of unsuccessful page loads from main memory.", "counter", undefined, "lower")],
  ["pageHit", metricInfo("Number of successful page loads from main memory.", "counter", undefined, "higher")],
  //typing
  ["typing", metricInfo("Total time to type the sample code.", "milliseconds")],
  ["typing#average_awt_delay", metricInfo("Average time to process a single empty AWT event in the queue during typing.", "milliseconds")],
  ["typing#max_awt_delay", metricInfo("Maximum time to process a single empty AWT event in the queue during typing.", "milliseconds")],
  [
    "typing#latency",
    metricInfo("Average time in ms to insert a letter in the editor, measured during typing — approximates the delay from key press to the letter appearing.", "milliseconds"),
  ],
  //refactorings
  ["performInlineRename#mean_value", metricInfo("Mean time to perform an inline rename. Find usages is not included.", "milliseconds")],
  ["startInlineRename#mean_value", metricInfo("Mean time to prepare rename template in the current file.", "milliseconds")],
  ["prepareForRename#mean_value", metricInfo("Mean time to perform find usages and other preparations for rename.", "milliseconds")],
  ["fus_refactoring_usages_searched", metricInfo("Mean time to search for usages during refactoring operations.", "milliseconds")],
  ["moveFiles#mean_value", metricInfo("Mean time to perform move files refactoring.", "milliseconds")],
  ["moveFiles_back#mean_value", metricInfo("Mean time to restore the project after move files.", "milliseconds")],
  ["moveDeclarations#mean_value", metricInfo("Mean time to perform move declarations refactoring.", "milliseconds")],
  ["moveDeclarations_back#mean_value", metricInfo("Mean time to restore the project after move declarations.", "milliseconds")],

  //editor actions
  ["execute_editor_optimizeimports", metricInfo("Time to execute optimize imports action in the editor.", "milliseconds")],
  ["execute_editor_paste", metricInfo("Time to execute paste action in the editor.", "milliseconds")],
  ["convertJavaToKotlin", metricInfo("Time to execute J2K action in the editor.", "milliseconds")],

  //GC
  // The GCViewer-sourced memory metrics below (freedMemoryByGC, freedMemoryByFullGC, totalHeapUsedMax)
  // are mebibytes, despite being stored raw. GCViewer keeps memory in KiB internally and steps down at
  // 1024 boundaries (B→K→M→G), so a value tagged "M" already equals bytes / 1024² — the same binary unit
  // MEM.*/JVM.* reach by dividing raw bytes by /1024/1024. GCLogAnalyzer keeps the "M" value as-is.
  [
    "freedMemoryByGC",
    metricInfo("Total memory freed by garbage collection. Tracks allocation churn; lower is better.", "mebibytes", "https://github.com/chewiebug/GCViewer#readme", "lower"),
  ],
  [
    "freedMemoryByFullGC",
    metricInfo("Memory reclaimed by full GC only. High values signal object promotion pressure.", "mebibytes", "https://github.com/chewiebug/GCViewer#readme", "lower"),
  ],
  ["fullGCPause", metricInfo("Time the IDE was fully paused during full GC.", "milliseconds", "https://github.com/chewiebug/GCViewer#readme")],
  ["gcPause", metricInfo("Total time spent in garbage collection, including minor collections.", "milliseconds", "https://github.com/chewiebug/GCViewer#readme")],
  ["gcPauseCount", metricInfo("Number of GC pause events.", "counter", "https://github.com/chewiebug/GCViewer#readme", "lower")],
  [
    "fullGcPauseCount",
    metricInfo("Number of full GC pause events. The full-GC-only counterpart of gcPauseCount.", "counter", "https://github.com/chewiebug/GCViewer#readme", "lower"),
  ],
  ["totalHeapUsedMax", metricInfo("Peak JVM heap usage during the test.", "mebibytes")],
  ["bsp.used.at.exit.mb", metricInfo("Heap used at exit.", "mebibytes")],
  ["bsp.used.after.sync.mb", metricInfo("Heap used after sync.", "mebibytes")],

  // GC — GCViewer summary (parsed from the GC log)
  ["throughput", metricInfo("Percentage of run time the application runs rather than pausing for GC.", "counter", "https://github.com/chewiebug/GCViewer#readme", "higher")],
  [
    "maxPause",
    metricInfo(
      "Longest single stop-the-world GC pause across all young, mixed, and full GC pauses combined. Represents the worst case GC latency",
      "milliseconds",
      "https://github.com/chewiebug/GCViewer#readme"
    ),
  ],
  [
    "maxFullGCPause",
    metricInfo(
      "Longest single full GC pause. Worst-case stop-the-world latency from a full GC; the full-GC-only counterpart of maxPause.",
      "milliseconds",
      "https://github.com/chewiebug/GCViewer#readme"
    ),
  ],
  [
    "accumPause",
    metricInfo("Accumulated GC pause time across collection kinds. gcPause and fullGCPause split it by kind.", "milliseconds", "https://github.com/chewiebug/GCViewer#readme"),
  ],
  [
    "avgfootprintAfterFullGC",
    metricInfo(
      "Average heap retained after a full GC — the steady-state live set. totalHeapUsedMax is peak usage, including uncollected garbage.",
      "mebibytes",
      "https://github.com/chewiebug/GCViewer#readme"
    ),
  ],
  ["totalPermUsedMax", metricInfo("Peak metaspace (class-metadata) usage.", "mebibytes", "https://github.com/chewiebug/GCViewer#readme")],
  //others
  ["searchEverywhere_*", "Time to fill all search everywhere results."],
  ["FileStructurePopup", "Time needed to display and fill a popup with information about the structure of a given file."],
  //indexing / VFS / VCS indexing
  [
    "dumbModeWithPauses",
    metricInfo(
      "Total time spent in dumb mode (indexes not ready), including pauses from GC and UI freezes.",
      "milliseconds",
      "https://plugins.jetbrains.com/docs/intellij/smart-mode-and-dumb-mode.html"
    ),
  ],
  ["dumbModeTimeWithPauses", metricInfo("Total time the IDE spent in dumb mode (indexes not ready), including pauses from GC and UI freezes.", "milliseconds")],
  ["pausedTimeInIndexingOrScanning", metricInfo("Time the indexing or scanning pipeline was paused.", "milliseconds")],
  ["numberOfRunsOfIndexing", metricInfo("Number of indexing passes. Higher means more incremental re-indexing.", "counter", undefined, "lower")],
  ["numberOfRunsOfScannning", metricInfo("Number of file system scanning passes. Extra passes mean wasted work.", "counter", undefined, "lower")],
  ["numberOfScannedFiles", metricInfo("Total files scanned during indexing. Compared with indexed files, shows filter efficiency.", "counter", undefined, "stable")],
  ["numberOfIndexedFilesWritingIndexValue", metricInfo("Files that produced index data this run. Should stay flat for a fixed project.", "counter", undefined, "stable")],
  ["numberOfIndexedFilesWithNothingToWrite", metricInfo("Files indexed but producing no index data. High values indicate wasted indexing work.", "counter", undefined, "stable")],
  ["numberOfFilesIndexedByExtensions", metricInfo("Files recognized by extension-based file type detection.", "counter", undefined, "stable")],
  ["numberOfFilesIndexedWithoutExtensions", metricInfo("Files indexed without a recognized file extension.", "counter", undefined, "stable")],
  ["vfs_initial_refresh", metricInfo("Duration of the initial VFS refresh that syncs on-disk files with the VFS cache on project opening.", "milliseconds")],
  ["vcs-log-indexing", metricInfo("Duration of VCS Log background indexing of commits for fast search and filtering.", "milliseconds")],

  // Write actions
  ["writeAction.count", metricInfo("Number of write actions executed.", "counter", undefined, "stable")],
  ["writeAction.wait.ms", metricInfo("Total time spent waiting for write actions.", "milliseconds")],
  ["writeAction.max.wait.ms", metricInfo("Longest single write action wait. Spikes indicate blocking.", "milliseconds")],
  ["writeAction.median.wait.ms", metricInfo("Median write action wait time. Represents typical contention.", "milliseconds")],

  // AWT
  ["AWTEventQueue.dispatchTimeTotal", metricInfo("Total AWT event queue dispatch time. High values indicate UI thread contention.", "milliseconds")],
  //build
  ["build_compilation_duration", metricInfo("Total elapsed time of a project build (module compile/rebuild/recompile via ProjectTaskManager).", "milliseconds")],
  //search everywhere
  ["searchEverywhere", metricInfo("End-to-end time of a Search Everywhere operation: popup invocation, query typing, and optional selection.", "milliseconds")],
  ["searchEverywhere_dialog_shown", metricInfo("Time from triggering Search Everywhere until the popup is fully displayed.", "milliseconds")],
  //vcs
  ["showFileHistory", metricInfo("Time from invoking Show File History until the initial data pack is loaded and rendered.", "milliseconds")],
  //menus
  ["%expandMainMenu", metricInfo("Time to recursively expand all actions in the main menu bar (GROUP_MAIN_MENU).", "milliseconds")],
  ["%expandProjectMenu", metricInfo("Time to expand the Project View context menu popup (GROUP_PROJECT_VIEW_POPUP).", "milliseconds")],
  ["%expandEditorMenu", metricInfo("Time to expand the editor right-click context menu (GROUP_EDITOR_POPUP).", "milliseconds")],
  //new file
  ["createKotlinFile", metricInfo("Time to create a new Kotlin file via the New File template action.", "milliseconds")],
  ["createJavaFile", metricInfo("Time to create a new Java class/file via the New File action (JavaDirectoryService.createClass).", "milliseconds")],
  //highlighting
  ["highlighting", metricInfo("Total time of background daemon syntax and semantic highlighting analysis on a file.", "milliseconds")],
  ["typingCodeAnalyzing", metricInfo("Time of daemon code analysis triggered after typing, from typing completion until daemon finishes.", "milliseconds")],
  //rename
  ["startInlineRename", metricInfo("Sum time of all inline rename invocations, from trigger through template preparation and editor setup.", "milliseconds")],
  //debug
  ["debugRunConfiguration", metricInfo("Time from launching a debug configuration to the first pause at a breakpoint.", "milliseconds")],
  ["debugStep_into", metricInfo("Time from invoking Step Into until the debugger pauses again.", "milliseconds")],
  //AI completion quality
  ["MatchedRatio", "Average length of accepted completion minus prefix, normalized by expected text (AI completion quality)."],
  ["SyntaxErrorsSessionRatio", "Ratio of completion sessions that left syntax errors in the resulting code."],
  ["EditSimilarity", "Maximum Levenshtein-based similarity between proposed completion suggestions and expected text."],
  //benchmark
  ["attempt.mean.ms", metricInfo("Mean duration in milliseconds across all benchmark attempt spans.", "milliseconds")],
  //GC — G1GC
  [
    "g1gcConcurrentMarkCycles",
    metricInfo(
      "Number of G1 concurrent marking cycles. Reflects old-gen GC pressure; lower is better.",
      "counter",
      "https://docs.oracle.com/en/java/javase/21/gctuning/garbage-first-g1-garbage-collector1.html#GUID-93D75606-CAC0-4D79-983B-2DCC7E22D13A",
      "lower"
    ),
  ],
  [
    "g1gcConcurrentMarkTimeMs",
    metricInfo(
      "Wall-clock time spent in G1 concurrent marking. Competes with UI threads for CPU.",
      "milliseconds",
      "https://docs.oracle.com/en/java/javase/21/gctuning/garbage-first-g1-garbage-collector1.html#GUID-93D75606-CAC0-4D79-983B-2DCC7E22D13A"
    ),
  ],
  [
    "g1gcHeapShrinkageCount",
    metricInfo("Number of times G1 returned heap memory to the OS. Zero means heap is pinned at maximum.", "counter", "https://openjdk.org/jeps/346", "lower"),
  ],
  ["g1gcHeapShrinkageMegabytes", metricInfo("Total heap memory returned to the OS by G1 shrinkage.", "mebibytes", "https://openjdk.org/jeps/346", "lower")],

  // JVM — GC (JVM-reported, independent of GCViewer)
  ["JVM.GC.collectionTimesMs", metricInfo("Total GC time as reported by the JVM. Independent cross-check against GCViewer-based `gcPause`.", "milliseconds")],
  ["JVM.GC.collections", metricInfo("Total number of GC collections as reported by the JVM.", "counter")],

  // JVM — Runtime
  ["JVM.maxThreadCount", metricInfo("Peak number of live JVM threads. Unexpected growth signals thread leaks.", "counter", undefined, "stable")],
  ["JVM.totalCpuTimeMs", metricInfo("Cumulative CPU time across all JVM threads. Compared with wall-clock time, reveals CPU saturation.", "milliseconds")],
  ["JVM.totalTimeToSafepointsMs", metricInfo("Time threads spent blocked waiting for JVM safepoints. Above 50 ms indicates contention.", "milliseconds")],
  [
    "JVM.totalTimeAtSafepointsMs",
    metricInfo(
      "Total stop-the-world time the application is paused at safepoints. JVM.totalTimeToSafepointsMs counts only the time to reach a safepoint, not the pause itself.",
      "milliseconds"
    ),
  ],
  ["JVM.totalSafepointCount", metricInfo("Number of safepoints reached. A superset of JVM.GC.collections, which counts only GC pauses.", "counter", undefined, "lower")],
  ["JVM.totalMegabytesAllocated", metricInfo("Total memory allocated across all threads — allocation pressure.", "mebibytes")],
  [
    "JVM.usedHeapMegabytes",
    metricInfo("Peak heap in use. JVM.maxHeapMegabytes is the configured ceiling; JVM.committedHeapMegabytes is the memory reserved from the OS.", "mebibytes"),
  ],
  ["JVM.usedNativeMegabytes", metricInfo("Peak off-heap native memory in use.", "mebibytes")],
  ["JVM.totalDirectByteBuffersMegabytes", metricInfo("Peak memory held by NIO direct byte buffers.", "mebibytes")],
  ["JVM.newThreadsCount", metricInfo("Total threads created. JVM.maxThreadCount captures only the live peak, so it hides thread churn.", "counter", undefined, "lower")],

  // OS
  ["OS.loadAverage", metricInfo("Peak OS load average — machine saturation.", "counter", undefined, "lower")],

  // Memory
  ["MEM.avgRamMegabytes", metricInfo("Average resident RAM of the IDE process during the test.", "mebibytes")],
  ["MEM.avgFileMappingsRamMegabytes", metricInfo("Average RAM consumed by memory-mapped files (JARs, indexes, OS pages).", "mebibytes")],
  ["MEM.avgRamMinusFileMappingsMegabytes", metricInfo("Average resident RAM excluding file mappings — the true heap and off-heap footprint.", "mebibytes")],
  ["JVM.committedHeapMegabytes", metricInfo("Heap memory committed from the OS. Compared with max heap, shows available headroom.", "mebibytes")],

  //gc
  ["freedMemory", metricInfo("Total memory freed by GC during the test, parsed from GCViewer output.", "mebibytes", "https://github.com/chewiebug/GCViewer#readme")],
  ["test#average_awt_delay", metricInfo("Average time to process a single empty AWT event in the queue during the test.", "milliseconds")],
  ["showQuickFixes", metricInfo("Time to show the quick fixes after calling Alt + Enter.", "milliseconds")],
  [
    "attempt.mad.ms",
    metricInfo(
      "MAD (Median Absolute Deviation) of attempt durations in ms. The MAD is a robust statistic, being more resilient to outliers in a data set than the standard deviation.",
      "milliseconds",
      "https://en.m.wikipedia.org/wiki/Median_absolute_deviation"
    ),
  ],
  ...GRADLE_METRICS_NEW_DASHBOARD,

  // Workspace Model
  ["workspaceModel.updates.count", "Total number of changes made to the project model, including modifying entities, changing configuration, and updating dependencies."],
  [
    "workspaceModel.updates.ms",
    metricInfo(
      "Total time to process workspace entity modifications: calling update handlers, collecting changes, initializing bridging operations, and generating snapshots.",
      "milliseconds"
    ),
  ],
  ["workspaceModel.mutableEntityStorage.to.snapshot.ms", metricInfo("Time to create a snapshot of the mutable entity storage.", "milliseconds")],
  ["workspaceModel.mutableEntityStorage.replace.by.source.ms", metricInfo("Time to replace entities in the mutable entity storage by source.", "milliseconds")],
  ["workspaceModel.mutableEntityStorage.add.diff.ms", metricInfo("Time to add differences to the mutable entity storage.", "milliseconds")],
  ["workspaceModel.loading.from.cache.ms", metricInfo("Time to load the workspace model from cache.", "milliseconds")],
  ["workspaceModel.do.save.caches.ms", metricInfo("Time to save the workspace model caches.", "milliseconds")],
  ["workspaceModel.mutableEntityStorage.collect.changes.ms", metricInfo("Time to collect changes in the mutable entity storage.", "milliseconds")],
  ["workspaceModel.mutableEntityStorage.add.entity.ms", metricInfo("Time to add an entity to the mutable entity storage.", "milliseconds")],

  // JVM — already plotted in AdditionalMetrics but missing from metricsDescription
  ["JVM.maxHeapMegabytes", metricInfo("Maximum JVM heap size configured for the IDE (-Xmx; constant). Observed peak is totalHeapUsedMax.", "mebibytes", undefined, "none")],
  ["jps.load.project.to.empty.storage.ms", metricInfo("Time to load a project into empty storage via JPS.", "milliseconds")],
  ["jps.project.serializers.load.ms", metricInfo("Time to load project serializers via JPS.", "milliseconds")],
  ["jps.project.serializers.save.ms", metricInfo("Time to save project serializers via JPS.", "milliseconds")],
  ["jps.facet.change.listener.process.change.events.ms", metricInfo("Time for the JPS facet change listener to process change events.", "milliseconds")],

  ["completion#median_value", metricInfo("Median completion time — typical warm latency.", "milliseconds")],
  ["completion#number#mean_value", metricInfo("Mean suggestion count. May shift if completion ranking or filtering changes.", "counter", undefined, "none")],
  ["completion#standard_deviation", metricInfo("Spread of completion time across invocations. Rising values mean inconsistent latency.", "milliseconds")],
  ["completion#firstElementShown#mean_value", metricInfo("Mean time to first suggestion (span instrument); includes the cold run.", "milliseconds")],
  ["completion_1#firstElementShown", metricInfo("Time to first suggestion on the cold first run (span instrument).", "milliseconds")],
  ["performCompletion_1", metricInfo("Cold-run cost of computing suggestions: resolve, type inference, stub-index lookups.", "milliseconds")],

  ["debugStep_out", metricInfo("Time from invoking Step Out until the debugger pauses again.", "milliseconds")],
  ["debugStep_over", metricInfo("Time from invoking Step Over until the debugger pauses again.", "milliseconds")],

  // GoLand — Data Flow Analysis
  ["go.dfa.general.total.time.ms", metricInfo("Total wall-clock time for data flow analysis across all files.", "milliseconds")],
  ["go.dfa.general.avg.time.ms", metricInfo("Average DFA analysis time per file, including summary loading.", "milliseconds")],
  ["go.dfa.general.avg.without.summary.load.time.ms", metricInfo("Average DFA analysis time per file excluding summary load — isolates pure analysis cost.", "milliseconds")],
  ["go.dfa.general.computed.file.gists.count", metricInfo("Number of file gists (abstracted function summaries) computed. Varies run-to-run.", "counter", undefined, "none")],
  ["go.dfa.general.files.count", metricInfo("Total files analyzed by DFA.", "counter", undefined, "stable")],
  ["go.dfa.general.functions.count", metricInfo("Total functions analyzed by DFA.", "counter", undefined, "stable")],
  [
    "go.dfa.general.resolve.issue.count",
    metricInfo("Deferred-call targets DFA could not resolve. A rise signals a resolve regression; lower is better.", "counter", undefined, "lower"),
  ],

  // Other
  ["Memory | IDE | RESIDENT SIZE (MB) 95th pctl", metricInfo("95th-percentile resident set size — near-peak real memory footprint.", "mebibytes")],
  ["ui.lagging#max_value", metricInfo("Maximum UI-thread lag during the run — the longest stretch the EDT stayed blocked and unresponsive.", "milliseconds")],
])

export interface MetricInfo {
  description: string
  url?: string
  /** The measure unit values of this metric are stored in. Authoritative for value formatting. */
  unit?: MeasureUnit
  betterDirection?: BetterDirection
}

function metricInfo(description: string, unit?: MeasureUnit, url?: string, betterDirection?: BetterDirection): MetricInfo {
  return { description, url, unit, betterDirection }
}

function extractMainPrefix(inputString: string): string {
  const regex = /^(\w+#).*/
  const match = regex.exec(inputString)
  return match ? match[1] : "non-matching"
}

export function getMetricDescription(metric: string | undefined): MetricInfo | null {
  if (metric == undefined) return null
  // A chart series' measureName is a composite of producer parts joined by SERIES_NAME_SEPARATOR
  // (machine/branch/dimension producers prepend non-metric tokens), so the metric is not always the
  // whole string nor at its start. Resolve from the first token matching an exact or "#*" wildcard entry.
  for (const token of metric.split(SERIES_NAME_SEPARATOR)) {
    const metricDescription = metricsDescription.get(token) ?? metricsDescription.get(extractMainPrefix(token) + "*")
    if (metricDescription != undefined) {
      return typeof metricDescription == "string" ? { description: metricDescription } : metricDescription
    }
  }
  return null
}

// The measure unit declared for `metric`, or undefined when the metric carries no declared unit.
// Consulted first by resolveMeasureUnit, so a declared unit overrides every name-based heuristic.
export function getMeasureUnit(metric: string | undefined): MeasureUnit | undefined {
  return getMetricDescription(metric)?.unit ?? undefined
}
