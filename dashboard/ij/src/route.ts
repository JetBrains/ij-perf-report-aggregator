import { ParentRouteRecord } from "shared/src/route"

export function getIjRoutes(): Array<ParentRouteRecord> {
  return [
    {
      title: "IJ Dashboard",
      children: [
        {
          path: "/ij/dashboard",
          component: () => import("./IntelliJDashboard.vue"),
          meta: {pageTitle: "IJ Dashboard", menuTitle: "Dashboard"},
        },
        {
          path: "/ij/explore",
          component: () => import("./IntelliJExplore.vue"),
          meta: {pageTitle: "IJ Explore", menuTitle: "Explore"},
        },
      ]
    },
    {
      title: null,
      children: [
        {
          path: "/sharedIndexes/dashboard",
          component: () => import("./SharedIndexesDashboard.vue"),
          meta: {pageTitle: "Shared Indexes Dashboard", menuTitle: "Shared Indexes Dashboard"},
        },
      ]
    },
  ]
}
