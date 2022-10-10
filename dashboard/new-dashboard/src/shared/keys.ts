import { InjectionKey } from "vue"
import { InfoSidebarVm } from "../components/InfoSidebarVm"

export const sidebarVmKey: InjectionKey<InfoSidebarVm> = Symbol("sidebarVm")