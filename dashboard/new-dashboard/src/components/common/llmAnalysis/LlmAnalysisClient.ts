import { ServerConfigurator } from "../dataQuery"
import { SpaceUploadAttachmentsResponse } from "../uploadAttachments/uploadAttachmentsUtils"

export interface LlmAnalysisRequest {
  project: string
  metric: string
  currentBuildId: string
  prevBuildId: string
  spaceAttachments: SpaceUploadAttachmentsResponse
  currentValue?: string
  previousValue?: string
  userName?: string
  firstCommitRevision?: string
  lastCommitRevision?: string
  testMethodName?: string
  ytIssueId?: string
  dashboardLink?: string
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

export interface LlmAnalysisDetails extends LlmAnalysisRun {
  project: string
  metric: string
  currentBuildId: string
  prevBuildId: string
  currentValue?: string
  previousValue?: string
  userName?: string
  userEmail?: string
  firstCommitRevision?: string
  lastCommitRevision?: string
  testMethodName?: string
  ytIssueId?: string
  llmGuiltyCommits?: string[]
  llmComment?: string
  totalCostUsd?: number
}

export interface AnalysisFeedback {
  id: number
  analysisId: number
  rate: number
  feedback?: string
  userEmail?: string
  createdAt: string
  updatedAt: string
}

export class LlmAnalysisClient {
  private readonly serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async sendLlmAnalysisRequest(request: LlmAnalysisRequest): Promise<LlmAnalysisRun> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/llm/analyses`
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

  async getLlmAnalysisRuns(project: string, metric: string, currentBuildId: string): Promise<LlmAnalysisRun[]> {
    const params = new URLSearchParams({ project, metric, currentBuildId })
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/llm/analyses?${params.toString()}`
    const response = await fetch(url)
    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to fetch LLM analysis runs: ${response.statusText} ${errorMessage}`)
    }
    return (await response.json()) as LlmAnalysisRun[]
  }

  async getLlmAnalysisById(id: number | string): Promise<LlmAnalysisDetails> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/llm/analyses/${encodeURIComponent(String(id))}`
    const response = await fetch(url)
    if (response.status === 404) {
      throw new Error("Analysis not found")
    }
    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to fetch LLM analysis: ${response.statusText} ${errorMessage}`)
    }
    return (await response.json()) as LlmAnalysisDetails
  }

  async submitAnalysisFeedback(analysisId: number | string, rate: number, feedback?: string): Promise<void> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/llm/analyses/${encodeURIComponent(String(analysisId))}/feedback`
    const body: { rate: number; feedback?: string } = { rate }
    if (feedback != null && feedback.trim() !== "") {
      body.feedback = feedback
    }
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    })
    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to submit feedback: ${response.statusText} ${errorMessage}`)
    }
  }

  async getAnalysisFeedback(analysisId: number | string): Promise<AnalysisFeedback[]> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/llm/analyses/${encodeURIComponent(String(analysisId))}/feedback`
    const response = await fetch(url)
    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to fetch feedback: ${response.statusText} ${errorMessage}`)
    }
    return (await response.json()) as AnalysisFeedback[]
  }
}
