import "@fontsource/jetbrains-mono/variable.css"
import "./main.css"
import PrimeVue from "primevue/config"
import ToastService from "primevue/toastservice"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"
import "primevue/resources/primevue.css"
import "primeicons/primeicons.css"
// Lara Light is not used because of https://github.com/primefaces/primevue/issues/2390 (arrow doesn't have borders for overlay panel)
import "primevue/resources/themes/saga-blue/theme.css"

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

