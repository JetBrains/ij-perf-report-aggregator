/* eslint-disable @typescript-eslint/prefer-literal-enum-member */

import { ParentRouteRecord } from "./components/common/route"

const enum ROUTE_PREFIX {
  Startup = "/ij",
  IntelliJ = "/intellij",
  IntelliJBuildTools = "/intellij/buildTools",
  IntelliJSharedIndices = "/intellij/sharedIndices",
  IntelliJUltimate = "/intellij/ultimate",
  IntelliJPackageChecker = "/intellij/packageChecker",
  IntelliJFus = "/intellij/fus",
  IntelliJExperiments = "/intellij/experiments",
  PhpStorm = "/phpstorm",
  GoLand = "/goland",
  RubyMine = "/rubymine",
  Kotlin = "/kotlin",
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
}
const TEST_ROUTE = "tests"
const DEV_TEST_ROUTE = "testsDev"
const DASHBOARD_ROUTE = "dashboard"
const COMPARE_ROUTE = "compare"

enum ROUTES {
  StartupPulse = `${ROUTE_PREFIX.Startup}/pulse`,
  StartupProgress = `${ROUTE_PREFIX.Startup}/progressOverTime`,
  StartupModuleLoading = `${ROUTE_PREFIX.Startup}/moduleLoading`,
  StartupGcAndMemory = `${ROUTE_PREFIX.Startup}/gcAndMemory`,
  StartupExplore = `${ROUTE_PREFIX.Startup}/explore`,
  StartupReport = `${ROUTE_PREFIX.Startup}/report`,
  IntelliJDashboard = `${ROUTE_PREFIX.IntelliJ}/${DASHBOARD_ROUTE}`,
  IntelliJIndexingDashboard = `${ROUTE_PREFIX.IntelliJ}/indexingDashboard`,
  IntelliJTinyDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/dashboardTiny`,
  IntelliJIncrementalCompilationDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/incrementalCompilationDashboard`,
  IntelliJScalabilityDashboard = `${ROUTE_PREFIX.IntelliJExperiments}/scalabilityDashboard`,
  IntelliJDevDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboardDev`,
  IntelliJTests = `${ROUTE_PREFIX.IntelliJ}/${TEST_ROUTE}`,
  IntelliJDevTests = `${ROUTE_PREFIX.IntelliJ}/${DEV_TEST_ROUTE}`,
  IntelliJCompare = `${ROUTE_PREFIX.IntelliJ}/${COMPARE_ROUTE}`,
  IntelliJGradleDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/gradleDashboard`,
  IntelliJMavenDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/mavenDashboard`,
  IntelliJJpsDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/jpsDashboard`,
  IntelliJBuildTests = `${ROUTE_PREFIX.IntelliJBuildTools}/${TEST_ROUTE}`,
  IntelliJUltimateDashboard = `${ROUTE_PREFIX.IntelliJUltimate}/${DASHBOARD_ROUTE}`,
  IntelliJUltimateTests = `${ROUTE_PREFIX.IntelliJUltimate}/${TEST_ROUTE}`,
  IntelliJSharedIndicesIndexingDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesIndexingDashboard`,
  IntelliJSharedIndicesScanningDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesScanningDashboard`,
  IntelliJSharedIndicesFindUsagesDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesFindUsagesDashboard`,
  IntelliJSharedIndicesCompletionDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesCompletionDashboard`,
  IntelliJSharedIndicesFirstCodeAnalysisDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesFirstCodeAnalysisDashboard`,
  IntelliJSharedIndicesNumberOfIndexedFilesDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesIndexedFilesDashboard`,
  IntelliJSharedIndicesNumberOfExtensionsDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/sharedIndexesNumberOfExtensionsDashboard`,
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
  PhpStormWithPluginsDashboard = `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests = `${ROUTE_PREFIX.PhpStorm}/${TEST_ROUTE}`,
  PhpStormWithPluginsTests = `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,
  PhpStormCompare = `${ROUTE_PREFIX.PhpStorm}/${COMPARE_ROUTE}`,
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
  KotlinK1VsK2ComparisonDev = `${ROUTE_PREFIX.Kotlin}/k1VsK2ComparisonDev`,
  KotlinCompare = `${ROUTE_PREFIX.Kotlin}/${COMPARE_ROUTE}`,
  GoLandDashboard = `${ROUTE_PREFIX.GoLand}/${DASHBOARD_ROUTE}`,
  GoLandTests = `${ROUTE_PREFIX.GoLand}/${TEST_ROUTE}`,
  GoLandCompare = `${ROUTE_PREFIX.GoLand}/${COMPARE_ROUTE}`,
  PyCharmDashboard = `${ROUTE_PREFIX.PyCharm}/${DASHBOARD_ROUTE}`,
  PyCharmTests = `${ROUTE_PREFIX.PyCharm}/${TEST_ROUTE}`,
  PyCharmCompare = `${ROUTE_PREFIX.PyCharm}/${COMPARE_ROUTE}`,
  WebStormDashboard = `${ROUTE_PREFIX.WebStorm}/${DASHBOARD_ROUTE}`,
  WebStormTests = `${ROUTE_PREFIX.WebStorm}/${TEST_ROUTE}`,
  WebStormCompare = `${ROUTE_PREFIX.WebStorm}/${COMPARE_ROUTE}`,
  RubyMineDashboard = `${ROUTE_PREFIX.RubyMine}/${DASHBOARD_ROUTE}`,
  RubyMineTests = `${ROUTE_PREFIX.RubyMine}/${TEST_ROUTE}`,
  RubyMineCompare = `${ROUTE_PREFIX.RubyMine}/${COMPARE_ROUTE}`,
  RustDashboard = `${ROUTE_PREFIX.Rust}/${DASHBOARD_ROUTE}`,
  RustTests = `${ROUTE_PREFIX.Rust}/${TEST_ROUTE}`,
  RustCompare = `${ROUTE_PREFIX.Rust}/${COMPARE_ROUTE}`,
  ScalaTests = `${ROUTE_PREFIX.Scala}/${TEST_ROUTE}`,
  ScalaCompare = `${ROUTE_PREFIX.Scala}/${COMPARE_ROUTE}`,
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
  ClionTest = `${ROUTE_PREFIX.Clion}/${TEST_ROUTE}`,
  ClionPerfDashboard = `${ROUTE_PREFIX.Clion}/perfDashboard`,
  ClionDetailedPerfDashboard = `${ROUTE_PREFIX.Clion}/detailedPerfDashboard`,
  VcsIdeaDashboard = `${ROUTE_PREFIX.Vcs}/idea`,
  VcsSpaceDashboard = `${ROUTE_PREFIX.Vcs}/space`,
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
const DASHBOARD_LABEL = "Dashboard"

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
          url: ROUTES.IntelliJDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.IntelliJDevDashboard,
          label: "Dashboard (Fast Installer)",
        },
        {
          url: ROUTES.IntelliJIndexingDashboard,
          label: "Indexing Dashboard",
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
          url: ROUTES.IntelliJCompare,
          label: COMPARE_BUILDS_LABEL,
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
          url: ROUTES.IntelliJMavenDashboard,
          label: "Maven Import",
        },
        {
          url: ROUTES.IntelliJJpsDashboard,
          label: "JPS Import",
        },
        {
          url: ROUTES.IntelliJBuildTests,
          label: TESTS_LABEL,
        },
      ],
    },
    {
      url: ROUTE_PREFIX.IntelliJSharedIndices,
      label: "Shared Indices",
      tabs: [
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
      label: "Vcs",
      tabs: [
        {
          url: ROUTES.VcsIdeaDashboard,
          label: "Performance dashboard idea project",
        },
        {
          url: ROUTES.VcsSpaceDashboard,
          label: "Performance dashboard space project",
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
          url: ROUTES.PhpStormDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.PhpStormTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.PhpStormCompare,
          label: COMPARE_BUILDS_LABEL,
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
      label: "",
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
          url: ROUTES.KotlinK1VsK2ComparisonDev,
          label: "K1 vs. K2 (dev)",
        },
        {
          url: ROUTES.KotlinCompare,
          label: COMPARE_BUILDS_LABEL,
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
      label: "",
      tabs: [
        {
          url: ROUTES.GoLandDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.GoLandTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.GoLandCompare,
          label: COMPARE_BUILDS_LABEL,
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
          url: ROUTES.RubyMineDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.RubyMineTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.RubyMineCompare,
          label: COMPARE_BUILDS_LABEL,
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
      label: "",
      tabs: [
        {
          url: ROUTES.PyCharmDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.PyCharmTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.PyCharmCompare,
          label: COMPARE_BUILDS_LABEL,
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
          url: ROUTES.WebStormDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.WebStormTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.WebStormCompare,
          label: COMPARE_BUILDS_LABEL,
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
          url: ROUTES.RustDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.RustTests,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.RustCompare,
          label: COMPARE_BUILDS_LABEL,
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
          url: ROUTES.ScalaCompare,
          label: COMPARE_BUILDS_LABEL,
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
      ],
    },
  ],
}

export const PRODUCTS = [IJ_STARTUP, IDEA, PHPSTORM, KOTLIN, GOLAND, RUBYMINE, PYCHARM, WEBSTORM, RUST, SCALA, JBR, FLEET, BAZEL, QODANA, CLION]
export function getNavigationElement(path: string): Product {
  return PRODUCTS.find((PRODUCTS) => path.startsWith(PRODUCTS.url)) ?? PRODUCTS[0]
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
          meta: { pageTitle: "Explore" },
        },
        {
          path: ROUTES.StartupReport,
          component: () => import("./report-visualizer/Report.vue"),
          meta: { pageTitle: "Startup Report" },
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
          component: () => import("./components/intelliJ/GradleImportPerformanceDashboard.vue"),
          meta: { pageTitle: "Gradle Import dashboard" },
        },
        {
          path: ROUTES.IntelliJMavenDashboard,
          component: () => import("./components/intelliJ/MavenImportPerformanceDashboard.vue"),
          meta: { pageTitle: "Maven Import dashboard" },
        },
        {
          path: ROUTES.IntelliJJpsDashboard,
          component: () => import("./components/intelliJ/JpsImportPerformanceDashboard.vue"),
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
          meta: { pageTitle: "Performance Tests For Shared Indices Dashboard: Indexing" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesScanningDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/ScanningDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indices Dashboard: Scanning" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesFindUsagesDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/FindUsagesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indices Dashboard: Finding Usages" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesCompletionDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/CompletionDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indices Dashboard: Completion" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesFirstCodeAnalysisDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/FirstCodeAnalysisDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indices Dashboard: Code Analysis" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesNumberOfIndexedFilesDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/NumberOfIndexedFilesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indices Dashboard: Number of indexed files" },
        },
        {
          path: ROUTES.IntelliJSharedIndicesNumberOfExtensionsDashboard,
          component: () => import("./components/intelliJ/sharedIndexes/NumberOfSharedIndexesDashboard.vue"),
          meta: { pageTitle: "Performance Tests For Shared Indices Dashboard: Number of indexed by shared indexes files" },
        },
        {
          path: ROUTES.IntelliJGCDashboard,
          component: () => import("./components/intelliJ/experiments/GarbageCollectorDashboard.vue"),
          meta: { pageTitle: "IntelliJ performance tests for different Garbage Collectors" },
        },
        {
          path: `${ROUTE_PREFIX.IntelliJ}/:subproject?/tests`,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "idea",
            initialMachine: "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
          },
          meta: { pageTitle: "IntelliJ Performance tests" },
        },
        {
          path: ROUTES.IntelliJDevTests,
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
          meta: { pageTitle: "IntelliJ Integration Performance Tests For Shared Indices" },
        },
        {
          path: ROUTES.IntelliJCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "idea",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.PhpStormDashboard,
          component: () => import("./components/phpstorm/PerformanceDashboard.vue"),
          meta: { pageTitle: "PhpStorm Performance dashboard" },
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
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.GoLandDashboard,
          component: () => import("./components/goland/PerformanceDashboard.vue"),
          meta: { pageTitle: "GoLand Performance dashboard" },
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
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "goland",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
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
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "pycharm",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
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
          path: ROUTES.WebStormCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "webstorm",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.RubyMineDashboard,
          component: () => import("./components/rubymine/PerformanceDashboard.vue"),
          meta: { pageTitle: "RubyMine Performance dashboard" },
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
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "ruby",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
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
          path: ROUTES.KotlinK1VsK2ComparisonDev,
          component: () => import("./components/kotlin/dev/K1VsK2ComparisonDashboard.vue"),
          meta: { pageTitle: "Kotlin K1 vs. K2 (dev/fast)" },
        },
        {
          path: ROUTES.KotlinCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
        },
        {
          path: ROUTES.RustDashboard,
          component: () => import("./components/rust/PerformanceDashboard.vue"),
          meta: { pageTitle: "Rust Performance dashboard" },
        },
        {
          path: ROUTES.RustTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
            initialMachine: "Linux EC2 m5d.xlarge or 5d.xlarge or m5ad.xlarge",
          },
          meta: { pageTitle: "Rust Performance tests" },
        },
        {
          path: ROUTES.RustCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
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
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "scala",
          },
          meta: { pageTitle: COMPARE_BUILDS_LABEL },
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
          path: ROUTES.VcsIdeaDashboard,
          component: () => import("./components/vcs/PerformanceDashboard.vue"),
          meta: { pageTitle: "Vcs Idea performance dashboard" },
        },
        {
          path: ROUTES.VcsSpaceDashboard,
          component: () => import("./components/vcs/PerformanceSpaceDashboard.vue"),
          meta: { pageTitle: "Vcs Space performance dashboard" },
        },
      ],
    },
  ]
}
