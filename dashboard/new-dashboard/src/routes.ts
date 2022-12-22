import { ParentRouteRecord } from "shared/src/route"

export interface NavigationItems {
  path: string
  name: string
}

const enum ROUTE_PREFIX {
  InteliJ = "/new/ij",
  PhpStorm = "/new/phpstorm",
}

const ROUTES = {
  InteliJDashboard:`${ROUTE_PREFIX.InteliJ}/dashboard`,
  InteliJTests: `${ROUTE_PREFIX.InteliJ}/tests`,
  InteliJCompare: `${ROUTE_PREFIX.InteliJ}/compare`,

  PhpStormDashboard: `${ROUTE_PREFIX.PhpStorm}/dashboard`,
  PhpStormWithPluginsDashboard: `${ROUTE_PREFIX.PhpStorm}/pluginsDashboard`,
  PhpStormTests: `${ROUTE_PREFIX.PhpStorm}/tests`,
  PhpStormWithPluginsTests: `${ROUTE_PREFIX.PhpStorm}/testsWithPlugins`,
}

export const topNavigationItems: NavigationItems[] = [
  {
    path: ROUTES.InteliJDashboard,
    name: "InteliJ",
  },
  {
    path: ROUTES.PhpStormDashboard,
    name: "PhpStorm",
  },
  {
    path: "/",
    name: "Back to old",
  },
]

export const intelijTabNavigationItems: NavigationItems[] = [
  {
    path: ROUTES.InteliJDashboard,
    name: "Performance dashboard",
  },
  {
    path: ROUTES.InteliJTests,
    name: "Performance tests",
  },
  // {
  //   path: ROUTES.Compare,
  //   name: "Compare branches",
  // },
]

export const phpStormNavigationItems: NavigationItems[] = [
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

export function getNavigationTabs(path: string): NavigationItems[] {
  if (path.startsWith(ROUTE_PREFIX.InteliJ)) {
    return intelijTabNavigationItems
  }

  if (path.startsWith(ROUTE_PREFIX.PhpStorm)) {
    return phpStormNavigationItems
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
          path: ROUTES.InteliJTests,
          component: () => import("./components/inteliJ/PerformanceTests.vue"),
          meta: {pageTitle: "InteliJ Performance tests"},
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
      ],
    },
  ]
}
