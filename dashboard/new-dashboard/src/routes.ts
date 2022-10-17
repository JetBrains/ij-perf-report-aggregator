import { ParentRouteRecord } from "shared/src/route"

export interface NavigationItems {
  path: string,
  name: string
}

const enum ROUTES {
  Dashboard = "/new/ij/dashboard",
  Tests = "/new/ij/tests",
  Compare = "/new/ij/compare",
}

export const topNavigationItems: NavigationItems[] = [
  {
    path: ROUTES.Dashboard,
    name: "InteliJ",
  },
  {
    path: "/",
    name: "Back to old",
  },
]

export function getDashboardMenuItems() {
  return [{
    path: ROUTES.Dashboard,
    name: "InteliJ",
  },
    {
      path: "/",
      name: "Back to old",
    }
  ]
}

export const tabNavigationItems: NavigationItems[] = [
  {
    path: ROUTES.Dashboard,
    name: "Performance dashboard",
  },
  {
    path: ROUTES.Tests,
    name: "Performance tests",
  },
  // {
  //   path: ROUTES.Compare,
  //   name: "Compare branches",
  // },
]

export function getNewDashboardRoutes(): Array<ParentRouteRecord> {
  return [
    {
      children: [
        {
          path: ROUTES.Dashboard,
          component: () => import("./components/IntelliJMainDashboard.vue"),
          meta: {pageTitle: "InteliJ Performance dashboard"},
        },
        {
          path: ROUTES.Tests,
          component: () => import("./components/IntelliJPerformanceTests.vue"),
          meta: {pageTitle: "InteliJ Performance tests"},
        },
        {
          path: ROUTES.Compare,
          component: () => import("./MainDashboard.vue"),
          meta: {pageTitle: "InteliJ Compare branches"},
        },
      ],
    },
  ]
}
