import { ParentRouteRecord } from "shared/src/route"

export interface NavigationItem {
  path: string
  name: string
  key?: string
}

const enum ROUTE_PREFIX {
  Startup = "/startup",
  IntelliJ = "/intellij",
  PhpStorm = "/phpstorm",
  GoLand = "/goland",
  RubyMine = "/rubymine",
  Kotlin = "/kotlin",
  Rust = "/rust",
  Scala = "/scala",
  JBR = "/jbr",
  Fleet = "/fleet",
}

const ROUTES = {
  StartupPulse: `${ROUTE_PREFIX.Startup}/pulse`,
  StartupProgress: `${ROUTE_PREFIX.Startup}/progressOverTime`,
  StartupModuleLoading: `${ROUTE_PREFIX.Startup}/moduleLoading`,
  StartupExplore: `${ROUTE_PREFIX.Startup}/explore`,
  StartupReport: `${ROUTE_PREFIX.Startup}/report`,

  IntelliJDashboard: `${ROUTE_PREFIX.IntelliJ}/dashboard`,
  IntelliJGradleDashboard: `${ROUTE_PREFIX.IntelliJ}/gradleDashboard`,
  IntelliJMavenDashboard: `${ROUTE_PREFIX.IntelliJ}/mavenDashboard`,
  IntelliJDevDashboard: `${ROUTE_PREFIX.IntelliJ}/devDashboard`,
  IntelliJTests: `${ROUTE_PREFIX.IntelliJ}/tests`,
  IntelliJDevTests: `${ROUTE_PREFIX.IntelliJ}/devTests`,
  IntelliJSharedIndicesDashboard: `${ROUTE_PREFIX.IntelliJ}/sharedIndicesDashboard`,
  IntelliJSharedIndicesTests: `${ROUTE_PREFIX.IntelliJ}/sharedIndicesTests`,
  IntelliJCompare: `${ROUTE_PREFIX.IntelliJ}/compare`,

  PhpStormDashboard: `${ROUTE_PREFIX.PhpStorm}/dashboard`,
  PhpStormWithPluginsDashboard: `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests: `${ROUTE_PREFIX.PhpStorm}/tests`,
  PhpStormWithPluginsTests: `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,

  KotlinDashboard: `${ROUTE_PREFIX.Kotlin}/dashboard`,
  KotlinDashboardDev: `${ROUTE_PREFIX.Kotlin}/dashboardDev`,
  KotlinExplore: `${ROUTE_PREFIX.Kotlin}/explore`,
  KotlinExploreDev: `${ROUTE_PREFIX.Kotlin}/exploreDev`,
  KotlinCompletion: `${ROUTE_PREFIX.Kotlin}/completion`,
  KotlinCompletionDev: `${ROUTE_PREFIX.Kotlin}/completionDev`,

  GoLandDashboard: `${ROUTE_PREFIX.GoLand}/dashboard`,
  GoLandTests: `${ROUTE_PREFIX.GoLand}/tests`,

  RubyMineDashboard: `${ROUTE_PREFIX.RubyMine}/dashboard`,
  RubyMineTests: `${ROUTE_PREFIX.RubyMine}/tests`,

  RustTests: `${ROUTE_PREFIX.Rust}/tests`,
  ScalaTests: `${ROUTE_PREFIX.Scala}/tests`,

  JBRTests: `${ROUTE_PREFIX.JBR}/tests`,
  MapBenchDashboard: `${ROUTE_PREFIX.JBR}/mapbenchDashboard`,
  DaCapoDashboard: `${ROUTE_PREFIX.JBR}/dacapoDashboard`,
  J2DBenchDashboard: `${ROUTE_PREFIX.JBR}/j2dDashboard`,
  JavaDrawDashboard: `${ROUTE_PREFIX.JBR}/javaDrawDashboard`,
  RenderDashboard: `${ROUTE_PREFIX.JBR}/renderDashboard`,

  FleetTest: `${ROUTE_PREFIX.Fleet}/tests`,
  FleetPerfDashboard: `${ROUTE_PREFIX.Fleet}/perfDashboard`,
  FleetStartupDashboard: `${ROUTE_PREFIX.Fleet}/startupDashboard`,
}

export const topNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.StartupPulse,
    name: "IntelliJ Startup",
    key: ROUTE_PREFIX.Startup,
  },
  {
    path: ROUTES.IntelliJDashboard,
    name: "IntelliJ",
    key: ROUTE_PREFIX.IntelliJ,
  },
  {
    path: ROUTES.PhpStormDashboard,
    name: "PhpStorm",
    key: ROUTE_PREFIX.PhpStorm,
  },
  {
    path: ROUTES.KotlinDashboard,
    name: "Kotlin",
    key: ROUTE_PREFIX.Kotlin,
  },
  {
    path: ROUTES.GoLandDashboard,
    name: "GoLand",
    key: ROUTE_PREFIX.GoLand,
  },
  {
    path: ROUTES.RubyMineDashboard,
    name: "RubyMine",
    key: ROUTE_PREFIX.RubyMine,
  },
  {
    path: ROUTES.JBRTests,
    name: "JBR",
    key: ROUTE_PREFIX.JBR,
  },
  {
    path: ROUTES.FleetStartupDashboard,
    name: "Fleet",
    key: ROUTE_PREFIX.Fleet,
  },
  {
    path: ROUTES.RustTests,
    name: "Rust",
    key: ROUTE_PREFIX.Rust,
  },
  {
    path: ROUTES.ScalaTests,
    name: "Scala",
    key: ROUTE_PREFIX.Scala,
  }
]
export const startupTabNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.StartupPulse,
    name: "Pulse",
  },
  {
    path: ROUTES.StartupProgress,
    name: "Progress Over Time",
  },
  {
    path: ROUTES.StartupModuleLoading,
    name: "Module Loading",
  },
  {
    path: ROUTES.StartupExplore,
    name: "Explore",
  },
  {
    path: ROUTES.StartupReport,
    name: "Report",
  },
]

export const intelliJTabNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.IntelliJDashboard,
    name: "Performance dashboard",
  },
  {
    path: ROUTES.IntelliJGradleDashboard,
    name: "Gradle Import dashboard",
  },
  {
    path: ROUTES.IntelliJMavenDashboard,
    name: "Maven Import dashboard",
  },
  {
    path: ROUTES.IntelliJDevDashboard,
    name: "Performance dashboard (Fast Installer)",
  },
  {
    path: ROUTES.IntelliJSharedIndicesDashboard,
    name: "Shared Indices dashboard",
  },
  {
    path: ROUTES.IntelliJTests,
    name: "Performance tests",
  },
  {
    path: ROUTES.IntelliJDevTests,
    name: "Performance Tests (Fast Installer)",
  },
  {
    path: ROUTES.IntelliJSharedIndicesTests,
    name: "Shared Indices",
  },
  // {
  //   path: ROUTES.Compare,
  //   name: "Compare branches",
  // },
]

export const phpStormNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.PhpStormDashboard,
    name: "Dashboard",
  },
  {
    path: ROUTES.PhpStormWithPluginsDashboard,
    name: "Dashboard with Plugins",
  },
  {
    path: ROUTES.PhpStormTests,
    name: "Tests",
  },
  {
    path: ROUTES.PhpStormWithPluginsTests,
    name: "Tests with Plugins",
  },
]

export const GoLandNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.GoLandDashboard,
    name: "Dashboard",
  },
  {
    path: ROUTES.GoLandTests,
    name: "Tests",
  },
]

export const RubyMineNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.RubyMineDashboard,
    name: "Dashboard",
  },
  {
    path: ROUTES.RubyMineTests,
    name: "Tests",
  },
]

export const kotlinNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.KotlinDashboard,
    name: "Dashboard",
  },
  {
    path: ROUTES.KotlinExplore,
    name: "Explore",
  },
  {
    path: ROUTES.KotlinDashboardDev,
    name: "Dashboard (dev/fast installer)",
  },
  {
    path: ROUTES.KotlinExploreDev,
    name: "Explore (dev/fast installer)",
  },
  {
    path: ROUTES.KotlinCompletion,
    name: "Completion dashboard",
  },
  {
    path: ROUTES.KotlinCompletionDev,
    name: "Completion dashboard (dev/fast installer)",
  },
]

export const RustNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.RustTests,
    name: "Tests",
  },
]

export const ScalaNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.ScalaTests,
    name: "Tests",
  },
]

export const JBRNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.DaCapoDashboard,
    name: "DaCapo",
  },
  {
    path: ROUTES.J2DBenchDashboard,
    name: "J2DBench",
  },
  {
    path: ROUTES.JavaDrawDashboard,
    name: "JavaDraw",
  },
  {
    path: ROUTES.RenderDashboard,
    name: "Render",
  },
  {
    path: ROUTES.MapBenchDashboard,
    name: "MapBench",
  },
  {
    path: ROUTES.JBRTests,
    name: "Tests",
  }
]

export const FleetNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.FleetStartupDashboard,
    name: "Startup Dashboard",
  },
  {
    path: ROUTES.FleetPerfDashboard,
    name: "Performance Dashboard",
  },
  {
    path: ROUTES.FleetTest,
    name: "Tests",
  },
]

export function getNavigationTabs(path: string): NavigationItem[] {
  if (path.startsWith(ROUTE_PREFIX.Startup)) {
    return startupTabNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.IntelliJ)) {
    return intelliJTabNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.PhpStorm)) {
    return phpStormNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.Kotlin)) {
    return kotlinNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.GoLand)) {
    return GoLandNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.RubyMine)) {
    return RubyMineNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.Rust)) {
    return RustNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.Scala)) {
    return ScalaNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.JBR)) {
    return JBRNavigationItems
  }
  if (path.startsWith(ROUTE_PREFIX.Fleet)) {
    return FleetNavigationItems
  }

  return []
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
          path: ROUTES.IntelliJDevDashboard,
          component: () => import("./components/intelliJ/PerformanceDevDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests On Fast Installer Dashboard"},
        },
        {
          path: ROUTES.IntelliJSharedIndicesDashboard,
          component: () => import("./components/intelliJ/SharedIndicesPerformanceDashboard.vue"),
          meta: {pageTitle: "IntelliJ Integration Performance Tests For Shared Indices Dashboard"},
        },
        {
          path: ROUTES.IntelliJTests,
          component: () => import("./components/intelliJ/PerformanceTests.vue"),
          meta: {pageTitle: "IntelliJ Performance tests"},
        },
        {
          path: ROUTES.IntelliJDevTests,
          component: () => import("./components/intelliJ/PerformanceTestsDev.vue"),
          meta: {pageTitle: "IntelliJ Integration Performance Tests On Fast Installer"},
        },
        {
          path: ROUTES.IntelliJSharedIndicesTests,
          component: () => import("./components/intelliJ/SharedIndicesTests.vue"),
          meta: {pageTitle: "IntelliJ Integration Performance Tests For Shared Indices"},
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
          component: () => import("./components/phpstorm/PerformanceTestsWithPlugins.vue"),
          meta: {pageTitle: "PhpStorm Performance tests with plugins"},
        },
        {
          path: ROUTES.PhpStormTests,
          component: () => import("./components/phpstorm/PerformanceTests.vue"),
          meta: {pageTitle: "PhpStorm Performance tests"},
        },
        {
          path: ROUTES.GoLandDashboard,
          component: () => import("./components/goland/PerformanceDashboard.vue"),
          meta: {pageTitle: "GoLand Performance dashboard"},
        },
        {
          path: ROUTES.GoLandTests,
          component: () => import("./components/goland/PerformanceTests.vue"),
          meta: {pageTitle: "GoLand Performance tests"},
        },

        {
          path: ROUTES.RubyMineDashboard,
          component: () => import("./components/rubymine/PerformanceDashboard.vue"),
          meta: {pageTitle: "RubyMine Performance dashboard"},
        },
        {
          path: ROUTES.RubyMineTests,
          component: () => import("./components/rubymine/PerformanceTests.vue"),
          meta: {pageTitle: "RubyMine Performance tests"},
        },

        {
          path: ROUTES.KotlinExplore,
          component: () => import("./components/kotlin/KotlinExplore.vue"),
          meta: {pageTitle: "Kotlin Performance tests explore"},
        },
        {
          path: ROUTES.KotlinExploreDev,
          component: () => import("./components/kotlin/dev/KotlinDevExplore.vue"),
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
          path: ROUTES.KotlinCompletion,
          component: () => import("./components/kotlin/KotlinCompletionDashboard.vue"),
          meta: {pageTitle: "Kotlin completion"},
        },
        {
          path: ROUTES.KotlinCompletionDev,
          component: () => import("./components/kotlin/dev/KotlinCompletionDevDashboard.vue"),
          meta: {pageTitle: "Kotlin completion (dev/fast)"},
        },
        {
          path: ROUTES.RustTests,
          component: () => import("./components/rust/PerformanceTests.vue"),
          meta: {pageTitle: "Rust Performance tests"},
        },
        {
          path: ROUTES.ScalaTests,
          component: () => import("./components/scala/PerformanceTests.vue"),
          meta: {pageTitle: "Scala Performance tests"},
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
          component: () => import("./components/fleet/PerformanceTests.vue"),
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
