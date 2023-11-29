import { InjectionKey, Ref } from "vue"
import { DataQueryConfigurator } from "../components/common/dataQuery"
import { InfoSidebar } from "../components/common/sideBar/InfoSidebar"
import { InfoDataPerformance } from "../components/common/sideBar/InfoSidebarPerformance"
import { InfoDataFromStartup } from "../components/common/sideBar/InfoSidebarStartup"
import { AccidentsConfigurator } from "../configurators/AccidentsConfigurator"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { FilterConfigurator } from "../configurators/filter"

export const sidebarVmKey: InjectionKey<InfoSidebar<InfoDataPerformance>> = Symbol("sidebarVm")
export const sidebarStartupKey: InjectionKey<InfoSidebar<InfoDataFromStartup>> = Symbol("sidebarStartup")
export const containerKey: InjectionKey<Ref<HTMLElement | undefined>> = Symbol("chartContainerKey")

export const serverConfiguratorKey: InjectionKey<ServerConfigurator> = Symbol("serverConfiguratorKey")
export const accidentsConfiguratorKey: InjectionKey<AccidentsConfigurator> = Symbol("accidentsKey")
export const dashboardConfiguratorsKey: InjectionKey<DataQueryConfigurator[] | FilterConfigurator[]> = Symbol("dashboardConfiguratorsKey")
