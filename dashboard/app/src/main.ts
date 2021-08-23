import "element-plus/packages/theme-chalk/src/reset.scss"
import "element-plus/packages/theme-chalk/src/base.scss"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"

const app = createApp(App)
app.use(createAndConfigureRouter())
app.mount("#app")
