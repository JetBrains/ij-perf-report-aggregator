import "@fontsource/jetbrains-mono"
import "@fontsource/inter"
import "floating-vue/dist/style.css"
import "./main.css"
import FloatingVue from "floating-vue"
import { createPinia } from "pinia"
import PrimeVue from "primevue/config"
import ToastService from "primevue/toastservice"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"
// get rid of color.png
// avoid tiff/svg/other deprecated stuff in a final build
import "./primeicons.css"
import "../prime-theme/themes/aura/aura-light/blue/theme.scss"

import "new-dashboard/src/primevue-theme/select.css"
import "new-dashboard/src/primevue-theme/select-panel.css"
import "new-dashboard/src/primevue-theme/misc.css"

async function initApp() {
  const app = createApp(App)
  const router = createAndConfigureRouter()
  const pinia = createPinia()
  app.use(router)
  app.use(PrimeVue)
  app.use(ToastService)
  app.use(pinia)
  app.use(FloatingVue, {
    themes: {
      info: {
        $extend: "tooltip",
        placement: "top-start",
      },
    },
  })

  await router.isReady().then(() => app.mount("#app"))
}

// eslint-disable-next-line unicorn/prefer-top-level-await
void initApp()
