import { ServerConfigurator } from "../dataQuery"

interface BisectRequest {
  buildId: string
  changes: string
  buildType: string
  testPatterns: string
  requester: string
  mode: string
  excludedCommits: string
  jpsCompilation: string
  dashboardLink?: string
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

interface BuildInfo {
  buildTypeId: string
  number: string
  branchName: string
  startDate: string
}

export interface ChangesGap {
  known: boolean
  hasGap: boolean
  gapCommitCount: number
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

  async fetchChangesGap(buildId: string, previousBuildId: string, currentFirstCommit: string): Promise<ChangesGap | null> {
    try {
      const response = await fetch(
        `${this.serverConfigurator?.serverUrl}/api/meta/teamcity/changesGap?buildId=${encodeURIComponent(buildId)}&previousBuildId=${encodeURIComponent(previousBuildId)}&currentFirstCommit=${encodeURIComponent(currentFirstCommit)}`
      )
      if (!response.ok) {
        const errorText = await response.text()
        console.log(`Failed to fetch changes gap: ${response.status} - ${errorText}`)
        return null
      }
      return (await response.json()) as ChangesGap
    } catch (error) {
      console.log("Error fetching changes gap:", error)
      return null
    }
  }

  async fetchBuildInfo(buildId: string): Promise<BuildInfo | null> {
    try {
      const response = await fetch(`${this.serverConfigurator?.serverUrl}/api/meta/teamcity/buildInfo?buildId=${encodeURIComponent(buildId)}`)
      if (!response.ok) {
        const errorText = await response.text()
        console.log(`Failed to fetch build info: ${response.status} - ${errorText}`)
        return null
      }
      return (await response.json()) as BuildInfo
    } catch (error) {
      console.log("Error fetching TeamCity build info:", error)
      return null
    }
  }
}
