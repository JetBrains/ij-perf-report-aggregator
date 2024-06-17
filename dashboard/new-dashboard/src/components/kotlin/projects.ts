import { ChartDefinition } from "../charts/DashboardCharts"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { computed, ComputedRef } from "vue"
import kotlinProjects from "../../../resources/projects/kotlin_projects.json"

export const KOTLIN_MAIN_METRICS = [
  "completion#mean_value",
  "completion#firstElementShown#mean_value",
  "localInspections#mean_value",
  "semanticHighlighting#mean_value",
  "findUsages#mean_value",
]

const MEASURES = {
  completionMeasures: [
    { name: "completion#mean_value", label: "completion mean value" },
    { name: "completion#firstElementShown#mean_value", label: "firstElementShown mean value" },
  ],
  highlightingMeasures: [{ name: "semanticHighlighting#mean_value", label: "semantic highlighting mean value" }],
  codeAnalysisMeasures: [{ name: "localInspections#mean_value", label: "Code Analysis mean value" }],
  refactoringMeasures: [
    { name: "performInlineRename#mean_value", label: "PerformInlineRename" },
    { name: "startInlineRename#mean_value", label: "StartInlineRename" },
    { name: "prepareForRename#mean_value", label: "PrepareForRename" },
  ],
  optimizeImportsMeasures: [{ name: "execute_editor_optimizeimports", label: "Optimize imports" }],
  insertCodeMeasures: [{ name: "execute_editor_paste", label: "Insert code" }],
  findUsagesMeasures: [{ name: "findUsages#mean_value", label: "findUsages mean value" }],
  evaluateExpressionMeasures: [{ name: "evaluateExpression#mean_value", label: "evaluate expression mean value" }],
  convertJavaToKotlinProjectsMeasures: [{ name: "convertJavaToKotlin", label: "convert java to kotlin" }],
  navigationToDeclarationMeasures: [
    { name: "localInspections_cold#mean_value", label: "Code Analysis cold cache" },
    { name: "localInspections_hot#mean_value", label: "Code Analysis  hot cache" },
    { name: "execute_editor_gotodeclaration_cold#mean_value", label: "Navigate to declaration cold cache" },
    { name: "execute_editor_gotodeclaration_hot#mean_value", label: "Navigate to declaration hot cache" },
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

  intelliJSources: buildCategory("Intellij Sources", "intellij_sources/"),
  intelliJTyping: buildCategory("IntelliJ with typing", ""),
  kotlinLang: buildCategory("Kotlin lang", "kotlin_lang/"),
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
  sqliter: buildCategory("SQLiter", "SQLiter/"),
  ktor: buildCategory("Ktor", "ktor_before_add_wasm_client/"),
  kotlinCoroutinesQG: buildCategory("Kotlin Coroutines QG", "kotlin_coroutines_qg/"),

  mppNativeAcceptance: buildCategory("Native-acceptance", "kotlin_kmp_native_acceptance/"),
  petClinic: buildCategory("Pet Clinic", "kotlin_petclinic/"),
  arrow: buildCategory("Arrow", "arrow/"),
  kotlinEmptyScript: buildCategory("Empty Script (.kts)", "kotlin_empty_kts/"),
}

export const KOTLIN_PROJECTS = kotlinProjects

export const KOTLIN_PROJECT_CONFIGURATOR = new SimpleMeasureConfigurator("project", null)
KOTLIN_PROJECT_CONFIGURATOR.initData(Object.values(PROJECT_CATEGORIES).flatMap((c) => c.label))

function buildCategory(label: string, prefix: string): ProjectCategory {
  return { label, prefix }
}

function projectsToDefinition(projectsByOS: ProjectsByOS[]): ComputedRef<ChartDefinition[]> {
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
      .filter((chart) => KOTLIN_PROJECT_CONFIGURATOR.selected.value?.some((selected) => chart.labels[0].startsWith(`'${selected}'`)))
  )
}

export const completionProjects = { ...KOTLIN_PROJECTS.linux.completion, ...KOTLIN_PROJECTS.mac.completion }

export const completionCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.completion,
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.completion,
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.mac],
  },
])

/**
 * Highlighting projects are also the projects for local inspections.
 */

export const highlightingProjects = { ...KOTLIN_PROJECTS.linux.highlighting, ...KOTLIN_PROJECTS.mac.highlighting }

export const highlightingCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.highlighting,
    measures: MEASURES.highlightingMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.highlighting,
    measures: MEASURES.highlightingMeasures,
    machines: [MACHINES.mac],
  },
])

export const codeAnalysisCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.highlighting,
    measures: MEASURES.codeAnalysisMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.highlighting,
    measures: MEASURES.codeAnalysisMeasures,
    machines: [MACHINES.mac],
  },
])

export const refactoringCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.refactoringRename,
    measures: MEASURES.refactoringMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.linux.refactoringInsertCode,
    measures: MEASURES.insertCodeMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.linux.optimizeImports,
    measures: MEASURES.optimizeImportsMeasures,
    machines: [MACHINES.linux],
  },
])

export const findUsagesProjects = { ...KOTLIN_PROJECTS.linux.findUsages, ...KOTLIN_PROJECTS.mac.findUsages }

export const findUsagesCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.findUsages,
    measures: MEASURES.findUsagesMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.mac.findUsages,
    measures: MEASURES.findUsagesMeasures,
    machines: [MACHINES.mac],
  },
])

export const evaluateExpressionChars = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.evaluateExpression,
    measures: MEASURES.evaluateExpressionMeasures,
    machines: [MACHINES.linux],
  },
  {
    projects: KOTLIN_PROJECTS.linux.completionInEvaluateExpression,
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.linux],
  },
])

export const convertJavaToKotlinProjectsChars = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.convertJavaToKotlin,
    measures: MEASURES.convertJavaToKotlinProjectsMeasures,
    machines: [MACHINES.linux],
  },
])

export const navigateToDeclarationCharts = projectsToDefinition([
  {
    projects: KOTLIN_PROJECTS.linux.navigationToDeclaration,
    measures: MEASURES.navigationToDeclarationMeasures,
    machines: [MACHINES.linux],
  },
])

const scriptHighlight = { kotlinScript: KOTLIN_PROJECTS.linux.highlighting.kotlinScript }
export const highlightingScriptCharts = projectsToDefinition([
  {
    projects: scriptHighlight,
    measures: MEASURES.highlightingMeasures,
    machines: [MACHINES.linux],
  },
])
export const codeAnalysisScriptCharts = projectsToDefinition([
  {
    projects: scriptHighlight,
    measures: MEASURES.codeAnalysisMeasures,
    machines: [MACHINES.linux],
  },
])

export const scriptCompletionCharts = projectsToDefinition([
  {
    projects: { kotlinScript: KOTLIN_PROJECTS.linux.completion.kotlinScript },
    measures: MEASURES.completionMeasures,
    machines: [MACHINES.linux],
  },
])

export const scriptFindUsagesCharts = projectsToDefinition([
  {
    projects: { kotlinScript: KOTLIN_PROJECTS.linux.findUsages.kotlinScript },
    measures: MEASURES.findUsagesMeasures,
    machines: [MACHINES.linux],
  },
])

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
