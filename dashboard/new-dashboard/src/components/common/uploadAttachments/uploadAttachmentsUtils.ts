import { ServerConfigurator } from "../dataQuery"

export interface UploadAttachmentsRequest {
  teamcityAttachmentInfo: {
    currentBuildId: number
    previousBuildId: number | undefined
  }
  projectName: string
  testType: string
}

export interface YoutrackUploadAttachmentsRequest extends UploadAttachmentsRequest {
  issueId: string
  chartPng?: string
}

export interface YoutrackUploadAttachmentsResponse {
  uploads: string[]
  exceptions: string[]
}

export interface SpaceUploadAttachmentsResponse {
  uploads: Record<number, string[]>
  exceptions: Record<number, string[]>
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

export async function fetchChartPngAsBase64(chartDataUrl: string): Promise<string> {
  const response = await fetch(chartDataUrl)
  const blob = await response.blob()
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader()
    reader.addEventListener("loadend", () => {
      if (typeof reader.result === "string") {
        resolve(reader.result.split(",")[1])
      } else {
        reject(new Error("FileReader result is not a string"))
      }
    })
    reader.addEventListener("error", () => {
      reject(new Error("Error reading blob as data URL"))
    })
    reader.readAsDataURL(blob)
  })
}
