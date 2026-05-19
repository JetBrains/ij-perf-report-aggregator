import type { ToastServiceMethods } from "primevue/toastservice"
import { LlmAnalysesConfigurator } from "../../../configurators/llmAnalyses/LlmAnalysesConfigurator"
import { InfoData } from "./InfoSidebar"

export async function startLlmAnalysisWithToast(configurator: LlmAnalysesConfigurator, data: InfoData, toast: ToastServiceMethods): Promise<boolean> {
  try {
    const { buildUrl: url } = await configurator.startRun(data)
    toast.add({
      severity: "success",
      summary: "LLM Analysis Started",
      detail: `View TC build: ${url}`,
      life: 15000,
    })
    return true
  } catch (error) {
    console.error("LLM Analysis start failed:", error)
    toast.add({
      severity: "error",
      summary: "LLM Analysis Failed",
      detail: `Failed to start LLM analysis: ${error instanceof Error ? error.message : String(error)}`,
      life: 8000,
    })
    return false
  }
}
