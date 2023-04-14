import { ParentRouteRecord } from "shared/src/route"

const enum ROUTE_PREFIX {
  Startup = "/ij",
  IntelliJ = "/intellij",
  IntelliJBuildTools = "/intellij/buildTools",
  IntelliJSharedIndices = "/intellij/sharedIndices",
  IntelliJUltimate = "/intellij/ultimate",
  IntelliJPackageChecker = "/intellij/packageChecker",
  PhpStorm = "/phpstorm",
  GoLand = "/goland",
  RubyMine = "/rubymine",
  Kotlin = "/kotlin",
  Rust = "/rust",
  Scala = "/scala",
  JBR = "/jbr",
  Fleet = "/fleet",
}

enum ROUTES {
  StartupPulse = `${ROUTE_PREFIX.Startup}/pulse`,
  StartupProgress = `${ROUTE_PREFIX.Startup}/progressOverTime`,
  StartupModuleLoading = `${ROUTE_PREFIX.Startup}/moduleLoading`,
  StartupExplore = `${ROUTE_PREFIX.Startup}/explore`,
  StartupReport = `${ROUTE_PREFIX.Startup}/report`,
  IntelliJDashboard = `${ROUTE_PREFIX.IntelliJ}/dashboard`,
  IntelliJDevDashboard = `${ROUTE_PREFIX.IntelliJ}/devDashboard`,
  IntelliJTests = `${ROUTE_PREFIX.IntelliJ}/tests`,
  IntelliJDevTests = `${ROUTE_PREFIX.IntelliJ}/devTests`,
  IntelliJCompare = `${ROUTE_PREFIX.IntelliJ}/compare`,
  IntelliJGradleDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/gradleDashboard`,
  IntelliJMavenDashboard = `${ROUTE_PREFIX.IntelliJBuildTools}/mavenDashboard`,
  IntelliJBuildTests = `${ROUTE_PREFIX.IntelliJBuildTools}/tests`,
  IntelliJUltimateDashboard = `${ROUTE_PREFIX.IntelliJUltimate}/dashboard`,
  IntelliJUltimateTests = `${ROUTE_PREFIX.IntelliJUltimate}/tests`,
  IntelliJSharedIndicesDashboard = `${ROUTE_PREFIX.IntelliJSharedIndices}/dashboard`,
  IntelliJSharedIndicesTests = `${ROUTE_PREFIX.IntelliJSharedIndices}/tests`,
  IntelliJPackageCheckerDashboard = `${ROUTE_PREFIX.IntelliJPackageChecker}/dashboard`,
  IntelliJPackageCheckerTests = `${ROUTE_PREFIX.IntelliJPackageChecker}/tests`,
  PhpStormDashboard = `${ROUTE_PREFIX.PhpStorm}/dashboard`,
  PhpStormWithPluginsDashboard = `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests = `${ROUTE_PREFIX.PhpStorm}/tests`,
  PhpStormWithPluginsTests = `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,
  PhpStormCompare = `${ROUTE_PREFIX.PhpStorm}/compare`,
  KotlinDashboard = `${ROUTE_PREFIX.Kotlin}/dashboard`,
  KotlinDashboardDev = `${ROUTE_PREFIX.Kotlin}/dashboardDev`,
  KotlinExplore = `${ROUTE_PREFIX.Kotlin}/explore`,
  KotlinExploreDev = `${ROUTE_PREFIX.Kotlin}/exploreDev`,
  KotlinCompletionDev = `${ROUTE_PREFIX.Kotlin}/completionDev`,
  KotlinHighlightingDev = `${ROUTE_PREFIX.Kotlin}/highlightingDev`,
  KotlinFindUsagesDev = `${ROUTE_PREFIX.Kotlin}/findUsagesDev`,
  KotlinRefactoringDev = `${ROUTE_PREFIX.Kotlin}/refactoringDev`,
  KotlinCompare = `${ROUTE_PREFIX.Kotlin}/compare`,
  GoLandDashboard = `${ROUTE_PREFIX.GoLand}/dashboard`,
  GoLandTests = `${ROUTE_PREFIX.GoLand}/tests`,
  GoLandCompare = `${ROUTE_PREFIX.GoLand}/compare`,
  RubyMineDashboard = `${ROUTE_PREFIX.RubyMine}/dashboard`,
  RubyMineTests = `${ROUTE_PREFIX.RubyMine}/tests`,
  RubyMineCompare = `${ROUTE_PREFIX.RubyMine}/compare`,
  RustTests = `${ROUTE_PREFIX.Rust}/tests`,
  RustCompare = `${ROUTE_PREFIX.Rust}/compare`,
  ScalaTests = `${ROUTE_PREFIX.Scala}/tests`,
  ScalaCompare = `${ROUTE_PREFIX.Scala}/compare`,
  JBRTests = `${ROUTE_PREFIX.JBR}/tests`,
  MapBenchDashboard = `${ROUTE_PREFIX.JBR}/mapbenchDashboard`,
  DaCapoDashboard = `${ROUTE_PREFIX.JBR}/dacapoDashboard`,
  J2DBenchDashboard = `${ROUTE_PREFIX.JBR}/j2dDashboard`,
  JavaDrawDashboard = `${ROUTE_PREFIX.JBR}/javaDrawDashboard`,
  RenderDashboard = `${ROUTE_PREFIX.JBR}/renderDashboard`,
  FleetTest = `${ROUTE_PREFIX.Fleet}/tests`,
  FleetPerfDashboard = `${ROUTE_PREFIX.Fleet}/perfDashboard`,
  FleetStartupDashboard = `${ROUTE_PREFIX.Fleet}/startupDashboard`,
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
  url: ROUTE_PREFIX|ROUTES
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
          url: ROUTES.StartupProgress,
          label: "Progress Over Time",
        },
        {
          url: ROUTES.StartupModuleLoading,
          label: "Module Loading",
        },
        {
          url: ROUTES.StartupExplore,
          label: "Explore",
        },
        {
          url: ROUTES.StartupReport,
          label: "Report",
        }],
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
        }],
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
        }],
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
          url: ROUTES.KotlinExplore,
          label: TESTS_LABEL,
        },
        {
          url: ROUTES.KotlinDashboardDev,
          label: "Dashboard (dev/fast installer)",
        },
        {
          url: ROUTES.KotlinExploreDev,
          label: "Explore (dev/fast installer)",
        },
        {
          url: ROUTES.KotlinCompletionDev,
          label: "Completion (dev)",
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
          url: ROUTES.KotlinCompare,
          label: COMPARE_BUILDS_LABEL,
        }],
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
const RUST: Product = {
  url: ROUTE_PREFIX.Rust,
  label: "Rust",
  children: [
    {
      url: ROUTE_PREFIX.Rust,
      label: "",
      tabs: [
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

export const PRODUCTS = [IJ_STARTUP, IDEA, PHPSTORM, KOTLIN, GOLAND, RUBYMINE, RUST, SCALA, JBR, FLEET]
export function getNavigationElement(path: string): Product {
  return PRODUCTS.find(PRODUCTS => path.startsWith(PRODUCTS.url)) ?? PRODUCTS[0]
}


export function getNewDashboardRoutes(): Array<ParentRouteRecord> {
  return [
    {
      children: [
        {
          path: ROUTES.StartupPulse,
          component: () => import("./components/startup/IntelliJPulse.vue"),
          meta: {pageTitle: "Pulse"},
        },
        {
          path: ROUTES.StartupProgress,
          component: () => import("./components/startup/IntelliJProgressOverTime.vue"),
          meta: {pageTitle: "Progress Over Time"},
        },
        {
          path: ROUTES.StartupModuleLoading,
          component: () => import("./components/startup/IntelliJModuleLoading.vue"),
          meta: {pageTitle: "Module Loading"},
        },
        {
          path: ROUTES.StartupExplore,
          component: () => import("./components/startup/IntelliJExplore.vue"),
          meta: {pageTitle: "Explore"},
        },
        {
          path: ROUTES.StartupReport,
          component: () => import("../../report-visualizer/src/Report.vue"),
          meta: {pageTitle: "Startup Report"},
        },
        {
          path: ROUTES.IntelliJDashboard,
          component: () => import("./components/intelliJ/PerformanceDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance dashboard"},
        },
        {
          path: ROUTES.IntelliJGradleDashboard,
          component: () => import("./components/intelliJ/GradleImportPerformanceDashboard.vue"),
          meta: {pageTitle: "Gradle Import dashboard"},
        },
        {
          path: ROUTES.IntelliJMavenDashboard,
          component: () => import("./components/intelliJ/MavenImportPerformanceDashboard.vue"),
          meta: {pageTitle: "Maven Import dashboard"},
        },
        {
          path: ROUTES.IntelliJUltimateDashboard,
          component: () => import("./components/intelliJ/UltimateProjectsDashboard.vue"),
          meta: {pageTitle: "Ultimate Projects"},
        },
        {
          path: ROUTES.IntelliJPackageCheckerDashboard,
          component: () => import("./components/intelliJ/PackageCheckerDashboard.vue"),
          meta: {pageTitle: "Package Checker"},
        },
        {
          path: ROUTES.IntelliJDevDashboard,
          component: () => import("./components/intelliJ/PerformanceDevDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance dashboard Fast Installer"},
        },
        {
          path: ROUTES.IntelliJSharedIndicesDashboard,
          component: () => import("./components/intelliJ/SharedIndicesPerformanceDashboard.vue"),
          meta: {pageTitle: "IntelliJ Integration Performance Tests For Shared Indices Dashboard"},
        },
        {
          path: `${ROUTE_PREFIX.IntelliJ}/:subproject?/tests`,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "idea",
            initialMachine: "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
          },
          meta: {pageTitle: "IntelliJ Performance tests"},
        },
        {
          path: ROUTES.IntelliJDevTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "idea",
            initialMachine: "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)",
            withInstaller: false
          },
          meta: {pageTitle: "IntelliJ Integration Performance Tests On Fast Installer"},
        },
        {
          path: ROUTES.IntelliJSharedIndicesTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "ideaSharedIndices"
          },
          meta: {pageTitle: "IntelliJ Integration Performance Tests For Shared Indices"},
        },
        {
          path: ROUTES.IntelliJCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "idea"
          },
          meta: {pageTitle: COMPARE_BUILDS_LABEL},
        },
        {
          path: ROUTES.PhpStormDashboard,
          component: () => import("./components/phpstorm/PerformanceDashboard.vue"),
          meta: {pageTitle: "PhpStorm Performance dashboard"},
        },
        {
          path: ROUTES.PhpStormWithPluginsDashboard,
          component: () => import("./components/phpstorm/PerformanceDashboardWithPlugins.vue"),
          meta: {pageTitle: "PhpStorm With Plugins Performance dashboard"},
        },
        {
          path: ROUTES.PhpStormWithPluginsTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "phpstormWithPlugins",
            initialMachine: "linux-blade-hetzner",
          },
          meta: {pageTitle: "PhpStorm Performance tests with plugins"},
        },
        {
          path: ROUTES.PhpStormTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm",
            initialMachine: "linux-blade-hetzner"
          },
          meta: {pageTitle: "PhpStorm Performance tests"},
        },
        {
          path: ROUTES.PhpStormCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm"
          },
          meta: {pageTitle: COMPARE_BUILDS_LABEL},
        },
        {
          path: ROUTES.GoLandDashboard,
          component: () => import("./components/goland/PerformanceDashboard.vue"),
          meta: {pageTitle: "GoLand Performance dashboard"},
        },
        {
          path: ROUTES.GoLandTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "goland",
            initialMachine: "linux-blade-hetzner"
          },
          meta: {pageTitle: "GoLand Performance tests"},
        },
        {
          path: ROUTES.GoLandCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "goland"
          },
          meta: {pageTitle: COMPARE_BUILDS_LABEL},
        },
        {
          path: ROUTES.RubyMineDashboard,
          component: () => import("./components/rubymine/PerformanceDashboard.vue"),
          meta: {pageTitle: "RubyMine Performance dashboard"},
        },
        {
          path: ROUTES.RubyMineTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "ruby",
            initialMachine: "Linux Munich i7-3770, 32 Gb"
          },
          meta: {pageTitle: "RubyMine Performance tests"},
        },
        {
          path: ROUTES.RubyMineCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "ruby"
          },
          meta: {pageTitle: COMPARE_BUILDS_LABEL},
        },

        {
          path: ROUTES.KotlinExplore,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin",
            initialMachine: "linux-blade-hetzner",
          },
          meta: {pageTitle: "Kotlin Performance tests explore"},
        },
        {
          path: ROUTES.KotlinExploreDev,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfintDev",
            table: "kotlin",
            initialMachine: "linux-blade-hetzner",
            withInstaller: false
          },
          meta: {pageTitle: "Kotlin Performance tests explore (dev/fast installer)"},
        },
        {
          path: ROUTES.KotlinDashboard,
          component: () => import("./components/kotlin/PerformanceDashboard.vue"),
          meta: {pageTitle: "Kotlin Performance dashboard"},
        },
        {
          path: ROUTES.KotlinDashboardDev,
          component: () => import("./components/kotlin/dev/PerformanceDevDashboard.vue"),
          meta: {pageTitle: "Kotlin Performance dashboard (dev/fast installer)"},
        },
        {
          path: ROUTES.KotlinCompletionDev,
          component: () => import("./components/kotlin/dev/CompletionDashboard.vue"),
          meta: {pageTitle: "Kotlin completion (dev/fast)"},
        },
        {
          path: ROUTES.KotlinHighlightingDev,
          component: () => import("./components/kotlin/dev/HighlightingDashboard.vue"),
          meta: {pageTitle: "Kotlin highlighting (dev/fast)"},
        },
        {
          path: ROUTES.KotlinFindUsagesDev,
          component: () => import("./components/kotlin/dev/FindUsagesDashboard.vue"),
          meta: {pageTitle: "Kotlin findUsages (dev/fast)"},
        },
        {
          path: ROUTES.KotlinRefactoringDev,
          component: () => import("./components/kotlin/dev/RefactoringDashboard.vue"),
          meta: {pageTitle: "Kotlin refactoring (dev/fast)"},
        },
        {
          path: ROUTES.KotlinCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin"
          },
          meta: {pageTitle: COMPARE_BUILDS_LABEL},
        },
        {
          path: ROUTES.RustTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
            initialMachine: "Linux EC2 m5d.xlarge or 5d.xlarge or m5ad.xlarge"
          },
          meta: {pageTitle: "Rust Performance tests"},
        },
        {
          path: ROUTES.RustCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "rust"
          },
          meta: {pageTitle: COMPARE_BUILDS_LABEL},
        },
        {
          path: ROUTES.ScalaTests,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "perfint",
            table: "scala",
            initialMachine: "linux-blade-hetzner"
          },
          meta: {pageTitle: "Scala Performance tests"},
        },
        {
          path: ROUTES.ScalaCompare,
          component: () => import("./components/common/CompareBuilds.vue"),
          props: {
            dbName: "perfint",
            table: "scala"
          },
          meta: {pageTitle: COMPARE_BUILDS_LABEL},
        },
        {
          path: ROUTES.JBRTests,
          component: () => import("./components/jbr/PerformanceTests.vue"),
          meta: {pageTitle: "JBR Performance tests"},
        },
        {
          path: ROUTES.MapBenchDashboard,
          component: () => import("./components/jbr/MapBenchDashboard.vue"),
          meta: {pageTitle: "MapBench Dashboard"},
        },
        {
          path: ROUTES.DaCapoDashboard,
          component: () => import("./components/jbr/DaCapoDashboard.vue"),
          meta: {pageTitle: "DaCapo Dashboard"},
        },
        {
          path: ROUTES.J2DBenchDashboard,
          component: () => import("./components/jbr/J2DBenchDashboard.vue"),
          meta: {pageTitle: "J2DBench Dashboard"},
        },
        {
          path: ROUTES.JavaDrawDashboard,
          component: () => import("./components/jbr/JavaDrawDashboard.vue"),
          meta: {pageTitle: "JavaDraw Dashboard"},
        },
        {
          path: ROUTES.RenderDashboard,
          component: () => import("./components/jbr/RenderDashboard.vue"),
          meta: {pageTitle: "Render Dashboard"},
        },
        {
          path: ROUTES.FleetTest,
          component: () => import("./components/common/PerformanceTests.vue"),
          props: {
            dbName: "fleet",
            table: "measure",
            initialMachine: "linux-blade-hetzner",
            withInstaller: false,
            unit: "ns"
          },
          meta: {pageTitle: "Fleet Performance tests"},
        },
        {
          path: ROUTES.FleetPerfDashboard,
          component: () => import("./components/fleet/PerformanceDashboard.vue"),
          meta: {pageTitle: "Fleet Performance dashboard"},
        },
        {
          path: ROUTES.FleetStartupDashboard,
          component: () => import("./components/fleet/FleetDashboard.vue"),
          meta: {pageTitle: "Fleet Startup dashboard"},
        },
      ],
    },
  ]
}
