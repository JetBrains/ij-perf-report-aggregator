import { InjectionKey, Ref } from "vue"
import { DataQueryConfigurator } from "../components/common/dataQuery"
import { InfoSidebar } from "../components/common/sideBar/InfoSidebar"
import { AccidentsConfigurator } from "../configurators/accidents/AccidentsConfigurator"
import { BranchConfigurator } from "../configurators/BranchConfigurator"
import { ServerWithCompressConfigurator } from "../configurators/ServerWithCompressConfigurator"
import { FilterConfigurator } from "../configurators/filter"
import { YoutrackClient } from "../components/common/youtrack/YoutrackClient"
import type { CompareSectionsRegistry, RenderMode } from "../components/charts/compareMode"

export const sidebarVmKey: InjectionKey<InfoSidebar> = Symbol("sidebarVm")
export const containerKey: InjectionKey<Ref<HTMLElement | null>> = Symbol("chartContainerKey")

export const serverConfiguratorKey: InjectionKey<ServerWithCompressConfigurator> = Symbol("serverConfiguratorKey")
export const accidentsConfiguratorKey: InjectionKey<AccidentsConfigurator> = Symbol("accidentsKey")
export const dashboardConfiguratorsKey: InjectionKey<DataQueryConfigurator[] | FilterConfigurator[]> = Symbol("dashboardConfiguratorsKey")
export const branchConfiguratorKey: InjectionKey<BranchConfigurator | null> = Symbol("branchConfiguratorKey")
export const renderModeKey: InjectionKey<Ref<RenderMode>> = Symbol("renderModeKey")
export const compareSectionsRegistryKey: InjectionKey<CompareSectionsRegistry> = Symbol("compareSectionsRegistryKey")
export const youtrackClientKey: InjectionKey<YoutrackClient> = Symbol("youtrackClientKey")
