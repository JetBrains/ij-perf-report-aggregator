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
          to: "/old/intellij/performanceTests",
        },
        {
          label: "Performance Tests (Fast Installer and Dev Server)",
          to: "/old/intellij/performanceTestsFastInstallerDevServer",
        },
        {
          label: "Performance Dashboard",
          to: "/old/intellij/dashboard",
        },
        {
          label: "Performance Dashboard (New)",
          to: "/",
        },
        {
          label: "Performance Dashboard (Fast Installer and Dev Server)",
          to: "/old/intellij/dashboardFastInstallerDevServer",
        },
        {
          label: "Gradle Import Dashboard",
          to: "/old/intellij/gradleImportDashboard",
        },
        {
          label: "Maven Import Dashboard",
          to: "/old/intellij/mavenImportDashboard",
        },
        {
          label: "Shared Indexes",
          to: "/old/intellij/sharedIndexes",
        },
        {
          label: "With Rust Plugin",
          to: "/old/intellij/rust/performanceTests",
        },
        {
          label: "Scala Plugin",
          to: "/old/intellij/scala/performanceTests",
        },
      ]
    },
    {
      label: "PhpStorm",
      items: [
        {
          label: "Dashboard",
          to: "/old/phpstorm/dashboard",
        },
        {
          label: "Explore",
          to: "/old/phpstorm/performanceTests",
        },
        {
          label: "Dashboard With Plugins",
          to: "/old/phpstorm/dashboardWithPlugins",
        },
        {
          label: "Explore With Plugins",
          to: "/old/phpstorm/performanceTestsWithPlugins",
        },
      ]
    },
    {
      label: "RubyMine",
      items: [
        {
          label: "Dashboard",
          to: "/old/rubymine/dashboard",
        },
        {
          label: "Explore",
          to: "/old/rubymine/performanceTests",
        },
      ]
    },
    {
      label: "GoLand",
      items: [
        {
          label: "Dashboard",
          to: "/old/goland/dashboard",
        },
        {
          label: "Explore",
          to: "/old/goland/performanceTests",
        },
      ]
    },
    {
      label: "DataGrip",
      items: [
        {
          label: "Explore",
          to: "/old/datagrip/performanceTests",
        },
      ]
    },
    {
      label: "Fleet",
      items: [
        {
          to: "/old/fleet/dashboard",
          label: "Startup Dashboard",
        },
        {
          to: "/old/fleet/perf/dashboard",
          label: "Performance Dashboard",
        },
        {
          to: "/old/fleet/perf",
          label: "Performance Tests",
        },
        {
          to: "/old/fleet/explore",
          label: "Explore",
        },
      ],
    },
    {
      label: "Kotlin",
      items: [
        {
          label: "Dashboard",
          to: "/old/kotlin/pluginDashboard",
        },
        {
          label: "Dashboard (Fast installer/Dev)",
          to: "/old/kotlin/pluginDashboardFastOrDev",
        },
        {
          label: "Explore",
          to: "/old/kotlin/performanceKotlinPluginTests",
        },
         {
          label: "Explore (Dev)",
          to: "/old/kotlin/performanceKotlinPluginTestsDev",
        },
        {
          label: "Build kts",
          to: "/old/kotlin/buildScript"
        },
        {
          label: "MPP projects",
          to: "/old/kotlin/mppProjects"
        },
      ]
    },
    {
      label: "JBR",
      items: [
        {
          label: "Explore",
          to: "/old/jbr/performanceTests",
        },
      ]
    },
    {
      label: "Aggregates",
      items: [
        {
          label: "PhpStorm",
          to: "/old/aggregates/phpstorm",
        },
        {
          label: "IDEA",
          to: "/old/aggregates/idea",
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
          path: "/old/intellij/sharedIndexes",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "ideaSharedIndices",
            defaultMeasures: [],
          },
          meta: {pageTitle: "IntelliJ Shared Indexes"},
        },
        {
          path: "/old/intellij/performanceTests",
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
          path: "/old/intellij/performanceTestsFastInstallerDevServer",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfintDev",
            table: "idea",
            defaultMeasures: [],
            installerExists: false,
          },
          meta: {pageTitle: "IntelliJ Integration Performance Tests On Fast Installer and Dev Server"},
        },
        {
          path: "/old/intellij/dashboardFastInstallerDevServer",
          component: () => import("./idea/IdeaPerformanceFastInstallerDevServerDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests On Fast Installer and Dev Server Dashboard"},
        },
        {
          path: "/old/intellij/dashboard",
          component: () => import("./idea/IdeaPerformanceDashboard.vue"),
          meta: {pageTitle: "IntelliJ Performance Tests Dashboard"},
        },
        {
          path: "/old/intellij/gradleImportDashboard",
          component: () => import("./idea/GradleImportPerformanceDashboard.vue"),
          meta: {pageTitle: "Performance Tests on Gradle Import"},
        },
        {
          path: "/old/intellij/mavenImportDashboard",
          component: () => import("./idea/MavenImportPerformanceDashboard.vue"),
          meta: {pageTitle: "Performance Tests on Maven Import"},
        },
        {
          path: "/old/intellij/rust/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "rust",
            defaultMeasures: [],
          },
          meta: {pageTitle: "IntelliJ with Rust Plugin"},
        },
        {
          path: "/old/intellij/scala/performanceTests",
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
          path: "/old/rubymine/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "ruby",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore RubyMine Tests"},
        },
        {
          path: "/old/rubymine/dashboard",
          component: () => import("./rubymine/RubyMineDashboard.vue"),
          meta: {pageTitle: "RubyMine Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/old/kotlin/performanceKotlinPluginTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "kotlin",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore Kotlin plugin Tests"},
        },
        {
          path: "/old/kotlin/performanceKotlinPluginTestsDev",
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
          path: "/old/kotlin/pluginDashboard",
          component: () => import("./kotlin/KotlinPluginDashboard.vue"),
          meta: {pageTitle: "Kotlin plugin Dashboard"},
        },
        {
          path: "/old/kotlin/pluginDashboardFastOrDev",
          component: () => import("./kotlin/KotlinPluginDashboardFastInstallerOrDev.vue"),
          props: {
            installerExists: false,
          },
          meta: {pageTitle: "Kotlin plugin Dashboard (Fast installer/Dev)"},
        },
        {
          path: "/old/kotlin/buildScript",
          component: () => import("./kotlin/KotlinBuildScriptDashboard.vue"),
          props: {
            installerExists: false,
          },
          meta: {pageTitle: "Kotlin build kts dashboard"},
        },
        {
          path: "/old/kotlin/mppProjects",
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
          path: "/old/datagrip/performanceTests",
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
          path: "/old/goland/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "goland",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore GoLand Tests"},
        },
        {
          path: "/old/goland/dashboard",
          component: () => import("./goland/GolandDashboard.vue"),
          meta: {pageTitle: "GoLand Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/old/phpstorm/performanceTests",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore PhpStorm Tests"},
        },
        {
          path: "/old/phpstorm/dashboard",
          component: () => import("./phpstorm/PhpStormDashboard.vue"),
          meta: {pageTitle: "PhpStorm Dashboard"},
        },
        {
          path: "/old/phpstorm/performanceTestsWithPlugins",
          component: () => import("shared/src/components/GenericMetricDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "phpstormWithPlugins",
            defaultMeasures: [],
          },
          meta: {pageTitle: "Explore PhpStorm Tests With Plugins"},
        },
        {
          path: "/old/phpstorm/dashboardWithPlugins",
          component: () => import("./phpstorm/PhpStormDashboardWithPlugins.vue"),
          meta: {pageTitle: "PhpStorm With Plugins Dashboard"},
        },
      ]
    },
    {
      children: [
        {
          path: "/old/fleet/dashboard",
          meta: {pageTitle: "Fleet Dashboard", menuTitle: "Dashboard"},
          component: () => import("./fleet/FleetDashboard.vue"),
        },
        {
          path: "/old/fleet/perf",
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
          path: "/old/fleet/perf/dashboard",
          component: () => import("./fleet/FleetPerformanceDashboard.vue"),
          meta: {pageTitle: "Fleet Performance Dashboard"},
        },
        {
          path: "/old/fleet/explore",
          component: () => import("./fleet/FleetExplore.vue"),
          meta: {pageTitle: "Fleet Explore"},
        },
      ],
    },
    {
      children: [
        {
          path: "/old/jbr/performanceTests",
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
          path: "/old/aggregates/phpstorm",
          component: () => import("shared/src/components/GenericAggregatedDashboard.vue"),
          props: {
            dbName: "perfint",
            table: "phpstorm"
          },
          meta: {pageTitle: "PhpStorm Aggregated Dashboard"},
        },
        {
          path: "/old/aggregates/idea",
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
