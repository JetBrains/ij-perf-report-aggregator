import { ParentRouteRecord } from "shared/src/route"

export function getIjRoutes(): Array<ParentRouteRecord> {
  return [
    {
      title: "IJ Dashboard",
      children: [
        {
          path: "/ij",
          component: () => import("./IntelliJDashboard.vue"),
          children: [
            {
              path: "/ij/dashboard",
              redirect: "/ij/pulse",
            },
            {
              path: "/ij/pulse",
              component: () => import("./Pulse.vue"),
              meta: {pageTitle: "IJ - Pulse", menuTitle: "Pulse"},
            },
            {
              path: "/ij/progressOverTime",
              component: () => import("./ProgressOverTime.vue"),
              meta: {pageTitle: "IJ - Progress Over Time", menuTitle: "Progress Over Time"},
            },
            {
              path: "/ij/moduleLoading",
              component: () => import("./ModuleLoading.vue"),
              meta: {pageTitle: "IJ - Module Loading", menuTitle: "Module Loading"},
            },
          ],
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
