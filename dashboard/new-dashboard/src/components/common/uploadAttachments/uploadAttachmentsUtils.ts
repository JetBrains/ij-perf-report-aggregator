import { ServerConfigurator } from "../dataQuery"

export interface UploadAttachmentsRequest {
  teamcityAttachmentInfo: {
    currentBuildId: number
    previousBuildId: number | undefined
  }
  projectName: string
  chartPng: string | undefined
  testType: string
}

export interface YoutrackUploadAttachmentsRequest extends UploadAttachmentsRequest {
  issueId: string
}

export interface YoutrackUploadAttachmentsResponse {
  uploads: string[]
  exceptions: string[] | undefined
}

export interface SpaceUploadAttachmentsResponse {
  uploads: Record<number, string[]>
  exceptions: Record<number, string[]> | undefined
}

export function uploadAttachmentsToYoutrack(serverConfigurator: ServerConfigurator | null, request: YoutrackUploadAttachmentsRequest): Promise<YoutrackUploadAttachmentsResponse> {
  return postUpload(serverConfigurator, "youtrack", request)
}

export function uploadAttachmentsToSpace(serverConfigurator: ServerConfigurator | null, request: UploadAttachmentsRequest): Promise<SpaceUploadAttachmentsResponse> {
  return postUpload(serverConfigurator, "space", request)
}

async function postUpload<TResponse>(serverConfigurator: ServerConfigurator | null, target: "youtrack" | "space", body: unknown): Promise<TResponse> {
  const response = await fetch(`${serverConfigurator?.serverUrl}/api/meta/${target}/uploadAttachments`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  })

  if (!response.ok) {
    throw new Error(`Failed to upload attachments to ${target}. HTTP error: ${response.status}`)
  }

  return (await response.json()) as TResponse
}
