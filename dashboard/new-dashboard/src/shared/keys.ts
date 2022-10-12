import { InjectionKey, Ref } from "vue"
import { InfoSidebarVm } from "../components/InfoSidebarVm"

export const sidebarVmKey: InjectionKey<InfoSidebarVm> = Symbol("sidebarVm")
export const containerKey: InjectionKey<Ref<HTMLElement>> = Symbol("chartContainerKey")