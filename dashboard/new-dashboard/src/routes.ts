/* eslint-disable @typescript-eslint/prefer-literal-enum-member */

import { ParentRouteRecord } from "./components/common/route"
import { KOTLIN_MAIN_METRICS } from "./components/kotlin/projects"
import { eap } from "./configurators/ReleaseNightlyConfigurator"

const enum ROUTE_PREFIX {
  Startup = "/ij",
  IntelliJ = "/intellij",
  IntelliJBuildTools = "/intellij/buildTools",
  IntelliJSharedIndices = "/intellij/sharedIndexes",
  IntelliJUltimate = "/intellij/ultimate",
  IntelliJPackageChecker = "/intellij/packageChecker",
  IntelliJFus = "/intellij/fus",
  IntelliJExperiments = "/intellij/experiments",
  PhpStorm = "/phpstorm",
  GoLand = "/goland",
  GoLandSharedIndices = "/goland/sharedIndexes",
  RubyMine = "/rubymine",
  Kotlin = "/kotlin",
  KotlinK1VsK2 = "/kotlinK1VsK2",
  KotlinMemory = "/kotlinMemory",
  KotlinMPP = "/kotlinMPP",
  Rust = "/rust",
  Scala = "/scala",
  JBR = "/jbr",
  Fleet = "/fleet",
  PyCharm = "/pycharm",
  PyCharmSharedIndices = "/pycharm/sharedIndexes",
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
}
const TEST_ROUTE = "tests"
const DEV_TEST_ROUTE = "testsDev"
const DASHBOARD_ROUTE = "dashboard"
const STARTUP_ROUTE = "startup"
const COMPARE_ROUTE = "compare"
const COMPARE_BRANCHES_ROUTE = "compareBranches"

enum ROUTES {
  StartupPulse = `${ROUTE_PREFIX.Startup}/pulse`,
  StartupProgress = `${ROUTE_PREFIX.Startup}/progressOverTime`,
  StartupModuleLoading = `${ROUTE_PREFIX.Startup}/moduleLoading`,
  StartupGcAndMemory = `${ROUTE_PREFIX.Startup}/gcAndMemory`,
  StartupExplore = `${ROUTE_PREFIX.Startup}/explore`,
  StartupExploreDev = `${ROUTE_PREFIX.Startup}/exploreDev`,
  StartupReport = `${ROUTE_PREFIX.Startup}/report`,
  IntelliJStartupDashboard = `${ROUTE_PREFIX.IntelliJ}/${STARTUP_ROUTE}`,
  IntelliJDashboard = `${ROUTE_PREFIX.IntelliJ}/${DASHBOARD_ROUTE}`,
  IntelliJIndexingDashboard = `${ROUTE_PREFIX.IntelliJ}/indexingDashboard`,
  IntelliJJBRDashboard = `${ROUTE_PREFIX.IntelliJ}/jbrPerformanceDashboard`,
  IntelliJTinyDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/dashboardTiny`,
  IntelliJIncrementalCompilationDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/incrementalCompilationDashboard`,
  IntelliJScalabilityDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/scalabilityDashboard`,
  IntelliJDevDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboardDev`,
  IntelliJFindUsagesDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboardFindUsages`,
  IntelliJSEDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboardSearchEverywhere`,
  IntelliJEmbeddingSearchDashboard = `${ROUTE_PREFIX.EmbeddingSearch}/dashboard`,
  IntelliJTests = `${ROUTE_PREFIX.IntelliJ}/${TEST_ROUTE}`,
  IntelliJDevTests = `${ROUTE_PREFIX.IntelliJ}/${DEV_TEST_ROUTE}`,
  IntelliJCompare = `${ROUTE_PREFIX.IntelliJ}/${COMPARE_ROUTE}`,
  IntelliJCompareBranches = `${ROUTE_PREFIX.IntelliJ}/${COMPARE_BRANCHES_ROUTE}`,
  IntelliJGradleDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/gradleDashboard`,
  IntelliJGradleDashboardFastInstallers = `${ROUTE_PREFIX.IntelliJBuildTools}/gradleDashboardFastInstallers`,
  IntelliJNewGradleDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/newGradleDashboard`,
  IntelliJNewGradleDashboardFastInstallers = `${ROUTE_PREFIX.IntelliJBuildTools}/newGradleDashboardFastInstallers`,
  IntelliJMavenDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/mavenDashboard`,
  IntelliJMavenDashboardFastInstallers = `${ROUTE_PREFIX.IntelliJBuildTools}/mavenDashboardFastInstallers`,
  IntelliJMavenImportersConfiguratorsDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/mavenImportersConfiguratorsDashboard`,
  IntelliJJpsDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/jpsDashboard`,
  IntelliJBuildTests = `${ROUTE_PREFIX.IntelliJBuildTools}/${TEST_ROUTE}`,
  IntelliJBuildTestsDev = `${ROUTE_PREFIX.IntelliJBuildTools}/${DEV_TEST_ROUTE}`,
  IntelliJUltimateDashboard = `${ROUTE_PREFIX.IntelliJUltimate}/${DASHBOARD_ROUTE}`,
  IntelliJUltimateTests = `${ROUTE_PREFIX.IntelliJUltimate}/${TEST_ROUTE}`,
  IntelliJSharedIndicesIndexingDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesIndexingDashboard`,
  IntelliJSharedIndicesScanningDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesScanningDashboard`,
  IntelliJSharedIndicesFindUsagesDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesFindUsagesDashboard`,
  IntelliJSharedIndicesCompletionDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesCompletionDashboard`,
  IntelliJSharedIndicesFirstCodeAnalysisDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesFirstCodeAnalysisDashboard`,
  IntelliJSharedIndicesNumberOfIndexedFilesDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesIndexedFilesDashboard`,
  IntelliJSharedIndicesNumberOfExtensionsDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesNumberOfExtensionsDashboard`,
  IntelliJSharedIndicesTypingDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesTypingDashboard`,
  IntelliJSharedIndicesDumbModeDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesDumbModeDashboard`,
  IntelliJGCDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/performanceGC`,
  IntelliJSharedIndicesTests = `${ROUTE_PREFIX.IntelliJSharedIndices}/${TEST_ROUTE}`,
  IntelliJPackageCheckerDashboard = `${ROUTE_PREFIX.IntelliJPackageChecker}/${DASHBOARD_ROUTE}`,
  IntelliJPackageCheckerTests = `${ROUTE_PREFIX.IntelliJPackageChecker}/${TEST_ROUTE}`,
  IntelliJFusDashboard = `${ROUTE_PREFIX.IntelliJFus}/${DASHBOARD_ROUTE}`,
  IntelliJFusDevDashboard = `${ROUTE_PREFIX.IntelliJFus}/dashboardDev`,
  IntelliJFusHetznerDashboard = `${ROUTE_PREFIX.IntelliJFus}/dashboardImport`,
  IntelliJFusStartupDashboard = `${ROUTE_PREFIX.IntelliJFus}/dashboardStartup`,
  IntelliJExperimentsGradleSyncDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/dashboardGradleSync`,
  IntelliJExperimentsMonorepoDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/dashboardMonorepo`,
  PhpStormDashboard = `${ROUTE_PREFIX.PhpStorm}/${DASHBOARD_ROUTE}`,
  PhpStormLLMDashboard = `${ROUTE_PREFIX.PhpStorm}/llmDashboard`,
  PhpStormStartupDashboard = `${ROUTE_PREFIX.PhpStorm}/${STARTUP_ROUTE}`,
  PhpStormWithPluginsDashboard = `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests = `${ROUTE_PREFIX.PhpStorm}/${TEST_ROUTE}`,
  PhpStormWithPluginsTests = `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,
  PhpStormCompare = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_ROUTE}`,
  PhpStormCompareBranches = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_BRANCHES_ROUTE}`,
  KotlinDashboard = `${ROUTE_PREFIX.Kotlin}/${DASHBOARD_ROUTE}`,
  KotlinCodeAnalysis = `${ROUTE_PREFIX.Kotlin}/codeAnalysis`,
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
  KotlinMemoryDashboard = `${ROUTE_PREFIX.KotlinMemory}/memoryDashboard`,
  KotlinMPPDashboard = `${ROUTE_PREFIX.KotlinMPP}/dashboard`,
  KotlinCompareBranches = `${ROUTE_PREFIX.Kotlin}/${COMPARE_BRANCHES_ROUTE}`,
  KotlinCompareBranchesDev = `${ROUTE_PREFIX.Kotlin}/${COMPARE_BRANCHES_ROUTE}Dev`,
  GoLandStartupDashboard = `${ROUTE_PREFIX.GoLand}/${STARTUP_ROUTE}`,
  GoLandIndexingDashboard = `${ROUTE_PREFIX.GoLand}/indexingDashboard`,
  GoLandScanningDashboard = `${ROUTE_PREFIX.GoLand}/scanningDashboard`,
  GoLandCompletionDashboard = `${ROUTE_PREFIX.GoLand}/completionDashboard`,
  GoLandInspectionDashboard = `${ROUTE_PREFIX.GoLand}/inspectionsDashboard`,
  GoLandTests = `${ROUTE_PREFIX.GoLand}/${TEST_ROUTE}`,
  GoLandCompare = `${ROUTE_PREFIX.GoLand}/${COMPARE_ROUTE}`,
  GoLandCompareBranches = `${ROUTE_PREFIX.GoLand}/${COMPARE_BRANCHES_ROUTE}`,
  GoLandSharedIndicesIndexingDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesIndexingDashboard`,
  GoLandSharedIndicesScanningDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesScanningDashboard`,
  GoLandSharedIndicesFindUsagesDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesFindUsagesDashboard`,
  GoLandSharedIndicesCompletionDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesCompletionDashboard`,
  GoLandSharedIndicesFirstCodeAnalysisDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesFirstCodeAnalysisDashboard`,
  GoLandSharedIndicesNumberOfIndexedFilesDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesIndexedFilesDashboard`,
  GoLandSharedIndicesNumberOfExtensionsDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesNumberOfExtensionsDashboard`,
  GoLandSharedIndicesTypingDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesTypingDashboard`,
  GoLandSharedIndicesDumbModeDashboard = `${ROUTE_PREFIX.GoLandSharedIndices}/sharedIndexesDumbModeDashboard`,
  PyCharmStartupDashboard = `${ROUTE_PREFIX.PyCharm}/${STARTUP_ROUTE}`,
  PyCharmDashboard = `${ROUTE_PREFIX.PyCharm}/${DASHBOARD_ROUTE}`,
  PyCharmTests = `${ROUTE_PREFIX.PyCharm}/${TEST_ROUTE}`,
  PyCharmCompare = `${ROUTE_PREFIX.PyCharm}/${COMPARE_ROUTE}`,
  PyCharmCompareBranches = `${ROUTE_PREFIX.PyCharm}/${COMPARE_BRANCHES_ROUTE}`,
  PyCharmSharedIndicesIndexingDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesIndexingDashboard`,
  PyCharmSharedIndicesScanningDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesScanningDashboard`,
  PyCharmSharedIndicesFindUsagesDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesFindUsagesDashboard`,
  PyCharmSharedIndicesCompletionDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesCompletionDashboard`,
  PyCharmSharedIndicesFirstCodeAnalysisDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesFirstCodeAnalysisDashboard`,
  PyCharmSharedIndicesNumberOfIndexedFilesDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesIndexedFilesDashboard`,
  PyCharmSharedIndicesNumberOfExtensionsDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesNumberOfExtensionsDashboard`,
  PyCharmSharedIndicesTypingDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesTypingDashboard`,
  PyCharmSharedIndicesDumbModeDashboard = `${ROUTE_PREFIX.PyCharmSharedIndices}/sharedIndexesDumbModeDashboard`,
  WebStormStartupDashboard = `${ROUTE_PREFIX.WebStorm}/${STARTUP_ROUTE}`,
  WebStormDashboard = `${ROUTE_PREFIX.WebStorm}/${DASHBOARD_ROUTE}`,
  WebStormDashboardNEXT = `${ROUTE_PREFIX.WebStorm}/dashboardNext`,
  WebStormTests = `${ROUTE_PREFIX.WebStorm}/${TEST_ROUTE}`,
  WebStormCompare = `${ROUTE_PREFIX.WebStorm}/${COMPARE_ROUTE}`,
  WebStormCompareBranches = `${ROUTE_PREFIX.WebStorm}/${COMPARE_BRANCHES_ROUTE}`,
  RubyStartupDashboard = `${ROUTE_PREFIX.RubyMine}/${STARTUP_ROUTE}`,
  RubyMineDashboard = `${ROUTE_PREFIX.RubyMine}/${DASHBOARD_ROUTE}`,
  RubyMineIndexingDashBoard = `${ROUTE_PREFIX.RubyMine}/indexingDashboard`,
  RubyMineTests = `${ROUTE_PREFIX.RubyMine}/${TEST_ROUTE}`,
  RubyMineCompare = `${ROUTE_PREFIX.RubyMine}/${COMPARE_ROUTE}`,
  RubyMineCompareBranches = `${ROUTE_PREFIX.RubyMine}/${COMPARE_BRANCHES_ROUTE}`,
  RustPluginDashboard = `${ROUTE_PREFIX.Rust}/rustPluginDashboard`,
  RustRoverDashboard = `${ROUTE_PREFIX.Rust}/rustRoverDashboard`,
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
  FleetTest = `${ROUTE_PREFIX.Fleet}/${TEST_ROUTE}`,
  FleetPerfDashboard = `${ROUTE_PREFIX.Fleet}/perfDashboard`,
  FleetStartupDashboard = `${ROUTE_PREFIX.Fleet}/startupDashboard`,
  BazelTest = `${ROUTE_PREFIX.Bazel}/${TEST_ROUTE}`,
  BazelBspDashboard = `${ROUTE_PREFIX.Bazel}/bazelBSPDashboard`,
  IntelliJBspDashboard = `${ROUTE_PREFIX.Bazel}/intellijBSPDashboard`,
  QodanaTest = `${ROUTE_PREFIX.Qodana}/${TEST_ROUTE}`,
  ClionStartupDashboard = `${ROUTE_PREFIX.Clion}/${STARTUP_ROUTE}`,
  ClionTest = `${ROUTE_PREFIX.Clion}/${TEST_ROUTE}`,
  ClionPerfDashboard = `${ROUTE_PREFIX.Clion}/perfDashboard`,
  ClionDetailedPerfDashboard = `${ROUTE_PREFIX.Clion}/detailedPerfDashboard`,
  ClionCompareBranches = `${ROUTE_PREFIX.Clion}/${COMPARE_BRANCHES_ROUTE}`,
  VcsIdeaDashboard = `${ROUTE_PREFIX.Vcs}/idea`,
  VcsSpaceDashboard = `${ROUTE_PREFIX.Vcs}/space`,
  VcsStarterDashboard = `${ROUTE_PREFIX.Vcs}/starter`,
  PerfUnitTests = `${ROUTE_PREFIX.PerfUnit}/${TEST_ROUTE}`,
  IJentBenchmarks = `${ROUTE_PREFIX.IJent}/benchmarks`,
  IJentPerfTests = `${ROUTE_PREFIX.IJent}/performance`,
  MLDevTests = `${ROUTE_PREFIX.ML}/dev/${TEST_ROUTE}`,
  LLMDevTests = `${ROUTE_PREFIX.ML}/dev/llmDashboardDev`,
  FullLineDevTests = `${ROUTE_PREFIX.ML}/dev/fullLineDashboardDev`,
  DataGripStartupDashboard = `${ROUTE_PREFIX.DataGrip}/${STARTUP_ROUTE}`,
  ReportDegradations = "/degradations/report",
  MetricsDescription = "/metrics/description",
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
const DASHBOARD_LABEL = "Dashboard"
const STARTUP_LABEL = "Startup"

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
          url: ROUTES.StartupExploreDev,
          label: "Explore (Dev)",
        },
        {
          url: ROUTES.StartupReport,
          label: "Report",
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
          url: ROUTES.IntelliJDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.IntelliJDevDashboard,
          label: "Dashboard (Fast Installer)",
        },
        {
          url: ROUTES.IntelliJFindUsagesDashboard,
          label: "Find Usages Dashboard",
        },
        {
          url: ROUTES.IntelliJSEDashboard,
          label: "Search Everywhere Dashboard",
        },
        {
          url: ROUTES.IntelliJIndexingDashboard,
          label: "Indexing Dashboard",
        },
        {
          url: ROUTES.IntelliJJBRDashboard,
          label: "JBR Performance tests Dashboard",
        },
        {
          url: ROUTES.IntelliJTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.IntelliJDevTests,
          label: "Tests (Fast Installer)",
        },
        {
          url: ROUTES.IntelliJCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJBuildTools,
      label: "Build Tools",
      tabs: [
        {
          url: ROUTES.IntelliJGradleDashboard,
          label: "Gradle Import",
        },
        {
          url: ROUTES.IntelliJGradleDashboardFastInstallers,
          label: "Gradle Import Fast Installers",
        },
        {
          url: ROUTES.IntelliJNewGradleDashboard,
          label: "New Gradle Import",
        },
        {
          url: ROUTES.IntelliJNewGradleDashboardFastInstallers,
          label: "New Gradle Import Fast Installers",
        },
        {
          url: ROUTES.IntelliJMavenDashboard,
          label: "Maven Import",
        },
        {
          url: ROUTES.IntelliJMavenDashboardFastInstallers,
          label: "Maven Import Fast Installers",
        },
        {
          url: ROUTES.IntelliJMavenImportersConfiguratorsDashboard,
          label: "Maven Importers and Configurators",
        },
        {
          url: ROUTES.IntelliJJpsDashboard,
          label: "JPS Import",
        },
        {
          url: ROUTES.IntelliJBuildTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.IntelliJBuildTestsDev,
          label: "Tests (Fast Installer)",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJSharedIndices,
      label: "Shared Indexes",
      tabs: [
        {
          url: ROUTES.IntelliJSharedIndicesDumbModeDashboard,
          label: "Dumb Mode Time",
        },
        {
          url: ROUTES.IntelliJSharedIndicesIndexingDashboard,
          label: "Indexing",
        },
        {
          url: ROUTES.IntelliJSharedIndicesScanningDashboard,
          label: "Scanning",
        },
        {
          url: ROUTES.IntelliJSharedIndicesFindUsagesDashboard,
          label: "FindUsages",
        },
        {
          url: ROUTES.IntelliJSharedIndicesCompletionDashboard,
          label: "Completion",
        },
        {
          url: ROUTES.IntelliJSharedIndicesFirstCodeAnalysisDashboard,
          label: "Code Analysis",
        },
        {
          url: ROUTES.IntelliJSharedIndicesTypingDashboard,
          label: "Typing",
        },
        {
          url: ROUTES.IntelliJSharedIndicesNumberOfIndexedFilesDashboard,
          label: "Indexed Files",
        },
        {
          url: ROUTES.IntelliJSharedIndicesNumberOfExtensionsDashboard,
          label: "Indexed by Extensions",
        },
        {
          url: ROUTES.IntelliJSharedIndicesTests,
          label: TESTS_LABEL,
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
          url: ROUTES.IntelliJUltimateTests,
          label: TESTS_LABEL,
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
      url: ROUTE_PREFIX.IntelliJFus,
      label: "FUS",
      tabs: [
        {
          url: ROUTES.IntelliJFusDashboard,
          label: "Dashboard 1",
        },
        {
          url: ROUTES.IntelliJFusDevDashboard,
          label: "Dashboard 2",
        },
        {
          url: ROUTES.IntelliJFusHetznerDashboard,
          label: "Dashboard Import",
        },
        {
          url: ROUTES.IntelliJFusStartupDashboard,
          label: "Dashboard Startup",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJExperiments,
      label: "Experiments",
      tabs: [
        {
          url: ROUTES.IntelliJExperimentsGradleSyncDashboard,
          label: "Gradle Sync Smart/Dumb",
        },
        {
          url: ROUTES.IntelliJExperimentsMonorepoDashboard,
          label: "IntelliJ + Dotnet dashboard",
        },
        {
          url: ROUTES.IntelliJTinyDashboard,
          label: "Dashboard (Tiny Agents)",
        },
        {
          url: ROUTES.IntelliJIncrementalCompilationDashboard,
          label: "Incremental Compilation",
        },
        {
          url: ROUTES.IntelliJGCDashboard,
          label: "Garbage Collectors",
        },
        {
          url: ROUTES.IntelliJScalabilityDashboard,
          label: "Scalability",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.Vcs,
      label: "VCS",
      tabs: [
        {
          url: ROUTES.VcsIdeaDashboard,
          label: "Performance dashboard idea project",
        },
        {
          url: ROUTES.VcsSpaceDashboard,
          label: "Performance dashboard space project",
        },
        {
          url: ROUTES.VcsStarterDashboard,
          label: "Performance dashboard starter project",
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
          url: ROUTES.PhpStormDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.PhpStormLLMDashboard,
          label: "LLM Dashboard",
        },
        {
          url: ROUTES.PhpStormTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.PhpStormCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
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
      url: ROUTE_PREFIX.KotlinK1VsK2,
      label: "K1 vs K2",
      tabs: [
        {
          url: ROUTES.KotlinDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.KotlinCodeAnalysis,
          label: "Code Analysis",
        },
        {
          url: ROUTES.KotlinTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.KotlinTestsDev,
          label: "Explore (dev/fast installer)",
        },
        {
          url: ROUTES.KotlinCompletionDev,
          label: "Completion (dev)",
        },
        {
          url: ROUTES.KotlinCodeAnalysisDev,
          label: "Code analysis (dev)",
        },
        {
          url: ROUTES.KotlinHighlightingDev,
          label: "Highlighting (dev)",
        },
        {
          url: ROUTES.KotlinFindUsagesDev,
          label: "FindUsages (dev)",
        },
        {
          url: ROUTES.KotlinRefactoringDev,
          label: "Refactoring (dev)",
        },
        {
          url: ROUTES.KotlinDebuggerDev,
          label: "Debugger (dev)",
        },
        {
          url: ROUTES.KotlinScriptDev,
          label: "Kts (dev)",
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
          label: COMPARE_BRANCHES_LABEL + "(dev)",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.KotlinMPP,
      label: "MPP projects",
      tabs: [
        {
          url: ROUTES.KotlinMPPDashboard,
          label: "k1 vs k2",
        },
      ],
    },
    {
      url: ROUTE_PREFIX.KotlinMemory,
      label: "Memory dashboards",
      tabs: [
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
          url: ROUTES.GoLandIndexingDashboard,
          label: "Indexing",
        },
        {
          url: ROUTES.GoLandScanningDashboard,
          label: "Scanning",
        },
        {
          url: ROUTES.GoLandCompletionDashboard,
          label: "Completion",
        },
        {
          url: ROUTES.GoLandInspectionDashboard,
          label: "Inspections & Analyzes",
        },
        {
          url: ROUTES.GoLandTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.GoLandCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.GoLandSharedIndices,
      label: "Shared Indexes",
      tabs: [
        {
          url: ROUTES.GoLandSharedIndicesDumbModeDashboard,
          label: "Dumb Mode Time",
        },
        {
          url: ROUTES.GoLandSharedIndicesIndexingDashboard,
          label: "Indexing",
        },
        {
          url: ROUTES.GoLandSharedIndicesScanningDashboard,
          label: "Scanning",
        },
        {
          url: ROUTES.GoLandSharedIndicesFindUsagesDashboard,
          label: "FindUsages",
        },
        {
          url: ROUTES.GoLandSharedIndicesCompletionDashboard,
          label: "Completion",
        },
        {
          url: ROUTES.GoLandSharedIndicesFirstCodeAnalysisDashboard,
          label: "Code Analysis",
        },
        {
          url: ROUTES.GoLandSharedIndicesTypingDashboard,
          label: "Typing",
        },
        {
          url: ROUTES.GoLandSharedIndicesNumberOfIndexedFilesDashboard,
          label: "Indexed Files",
        },
        {
          url: ROUTES.GoLandSharedIndicesNumberOfExtensionsDashboard,
          label: "Indexed by Extensions",
        },
        {
          url: ROUTES.GoLandTests,
          label: TESTS_LABEL,
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
          url: ROUTES.RubyMineDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.RubyMineIndexingDashBoard,
          label: "Indexing",
        },
        {
          url: ROUTES.RubyMineTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.RubyMineCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
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
          url: ROUTES.PyCharmDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.PyCharmTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.PyCharmCompareBranches,
          label: COMPARE_BRANCHES_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.PyCharmSharedIndices,
      label: "Shared Indexes",
      tabs: [
        {
          url: ROUTES.PyCharmSharedIndicesDumbModeDashboard,
          label: "Dumb Mode Time",
        },
        {
          url: ROUTES.PyCharmSharedIndicesIndexingDashboard,
          label: "Indexing",
        },
        {
          url: ROUTES.PyCharmSharedIndicesScanningDashboard,
          label: "Scanning",
        },
        {
          url: ROUTES.PyCharmSharedIndicesFindUsagesDashboard,
          label: "FindUsages",
        },
        {
          url: ROUTES.PyCharmSharedIndicesCompletionDashboard,
          label: "Completion",
        },
        {
          url: ROUTES.PyCharmSharedIndicesFirstCodeAnalysisDashboard,
          label: "Code Analysis",
        },
        {
          url: ROUTES.PyCharmSharedIndicesTypingDashboard,
          label: "Typing",
        },
        {
          url: ROUTES.PyCharmSharedIndicesNumberOfIndexedFilesDashboard,
          label: "Indexed Files",
        },
        {
          url: ROUTES.PyCharmSharedIndicesNumberOfExtensionsDashboard,
          label: "Indexed by Extensions",
        },
        {
          url: ROUTES.PyCharmTests,
          label: TESTS_LABEL,
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
          url: ROUTES.WebStormDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.WebStormDashboardNEXT,
          label: "WebStorm NEXT",
        },
        {
          url: ROUTES.WebStormTests,
          label: TESTS_LABEL,
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
          url: ROUTES.RustPluginDashboard,
          label: "Rust Plugin Dashboard",
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
          url: ROUTES.FleetPerfDashboard,
          label: "Performance Dashboard",
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
          url: ROUTES.BazelBspDashboard,
          label: "Bazel BSP Dashboard",
        },
        {
          url: ROUTES.IntelliJBspDashboard,
          label: "IntelliJ BSP Dashboard",
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
          url: ROUTES.ClionStartupDashboard,
          label: STARTUP_LABEL,
        },
        {
          url: ROUTES.ClionPerfDashboard,
          label: "Performance",
        },
        {
          url: ROUTES.ClionDetailedPerfDashboard,
          label: "Detailed Performance",
        },
        {
          url: ROUTES.ClionTest,
          label: TESTS_LABEL,
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
          url: ROUTES.IJentBenchmarks,
          label: "Benchmarks",
        },
        {
          url: ROUTES.IJentPerfTests,
          label: "Tests",
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
          url: ROUTES.LLMDevTests,
          label: "LLM test dashboard dev-server/fast-installer",
        },
        {
          url: ROUTES.FullLineDevTests,
          label: "FullLine test dashboard dev-server/fast-installer",
        },
        {
          url: ROUTES.MLDevTests,
          label: "ML Tests on dev-server/fast-installer",
        },
      ],
    },
  ],
}

export const PRODUCTS = [
  IDEA,
  PHPSTORM,
  KOTLIN,
  GOLAND,
  RUBYMINE,
  PYCHARM,
  WEBSTORM,
  CLION,
  DATAGRIP,
  RUST,
  FLEET,
  BAZEL,
  QODANA,
  IJ_STARTUP,
  SCALA,
  JBR,
  PERF_UNIT,
  IJENT,
  ML_TESTS,
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
            withInstaller: true,
          },
          meta: { pageTitle: "Explore" },
        },
        {
          path: ROUTES.StartupExploreDev,
          component: () => import("./components/startup/IntelliJExplore.vue"),
          props: {
            withInstaller: false,
          },
          meta: { pageTitle: "Explore" },
        },
        {
          path: ROUTES.StartupReport,
          component: () => import("./report-visualizer/Report.vue"),
          meta: { pageTitle: "Startup Report" },
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
          path: ROUTES.IntelliJDashboard,
          component: () => import("./components/intelliJ/PerformanceDashboard.vue"),
          meta: { pageTitle: "IntelliJ Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJIndexingDashboard,
          component: () => import("./components/intelliJ/IndexingDashboard.vue"),
          meta: { pageTitle: "IntelliJ Indexing Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJJBRDashboard,
          component: () => import("./components/intelliJ/Jbr21PerformanceDashboard.vue"),
          meta: { pageTitle: "JBR Performance dashboard" },
        },
        {
          path: ROUTES.IntelliJIncrementalCompilationDashboard,
          component: () => import("./components/intelliJ/experiments/IncrementalCompilationDashboard.vue"),
          meta: { pageTitle: "IntelliJ Incremental Compilation dashboard" },
        },
        {
          path: ROUTES.IntelliJScalabilityDashboard,
          component: () => import("./components/intelliJ/experiments/ScalabilityDashboard.vue"),
          meta: { pageTitle: "IntelliJ Scalability dashboard" },
        },
        {
          path: ROUTES.IntelliJGradleDashboard,
          component: () => import("./components/intelliJ/build-tools/gradle/GradleImportPerformanceDashboard.vue"),
          meta: { pageTitle: "Gradle Import dashboard" },
        },
        {
          path: ROUTES.IntelliJGradleDashboardFastInstallers,
          component: () => import("./components/intelliJ/build-tools/gradle/GradleImportPerformanceDashboardFastInstallers.vue"),
          meta: { pageTitle: "Gradle Import dashboard fast installers" },
        },
        {
          path: ROUTES.IntelliJNewGradleDashboard,
          component: () => import("./components/intelliJ/build-tools/gradle/NewGradleImportPerformanceDashboard.vue"),
          meta: { pageTitle: "New gradle Import dashboard" },
        },
        {
          path: ROUTES.IntelliJNewGradleDashboardFastInstallers,
          component: () => import("./components/intelliJ/build-tools/gradle/NewGradleImportPerformanceDashboardFastInstallers.vue"),
          meta: { pageTitle: "New gradle Import dashboard fast installers" },
        },
        {
          path: ROUTES.IntelliJMavenDashboard,
          component: () => import("./components/intelliJ/build-tools/maven/MavenImportPerformanceDashboard.vue"),
          meta: { pageTitle: "Maven Import dashboard" },
        },
        {
          path: ROUTES.IntelliJMavenDashboardFastInstallers,
          component: () => import("./components/intelliJ/build-tools/maven/MavenImportPerformanceDashboardFastInstallers.vue"),
          meta: { pageTitle: "Maven Import dashboard fast installers" },
        },
        {
          path: ROUTES.IntelliJMavenImportersConfiguratorsDashboard,
          component: () => import("./components/intelliJ/build-tools/maven/MavenImportersAndConfiguratorsPerformanceDashboard.vue"),
          meta: { pageTitle: "Maven Importers And Configurators dashboard" },
        },
        {
          path: ROUTES.IntelliJJpsDashboard,
          component: () => import("./components/intelliJ/build-tools/jps/JpsImportPerformanceDashboard.vue"),
          meta: { pageTitle: "JPS Import dashboard" },
        },
        {
          path: ROUTES.IntelliJUltimateDashboard,
          component: () => import("./components/intelliJ/UltimateProjectsDashboard.vue"),
          meta: { pageTitle: "Ultimate Projects" },
        },
        {
          path: ROUTES.IntelliJPackageCheckerDashboard,
          component: () => import("./components/intelliJ/PackageCheckerDashboard.vue"),
          meta: { pageTitle: "Package Checker" },
        },
        {
          path: ROUTES.IntelliJFusDashboard,
          component: () => import("./components/intelliJ/fus/FUSDashboard.vue"),
          meta: { pageTitle: "FUS 1 dashboard" },
        },
        {
          path: ROUTES.IntelliJFusDevDashboard,
          component: () => import("./components/intelliJ/fus/FUSDevDashboard.vue"),
          meta: { pageTitle: "FUS 2 dashboard" },
        },
        {
          path: ROUTES.IntelliJFusHetznerDashboard,
          component: () => import("./components/intelliJ/fus/FUSHetznerDashboard.vue"),
          meta: { pageTitle: "FUS Import dashboard" },
        },
        {
          path: ROUTES.IntelliJFusStartupDashboard,
          component: () => import("./components/intelliJ/fus/FUSStartupDashboard.vue"),
          meta: { pageTitle: "FUS Startup dashboard" },
        },
        {
          path: ROUTES.IntelliJDevDashboard,
          component: () => import("./components/intelliJ/PerformanceDevDashboard.vue"),
          meta: { pageTitle: "IntelliJ Performance dashboard Fast Installer" },
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
          path: ROUTES.IntelliJExperimentsMonorepoDashboard,
          component: () => import("./components/intelliJ/experiments/IntelliJDotnetDashboard.vue"),
          meta: { pageTitle: "IntelliJ + Dotnet performance dashboard" },
        },
        {
          path: ROUTES.IntelliJExperimentsGradleSyncDashboard,
          component: () => import("./components/intelliJ/experiments/GradleSyncDashboard.vue"),
          meta: { pageTitle: "Performance of Gradle Sync in Smart and Dumb modes" },
        },
        {
          path: ROUTES.IntelliJTinyDashboard,
          component: () => import("./components/intelliJ/experiments/PerformanceTinyDashboard.vue"),
          meta: { pageTitle: "IntelliJ Performance dashboard (Tiny)" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesIndexingDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/IndexingDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Indexing" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesScanningDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/ScanningDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Scanning" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesFindUsagesDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/FindUsagesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Finding Usages" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesCompletionDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/CompletionDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Completion" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesFirstCodeAnalysisDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/FirstCodeAnalysisDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Code Analysis" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesNumberOfIndexedFilesDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/NumberOfIndexedFilesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Number of indexed files" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesNumberOfExtensionsDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/NumberOfSharedIndexesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Number of indexed by shared indexes files" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesTypingDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/TypingDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Typing max awt delay" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesDumbModeDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/DumbModeDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Dumb Mode Time" },
        },
        {
          path: ROUTES.IntelliJGCDashboard,
          component: () => import("./components/intelliJ/experiments/GarbageCollectorDashboard.vue"),
          meta: { pageTitle: "IntelliJ performance tests for different Garbage Collectors" },
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
          meta: { pageTitle: "IntelliJ Integration Performance Tests On Fast Installer" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "ideaSharedIndices",
          },
          meta: { pageTitle: "IntelliJ Integration Performance Tests For Shared Indexes" },
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
            dbName: "perfint",
            table: "idea",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
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
          path: ROUTES.PhpStormDashboard,
          component: () => import("./components/phpstorm/PerformanceDashboard.vue"),
          meta: { pageTitle: "PhpStorm Performance dashboard" },
        },
        {
          path: ROUTES.PhpStormLLMDashboard,
          component: () => import("./components/phpstorm/MLDashboard.vue"),
          meta: { pageTitle: "PhpStorm LLM Performance dashboard" },
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
          path: ROUTES.PhpStormCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
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
          path: ROUTES.GoLandInspectionDashboard,
          component: () => import("./components/goland/InspectionCodeAnalyzesDashboard.vue"),
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
          path: ROUTES.GoLandIndexingDashboard,
          component: () => import("./components/goland/IndexingDashboard.vue"),
          meta: { pageTitle: "GoLand Indexing dashboard" },
        },
        {
          path: ROUTES.GoLandScanningDashboard,
          component: () => import("./components/goland/ScanningDashboard.vue"),
          meta: { pageTitle: "GoLand Scanning dashboard" },
        },
        {
          path: ROUTES.GoLandCompletionDashboard,
          component: () => import("./components/goland/CompletionDashboard.vue"),
          meta: { pageTitle: "GoLand Completion dashboard" },
        },
        {
          path: ROUTES.GoLandTests,
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
            dbName: "perfint",
            table: "goland",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.GoLandCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
            table: "goland",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.GoLandSharedIndicesIndexingDashboard,
          component: () => import("./components/goland/sharedIndexes/IndexingDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Indexing" },
        },
        {
          path: ROUTES.GoLandSharedIndicesScanningDashboard,
          component: () => import("./components/goland/sharedIndexes/ScanningDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Scanning" },
        },
        {
          path: ROUTES.GoLandSharedIndicesFindUsagesDashboard,
          component: () => import("./components/goland/sharedIndexes/FindUsagesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Finding Usages" },
        },
        {
          path: ROUTES.GoLandSharedIndicesCompletionDashboard,
          component: () => import("./components/goland/sharedIndexes/CompletionDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Completion" },
        },
        {
          path: ROUTES.GoLandSharedIndicesFirstCodeAnalysisDashboard,
          component: () => import("./components/goland/sharedIndexes/FirstCodeAnalysisDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Code Analysis" },
        },
        {
          path: ROUTES.GoLandSharedIndicesNumberOfIndexedFilesDashboard,
          component: () => import("./components/goland/sharedIndexes/NumberOfIndexedFilesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Number of indexed files" },
        },
        {
          path: ROUTES.GoLandSharedIndicesNumberOfExtensionsDashboard,
          component: () => import("./components/goland/sharedIndexes/NumberOfSharedIndexesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Number of indexed by shared indexes files" },
        },
        {
          path: ROUTES.GoLandSharedIndicesTypingDashboard,
          component: () => import("./components/goland/sharedIndexes/TypingDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Typing max awt delay" },
        },
        {
          path: ROUTES.GoLandSharedIndicesDumbModeDashboard,
          component: () => import("./components/goland/sharedIndexes/DumbModeDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Dumb Mode Time" },
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
          path: ROUTES.PyCharmDashboard,
          component: () => import("./components/pycharm/PerformanceDashboard.vue"),
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
            dbName: "perfint",
            table: "pycharm",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.PyCharmSharedIndicesIndexingDashboard,
          component: () => import("./components/pycharm/sharedIndexes/IndexingDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Indexing" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesScanningDashboard,
          component: () => import("./components/pycharm/sharedIndexes/ScanningDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Scanning" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesFindUsagesDashboard,
          component: () => import("./components/pycharm/sharedIndexes/FindUsagesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Finding Usages" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesCompletionDashboard,
          component: () => import("./components/pycharm/sharedIndexes/CompletionDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Completion" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesFirstCodeAnalysisDashboard,
          component: () => import("./components/pycharm/sharedIndexes/FirstCodeAnalysisDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Code Analysis" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesNumberOfIndexedFilesDashboard,
          component: () => import("./components/pycharm/sharedIndexes/NumberOfIndexedFilesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Number of indexed files" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesNumberOfExtensionsDashboard,
          component: () => import("./components/pycharm/sharedIndexes/NumberOfSharedIndexesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Number of indexed by shared indexes files" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesTypingDashboard,
          component: () => import("./components/pycharm/sharedIndexes/TypingDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Typing max awt delay" },
        },
        {
          path: ROUTES.PyCharmSharedIndicesDumbModeDashboard,
          component: () => import("./components/pycharm/sharedIndexes/DumbModeDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indexes Dashboard: Dumb Mode Time" },
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
          path: ROUTES.WebStormDashboard,
          component: () => import("./components/webstorm/PerformanceDashboard.vue"),
          meta: { pageTitle: "WebStorm Performance dashboard" },
        },
        {
          path: ROUTES.WebStormTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "webstorm",
            initialMachine: "linux-blade-hetzner",
          },
          meta: { pageTitle: "WebStorm Performance tests" },
        },
        {
          path: ROUTES.WebStormDashboardNEXT,
          component: () => import("./components/webstorm/PerformanceDashboardNEXT.vue"),
          meta: { pageTitle: "WebStorm NEXT" },
        },
        {
          path: ROUTES.WebStormCompare,
          component: () => import("./components/common/compare/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "webstorm",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.WebStormCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
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
          path: ROUTES.RubyMineDashboard,
          component: () => import("./components/rubymine/PerformanceDashboard.vue"),
          meta: { pageTitle: "RubyMine Performance Dashboard" },
        },
        {
          path: ROUTES.RubyMineIndexingDashBoard,
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
            dbName: "perfint",
            table: "ruby",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
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
          path: ROUTES.KotlinCodeAnalysis,
          component: () => import("./components/kotlin/KotlinCodeAnalysisDashboard.vue"),
          meta: { pageTitle: "Code analysis" },
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
          path: ROUTES.KotlinHighlightingDev,
          component: () => import("./components/kotlin/dev/HighlightingDashboard.vue"),
          meta: { pageTitle: "Kotlin highlighting (dev/fast)" },
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
          path: ROUTES.KotlinMPPDashboard,
          component: () => import("./components/kotlin/mpp/Dashboard.vue"),
          meta: { pageTitle: "Mpp projects" },
        },
        {
          path: ROUTES.RustPluginDashboard,
          component: () => import("./components/rust/PerformanceDashboardIdeWithRustPlugin.vue"),
          props: {
            releaseConfigurator: eap,
          },
          meta: { pageTitle: "Rust Plugin Performance dashboard" },
        },
        {
          path: ROUTES.RustRoverDashboard,
          component: () => import("./components/rust/PerformanceDashboardRustRover.vue"),
          meta: { pageTitle: "RustRover Performance dashboard" },
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
          path: ROUTES.FleetTest,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "fleet",
            table: "measure",
            initialMachine: "linux-blade-hetzner",
            withInstaller: false,
            unit: "ns",
          },
          meta: { pageTitle: "Fleet Performance tests" },
        },
        {
          path: ROUTES.FleetPerfDashboard,
          component: () => import("./components/fleet/PerformanceDashboard.vue"),
          meta: { pageTitle: "Fleet Performance dashboard" },
        },
        {
          path: ROUTES.FleetStartupDashboard,
          component: () => import("./components/fleet/FleetDashboard.vue"),
          meta: { pageTitle: "Fleet Startup dashboard" },
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
          path: ROUTES.BazelBspDashboard,
          component: () => import("./components/bazel/BazelBSPDashboard.vue"),
          meta: { pageTitle: "Bazel BSP dashboard" },
        },
        {
          path: ROUTES.IntelliJBspDashboard,
          component: () => import("./components/bazel/IntelliJBSPDashboard.vue"),
          meta: { pageTitle: "IntelliJ BSP dashboard" },
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
            dbName: "perfint",
            table: "clion",
            initialMachine: "Linux EC2 c5a(d).xlarge (4 vCPU, 8 GB)",
          },
          meta: { pageTitle: "CLion tests" },
        },
        {
          path: ROUTES.ClionStartupDashboard,
          component: () => import("./components/common/StartupProductDashboard.vue"),
          props: {
            product: "CL",
            defaultProject: "cmake",
          },
          meta: { pageTitle: "CLion Startup dashboard" },
        },
        {
          path: ROUTES.ClionPerfDashboard,
          component: () => import("./components/clion/PerformanceDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "clion",
            initialMachine: "Linux EC2 c5a(d).xlarge (4 vCPU, 8 GB)",
          },
          meta: { pageTitle: "CLion dashboard" },
        },
        {
          path: ROUTES.ClionDetailedPerfDashboard,
          component: () => import("./components/clion/DetailedPerformanceDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "clion",
            initialMachine: "Linux EC2 c5a(d).xlarge (4 vCPU, 8 GB)",
          },
          meta: { pageTitle: "CLion Detailed Performance dashboard" },
        },
        {
          path: ROUTES.ClionCompareBranches,
          component: () => import("./components/common/compare/CompareBranches.vue"),
          props: {
            dbName: "perfint",
            table: "clion",
          },
          meta: { pageTitle: COMPARE_BRANCHES_LABEL },
        },
        {
          path: ROUTES.VcsIdeaDashboard,
          component: () => import("./components/vcs/PerformanceDashboard.vue"),
          meta: { pageTitle: "Vcs Idea performance dashboard" },
        },
        {
          path: ROUTES.VcsSpaceDashboard,
          component: () => import("./components/vcs/PerformanceSpaceDashboard.vue"),
          meta: { pageTitle: "Vcs Space performance dashboard" },
        },
        {
          path: ROUTES.VcsStarterDashboard,
          component: () => import("./components/vcs/PerformanceStarterDashboard.vue"),
          meta: { pageTitle: "Vcs Starer performance dashboard" },
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
          path: ROUTES.IJentBenchmarks,
          component: () => import("./components/ijent/IJentBenchmarkTests.vue"),
          props: {
            dbName: "perfintDev_ijent",
            table: "report",
            initialMachine: "windows-blade-hetzner",
            withInstaller: false,
          },
          meta: { pageTitle: "IJent Benchmarks" },
        },
        {
          path: ROUTES.IJentPerfTests,
          component: () => import("./components/ijent/IJentPerformanceTests.vue"),
          meta: { pageTitle: "IJent Performance" },
        },
        {
          path: ROUTES.LLMDevTests,
          component: () => import("./components/ml/dev/LLMDashboard.vue"),
          meta: { pageTitle: "LLM dashboard" },
        },
        {
          path: ROUTES.FullLineDevTests,
          component: () => import("./components/ml/dev/FullLineDashboard.vue"),
          meta: { pageTitle: "FullLine dashboard" },
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
            product: "DG",
            defaultProject: "empty project",
          },
          meta: { pageTitle: "DataGrip Startup dashboard" },
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
      ],
    },
  ]
}
