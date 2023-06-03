import { ref } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"

export const reportData = ref("")
export const recentlyUsedIdePort = ref(63342)

const stateManager = new PersistentStateManager("ij-report-visualizer")
stateManager.add("data", reportData)
stateManager.add("recentlyUsedIdePort", recentlyUsedIdePort)
