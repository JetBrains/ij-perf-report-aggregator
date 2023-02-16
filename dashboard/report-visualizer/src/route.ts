import { MenuItem } from "primevue/menuitem"
import { ParentRouteRecord } from "shared/src/route"

export function getReportVisualizerItems(): Array<MenuItem> {
  return [
    {
      label: "Report Analyzer",
      to: "/report"
    },
  ]
}

export function getReportVisualizerRoutes(): Array<ParentRouteRecord> {
  return [
    {
      children: [
        {
          path: "/report",
          component: () => import("./Report.vue"),
          meta: {pageTitle: "Report Analyzer"},
        },
      ],
    },
  ]
}
