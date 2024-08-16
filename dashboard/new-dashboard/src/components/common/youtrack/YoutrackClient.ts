import { ServerConfigurator } from "../dataQuery"

export class YoutrackClient {
  private availableProjects: Project[] = [
    { name: "Automation Testing", id: "22-570" },
    { name: "IDEA", id: "22-22" },
    { name: "Kotlin", id: "22-68" },
    { name: "Kotlin Plugin", id: "22-414" },
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

    if (!response.ok) {
      throw new Error(`Failed to upload attachments. Errors: ${attachmentsResponse.exceptions?.join("\n") ?? ""}`)
    }
  }

  getProjects(): Project[] {
    return this.availableProjects
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
