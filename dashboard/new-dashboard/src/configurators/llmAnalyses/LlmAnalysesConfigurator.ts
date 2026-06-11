import { useIntervalFn } from "@vueuse/core"
import { Ref, ref, watch } from "vue"
import { Router } from "vue-router"
import { ServerConfigurator } from "../../components/common/dataQuery"
import { LlmAnalysisClient, LlmAnalysisRequest, LlmAnalysisRun, LlmAnalysisState } from "../../components/common/llmAnalysis/LlmAnalysisClient"
import { buildUrl, getNavigateToTestUrl, InfoData } from "../../components/common/sideBar/InfoSidebar"
import { getPersistentLink } from "../../components/settings/CopyLink"
import { UploadAttachmentsRequest, uploadAttachmentsToSpace } from "../../components/common/uploadAttachments/uploadAttachmentsUtils"
import { dbTypeStore } from "../../shared/dbTypes"
import { useUserStore } from "../../shared/useUserStore"
import { getFirstAndLastCommit } from "../../util/changes"
import { TimeRangeConfigurator } from "../TimeRangeConfigurator"

export interface RunLlmAnalysisResult {
  run: LlmAnalysisRun
  buildUrl: string
}

export class LlmAnalysesConfigurator {
  readonly value: Ref<LlmAnalysisRun[]> = ref([])
  readonly data: Ref<InfoData | null> = ref(null)

  private readonly client: LlmAnalysisClient

  constructor(
    private readonly serverConfigurator: ServerConfigurator | null,
    private readonly router: Router,
    private readonly timerangeConfigurator: TimeRangeConfigurator
  ) {
    this.client = new LlmAnalysisClient(serverConfigurator)

    const { pause, resume } = useIntervalFn(
      () => {
        const d = this.data.value
        if (this.canLoad(d)) {
          void this.loadRuns(d as InfoData)
        }
      },
      30_000,
      { immediate: false }
    )

    watch(this.data, (d) => {
      if (this.canLoad(d)) {
        void this.loadRuns(d as InfoData)
      } else {
        this.value.value = []
        pause()
      }
    })

    watch(
      this.value,
      (runs) => {
        const sidebarOpen = this.data.value != null
        const hasInProgress = runs.some((r) => r.state === LlmAnalysisState.InProgress)
        if (sidebarOpen && hasInProgress) {
          resume()
        } else {
          pause()
        }
      },
      { immediate: true }
    )
  }

  canLoad(data: InfoData | null): boolean {
    return this.serverConfigurator != null && data != null && data.series[0]?.metricName != null
  }

  canStart(data: InfoData | null): boolean {
    return this.canLoad(data) && data?.buildIdPrevious != null
  }

  async loadRuns(data: InfoData): Promise<void> {
    const metric = data.series[0]?.metricName
    if (metric == null) {
      this.value.value = []
      return
    }
    this.value.value = await this.client.getLlmAnalysisRuns(data.projectName, metric, String(data.buildId))
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
      methodName: data.description.value?.methodName?.replaceAll("#", "."),
    }
    const spaceAttachments = await uploadAttachmentsToSpace(serverConfigurator, attachmentsInfo)
    const { firstCommit, lastCommit } = await getFirstAndLastCommit(serverConfigurator.db, data.installerId ?? data.buildId)
    const dashboardLink = getPersistentLink(getNavigateToTestUrl(data, this.router), this.timerangeConfigurator)
    const request: LlmAnalysisRequest = {
      project: data.projectName,
      metric,
      currentBuildId: String(data.buildId),
      prevBuildId: String(data.buildIdPrevious),
      spaceAttachments,
      currentValue: data.formattedCurrentValue ?? undefined,
      previousValue: data.formattedPreviousValue ?? undefined,
      userName: useUserStore().user?.name ?? undefined,
      firstCommitRevision: firstCommit ?? undefined,
      lastCommitRevision: lastCommit ?? undefined,
      testMethodName: data.description.value?.methodName?.replaceAll("#", "."),
      ytIssueId: ytIssueId ?? undefined,
      dashboardLink,
    }
    const run = await this.client.sendLlmAnalysisRequest(request)
    this.value.value = [...this.value.value, run]
    return { run, buildUrl: buildUrl(Number(run.runBuildId)) }
  }
}
