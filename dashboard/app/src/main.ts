import "element-plus/theme-chalk/el-reset.css"
import "element-plus/theme-chalk/index.css"
// import "element-plus/theme-chalk/src/base.scss"

import "./main.css"
import { createApp } from "vue"
import App from "./App.vue"
import { createAndConfigureRouter } from "./route"

const app = createApp(App)
app.use(createAndConfigureRouter())
app.mount("#app")
