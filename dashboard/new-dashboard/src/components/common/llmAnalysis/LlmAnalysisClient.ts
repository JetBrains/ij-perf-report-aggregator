import { ServerConfigurator } from "../dataQuery"
import { SpaceUploadAttachmentsResponse } from "../uploadAttachments/uploadAttachmentsUtils"

export interface LlmAnalysisRequest {
project: string
  metric: string
  currentBuildId: string
  prevBuildId: string
  currentValue?: string
  previousValue?: string
  userName?: string
  firstCommitRevision?: string
  lastCommitRevision?: string
  testMethodName?: string
  ytIssueId?: string
  spaceAttachments?: SpaceUploadAttachmentsResponse
}

export enum LlmAnalysisState {
  InProgress = "in_progress",
  Success = "success",
  Failed = "failed",
}

export interface LlmAnalysisRun {
  id: number
  createdAt: string
  runBuildId: string
  state: LlmAnalysisState
}

export interface LlmAnalysisRunsQuery {
  project: string
  metric: string
  currentBuildId: string
}

export class LlmAnalysisClient {
  private readonly serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async sendLlmAnalysisRequest(request: LlmAnalysisRequest): Promise<LlmAnalysisRun> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/llm/startAnalysis`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(request),
    })

    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to send LLM analysis request: ${response.statusText} ${errorMessage}`)
    }
    return (await response.json()) as LlmAnalysisRun
  }

  async getLlmAnalysisRuns(query: LlmAnalysisRunsQuery): Promise<LlmAnalysisRun[]> {
    const params = new URLSearchParams(query as unknown as Record<string, string>)
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/llm/analysisRuns?${params.toString()}`
    const response = await fetch(url)
    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to fetch LLM analysis runs: ${response.statusText} ${errorMessage}`)
    }
    return (await response.json()) as LlmAnalysisRun[]
  }
}
