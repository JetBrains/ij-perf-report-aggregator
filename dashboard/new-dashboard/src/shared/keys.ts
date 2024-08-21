import { InjectionKey, Ref } from "vue"
import { DataQueryConfigurator } from "../components/common/dataQuery"
import { InfoSidebar } from "../components/common/sideBar/InfoSidebar"
import { AccidentsConfigurator } from "../configurators/AccidentsConfigurator"
import { ServerWithCompressConfigurator } from "../configurators/ServerWithCompressConfigurator"
import { FilterConfigurator } from "../configurators/filter"
import { YoutrackClient } from "../components/common/youtrack/YoutrackClient"

export const sidebarVmKey: InjectionKey<InfoSidebar> = Symbol("sidebarVm")
export const containerKey: InjectionKey<Ref<HTMLElement | undefined>> = Symbol("chartContainerKey")

export const serverConfiguratorKey: InjectionKey<ServerWithCompressConfigurator> = Symbol("serverConfiguratorKey")
export const accidentsConfiguratorKey: InjectionKey<AccidentsConfigurator> = Symbol("accidentsKey")
export const dashboardConfiguratorsKey: InjectionKey<DataQueryConfigurator[] | FilterConfigurator[]> = Symbol("dashboardConfiguratorsKey")
export const youtrackClientKey: InjectionKey<YoutrackClient> = Symbol("youtrackClientKey")
