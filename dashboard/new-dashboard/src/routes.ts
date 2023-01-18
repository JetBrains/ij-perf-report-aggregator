import { ParentRouteRecord } from "shared/src/route"

export interface NavigationItem {
  path: string
  name: string
  key?: string
}

const enum ROUTE_PREFIX {
  InteliJ = "/new/ij",
  PhpStorm = "/new/phpstorm",
  GoLand = "/new/goland",
  RubyMine = "/new/rubymine",
  Kotlin = "/new/kotlin",
  Rust = "/new/rust",
  Scala = "/new/scala",
  JBR = "/new/jbr",
  Fleet = "/new/fleet",
}

const ROUTES = {
  InteliJDashboard: `${ROUTE_PREFIX.InteliJ}/dashboard`,
  InteliJGradleDashboard: `${ROUTE_PREFIX.InteliJ}/gradleDashboard`,
  InteliJMavenDashboard: `${ROUTE_PREFIX.InteliJ}/mavenDashboard`,
  InteliJDevDashboard: `${ROUTE_PREFIX.InteliJ}/devDashboard`,
  InteliJTests: `${ROUTE_PREFIX.InteliJ}/tests`,
  InteliJDevTests: `${ROUTE_PREFIX.InteliJ}/devTests`,
  InteliJSharedIndicesTests: `${ROUTE_PREFIX.InteliJ}/sharedIndicesTests`,
  InteliJCompare: `${ROUTE_PREFIX.InteliJ}/compare`,

  PhpStormDashboard: `${ROUTE_PREFIX.PhpStorm}/dashboard`,
  PhpStormWithPluginsDashboard: `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests: `${ROUTE_PREFIX.PhpStorm}/tests`,
  PhpStormWithPluginsTests: `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,

  KotlinDashboard: `${ROUTE_PREFIX.Kotlin}/dashboard`,
  KotlinDashboardDev: `${ROUTE_PREFIX.Kotlin}/dashboardDev`,
  KotlinExplore: `${ROUTE_PREFIX.Kotlin}/explore`,
  KotlinExploreDev: `${ROUTE_PREFIX.Kotlin}/exploreDev`,

  GoLandDashboard: `${ROUTE_PREFIX.GoLand}/dashboard`,
  GoLandTests: `${ROUTE_PREFIX.GoLand}/tests`,

  RubyMineDashboard: `${ROUTE_PREFIX.RubyMine}/dashboard`,
  RubyMineTests: `${ROUTE_PREFIX.RubyMine}/tests`,

  RustTests: `${ROUTE_PREFIX.Rust}/tests`,
  ScalaTests: `${ROUTE_PREFIX.Scala}/tests`,
  JBRTests: `${ROUTE_PREFIX.JBR}/tests`,
  FleetTest: `${ROUTE_PREFIX.Fleet}/tests`,
  FleetPerfDashboard: `${ROUTE_PREFIX.Fleet}/perfDashboard`,
  FleetStartupDashboard: `${ROUTE_PREFIX.Fleet}/startupDashboard`,
}

export const topNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.InteliJDashboard,
    name: "InteliJ",
    key: ROUTE_PREFIX.InteliJ,
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
  },
  {
    path: "/",
    name: "Back to old",
  },
]

export const intelijTabNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.InteliJDashboard,
    name: "Performance dashboard",
  },
  {
    path: ROUTES.InteliJGradleDashboard,
    name: "Gradle Import dashboard",
  },
  {
    path: ROUTES.InteliJMavenDashboard,
    name: "Maven Import dashboard",
  },
  {
    path: ROUTES.InteliJDevDashboard,
    name: "Performance dashboard (Fast Installer)",
  },
  {
    path: ROUTES.InteliJTests,
    name: "Performance tests",
  },
  {
    path: ROUTES.InteliJDevTests,
    name: "Performance Tests (Fast Installer)",
  },
  {
    path: ROUTES.InteliJSharedIndicesTests,
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
    path: ROUTES.KotlinDashboardDev,
    name: "Dashboard (dev/fast installer)",
  },
  {
    path: ROUTES.KotlinExplore,
    name: "Explore",
  },
  {
    path: ROUTES.KotlinExploreDev,
    name: "Explore (dev/fast installer)",
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
    path: ROUTES.JBRTests,
    name: "Tests",
  },
]

export const FleetNavigationItems: NavigationItem[] = [
  {
    path: ROUTES.FleetStartupDashboard,
    name: "Starup Dashboard",
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
  if (path.startsWith(ROUTE_PREFIX.InteliJ)) {
    return intelijTabNavigationItems
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
          path: ROUTES.InteliJDashboard,
          component: () => import("./components/inteliJ/PerformanceDashboard.vue"),
          meta: {pageTitle: "InteliJ Performance dashboard"},
        },
        {
          path: ROUTES.InteliJGradleDashboard,
          component: () => import("./components/inteliJ/GradleImportPerformanceDashboard.vue"),
          meta: {pageTitle: "Gradle Import dashboard"},
        },
        {
          path: ROUTES.InteliJMavenDashboard,
          component: () => import("./components/inteliJ/MavenImportPerformanceDashboard.vue"),
          meta: {pageTitle: "Maven Import dashboard"},
        },
        {
          path: ROUTES.InteliJDevDashboard,
          component: () => import("./components/inteliJ/PerformanceDevDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests On Fast Installer Dashboard"},
        },
        {
          path: ROUTES.InteliJTests,
          component: () => import("./components/inteliJ/PerformanceTests.vue"),
          meta: {pageTitle: "InteliJ Performance tests"},
        },
        {
          path: ROUTES.InteliJDevTests,
          component: () => import("./components/inteliJ/PerformanceTestsDev.vue"),
          meta: {pageTitle: "IntelliJ Integration Performance Tests On Fast Installer"},
        },
        {
          path: ROUTES.InteliJSharedIndicesTests,
          component: () => import("./components/inteliJ/SharedIndicesTests.vue"),
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
          component: () => import("./components/kotlin/KotlinDevExplore.vue"),
          meta: {pageTitle: "Kotlin Performance tests explore (dev/fast installer)"},
        },
        {
          path: ROUTES.KotlinDashboard,
          component: () => import("./components/kotlin/PerformanceDashboard.vue"),
          meta: {pageTitle: "Kotlin Performance dashboard"},
        },
        {
          path: ROUTES.KotlinDashboardDev,
          component: () => import("./components/kotlin/PerformanceDevDashboard.vue"),
          meta: {pageTitle: "Kotlin Performance dashboard (dev/fast installer)"},
        },
        {
          path: ROUTES.RustTests,
          component: () => import("./components/rust/PerformanceTests.vue"),
          meta: {pageTitle: "Rust Performance tests"},
        },
        {
          path: ROUTES.RustTests,
          component: () => import("./components/scala/PerformanceTests.vue"),
          meta: {pageTitle: "Scala Performance tests"},
        },
        {
          path: ROUTES.JBRTests,
          component: () => import("./components/jbr/PerformanceTests.vue"),
          meta: {pageTitle: "JBR Performance tests"},
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
