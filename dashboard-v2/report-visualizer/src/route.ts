import { ParentRouteRecord } from "shared/src/route"

export function getReportVisualizerRoutes(): Array<ParentRouteRecord> {
  return [
    {
      title: null,
      children: [
        {
          path: "/report",
          component: () => import("./Report.vue"),
          meta: {pageTitle: "Report Analyzer", menuTitle: "Report Analyzer"},
        },
      ],
    },
  ]
}
