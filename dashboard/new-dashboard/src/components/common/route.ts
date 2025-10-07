import { RouteRecordRaw } from "vue-router"

export interface ParentRouteRecord {
  children: RouteRecordRaw[]
}

export interface TypedRouteRecord<Props = Record<string, unknown>> extends Omit<RouteRecordRaw, "props"> {
  props: Props
}
