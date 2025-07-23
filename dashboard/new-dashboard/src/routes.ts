/* eslint-disable @typescript-eslint/prefer-literal-enum-member */

import { ParentRouteRecord } from "./components/common/route"
import { KOTLIN_MAIN_METRICS } from "./components/kotlin/projects"
import { eap } from "./configurators/ReleaseNightlyConfigurator"

const enum ROUTE_PREFIX {
  Startup = "/ij",
  IntelliJ = "/intellij",
  IntelliJBuildTools = "/intellij/buildTools",
  IntelliJUltimate = "/intellij/ultimate",
  IntelliJJava = "/intellij/java",
  IntelliJSharedIndexes = "/intellij/sharedIndexes",
  IntelliJIncrementalCompilation = "/intellij/incrementalCompilation",
  IntelliJKotlinK2Performance = "/intellij/kotlinK2Performance",
  IntelliJPackageChecker = "/intellij/packageChecker",
  IntelliJFus = "/intellij/fus",
  PhpStorm = "/phpstorm",
  GoLand = "/goland",
  RubyMine = "/rubymine",
  Kotlin = "/kotlin",
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
}

const TEST_ROUTE = "tests"
const DEV_TEST_ROUTE = "testsDev"
const DASHBOARD_ROUTE = "dashboard"
const STARTUP_ROUTE = "startup"
const STARTUP_METRICS_ROUTE = "startup-metrics"
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
  IntelliJNewStartupDashboard = `${ROUTE_PREFIX.IntelliJ}/${STARTUP_METRICS_ROUTE}`,
  IntelliJProductMetricsDashboard = `${ROUTE_PREFIX.IntelliJ}/${PRODUCT_METRICS_ROUTE}`,
  IntelliJDashboard = `${ROUTE_PREFIX.IntelliJ}/${DASHBOARD_ROUTE}`,
  IntelliJPopupsDashboard = `${ROUTE_PREFIX.IntelliJ}/popupsDashboard`,
  IntelliJLaggingLatencyDashboard = `${ROUTE_PREFIX.IntelliJ}/laggingLatencyDashboard`,
  IntelliJIndexingDashboard = `${ROUTE_PREFIX.IntelliJ}/indexingDashboard`,
  IntelliJIncrementalCompilationDashboard = `${ROUTE_PREFIX.IntelliJIncrementalCompilation}/${DASHBOARD_ROUTE}`,
  IntelliJFindUsagesDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboardFindUsages`,
  IntelliJSEDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboardSearchEverywhere`,
  IntelliJWSLDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboardWSL`,
  IntelliJEmbeddingSearchDashboard = `${ROUTE_PREFIX.EmbeddingSearch}/dashboard`,
  IntelliJK2Dashboard = `${ROUTE_PREFIX.IntelliJKotlinK2Performance}/${DASHBOARD_ROUTE}`,
  IntelliJTests = `${ROUTE_PREFIX.IntelliJ}/${TEST_ROUTE}`,
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
  IntelliJUltimateDashboard = `${ROUTE_PREFIX.IntelliJUltimate}/${DASHBOARD_ROUTE}`,
  IntelliJUltimateDashboardOld = `${ROUTE_PREFIX.IntelliJUltimate}/ultimateDashboardOld`,
  IntelliJJavaDashboard = `${ROUTE_PREFIX.IntelliJJava}/${DASHBOARD_ROUTE}`,
  IntelliJUltimateTests = `${ROUTE_PREFIX.IntelliJUltimate}/${DEV_TEST_ROUTE}`,
  IntelliJSharedIndicesDashboard = `${ROUTE_PREFIX.IntelliJSharedIndexes}/${DASHBOARD_ROUTE}`,
  IntelliJSharedIndicesTests = `${ROUTE_PREFIX.IntelliJSharedIndexes}/${TEST_ROUTE}`,
  IntelliJPackageCheckerDashboard = `${ROUTE_PREFIX.IntelliJPackageChecker}/${DASHBOARD_ROUTE}`,
  IntelliJPackageCheckerTests = `${ROUTE_PREFIX.IntelliJPackageChecker}/${TEST_ROUTE}`,
  PhpStormProductMetricsDashboard = `${ROUTE_PREFIX.PhpStorm}/${PRODUCT_METRICS_ROUTE}`,
  PhpStormLLMDashboard = `${ROUTE_PREFIX.PhpStorm}/llmDashboard`,
  PhpStormIndexingDashboard = `${ROUTE_PREFIX.PhpStorm}/indexingDashboard`,
  PhpStormInspectionsDashboard = `${ROUTE_PREFIX.PhpStorm}/inspectionsDashboard`,
  PhpStormCodeEditingDashboard = `${ROUTE_PREFIX.PhpStorm}/codeEditingDashboard`,
  PhpStormStartupDashboard = `${ROUTE_PREFIX.PhpStorm}/${STARTUP_ROUTE}`,
  PhpStormWithPluginsDashboard = `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests = `${ROUTE_PREFIX.PhpStorm}/${TEST_ROUTE}`,
  PhpStormDevTests = `${ROUTE_PREFIX.PhpStorm}/${DEV_TEST_ROUTE}`,
  PhpStormWithPluginsTests = `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,
  PhpStormCompare = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_ROUTE}`,
  PhpStormCompareBranches = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_BRANCHES_ROUTE}`,
  PhpStormCompareModes = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_MODES_ROUTE}`,
  KotlinDashboard = `${ROUTE_PREFIX.Kotlin}/${DASHBOARD_ROUTE}`,
  KotlinDashboardDev = `${ROUTE_PREFIX.Kotlin}/${DASHBOARD_ROUTE}Dev`,
  KotlinUserScenariosDashboardDev = `${ROUTE_PREFIX.Kotlin}/Scenarios${DASHBOARD_ROUTE}Dev`,
  KotlinCodeAnalysisDev = `${ROUTE_PREFIX.Kotlin}/codeAnalysisDev `,
  KotlinTests = `${ROUTE_PREFIX.Kotlin}/${TEST_ROUTE}`,
  KotlinTestsDev = `${ROUTE_PREFIX.Kotlin}/${DEV_TEST_ROUTE}`,
  KotlinCompletionDev = `${ROUTE_PREFIX.Kotlin}/completionDev`,
  KotlinHighlightingDev = `${ROUTE_PREFIX.Kotlin}/highlightingDev`,
  KotlinFindUsagesDev = `${ROUTE_PREFIX.Kotlin}/findUsagesDev`,
  KotlinRefactoringDev = `${ROUTE_PREFIX.Kotlin}/refactoringDev`,
  KotlinDebuggerDev = `${ROUTE_PREFIX.Kotlin}/debuggerDev`,
  KotlinScriptDev = `${ROUTE_PREFIX.Kotlin}/scriptDev`,
  KotlinK1VsK2Comparison = `${ROUTE_PREFIX.Kotlin}/k1VsK2Comparison`,
  KotlinK1VsK2ComparisonDev = `${ROUTE_PREFIX.Kotlin}/k1VsK2ComparisonDev`,
  KotlinCompare = `${ROUTE_PREFIX.Kotlin}/${COMPARE_ROUTE}`,
  KotlinMemoryDashboard = `${ROUTE_PREFIX.KotlinMemory}/dashboard`,
  KotlinMemoryDashboardDev = `${ROUTE_PREFIX.KotlinMemory}/dashboardDev`,
  KotlinCompareBranches = `${ROUTE_PREFIX.Kotlin}/${COMPARE_BRANCHES_ROUTE}`,
  KotlinCompareBranchesDev = `${ROUTE_PREFIX.Kotlin}/${COMPARE_BRANCHES_ROUTE}Dev`,
  GoLandStartupDashboard = `${ROUTE_PREFIX.GoLand}/${STARTUP_ROUTE}`,
  GoLandProductMetricsDashboard = `${ROUTE_PREFIX.GoLand}/${PRODUCT_METRICS_ROUTE}`,
  GoLandProductMetricsDashboardOld = `${ROUTE_PREFIX.GoLand}/${PRODUCT_METRICS_ROUTE}Old`,
  GoLandIndexingDashboard = `${ROUTE_PREFIX.GoLand}/indexingDashboard`,
  GoLandIndexingDashboardOld = `${ROUTE_PREFIX.GoLand}/indexingDashboardOld`,
  GoLandScanningDashboard = `${ROUTE_PREFIX.GoLand}/scanningDashboard`,
  GoLandScanningDashboardOld = `${ROUTE_PREFIX.GoLand}/scanningDashboardOld`,
  GoLandCompletionDashboard = `${ROUTE_PREFIX.GoLand}/completionDashboard`,
  GoLandCompletionDashboardOld = `${ROUTE_PREFIX.GoLand}/completionDashboardOld`,
  GoLandInspectionDashboard = `${ROUTE_PREFIX.GoLand}/inspectionsDashboard`,
  GoLandInspectionDashboardOld = `${ROUTE_PREFIX.GoLand}/inspectionsDashboardOld`,
  GoLandDebuggerDashboard = `${ROUTE_PREFIX.GoLand}/debuggerDashboard`,
  GoLandDebuggerDashboardOld = `${ROUTE_PREFIX.GoLand}/debuggerDashboardOld`,
  GoLandFindUsagesDashboard = `${ROUTE_PREFIX.GoLand}/findUsagesDashboard`,
  GoLandFindUsagesDashboardOld = `${ROUTE_PREFIX.GoLand}/findUsagesDashboardOld`,
  GoLandDFADashboard = `${ROUTE_PREFIX.GoLand}/dfaDashboard`,
  GoLandDFADashboardOld = `${ROUTE_PREFIX.GoLand}/dfaDashboardOld`,
  GoLandDistributiveSizeDashboard = `${ROUTE_PREFIX.GoLand}/distributiveDashboard`,
  GoLandTests = `${ROUTE_PREFIX.GoLand}/${TEST_ROUTE}Dev`,
  GoLandTestsOld = `${ROUTE_PREFIX.GoLand}/${TEST_ROUTE}`,
  GoLandCompare = `${ROUTE_PREFIX.GoLand}/${COMPARE_ROUTE}`,
  GoLandCompareBranches = `${ROUTE_PREFIX.GoLand}/${COMPARE_BRANCHES_ROUTE}`,
  PyCharmStartupDashboard = `${ROUTE_PREFIX.PyCharm}/${STARTUP_ROUTE}`,
  PyCharmProductMetricsDashboard = `${ROUTE_PREFIX.PyCharm}/${PRODUCT_METRICS_ROUTE}`,
  PyCharmDashboard = `${ROUTE_PREFIX.PyCharm}/${DASHBOARD_ROUTE}Dev`,
  PyCharmOldDashboard = `${ROUTE_PREFIX.PyCharm}/${DASHBOARD_ROUTE}`,
  PyCharmTests = `${ROUTE_PREFIX.PyCharm}/${TEST_ROUTE}`,
  PyCharmDevTests = `${ROUTE_PREFIX.PyCharm}/${DEV_TEST_ROUTE}`,
  PyCharmCompare = `${ROUTE_PREFIX.PyCharm}/${COMPARE_ROUTE}`,
  PyCharmCompareBranches = `${ROUTE_PREFIX.PyCharm}/${COMPARE_BRANCHES_ROUTE}`,
  WebStormStartupDashboard = `${ROUTE_PREFIX.WebStorm}/${STARTUP_ROUTE}`,
  WebStormProductMetricsDashboard = `${ROUTE_PREFIX.WebStorm}/${PRODUCT_METRICS_ROUTE}`,
  WebStormProductMetricsDashboardOld = `${ROUTE_PREFIX.WebStorm}/${PRODUCT_METRICS_ROUTE}Old`,
  WebStormDashboard = `${ROUTE_PREFIX.WebStorm}/${DASHBOARD_ROUTE}`,
  WebStormDashboardOld = `${ROUTE_PREFIX.WebStorm}/${DASHBOARD_ROUTE}Old`,
  WebStormDashboardBuiltInVsNEXT = `${ROUTE_PREFIX.WebStorm}/dashboardBuiltInVsNext`,
  WebStormDashboardBuiltInVsNEXTOld = `${ROUTE_PREFIX.WebStorm}/dashboardBuiltInVsNextOld`,
  WebStormDashboardDelicateProjects = `${ROUTE_PREFIX.WebStorm}/dashboardDelicateProjects`,
  WebStormTests = `${ROUTE_PREFIX.WebStorm}/${TEST_ROUTE}Dev`,
  WebStormTestsOld = `${ROUTE_PREFIX.WebStorm}/${TEST_ROUTE}`,
  WebStormCompare = `${ROUTE_PREFIX.WebStorm}/${COMPARE_ROUTE}`,
  WebStormCompareBranches = `${ROUTE_PREFIX.WebStorm}/${COMPARE_BRANCHES_ROUTE}`,
  RubyStartupDashboard = `${ROUTE_PREFIX.RubyMine}/${STARTUP_ROUTE}`,
  RubyMineProductMetricsDashboard = `${ROUTE_PREFIX.RubyMine}/${PRODUCT_METRICS_ROUTE}Dev`,
  RubyMineProductMetricsDashboardOld = `${ROUTE_PREFIX.RubyMine}/${PRODUCT_METRICS_ROUTE}`,
  RubyMineDashboardOld = `${ROUTE_PREFIX.RubyMine}/${DASHBOARD_ROUTE}`,
  RubyMineDashboard = `${ROUTE_PREFIX.RubyMine}/${DASHBOARD_ROUTE}Dev`,
  RubyMineIndexingDashBoardOld = `${ROUTE_PREFIX.RubyMine}/indexingDashboard`,
  RubyMineIndexingDashBoard = `${ROUTE_PREFIX.RubyMine}/indexingDashboardDev`,
  RubyMineInspectionsDashBoardOld = `${ROUTE_PREFIX.RubyMine}/inspectionsDashboard`,
  RubyMineInspectionsDashBoard = `${ROUTE_PREFIX.RubyMine}/inspectionsDashboardDev`,
  RubyMineTests = `${ROUTE_PREFIX.RubyMine}/${TEST_ROUTE}`,
  RubyMineTestsDev = `${ROUTE_PREFIX.RubyMine}/${DEV_TEST_ROUTE}`,
  RubyMineCompare = `${ROUTE_PREFIX.RubyMine}/${COMPARE_ROUTE}`,
  RubyMineCompareBranches = `${ROUTE_PREFIX.RubyMine}/${COMPARE_BRANCHES_ROUTE}`,
  RubyMineCompareModes = `${ROUTE_PREFIX.RubyMine}/${COMPARE_MODES_ROUTE}`,
  RustRoverDashboard = `${ROUTE_PREFIX.Rust}/rustPluginDashboard`,
  RustRoverProductMetricsDashboard = `${ROUTE_PREFIX.Rust}/${PRODUCT_METRICS_ROUTE}`,
  RustRoverFirstStartupDashboard = `${ROUTE_PREFIX.Rust}/rustRoverFirstStartupDashboard`,
  RustTests = `${ROUTE_PREFIX.Rust}/${TEST_ROUTE}`,
  RustCompare = `${ROUTE_PREFIX.Rust}/${COMPARE_ROUTE}`,
  RustCompareBranches = `${ROUTE_PREFIX.Rust}/${COMPARE_BRANCHES_ROUTE}`,
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
  ClionClassicStartupDashboard = `${ROUTE_PREFIX.Clion}/${STARTUP_ROUTE}`,
  ClionNovaStartupDashboard = `${ROUTE_PREFIX.Clion}/nova_${STARTUP_ROUTE}`,
  ClionProductMetricsDashboard = `${ROUTE_PREFIX.Clion}/${PRODUCT_METRICS_ROUTE}`,
  ClionProductMetricsDashboardOld = `${ROUTE_PREFIX.Clion}/${PRODUCT_METRICS_ROUTE}Old`,
  ClionTest = `${ROUTE_PREFIX.Clion}/${DEV_TEST_ROUTE}`,
  ClionTestOld = `${ROUTE_PREFIX.Clion}/${TEST_ROUTE}`,
  ClionPerfDashboard = `${ROUTE_PREFIX.Clion}/perfDashboard`,
  ClionPerfDashboardOld = `${ROUTE_PREFIX.Clion}/perfDashboardOld`,
  ClionDetailedPerfDashboard = `${ROUTE_PREFIX.Clion}/detailedPerfDashboard`,
  ClionDetailedPerfDashboardOld = `${ROUTE_PREFIX.Clion}/detailedPerfDashboardOld`,
  ClionMemoryDashboard = `${ROUTE_PREFIX.Clion}/memoryDashboard`,
  ClionMemoryDashboardOld = `${ROUTE_PREFIX.Clion}/memoryDashboardOld`,
  ClionProjectModelDashboard = `${ROUTE_PREFIX.Clion}/projectModelDashboard`,
  ClionProjectModelDashboardOld = `${ROUTE_PREFIX.Clion}/projectModelDashboardOld`,
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
  IJentBenchmarks = `${ROUTE_PREFIX.IJent}/benchmarks`,
  MLDevTests = `${ROUTE_PREFIX.ML}/dev/${DEV_TEST_ROUTE}`,
  AIAssistantApiTests = `${ROUTE_PREFIX.ML}/dev/apiTests`,
  AIAssistantTestGeneration = `${ROUTE_PREFIX.ML}/dev/testGeneration`,
  LLMDevTests = `${ROUTE_PREFIX.ML}/dev/llmDashboardDev`,
  AIAPrivacyDashboard = `${ROUTE_PREFIX.ML}/dev/aiaPrivacyDashboard`,
  DataGripStartupDashboard = `${ROUTE_PREFIX.DataGrip}/${STARTUP_ROUTE}`,
  DataGripProductMetricsDashboard = `${ROUTE_PREFIX.DataGrip}/${PRODUCT_METRICS_ROUTE}`,
  DataGripIndexingDashboard = `${ROUTE_PREFIX.DataGrip}/indexingDashboard`,
  AIATests = `${ROUTE_PREFIX.AIA}/${TEST_ROUTE}`,
  AIACompletionDashboard = `${ROUTE_PREFIX.AIA}/completion`,
  AIACodeGenerationDashboard = `${ROUTE_PREFIX.AIA}/codeGeneration`,
  AIAChatCodeGenerationDashboard = `${ROUTE_PREFIX.AIA}/chatCodeGeneration`,
  AIANameSuggestionDashboard = `${ROUTE_PREFIX.AIA}/nameSuggestion`,
  AIATestGenerationDashboard = `${ROUTE_PREFIX.AIA}/testGeneration`,
  KMTTests = `${ROUTE_PREFIX.KMT}/unitTests`,
  KMTIntegrationTests = `${ROUTE_PREFIX.KMT}/${TEST_ROUTE}`,
  KMTDashboard = `${ROUTE_PREFIX.KMT}/${DASHBOARD_ROUTE}`,
  DiogenTests = `${ROUTE_PREFIX.Diogen}/${TEST_ROUTE}`,
  ToolboxTests = `${ROUTE_PREFIX.Toolbox}/${TEST_ROUTE}`,
  ReportDegradations = "/degradations/report",
  MetricsDescription = "/metrics/description",
  BisectLauncher = "/bisect/launcher",
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
const STARTUP_NEW_LABEL = "Startup NEW"
const PRODUCT_METRICS_LABEL = "Product Metrics"

const IJ_STARTUP: Product = {
  url: ROUTE_PREFIX.Startup,
  label: "IntelliJ Startup",
  children: [
    {
      url: ROUTE_PREFIX.Startup,
      label: "",
      tabs: [
        {
          url: ROUTES.StartupPulse,
          label: "Pulse",
        },
        {
          url: ROUTES.StartupPulseInstaller,
          label: "Pulse (Installer)",
        },
        {
          url: ROUTES.StartupModuleLoading,
          label: "Module Loading",
        },
        {
          url: ROUTES.StartupGcAndMemory,
          label: "GC and Memory",
        },
        {
          url: ROUTES.StartupProgress,
          label: "Progress Over Time",
        },
        {
          url: ROUTES.StartupExplore,
          label: "Explore",
        },
        {
          url: ROUTES.StartupExploreInstaller,
          label: "Explore (Installer)",
        },
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
        {
          url: ROUTES.IntelliJStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.IntelliJNewStartupDashboard,
          label: STARTUP_NEW_LABEL,
        },
        {
          url: ROUTES.IntelliJProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.IntelliJDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.IntelliJPopupsDashboard,
          label: "Popups",
        },
        {
          url: ROUTES.IntelliJLaggingLatencyDashboard,
          label: "Lagging/Latency",
        },
        {
          url: ROUTES.IntelliJFindUsagesDashboard,
          label: "Find Usages",
        },
        {
          url: ROUTES.IntelliJSEDashboard,
          label: "Search Everywhere",
        },
        {
          url: ROUTES.IntelliJIndexingDashboard,
          label: "Indexes",
        },
        {
          url: ROUTES.IntelliJWSLDashboard,
          label: "WSL",
        },
        {
          url: ROUTES.IntelliJTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.IntelliJDevTests,
          label: "Tests (Dev)",
        },
        {
          url: ROUTES.IntelliJCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
        {
          url: ROUTES.IntelliJCompareModes,
          label: COMPARE_MODES_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJBuildTools,
      label: "Build Tools",
      tabs: [
        {
          url: ROUTES.IntelliJGradleDashboardDev,
          label: "Gradle Import DevServer",
        },
        {
          url: ROUTES.IntelliJMavenDashboardDev,
          label: "Maven Import DevServer",
        },
        {
          url: ROUTES.IntelliJMavenImportersConfiguratorsDashboardDev,
          label: "Maven Importers and Configurators DevServer",
        },
        {
          url: ROUTES.IntelliJJpsDashboardDev,
          label: "JPS Import DevServer",
        },
        {
          url: ROUTES.IntelliJBuildTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.IntelliJBuildTestsDev,
          label: "Tests (DevServer)",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJUltimate,
      label: "Ultimate",
      tabs: [
        {
          url: ROUTES.IntelliJUltimateDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.IntelliJUltimateDashboardOld,
          label: "Dashboard (<=241)",
        },
        {
          url: ROUTES.IntelliJUltimateTests,
          label: TESTS_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJJava,
      label: "Java",
      tabs: [
        {
          url: ROUTES.IntelliJJavaDashboard,
          label: DASHBOARD_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJSharedIndexes,
      label: "Shared Indexes",
      tabs: [
        {
          url: ROUTES.IntelliJSharedIndicesDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.IntelliJSharedIndicesTests,
          label: TESTS_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJIncrementalCompilation,
      label: "Incremental Compilation",
      tabs: [
        {
          url: ROUTES.IntelliJIncrementalCompilationDashboard,
          label: DASHBOARD_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJKotlinK2Performance,
      label: "Performance K2",
      tabs: [
        {
          url: ROUTES.IntelliJK2Dashboard,
          label: DASHBOARD_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJPackageChecker,
      label: "Package Checker",
      tabs: [
        {
          url: ROUTES.IntelliJPackageCheckerDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.IntelliJPackageCheckerTests,
          label: TESTS_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.Vcs,
      label: "VCS",
      tabs: [
        {
          url: ROUTES.VcsIdeaDashboardDev,
          label: "Performance dashboard idea project DevServer",
        },
        {
          url: ROUTES.VcsSpaceDashboardDev,
          label: "Performance dashboard space project DevServer",
        },
        {
          url: ROUTES.VcsStarterDashboardDev,
          label: "Performance dashboard starter project DevServer",
        },
        {
          url: ROUTES.VcsIdeaDashboard,
          label: "Performance dashboard idea project (obsolete)",
        },
        {
          url: ROUTES.VcsSpaceDashboard,
          label: "Performance dashboard space project (obsolete)",
        },
        {
          url: ROUTES.VcsStarterDashboard,
          label: "Performance dashboard starter project (obsolete)",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.EmbeddingSearch,
      label: "Embedding Search",
      tabs: [
        {
          url: ROUTES.IntelliJEmbeddingSearchDashboard,
          label: "Embedding Search",
        },
      ],
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
        {
          url: ROUTES.PhpStormStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.PhpStormProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.PhpStormLLMDashboard,
          label: "LLM Dashboard",
        },
        {
          url: ROUTES.PhpStormIndexingDashboard,
          label: "Indexing Dashboard",
        },
        {
          url: ROUTES.PhpStormInspectionsDashboard,
          label: "Inspections Dashboard",
        },
        {
          url: ROUTES.PhpStormCodeEditingDashboard,
          label: "Code Editing Dashboard",
        },
        {
          url: ROUTES.PhpStormDevTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.PhpStormCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
        {
          url: ROUTES.PhpStormCompareModes,
          label: COMPARE_MODES_LABEL,
        },
        {
          url: ROUTES.PhpStormWithPluginsDashboard,
          label: "Dashboard with Plugins",
        },
        {
          url: ROUTES.PhpStormWithPluginsTests,
          label: "Tests with Plugins",
        },
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
      label: "K1 vs K2",
      tabs: [
        {
          url: ROUTES.KotlinDashboardDev,
          label: DASHBOARD_LABEL + " (dev)",
        },
        {
          url: ROUTES.KotlinUserScenariosDashboardDev,
          label: "User Scenarios(dev)",
        },
        {
          url: ROUTES.KotlinDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.KotlinTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.KotlinTestsDev,
          label: "Tests (dev)",
        },
        {
          url: ROUTES.KotlinK1VsK2Comparison,
          label: "K1 vs. K2",
        },
        {
          url: ROUTES.KotlinK1VsK2ComparisonDev,
          label: "K1 vs. K2 (dev)",
        },
        {
          url: ROUTES.KotlinCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
        {
          url: ROUTES.KotlinCompareBranchesDev,
          label: COMPARE_BRANCHES_LABEL + " (dev)",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.KotlinMemory,
      label: "Memory dashboards",
      tabs: [
        {
          url: ROUTES.KotlinMemoryDashboardDev,
          label: "Memory k1 vs k2 (dev)",
        },
        {
          url: ROUTES.KotlinMemoryDashboard,
          label: "Memory k1 vs k2",
        },
      ],
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
        {
          url: ROUTES.GoLandStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.GoLandProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.GoLandProductMetricsDashboardOld,
          label: PRODUCT_METRICS_LABEL + " (Old)",
        },
        {
          url: ROUTES.GoLandIndexingDashboard,
          label: "Indexing",
        },
        {
          url: ROUTES.GoLandIndexingDashboardOld,
          label: "Indexing (Old)",
        },
        {
          url: ROUTES.GoLandScanningDashboard,
          label: "Scanning",
        },
        {
          url: ROUTES.GoLandScanningDashboardOld,
          label: "Scanning (Old)",
        },
        {
          url: ROUTES.GoLandCompletionDashboard,
          label: "Completion",
        },
        {
          url: ROUTES.GoLandCompletionDashboardOld,
          label: "Completion (Old)",
        },
        {
          url: ROUTES.GoLandInspectionDashboard,
          label: "Inspections",
        },
        {
          url: ROUTES.GoLandInspectionDashboardOld,
          label: "Inspections (Old)",
        },
        {
          url: ROUTES.GoLandDebuggerDashboard,
          label: "Debugger",
        },
        {
          url: ROUTES.GoLandDebuggerDashboardOld,
          label: "Debugger (Old)",
        },
        {
          url: ROUTES.GoLandFindUsagesDashboard,
          label: "Find Usages",
        },
        {
          url: ROUTES.GoLandFindUsagesDashboardOld,
          label: "Find Usages (Old)",
        },
        {
          url: ROUTES.GoLandDFADashboard,
          label: "DFA",
        },
        {
          url: ROUTES.GoLandDFADashboardOld,
          label: "DFA (Old)",
        },
        {
          url: ROUTES.GoLandDistributiveSizeDashboard,
          label: "Distributive Size",
        },
        {
          url: ROUTES.GoLandTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.GoLandTestsOld,
          label: TESTS_LABEL + " (Old)",
        },
        {
          url: ROUTES.GoLandCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
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
        {
          url: ROUTES.RubyStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.RubyMineProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.RubyMineProductMetricsDashboardOld,
          label: PRODUCT_METRICS_LABEL + " (Old)",
        },
        {
          url: ROUTES.RubyMineDashboardOld,
          label: DASHBOARD_LABEL + " (Old)",
        },
        {
          url: ROUTES.RubyMineDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.RubyMineInspectionsDashBoard,
          label: "Inspections",
        },
        {
          url: ROUTES.RubyMineInspectionsDashBoardOld,
          label: "Inspections (Old)",
        },
        {
          url: ROUTES.RubyMineIndexingDashBoard,
          label: "Indexing",
        },
        {
          url: ROUTES.RubyMineIndexingDashBoardOld,
          label: "Indexing (Old)",
        },
        {
          url: ROUTES.RubyMineTests,
          label: TESTS_LABEL + " (Old)",
        },
        {
          url: ROUTES.RubyMineTestsDev,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.RubyMineCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
        {
          url: ROUTES.RubyMineCompareModes,
          label: COMPARE_MODES_LABEL,
        },
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
        {
          url: ROUTES.PyCharmStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.PyCharmProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.PyCharmDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.PyCharmOldDashboard,
          label: DASHBOARD_LABEL + " (Old)",
        },
        {
          url: ROUTES.PyCharmDevTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.PyCharmTests,
          label: TESTS_LABEL + " (Old)",
        },
        {
          url: ROUTES.PyCharmCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
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
        {
          url: ROUTES.WebStormStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.WebStormProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.WebStormProductMetricsDashboardOld,
          label: PRODUCT_METRICS_LABEL + " (Old)",
        },
        {
          url: ROUTES.WebStormDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.WebStormDashboardOld,
          label: DASHBOARD_LABEL + " (Old)",
        },
        {
          url: ROUTES.WebStormDashboardBuiltInVsNEXT,
          label: "Built-in vs NEXT",
        },
        {
          url: ROUTES.WebStormDashboardBuiltInVsNEXTOld,
          label: "Built-in vs NEXT" + " (Old)",
        },
        {
          url: ROUTES.WebStormDashboardDelicateProjects,
          label: "Delicate Projects",
        },
        {
          url: ROUTES.WebStormTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.WebStormTestsOld,
          label: TESTS_LABEL + " (Old)",
        },
        {
          url: ROUTES.WebStormCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
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
        {
          url: ROUTES.RustRoverFirstStartupDashboard,
          label: "RustRover First Startup Dashboard",
        },
        {
          url: ROUTES.RustRoverProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.RustRoverDashboard,
          label: "RustRover Dashboard",
        },
        {
          url: ROUTES.RustTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.RustCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
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
      tabs: [
        {
          url: ROUTES.ScalaTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.ScalaCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
      ],
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
        {
          url: ROUTES.DaCapoDashboard,
          label: "DaCapo",
        },
        {
          url: ROUTES.J2DBenchDashboard,
          label: "J2DBench",
        },
        {
          url: ROUTES.JavaDrawDashboard,
          label: "JavaDraw",
        },
        {
          url: ROUTES.RenderDashboard,
          label: "Render",
        },
        {
          url: ROUTES.SPECjbb2015Dashboard,
          label: "SPECjbb2015",
        },
        {
          url: ROUTES.SwingMarkDashboard,
          label: "SwingMark",
        },
        {
          url: ROUTES.MapBenchDashboard,
          label: "MapBench",
        },
        {
          url: ROUTES.JBRTests,
          label: TESTS_LABEL,
        },
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
        {
          url: ROUTES.FleetStartupDashboard,
          label: "Startup Dashboard",
        },
        {
          url: ROUTES.FleetStartupExplore,
          label: "Startup Explore",
        },
        {
          url: ROUTES.FleetPerfDashboard,
          label: "Performance Dashboard",
        },
        {
          url: ROUTES.FleetPerfStartupComparisonDashboard,
          label: "Startup Comparison Dashboard",
        },
        {
          url: ROUTES.FleetTest,
          label: TESTS_LABEL,
        },
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
      tabs: [
        {
          url: ROUTES.BazelPluginDashboard,
          label: "Bazel Plugin Dashboard",
        },
        {
          url: ROUTES.BazelTest,
          label: TESTS_LABEL,
        },
      ],
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
      tabs: [
        {
          url: ROUTES.QodanaTest,
          label: TESTS_LABEL,
        },
      ],
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
        {
          url: ROUTES.ClionClassicStartupDashboard,
          label: "CLion Classic Startup",
        },
        {
          url: ROUTES.ClionNovaStartupDashboard,
          label: "CLion Nova Startup",
        },
        {
          url: ROUTES.ClionProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.ClionProductMetricsDashboardOld,
          label: `${PRODUCT_METRICS_LABEL} (Old)`,
        },
        {
          url: ROUTES.ClionPerfDashboard,
          label: "Performance",
        },
        {
          url: ROUTES.ClionPerfDashboardOld,
          label: "Performance (Old)",
        },
        {
          url: ROUTES.ClionDetailedPerfDashboard,
          label: "Detailed Performance",
        },
        {
          url: ROUTES.ClionDetailedPerfDashboardOld,
          label: "Detailed Performance (Old)",
        },
        {
          url: ROUTES.ClionMemoryDashboard,
          label: "Memory",
        },
        {
          url: ROUTES.ClionMemoryDashboardOld,
          label: "Memory (Old)",
        },
        {
          url: ROUTES.ClionProjectModelDashboard,
          label: "Project Model",
        },
        {
          url: ROUTES.ClionProjectModelDashboardOld,
          label: "Project Model (Old)",
        },
        {
          url: ROUTES.ClionTest,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.ClionTestOld,
          label: TESTS_LABEL + "(Old)",
        },
        {
          url: ROUTES.ClionCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
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
        {
          url: ROUTES.DataGripStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.DataGripProductMetricsDashboard,
          label: PRODUCT_METRICS_LABEL,
        },
        {
          url: ROUTES.DataGripIndexingDashboard,
          label: "Indexing",
        },
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
      tabs: [
        {
          url: ROUTES.PerfUnitTests,
          label: "Tests",
        },
      ],
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
        {
          url: ROUTES.IJentBenchmarksDashboard,
          label: "Benchmarks Dashboard",
        },
        {
          url: ROUTES.IJentPerfTestsDashboard,
          label: "Performance Dashboard",
        },
        {
          url: ROUTES.IJentBenchmarks,
          label: "Benchmarks",
        },
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
        {
          url: ROUTES.AIAssistantApiTests,
          label: "AI Assistant Api Tests",
        },
        {
          url: ROUTES.AIAssistantTestGeneration,
          label: "Test generation",
        },
        {
          url: ROUTES.LLMDevTests,
          label: "AIA Dashboard",
        },
        {
          url: ROUTES.AIAPrivacyDashboard,
          label: "AIA Privacy Dashboard",
        },
        {
          url: ROUTES.MLDevTests,
          label: "ML Tests on dev-server/fast-installer",
        },
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
        {
          url: ROUTES.AIATests,
          label: "All",
        },
        {
          url: ROUTES.AIACompletionDashboard,
          label: "Completion",
        },
        {
          url: ROUTES.AIACodeGenerationDashboard,
          label: "Code Generation",
        },
        {
          url: ROUTES.AIAChatCodeGenerationDashboard,
          label: "Chat Code Generation",
        },
        {
          url: ROUTES.AIANameSuggestionDashboard,
          label: "Name Suggestion",
        },
        {
          url: ROUTES.AIATestGenerationDashboard,
          label: "Test Generation",
        },
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
      tabs: [
        {
          url: ROUTES.KMTDashboard,
          label: "Dashboard",
        },
        {
          url: ROUTES.KMTIntegrationTests,
          label: "Integration Tests",
        },
        {
          url: ROUTES.KMTTests,
          label: "Unit Tests",
        },
      ],
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
      tabs: [
        {
          url: ROUTES.DiogenTests,
          label: "All",
        },
      ],
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
      tabs: [
        {
          url: ROUTES.ToolboxTests,
          label: "All",
        },
      ],
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

export function getNewDashboardRoutes(): ParentRouteRecord[] {
  return [
    {
      children: [
        {
          path: ROUTES.StartupPulse,
          component: () => import("./components/startup/IntelliJPulse.vue"),
          meta: { pageTitle: "Pulse" },
        },
        {
          path: ROUTES.StartupPulseInstaller,
          component: () => import("./components/startup/IntelliJPulse.vue"),
          props: {
            withInstaller: true,
          },
          meta: { pageTitle: "Pulse" },
        },
        {
          path: ROUTES.StartupModuleLoading,
          component: () => import("./components/startup/IntelliJModuleLoading.vue"),
          meta: { pageTitle: "Module Loading" },
        },
        {
          path: ROUTES.StartupGcAndMemory,
          component: () => import("./components/startup/GcAndMemory.vue"),
          meta: { pageTitle: "GC and Memory" },
        },
        {
          path: ROUTES.StartupProgress,
          component: () => import("./components/startup/IntelliJProgressOverTime.vue"),
          meta: { pageTitle: "Progress Over Time" },
        },
        {
          path: ROUTES.StartupExplore,
          component: () => import("./components/startup/IntelliJExplore.vue"),
          props: {
            withInstaller: false,
          },
          meta: { pageTitle: "Explore" },
        },
        {
          path: ROUTES.StartupExploreInstaller,
          component: () => import("./components/startup/IntelliJExplore.vue"),
          props: {
            withInstaller: true,
          },
          meta: { pageTitle: "Explore (Installer)" },
        },
        {
          path: ROUTES.IntelliJStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "IU",
            defaultProject: "idea",
          },
          meta: { pageTitle: "IDEA Startup dashboard" },
        },
        {
          path: ROUTES.IntelliJNewStartupDashboard,
          component: () => import("./components/intelliJ/StartupMetricsDashboard.vue"),
          meta: { pageTitle: "IDEA NEW Startup dashboard" },
        },
        {
          path: ROUTES.IntelliJProductMetricsDashboard,
          component: () => import("./components/intelliJ/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "IDEA product metrics" },
        },
        {
          path: ROUTES.IntelliJDashboard,
          component: () => import("./components/intelliJ/PerformanceDashboard.vue"),
          meta: { pageTitle: "IntelliJ Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJPopupsDashboard,
          component: () => import("./components/intelliJ/PerformancePopupsDashboard.vue"),
          meta: { pageTitle: "IntelliJ Popups Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJLaggingLatencyDashboard,
          component: () => import("./components/intelliJ/PerformanceLaggingLatencyDashboard.vue"),
          meta: { pageTitle: "IntelliJ Lagging/Latency Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJIndexingDashboard,
          component: () => import("./components/intelliJ/IndexingDashboard.vue"),
          meta: { pageTitle: "IntelliJ Indexing Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJIncrementalCompilationDashboard,
          component: () => import("./components/intelliJ/IncrementalCompilationDashboard.vue"),
          meta: { pageTitle: "IntelliJ Incremental Compilation dashboard" },
        },
        {
          path: ROUTES.IntelliJK2Dashboard,
          component: () => import("./components/intelliJ/PerformanceK2Dashboard.vue"),
          meta: { pageTitle: "IntelliJ Performance K2 dashboard" },
        },
        {
          path: ROUTES.IntelliJGradleDashboardDev,
          component: () => import("./components/intelliJ/build-tools/gradle/GradleImportPerformanceDashboardDevServer.vue"),
          meta: { pageTitle: "Gradle Import DevServer" },
        },
        {
          path: ROUTES.IntelliJMavenDashboardDev,
          component: () => import("./components/intelliJ/build-tools/maven/MavenImportPerformanceDashboardDevServer.vue"),
          meta: { pageTitle: "Maven Import dashboard DevServer" },
        },
        {
          path: ROUTES.IntelliJMavenImportersConfiguratorsDashboardDev,
          component: () => import("./components/intelliJ/build-tools/maven/MavenImportersAndConfiguratorsPerformanceDashboardDevServer.vue"),
          meta: { pageTitle: "Maven Importers And Configurators dashboard DevServer" },
        },
        {
          path: ROUTES.IntelliJJpsDashboardDev,
          component: () => import("./components/intelliJ/build-tools/jps/JpsImportPerformanceDashboardDevServer.vue"),
          meta: { pageTitle: "JPS Import dashboard DevServer" },
        },
        {
          path: ROUTES.IntelliJUltimateDashboard,
          component: () => import("./components/intelliJ/UltimateProjectsDashboard.vue"),
          meta: { pageTitle: "Ultimate Projects" },
        },
        {
          path: ROUTES.IntelliJJavaDashboard,
          component: () => import("./components/intelliJ/JavaProjectsDashboard.vue"),
          meta: { pageTitle: "Java Projects" },
        },
        {
          path: ROUTES.IntelliJPackageCheckerDashboard,
          component: () => import("./components/intelliJ/PackageCheckerDashboard.vue"),
          meta: { pageTitle: "Package Checker" },
        },
        {
          path: ROUTES.IntelliJFindUsagesDashboard,
          component: () => import("./components/intelliJ/PerformanceFindUsagesDashboard.vue"),
          meta: { pageTitle: "Find Usages IntelliJ Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJSEDashboard,
          component: () => import("./components/intelliJ/PerformanceSEDashboard.vue"),
          meta: { pageTitle: "Search Everywhere IntelliJ Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesDashboard,
          component: () => import("./components/intelliJ/SharedIndexesDashboard.vue"),
          meta: { pageTitle: "Shared Indexes Performance Dashboard" },
        },
        {
          path: ROUTES.IntelliJWSLDashboard,
          component: () => import("./components/intelliJ/PerformanceWSLDashboard.vue"),
          meta: { pageTitle: "WSL Performance Dashboard" },
        },
        {
          path: ROUTES.IntelliJEmbeddingSearchDashboard,
          component: () => import("./components/intelliJ/embeddingSearch/Dashboard.vue"),
          meta: { pageTitle: "IntelliJ performance tests for embedding search" },
        },
        {
          path: `${ROUTE_PREFIX.IntelliJ}/:subproject?/${TEST_ROUTE}`,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "idea",
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "IntelliJ Performance tests" },
        },
        {
          path: `${ROUTE_PREFIX.IntelliJ}/:subproject?/${DEV_TEST_ROUTE}`,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "idea",
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
            withInstaller: false,
          },
          meta: { pageTitle: "IntelliJ Integration Performance Tests On DevServer" },
        },
        {
          path: ROUTES.IntelliJCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "idea",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.IntelliJCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfintDev",
            table: "idea",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.IntelliJCompareModes,
          component: () => import("./components/common/compare/CompareModes.vue"),
          props: {
            dbName: "perfintDev",
            table: "idea",
          },
          meta: { pageTitle: COMPARE_MODES_LABEL },
        },
        {
          path: ROUTES.PhpStormStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "PS",
            defaultProject: "stitcher with composer",
          },
          meta: { pageTitle: "PhpStorm Startup dashboard" },
        },
        {
          path: ROUTES.PhpStormProductMetricsDashboard,
          component: () => import("./components/phpstorm/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "PhpStorm product metrics" },
        },
        {
          path: ROUTES.PhpStormLLMDashboard,
          component: () => import("./components/phpstorm/MLDashboard.vue"),
          meta: { pageTitle: "PhpStorm LLM Performance dashboard" },
        },
        {
          path: ROUTES.PhpStormIndexingDashboard,
          component: () => import("./components/phpstorm/IndexingDashboard.vue"),
          meta: { pageTitle: "PhpStorm Indexing Dashboard" },
        },
        {
          path: ROUTES.PhpStormInspectionsDashboard,
          component: () => import("./components/phpstorm/InspectionsDashboard.vue"),
          meta: { pageTitle: "PhpStorm Inspections Dashboard" },
        },
        {
          path: ROUTES.PhpStormCodeEditingDashboard,
          component: () => import("./components/phpstorm/CodeEditingDashboard.vue"),
          meta: { pageTitle: "PhpStorm Code Editing Dashboard" },
        },
        {
          path: ROUTES.PhpStormWithPluginsDashboard,
          component: () => import("./components/phpstorm/PerformanceDashboardWithPlugins.vue"),
          meta: { pageTitle: "PhpStorm With Plugins Performance dashboard" },
        },
        {
          path: ROUTES.PhpStormWithPluginsTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "phpstormWithPlugins",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "PhpStorm Performance tests with plugins" },
        },
        {
          path: ROUTES.PhpStormTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "PhpStorm Performance tests" },
        },
        {
          path: ROUTES.PhpStormDevTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "phpstorm",
            initialMachine: "linux-blade-hetzner",
            withInstaller: false,
          },
          meta: { pageTitle: "PhpStorm Performance tests" },
        },
        {
          path: ROUTES.PhpStormCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.PhpStormCompareModes,
          component: () => import("./components/common/compare/CompareModes.vue"),
          props: {
            dbName: "perfintDev",
            table: "phpstorm",
          },
          meta: { pageTitle: COMPARE_MODES_LABEL },
        },
        {
          path: ROUTES.GoLandInspectionDashboard,
          component: () => import("./components/goland/InspectionsDashboard.vue"),
          meta: { pageTitle: "GoLand Inspections dashboard" },
        },
        {
          path: ROUTES.GoLandInspectionDashboardOld,
          component: () => import("./components/goland/InspectionsDashboardOld.vue"),
          meta: { pageTitle: "GoLand Inspections dashboard" },
        },
        {
          path: ROUTES.GoLandStartupDashboard,
          component: () => import("./components/goland/StartupDashboard.vue"),
          props: {
            product: "GO",
            defaultProject: "kratos",
          },
          meta: { pageTitle: "GoLand Startup dashboard" },
        },
        {
          path: ROUTES.GoLandProductMetricsDashboard,
          component: () => import("./components/goland/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "GoLand product metrics" },
        },
        {
          path: ROUTES.GoLandProductMetricsDashboardOld,
          component: () => import("./components/goland/ProductMetricsDashboardOld.vue"),
          meta: { pageTitle: "GoLand product metrics" },
        },
        {
          path: ROUTES.GoLandIndexingDashboard,
          component: () => import("./components/goland/IndexingDashboard.vue"),
          meta: { pageTitle: "GoLand Indexing dashboard" },
        },
        {
          path: ROUTES.GoLandIndexingDashboardOld,
          component: () => import("./components/goland/IndexingDashboardOld.vue"),
          meta: { pageTitle: "GoLand Indexing dashboard" },
        },
        {
          path: ROUTES.GoLandScanningDashboard,
          component: () => import("./components/goland/ScanningDashboard.vue"),
          meta: { pageTitle: "GoLand Scanning dashboard" },
        },
        {
          path: ROUTES.GoLandScanningDashboardOld,
          component: () => import("./components/goland/ScanningDashboardOld.vue"),
          meta: { pageTitle: "GoLand Scanning dashboard" },
        },
        {
          path: ROUTES.GoLandCompletionDashboard,
          component: () => import("./components/goland/CompletionDashboard.vue"),
          meta: { pageTitle: "GoLand Completion dashboard" },
        },
        {
          path: ROUTES.GoLandCompletionDashboardOld,
          component: () => import("./components/goland/CompletionDashboardOld.vue"),
          meta: { pageTitle: "GoLand Completion dashboard" },
        },
        {
          path: ROUTES.GoLandDebuggerDashboard,
          component: () => import("./components/goland/DebuggerDashboard.vue"),
          meta: { pageTitle: "GoLand Debugger dashboard" },
        },
        {
          path: ROUTES.GoLandDebuggerDashboardOld,
          component: () => import("./components/goland/DebuggerDashboardOld.vue"),
          meta: { pageTitle: "GoLand Debugger dashboard" },
        },
        {
          path: ROUTES.GoLandFindUsagesDashboard,
          component: () => import("./components/goland/FindUsagesDashboard.vue"),
          meta: { pageTitle: "GoLand Find Usages dashboard" },
        },
        {
          path: ROUTES.GoLandFindUsagesDashboardOld,
          component: () => import("./components/goland/FindUsagesDashboardOld.vue"),
          meta: { pageTitle: "GoLand Find Usages dashboard" },
        },
        {
          path: ROUTES.GoLandDFADashboard,
          component: () => import("./components/goland/DataFlowAnalysisDashboard.vue"),
          meta: { pageTitle: "GoLand DFA dashboard" },
        },
        {
          path: ROUTES.GoLandDFADashboardOld,
          component: () => import("./components/goland/DataFlowAnalysisDashboardOld.vue"),
          meta: { pageTitle: "GoLand DFA dashboard" },
        },
        {
          path: ROUTES.GoLandDistributiveSizeDashboard,
          component: () => import("./components/goland/DistributionSizeDashboard.vue"),
          meta: { pageTitle: "GoLand Distribuvite Size dashboard" },
        },
        {
          path: ROUTES.GoLandTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "goland",
            withInstaller: false,
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "GoLand Performance tests" },
        },
        {
          path: ROUTES.GoLandTestsOld,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "goland",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "GoLand Performance tests" },
        },
        {
          path: ROUTES.GoLandCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfintDev",
            table: "goland",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.GoLandCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfintDev",
            table: "goland",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.PyCharmStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "PY",
            defaultProject: "tensorflow",
          },
          meta: { pageTitle: "PyCharm Startup dashboard" },
        },
        {
          path: ROUTES.PyCharmProductMetricsDashboard,
          component: () => import("./components/pycharm/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "PyCharm product metrics" },
        },
        {
          path: ROUTES.PyCharmDashboard,
          component: () => import("./components/pycharm/PerformanceDashboard.vue"),
          meta: { pageTitle: "PyCharm Performance dashboard" },
        },
        {
          path: ROUTES.PyCharmOldDashboard,
          component: () => import("./components/pycharm/PerformanceDashboardOld.vue"),
          meta: { pageTitle: "PyCharm Performance dashboard" },
        },
        {
          path: ROUTES.PyCharmTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "pycharm",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "PyCharm Performance tests" },
        },
        {
          path: ROUTES.PyCharmDevTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "pycharm",
            initialMachine: "linux-blade-hetzner",
            withInstaller: false,
          },
          meta: { pageTitle: "PyCharm Performance tests" },
        },
        {
          path: ROUTES.PyCharmCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "pycharm",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.PyCharmCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfintDev",
            table: "pycharm",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.WebStormStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "WS",
            defaultProject: "angular",
          },
          meta: { pageTitle: "WebStorm Startup dashboard" },
        },
        {
          path: ROUTES.WebStormProductMetricsDashboard,
          component: () => import("./components/webstorm/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "WebStorm product metrics" },
        },
        {
          path: ROUTES.WebStormProductMetricsDashboardOld,
          component: () => import("./components/webstorm/ProductMetricsDashboardOld.vue"),
          meta: { pageTitle: "WebStorm product metrics (Old)" },
        },
        {
          path: ROUTES.WebStormDashboard,
          component: () => import("./components/webstorm/PerformanceDashboard.vue"),
          meta: { pageTitle: "WebStorm Performance dashboard" },
        },
        {
          path: ROUTES.WebStormDashboardOld,
          component: () => import("./components/webstorm/PerformanceDashboardOld.vue"),
          meta: { pageTitle: "WebStorm Performance dashboard (Old)" },
        },
        {
          path: ROUTES.WebStormTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            withInstaller: false,
            table: "webstorm",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "WebStorm Performance tests" },
        },
        {
          path: ROUTES.WebStormTestsOld,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "webstorm",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "WebStorm Performance tests" },
        },
        {
          path: ROUTES.WebStormDashboardBuiltInVsNEXT,
          component: () => import("./components/webstorm/PerformanceDashboardBuiltInVsNEXT.vue"),
          meta: { pageTitle: "Built-in vs NEXT" },
        },
        {
          path: ROUTES.WebStormDashboardBuiltInVsNEXTOld,
          component: () => import("./components/webstorm/PerformanceDashboardBuiltInVsNEXTOld.vue"),
          meta: { pageTitle: "Built-in vs NEXT (Old)" },
        },
        {
          path: ROUTES.WebStormDashboardDelicateProjects,
          component: () => import("./components/webstorm/PerformanceDashboardDelicateProjects.vue"),
          meta: { pageTitle: "Delicate Projects" },
        },
        {
          path: ROUTES.WebStormCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfintDev",
            table: "webstorm",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.WebStormCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfintDev",
            table: "webstorm",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.RubyStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "RM",
            defaultProject: "diaspora",
          },
          meta: { pageTitle: "Ruby Startup dashboard" },
        },
        {
          path: ROUTES.RubyMineProductMetricsDashboard,
          component: () => import("./components/rubymine/ProductMetricsDevDashboard.vue"),
          meta: { pageTitle: "RubyMine product metrics" },
        },
        {
          path: ROUTES.RubyMineProductMetricsDashboardOld,
          component: () => import("./components/rubymine/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "RubyMine product metrics" },
        },
        {
          path: ROUTES.RubyMineDashboard,
          component: () => import("./components/rubymine/PerformanceDevDashboard.vue"),
          meta: { pageTitle: "RubyMine Performance Dashboard" },
        },
        {
          path: ROUTES.RubyMineDashboardOld,
          component: () => import("./components/rubymine/PerformanceDashboard.vue"),
          meta: { pageTitle: "RubyMine Performance Dashboard" },
        },
        {
          path: ROUTES.RubyMineInspectionsDashBoard,
          component: () => import("./components/rubymine/InspectionsDevDashboard.vue"),
          meta: { pageTitle: "RubyMine Inspections Dashboard" },
        },
        {
          path: ROUTES.RubyMineInspectionsDashBoardOld,
          component: () => import("./components/rubymine/InspectionsDashboard.vue"),
          meta: { pageTitle: "RubyMine Inspections Dashboard" },
        },
        {
          path: ROUTES.RubyMineIndexingDashBoard,
          component: () => import("./components/rubymine/IndexingDevDashboard.vue"),
          meta: { pageTitle: "RubyMine Indexing Dashboard" },
        },
        {
          path: ROUTES.RubyMineIndexingDashBoardOld,
          component: () => import("./components/rubymine/IndexingDashboard.vue"),
          meta: { pageTitle: "RubyMine Indexing Dashboard" },
        },
        {
          path: ROUTES.RubyMineTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "ruby",
            initialMachine: "Linux Munich i7-3770, 32 Gb",
          },
          meta: { pageTitle: "RubyMine Performance tests" },
        },
        {
          path: ROUTES.RubyMineTestsDev,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "ruby",
            initialMachine: "Linux Munich i7-3770, 32 Gb",
            withInstaller: false,
          },
          meta: { pageTitle: "RubyMine Performance tests" },
        },
        {
          path: ROUTES.RubyMineCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "ruby",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.RubyMineCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfintDev",
            table: "ruby",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.RubyMineCompareModes,
          component: () => import("./components/common/compare/CompareModes.vue"),
          props: {
            dbName: "perfintDev",
            table: "ruby",
          },
          meta: { pageTitle: COMPARE_MODES_LABEL },
        },
        {
          path: ROUTES.RustCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
            table: "ruby",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },

        {
          path: ROUTES.KotlinTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "Kotlin Performance tests explore" },
        },
        {
          path: ROUTES.KotlinTestsDev,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "kotlin",
            initialMachine: "linux-blade-hetzner",
            withInstaller: false,
          },
          meta: { pageTitle: "Kotlin Performance tests explore (dev/fast installer)" },
        },
        {
          path: ROUTES.KotlinDashboard,
          component: () => import("./components/kotlin/PerformanceDashboard.vue"),
          meta: { pageTitle: "Kotlin Performance dashboard" },
        },
        {
          path: ROUTES.KotlinDashboardDev,
          component: () => import("./components/kotlin/dev/PerformanceDashboard.vue"),
          meta: { pageTitle: "Kotlin Performance dashboard (dev)" },
        },
        {
          path: ROUTES.KotlinUserScenariosDashboardDev,
          component: () => import("./components/kotlin/dev/UserScenariosDashboard.vue"),
          meta: { pageTitle: "User scenarios dashboard (dev)" },
        },
        {
          path: ROUTES.KotlinCodeAnalysisDev,
          component: () => import("./components/kotlin/dev/KotlinCodeAnalysisChartsDashboard.vue"),
          meta: { pageTitle: "Code analysis (dev)" },
        },
        {
          path: ROUTES.KotlinCompletionDev,
          component: () => import("./components/kotlin/dev/CompletionDashboard.vue"),
          meta: { pageTitle: "Kotlin completion (dev/fast)" },
        },
        {
          path: ROUTES.KotlinFindUsagesDev,
          component: () => import("./components/kotlin/dev/FindUsagesDashboard.vue"),
          meta: { pageTitle: "Kotlin findUsages (dev/fast)" },
        },
        {
          path: ROUTES.KotlinRefactoringDev,
          component: () => import("./components/kotlin/dev/RefactoringDashboard.vue"),
          meta: { pageTitle: "Kotlin refactoring (dev/fast)" },
        },
        {
          path: ROUTES.KotlinDebuggerDev,
          component: () => import("./components/kotlin/dev/DebuggerDashboard.vue"),
          meta: { pageTitle: "Kotlin debugger (dev/fast)" },
        },
        {
          path: ROUTES.KotlinScriptDev,
          component: () => import("./components/kotlin/dev/ScriptDashboard.vue"),
          meta: { pageTitle: "Kts (dev/fast)" },
        },
        {
          path: ROUTES.KotlinK1VsK2Comparison,
          component: () => import("./components/kotlin/K1VsK2ComparisonDashboard.vue"),
          meta: { pageTitle: "Kotlin K1 vs. K2" },
        },
        {
          path: ROUTES.KotlinK1VsK2ComparisonDev,
          component: () => import("./components/kotlin/dev/K1VsK2ComparisonDevDashboard.vue"),
          meta: { pageTitle: "Kotlin K1 vs. K2 (dev/fast)" },
        },
        {
          path: ROUTES.KotlinCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.KotlinCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin",
            metricsNames: KOTLIN_MAIN_METRICS,
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.KotlinCompareBranchesDev,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfintDev",
            table: "kotlin",
            metricsNames: KOTLIN_MAIN_METRICS,
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL + "(dev/fast)" },
        },
        {
          path: ROUTES.KotlinMemoryDashboard,
          component: () => import("./components/kotlin/MemoryPerformanceDashboard.vue"),
          meta: { pageTitle: "Memory" },
        },
        {
          path: ROUTES.KotlinMemoryDashboardDev,
          component: () => import("./components/kotlin/dev/MemoryPerformanceDashboard.vue"),
          meta: { pageTitle: "Memory (dev)" },
        },
        {
          path: ROUTES.RustRoverProductMetricsDashboard,
          component: () => import("./components/rust/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "RustRover product metrics" },
        },
        {
          path: ROUTES.RustRoverDashboard,
          component: () => import("./components/rust/PerformanceDashboardRustRover.vue"),
          props: {
            releaseConfigurator: eap,
          },
          meta: { pageTitle: "RustRover Performance dashboard" },
        },
        {
          path: ROUTES.RustRoverFirstStartupDashboard,
          component: () => import("./components/rust/PerformanceDashboardRustRoverFirstStartup.vue"),
          meta: { pageTitle: "RustRover First Startup Performance dashboard" },
        },
        {
          path: ROUTES.RustTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
            releaseConfigurator: eap,
          },
          meta: { pageTitle: "Rust Performance tests" },
        },
        {
          path: ROUTES.RustCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.RustCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.ScalaTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "scala",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "Scala Performance tests" },
        },
        {
          path: ROUTES.ScalaCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "scala",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.ScalaCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
            table: "scala",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.JBRTests,
          component: () => import("./components/jbr/PerformanceTests.vue"),
          meta: { pageTitle: "JBR Performance tests" },
        },
        {
          path: ROUTES.MapBenchDashboard,
          component: () => import("./components/jbr/MapBenchDashboard.vue"),
          meta: { pageTitle: "MapBench Dashboard" },
        },
        {
          path: ROUTES.DaCapoDashboard,
          component: () => import("./components/jbr/DaCapoDashboard.vue"),
          meta: { pageTitle: "DaCapo Dashboard" },
        },
        {
          path: ROUTES.J2DBenchDashboard,
          component: () => import("./components/jbr/J2DBenchDashboard.vue"),
          meta: { pageTitle: "J2DBench Dashboard" },
        },
        {
          path: ROUTES.JavaDrawDashboard,
          component: () => import("./components/jbr/JavaDrawDashboard.vue"),
          meta: { pageTitle: "JavaDraw Dashboard" },
        },
        {
          path: ROUTES.RenderDashboard,
          component: () => import("./components/jbr/RenderDashboard.vue"),
          meta: { pageTitle: "Render Dashboard" },
        },
        {
          path: ROUTES.SPECjbb2015Dashboard,
          component: () => import("./components/jbr/SPECjbb2015Dashboard.vue"),
          meta: { pageTitle: "Spec Dashboard" },
        },
        {
          path: ROUTES.SwingMarkDashboard,
          component: () => import("./components/jbr/SwingMarkDashboard.vue"),
          meta: { pageTitle: "SwingMark Dashboard" },
        },
        {
          path: ROUTES.FleetTest,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "fleet",
            table: "measure_new",
            initialMachine: "linux-blade-hetzner",
            withInstaller: false,
          },
          meta: { pageTitle: "Fleet Performance tests" },
        },
        {
          path: ROUTES.FleetPerfDashboard,
          component: () => import("./components/fleet/PerformanceDashboard.vue"),
          meta: { pageTitle: "Fleet Performance dashboard" },
        },
        {
          path: ROUTES.FleetPerfStartupComparisonDashboard,
          component: () => import("./components/fleet/StartupComparisonDashboard.vue"),
          meta: { pageTitle: "Fleet Startup Comparison dashboard" },
        },
        {
          path: ROUTES.FleetStartupDashboard,
          component: () => import("./components/fleet/FleetDashboard.vue"),
          meta: { pageTitle: "Fleet Startup dashboard" },
        },
        {
          path: ROUTES.FleetStartupExplore,
          component: () => import("./components/fleet/FleetExplore.vue"),
          meta: { pageTitle: "Fleet Startup Explore" },
          props: {
            withInstaller: true,
          },
        },
        {
          path: ROUTES.BazelTest,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "bazel",
            table: "report",
            initialMachine: "Linux EC2 m5ad.2xlarge (8 vCPU Xeon, 32 GB)",
            withInstaller: false,
          },
          meta: { pageTitle: "Bazel Performance tests" },
        },
        {
          path: ROUTES.BazelPluginDashboard,
          component: () => import("./components/bazel/BazelPluginDashboard.vue"),
          meta: { pageTitle: "Bazel Plugin Dashboard" },
        },
        {
          path: ROUTES.QodanaTest,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "qodana",
            table: "report",
            initialMachine: "Linux EC2 c5a(d).xlarge (4 vCPU, 8 GB)",
            withInstaller: false,
          },
          meta: { pageTitle: "Qodana tests" },
        },
        {
          path: ROUTES.ClionTest,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "clion",
            withInstallers: false,
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion tests" },
        },
        {
          path: ROUTES.ClionTestOld,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "clion",
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion tests" },
        },
        {
          path: ROUTES.ClionClassicStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "CL",
            defaultProject: "clion/cmake",
          },
          meta: { pageTitle: "CLion Classic Startup dashboard" },
        },
        {
          path: ROUTES.ClionNovaStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "CL",
            defaultProject: "radler/cmake",
            persistentId: "nova-startup-dashboard",
            withInstaller: true,
          },
          meta: { pageTitle: "CLion Nova Startup dashboard" },
        },
        {
          path: ROUTES.ClionProductMetricsDashboard,
          component: () => import("./components/clion/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "CLion product metrics" },
        },
        {
          path: ROUTES.ClionProductMetricsDashboard,
          component: () => import("./components/clion/ProductMetricsDashboardOld.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion product metrics" },
        },
        {
          path: ROUTES.ClionPerfDashboard,
          component: () => import("./components/clion/PerformanceDashboard.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion dashboard" },
        },
        {
          path: ROUTES.ClionPerfDashboardOld,
          component: () => import("./components/clion/PerformanceDashboardOld.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion dashboard" },
        },
        {
          path: ROUTES.ClionDetailedPerfDashboard,
          component: () => import("./components/clion/DetailedPerformanceDashboard.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion Detailed Performance dashboard" },
        },
        {
          path: ROUTES.ClionDetailedPerfDashboardOld,
          component: () => import("./components/clion/DetailedPerformanceDashboardOld.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion Detailed Performance dashboard" },
        },
        {
          path: ROUTES.ClionMemoryDashboard,
          component: () => import("./components/clion/MemoryDashboard.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion Memory dashboard" },
        },
        {
          path: ROUTES.ClionMemoryDashboardOld,
          component: () => import("./components/clion/MemoryDashboardOld.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion Memory dashboard" },
        },
        {
          path: ROUTES.ClionProjectModelDashboard,
          component: () => import("./components/clion/ProjectModelDashboard.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion Project Model dashboard" },
        },
        {
          path: ROUTES.ClionProjectModelDashboardOld,
          component: () => import("./components/clion/ProjectModelDashboardOld.vue"),
          props: {
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "CLion Project Model dashboard" },
        },
        {
          path: ROUTES.ClionCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfintDev",
            table: "clion",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.VcsIdeaDashboardDev,
          component: () => import("./components/vcs/PerformanceDashboardDev.vue"),
          meta: { pageTitle: "Vcs Idea performance dashboard DevServer" },
        },
        {
          path: ROUTES.VcsSpaceDashboardDev,
          component: () => import("./components/vcs/PerformanceSpaceDashboardDev.vue"),
          meta: { pageTitle: "Vcs Space performance dashboard DevServer" },
        },
        {
          path: ROUTES.VcsStarterDashboardDev,
          component: () => import("./components/vcs/PerformanceStarterDashboardDev.vue"),
          meta: { pageTitle: "Vcs Starer performance dashboard DevServer" },
        },
        {
          path: ROUTES.VcsIdeaDashboard,
          component: () => import("./components/vcs/PerformanceDashboard.vue"),
          meta: { pageTitle: "Vcs Idea performance dashboard (obsolete)" },
        },
        {
          path: ROUTES.VcsSpaceDashboard,
          component: () => import("./components/vcs/PerformanceSpaceDashboard.vue"),
          meta: { pageTitle: "Vcs Space performance dashboard (obsolete)" },
        },
        {
          path: ROUTES.VcsStarterDashboard,
          component: () => import("./components/vcs/PerformanceStarterDashboard.vue"),
          meta: { pageTitle: "Vcs Starer performance dashboard (obsolete)" },
        },
        {
          path: ROUTES.PerfUnitTests,
          component: () => import("./components/common/PerformanceUnitTests.vue"),
          props: {
            dbName: "perfUnitTests",
            table: "report",
            initialMachine: "Linux EC2 C5ad.xlarge (4 vCPU AMD EPYC 7002, 8 GB)",
            withInstaller: false,
          },
          meta: { pageTitle: "Perf Unit Tests" },
        },
        {
          path: ROUTES.IJentBenchmarksDashboard,
          component: () => import("./components/ijent/IJentBenchmarskDashboard.vue"),
          meta: { pageTitle: "IJent Benchmarks Dashboard" },
        },
        {
          path: ROUTES.IJentPerfTestsDashboard,
          component: () => import("./components/ijent/IJentPerformanceTestsDashboard.vue"),
          meta: { pageTitle: "IJent Performance Tests Dashboard" },
        },
        {
          path: ROUTES.IJentBenchmarks,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "ijent",
            initialMachine: "windows-azure",
            withInstaller: false,
          },
          meta: { pageTitle: "IJent Benchmarks" },
        },
        {
          path: ROUTES.AIAssistantApiTests,
          component: () => import("./components/ml/dev/AiAssistantApiTests.vue"),
          meta: { pageTitle: "AI API Tests" },
        },
        {
          path: ROUTES.AIAssistantTestGeneration,
          component: () => import("./components/ml/dev/TestGenerationDashboard.vue"),
          meta: { pageTitle: "Test generation" },
        },
        {
          path: ROUTES.LLMDevTests,
          component: () => import("./components/ml/dev/LLMDashboard.vue"),
          meta: { pageTitle: "AIA dashboard" },
        },
        {
          path: ROUTES.AIAPrivacyDashboard,
          component: () => import("./components/ml/dev/AIAPrivacyDashboard.vue"),
          meta: { pageTitle: "AIA Privacy" },
        },
        {
          path: ROUTES.MLDevTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "ml",
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
            withInstaller: false,
          },
          meta: { pageTitle: "ML Tests dev-server" },
        },
        {
          path: ROUTES.DataGripStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "DB",
            defaultProject: "empty project",
          },
          meta: { pageTitle: "DataGrip Startup dashboard" },
        },
        {
          path: ROUTES.DataGripProductMetricsDashboard,
          component: () => import("./components/datagrip/ProductMetricsDashboard.vue"),
          meta: { pageTitle: "DataGrip product metrics" },
        },
        {
          path: ROUTES.DataGripIndexingDashboard,
          component: () => import("./components/datagrip/IndexingDashboard.vue"),
          meta: { pageTitle: "DataGrip Indexing dashboard" },
        },
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
          path: ROUTES.AIATests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "mlEvaluation",
            table: "report",
            withInstaller: false,
            branch: null,
            initialMachine: null,
          },
          meta: { pageTitle: "AIA Tests" },
        },
        {
          path: ROUTES.AIACompletionDashboard,
          component: () => import("./components/aia/AIACompletionDashboard.vue"),
          meta: { pageTitle: "AIA completion dashboard" },
        },
        {
          path: ROUTES.AIACodeGenerationDashboard,
          component: () => import("./components/aia/AIACodeGeneration.vue"),
          meta: { pageTitle: "AIA code generation dashboard" },
        },
        {
          path: ROUTES.AIAChatCodeGenerationDashboard,
          component: () => import("./components/aia/ChatCodeGeneration.vue"),
          meta: { pageTitle: "AIA chat code generation dashboard" },
        },
        {
          path: ROUTES.AIANameSuggestionDashboard,
          component: () => import("./components/aia/AIANameSuggestion.vue"),
          meta: { pageTitle: "AIA name suggestion dashboard" },
        },
        {
          path: ROUTES.AIATestGenerationDashboard,
          component: () => import("./components/aia/AIATestGeneration.vue"),
          meta: { pageTitle: "AIA test generation dashboard" },
        },
        {
          path: ROUTES.KMTTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "swift",
            withInstaller: false,
            branch: null,
            initialMachine: "Mac Cidr Performance",
          },
          meta: { pageTitle: "KMT Unit Tests" },
        },
        {
          path: ROUTES.KMTIntegrationTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "kmt",
            withInstaller: false,
            branch: "master",
            initialMachine: "Mac Cidr Performance",
          },
          meta: { pageTitle: "KMT Integration Tests" },
        },
        {
          path: ROUTES.KMTDashboard,
          component: () => import("./components/kmt/PerformanceDashboard.vue"),
          meta: { pageTitle: "KMT Dashboard" },
        },
        {
          path: ROUTES.DiogenTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "diogen",
            table: "report",
            withInstaller: false,
            branch: "refs/heads/main",
            initialMachine: null,
            withoutAccidents: true,
          },
          meta: { pageTitle: "Diogen" },
        },
        {
          path: ROUTES.ToolboxTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "toolbox",
            table: "report",
            withInstaller: false,
            branch: "refs/heads/main",
            initialMachine: "Linux EC2 M5d.xlarge (4 vCPU Xeon, 16 GB)",
            withoutAccidents: true,
          },
          meta: { pageTitle: "Toolbox" },
        },
      ],
    },
  ]
}
