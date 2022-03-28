import "element-plus/theme-chalk/el-reset.css"
import "element-plus/theme-chalk/index.css"
// import "element-plus/theme-chalk/src/base.scss"

import "./main.css"
import { createPinia } from "pinia"
import PrimeVue from "primevue/config"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"
import "primevue/resources/themes/tailwind-light/theme.css"
import "primevue/resources/primevue.min.css"
import "primeicons/primeicons.css"

const app = createApp(App)
app.use(createAndConfigureRouter())
app.use(PrimeVue)
const pinia = createPinia()
app.use(pinia)
app.mount("#app")
