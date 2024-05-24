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
    scrollBehavior(to) {
      if (to.hash) {
        return new Promise((resolve) => {
          setTimeout(() => {
            const element = document.querySelector(to.hash)
            const yOffset = -60 // Adjust this value as needed for your fixed header or other elements
            const y = (element?.getBoundingClientRect().top ?? 0) + window.pageYOffset + yOffset
            window.scrollTo({ top: y, behavior: "smooth" })
            resolve({ left: 0, top: y })
          }, 600) // Delay to ensure the element is rendered
        })
      }
      return { left: 0, top: 0 }
    },
  })
  router.afterEach((to, _from) => {
    void nextTick(() => {
      document.title = (to.meta["pageTitle"] as string | null) ?? ""
    })
  })
  return router
}
