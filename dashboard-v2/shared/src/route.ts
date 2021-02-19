import { RouteRecordRaw } from "vue-router"

export interface ParentRouteRecord {
  title: string | null
  children: Array<RouteRecordRaw>
}
