import { MenuItem } from "primevue/menuitem"
import { ParentRouteRecord } from "shared/src/route"

export function getIjItems(): Array<MenuItem> {
  return [
    {
      label: "IJ Startup",
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
      label: "IntelliJ",
      items: [
        {
          label: "Shared Indexes",
          to: "/intellij/sharedIndexes/dashboard",
        },
        {
          label: "Performance Tests",
          to: "/intellij/performanceTests",
        },
        {
          label: "Performance Dashboard",
          to: "/intellij/dashboard",
        },
      ]
    },
    {
      label: "PhpStorm",
      items: [
        {
          label: "Dashboard",
          to: "/phpstorm/dashboard",
        },
        {
          label: "Explore",
          to: "/phpstorm/performanceTests",
        },
      ]
    },
    {
      label: "RubyMine",
      to: "/rubymine/performanceTests",
    },
    {
      label: "GoLand",
      to: "/goland/performanceTests",
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
      children:[
        {
          path: "/intellij/sharedIndexes/dashboard",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "ideaSharedIndices",
            defaultMeasures: [],
          },
          meta: {pageTitle: "IntelliJ Shared Indexes Dashboard"},
        },
        {
          path: "/intellij/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "idea",
            defaultMeasures: [],
          },
          meta: {pageTitle: "IntelliJ Integration Performance Tests"},
        },
        {
          path: "/intellij/dashboard",
          component: () => import("./idea/IdeaPerformanceDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests Dashboard"},
        },
      ],
    },
    {
      children: [
        {
          path: "/rubymine/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "ruby",
            defaultMeasures: [],
          },
          meta: {pageTitle: "RubyMine Integration Performance Tests"},
        },
      ]
    },
    {
      children: [
        {
          path: "/goland/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "goland",
            defaultMeasures: [],
          },
          meta: {pageTitle: "GoLand Integration Performance Tests"},
        },
      ]
    },
    {
      children: [
        {
          path: "/phpstorm/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore PhpStorm Tests"},
        },
        {
          path: "/phpstorm/dashboard",
          component: () => import("./phpstorm/PhpStormDashboard.vue"),
          meta: {pageTitle: "PhpStorm Dashboard"},
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
            installerExists: false,
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
