import "@fontsource/inter/variable.css"
import "@fontsource/jetbrains-mono/variable.css"
import "./main.css"
import { createPinia } from "pinia"
import PrimeVue from "primevue/config"
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
  const pinia = createPinia()
  app.use(pinia)
  await router.isReady()
    .then(() => app.mount("#app"))
}

void initApp()

