import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { InjectionKey, Ref } from "vue"
import { InfoSidebarVm } from "../components/InfoSidebarVm"

export const sidebarVmKey: InjectionKey<InfoSidebarVm> = Symbol("sidebarVm")
export const containerKey: InjectionKey<Ref<HTMLElement|undefined>> = Symbol("chartContainerKey")

export const serverConfiguratorKey: InjectionKey<ServerConfigurator> = Symbol("serverConfigurator")