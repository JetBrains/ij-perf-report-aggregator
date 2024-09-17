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
    },
  })
  app.use(ToastService)
  app.use(pinia)
  app.directive("tooltip", Tooltip)

  await router.isReady().then(() => app.mount("#app"))
}

// eslint-disable-next-line unicorn/prefer-top-level-await
void initApp()
