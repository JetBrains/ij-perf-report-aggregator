import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { Accident } from "shared/src/meta"
import { InjectionKey, Ref } from "vue"
import { InfoSidebarVm } from "../components/InfoSidebarVm"

export const sidebarVmKey: InjectionKey<InfoSidebarVm> = Symbol("sidebarVm")
export const containerKey: InjectionKey<Ref<HTMLElement|undefined>> = Symbol("chartContainerKey")

export const serverConfiguratorKey: InjectionKey<ServerConfigurator> = Symbol("serverConfiguratorKey")
export const accidentsKeys: InjectionKey<Ref<Array<Accident>>> = Symbol("accidentsKey")