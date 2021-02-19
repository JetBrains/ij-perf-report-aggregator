import { createApp, nextTick } from "vue"
import App from "./App.vue"
import "element-plus/lib/theme-chalk/el-reset.css"
import "element-plus/lib/theme-chalk/index.css"
import ElButton from "element-plus/es/el-button"
import ElCard from "element-plus/es/el-card"
import ElCol from "element-plus/es/el-col"
import ElContainer from "element-plus/es/el-container"
import ElForm from "element-plus/es/el-form"
import ElFormItem from "element-plus/es/el-form-item"
import ElIcon from "element-plus/es/el-icon"
import ElImage from "element-plus/es/el-image"
import ElInput from "element-plus/es/el-input"
import ElOption from "element-plus/es/el-option"
import ElOptionGroup from "element-plus/es/el-option-group"
import ElRow from "element-plus/es/el-row"
import ElSelect from "element-plus/es/el-select"
import ElDivider from "element-plus/es/el-divider"
import { createRouter, createWebHistory } from "vue-router"

const components = [
  ElButton,
  ElCard,
  ElCol,
  ElContainer,
  ElForm,
  ElFormItem,
  ElIcon,
  ElImage,
  ElInput,
  ElOption,
  ElOptionGroup,
  ElRow,
  ElSelect,
  ElDivider,
]

const app = createApp(App)

for (const component of components) {
  app.component(component.name, component)
}

const router = createRouter({
  history: createWebHistory("/v2/"),
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
