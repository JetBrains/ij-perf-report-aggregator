import { ChartManager } from "./charts/ChartComponent"
import { CommonItem, InputDataV20, UnitConverter } from "./data"
import { TimeLineChartManager } from "./charts/TimeLineChartManager"
import { TimeDistributionChartManager } from "./charts/TimeDistributionChartManager"
import { PluginClassCountTreeMapChartManager } from "./charts/PluginClassCountTreeMapChartManager"
import { StatsChartManager } from "./charts/StatsChartManager"
import { ActivityBarChartManager } from "./charts/ActivityBarChartManager"

export interface ActivityChartDescriptor {
  readonly label: string
  readonly id: string

  readonly isInfoChart?: boolean

  readonly sourceNames?: string[]

  readonly rotatedLabels?: boolean
  readonly groupByThread?: boolean
  readonly sourceHasPluginInformation?: boolean

  readonly chartManagerProducer?: (container: HTMLElement, sourceNames: string[], descriptor: ActivityChartDescriptor) => ChartManager
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

export const chartDescriptors: ActivityChartDescriptor[] = [
  {
    label: "Services",
    id: "services",
    sourceNames: serviceSourceNames,
    shortNameProducer: getShortName,
    chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): ChartManager {
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
    label: "Prepare App Init",
    id: "prepareAppInitActivities",
    groupByThread: true,
    sourceHasPluginInformation: false,
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
    chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): ChartManager {
      return new TimeLineChartManager(
        container,
        (dataManager) => {
          return [
            {
              category: "items",
              items: dataManager.items,
            },
            {
              category: "prepareAppInitActivities",
              items: dataManager.data.prepareAppInitActivities,
            },
          ].filter((it) => {
            // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
            return it.items != null
          })
        },
        {
          unitConverter: UnitConverter.MILLISECONDS,
          threshold: 0,
          shortenName: false,
        }
      )
    },
  },
  {
    label: "Service Timelines",
    isInfoChart: true,
    id: "serviceTimeline",
    chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): ChartManager {
      return new TimeLineChartManager(
        container,
        (dataManager) => {
          return [...dataManager.getServiceItems(), { category: "service waiting", items: (dataManager.data as InputDataV20).serviceWaiting ?? [] }]
        },
        {
          unitConverter: UnitConverter.MICROSECONDS,
        }
      )
    },
  },
  {
    label: "Time Distribution",
    isInfoChart: true,
    id: "timeDistribution",
    chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): ChartManager {
      return new TimeDistributionChartManager(container)
    },
  },
  {
    label: "Plugin Classes",
    isInfoChart: true,
    id: "pluginClassCount",
    chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): ChartManager {
      return new PluginClassCountTreeMapChartManager(container)
    },
  },
  {
    label: "Stats",
    isInfoChart: true,
    id: "stats",
    chartManagerProducer(container: HTMLElement, _sourceNames: string[], _descriptor: ActivityChartDescriptor): ChartManager {
      return new StatsChartManager(container)
    },
  },
]
