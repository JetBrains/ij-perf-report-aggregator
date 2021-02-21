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
} from "element-plus"
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
      component: () => import("./IntelliJDashboard.vue"),
      name: "IJ Dashboard",
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
