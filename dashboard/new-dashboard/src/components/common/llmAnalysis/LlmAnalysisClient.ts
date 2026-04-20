import { ServerConfigurator } from "../dataQuery"

export interface LlmAnalysisRequest {
  currentBuildId: string
  currentValue: string | undefined
  previousValue: string | undefined
  affectedMetric: string
  testMethodName: string | undefined
  youtrackIssueReadableId: string
  youtrackIssueId: string
  spaceUploadedFiles: string[]
}

export class LlmAnalysisClient {
  private readonly serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async sendLlmAnalysisRequest(request: LlmAnalysisRequest): Promise<string> {
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
    return response.text()
  }
}
