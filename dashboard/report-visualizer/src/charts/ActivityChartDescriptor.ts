import { CommonItem } from "../state/data"
import { ChartManager } from "./ChartManager"

export interface ActivityChartDescriptor {
  readonly label: string
  readonly id: string

  readonly isInfoChart?: boolean

  readonly sourceNames?: Array<string>

  readonly rotatedLabels?: boolean
  readonly groupByThread?: boolean
  readonly sourceHasPluginInformation?: boolean

  readonly chartManagerProducer?: (container: HTMLElement, sourceNames: Array<string>, descriptor: ActivityChartDescriptor) => Promise<ChartManager>
  readonly shortNameProducer?: (item: CommonItem) => string
}

export function getShortName(item: { name?: string; n?: string }): string {
  // eslint-disable-next-line
  const name = item.name || item.n!
  const lastDotIndex = name.lastIndexOf(".")
  return lastDotIndex < 0 ? name : name.substring(lastDotIndex + 1)
}

// not as part of ItemChartManager.ts to reduce scope of changes on change
// (make sure that hot reloading will not reload all modules where `chartDescriptors` is used - especially `router`)
export const serviceSourceNames = ["appServices", "projectServices", "moduleServices", "appComponents", "projectComponents", "moduleComponents"]

export const chartDescriptors: Array<ActivityChartDescriptor> = [
  {
    label: "Services",
    id: "services",
    sourceNames: serviceSourceNames,
    shortNameProducer: getShortName,
    async chartManagerProducer(container: HTMLElement, sourceNames: Array<string>, descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      return new ((await import("./ServiceChartManager")).ServiceChartManager)(container, sourceNames, descriptor)
    }
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
    label: "Time Distribution",
    isInfoChart: true,
    id: "timeDistribution",
    sourceNames: [],
    async chartManagerProducer(container: HTMLElement, _sourceNames: Array<string>, _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      return new (await import("./TimeDistributionChartManager")).TimeDistributionChartManager(container)
    },
  },
  {
    label: "Plugin Classes",
    isInfoChart: true,
    id: "pluginClassCount",
    sourceNames: [],
    async chartManagerProducer(container: HTMLElement, _sourceNames: Array<string>, _descriptor: ActivityChartDescriptor): Promise<ChartManager> {
      return new (await import("./PluginClassCountTreeMapChartManager")).PluginClassCountTreeMapChartManager(container)
    },
  },
]
