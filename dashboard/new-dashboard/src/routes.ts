import { ParentRouteRecord } from "shared/src/route"

export function getDashboardMenuItems() {
  return [{
    path: "/new/ij/dashboard",
    name: "InteliJ"
  }]
}

export function getNewDashboardRoutes(): ParentRouteRecord[] {
  return [
    {
      children: [
        {
          path: "/new/ij/dashboard",
          name: "InteliJ",
          component: () => import("./components/IntelliJMainDashboard.vue"),
        },
      ]
    },
  ]
}