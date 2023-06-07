import { ParentRouteRecord } from "new-dashboard/src/components/common/route"
import { getNewDashboardRoutes } from "new-dashboard/src/routes"
import { nextTick } from "vue"
import { createRouter, createWebHistory, Router, RouteRecordRaw } from "vue-router"

function addRoutes(routes: ParentRouteRecord[], result: RouteRecordRaw[]) {
  for (const route of routes) {
    result.push(...route.children)
  }
}

export function createAndConfigureRouter(): Router {
  const routes: RouteRecordRaw[] = [
    {
      path: "",
      redirect: "/intellij/dashboard",
    },
    {
      path: "/:catchAll(.*)",
      name: "Page Not Found",
      component: () => import("new-dashboard/src/components/charts/PageNotFound.vue"),
    },
  ]

  addRoutes(getNewDashboardRoutes(), routes)

  const router = createRouter({
    history: createWebHistory("/"),
    routes,
  })
  router.afterEach((to, _from) => {
    void nextTick(() => {
      document.title = (to.meta["pageTitle"] as string | null) ?? ""
    })
  })
  return router
}
