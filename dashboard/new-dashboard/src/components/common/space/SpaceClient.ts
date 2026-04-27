import { ServerConfigurator } from "../dataQuery"
import { UploadAttachmentsRequest, UploadAttachmentsResponse } from "../youtrack/YoutrackClient"

export class SpaceClient {
  private readonly serverConfigurator: ServerConfigurator | null

  constructor(serverConfigurator: ServerConfigurator | null) {
    this.serverConfigurator = serverConfigurator
  }

  async uploadAttachments(request: UploadAttachmentsRequest): Promise<UploadAttachmentsResponse> {
    const url = `${this.serverConfigurator?.serverUrl}/api/meta/space/uploadAttachments`
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(request),
    })

    if (!response.ok) {
      throw new Error(`Failed to upload attachments to Space. HTTP error: ${response.status}`)
    }

    return (await response.json()) as UploadAttachmentsResponse
  }
}
