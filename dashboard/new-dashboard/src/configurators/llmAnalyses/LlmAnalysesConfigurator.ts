import { Ref, ref, watchEffect } from "vue"
import { ServerConfigurator } from "../../components/common/dataQuery"
import { LlmAnalysisClient, LlmAnalysisRequest, LlmAnalysisRun } from "../../components/common/llmAnalysis/LlmAnalysisClient"
import { buildUrl, InfoData } from "../../components/common/sideBar/InfoSidebar"
import { UploadAttachmentsRequest, uploadAttachmentsToSpace } from "../../components/common/uploadAttachments/uploadAttachmentsUtils"
import { dbTypeStore } from "../../shared/dbTypes"
import { useUserStore } from "../../shared/useUserStore"
import { getFirstAndLastCommit } from "../../util/changes"

export interface RunLlmAnalysisResult {
  run: LlmAnalysisRun
  buildUrl: string
}

export class LlmAnalysesConfigurator {
  readonly value: Ref<LlmAnalysisRun[]> = ref([])
  readonly data: Ref<InfoData | null> = ref(null)

  private readonly client: LlmAnalysisClient

  constructor(private readonly serverConfigurator: ServerConfigurator | null) {
    this.client = new LlmAnalysisClient(serverConfigurator)
    watchEffect(() => {
      const d = this.data.value
      if (this.canStart(d)) {
        void this.loadRuns(d as InfoData)
      } else {
        this.value.value = []
      }
    })
  }

  canStart(data: InfoData | null): boolean {
    return this.serverConfigurator != null && data != null && data.buildIdPrevious != null && data.series[0]?.metricName != null
  }

  async loadRuns(data: InfoData): Promise<void> {
    const metric = data.series[0]?.metricName
    if (metric == null || data.buildIdPrevious == null) {
      this.value.value = []
      return
    }
    this.value.value = await this.client.getLlmAnalysisRuns({
      date: data.date,
      project: data.projectName,
      metric,
      currentBuildId: String(data.buildId),
      prevBuildId: String(data.buildIdPrevious),
    })
  }

  async startRun(data: InfoData, ytIssueId?: string): Promise<RunLlmAnalysisResult> {
    if (data.buildIdPrevious == null) {
      throw new Error("Previous build is required to run LLM analysis")
    }
    if (this.serverConfigurator == null) {
      throw new Error("Server configurator is not available")
    }
    const metric = data.series[0]?.metricName
    if (metric == null) {
      throw new Error("Metric is required to run LLM analysis")
    }
    const serverConfigurator = this.serverConfigurator
    const attachmentsInfo: UploadAttachmentsRequest = {
      teamcityAttachmentInfo: {
        currentBuildId: data.buildId,
        previousBuildId: data.buildIdPrevious,
      },
      projectName: data.projectName,
      testType: dbTypeStore().dbType,
    }
    const spaceAttachments = await uploadAttachmentsToSpace(serverConfigurator, attachmentsInfo)
    const { firstCommit, lastCommit } = await getFirstAndLastCommit(serverConfigurator.db, data.installerId ?? data.buildId)
    const request: LlmAnalysisRequest = {
      date: data.date,
      project: data.projectName,
      metric,
      currentBuildId: String(data.buildId),
      prevBuildId: String(data.buildIdPrevious),
      currentValue: data.formattedCurrentValue ?? undefined,
      previousValue: data.formattedPreviousValue ?? undefined,
      userName: useUserStore().user?.name ?? undefined,
      firstCommitRevision: firstCommit ?? undefined,
      lastCommitRevision: lastCommit ?? undefined,
      testMethodName: data.description.value?.methodName?.replaceAll("#", "."),
      ytIssueId: ytIssueId ?? undefined,
      spaceAttachments,
    }
    const run = await this.client.sendLlmAnalysisRequest(request)
    this.value.value = [...this.value.value, run]
    return { run, buildUrl: buildUrl(Number(run.runBuildId)) }
  }
}
