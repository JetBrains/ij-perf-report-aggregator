import { ServerConfigurator } from "../dataQuery"

export enum UploadTarget {
  YouTrack = "youtrack",
  Space = "space",
}

export interface UploadAttachmentsRequest {
  issueId: string
  teamcityAttachmentInfo: {
    currentBuildId: number
    previousBuildId: number | undefined
  }
  chartPng: string | undefined
  affectedTest: string
  testType: string
}

export interface UploadAttachmentsResponse {
  uploads: string[]
  exceptions: string[] | undefined
}

export async function uploadAttachments(
  serverConfigurator: ServerConfigurator | null,
  request: UploadAttachmentsRequest,
  target: UploadTarget
): Promise<UploadAttachmentsResponse> {
  const url = `${serverConfigurator?.serverUrl}/api/meta/${target}/uploadAttachments`
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  })

  if (!response.ok) {
    const serviceName = target === UploadTarget.YouTrack ? "YouTrack" : "Space"
    throw new Error(`Failed to upload attachments to ${serviceName}. HTTP error: ${response.status}`)
  }

  return (await response.json()) as UploadAttachmentsResponse
}
