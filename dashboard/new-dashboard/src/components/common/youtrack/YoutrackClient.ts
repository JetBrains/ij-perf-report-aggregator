import { ServerConfigurator } from "../dataQuery"

export class YoutrackClient {
  private commonProjects: Project[] = [
    { name: "IntelliJ Platform", id: "22-619" },
    { name: "Automation Testing", id: "22-570" },
  ]
  private serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async createIssue(issueInfo: CreateIssueRequest): Promise<IssueResponse> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/youtrack/createIssue`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(issueInfo),
    })

    const issueResponse = (await response.json()) as IssueResponse

    if (issueResponse.issue.id.length === 0) {
      const error = `Failed to create issue. Errors: ${issueResponse.exceptions?.join("\n") ?? ""}`
      console.error(error)
      throw new Error(error)
    } else {
      return issueResponse
    }
  }

  async uploadAttachments(attachmentsInfo: UploadAttachmentsRequest) {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/youtrack/uploadAttachments`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(attachmentsInfo),
    })

    const attachmentsResponse = (await response.json()) as AttachmentsResponse

    if (attachmentsResponse.exceptions != undefined || !response.ok) {
      throw new Error(`Failed to upload attachments. Errors: ${attachmentsResponse.exceptions?.join("\n") ?? ""}`)
    }
  }

  private static readonly PROJECT_MAP: Record<string, Project[]> = {
    webstorm: [{ name: "WebStorm", id: "22-96" }],
    phpstorm: [{ name: "PhpStorm", id: "22-19" }],
    pycharm: [{ name: "PyCharm", id: "22-36" }],
    clion: [{ name: "CLion", id: "22-139" }],
    goland: [{ name: "GoLand", id: "22-211" }],
    ruby: [{ name: "RubyMine", id: "22-25" }],
    rust: [{ name: "RustRover", id: "22-725" }],
    kotlin: [
      { name: "Kotlin", id: "22-68" },
      { name: "Kotlin Plugin", id: "22-414" },
    ],
    idea: [
      { name: "IDEA", id: "22-22" },
      { name: "Kotlin", id: "22-68" },
      { name: "Kotlin Plugin", id: "22-414" },
    ],
    bazel: [{ name: "Bazel", id: "22-541" }],
    jbr: [{ name: "JetBrains Runtime", id: "22-202" }],
    qodana: [{ name: "Qodana", id: "22-332" }],
    fleet: [{ name: "Fleet", id: "22-520" }],
  }

  private getConfiguratorId(): string {
    if (!this.serverConfigurator) return ""
    return this.serverConfigurator.db === "perfint" || this.serverConfigurator.db === "perfintDev" ? this.serverConfigurator.table : this.serverConfigurator.db
  }

  getProjects(): Project[] {
    const id = this.getConfiguratorId()
    const relatedProjects = YoutrackClient.PROJECT_MAP[id] ?? []
    return [...relatedProjects, ...this.commonProjects]
  }
}

export interface IssueResponse {
  exceptions: string[] | undefined
  issue: Issue
}

export interface AttachmentsResponse {
  exceptions: string[] | undefined
}

interface Issue {
  id: string
  idReadable: string
}

export interface UploadAttachmentsRequest {
  issueId: string
  teamcityAttachmentInfo: {
    buildTypeId: string
    currentBuildId: number
    previousBuildId: number | undefined
  }
  chartPng: string | undefined
  affectedTest: string
}

export interface CreateIssueRequest {
  accidentId: string
  ticketLabel: string
  projectId: string
  buildLink: string
  changesLink: string
  customFields: {
    name: string
    $type: string
    value: { name: string }
  }[]
  dashboardLink: string
  affectedMetric: string
  delta: string
  testMethodName: string | undefined
}

export interface Project {
  id: string
  name: string
}
