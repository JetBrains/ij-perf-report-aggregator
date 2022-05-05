import { MenuItem } from "primevue/menuitem"
import { ParentRouteRecord } from "shared/src/route"

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
    {
      label: "Fleet",
      items: [
        {
          to: "/fleet/dashboard",
          label: "Dashboard",
        },
        {
          to: "/fleet/perf",
          label: "Performance Tests",
        },
        {
          to: "/fleet/explore",
          label: "Explore",
        },
      ],
    },
  ]
}

export function getIjRoutes(): Array<ParentRouteRecord> {
  return [
    {
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
              component: () => import("./IntelliJPulse.vue"),
              meta: {pageTitle: "IJ - Pulse"},
            },
            {
              path: "/ij/progressOverTime",
              component: () => import("./ProgressOverTime.vue"),
              meta: {pageTitle: "IJ - Progress Over Time"},
            },
            {
              path: "/ij/moduleLoading",
              component: () => import("./ModuleLoading.vue"),
              meta: {pageTitle: "IJ - Module Loading"},
            },
          ],
        },
        {
          path: "/ij/explore",
          component: () => import("./IntelliJExplore.vue"),
          meta: {pageTitle: "IJ Explore"},
        },
      ]
    },
    {
      children: [
        {
          path: "/sharedIndexes/dashboard",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "sharedIndexes",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Shared Indexes Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/performanceIntegration/dashboard",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Integration Performance Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/rubyMinePerformanceIntegration/dashboard",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "rubymineperfint",
            defaultMeasures: [],
          },
          meta: {pageTitle: "RubyMine Integration Performance Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/fleet/dashboard",
          meta: {pageTitle: "Fleet Dashboard", menuTitle: "Dashboard"},
          component: () => import("./fleet/FleetDashboard.vue"),
        },
        {
          path: "/fleet/perf",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "fleet",
            table: "measure",
            defaultMeasures: [],
            urlEnabled: false,
            compoundTooltip: false,
            chartType: "scatter",
            valueUnit: "ns",
          },
          meta: {pageTitle: "Fleet Performance Tests"},
        },
        {
          path: "/fleet/explore",
          component: () => import("./fleet/FleetExplore.vue"),
          meta: {pageTitle: "Fleet Explore"},
        },
      ],
    },
  ]
}
