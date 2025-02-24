import "@fontsource/jetbrains-mono"
import "@fontsource/inter"
import "./main.css"
import { createPinia } from "pinia"
import PrimeVue from "primevue/config"
import ToastService from "primevue/toastservice"
import Tooltip from "primevue/tooltip"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"
// get rid of color.png
// avoid tiff/svg/other deprecated stuff in a final build
import "../theme/primeicons.css"
import { MyPreset } from "../theme/theme"
import "./main.css"

import "../theme/button.css"
import "../theme/menubar.css"
import "../theme/misc.css"
import "../theme/select.css"
import "../theme/select-panel.css"
import "../theme/toolbar.css"
import * as echarts from "echarts"
import { chartTheme } from "../theme/chartTheme"

async function initApp() {
  const app = createApp(App)
  const router = createAndConfigureRouter()
  const pinia = createPinia()
  app.use(router)
  app.use(PrimeVue, {
    theme: {
      preset: MyPreset,
      cssLayer: {
        name: "primevue",
        order: "tailwind-base, primevue, tailwind-utilities",
      },
      options: {
        darkModeSelector: ".dark-mode",
      },
    },
  })
  app.use(ToastService)
  app.use(pinia)
  app.directive("tooltip", Tooltip)

  echarts.registerTheme("chalk", chartTheme)

  await router.isReady().then(() => app.mount("#app"))
}

void initApp()
