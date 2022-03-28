import { MenuItem } from "primevue/menuitem"
import { ParentRouteRecord } from "shared/src/route"
import IntelliJDashboard from "./IntelliJDashboard.vue"
import IntelliJExplore from "./IntelliJExplore.vue"
import ModuleLoading from "./ModuleLoading.vue"
import ProgressOverTime from "./ProgressOverTime.vue"
import Pulse from "./Pulse.vue"
import SharedIndexesDashboard from "./SharedIndexesDashboard.vue"

export function getIjItems(): Array<MenuItem> {
  return [
    {
      label: "IJ",
      items: [
        {
          to: "/ij/pulse",
          label: "Pulse",
        },
        {
          to: "/ij/progressOverTime",
          label: "Progress Over Time",
        },
        {
          to: "/ij/moduleLoading",
          label: "Module Loading",
        },
        {
          to: "/ij/explore",
          label: "Explore",
        },
      ],
    },
    {
      label: "Shared Indexes",
      to: "/sharedIndexes/dashboard",
    },
    {
      label: "Integration Performance",
      to: "/performanceIntegration/dashboard",
    },
    {
      label: "RubyMine Integration Performance",
      to: "/rubyMinePerformanceIntegration/dashboard",
    },
  ]
}

export function getIjRoutes(): Array<ParentRouteRecord> {
  return [
    {
      children: [
        {
          path: "/ij",
          component: () => IntelliJDashboard,
          children: [
            {
              path: "/ij/dashboard",
              redirect: "/ij/pulse",
            },
            {
              path: "/ij/pulse",
              component: () => Pulse,
              meta: {pageTitle: "IJ - Pulse"},
            },
            {
              path: "/ij/progressOverTime",
              component: () => ProgressOverTime,
              meta: {pageTitle: "IJ - Progress Over Time"},
            },
            {
              path: "/ij/moduleLoading",
              component: () => ModuleLoading,
              meta: {pageTitle: "IJ - Module Loading"},
            },
          ],
        },
        {
          path: "/ij/explore",
          component: () => IntelliJExplore,
          meta: {pageTitle: "IJ Explore"},
        },
      ]
    },
    {
      children: [
        {
          path: "/sharedIndexes/dashboard",
          component: () => SharedIndexesDashboard,
          props: {
            dbName: "sharedIndexes",
            defaultMeasures: ["indexing", "scanning"],
          },
          meta: {pageTitle: "Shared Indexes Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/performanceIntegration/dashboard",
          component: () => SharedIndexesDashboard,
          props: {
            dbName: "perfint",
            defaultMeasures: ["indexing", "numberOfIndexedFiles", "numberOfIndexingRuns", "scanning", "updatingTime"],
          },
          meta: {pageTitle: "Integration Performance Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/rubyMinePerformanceIntegration/dashboard",
          component: () => SharedIndexesDashboard,
          props: {
            dbName: "rubymineperfint",
            defaultMeasures: [],
          },
          meta: {pageTitle: "RubyMine Integration Performance Dashboard"},
        },
      ]
    },
  ]
}
