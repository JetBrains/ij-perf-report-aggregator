import { ServerConfigurator } from "../dataQuery"

export enum UploadTarget {
  YOUTRACK = "youtrack",
  SPACE = "space",
}

export interface UploadAttachmentsRequest {
  targets: UploadTarget[]
  issueId: string
  teamcityAttachmentInfo: {
    currentBuildId: number
    previousBuildId: number | undefined
  }
  chartPng: string | undefined
  affectedTest: string
  testType: string
}

export interface AttachmentsResponse {
  uploads: Partial<Record<UploadTarget, string[]>>
  exceptions: string[] | undefined
}

export class UploadAttachmentsClient {
  private readonly serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async uploadAttachments(attachmentsInfo: UploadAttachmentsRequest) {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/uploadAttachments`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(attachmentsInfo),
    })

    if (!response.ok) {
      throw new Error(`Failed to upload attachments. HTTP error: ${response.status}`)
    }

    return (await response.json()) as AttachmentsResponse
  }
}
