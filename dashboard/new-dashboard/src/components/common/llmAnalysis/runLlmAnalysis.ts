import type { ToastServiceMethods } from "primevue/toastservice"
import { dbTypeStore } from "../../../shared/dbTypes"
import { useUserStore } from "../../../shared/useUserStore"
import { getFirstAndLastCommit } from "../../../util/changes"
import { ServerConfigurator } from "../dataQuery"
import { buildUrl, InfoData } from "../sideBar/InfoSidebar"
import { uploadAttachmentsToSpace, UploadAttachmentsRequest } from "../uploadAttachments/uploadAttachmentsUtils"
import { LlmAnalysisClient, LlmAnalysisRequest, LlmAnalysisRun } from "./LlmAnalysisClient"

export interface RunLlmAnalysisResult {
  run: LlmAnalysisRun
  buildUrl: string
}

export async function runLlmAnalysis(serverConfigurator: ServerConfigurator, data: InfoData, attachmentsInfo: UploadAttachmentsRequest): Promise<RunLlmAnalysisResult> {
  if (data.buildIdPrevious == null) {
    throw new Error("Previous build is required to run LLM analysis")
  }
  await uploadAttachmentsToSpace(serverConfigurator, attachmentsInfo)
  const { firstCommit, lastCommit } = await getFirstAndLastCommit(serverConfigurator.db, data.installerId ?? data.buildId)
  const request: LlmAnalysisRequest = {
    date: data.date,
    project: data.projectName,
    metric: data.series[0]?.metricName ?? "",
    currentBuildId: String(data.buildId),
    prevBuildId: String(data.buildIdPrevious),
    currentValue: data.formattedCurrentValue ?? undefined,
    previousValue: data.formattedPreviousValue ?? undefined,
    userName: useUserStore().user?.name ?? undefined,
    firstCommitRevision: firstCommit ?? undefined,
    lastCommitRevision: lastCommit ?? undefined,
    testMethodName: data.description.value?.methodName?.replaceAll("#", "."),
  }
  const run = await new LlmAnalysisClient(serverConfigurator).sendLlmAnalysisRequest(request)
  return { run, buildUrl: buildUrl(Number(run.runBuildId)) }
}

export async function runLlmAnalysisWithToast(serverConfigurator: ServerConfigurator, data: InfoData, toast: ToastServiceMethods): Promise<boolean> {
  const attachmentsInfo: UploadAttachmentsRequest = {
    teamcityAttachmentInfo: {
      currentBuildId: data.buildId,
      previousBuildId: data.buildIdPrevious,
    },
    projectName: data.projectName,
    testType: dbTypeStore().dbType,
  }
  try {
    const { buildUrl: url } = await runLlmAnalysis(serverConfigurator, data, attachmentsInfo)
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
