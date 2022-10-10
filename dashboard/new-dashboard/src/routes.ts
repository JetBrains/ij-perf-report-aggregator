import { ParentRouteRecord } from "shared/src/route"

export function getNewDashboardRoutes(): ParentRouteRecord {
  return {
    children: [
      {
        path: "/dashboard/ij",
        name: "InteliJ",
        component: () => import("./components/IntelliJMainDashboard.vue"),
      },
    ]
  }
}