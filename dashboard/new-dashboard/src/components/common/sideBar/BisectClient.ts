import { ServerConfigurator } from "../dataQuery"

interface BisectRequest {
  buildId: string
  changes: string
  buildType: string
  className: string
  requester: string
  mode: string
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

interface CommitRevisions {
  firstCommit: string
  lastCommit: string
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

  async fetchBuildType(buildId: string): Promise<string> {
    try {
      const response = await fetch(`${this.serverConfigurator?.serverUrl}/api/meta/teamcity/buildType?buildId=${encodeURIComponent(buildId)}`)
      if (!response.ok) {
        const errorText = await response.text()
        console.log(`Failed to fetch changes: ${response.status} - ${errorText}`)
        return ""
      }
      return (await response.json()) as string
    } catch (error) {
      console.log("Error fetching TeamCity changes:", error)
      return ""
    }
  }

  async fetchTeamCityChanges(buildId: string): Promise<CommitRevisions> {
    try {
      const response = await fetch(`${this.serverConfigurator?.serverUrl}/api/meta/teamcity/changes?buildId=${encodeURIComponent(buildId)}`)
      if (!response.ok) {
        const errorText = await response.text()
        console.log(`Failed to fetch changes: ${response.status} - ${errorText}`)
        return { firstCommit: "", lastCommit: "" }
      }
      return (await response.json()) as CommitRevisions
    } catch (error) {
      console.log("Error fetching TeamCity changes:", error)
      return { firstCommit: "", lastCommit: "" }
    }
  }
}
