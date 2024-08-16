import { ServerConfigurator } from "../dataQuery"

export class YoutrackClient {
  private availableProjects: Project[] = [
    { name: "Automation Testing", id: "22-570" },
    { name: "IDEA", id: "22-22" },
    { name: "Kotlin", id: "22-68" },
    { name: "Kotlin Plugin", id: "22-414" },
  ]
  private serverConfigurator: ServerConfigurator

  constructor(serverConfigurator: ServerConfigurator) {
    this.serverConfigurator = serverConfigurator
  }

  async createIssue(issueInfo: CreateIssueRequest): Promise<IssueResponse> {
    const url = `${this.serverConfigurator.serverUrl}/api/meta/youtrack/createIssue`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(issueInfo),
    })

    const issueResponse = await response.json()

    if (issueResponse.issue?.id != undefined) {
      return issueResponse
    } else {
      const error = `Failed to create issue. Errors: ${issueResponse.exceptions?.join("\n") ?? issueResponse}`
      console.error(error)
      throw new Error(error)
    }
  }

  async uploadAttachments(attachmentsInfo: UploadAttachmentsRequest) {
    const url = `${this.serverConfigurator.serverUrl}/api/meta/youtrack/uploadAttachments`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(attachmentsInfo),
    })

    const attachmentsResponse = await response.json()

    if (!response.ok) {
      throw new Error(`Failed to upload attachments. Errors: ${attachmentsResponse.exceptions?.join("\n") ?? attachmentsResponse}`)
    }
  }

  getProjects(): Project[] {
    return this.availableProjects
  }
}

export interface IssueResponse {
  exceptions: string[]
  issue: Issue
}

interface Issue {
  id: string
  idReadable: string
}

interface UploadAttachmentsRequest {
  issueId: string
  teamcityAttachmentInfo: {
    buildTypeId: string
    currentBuildId: number
    previousBuildId: number
  }
  chartPng: string | undefined
  affectedTest: string
}

interface CreateIssueRequest {
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
}

export interface Project {
  id: string
  name: string
}
