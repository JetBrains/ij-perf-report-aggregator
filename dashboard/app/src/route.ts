import { ParentRouteRecord } from "new-dashboard/src/components/common/route"
import { getNewDashboardRoutes } from "new-dashboard/src/routes"
import { nextTick } from "vue"
import { createRouter, createWebHistory, RouteLocationNormalized, Router, RouteRecordRaw } from "vue-router"
import { SCROLL_RESTORATION_RESTORE_DEFAULT, ScrollRestorationSessionEntry } from "vue-router/experimental"

function addRoutes(routes: ParentRouteRecord[], result: RouteRecordRaw[]) {
  for (const route of routes) {
    result.push(...route.children)
  }
}

export function createAndConfigureRouter(): Router {
  const routes: RouteRecordRaw[] = [
    {
      path: "",
      redirect: "/intellij/product-metrics",
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

export function restoreScroll(entry: ScrollRestorationSessionEntry | null | undefined, to: RouteLocationNormalized) {
  if (to.hash) {
    setTimeout(() => {
      const element = document.querySelector(to.hash)
      const yOffset = -60 // Adjust this value as needed for your fixed header or other elements
      const y = (element?.getBoundingClientRect().top ?? 0) + window.scrollY + yOffset
      window.scrollTo({ top: y, behavior: "smooth" })
    }, 600)
    return
  }
  // restore runs from App.vue (outside RouterView) before the new route renders, so wait for it
  void nextTick(() => {
    SCROLL_RESTORATION_RESTORE_DEFAULT(entry, to)
  })
}
