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
import "./primeicons.css"
import { MyPreset } from "./theme"
import "./main.css"

import "new-dashboard/src/primevue-theme/button.css"
import "new-dashboard/src/primevue-theme/menubar.css"
import "new-dashboard/src/primevue-theme/misc.css"
import "new-dashboard/src/primevue-theme/select.css"
import "new-dashboard/src/primevue-theme/select-panel.css"
import "new-dashboard/src/primevue-theme/toolbar.css"
import * as echarts from "echarts"
import { chartTheme } from "../theme/chartTheme"

async function initApp() {
  const app = createApp(App)
  const router = createAndConfigureRouter()
  const pinia = createPinia()
  app.use(router)
  app.use(PrimeVue, {
    theme: {
      //eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
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

// eslint-disable-next-line unicorn/prefer-top-level-await
void initApp()
