import { createApp, nextTick } from "vue"
import App from "./App.vue"
import "element-plus/lib/theme-chalk/el-reset.css"
import "element-plus/lib/theme-chalk/index.css"
import {
  ElButton,
  ElCard,
  ElCascader,
  ElCascaderPanel,
  ElCheckbox,
  ElCheckboxButton,
  ElCheckboxGroup,
  ElCol,
  ElContainer,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElMenu,
  ElMenuItem,
  ElMenuItemGroup,
  ElSubmenu,
  ElFooter,
  ElForm,
  ElFormItem,
  ElHeader,
  ElIcon,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMain,
  ElOption,
  ElOptionGroup,
  ElPageHeader,
  ElRadio,
  ElRadioButton,
  ElRadioGroup,
  ElRow,
  ElScrollbar,
  ElSelect,
  ElSwitch,
  ElTooltip,
  ElDivider,
} from "element-plus/es"
import { createRouter, createWebHistory } from "vue-router"

const components = [
  ElButton,
  ElCard,
  ElCascader,
  ElCascaderPanel,
  ElCheckbox,
  ElCheckboxButton,
  ElCheckboxGroup,
  ElCol,
  ElContainer,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElMenu,
  ElMenuItem,
  ElMenuItemGroup,
  ElSubmenu,
  ElFooter,
  ElForm,
  ElFormItem,
  ElHeader,
  ElIcon,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMain,
  ElOption,
  ElOptionGroup,
  ElPageHeader,
  ElRadio,
  ElRadioButton,
  ElRadioGroup,
  ElRow,
  ElScrollbar,
  ElSelect,
  ElSwitch,
  ElTooltip,
  ElDivider,
]

const app = createApp(App)

for (const component of components) {
  app.component(component.name, component)
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/ij/dashboard",
      component: () => import("./ij/IntelliJDashboard.vue"),
      name: "IJ Dashboard",
    },
    {
      path: "/ij/explore",
      component: () => import("./ij/IntelliJExplore.vue"),
      name: "IJ Explore",
    },
    {
      path: "",
      redirect: "/ij/dashboard",
    },
    {
      path: "/:catchAll(.*)",
      name: "Page Not Found",
      component: () => import("./components/PageNotFound.vue"),
    }
  ],
})
router.afterEach((to, _from) => {
  // noinspection JSIgnoredPromiseFromCall
  nextTick(() => {
    document.title = to.name as string ?? ""
  })
})

app.use(router)

app.mount("#app")
