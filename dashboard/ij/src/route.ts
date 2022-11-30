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
          label: "Performance Tests",
          to: "/intellij/performanceTests",
        },
        {
          label: "Performance Tests (Dev Server)",
          to: "/intellij/performanceTestsDev",
        },
        {
          label: "Performance Dashboard",
          to: "/intellij/dashboard",
        },
        {
          label: "Performance Dashboard (New)",
          to: "/new",
        },
        {
          label: "Performance Dashboard (Dev Server)",
          to: "/intellij/dashboardDev",
        },
        {
          label: "Import Dashboard",
          to: "/intellij/importDashboard",
        },
        {
          label: "Shared Indexes",
          to: "/intellij/sharedIndexes",
        },
        {
          label: "With Rust Plugin",
          to: "/intellij/rust/performanceTests",
        },
        {
          label: "Scala Plugin",
          to: "/intellij/scala/performanceTests",
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
        {
          label: "Dashboard With Plugins",
          to: "/phpstorm/dashboardWithPlugins",
        },
        {
          label: "Explore With Plugins",
          to: "/phpstorm/performanceTestsWithPlugins",
        },
      ]
    },
    {
      label: "RubyMine",
      items: [
        {
          label: "Dashboard",
          to: "/rubymine/dashboard",
        },
        {
          label: "Explore",
          to: "/rubymine/performanceTests",
        },
      ]
    },
    {
      label: "GoLand",
      items: [
        {
          label: "Dashboard",
          to: "/goland/dashboard",
        },
        {
          label: "Explore",
          to: "/goland/performanceTests",
        },
      ]
    },
    {
      label: "DataGrip",
      items: [
        {
          label: "Explore",
          to: "/datagrip/performanceTests",
        },
      ]
    },
    {
      label: "Fleet",
      items: [
        {
          to: "/fleet/dashboard",
          label: "Startup Dashboard",
        },
        {
          to: "/fleet/perf/dashboard",
          label: "Performance Dashboard",
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
    {
      label: "Kotlin",
      items: [
        {
          label: "Dashboard",
          to: "/kotlin/pluginDashboard",
        },
        {
          label: "Explore",
          to: "/kotlin/performanceKotlinPluginTests",
        },
         {
          label: "Explore (Dev)",
          to: "/kotlin/performanceKotlinPluginTestsDev",
        },
        {
          label: "Build kts",
          to: "/kotlin/buildScript"
        },
        {
          label: "MPP projects",
          to: "/kotlin/mppProjects"
        },
      ]
    },
    {
      label: "JBR",
      items: [
        {
          label: "Explore",
          to: "/jbr/performanceTests",
        },
      ]
    },
    {
      label: "Aggregates",
      items: [
        {
          label: "PhpStorm",
          to: "/aggregates/phpstorm",
        },
        {
          label: "IDEA",
          to: "/aggregates/idea",
        },
      ]
    }
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
          path: "/intellij/sharedIndexes",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "ideaSharedIndices",
            defaultMeasures: [],
          },
          meta: {pageTitle: "IntelliJ Shared Indexes"},
        },
        {
          path: "/intellij/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "idea",
            defaultMeasures: [],
            supportReleases: true,
          },
          meta: {pageTitle: "IntelliJ Integration Performance Tests"},
        },
        {
          path: "/intellij/performanceTestsDev",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfintDev",
            table: "idea",
            defaultMeasures: [],
            installerExists: false,
          },
          meta: {pageTitle: "IntelliJ Integration Performance Tests On Dev Server"},
        },
        {
          path: "/intellij/dashboardDev",
          component: () => import("./idea/IdeaPerformanceDevDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests On Dev Server Dashboard"},
        },
        {
          path: "/intellij/dashboard",
          component: () => import("./idea/IdeaPerformanceDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests Dashboard"},
        },
        {
          path: "/intellij/importDashboard",
          component: () => import("./idea/ImportPerformanceDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests on Import Dashboard"},
        },
        {
          path: "/intellij/rust/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
            defaultMeasures: [],
          },
          meta: {pageTitle: "IntelliJ with Rust Plugin"},
        },
        {
          path: "/intellij/scala/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "scala",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Scala Plugin"},
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
          meta: {pageTitle: "Explore RubyMine Tests"},
        },
        {
          path: "/rubymine/dashboard",
          component: () => import("./rubymine/RubyMineDashboard.vue"),
          meta: {pageTitle: "RubyMine Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/kotlin/performanceKotlinPluginTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore Kotlin plugin Tests"},
        },
        {
          path: "/kotlin/performanceKotlinPluginTestsDev",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfintDev",
            table: "kotlin",
            defaultMeasures: [],
            installerExists: false,
          },
          meta: {pageTitle: "Explore Kotlin plugin Tests (Dev)"},
        },
        {
          path: "/kotlin/pluginDashboard",
          component: () => import("./kotlin/KotlinPluginDashboard.vue"),
          meta: {pageTitle: "Kotlin plugin Dashboard"},
        },
        {
          path: "/kotlin/buildScript",
          component: () => import("./kotlin/KotlinBuildScriptDashboard.vue"),
          props: {
            installerExists: false,
          },
          meta: {pageTitle: "Kotlin build kts dashboard"},
        },
        {
          path: "/kotlin/mppProjects",
          component: () => import("./kotlin/MppProjectsDashboard.vue"),
          props: {
            installerExists: false,
          },
          meta: {pageTitle: "Kotlin MPP projects dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/datagrip/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "datagrip",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore DataGrip Tests"},
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
          meta: {pageTitle: "Explore GoLand Tests"},
        },
        {
          path: "/goland/dashboard",
          component: () => import("./goland/GolandDashboard.vue"),
          meta: {pageTitle: "GoLand Dashboard"},
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
        {
          path: "/phpstorm/performanceTestsWithPlugins",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "phpstormWithPlugins",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore PhpStorm Tests With Plugins"},
        },
        {
          path: "/phpstorm/dashboardWithPlugins",
          component: () => import("./phpstorm/PhpStormDashboardWithPlugins.vue"),
          meta: {pageTitle: "PhpStorm With Plugins Dashboard"},
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
          path: "/fleet/perf/dashboard",
          component: () => import("./fleet/FleetPerformanceDashboard.vue"),
          meta: {pageTitle: "Fleet Performance Dashboard"},
        },
        {
          path: "/fleet/explore",
          component: () => import("./fleet/FleetExplore.vue"),
          meta: {pageTitle: "Fleet Explore"},
        },
      ],
    },
    {
      children: [
        {
          path: "/jbr/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "jbr",
            table: "report",
            installerExists: false,
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore JBR Tests"},
        },
      ]
    },
    {
      children: [
        {
          path: "/aggregates/phpstorm",
          component: () => import("shared/src/components/GenericAggregatedDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm"
          },
          meta: {pageTitle: "PhpStorm Aggregated Dashboard"},
        },
        {
          path: "/aggregates/idea",
          component: () => import("shared/src/components/GenericAggregatedDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "idea"
          },
          meta: {pageTitle: "IDEA Aggregated Dashboard"},
        },
      ],
    },
  ]
}
