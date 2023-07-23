import { ChartManager } from "./charts/ChartComponent"
import { CommonItem, UnitConverter } from "./data"

export interface ActivityChartDescriptor {
  readonly label: string
  readonly id: string

  readonly isInfoChart?: boolean

  readonly sourceNames?: string[]

  readonly rotatedLabels?: boolean
  readonly groupByThread?: boolean
  readonly sourceHasPluginInformation?: boolean

  readonly chartManagerProducer?: (container: HTMLElement, sourceNames: string[], descriptor: ActivityChartDescriptor) => Promise<ChartManager>
  readonly shortNameProducer?: (item: CommonItem) => string
}

export function getShortName(item: { name?: string; n?: string }): string {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const name = item.n!
  const lastDotIndex = name.lastIndexOf(".")
  return lastDotIndex < 0 ? name : name.slice(lastDotIndex + 1)
}

// not as part of ItemChartManager.ts to reduce scope of changes on change
// (make sure that hot reloading will not reload all modules where `chartDescriptors` is used - especially `router`)
export const serviceSourceNames = ["appServices", "projectServices", "moduleServices", "appComponents", "projectComponents", "moduleComponents"]

// prettier-ignore
export const chartDescriptors: ActivityChartDescriptor[] = [
  {
    label: "Services",
    id: "services",
    sourceNames: serviceSourceNames,
    shortNameProducer: getShortName,
    async chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      const { ActivityBarChartManager: ActivityBarChartManager } = await import("./charts/ActivityBarChartManager")
      return new ActivityBarChartManager(container, (dataManager) => dataManager.getServiceItems(), {
        unitConverter: UnitConverter.MICROSECONDS,
      })
    },
  },
  {
    label: "Extensions",
    id: "extensions",
    sourceNames: ["appExtensions", "projectExtensions", "moduleExtensions"],
    shortNameProducer: getShortName,
  },
  {
    label: "Options Top Hit Providers",
    id: "topHitProviders",
    sourceNames: ["appOptionsTopHitProviders", "projectOptionsTopHitProviders"],
    shortNameProducer: getShortName,
  },
  {
    label: "Preload",
    id: "preloadActivities",
    shortNameProducer: getShortName,
  },
  {
    label: "Project Post-Startup",
    id: "projectPostStartupActivities",
    shortNameProducer: getShortName,
  },
  {
    label: "Reopening Editors",
    id: "reopeningEditors",
    sourceHasPluginInformation: false,
  },
  {
    label: "GCs",
    id: "GCs",
    rotatedLabels: false,
  },
  {
    label: "Timeline",
    isInfoChart: true,
    id: "timeline",
    async chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      const { TimeLineChartManager: TimeLineChartManager } = await import("./charts/TimeLineChartManager")
      return new TimeLineChartManager(container, {
        hasParent: true,
        dataProvider(dataManager) {
          return [
            {
              category: "items",
              items: dataManager.items,
            },
          ].filter((it) => {
            // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
            return it.items != null
          })
        },
        groupData(item) {
          const name = item.n
          if (name.endsWith(": scheduled")) {
            return "schedule"
          }
          else if (name.endsWith(": completing")) {
            return "completing"
          }
          else {
            return "all"
          }
        },
        sortGroups(names) {
          const w = (n: string) => (n == "all" ? 0 : n == "completing" ? 1 : 2)
          names.sort((a, b) => w(a) - w(b))
        },
        dataDescriptor: {
          unitConverter: UnitConverter.MILLISECONDS,
          threshold: 0,
          shortenName: false,
        },
      })
    },
  },
  {
    label: "Service Timelines",
    isInfoChart: true,
    id: "serviceTimeline",
    async chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      const { TimeLineChartManager: TimeLineChartManager } = await import("./charts/TimeLineChartManager")
      return new TimeLineChartManager(container, {
        hasParent: false,
        groupData: (_, category) => (category == "service waiting" ? "service waiting" : "all"),
        dataDescriptor: {
          unitConverter: UnitConverter.MICROSECONDS,
        },
        sortGroups(names) {
          const w = (n: string) => (n == "all" ? 0 : 1)
          names.sort((a, b) => w(a) - w(b))
        },
        dataProvider(dataManager) {
          return [...dataManager.getServiceItems(), { category: "service waiting", items: dataManager.data.serviceWaiting ?? [] }]
        },
      })
    },
  },
  {
    label: "Time Distribution",
    isInfoChart: true,
    id: "timeDistribution",
    async chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      const { TimeDistributionChartManager: TimeDistributionChartManager } = await import("./charts/TimeDistributionChartManager")
      return new TimeDistributionChartManager(container)
    },
  },
  {
    label: "Plugin Classes",
    isInfoChart: true,
    id: "pluginClassCount",
    async chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      const { PluginClassCountTreeMapChartManager: PluginClassCountTreeMapChartManager } = await import("./charts/PluginClassCountTreeMapChartManager")
      return new PluginClassCountTreeMapChartManager(container)
    },
  },
  {
    label: "Stats",
    isInfoChart: true,
    id: "stats",
    async chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      const { StatsChartManager: StatsChartManager } = await import("./charts/StatsChartManager")
      return new StatsChartManager(container)
    },
  },
]
