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
import "./tailwind-primevue-theme.css"

const app = createApp(App)
app.use(createAndConfigureRouter())
app.use(PrimeVue)
const pinia = createPinia()
app.use(pinia)
app.mount("#app")
