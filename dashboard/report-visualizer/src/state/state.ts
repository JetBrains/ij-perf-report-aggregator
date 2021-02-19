import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { ref } from "vue"

export const markerNames = ["app initialized callback", "module loading"]

export const reportData = ref("")
export const recentlyUsedIdePort = ref(63342)

const stateManager = new PersistentStateManager("ij-report-visualizer")
stateManager.add("data", reportData)
stateManager.add("recentlyUsedIdePort", recentlyUsedIdePort)
stateManager.init()