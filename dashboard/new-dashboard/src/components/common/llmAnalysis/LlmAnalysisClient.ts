import { ServerConfigurator } from "../dataQuery"

export interface LlmAnalysisRequest {
  date: string
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
}

export enum LlmAnalysisState {
  NotStarted = "not_started",
  Queued = "queued",
  InProgress = "in_progress",
  Success = "success",
  Failed = "failed",
}

export interface LlmAnalysisRun {
  id: number
  date: string
  createdAt: string
  runBuildId: string
  state: LlmAnalysisState
}

export interface LlmAnalysisRunsQuery {
  date: string
  project: string
  metric: string
  currentBuildId: string
  prevBuildId: string
}

export class LlmAnalysisClient {
  private readonly serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async sendLlmAnalysisRequest(request: LlmAnalysisRequest): Promise<LlmAnalysisRun> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/teamcity/startLlmAnalysis`
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
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/teamcity/llmAnalysisRuns?${params.toString()}`
    const response = await fetch(url)
    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to fetch LLM analysis runs: ${response.statusText} ${errorMessage}`)
    }
    return (await response.json()) as LlmAnalysisRun[]
  }
}
