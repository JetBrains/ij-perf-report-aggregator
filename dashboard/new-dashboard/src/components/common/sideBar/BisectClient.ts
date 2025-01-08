import { ServerConfigurator } from "../dataQuery"

interface BisectRequest {
  changes: string
  buildType: string
  className: string
}

interface PerformanceBisectRequest extends BisectRequest {
  targetValue: string
  direction: string
  metric: string
  test: string
}

interface FunctionalBisectRequest extends BisectRequest {
  errorMessage: string
}

export class BisectClient {
  private readonly serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async sendBisectRequest(request: PerformanceBisectRequest | FunctionalBisectRequest): Promise<string> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/teamcity/startBisect`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(request),
    })

    if (!response.ok) {
      const errorMessage = await response.text()
      throw new Error(`Failed to send bisect request: ${response.statusText} ${errorMessage}`)
    }
    return response.text()
  }
}
