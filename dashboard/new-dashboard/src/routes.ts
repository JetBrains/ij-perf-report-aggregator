/* eslint-disable @typescript-eslint/prefer-literal-enum-member */

import { RouteRecordRaw } from "vue-router"
import { ParentRouteRecord, TypedRouteRecord } from "./components/common/route"
import { KOTLIN_MAIN_METRICS } from "./components/kotlin/projects"
import type { PerformanceTestsProps } from "./components/common/PerformanceTests.props"
import type { PerformanceUnitTestsProps } from "./components/common/PerformanceUnitTests.vue"

const COMPONENTS = {
  perfTests: () => import("./components/common/PerformanceTests.vue"),
  compareBuilds: () => import("./components/common/compare/CompareBuilds.vue"),
  startupDashboard: () => import("./components/common/StartupMetricsDashboard.vue"),
  compareBranches: () => import("./components/common/compare/CompareBranches.vue"),
  compareModes: () => import("./components/common/compare/CompareModes.vue"),
}

const MACHINES = {
  HETZNER: "linux-blade-hetzner",
  AWS_LINUX: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
}

enum ROUTE_PREFIX {
  Startup = "/ij",
  IntelliJ = "/intellij",
  IntelliJBuildTools = "/intellij/buildTools",
  IntelliJSharedIndexes = "/intellij/sharedIndexes",
  IntelliJKotlinK2Performance = "/intellij/kotlinK2Performance",
  IntelliJPackageChecker = "/intellij/packageChecker",
  PhpStorm = "/phpstorm",
  GoLand = "/goland",
  RubyMine = "/rubymine",
  Kotlin = "/kotlin",
  KotlinBuildTools = "/kotlinBuildTools",
  KotlinMemory = Kotlin + "/memory",
  Rust = "/rust",
  Scala = "/scala",
  JBR = "/jbr",
  Fleet = "/fleet",
  PyCharm = "/pycharm",
  WebStorm = "/webstorm",
  Bazel = "/bazel",
  Qodana = "/qodana",
  Clion = "/clion",
  Vcs = IntelliJ + "/vcs",
  EmbeddingSearch = IntelliJ + "/embeddingSearch",
  PerfUnit = "/perfUnit",
  IJent = "/ijent",
  ML = "/ml",
  DataGrip = "/datagrip",
  AIA = "/aia",
  KMT = "/kmt",
  Diogen = "/diogen",
  Toolbox = "/toolbox",
  LSP = "/lsp",
  KotlinNotebooks = "/kotlinNotebooks",
}

const TEST_ROUTE = "tests"
const DEV_TEST_ROUTE = "testsDev"
const DASHBOARD_ROUTE = "dashboard"
const STARTUP_ROUTE = "startup"
const PRODUCT_METRICS_ROUTE = "product-metrics"
const COMPARE_ROUTE = "compare"
const COMPARE_BRANCHES_ROUTE = "compareBranches"
const COMPARE_MODES_ROUTE = "compareModes"

enum ROUTES {
  StartupPulse = `${ROUTE_PREFIX.Startup}/pulse`,
  StartupPulseInstaller = `${ROUTE_PREFIX.Startup}/pulseInstaller`,
  StartupProgress = `${ROUTE_PREFIX.Startup}/progressOverTime`,
  StartupModuleLoading = `${ROUTE_PREFIX.Startup}/moduleLoading`,
  StartupGcAndMemory = `${ROUTE_PREFIX.Startup}/gcAndMemory`,
  StartupExplore = `${ROUTE_PREFIX.Startup}/explore`,
  StartupExploreInstaller = `${ROUTE_PREFIX.Startup}/exploreInstaller`,
  IntelliJStartupDashboard = `${ROUTE_PREFIX.IntelliJ}/${STARTUP_ROUTE}`,
  IntelliJProductMetricsDashboard = `${ROUTE_PREFIX.IntelliJ}/${PRODUCT_METRICS_ROUTE}`,
  IntelliJIndexingDashboard = `${ROUTE_PREFIX.IntelliJ}/indexingDashboard`,
  IntelliJJavaDashboard = `${ROUTE_PREFIX.IntelliJ}/javaDashboard`,
  IntelliJKotlinDashboard = `${ROUTE_PREFIX.IntelliJ}/kotlinDashboard`,
  IntelliJUltimateDashboard = `${ROUTE_PREFIX.IntelliJ}/ultimateDashboard`,
  IntelliJUIDashboard = `${ROUTE_PREFIX.IntelliJ}/uiDashboard`,
  IntelliJSearchEverywhereExDashboard = `${ROUTE_PREFIX.IntelliJ}/searchEverywhereExDashboard`,
  IntelliJEmbeddingSearchDashboard = `${ROUTE_PREFIX.EmbeddingSearch}/dashboard`,
  IntelliJK2Dashboard = `${ROUTE_PREFIX.IntelliJKotlinK2Performance}/${DASHBOARD_ROUTE}`,
  IntelliJDevTests = `${ROUTE_PREFIX.IntelliJ}/${DEV_TEST_ROUTE}`,
  IntelliJCompare = `${ROUTE_PREFIX.IntelliJ}/${COMPARE_ROUTE}`,
  IntelliJCompareBranches = `${ROUTE_PREFIX.IntelliJ}/${COMPARE_BRANCHES_ROUTE}`,
  IntelliJCompareModes = `${ROUTE_PREFIX.IntelliJ}/${COMPARE_MODES_ROUTE}`,
  IntelliJGradleDashboardDev = `${ROUTE_PREFIX.IntelliJBuildTools}/gradleDashboardDev`,
  IntelliJMavenDashboardDev = `${ROUTE_PREFIX.IntelliJBuildTools}/mavenDashboardDev`,
  IntelliJMavenImportersConfiguratorsDashboardDev = `${ROUTE_PREFIX.IntelliJBuildTools}/mavenImportersConfiguratorsDashboardDev`,
  IntelliJJpsDashboardDev = `${ROUTE_PREFIX.IntelliJBuildTools}/jpsDashboardDev`,
  IntelliJBuildTests = `${ROUTE_PREFIX.IntelliJBuildTools}/${TEST_ROUTE}`,
  IntelliJBuildTestsDev = `${ROUTE_PREFIX.IntelliJBuildTools}/${DEV_TEST_ROUTE}`,
  IntelliJGradleBenchmarks = `${ROUTE_PREFIX.IntelliJBuildTools}/gradleBenchmarks`,
  IntelliJSharedIndicesDashboard = `${ROUTE_PREFIX.IntelliJSharedIndexes}/${DASHBOARD_ROUTE}`,
  IntelliJSharedIndicesTests = `${ROUTE_PREFIX.IntelliJSharedIndexes}/${TEST_ROUTE}`,
  IntelliJPackageCheckerDashboard = `${ROUTE_PREFIX.IntelliJPackageChecker}/${DASHBOARD_ROUTE}`,
  IntelliJPackageCheckerTests = `${ROUTE_PREFIX.IntelliJPackageChecker}/${TEST_ROUTE}`,
  PhpStormProductMetricsDashboard = `${ROUTE_PREFIX.PhpStorm}/${PRODUCT_METRICS_ROUTE}`,
  PhpStormLLMDashboard = `${ROUTE_PREFIX.PhpStorm}/llmDashboard`,
  PhpStormIndexingDashboard = `${ROUTE_PREFIX.PhpStorm}/indexingDashboard`,
  PhpStormInspectionsDashboard = `${ROUTE_PREFIX.PhpStorm}/inspectionsDashboard`,
  PhpStormCodeEditingDashboard = `${ROUTE_PREFIX.PhpStorm}/codeEditingDashboard`,
  PhpStormUnitTestsDashboard = `${ROUTE_PREFIX.PhpStorm}/unitTestsDashboard`,
  PhpStormStartupDashboard = `${ROUTE_PREFIX.PhpStorm}/${STARTUP_ROUTE}`,
  PhpStormWithPluginsDashboard = `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests = `${ROUTE_PREFIX.PhpStorm}/${TEST_ROUTE}`,
  PhpStormDevTests = `${ROUTE_PREFIX.PhpStorm}/${DEV_TEST_ROUTE}`,
  PhpStormWithPluginsTests = `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,
  PhpStormCompareBranches = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_BRANCHES_ROUTE}`,
  PhpStormCompareModes = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_MODES_ROUTE}`,
  KotlinDashboard = `${ROUTE_PREFIX.Kotlin}/${DASHBOARD_ROUTE}`,
  KotlinDashboardDev = `${ROUTE_PREFIX.Kotlin}/${DASHBOARD_ROUTE}Dev`,
  KotlinUserScenariosDashboardDev = `${ROUTE_PREFIX.Kotlin}/Scenarios${DASHBOARD_ROUTE}Dev`,
  KotlinCodeAnalysisDev = `${ROUTE_PREFIX.Kotlin}/codeAnalysisDev `,
  KotlinTests = `${ROUTE_PREFIX.Kotlin}/${TEST_ROUTE}`,
  KotlinTestsDev = `${ROUTE_PREFIX.Kotlin}/${DEV_TEST_ROUTE}`,
  KotlinCompletionDev = `${ROUTE_PREFIX.Kotlin}/completionDev`,
  KotlinFindUsagesDev = `${ROUTE_PREFIX.Kotlin}/findUsagesDev`,
  KotlinRefactoringDev = `${ROUTE_PREFIX.Kotlin}/refactoringDev`,
  KotlinDebuggerDev = `${ROUTE_PREFIX.Kotlin}/debuggerDev`,
  KotlinScriptDev = `${ROUTE_PREFIX.Kotlin}/scriptDev`,
  KotlinCompare = `${ROUTE_PREFIX.Kotlin}/${COMPARE_ROUTE}`,
  KotlinMemoryDashboard = `${ROUTE_PREFIX.KotlinMemory}/dashboard`,
  KotlinMemoryDashboardDev = `${ROUTE_PREFIX.KotlinMemory}/dashboardDev`,
  KotlinCompareBranches = `${ROUTE_PREFIX.Kotlin}/${COMPARE_BRANCHES_ROUTE}`,
  KotlinCompareBranchesDev = `${ROUTE_PREFIX.Kotlin}/${COMPARE_BRANCHES_ROUTE}Dev`,
  GoLandStartupDashboard = `${ROUTE_PREFIX.GoLand}/${STARTUP_ROUTE}`,
  GoLandProductMetricsDashboard = `${ROUTE_PREFIX.GoLand}/${PRODUCT_METRICS_ROUTE}`,
  GoLandIndexingDashboard = `${ROUTE_PREFIX.GoLand}/indexingDashboard`,
  GoLandCompletionDashboard = `${ROUTE_PREFIX.GoLand}/completionDashboard`,
  GoLandHighlightingDashboard = `${ROUTE_PREFIX.GoLand}/highlightingDashboard`,
  GoLandInspectionDashboard = `${ROUTE_PREFIX.GoLand}/inspectionsDashboard`,
  GoLandDebuggerDashboard = `${ROUTE_PREFIX.GoLand}/debuggerDashboard`,
  GoLandFindUsagesDashboard = `${ROUTE_PREFIX.GoLand}/findUsagesDashboard`,
  GoLandDFADashboard = `${ROUTE_PREFIX.GoLand}/dfaDashboard`,
  GoLandTests = `${ROUTE_PREFIX.GoLand}/${TEST_ROUTE}Dev`,
  GoLandCompare = `${ROUTE_PREFIX.GoLand}/${COMPARE_ROUTE}`,
  GoLandCompareBranches = `${ROUTE_PREFIX.GoLand}/${COMPARE_BRANCHES_ROUTE}`,
  GoLandCompareModes = `${ROUTE_PREFIX.GoLand}/${COMPARE_MODES_ROUTE}`,
  PyCharmStartupDashboard = `${ROUTE_PREFIX.PyCharm}/${STARTUP_ROUTE}`,
  PyCharmProductMetricsDashboard = `${ROUTE_PREFIX.PyCharm}/${PRODUCT_METRICS_ROUTE}`,
  PyCharmDashboard = `${ROUTE_PREFIX.PyCharm}/${DASHBOARD_ROUTE}Dev`,
  PyCharmExternalTypeProviders = `${ROUTE_PREFIX.PyCharm}/externalTypeProviders`,
  PyCharmExternalTypeProvidersUnitPerfTests = `${ROUTE_PREFIX.PyCharm}/externalTypeProvidersUnitPerfTests`,
  PyCharmPerfUnitTests = `${ROUTE_PREFIX.PyCharm}/perfUnitTests`,
  PyCharmOldDashboard = `${ROUTE_PREFIX.PyCharm}/${DASHBOARD_ROUTE}`,
  PyCharmTests = `${ROUTE_PREFIX.PyCharm}/${TEST_ROUTE}`,
  PyCharmDevTests = `${ROUTE_PREFIX.PyCharm}/${DEV_TEST_ROUTE}`,
  PyCharmCompare = `${ROUTE_PREFIX.PyCharm}/${COMPARE_ROUTE}`,
  PyCharmCompareBranches = `${ROUTE_PREFIX.PyCharm}/${COMPARE_BRANCHES_ROUTE}`,
  WebStormStartupDashboard = `${ROUTE_PREFIX.WebStorm}/${STARTUP_ROUTE}`,
  WebStormProductMetricsDashboard = `${ROUTE_PREFIX.WebStorm}/${PRODUCT_METRICS_ROUTE}`,
  WebStormDashboardBuiltInVsNEXT = `${ROUTE_PREFIX.WebStorm}/dashboardBuiltInVsNext`,
  WebStormDashboardDelicateProjects = `${ROUTE_PREFIX.WebStorm}/dashboardDelicateProjects`,
  WebStormTests = `${ROUTE_PREFIX.WebStorm}/${TEST_ROUTE}Dev`,
  WebStormCompare = `${ROUTE_PREFIX.WebStorm}/${COMPARE_ROUTE}`,
  WebStormCompareBranches = `${ROUTE_PREFIX.WebStorm}/${COMPARE_BRANCHES_ROUTE}`,
  RubyStartupDashboard = `${ROUTE_PREFIX.RubyMine}/${STARTUP_ROUTE}`,
  RubyMineProductMetricsDashboard = `${ROUTE_PREFIX.RubyMine}/${PRODUCT_METRICS_ROUTE}Dev`,
  RubyMineDashboard = `${ROUTE_PREFIX.RubyMine}/${DASHBOARD_ROUTE}Dev`,
  RubyMineIndexingDashBoard = `${ROUTE_PREFIX.RubyMine}/indexingDashboardDev`,
  RubyMineInspectionsDashBoard = `${ROUTE_PREFIX.RubyMine}/inspectionsDashboardDev`,
  RubyMineTestsDev = `${ROUTE_PREFIX.RubyMine}/${DEV_TEST_ROUTE}`,
  RubyMineCompare = `${ROUTE_PREFIX.RubyMine}/${COMPARE_ROUTE}`,
  RubyMineCompareBranches = `${ROUTE_PREFIX.RubyMine}/${COMPARE_BRANCHES_ROUTE}`,
  RubyMineCompareModes = `${ROUTE_PREFIX.RubyMine}/${COMPARE_MODES_ROUTE}`,
  RustRoverDashboardDev = `${ROUTE_PREFIX.Rust}/rustPluginDashboardDev`,
  RustRoverProductMetricsDashboardDev = `${ROUTE_PREFIX.Rust}/${PRODUCT_METRICS_ROUTE}Dev`,
  RustRoverFirstStartupDashboardDev = `${ROUTE_PREFIX.Rust}/rustRoverFirstStartupDashboardDev`,
  RustRoverDebuggerDashboardDev = `${ROUTE_PREFIX.Rust}/debuggerDashboardDev`,
  RustRoverRefactoringDashboardDev = `${ROUTE_PREFIX.Rust}/refactoringDashboardDev`,
  RustRoverUnitTestsDashboardDev = `${ROUTE_PREFIX.Rust}/unitTestsDashboardDev`,
  RustTestsDev = `${ROUTE_PREFIX.Rust}/${TEST_ROUTE}Dev`,
  RustCompareBranchesDev = `${ROUTE_PREFIX.Rust}/${COMPARE_BRANCHES_ROUTE}Dev`,
  ScalaTests = `${ROUTE_PREFIX.Scala}/${TEST_ROUTE}`,
  ScalaCompare = `${ROUTE_PREFIX.Scala}/${COMPARE_ROUTE}`,
  ScalaCompareBranches = `${ROUTE_PREFIX.Scala}/${COMPARE_BRANCHES_ROUTE}`,
  JBRTests = `${ROUTE_PREFIX.JBR}/${TEST_ROUTE}`,
  MapBenchDashboard = `${ROUTE_PREFIX.JBR}/mapbenchDashboard`,
  DaCapoDashboard = `${ROUTE_PREFIX.JBR}/dacapoDashboard`,
  J2DBenchDashboard = `${ROUTE_PREFIX.JBR}/j2dDashboard`,
  JavaDrawDashboard = `${ROUTE_PREFIX.JBR}/javaDrawDashboard`,
  RenderDashboard = `${ROUTE_PREFIX.JBR}/renderDashboard`,
  SPECjbb2015Dashboard = `${ROUTE_PREFIX.JBR}/specDashboard`,
  SwingMarkDashboard = `${ROUTE_PREFIX.JBR}/swingmarkDashboard`,
  FleetTest = `${ROUTE_PREFIX.Fleet}/${TEST_ROUTE}`,
  FleetPerfDashboard = `${ROUTE_PREFIX.Fleet}/perfDashboard`,
  FleetPerfStartupComparisonDashboard = `${ROUTE_PREFIX.Fleet}/startupComparisonDashboard`,
  FleetStartupDashboard = `${ROUTE_PREFIX.Fleet}/startupDashboard`,
  FleetStartupExplore = `${ROUTE_PREFIX.Fleet}/startupExplore`,
  BazelTest = `${ROUTE_PREFIX.Bazel}/${TEST_ROUTE}`,
  BazelPluginDashboard = `${ROUTE_PREFIX.Bazel}/BazelPluginDashboard`,
  QodanaTest = `${ROUTE_PREFIX.Qodana}/${TEST_ROUTE}`,
  ClionStartupDashboard = `${ROUTE_PREFIX.Clion}/${STARTUP_ROUTE}`,
  ClionProductMetricsDashboard = `${ROUTE_PREFIX.Clion}/${PRODUCT_METRICS_ROUTE}`,
  ClionTest = `${ROUTE_PREFIX.Clion}/${DEV_TEST_ROUTE}`,
  ClionPerfDashboard = `${ROUTE_PREFIX.Clion}/perfDashboard`,
  ClionDetailedPerfDashboard = `${ROUTE_PREFIX.Clion}/detailedPerfDashboard`,
  ClionFindUsageDashboard = `${ROUTE_PREFIX.Clion}/findUsageDashboard`,
  ClionMemoryDashboard = `${ROUTE_PREFIX.Clion}/memoryDashboard`,
  ClionProjectModelDashboard = `${ROUTE_PREFIX.Clion}/projectModelDashboard`,
  ClionLaggingLatencyDashboard = `${ROUTE_PREFIX.Clion}/laggingLatencyDashboard`,
  CLionOldVsNewSeDashboard = `${ROUTE_PREFIX.Clion}/oldVsNewSeDashboard`,
  ClionCompareBranches = `${ROUTE_PREFIX.Clion}/${COMPARE_BRANCHES_ROUTE}`,
  VcsIdeaDashboard = `${ROUTE_PREFIX.Vcs}/idea`,
  VcsSpaceDashboard = `${ROUTE_PREFIX.Vcs}/space`,
  VcsStarterDashboard = `${ROUTE_PREFIX.Vcs}/starter`,
  VcsIdeaDashboardDev = `${ROUTE_PREFIX.Vcs}/ideaDev`,
  VcsSpaceDashboardDev = `${ROUTE_PREFIX.Vcs}/spaceDev`,
  VcsStarterDashboardDev = `${ROUTE_PREFIX.Vcs}/starterDev`,
  PerfUnitTests = `${ROUTE_PREFIX.PerfUnit}/${TEST_ROUTE}`,
  IJentBenchmarksDashboard = `${ROUTE_PREFIX.IJent}/benchmarksDashboard`,
  IJentPerfTestsDashboard = `${ROUTE_PREFIX.IJent}/performanceDashboard`,
  IJentProjectLoadingDashboard = `${ROUTE_PREFIX.IJent}/projectLoadingDashboard`,
  IJentRuntimeDashboard = `${ROUTE_PREFIX.IJent}/runtimeDashboard`,
  IJentRawPerfData = `${ROUTE_PREFIX.IJent}/rawPerfData`,
  MLDevTests = `${ROUTE_PREFIX.ML}/dev/${DEV_TEST_ROUTE}`,
  AIAssistantApiTests = `${ROUTE_PREFIX.ML}/dev/apiTests`,
  AIAssistantTestGeneration = `${ROUTE_PREFIX.ML}/dev/testGeneration`,
  LLMDevTests = `${ROUTE_PREFIX.ML}/dev/llmDashboardDev`,
  AIAPrivacyDashboard = `${ROUTE_PREFIX.ML}/dev/aiaPrivacyDashboard`,
  DataGripTests = `${ROUTE_PREFIX.DataGrip}/${TEST_ROUTE}`,
  DataGripProductMetricsDashboard = `${ROUTE_PREFIX.DataGrip}/${PRODUCT_METRICS_ROUTE}`,
  DataGripIndexingDashboard = `${ROUTE_PREFIX.DataGrip}/indexingDashboard`,
  DataGripDataGridRenderingDashboard = `${ROUTE_PREFIX.DataGrip}/dataGridRendering`,
  AIATests = `${ROUTE_PREFIX.AIA}/${TEST_ROUTE}`,
  AIACompletionDashboard = `${ROUTE_PREFIX.AIA}/completion`,
  AIACodeGenerationDashboard = `${ROUTE_PREFIX.AIA}/codeGeneration`,
  AIAChatCodeGenerationDashboard = `${ROUTE_PREFIX.AIA}/chatCodeGeneration`,
  AIANameSuggestionDashboard = `${ROUTE_PREFIX.AIA}/nameSuggestion`,
  AIATestGenerationDashboard = `${ROUTE_PREFIX.AIA}/testGeneration`,
  KMTTests = `${ROUTE_PREFIX.KMT}/unitTests`,
  KMTIntegrationTests = `${ROUTE_PREFIX.KMT}/${DEV_TEST_ROUTE}`,
  KMTDashboard = `${ROUTE_PREFIX.KMT}/${DASHBOARD_ROUTE}`,
  DiogenTests = `${ROUTE_PREFIX.Diogen}/${TEST_ROUTE}`,
  DiogenDashboard = `${ROUTE_PREFIX.Diogen}/${DASHBOARD_ROUTE}`,
  ToolboxTests = `${ROUTE_PREFIX.Toolbox}/${TEST_ROUTE}`,
  KotlinBuildToolsTests = `${ROUTE_PREFIX.KotlinBuildTools}/${TEST_ROUTE}`,
  ToolboxTestsGwDeployDashboard = `${ROUTE_PREFIX.Toolbox}/gw-deploy`,
  LSPTests = `${ROUTE_PREFIX.LSP}/${TEST_ROUTE}`,
  LSPDashboard = `${ROUTE_PREFIX.LSP}/${DASHBOARD_ROUTE}`,
  KotlinNotebooksTests = `${ROUTE_PREFIX.KotlinNotebooks}/${TEST_ROUTE}`,
  KotlinNotebooksDashboard = `${ROUTE_PREFIX.KotlinNotebooks}/${DASHBOARD_ROUTE}`,
  ReportDegradations = "/degradations/report",
  MetricsDescription = "/metrics/description",
  BisectLauncher = "/bisect/launcher",
  OwnersTest = "/owners/test",
  LlmAnalyses = "/analyses",
}

export interface Tab {
  url: ROUTES
  label: string
}

export interface SubProject {
  url: ROUTE_PREFIX
  label: string
  tabs: Tab[]
}

interface Product {
  url: ROUTE_PREFIX | ROUTES
  label: string
  children: SubProject[]
}

const TESTS_LABEL = "Tests"
const COMPARE_BUILDS_LABEL = "Compare Builds"
const COMPARE_BRANCHES_LABEL = "Compare Branches"
const COMPARE_MODES_LABEL = "Compare Modes"
const DASHBOARD_LABEL = "Dashboard"
const STARTUP_LABEL = "Startup"
const PRODUCT_METRICS_LABEL = "Product Metrics"

type LazyComponent = NonNullable<RouteRecordRaw["component"]>

// A single navigation tab.
function tab(url: ROUTES, label: string): Tab {
  return { url, label }
}

// A dashboard route: { path, component, meta: { pageTitle } }. `props` covers the few dashboards
// that pass static props to their component.
function dashboard(path: ROUTES, component: LazyComponent, pageTitle: string, props?: Record<string, unknown>): RouteRecordRaw {
  return props === undefined ? { path, component, meta: { pageTitle } } : { path, component, props, meta: { pageTitle } }
}

// Performance-tests page, typed against the shared PerformanceTests component props.
function perfTests(path: ROUTES, props: PerformanceTestsProps, pageTitle: string) {
  return { path, component: COMPONENTS.perfTests, props, meta: { pageTitle } } satisfies TypedRouteRecord<PerformanceTestsProps>
}

// Startup metrics dashboard page.
function startupDashboard(path: ROUTES, props: { table: string; defaultProject?: string }, pageTitle: string): RouteRecordRaw {
  return { path, component: COMPONENTS.startupDashboard, props, meta: { pageTitle } }
}

// Compare-builds / compare-branches / compare-modes pages; the page title defaults to the standard label.
function compareBuilds(path: ROUTES, props: { dbName: string; table: string }, pageTitle: string = COMPARE_BUILDS_LABEL): RouteRecordRaw {
  return { path, component: COMPONENTS.compareBuilds, props, meta: { pageTitle } }
}

function compareBranches(path: ROUTES, props: { dbName: string; table: string; metricsNames?: string[] }, pageTitle: string = COMPARE_BRANCHES_LABEL): RouteRecordRaw {
  return { path, component: COMPONENTS.compareBranches, props, meta: { pageTitle } }
}

function compareModes(path: ROUTES, props: { dbName: string; table: string }, pageTitle: string = COMPARE_MODES_LABEL): RouteRecordRaw {
  return { path, component: COMPONENTS.compareModes, props, meta: { pageTitle } }
}

const IJ_STARTUP: Product = {
  url: ROUTE_PREFIX.Startup,
  label: "IntelliJ Startup (deprecated)",
  children: [
    {
      url: ROUTE_PREFIX.Startup,
      label: "",
      tabs: [
        tab(ROUTES.StartupPulse, "Pulse"),
        tab(ROUTES.StartupPulseInstaller, "Pulse (Installer)"),
        tab(ROUTES.StartupModuleLoading, "Module Loading"),
        tab(ROUTES.StartupGcAndMemory, "GC and Memory"),
        tab(ROUTES.StartupProgress, "Progress Over Time"),
        tab(ROUTES.StartupExplore, "Explore"),
        tab(ROUTES.StartupExploreInstaller, "Explore (Installer)"),
      ],
    },
  ],
}
const IDEA: Product = {
  url: ROUTE_PREFIX.IntelliJ,
  label: "IDEA",
  children: [
    {
      url: ROUTE_PREFIX.IntelliJ,
      label: "Primary Functionality",
      tabs: [
        tab(ROUTES.IntelliJStartupDashboard, STARTUP_LABEL),
        tab(ROUTES.IntelliJProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.IntelliJIndexingDashboard, "Indexes"),
        tab(ROUTES.IntelliJJavaDashboard, "Java"),
        tab(ROUTES.IntelliJKotlinDashboard, "Kotlin"),
        tab(ROUTES.IntelliJUltimateDashboard, "Ultimate"),
        tab(ROUTES.IntelliJUIDashboard, "UI"),
        tab(ROUTES.IntelliJSearchEverywhereExDashboard, "Search Everywhere Ex"),
        tab(ROUTES.IntelliJDevTests, "Tests (Dev)"),
        tab(ROUTES.IntelliJCompareBranches, COMPARE_BRANCHES_LABEL),
        tab(ROUTES.IntelliJCompareModes, COMPARE_MODES_LABEL),
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJBuildTools,
      label: "Build Tools",
      tabs: [
        tab(ROUTES.IntelliJGradleDashboardDev, "Gradle Import DevServer"),
        tab(ROUTES.IntelliJMavenDashboardDev, "Maven Import DevServer"),
        tab(ROUTES.IntelliJMavenImportersConfiguratorsDashboardDev, "Maven Importers and Configurators DevServer"),
        tab(ROUTES.IntelliJJpsDashboardDev, "JPS Import DevServer"),
        tab(ROUTES.IntelliJGradleBenchmarks, "Gradle Benchmark"),
        tab(ROUTES.IntelliJBuildTests, TESTS_LABEL),
        tab(ROUTES.IntelliJBuildTestsDev, "Tests (DevServer)"),
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJSharedIndexes,
      label: "Shared Indexes",
      tabs: [tab(ROUTES.IntelliJSharedIndicesDashboard, DASHBOARD_LABEL), tab(ROUTES.IntelliJSharedIndicesTests, TESTS_LABEL)],
    },
    {
      url: ROUTE_PREFIX.IntelliJKotlinK2Performance,
      label: "Performance K2",
      tabs: [tab(ROUTES.IntelliJK2Dashboard, DASHBOARD_LABEL)],
    },
    {
      url: ROUTE_PREFIX.IntelliJPackageChecker,
      label: "Package Checker",
      tabs: [tab(ROUTES.IntelliJPackageCheckerDashboard, DASHBOARD_LABEL), tab(ROUTES.IntelliJPackageCheckerTests, TESTS_LABEL)],
    },
    {
      url: ROUTE_PREFIX.Vcs,
      label: "VCS",
      tabs: [
        tab(ROUTES.VcsIdeaDashboardDev, "Performance dashboard idea project DevServer"),
        tab(ROUTES.VcsSpaceDashboardDev, "Performance dashboard space project DevServer"),
        tab(ROUTES.VcsStarterDashboardDev, "Performance dashboard starter project DevServer"),
        tab(ROUTES.VcsIdeaDashboard, "Performance dashboard idea project (obsolete)"),
        tab(ROUTES.VcsSpaceDashboard, "Performance dashboard space project (obsolete)"),
        tab(ROUTES.VcsStarterDashboard, "Performance dashboard starter project (obsolete)"),
      ],
    },
    {
      url: ROUTE_PREFIX.EmbeddingSearch,
      label: "Embedding Search",
      tabs: [tab(ROUTES.IntelliJEmbeddingSearchDashboard, "Embedding Search")],
    },
  ],
}
const PHPSTORM: Product = {
  url: ROUTE_PREFIX.PhpStorm,
  label: "PhpStorm",
  children: [
    {
      url: ROUTE_PREFIX.PhpStorm,
      label: "",
      tabs: [
        tab(ROUTES.PhpStormStartupDashboard, STARTUP_LABEL),
        tab(ROUTES.PhpStormProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.PhpStormLLMDashboard, "LLM Dashboard"),
        tab(ROUTES.PhpStormIndexingDashboard, "Indexing"),
        tab(ROUTES.PhpStormInspectionsDashboard, "Inspections"),
        tab(ROUTES.PhpStormCodeEditingDashboard, "Code Editing"),
        tab(ROUTES.PhpStormUnitTestsDashboard, "Unit Tests"),
        tab(ROUTES.PhpStormDevTests, TESTS_LABEL),
        tab(ROUTES.PhpStormCompareBranches, COMPARE_BRANCHES_LABEL),
        tab(ROUTES.PhpStormCompareModes, COMPARE_MODES_LABEL),
        tab(ROUTES.PhpStormWithPluginsDashboard, "Dashboard with Plugins"),
        tab(ROUTES.PhpStormWithPluginsTests, "Tests with Plugins"),
      ],
    },
  ],
}
const KOTLIN: Product = {
  url: ROUTE_PREFIX.Kotlin,
  label: "Kotlin",
  children: [
    {
      url: ROUTE_PREFIX.Kotlin,
      label: "Performance dashboards",
      tabs: [
        tab(ROUTES.KotlinDashboardDev, DASHBOARD_LABEL + " (dev)"),
        tab(ROUTES.KotlinUserScenariosDashboardDev, "User Scenarios(dev)"),
        tab(ROUTES.KotlinDashboard, DASHBOARD_LABEL),
        tab(ROUTES.KotlinTests, TESTS_LABEL),
        tab(ROUTES.KotlinTestsDev, "Tests (dev)"),
        tab(ROUTES.KotlinCompareBranches, COMPARE_BRANCHES_LABEL),
        tab(ROUTES.KotlinCompareBranchesDev, COMPARE_BRANCHES_LABEL + " (dev)"),
      ],
    },
    {
      url: ROUTE_PREFIX.KotlinMemory,
      label: "Memory dashboards",
      tabs: [tab(ROUTES.KotlinMemoryDashboardDev, "Memory (dev)"), tab(ROUTES.KotlinMemoryDashboard, "Memory")],
    },
  ],
}
const GOLAND: Product = {
  url: ROUTE_PREFIX.GoLand,
  label: "GoLand",
  children: [
    {
      url: ROUTE_PREFIX.GoLand,
      label: "Primary Functionality",
      tabs: [
        tab(ROUTES.GoLandStartupDashboard, STARTUP_LABEL),
        tab(ROUTES.GoLandProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.GoLandIndexingDashboard, "Indexing"),
        tab(ROUTES.GoLandCompletionDashboard, "Completion"),
        tab(ROUTES.GoLandHighlightingDashboard, "Highlighting"),
        tab(ROUTES.GoLandInspectionDashboard, "Inspections"),
        tab(ROUTES.GoLandDebuggerDashboard, "Debugger"),
        tab(ROUTES.GoLandFindUsagesDashboard, "Find Usages"),
        tab(ROUTES.GoLandDFADashboard, "DFA"),
        tab(ROUTES.GoLandTests, TESTS_LABEL),
        tab(ROUTES.GoLandCompareBranches, COMPARE_BRANCHES_LABEL),
        tab(ROUTES.GoLandCompareModes, COMPARE_MODES_LABEL),
      ],
    },
  ],
}
const RUBYMINE: Product = {
  url: ROUTE_PREFIX.RubyMine,
  label: "RubyMine",
  children: [
    {
      url: ROUTE_PREFIX.RubyMine,
      label: "",
      tabs: [
        tab(ROUTES.RubyStartupDashboard, STARTUP_LABEL),
        tab(ROUTES.RubyMineProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.RubyMineDashboard, DASHBOARD_LABEL),
        tab(ROUTES.RubyMineInspectionsDashBoard, "Inspections"),
        tab(ROUTES.RubyMineIndexingDashBoard, "Indexing"),
        tab(ROUTES.RubyMineTestsDev, TESTS_LABEL),
        tab(ROUTES.RubyMineCompareBranches, COMPARE_BRANCHES_LABEL),
        tab(ROUTES.RubyMineCompareModes, COMPARE_MODES_LABEL),
      ],
    },
  ],
}

const PYCHARM: Product = {
  url: ROUTE_PREFIX.PyCharm,
  label: "PyCharm",
  children: [
    {
      url: ROUTE_PREFIX.PyCharm,
      label: "Primary Functionality",
      tabs: [
        tab(ROUTES.PyCharmStartupDashboard, STARTUP_LABEL),
        tab(ROUTES.PyCharmProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.PyCharmDashboard, DASHBOARD_LABEL),
        tab(ROUTES.PyCharmExternalTypeProviders, "External Type Providers"),
        tab(ROUTES.PyCharmExternalTypeProvidersUnitPerfTests, "External Type Providers Unit Performance Tests"),
        tab(ROUTES.PyCharmPerfUnitTests, "Performance Unit Tests"),
        tab(ROUTES.PyCharmOldDashboard, DASHBOARD_LABEL + " (Old)"),
        tab(ROUTES.PyCharmDevTests, TESTS_LABEL),
        tab(ROUTES.PyCharmTests, TESTS_LABEL + " (Old)"),
        tab(ROUTES.PyCharmCompareBranches, COMPARE_BRANCHES_LABEL),
      ],
    },
  ],
}

const WEBSTORM: Product = {
  url: ROUTE_PREFIX.WebStorm,
  label: "WebStorm",
  children: [
    {
      url: ROUTE_PREFIX.WebStorm,
      label: "",
      tabs: [
        tab(ROUTES.WebStormStartupDashboard, STARTUP_LABEL),
        tab(ROUTES.WebStormProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.WebStormDashboardBuiltInVsNEXT, "Built-in vs NEXT"),
        tab(ROUTES.WebStormDashboardDelicateProjects, "Delicate Projects"),
        tab(ROUTES.WebStormTests, TESTS_LABEL),
        tab(ROUTES.WebStormCompareBranches, COMPARE_BRANCHES_LABEL),
      ],
    },
  ],
}

const RUST: Product = {
  url: ROUTE_PREFIX.Rust,
  label: "Rust",
  children: [
    {
      url: ROUTE_PREFIX.Rust,
      label: "",
      tabs: [
        tab(ROUTES.RustRoverFirstStartupDashboardDev, "RustRover First Startup Dashboard"),
        tab(ROUTES.RustRoverProductMetricsDashboardDev, PRODUCT_METRICS_LABEL),
        tab(ROUTES.RustRoverDashboardDev, "RustRover Dashboard"),
        tab(ROUTES.RustRoverDebuggerDashboardDev, "Debugger"),
        tab(ROUTES.RustRoverRefactoringDashboardDev, "Refactoring"),
        tab(ROUTES.RustRoverUnitTestsDashboardDev, "Unit Tests"),
        tab(ROUTES.RustTestsDev, TESTS_LABEL),
        tab(ROUTES.RustCompareBranchesDev, COMPARE_BRANCHES_LABEL),
      ],
    },
  ],
}
const SCALA: Product = {
  url: ROUTE_PREFIX.Scala,
  label: "Scala",
  children: [
    {
      url: ROUTE_PREFIX.Scala,
      label: "",
      tabs: [tab(ROUTES.ScalaTests, TESTS_LABEL), tab(ROUTES.ScalaCompareBranches, COMPARE_BRANCHES_LABEL)],
    },
  ],
}
const JBR: Product = {
  url: ROUTE_PREFIX.JBR,
  label: "JBR",
  children: [
    {
      url: ROUTE_PREFIX.JBR,
      label: "",
      tabs: [
        tab(ROUTES.DaCapoDashboard, "DaCapo"),
        tab(ROUTES.J2DBenchDashboard, "J2DBench"),
        tab(ROUTES.JavaDrawDashboard, "JavaDraw"),
        tab(ROUTES.RenderDashboard, "Render"),
        tab(ROUTES.SPECjbb2015Dashboard, "SPECjbb2015"),
        tab(ROUTES.SwingMarkDashboard, "SwingMark"),
        tab(ROUTES.MapBenchDashboard, "MapBench"),
        tab(ROUTES.JBRTests, TESTS_LABEL),
      ],
    },
  ],
}
const FLEET: Product = {
  url: ROUTE_PREFIX.Fleet,
  label: "Fleet",
  children: [
    {
      url: ROUTE_PREFIX.Fleet,
      label: "",
      tabs: [
        tab(ROUTES.FleetStartupDashboard, "Startup Dashboard"),
        tab(ROUTES.FleetStartupExplore, "Startup Explore"),
        tab(ROUTES.FleetPerfDashboard, "Performance Dashboard"),
        tab(ROUTES.FleetPerfStartupComparisonDashboard, "Startup Comparison Dashboard"),
        tab(ROUTES.FleetTest, TESTS_LABEL),
      ],
    },
  ],
}

const BAZEL: Product = {
  url: ROUTE_PREFIX.Bazel,
  label: "Bazel",
  children: [
    {
      url: ROUTE_PREFIX.Bazel,
      label: "",
      tabs: [tab(ROUTES.BazelPluginDashboard, "Bazel Plugin Dashboard"), tab(ROUTES.BazelTest, TESTS_LABEL)],
    },
  ],
}

const QODANA: Product = {
  url: ROUTE_PREFIX.Qodana,
  label: "Qodana",
  children: [
    {
      url: ROUTE_PREFIX.Qodana,
      label: "",
      tabs: [tab(ROUTES.QodanaTest, TESTS_LABEL)],
    },
  ],
}

const CLION: Product = {
  url: ROUTE_PREFIX.Clion,
  label: "CLion",
  children: [
    {
      url: ROUTE_PREFIX.Clion,
      label: "",
      tabs: [
        tab(ROUTES.ClionStartupDashboard, "CLion Startup"),
        tab(ROUTES.ClionProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.ClionPerfDashboard, "Performance"),
        tab(ROUTES.ClionDetailedPerfDashboard, "Detailed Performance"),
        tab(ROUTES.ClionFindUsageDashboard, "Find Usages"),
        tab(ROUTES.ClionMemoryDashboard, "Memory"),
        tab(ROUTES.ClionProjectModelDashboard, "Project Model"),
        tab(ROUTES.ClionLaggingLatencyDashboard, "Lagging/Latency"),
        tab(ROUTES.CLionOldVsNewSeDashboard, "Old vs New SE"),
        tab(ROUTES.ClionTest, TESTS_LABEL),
        tab(ROUTES.ClionCompareBranches, COMPARE_BRANCHES_LABEL),
      ],
    },
  ],
}

const DATAGRIP: Product = {
  url: ROUTE_PREFIX.DataGrip,
  label: "DataGrip",
  children: [
    {
      url: ROUTE_PREFIX.DataGrip,
      label: "",
      tabs: [
        tab(ROUTES.DataGripProductMetricsDashboard, PRODUCT_METRICS_LABEL),
        tab(ROUTES.DataGripIndexingDashboard, "Indexing"),
        tab(ROUTES.DataGripDataGridRenderingDashboard, "Data Grid Rendering"),
        tab(ROUTES.DataGripTests, TESTS_LABEL),
      ],
    },
  ],
}

const PERF_UNIT: Product = {
  url: ROUTE_PREFIX.PerfUnit,
  label: "Perf Unit Tests",
  children: [
    {
      url: ROUTE_PREFIX.PerfUnit,
      label: "",
      tabs: [tab(ROUTES.PerfUnitTests, "Tests")],
    },
  ],
}

const IJENT: Product = {
  url: ROUTE_PREFIX.IJent,
  label: "IJent",
  children: [
    {
      url: ROUTE_PREFIX.IJent,
      label: "",
      tabs: [
        tab(ROUTES.IJentBenchmarksDashboard, "Benchmarks Dashboard"),
        tab(ROUTES.IJentPerfTestsDashboard, "Performance Dashboard"),
        tab(ROUTES.IJentProjectLoadingDashboard, "Project Loading (Community)"),
        tab(ROUTES.IJentRuntimeDashboard, "Runtime (Community)"),
        tab(ROUTES.IJentRawPerfData, "Raw Performance Data"),
      ],
    },
  ],
}

const ML_TESTS: Product = {
  url: ROUTE_PREFIX.ML,
  label: "ML Tests",
  children: [
    {
      url: ROUTE_PREFIX.ML,
      label: "",
      tabs: [
        tab(ROUTES.AIAssistantApiTests, "AI Assistant Api Tests"),
        tab(ROUTES.AIAssistantTestGeneration, "Test generation"),
        tab(ROUTES.LLMDevTests, "AIA Dashboard"),
        tab(ROUTES.AIAPrivacyDashboard, "AIA Privacy Dashboard"),
        tab(ROUTES.MLDevTests, "ML Tests on dev-server/fast-installer"),
      ],
    },
  ],
}

const AIA: Product = {
  url: ROUTE_PREFIX.AIA,
  label: "AIA",
  children: [
    {
      url: ROUTE_PREFIX.AIA,
      label: "",
      tabs: [
        tab(ROUTES.AIATests, "All"),
        tab(ROUTES.AIACompletionDashboard, "Completion"),
        tab(ROUTES.AIACodeGenerationDashboard, "Code Generation"),
        tab(ROUTES.AIAChatCodeGenerationDashboard, "Chat Code Generation"),
        tab(ROUTES.AIANameSuggestionDashboard, "Name Suggestion"),
        tab(ROUTES.AIATestGenerationDashboard, "Test Generation"),
      ],
    },
  ],
}

const KMT: Product = {
  url: ROUTE_PREFIX.KMT,
  label: "KMT",
  children: [
    {
      url: ROUTE_PREFIX.KMT,
      label: "",
      tabs: [tab(ROUTES.KMTDashboard, "Dashboard"), tab(ROUTES.KMTIntegrationTests, "Integration Tests"), tab(ROUTES.KMTTests, "Unit Tests")],
    },
  ],
}

const DIOGEN: Product = {
  url: ROUTE_PREFIX.Diogen,
  label: "Diogen",
  children: [
    {
      url: ROUTE_PREFIX.Diogen,
      label: "",
      tabs: [tab(ROUTES.DiogenDashboard, "Pipeline"), tab(ROUTES.DiogenTests, "All")],
    },
  ],
}

const TOOLBOX: Product = {
  url: ROUTE_PREFIX.Toolbox,
  label: "Toolbox",
  children: [
    {
      url: ROUTE_PREFIX.Toolbox,
      label: "",
      tabs: [tab(ROUTES.ToolboxTests, "All"), tab(ROUTES.ToolboxTestsGwDeployDashboard, "GW Deploy")],
    },
  ],
}

const KOTLIN_BUILD_TOOLS: Product = {
  url: ROUTE_PREFIX.KotlinBuildTools,
  label: "Kotlin Build Tools",
  children: [
    {
      url: ROUTE_PREFIX.KotlinBuildTools,
      label: "",
      tabs: [tab(ROUTES.KotlinBuildToolsTests, "All")],
    },
  ],
}

const LSP: Product = {
  url: ROUTE_PREFIX.LSP,
  label: "LSP",
  children: [
    {
      url: ROUTE_PREFIX.LSP,
      label: "",
      tabs: [tab(ROUTES.LSPDashboard, "Dashboard"), tab(ROUTES.LSPTests, "Tests")],
    },
  ],
}

const KOTLIN_NOTEBOOKS: Product = {
  url: ROUTE_PREFIX.KotlinNotebooks,
  label: "Kotlin Notebooks",
  children: [
    {
      url: ROUTE_PREFIX.KotlinNotebooks,
      label: "",
      tabs: [tab(ROUTES.KotlinNotebooksDashboard, DASHBOARD_LABEL), tab(ROUTES.KotlinNotebooksTests, TESTS_LABEL)],
    },
  ],
}

export const PRODUCTS = [
  AIA,
  BAZEL,
  CLION,
  DATAGRIP,
  DIOGEN,
  FLEET,
  GOLAND,
  IDEA,
  IJENT,
  IJ_STARTUP,
  JBR,
  KMT,
  KOTLIN,
  KOTLIN_BUILD_TOOLS,
  KOTLIN_NOTEBOOKS,
  LSP,
  ML_TESTS,
  PERF_UNIT,
  PHPSTORM,
  PYCHARM,
  QODANA,
  RUBYMINE,
  RUST,
  SCALA,
  TOOLBOX,
  WEBSTORM,
]

export function getNavigationElement(path: string): Product {
  const prefix = "/" + path.split("/")[1]
  return PRODUCTS.find((PRODUCTS) => prefix == PRODUCTS.url) ?? PRODUCTS[0]
}

const startupRoutes = [
  dashboard(ROUTES.StartupPulse, () => import("./components/startup/IntelliJPulse.vue"), "Pulse"),
  dashboard(ROUTES.StartupPulseInstaller, () => import("./components/startup/IntelliJPulse.vue"), "Pulse", { withInstaller: true }),
  dashboard(ROUTES.StartupModuleLoading, () => import("./components/startup/IntelliJModuleLoading.vue"), "Module Loading"),
  dashboard(ROUTES.StartupGcAndMemory, () => import("./components/startup/GcAndMemory.vue"), "GC and Memory"),
  dashboard(ROUTES.StartupProgress, () => import("./components/startup/IntelliJProgressOverTime.vue"), "Progress Over Time"),
  dashboard(ROUTES.StartupExplore, () => import("./components/startup/IntelliJExplore.vue"), "Explore", { withInstaller: false }),
  dashboard(ROUTES.StartupExploreInstaller, () => import("./components/startup/IntelliJExplore.vue"), "Explore (Installer)", { withInstaller: true }),
]

const intellijRoutes = [
  startupDashboard(ROUTES.IntelliJStartupDashboard, { table: "idea", defaultProject: "idea" }, "IDEA Startup dashboard"),
  dashboard(ROUTES.IntelliJProductMetricsDashboard, () => import("./components/intelliJ/ProductMetricsDashboard.vue"), "IDEA product metrics"),
  dashboard(ROUTES.IntelliJIndexingDashboard, () => import("./components/intelliJ/IndexingDashboard.vue"), "IntelliJ Indexing Performance dashboard"),
  dashboard(ROUTES.IntelliJJavaDashboard, () => import("./components/intelliJ/JavaDashboard.vue"), "IntelliJ Java Performance dashboard"),
  dashboard(ROUTES.IntelliJKotlinDashboard, () => import("./components/intelliJ/KotlinDashboard.vue"), "IntelliJ Kotlin Performance dashboard"),
  dashboard(ROUTES.IntelliJUltimateDashboard, () => import("./components/intelliJ/UltimateDashboard.vue"), "IntelliJ Ultimate Performance dashboard"),
  dashboard(ROUTES.IntelliJUIDashboard, () => import("./components/intelliJ/UIDashboard.vue"), "IntelliJ UI Performance dashboard"),
  dashboard(ROUTES.IntelliJSearchEverywhereExDashboard, () => import("./components/intelliJ/SearchEverywhereExDashboard.vue"), "IntelliJ Search Everywhere Ex dashboard"),
  dashboard(ROUTES.IntelliJK2Dashboard, () => import("./components/intelliJ/PerformanceK2Dashboard.vue"), "IntelliJ Performance K2 dashboard"),
  dashboard(ROUTES.IntelliJGradleDashboardDev, () => import("./components/intelliJ/build-tools/gradle/GradleImportPerformanceDashboardDevServer.vue"), "Gradle Import DevServer"),
  dashboard(
    ROUTES.IntelliJMavenDashboardDev,
    () => import("./components/intelliJ/build-tools/maven/MavenImportPerformanceDashboardDevServer.vue"),
    "Maven Import dashboard DevServer"
  ),
  dashboard(
    ROUTES.IntelliJMavenImportersConfiguratorsDashboardDev,
    () => import("./components/intelliJ/build-tools/maven/MavenImportersAndConfiguratorsPerformanceDashboardDevServer.vue"),
    "Maven Importers And Configurators dashboard DevServer"
  ),
  dashboard(ROUTES.IntelliJJpsDashboardDev, () => import("./components/intelliJ/build-tools/jps/JpsImportPerformanceDashboardDevServer.vue"), "JPS Import dashboard DevServer"),
  {
    path: ROUTES.IntelliJGradleBenchmarks,
    component: () => import("./components/common/PerformanceUnitTests.vue"),
    props: {
      dbName: "perfUnitTests",
      table: "report",
      initialMachine: MACHINES.HETZNER,
      withInstaller: false,
      projectFilter: "com.intellij.gradle.%benchmark%",
      persistentId: "gradleBenchmarks-dashboard",
      preselectAll: true,
    },
    meta: { pageTitle: "Gradle Benchmark" },
  } satisfies TypedRouteRecord<PerformanceUnitTestsProps>,
  dashboard(ROUTES.IntelliJPackageCheckerDashboard, () => import("./components/intelliJ/PackageCheckerDashboard.vue"), "Package Checker"),
  dashboard(ROUTES.IntelliJSharedIndicesDashboard, () => import("./components/intelliJ/SharedIndexesDashboard.vue"), "Shared Indexes Performance Dashboard"),
  dashboard(ROUTES.IntelliJEmbeddingSearchDashboard, () => import("./components/intelliJ/embeddingSearch/Dashboard.vue"), "IntelliJ performance tests for embedding search"),
  {
    path: `${ROUTE_PREFIX.IntelliJ}/:subproject?/${TEST_ROUTE}`,
    component: COMPONENTS.perfTests,
    props: {
      dbName: "perfint",
      table: "idea",
      initialMachine: MACHINES.AWS_LINUX,
    },
    meta: { pageTitle: "IntelliJ Performance tests" },
  } satisfies TypedRouteRecord<PerformanceTestsProps>,
  {
    path: `${ROUTE_PREFIX.IntelliJ}/:subproject?/${DEV_TEST_ROUTE}`,
    component: COMPONENTS.perfTests,
    props: {
      dbName: "perfintDev",
      table: "idea",
      initialMachine: MACHINES.AWS_LINUX,
      withInstaller: false,
    },
    meta: { pageTitle: "IntelliJ Integration Performance Tests On DevServer" },
  } satisfies TypedRouteRecord<PerformanceTestsProps>,
  {
    path: ROUTES.IntelliJCompare,
    component: () => COMPONENTS.compareBuilds,
    props: {
      dbName: "perfint",
      table: "idea",
    },
    meta: { pageTitle: COMPARE_BUILDS_LABEL },
  },
  compareBranches(ROUTES.IntelliJCompareBranches, { dbName: "perfintDev", table: "idea" }),
  compareModes(ROUTES.IntelliJCompareModes, { dbName: "perfintDev", table: "idea" }),
]

const phpstormRoutes = [
  startupDashboard(ROUTES.PhpStormStartupDashboard, { table: "phpstorm", defaultProject: "stitcher_with_composer" }, "PhpStorm Startup dashboard"),
  dashboard(ROUTES.PhpStormProductMetricsDashboard, () => import("./components/phpstorm/ProductMetricsDashboard.vue"), "PhpStorm product metrics"),
  dashboard(ROUTES.PhpStormLLMDashboard, () => import("./components/phpstorm/MLDashboard.vue"), "PhpStorm LLM Performance dashboard"),
  dashboard(ROUTES.PhpStormIndexingDashboard, () => import("./components/phpstorm/IndexingDashboard.vue"), "PhpStorm Indexing Dashboard"),
  dashboard(ROUTES.PhpStormInspectionsDashboard, () => import("./components/phpstorm/InspectionsDashboard.vue"), "PhpStorm Inspections Dashboard"),
  dashboard(ROUTES.PhpStormCodeEditingDashboard, () => import("./components/phpstorm/CodeEditingDashboard.vue"), "PhpStorm Code Editing Dashboard"),
  dashboard(ROUTES.PhpStormUnitTestsDashboard, () => import("./components/phpstorm/UnitTestsDashboard.vue"), "PhpStorm Unit Tests Dashboard"),
  dashboard(ROUTES.PhpStormWithPluginsDashboard, () => import("./components/phpstorm/PerformanceDashboardWithPlugins.vue"), "PhpStorm With Plugins Performance dashboard"),
  perfTests(ROUTES.PhpStormWithPluginsTests, { dbName: "perfint", table: "phpstormWithPlugins", initialMachine: MACHINES.HETZNER }, "PhpStorm Performance tests with plugins"),
  perfTests(ROUTES.PhpStormTests, { dbName: "perfint", table: "phpstorm", initialMachine: MACHINES.HETZNER }, "PhpStorm Performance tests"),
  perfTests(ROUTES.PhpStormDevTests, { dbName: "perfintDev", table: "phpstorm", initialMachine: MACHINES.HETZNER, withInstaller: false }, "PhpStorm Performance tests"),
  compareBranches(ROUTES.PhpStormCompareBranches, { dbName: "perfintDev", table: "phpstorm" }),
  compareModes(ROUTES.PhpStormCompareModes, { dbName: "perfintDev", table: "phpstorm" }),
]

const golandRoutes = [
  dashboard(ROUTES.GoLandInspectionDashboard, () => import("./components/goland/InspectionsDashboard.vue"), "GoLand Inspections dashboard"),
  startupDashboard(ROUTES.GoLandStartupDashboard, { table: "goland", defaultProject: "pocketbase" }, "GoLand Startup dashboard"),
  dashboard(ROUTES.GoLandProductMetricsDashboard, () => import("./components/goland/ProductMetricsDashboard.vue"), "GoLand product metrics"),
  dashboard(ROUTES.GoLandIndexingDashboard, () => import("./components/goland/IndexingDashboard.vue"), "GoLand Indexing dashboard"),
  dashboard(ROUTES.GoLandCompletionDashboard, () => import("./components/goland/CompletionDashboard.vue"), "GoLand Completion dashboard"),
  dashboard(ROUTES.GoLandHighlightingDashboard, () => import("./components/goland/HighlightingDashboard.vue"), "GoLand Highlighting dashboard"),
  dashboard(ROUTES.GoLandDebuggerDashboard, () => import("./components/goland/DebuggerDashboard.vue"), "GoLand Debugger dashboard"),
  dashboard(ROUTES.GoLandFindUsagesDashboard, () => import("./components/goland/FindUsagesDashboard.vue"), "GoLand Find Usages dashboard"),
  dashboard(ROUTES.GoLandDFADashboard, () => import("./components/goland/DataFlowAnalysisDashboard.vue"), "GoLand DFA dashboard"),
  perfTests(ROUTES.GoLandTests, { dbName: "perfintDev", table: "goland", withInstaller: false, initialMachine: MACHINES.HETZNER }, "GoLand Performance tests"),
  compareBuilds(ROUTES.GoLandCompare, { dbName: "perfintDev", table: "goland" }),
  compareBranches(ROUTES.GoLandCompareBranches, { dbName: "perfintDev", table: "goland" }),
  compareModes(ROUTES.GoLandCompareModes, { dbName: "perfintDev", table: "goland" }),
]

const pycharmRoutes = [
  startupDashboard(ROUTES.PyCharmStartupDashboard, { table: "pycharm", defaultProject: "tensorflow" }, "PyCharm Startup dashboard"),
  dashboard(ROUTES.PyCharmProductMetricsDashboard, () => import("./components/pycharm/ProductMetricsDashboard.vue"), "PyCharm product metrics"),
  dashboard(ROUTES.PyCharmDashboard, () => import("./components/pycharm/PerformanceDashboard.vue"), "PyCharm Performance dashboard"),
  dashboard(ROUTES.PyCharmExternalTypeProviders, () => import("./components/pycharm/ExternalTypeProviders.vue"), "PyCharm External Type Providers dashboard"),
  dashboard(
    ROUTES.PyCharmExternalTypeProvidersUnitPerfTests,
    () => import("./components/pycharm/ExternalTypeProvidersUnitPerfTests.vue"),
    "PyCharm External Type Providers Unit Performance Tests dashboard"
  ),
  dashboard(ROUTES.PyCharmPerfUnitTests, () => import("./components/pycharm/PyPerfUnitTests.vue"), "PyCharm Performance Unit Tests dashboard"),
  dashboard(ROUTES.PyCharmOldDashboard, () => import("./components/pycharm/PerformanceDashboardOld.vue"), "PyCharm Performance dashboard"),
  perfTests(ROUTES.PyCharmTests, { dbName: "perfint", table: "pycharm", initialMachine: MACHINES.HETZNER }, "PyCharm Performance tests"),
  perfTests(ROUTES.PyCharmDevTests, { dbName: "perfintDev", table: "pycharm", initialMachine: MACHINES.HETZNER, withInstaller: false }, "PyCharm Performance tests"),
  compareBuilds(ROUTES.PyCharmCompare, { dbName: "perfint", table: "pycharm" }),
  compareBranches(ROUTES.PyCharmCompareBranches, { dbName: "perfintDev", table: "pycharm" }),
]

const webstormRoutes = [
  startupDashboard(ROUTES.WebStormStartupDashboard, { table: "webstorm", defaultProject: "angular" }, "WebStorm Startup dashboard"),
  dashboard(ROUTES.WebStormProductMetricsDashboard, () => import("./components/webstorm/ProductMetricsDashboard.vue"), "WebStorm product metrics"),
  perfTests(ROUTES.WebStormTests, { dbName: "perfintDev", withInstaller: false, table: "webstorm", initialMachine: MACHINES.HETZNER }, "WebStorm Performance tests"),
  dashboard(ROUTES.WebStormDashboardBuiltInVsNEXT, () => import("./components/webstorm/PerformanceDashboardBuiltInVsNEXT.vue"), "Built-in vs NEXT"),
  dashboard(ROUTES.WebStormDashboardDelicateProjects, () => import("./components/webstorm/PerformanceDashboardDelicateProjects.vue"), "Delicate Projects"),
  compareBuilds(ROUTES.WebStormCompare, { dbName: "perfintDev", table: "webstorm" }),
  compareBranches(ROUTES.WebStormCompareBranches, { dbName: "perfintDev", table: "webstorm" }),
]

const rubymineRoutes = [
  startupDashboard(ROUTES.RubyStartupDashboard, { table: "ruby", defaultProject: "diaspora" }, "Ruby Startup dashboard"),
  dashboard(ROUTES.RubyMineProductMetricsDashboard, () => import("./components/rubymine/ProductMetricsDevDashboard.vue"), "RubyMine product metrics"),
  dashboard(ROUTES.RubyMineDashboard, () => import("./components/rubymine/PerformanceDevDashboard.vue"), "RubyMine Performance Dashboard"),
  dashboard(ROUTES.RubyMineInspectionsDashBoard, () => import("./components/rubymine/InspectionsDevDashboard.vue"), "RubyMine Inspections Dashboard"),
  dashboard(ROUTES.RubyMineIndexingDashBoard, () => import("./components/rubymine/IndexingDevDashboard.vue"), "RubyMine Indexing Dashboard"),
  perfTests(ROUTES.RubyMineTestsDev, { dbName: "perfintDev", table: "ruby", initialMachine: MACHINES.AWS_LINUX, withInstaller: false }, "RubyMine Performance tests"),
  compareBuilds(ROUTES.RubyMineCompare, { dbName: "perfint", table: "ruby" }),
  compareBranches(ROUTES.RubyMineCompareBranches, { dbName: "perfintDev", table: "ruby" }),
  compareModes(ROUTES.RubyMineCompareModes, { dbName: "perfintDev", table: "ruby" }),
]

const rustRoutes = [
  dashboard(ROUTES.RustRoverProductMetricsDashboardDev, () => import("./components/rust/ProductMetricsDashboardDev.vue"), "RustRover product metrics"),
  dashboard(ROUTES.RustRoverDashboardDev, () => import("./components/rust/PerformanceDashboardRustRoverDev.vue"), "RustRover Performance dashboard"),
  dashboard(
    ROUTES.RustRoverFirstStartupDashboardDev,
    () => import("./components/rust/PerformanceDashboardRustRoverFirstStartupDev.vue"),
    "RustRover First Startup Performance dashboard"
  ),
  dashboard(ROUTES.RustRoverDebuggerDashboardDev, () => import("./components/rust/DebuggerDashboard.vue"), "RustRover Debugger dashboard"),
  dashboard(ROUTES.RustRoverRefactoringDashboardDev, () => import("./components/rust/RefactoringDashboard.vue"), "RustRover Refactoring dashboard"),
  dashboard(ROUTES.RustRoverUnitTestsDashboardDev, () => import("./components/rust/UnitTestsDashboard.vue"), "Rust Unit Tests dashboard"),
  perfTests(ROUTES.RustTestsDev, { dbName: "perfintDev", table: "rust", initialMachine: MACHINES.AWS_LINUX, withInstaller: false }, "Rust Performance tests"),
  compareBranches(ROUTES.RustCompareBranchesDev, { dbName: "perfintDev", table: "rust" }),
]

const kotlinRoutes = [
  perfTests(ROUTES.KotlinTests, { dbName: "perfint", table: "kotlin", initialMachine: MACHINES.HETZNER }, "Kotlin Performance tests explore"),
  perfTests(
    ROUTES.KotlinTestsDev,
    { dbName: "perfintDev", table: "kotlin", initialMachine: MACHINES.HETZNER, withInstaller: false },
    "Kotlin Performance tests explore (dev/fast installer)"
  ),
  dashboard(ROUTES.KotlinDashboard, () => import("./components/kotlin/PerformanceDashboard.vue"), "Kotlin Performance dashboard"),
  dashboard(ROUTES.KotlinDashboardDev, () => import("./components/kotlin/dev/PerformanceDashboard.vue"), "Kotlin Performance dashboard (dev)"),
  dashboard(ROUTES.KotlinUserScenariosDashboardDev, () => import("./components/kotlin/dev/UserScenariosDashboard.vue"), "User scenarios dashboard (dev)"),
  dashboard(ROUTES.KotlinCodeAnalysisDev, () => import("./components/kotlin/dev/KotlinCodeAnalysisChartsDashboard.vue"), "Code analysis (dev)"),
  dashboard(ROUTES.KotlinCompletionDev, () => import("./components/kotlin/dev/CompletionDashboard.vue"), "Kotlin completion (dev/fast)"),
  dashboard(ROUTES.KotlinFindUsagesDev, () => import("./components/kotlin/dev/FindUsagesDashboard.vue"), "Kotlin findUsages (dev/fast)"),
  dashboard(ROUTES.KotlinRefactoringDev, () => import("./components/kotlin/dev/RefactoringDashboard.vue"), "Kotlin refactoring (dev/fast)"),
  dashboard(ROUTES.KotlinDebuggerDev, () => import("./components/kotlin/dev/DebuggerDashboard.vue"), "Kotlin debugger (dev/fast)"),
  dashboard(ROUTES.KotlinScriptDev, () => import("./components/kotlin/dev/ScriptDashboard.vue"), "Kts (dev/fast)"),
  compareBuilds(ROUTES.KotlinCompare, { dbName: "perfint", table: "kotlin" }),
  compareBranches(ROUTES.KotlinCompareBranches, { dbName: "perfint", table: "kotlin", metricsNames: KOTLIN_MAIN_METRICS }),
  compareBranches(ROUTES.KotlinCompareBranchesDev, { dbName: "perfintDev", table: "kotlin", metricsNames: KOTLIN_MAIN_METRICS }, COMPARE_BRANCHES_LABEL + "(dev/fast)"),
  dashboard(ROUTES.KotlinMemoryDashboard, () => import("./components/kotlin/MemoryPerformanceDashboard.vue"), "Memory"),
  dashboard(ROUTES.KotlinMemoryDashboardDev, () => import("./components/kotlin/dev/MemoryPerformanceDashboard.vue"), "Memory (dev)"),
]

const scalaRoutes = [
  perfTests(ROUTES.ScalaTests, { dbName: "perfint", table: "scala", initialMachine: MACHINES.HETZNER }, "Scala Performance tests"),
  compareBuilds(ROUTES.ScalaCompare, { dbName: "perfint", table: "scala" }),
  compareBranches(ROUTES.ScalaCompareBranches, { dbName: "perfint", table: "scala" }),
]

const jbrRoutes = [
  dashboard(ROUTES.JBRTests, () => import("./components/jbr/PerformanceTests.vue"), "JBR Performance tests"),
  dashboard(ROUTES.MapBenchDashboard, () => import("./components/jbr/MapBenchDashboard.vue"), "MapBench Dashboard"),
  dashboard(ROUTES.DaCapoDashboard, () => import("./components/jbr/DaCapoDashboard.vue"), "DaCapo Dashboard"),
  dashboard(ROUTES.J2DBenchDashboard, () => import("./components/jbr/J2DBenchDashboard.vue"), "J2DBench Dashboard"),
  dashboard(ROUTES.JavaDrawDashboard, () => import("./components/jbr/JavaDrawDashboard.vue"), "JavaDraw Dashboard"),
  dashboard(ROUTES.RenderDashboard, () => import("./components/jbr/RenderDashboard.vue"), "Render Dashboard"),
  dashboard(ROUTES.SPECjbb2015Dashboard, () => import("./components/jbr/SPECjbb2015Dashboard.vue"), "Spec Dashboard"),
  dashboard(ROUTES.SwingMarkDashboard, () => import("./components/jbr/SwingMarkDashboard.vue"), "SwingMark Dashboard"),
]

const fleetRoutes = [
  perfTests(ROUTES.FleetTest, { dbName: "fleet", table: "measure_new", initialMachine: MACHINES.HETZNER, withInstaller: false }, "Fleet Performance tests"),
  dashboard(ROUTES.FleetPerfDashboard, () => import("./components/fleet/PerformanceDashboard.vue"), "Fleet Performance dashboard"),
  dashboard(ROUTES.FleetPerfStartupComparisonDashboard, () => import("./components/fleet/StartupComparisonDashboard.vue"), "Fleet Startup Comparison dashboard"),
  dashboard(ROUTES.FleetStartupDashboard, () => import("./components/fleet/FleetDashboard.vue"), "Fleet Startup dashboard"),
  dashboard(ROUTES.FleetStartupExplore, () => import("./components/fleet/FleetExplore.vue"), "Fleet Startup Explore", { withInstaller: true }),
]

const bazelRoutes = [
  perfTests(ROUTES.BazelTest, { dbName: "bazel", table: "report", initialMachine: "Linux EC2 m5ad.2xlarge (8 vCPU Xeon, 32 GB)", withInstaller: false }, "Bazel Performance tests"),
  dashboard(ROUTES.BazelPluginDashboard, () => import("./components/bazel/BazelPluginDashboard.vue"), "Bazel Plugin Dashboard"),
]

const qodanaRoutes = [
  perfTests(ROUTES.QodanaTest, { dbName: "qodana", table: "report", initialMachine: "Linux EC2 c5.xlarge fleet (4 vCPU, 8 GB)", withInstaller: false }, "Qodana tests"),
]

const clionRoutes = [
  perfTests(ROUTES.ClionTest, { dbName: "perfintDev", table: "clion", withInstaller: false, initialMachine: MACHINES.AWS_LINUX }, "CLion tests"),
  startupDashboard(ROUTES.ClionStartupDashboard, { table: "clion", defaultProject: "radler/radler/cmake" }, "CLion Startup dashboard"),
  dashboard(ROUTES.ClionProductMetricsDashboard, () => import("./components/clion/ProductMetricsDashboard.vue"), "CLion product metrics"),
  dashboard(ROUTES.ClionPerfDashboard, () => import("./components/clion/PerformanceDashboard.vue"), "CLion dashboard", { initialMachine: MACHINES.AWS_LINUX }),
  dashboard(ROUTES.ClionDetailedPerfDashboard, () => import("./components/clion/DetailedPerformanceDashboard.vue"), "CLion Detailed Performance dashboard", {
    initialMachine: MACHINES.AWS_LINUX,
  }),
  dashboard(ROUTES.ClionFindUsageDashboard, () => import("./components/clion/FindUsageDashboard.vue"), "CLion Find Usages dashboard", { initialMachine: MACHINES.AWS_LINUX }),
  dashboard(ROUTES.ClionMemoryDashboard, () => import("./components/clion/MemoryDashboard.vue"), "CLion Memory dashboard", { initialMachine: MACHINES.AWS_LINUX }),
  dashboard(ROUTES.ClionProjectModelDashboard, () => import("./components/clion/ProjectModelDashboard.vue"), "CLion Project Model dashboard", {
    initialMachine: MACHINES.AWS_LINUX,
  }),
  dashboard(ROUTES.ClionLaggingLatencyDashboard, () => import("./components/clion/CLionLaggingLatencyDashboard.vue"), "CLion Lagging/Latency dashboard"),
  dashboard(ROUTES.CLionOldVsNewSeDashboard, () => import("./components/clion/oldVsNewSeDashboard.vue"), "CLion Old vs New SE dashboard"),
  compareBranches(ROUTES.ClionCompareBranches, { dbName: "perfintDev", table: "clion" }),
]

const vcsRoutes = [
  dashboard(ROUTES.VcsIdeaDashboardDev, () => import("./components/vcs/PerformanceDashboardDev.vue"), "Vcs Idea performance dashboard DevServer"),
  dashboard(ROUTES.VcsSpaceDashboardDev, () => import("./components/vcs/PerformanceSpaceDashboardDev.vue"), "Vcs Space performance dashboard DevServer"),
  dashboard(ROUTES.VcsStarterDashboardDev, () => import("./components/vcs/PerformanceStarterDashboardDev.vue"), "Vcs Starer performance dashboard DevServer"),
  dashboard(ROUTES.VcsIdeaDashboard, () => import("./components/vcs/PerformanceDashboard.vue"), "Vcs Idea performance dashboard (obsolete)"),
  dashboard(ROUTES.VcsSpaceDashboard, () => import("./components/vcs/PerformanceSpaceDashboard.vue"), "Vcs Space performance dashboard (obsolete)"),
  dashboard(ROUTES.VcsStarterDashboard, () => import("./components/vcs/PerformanceStarterDashboard.vue"), "Vcs Starer performance dashboard (obsolete)"),
]

const datagripRoutes = [
  perfTests(
    ROUTES.DataGripTests,
    { dbName: "perfintDev", table: "datagrip", initialMachine: "Linux EC2 m5d.xlarge (4 vCPU Xeon, 16 GB)", withInstaller: false },
    "DataGrip Performance tests"
  ),
  dashboard(ROUTES.DataGripProductMetricsDashboard, () => import("./components/datagrip/ProductMetricsDashboard.vue"), "DataGrip product metrics"),
  dashboard(ROUTES.DataGripIndexingDashboard, () => import("./components/datagrip/IndexingDashboard.vue"), "DataGrip Indexing dashboard"),
  dashboard(ROUTES.DataGripDataGridRenderingDashboard, () => import("./components/datagrip/DataGridRenderingDashboard.vue"), "DataGrip Data Grid Rendering dashboard"),
]

const toolboxRoutes = [
  perfTests(
    ROUTES.ToolboxTests,
    {
      dbName: "toolbox",
      table: "report",
      withInstaller: false,
      branch: "refs/heads/main",
      initialMachine: "Linux EC2 M5d.xlarge (4 vCPU Xeon, 16 GB)",
      withoutAccidents: true,
    },
    "Toolbox"
  ),
  dashboard(ROUTES.ToolboxTestsGwDeployDashboard, () => import("./components/toolbox/GwDeployMetricsDashboard.vue"), "GW Dashboard", {
    dbName: "toolbox",
    table: "report",
    withInstaller: false,
    branch: "refs/heads/master",
    withoutAccidents: true,
  }),
]

const ijentRoutes = [
  dashboard(ROUTES.IJentBenchmarksDashboard, () => import("./components/ijent/IJentBenchmarskDashboard.vue"), "IJent Benchmarks Dashboard"),
  dashboard(ROUTES.IJentPerfTestsDashboard, () => import("./components/ijent/IJentPerformanceTestsDashboard.vue"), "IJent Performance Tests Dashboard"),
  dashboard(ROUTES.IJentProjectLoadingDashboard, () => import("./components/ijent/IJentProjectLoadingDashboard.vue"), "IJent Project Loading (Community)"),
  dashboard(ROUTES.IJentRuntimeDashboard, () => import("./components/ijent/IJentRuntimeDashboard.vue"), "IJent Runtime (Community)"),
  perfTests(ROUTES.IJentRawPerfData, { dbName: "perfintDev", table: "ijent", initialMachine: "Linux Munich i7-13700, 64 Gb", withInstaller: false }, "IJent Raw Performance Data"),
]

const mlRoutes = [
  dashboard(ROUTES.AIAssistantApiTests, () => import("./components/ml/dev/AiAssistantApiTests.vue"), "AI API Tests"),
  dashboard(ROUTES.AIAssistantTestGeneration, () => import("./components/ml/dev/TestGenerationDashboard.vue"), "Test generation"),
  dashboard(ROUTES.LLMDevTests, () => import("./components/ml/dev/LLMDashboard.vue"), "AIA dashboard"),
  dashboard(ROUTES.AIAPrivacyDashboard, () => import("./components/ml/dev/AIAPrivacyDashboard.vue"), "AIA Privacy"),
  perfTests(ROUTES.MLDevTests, { dbName: "perfintDev", table: "ml", initialMachine: MACHINES.AWS_LINUX, withInstaller: false }, "ML Tests dev-server"),
]

const aiaRoutes = [
  perfTests(ROUTES.AIATests, { dbName: "mlEvaluation", table: "report", withInstaller: false, branch: null, initialMachine: null }, "AIA Tests"),
  dashboard(ROUTES.AIACompletionDashboard, () => import("./components/aia/AIACompletionDashboard.vue"), "AIA completion dashboard"),
  dashboard(ROUTES.AIACodeGenerationDashboard, () => import("./components/aia/AIACodeGeneration.vue"), "AIA code generation dashboard"),
  dashboard(ROUTES.AIAChatCodeGenerationDashboard, () => import("./components/aia/ChatCodeGeneration.vue"), "AIA chat code generation dashboard"),
  dashboard(ROUTES.AIANameSuggestionDashboard, () => import("./components/aia/AIANameSuggestion.vue"), "AIA name suggestion dashboard"),
  dashboard(ROUTES.AIATestGenerationDashboard, () => import("./components/aia/AIATestGeneration.vue"), "AIA test generation dashboard"),
]

const kmtRoutes = [
  perfTests(ROUTES.KMTTests, { dbName: "perfintDev", table: "swift", withInstaller: false, branch: null, initialMachine: "Mac Cidr Performance" }, "KMT Unit Tests"),
  perfTests(
    ROUTES.KMTIntegrationTests,
    { dbName: "perfintDev", table: "kmt", withInstaller: false, branch: "master", initialMachine: "Mac Cidr Performance" },
    "KMT Integration Tests"
  ),
  dashboard(ROUTES.KMTDashboard, () => import("./components/kmt/PerformanceDashboard.vue"), "KMT Dashboard"),
]

const diogenRoutes = [
  perfTests(ROUTES.DiogenTests, { dbName: "diogen", table: "report", withInstaller: false, branch: "refs/heads/main", initialMachine: null, withoutAccidents: true }, "Diogen"),
  dashboard(ROUTES.DiogenDashboard, () => import("./components/diogen/DiogenPipelineDashboard.vue"), "Diogen Pipeline"),
]

const perfUnitTestsRoutes = [
  {
    path: ROUTES.PerfUnitTests,
    component: () => import("./components/common/PerformanceUnitTests.vue"),
    props: {
      dbName: "perfUnitTests",
      table: "report",
      initialMachine: MACHINES.HETZNER,
      withInstaller: false,
    },
    meta: { pageTitle: "Perf Unit Tests" },
  } satisfies TypedRouteRecord<PerformanceTestsProps>,
]

const kotlinBuildToolsRoutes = [
  perfTests(
    ROUTES.KotlinBuildToolsTests,
    { dbName: "perfintDev", table: "kotlinBuildTools", withInstaller: false, branch: "master", initialMachine: MACHINES.HETZNER },
    "Kotlin Build Tools Tests"
  ),
]

const lspRoutes = [
  {
    path: ROUTES.LSPTests,
    component: () => import("./components/common/PerformanceUnitTests.vue"),
    props: {
      dbName: "perfUnitTests",
      table: "report",
      initialMachine: MACHINES.HETZNER,
      withInstaller: false,
    },
    meta: { pageTitle: "LSP Tests" },
  } satisfies TypedRouteRecord<PerformanceTestsProps>,
  dashboard(ROUTES.LSPDashboard, () => import("./components/lsp/LSPDashboard.vue"), "LSP Dashboard"),
]

const kotlinNotebooksRoutes = [
  perfTests(
    ROUTES.KotlinNotebooksTests,
    { dbName: "perfintDev", table: "kotlinNotebooks", initialMachine: MACHINES.AWS_LINUX, withInstaller: false },
    "Kotlin Notebooks Performance tests"
  ),
  dashboard(ROUTES.KotlinNotebooksDashboard, () => import("./components/kotlinNotebooks/PerformanceDashboard.vue"), "Kotlin Notebooks Dashboard"),
]

export function getNewDashboardRoutes(): ParentRouteRecord[] {
  return [
    {
      children: [
        ...startupRoutes,
        ...intellijRoutes,
        ...phpstormRoutes,
        ...golandRoutes,
        ...pycharmRoutes,
        ...webstormRoutes,
        ...rubymineRoutes,
        ...rustRoutes,
        ...kotlinRoutes,
        ...scalaRoutes,
        ...jbrRoutes,
        ...fleetRoutes,
        ...bazelRoutes,
        ...qodanaRoutes,
        ...clionRoutes,
        ...vcsRoutes,
        ...datagripRoutes,
        ...toolboxRoutes,
        ...ijentRoutes,
        ...mlRoutes,
        ...aiaRoutes,
        ...kmtRoutes,
        ...diogenRoutes,
        ...perfUnitTestsRoutes,
        ...kotlinBuildToolsRoutes,
        ...lspRoutes,
        ...kotlinNotebooksRoutes,
        {
          path: ROUTES.ReportDegradations,
          component: () => import("./components/degradations/ReportDegradation.vue"),
          meta: { pageTitle: "Report degradations" },
          props: (route) => ({
            tests: route.query["tests"],
            build: route.query["build"],
            date: route.query["date"],
          }),
        },
        {
          path: ROUTES.MetricsDescription,
          component: () => import("./components/metrics/MetricDescriptions.vue"),
          meta: { pageTitle: "Metrics description" },
        },
        {
          path: ROUTES.BisectLauncher,
          component: () => import("./components/bisect/BisectLauncher.vue"),
          meta: { pageTitle: "Bisect launcher" },
          props: (route) => ({ ...route.query }),
        },
        {
          path: ROUTES.OwnersTest,
          component: () => import("./components/common/OwnerTests.vue"),
          meta: { pageTitle: "Performance tests" },
        },
        {
          path: ROUTES.LlmAnalyses,
          component: () => import("./components/analyses/AnalysesListPage.vue"),
          meta: { pageTitle: "LLM Analyses" },
        },
      ],
    },
  ]
}
