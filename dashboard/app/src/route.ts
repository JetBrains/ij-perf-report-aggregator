import { getNewDashboardRoutes } from "new-dashboard/src/routes"
import { ParentRouteRecord } from "shared/src/route"
import { nextTick } from "vue"
import { createRouter, createWebHistory, Router, RouteRecordRaw } from "vue-router"

function addRoutes(routes: Array<ParentRouteRecord>, result: Array<RouteRecordRaw>) {
  for (const route of routes) {
    result.push(...route.children)
  }
}

export function createAndConfigureRouter(): Router {
  const routes: Array<RouteRecordRaw> = [
    {
      path: "",
      redirect: "/intellij/dashboard",
    },
    {
      path: "/:catchAll(.*)",
      name: "Page Not Found",
      component: () => import("shared/src/components/PageNotFound.vue"),
    },
  ]

  addRoutes(getNewDashboardRoutes(), routes)

  const router = createRouter({
    history: createWebHistory("/"),
    routes,
  })
  router.afterEach((to, _from) => {
    void nextTick(() => {
      document.title = to.meta["pageTitle"] as string ?? ""
    })
  })
  return router
}