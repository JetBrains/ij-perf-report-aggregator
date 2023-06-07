import { InjectionKey, Ref } from "vue"
import { InfoSidebarVm } from "../components/InfoSidebarVm"
import { DataQueryConfigurator } from "../components/common/dataQuery"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { FilterConfigurator } from "../configurators/filter"
import { Accident } from "../util/meta"

export const sidebarVmKey: InjectionKey<InfoSidebarVm> = Symbol("sidebarVm")
export const containerKey: InjectionKey<Ref<HTMLElement | undefined>> = Symbol("chartContainerKey")

export const serverConfiguratorKey: InjectionKey<ServerConfigurator> = Symbol("serverConfiguratorKey")
export const accidentsKeys: InjectionKey<Ref<Accident[]>> = Symbol("accidentsKey")
export const dashboardConfiguratorsKey: InjectionKey<DataQueryConfigurator[] | FilterConfigurator[]> = Symbol("dashboardConfiguratorsKey")
