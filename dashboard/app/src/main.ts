import "@fontsource/inter/variable.css"
import "@fontsource/jetbrains-mono/variable.css"
import "./main.css"
import PrimeVue from "primevue/config"
import ToastService from "primevue/toastservice"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"
import "primevue/resources/primevue.css"
import "primeicons/primeicons.css"
// we use variable inter font, so, patched version of tailwind-light theme
import "primevue/resources/themes/lara-light-blue/theme.css"

async function initApp() {
  const app = createApp(App)
  const router = createAndConfigureRouter()
  app.use(router)
  app.use(PrimeVue)
  app.use(ToastService)
  await router.isReady()
    .then(() => app.mount("#app"))
}

void initApp()

