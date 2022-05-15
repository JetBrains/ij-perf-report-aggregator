import "@fontsource/jetbrains-mono/variable.css"
import "@fontsource/inter/variable.css"
import "./main.css"
import PrimeVue from "primevue/config"
import ToastService from "primevue/toastservice"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"
// get rid of color.png
import "./primevue.css"
// due to https://github.com/primefaces/primeicons/issues/301, we don't use primeicons package, but instead svg source was converted into woff2
import "./primeicons.css"
// import "../../../jb/prime-theme/themes/saga/saga-blue/theme.scss"
import "../../../jb/prime-theme/themes/lara/lara-light/blue/theme.scss"

import "shared/src/primevue-theme/select.css"
import "shared/src/primevue-theme/select-panel.css"
import "shared/src/primevue-theme/misc.css"

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

