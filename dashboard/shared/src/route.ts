import { RouteRecordRaw } from "vue-router"

export interface ParentRouteRecord {
  children: Array<RouteRecordRaw>
}
