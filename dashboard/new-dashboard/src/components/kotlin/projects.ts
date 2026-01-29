import { ChartDefinition } from "../charts/DashboardCharts"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { computed, ComputedRef } from "vue"
import kotlinProjects from "../../../resources/projects/kotlin_projects.json"

export const KOTLIN_MAIN_METRICS = ["completion#mean_value", "completion#firstElementShown#mean_value", "localInspections#mean_value", "findUsages#mean_value"]

const MEASURES = {
  completionMeasures: [
    { name: "completion#mean_value", label: "completion mean value" },
    { name: "completion#firstElementShown#mean_value", label: "firstElementShown mean value" },
  ],
  codeAnalysisMeasures: [{ name: "localInspections#mean_value", label: "Code Analysis mean value" }],
  refactoringMeasures: [
    { name: "performInlineRename#mean_value", label: "PerformInlineRename" },
    { name: "startInlineRename#mean_value", label: "StartInlineRename" },
    { name: "prepareForRename#mean_value", label: "PrepareForRename" },
    { name: "fus_refactoring_usages_searched", label: "FindUsagesForRename" },
  ],
  codeTypingMeasures: [{ name: "codeTyping#mean_value", label: "Code Typing mean value" }],
  moveFilesMeasure: [
    { name: "moveFiles#mean_value", label: "Move files" },
    { name: "moveFiles_back#mean_value", label: "Move files back" },
  ],
  moveDeclarationsMeasure: [
    { name: "moveDeclarations#mean_value", label: "Move declarations" },
    { name: "moveDeclarations_back#mean_value", label: "Move declarations back" },
  ],
  optimizeImportsMeasures: [{ name: "execute_editor_optimizeimports", label: "Optimize imports" }],
  insertCodeMeasures: [{ name: "execute_editor_paste", label: "Insert code" }],
  findUsagesMeasures: [{ name: "findUsages#mean_value", label: "findUsages mean value" }],
  findUsagesFirstUsageMeasures: [{ name: "findUsages_firstUsage#mean_value", label: "findUsages first usage mean value" }],
  evaluateExpressionMeasures: [{ name: "evaluateExpression#mean_value", label: "evaluate expression mean value" }],
  convertJavaToKotlinProjectsMeasures: [{ name: "convertJavaToKotlin", label: "convert java to kotlin" }],
  navigationToDeclarationMeasures: [
    { name: "localInspections_cold#mean_value", label: "Code Analysis cold cache" },
    { name: "localInspections_hot#mean_value", label: "Code Analysis hot cache" },
    { name: "execute_editor_gotodeclaration_cold#mean_value", label: "Navigate to declaration cold cache" },
    { name: "execute_editor_gotodeclaration_hot#mean_value", label: "Navigate to declaration hot cache" },
    { name: "freedMemoryByGC", label: "Freed memory by GC" },
  ],
  sequenceHighlightingMeasures: [
    { name: "localInspections_cold#mean_value", label: "Code Analysis cold cache" },
    { name: "localInspections_hot#mean_value", label: "Code Analysis hot cache" },
    { name: "freedMemoryByGC", label: "Freed memory by GC" },
  ],
  findUsagesScenarioMeasures: [
    { name: "findUsagesInBackground_firstUsage_1", label: "First usage found first iteration" },
    { name: "findUsagesInBackground_firstUsage#mean_value", label: "First usage found mean value" },
    { name: "FindUsagesTotal#mean_value", label: "Total find usages time mean value" },
  ],
  completionCausingModificationMeasures: [{ name: "total_test_step#mean_value", label: "Total test time" }],
  renameAndCompletionMeasures: [{ name: "total_test_step#mean_value", label: "Total test time" }],
  errorCodeModificationMeasures: [{ name: "typing_and_highlighting#mean_value", label: "Make/fix error and wait for highlighting" }],
  goToImplementationScenarioMeasures: (tag: string) => [
    { name: "execute_editor_gotoimplementation_1", label: `Go to implementation first iteration ${tag}` },
    { name: "execute_editor_gotoimplementation#mean_value", label: `Go to implementation mean value ${tag}` },
  ],
  intelliJFindUsagesAndHighlightingMeasures: (tag: string) => [
    { name: "highlighting_IdeResourcesUtil.kt_1", label: `Code Analysis IdeResourcesUtil first iteration ${tag}` },
    { name: "highlighting_IdeResourcesUtil.kt#mean_value", label: `Code Analysis IdeResourcesUtil mean value ${tag}` },
    { name: "freedMemoryByGC", label: `Freed memory by GC ${tag}` },
  ],
  kotlinFindUsagesAndHighlightingMeasures: (tag: string) => [
    { name: "highlighting_RawFirBuilder.kt_1", label: `Code Analysis RawFirBuilder first iteration ${tag}` },
    { name: "highlighting_RawFirBuilder.kt#mean_value", label: `Code Analysis RawFirBuilder mean value ${tag}` },
    { name: "freedMemoryByGC", label: `Freed memory by GC ${tag}` },
  ],
  intelliJFindUsagesAndGoToImplementationMeasures: (tag: string) => [
    { name: "highlighting_ReadActionCacheImpl.kt_1", label: `Code Analysis ReadActionCacheImpl first iteration ${tag}` },
    { name: "highlighting_ReadActionCacheImpl.kt#mean_value", label: `Code Analysis ReadActionCacheImpl mean value ${tag}` },
    { name: "freedMemoryByGC", label: `Freed memory by GC ${tag}` },
  ],
  kotlinFindUsagesAndGoToImplementationMeasures: (tag: string) => [
    { name: "highlighting_KaptExtension.kt_1", label: `Code Analysis KaptExtension first iteration ${tag}` },
    { name: "highlighting_KaptExtension.kt#mean_value", label: `Code Analysis KaptExtension mean value ${tag}` },
    { name: "freedMemoryByGC", label: `Freed memory by GC ${tag}` },
  ],
  deleteAllImportsMeasures: [
    { name: "localInspections#mean_value", label: "Code Analysis" },
    { name: "completion#mean_value", label: "Completion" },
    { name: "completion#firstElementShown#mean_value", label: "First element shown" },
    { name: "freedMemoryByGC", label: "Freed memory by GC" },
  ],
}

export const MACHINES = {
  linux: "linux-blade-hetzner",
  mac: "Mac Mini M2 Pro (10 vCPU, 32 GB)",
}

export const PROJECT_CATEGORIES: Record<string, ProjectCategory> = {
  kotlinEmpty: buildCategory("Kotlin empty", "kotlin_empty/"),
  intelliJ: buildCategory("IntelliJ", "intellij_commit/"),
  // Same intelliJ. Need to avoid lot of lines on chart
  intelliJ2: buildCategory("IntelliJ suite 2", ""),
  intelliJ3: buildCategory("IntelliJ suite 3", ""),

  intelliJSources: buildCategory("IntelliJ Sources", "intellij_sources/"),
  intelliJTyping: buildCategory("IntelliJ with typing", ""),
  kotlinLang: buildCategory("Kotlin lang", "kotlin_lang/"),
  kotlinLang_slow: buildCategory("Kotlin lang (slow)", "kotlin_lang/"),
  // Same kotlinLang. Need to avoid lot of lines on chart
  kotlinScript: buildCategory("Kotlin script", ""),

  kotlinLanguageServer: buildCategory("Kotlin language server", "kotlin_language_server/"),
  tbe: buildCategory("Toolbox Enterprise (TBE)", "toolbox_enterprise/"),
  tbeCaseWithAssert: buildCategory("Toolbox Enterprise (TBE) different length", ""),

  ktorSamples: buildCategory("Ktor samples", "ktor_samples"),
  androidCanaryLeak: buildCategory("Android canary leak", "leak-canary-android/"),
  anki: buildCategory("Android anki project", "anki-android/"),
  removedImports: buildCategory("Files with removed imports", ""),
  springFramework: buildCategory("Spring framework", "spring-framework/"),
  rustPlugin: buildCategory("Rust plugin", "rust_commit/"),
  syntheticFiles: buildCategory("Synthetic files", "kotlin_synthetic_files/"),
  sqliter: buildCategory("SQLiter", "SQLiter/"),
  sqliter_slow: buildCategory("SQLiter (slow)", "SQLiter/"),
  ktor: buildCategory("Ktor", "ktor_before_add_wasm_client/"),
  kotlinCoroutinesQG: buildCategory("Kotlin Coroutines QG", "kotlin_coroutines_qg/"),
  kotlinCoroutinesQG_slow: buildCategory("Kotlin Coroutines QG (slow)", "kotlin_coroutines_qg/"),

  mppNativeAcceptance: buildCategory("Native-acceptance", "kotlin_kmp_native_acceptance/"),
  petClinic: buildCategory("Pet Clinic", "kotlin_petclinic/"),
  arrow: buildCategory("Arrow", "arrow/"),
  jooq: buildCategory("JOOQ", "jooq-k2/"),
  kotlinEmptyScript: buildCategory("Empty Script (.kts)", "kotlin_empty_kts/"),

  qaRefactorMove: buildCategory("QA Refactor / Move", "refactor-move/"),
}

export const KOTLIN_PROJECTS = kotlinProjects

const codeTypingProjects = KOTLIN_PROJECTS.linux.codeTyping as Record<string, string[]>
const findUsagesFirstUsageProjects = KOTLIN_PROJECTS.linux.findUsagesFirstUsage as Record<string, string[]>
const evaluateExpressionProjects = KOTLIN_PROJECTS.linux.evaluateExpression as Record<string, string[]>
const convertJavaToKotlinProjects = KOTLIN_PROJECTS.linux.convertJavaToKotlin as Record<string, string[]>
const navigationToDeclarationProjects = KOTLIN_PROJECTS.linux.navigationToDeclaration as Record<string, string[]>
const sequenceHighlightingProjects = KOTLIN_PROJECTS.linux.sequenceHighlighting as Record<string, string[]>

const scriptHighlight = { kotlinScript: (KOTLIN_PROJECTS.linux.highlighting as Record<string, string[]>)["kotlinScript"] }
const scriptCompletion = { kotlinScript: (KOTLIN_PROJECTS.linux.completion as Record<string, string[]>)["kotlinScript"] }
const scriptFindUsages = { kotlinScript: (KOTLIN_PROJECTS.linux.findUsages as Record<string, string[]>)["kotlinScript"] }

function buildCategory(label: string, prefix: string): ProjectCategory {
  return { label, prefix }
}

function projectsToDefinition(projectsByOS: ProjectsByOS[], configurator: SimpleMeasureConfigurator): ComputedRef<ChartDefinition[]> {
  return computed(() =>
    projectsByOS
      .flatMap(({ projects, measures, machines }) =>
        measures.flatMap(({ name, label }) =>
          Object.entries(projects).flatMap(([key, value]) => {
            return {
              labels: [`'${PROJECT_CATEGORIES[key].label}' ${label}`],
              measures: [name],
              projects: value,
              machines,
            }
          })
        )
      )
      .filter((chart) => configurator.selected.value?.some((selected) => chart.labels[0].startsWith(`'${selected}'`)))
  )
}

export const completionProjects = { ...KOTLIN_PROJECTS.linux.completion, ...KOTLIN_PROJECTS.mac.completion } as Record<string, string[]>

export const completionChartsDescription =
  "A Completion test invokes completion (CTRL + Space) explicitly at a specific position in a file. Depending on the test, the completion can be invoked with or without already typed text."

/**
 * Highlighting projects are also the projects for local inspections.
 */
export const highlightingProjects = { ...KOTLIN_PROJECTS.linux.highlighting, ...KOTLIN_PROJECTS.mac.highlighting } as Record<string, string[]>

export const codeAnalysisChartsDescription = "A Code Analysis test measures the full highlighting of a single file, including inspections."

export const refactoringChartsDescription = undefined

export const codeTypingChartsDescription =
  "A Code Typing test types a piece of code with medium complexity from beginning to end. It invokes completion and waits for code analysis at predetermined points."

export const findUsagesProjects = { ...KOTLIN_PROJECTS.linux.findUsages, ...KOTLIN_PROJECTS.mac.findUsages }

export const findUsagesChartsDescription =
  "A Find Usages test invokes Find Usages (Option + F7) on a specific declaration in a file. It waits until all usages have been found and compares the number of usages with the" +
  " expected value.\n\nThe test does not wait for the file's code analysis to finish but rather runs Find Usages in parallel with it. This matches user behavior more closely."

export const evaluateExpressionChartsDescription = undefined

export const convertJavaToKotlinProjectsChartsDescription =
  "A Java to Kotlin (J2K) test converts a Java file to a Kotlin file. It measures the full time it takes to complete the conversion."

export const scriptChartsDescription = "This category contains various kinds of performance tests that are performed on script files (.kts)."

// Deprecated: Only used for scenario charts. New dashboards should create their own local configurator instance.
// Initialized immediately with all projects so scenario charts display all data by default.
const KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED = new SimpleMeasureConfigurator("project", null)
KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED.initData(Object.values(PROJECT_CATEGORIES).flatMap((c) => c.label))

// Scenario charts for USER_SCENARIOS (using deprecated configurator for now)
const sequenceHighlightingScenarioCharts = projectsToDefinition(
  [{ projects: KOTLIN_PROJECTS.linux.sequenceHighlighting, measures: MEASURES.sequenceHighlightingMeasures, machines: [MACHINES.linux] }],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

const navigateToDeclarationScenarioCharts = projectsToDefinition(
  [{ projects: KOTLIN_PROJECTS.linux.navigationToDeclaration, measures: MEASURES.navigationToDeclarationMeasures, machines: [MACHINES.linux] }],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

const deleteAllImportsScenarioCharts = projectsToDefinition(
  [{ projects: KOTLIN_PROJECTS.linux.deleteAllImports, measures: MEASURES.deleteAllImportsMeasures, machines: [MACHINES.linux] }],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

const completionCausingModificationScenarioCharts = projectsToDefinition(
  [{ projects: KOTLIN_PROJECTS.linux.completionCausingModification, measures: MEASURES.completionCausingModificationMeasures, machines: [MACHINES.linux] }],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

const errorCodeModificationScenarioCharts = projectsToDefinition(
  [{ projects: KOTLIN_PROJECTS.linux.errorCodeModification, measures: MEASURES.errorCodeModificationMeasures, machines: [MACHINES.linux] }],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

const renameAndCompletionScenarioCharts = projectsToDefinition(
  [{ projects: KOTLIN_PROJECTS.linux.renameAndCompletion, measures: MEASURES.renameAndCompletionMeasures, machines: [MACHINES.linux] }],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

const findUsagesAndHighlightingScenarioCharts = projectsToDefinition(
  [
    {
      projects: KOTLIN_PROJECTS.linux.intelliJFindUsagesAndHighlighting,
      measures: [...MEASURES.intelliJFindUsagesAndHighlightingMeasures("with usages"), ...MEASURES.findUsagesScenarioMeasures],
      machines: [MACHINES.linux],
    },
    {
      projects: KOTLIN_PROJECTS.linux.intelliJFindUsagesAndHighlightingNoFindUsages,
      measures: MEASURES.intelliJFindUsagesAndHighlightingMeasures("without usages"),
      machines: [MACHINES.linux],
    },
    {
      projects: KOTLIN_PROJECTS.linux.kotlinFindUsagesAndHighlighting,
      measures: [...MEASURES.kotlinFindUsagesAndHighlightingMeasures("with usages"), ...MEASURES.findUsagesScenarioMeasures],
      machines: [MACHINES.linux],
    },
    {
      projects: KOTLIN_PROJECTS.linux.kotlinFindUsagesAndHighlightingKotlinNoFindUsages,
      measures: MEASURES.kotlinFindUsagesAndHighlightingMeasures("without usages"),
      machines: [MACHINES.linux],
    },
  ],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

const findUsagesAndGoToImplementationScenarioCharts = projectsToDefinition(
  [
    {
      projects: KOTLIN_PROJECTS.linux.intelliJFindUsagesAndGotoImplementation,
      measures: [
        ...MEASURES.intelliJFindUsagesAndGoToImplementationMeasures("with usages"),
        ...MEASURES.findUsagesScenarioMeasures,
        ...MEASURES.goToImplementationScenarioMeasures("with usages"),
      ],
      machines: [MACHINES.linux],
    },
    {
      projects: KOTLIN_PROJECTS.linux.intelliJFindUsagesAndGotoImplementationNoUsages,
      measures: [...MEASURES.intelliJFindUsagesAndGoToImplementationMeasures("without usages"), ...MEASURES.goToImplementationScenarioMeasures("without usages")],
      machines: [MACHINES.linux],
    },
    {
      projects: KOTLIN_PROJECTS.linux.kotlinFindUsagesAndGotoImplementation,
      measures: [
        ...MEASURES.kotlinFindUsagesAndGoToImplementationMeasures("with usages"),
        ...MEASURES.findUsagesScenarioMeasures,
        ...MEASURES.goToImplementationScenarioMeasures("with usages"),
      ],
      machines: [MACHINES.linux],
    },
    {
      projects: KOTLIN_PROJECTS.linux.kotlinFindUsagesAndGotoImplementationNoUsages,
      measures: [...MEASURES.kotlinFindUsagesAndGoToImplementationMeasures("without usages"), ...MEASURES.goToImplementationScenarioMeasures("without usages")],
      machines: [MACHINES.linux],
    },
  ],
  KOTLIN_PROJECT_CONFIGURATOR_DEPRECATED
)

export const USER_SCENARIOS: Record<string, ScenarioData> = {
  navigateToDeclaration: { label: "Navigate to declaration (one file per test)", charts: navigateToDeclarationScenarioCharts },
  sequenceNavigateToDeclaration: { label: "Sequence highlighting", charts: sequenceHighlightingScenarioCharts },
  deleteAllImports: { label: "Delete all imports", charts: deleteAllImportsScenarioCharts },
  findUsagesAndHighlighting: { label: "Find usages and Highlighting", charts: findUsagesAndHighlightingScenarioCharts },
  findUsagesAndGoToImplementation: { label: "Find usages and Goto implementation", charts: findUsagesAndGoToImplementationScenarioCharts },
  completionCausingModification: { label: "Completion causing modification", charts: completionCausingModificationScenarioCharts },
  errorCodeModification: { label: "Make/fix error and highlighting", charts: errorCodeModificationScenarioCharts },
  renameAndCompletion: { label: "Rename symbol (file A) and use this name in completion (file B)", charts: renameAndCompletionScenarioCharts },
}

export function createKotlinCharts(configurator: SimpleMeasureConfigurator) {
  return {
    completionCharts: projectsToDefinition(
      [
        { projects: KOTLIN_PROJECTS.linux.completion as Record<string, string[]>, measures: MEASURES.completionMeasures, machines: [MACHINES.linux] },
        { projects: KOTLIN_PROJECTS.mac.completion as Record<string, string[]>, measures: MEASURES.completionMeasures, machines: [MACHINES.mac] },
        { projects: navigationToDeclarationProjects, measures: MEASURES.navigationToDeclarationMeasures, machines: [MACHINES.linux] },
      ],
      configurator
    ),
    codeAnalysisCharts: projectsToDefinition(
      [
        { projects: KOTLIN_PROJECTS.linux.highlighting as Record<string, string[]>, measures: MEASURES.codeAnalysisMeasures, machines: [MACHINES.linux] },
        { projects: KOTLIN_PROJECTS.mac.highlighting as Record<string, string[]>, measures: MEASURES.codeAnalysisMeasures, machines: [MACHINES.mac] },
        { projects: sequenceHighlightingProjects, measures: MEASURES.sequenceHighlightingMeasures, machines: [MACHINES.linux] },
      ],
      configurator
    ),
    refactoringCharts: projectsToDefinition(
      [
        { projects: KOTLIN_PROJECTS.linux.refactoringRename as Record<string, string[]>, measures: MEASURES.refactoringMeasures, machines: [MACHINES.linux] },
        { projects: KOTLIN_PROJECTS.linux.refactoringInsertCode as Record<string, string[]>, measures: MEASURES.insertCodeMeasures, machines: [MACHINES.linux] },
        { projects: KOTLIN_PROJECTS.linux.optimizeImports as Record<string, string[]>, measures: MEASURES.optimizeImportsMeasures, machines: [MACHINES.linux] },
        { projects: KOTLIN_PROJECTS.linux.moveFiles as Record<string, string[]>, measures: MEASURES.moveFilesMeasure, machines: [MACHINES.linux] },
        { projects: KOTLIN_PROJECTS.linux.moveDeclarations as Record<string, string[]>, measures: MEASURES.moveDeclarationsMeasure, machines: [MACHINES.linux] },
      ],
      configurator
    ),
    codeTypingCharts: projectsToDefinition([{ projects: codeTypingProjects, measures: MEASURES.codeTypingMeasures, machines: [MACHINES.linux] }], configurator),
    findUsagesCharts: projectsToDefinition(
      [
        { projects: KOTLIN_PROJECTS.linux.findUsages as Record<string, string[]>, measures: MEASURES.findUsagesMeasures, machines: [MACHINES.linux] },
        { projects: KOTLIN_PROJECTS.mac.findUsages as Record<string, string[]>, measures: MEASURES.findUsagesMeasures, machines: [MACHINES.mac] },
        { projects: findUsagesFirstUsageProjects, measures: MEASURES.findUsagesFirstUsageMeasures, machines: [MACHINES.linux] },
      ],
      configurator
    ),
    evaluateExpressionCharts: projectsToDefinition(
      [{ projects: evaluateExpressionProjects, measures: MEASURES.evaluateExpressionMeasures, machines: [MACHINES.linux] }],
      configurator
    ),
    convertJavaToKotlinProjectsCharts: projectsToDefinition(
      [{ projects: convertJavaToKotlinProjects, measures: MEASURES.convertJavaToKotlinProjectsMeasures, machines: [MACHINES.linux] }],
      configurator
    ),
    codeAnalysisScriptCharts: projectsToDefinition([{ projects: scriptHighlight, measures: MEASURES.codeAnalysisMeasures, machines: [MACHINES.linux] }], configurator),
    scriptCompletionCharts: projectsToDefinition([{ projects: scriptCompletion, measures: MEASURES.completionMeasures, machines: [MACHINES.linux] }], configurator),
    scriptFindUsagesCharts: projectsToDefinition([{ projects: scriptFindUsages, measures: MEASURES.findUsagesMeasures, machines: [MACHINES.linux] }], configurator),
  }
}

export const KOTLIN_SCENARIO_CONFIGURATOR = new SimpleMeasureConfigurator("scenario", null)

type Projects = Record<string, string[]>

interface ProjectsByOS {
  projects: Projects
  measures: Measure[]
  machines: string[]
}

/**
 * A project category is a project name prefix such as "intellij_commit/" and "kotlin_lang/" with an associated, human-readable label.
 */
interface ProjectCategory {
  label: string
  prefix: string
}

interface Measure {
  name: string
  label: string
}

interface ScenarioData {
  label: string
  charts: ComputedRef<ChartDefinition[]>
}
