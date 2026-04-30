<template>
  <Dialog
    v-model:visible="showYoutrackDialog"
    modal
    header="Youtrack"
    :style="{ minWidth: '30vw', maxWidth: '50vw' }"
  >
    <div class="flex items-center space-x-4 mb-4 mt-6">
      <FloatLabel class="w-full">
        <InputText
          id="label"
          v-model="label"
          class="w-full"
          :disabled="progressState != ProgressState.NOT_STARTED"
        />
        <label for="label">Label</label>
      </FloatLabel>
    </div>
    <Select
      v-model="project"
      placeholder="Project"
      :options="projects"
      option-label="name"
      :disabled="progressState != ProgressState.NOT_STARTED"
    >
      <template #value="{ value }">
        <div class="group inline-flex justify-center font-medium">
          {{ value.name }}
          <ChevronDownIcon
            class="-mr-1 ml-1 h-5 w-5 shrink-0"
            aria-hidden="true"
          />
        </div>
      </template>
      <template #dropdownicon>
        <!-- empty element to avoid ignoring override of slot -->
        <span />
      </template>
    </Select>
    <!-- Footer buttons -->
    <template #footer>
      <div
        v-if="progressState == ProgressState.NOT_STARTED"
        class="flex justify-end space-x-2"
      >
        <Button
          label="Cancel"
          icon="pi pi-times"
          severity="secondary"
          @click="showYoutrackDialog = false"
        />
        <Button
          label="Create"
          icon="pi pi-check"
          autofocus
          @click="createTicket"
        />
      </div>
      <div
        v-else
        class="flex flex-1 justify-center align-middle flex-col"
      >
        <div class="flex justify-between items-center">
          <div>
            Creating issue:
            <a
              v-if="createdTicket.length > 0"
              target="_blank"
              class="link-like-text"
              :href="`https://youtrack.jetbrains.com/issue/${createdTicket}`"
            >
              {{ createdTicket }}
            </a>
          </div>
          <div v-if="createdTicket.length <= 0 && !createException">
            <i class="pi pi-spin pi-spinner"></i>
          </div>
          <div
            v-else-if="createdTicket.length > 0"
            class="icon-wrapper"
          >
            <i :class="{ pi: true, 'pi-verified': true, exception: createException }"></i>
            <span
              v-if="createException"
              class="tooltip-text"
              >Ticket was created but with problems. See console for details.</span
            >
          </div>
          <div
            v-else
            class="icon-wrapper"
          >
            <i class="pi pi-times-circle"></i>
            <span class="tooltip-text">Ticket was not created. See console for details</span>
          </div>
        </div>
        <div class="flex justify-between items-center mt-10">
          <div>Uploading attachments</div>
          <div v-if="progressState == ProgressState.UPLOADING_ATTACHMENTS">
            <i class="pi pi-spin pi-spinner"></i>
          </div>
          <div
            v-else-if="createdTicket.length > 0"
            class="icon-wrapper"
          >
            <i :class="{ pi: true, 'pi-verified': true, exception: attachmentException }"></i>
            <span
              v-if="attachmentException"
              class="tooltip-text"
              >Attachments were not uploaded (fully or partially). See console for details.</span
            >
          </div>
          <div
            v-else-if="progressState == ProgressState.FINISHED"
            class="icon-wrapper"
          >
            <i class="pi pi-times-circle"></i>
            <span class="tooltip-text">Ticket was not created. See console for details</span>
          </div>
        </div>
        <div
          v-if="llmAnalysisDisplay.active"
          class="flex justify-between items-center mt-10"
        >
          <div>
            LLM Analysis:
            <template v-if="llmAnalysisDisplay.preparing"> preparation... </template>
            <a
              v-else-if="llmAnalysisDisplay.done && llmAnalysisBuildUrl.length > 0"
              target="_blank"
              class="link-like-text"
              :href="llmAnalysisBuildUrl"
            >
              View TC Build
            </a>
            <template v-else-if="llmAnalysisDisplay.failed"> failed </template>
          </div>
          <div class="icon-wrapper">
            <i :class="llmAnalysisDisplay.icon"></i>
            <span
              v-if="llmAnalysisDisplay.done"
              class="tooltip-text"
              >Experimental feature, the results will be published to YT ticket</span
            >
          </div>
        </div>
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { computed, Ref, ref } from "vue"
import { useToast } from "primevue/usetoast"
import { getNavigateToTestUrl, getSpaceUrl, InfoData } from "../sideBar/InfoSidebar"
import { generateDefaultReason } from "../sideBar/AccidentUtils"
import { CreateIssueRequest, IssueResponse, Project } from "./YoutrackClient"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../../configurators/accidents/AccidentsConfigurator"
import { serverConfiguratorKey, youtrackClientKey } from "../../../shared/keys"
import { injectOrError } from "../../../shared/injectionKeys"
import { useRouter } from "vue-router"
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { getPersistentLink } from "../../settings/CopyLink"
import { TimeRangeConfigurator } from "../../../configurators/TimeRangeConfigurator"
import { dbTypeStore } from "../../../shared/dbTypes"
import { LlmAnalysisClient, LlmAnalysisRequest } from "../llmAnalysis/LlmAnalysisClient"
import { uploadAttachments, UploadAttachmentsRequest, UploadTarget } from "../uploadAttachments/uploadAttachmentsUtils"

enum ProgressState {
  NOT_STARTED,
  CREATING_ISSUE,
  UPLOADING_ATTACHMENTS,
  FINISHED,
}

enum LlmAnalysisState {
  NOT_STARTED,
  PREPARING,
  DONE,
  FAILED,
}

const router = useRouter()

const { data, accident, accidentConfigurator, timerangeConfigurator } = defineProps<{
  data: InfoData | null
  accident: Accident | null
  accidentConfigurator: AccidentsConfigurator | null
  timerangeConfigurator: TimeRangeConfigurator
}>()

const youtrackClient = injectOrError(youtrackClientKey)
const serverConfigurator = injectOrError(serverConfiguratorKey)
const llmAnalysisClient = new LlmAnalysisClient(serverConfigurator)
const toast = useToast()
const showYoutrackDialog = defineModel<boolean>()
const createdTicket = ref("")
const createException = ref(false)
const attachmentException = ref(false)
const llmAnalysisBuildUrl = ref("")
const llmAnalysisState = ref(LlmAnalysisState.NOT_STARTED)
const progressState = ref(ProgressState.NOT_STARTED)
const label = ref(generateLabel())

const llmAnalysisDisplay = computed(() => ({
  active: llmAnalysisState.value !== LlmAnalysisState.NOT_STARTED,
  preparing: llmAnalysisState.value === LlmAnalysisState.PREPARING,
  done: llmAnalysisState.value === LlmAnalysisState.DONE,
  failed: llmAnalysisState.value === LlmAnalysisState.FAILED,
  icon:
    llmAnalysisState.value === LlmAnalysisState.PREPARING ? "pi pi-spin pi-spinner" : llmAnalysisState.value === LlmAnalysisState.DONE ? "pi pi-verified" : "pi pi-times-circle",
}))

function generateLabel(): string {
  if (data == null || accident == null) return ""
  const defaultReason = generateDefaultReason(data)
  if (accident.reason !== defaultReason) return accident.reason
  return `${accident.kind} ${defaultReason} ${data.mode ? `on ${data.mode} mode` : ""}`
}

function reportAttachmentFailure(message: string, error?: unknown) {
  let detail = message
  if (error !== undefined) {
    console.error(message, error)
    detail = `${message}: ${error instanceof Error ? error.message : String(error)}`
  }
  toast.add({
    severity: "error",
    summary: "Attachment Upload Failed",
    detail,
    life: 8000,
  })
  attachmentException.value = true
  progressState.value = ProgressState.FINISHED
}
const projects: Ref<Project[]> = ref(youtrackClient.getProjects())
const project = ref(projects.value[0])

async function createTicket() {
  if (label.value.trim().length < 5) {
    toast.add({
      severity: "error",
      summary: "Validation Error",
      detail: "Label must be at least 5 characters long",
      life: 5000,
    })
    return
  }
  if (data == null) throw new Error("There is no info data")
  if (accident == null) throw new Error("There is no accident")
  if (accidentConfigurator == null) throw new Error("There is no accidentConfigurator")

  progressState.value = ProgressState.CREATING_ISSUE

  const buildId = data.buildId
  const affectedMetric = data.series[0].metricName ?? ""
  const spaceUrls = await getSpaceUrl(data, serverConfigurator)

  const issueInfo: CreateIssueRequest = {
    accidentId: `${accident.id}`,
    ticketLabel: label.value,
    projectId: project.value.id,
    buildLink: data.artifactsUrl,
    changesLink: spaceUrls.length > 0 ? spaceUrls.join(",") : data.changesUrl,
    testMethodName: data.description.value?.methodName?.replaceAll("#", "."),
    dashboardLink: `${window.location.origin}${getPersistentLink(getNavigateToTestUrl(data, router), timerangeConfigurator)}`,
    affectedMetric,
    delta: data.deltaPrevious?.replaceAll(/[+-]/g, (match) => (match === "+" ? "-" : "+")) ?? "",
    currentValue: data.formattedCurrentValue ?? "",
    previousValue: data.formattedPreviousValue ?? "",
    testType: dbTypeStore().dbType,
  }

  let issueResponse: IssueResponse
  try {
    issueResponse = await youtrackClient.createIssue(issueInfo)
    createdTicket.value = issueResponse.issue.idReadable
    if (issueResponse.exceptions) {
      console.error(`Issue was created, but with some problems:\n ${issueResponse.exceptions.join("\n")}`)
      toast.add({
        severity: "warn",
        summary: "Issue Created with Problems",
        detail: `YouTrack issue was created but with some problems:\n${issueResponse.exceptions.join("\n")}`,
        life: 8000,
      })
      createException.value = true
    }
  } catch (error: unknown) {
    console.error(error)
    const errorMessage = error instanceof Error ? error.message : String(error)
    toast.add({
      severity: "error",
      summary: "Issue Creation Failed",
      detail: `Failed to create YouTrack issue: ${errorMessage}`,
      life: 8000,
    })
    createException.value = true
    return
  }

  try {
    await accidentConfigurator.reloadAccidentData(accident.id)
  } catch (error: unknown) {
    console.error(error)
    const errorMessage = error instanceof Error ? error.message : String(error)
    toast.add({
      severity: "error",
      summary: "Accident Data Reload Failed",
      detail: `Failed to reload accident data: ${errorMessage}`,
      life: 8000,
    })
    createException.value = true
  }

  let affectedTest = accident.affectedTest

  if (affectedTest.endsWith(affectedMetric)) {
    affectedTest = affectedTest.slice(0, -affectedMetric.length - 1)
  }
  const attachmentsInfo: UploadAttachmentsRequest = {
    issueId: issueResponse.issue.id,
    teamcityAttachmentInfo: {
      currentBuildId: buildId,
      previousBuildId: undefined,
    },
    affectedTest,
    chartPng: undefined,
    testType: dbTypeStore().dbType,
  }
  if (accident.kind != AccidentKind.Exception) {
    attachmentsInfo.teamcityAttachmentInfo.previousBuildId = data.buildIdPrevious
    try {
      attachmentsInfo.chartPng = await fetch(data.chartDataUrl)
        .then((res) => res.blob())
        .then((blob) => {
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
        })
    } catch (error: unknown) {
      reportAttachmentFailure("Failed to prepare chart for upload", error)
      return
    }
  }

  progressState.value = ProgressState.UPLOADING_ATTACHMENTS
  uploadAttachments(serverConfigurator, attachmentsInfo, UploadTarget.YouTrack)
    .then((response) => {
      if (response.exceptions?.length) {
        reportAttachmentFailure(`Failed to upload attachments. Errors: ${response.exceptions.join("\n")}`)
      } else {
        progressState.value = ProgressState.FINISHED
      }
      return response
    })
    .catch((error: unknown) => {
      reportAttachmentFailure("Failed to upload attachments to YouTrack", error)
    })

  if (accident.kind === AccidentKind.Regression || accident.kind === AccidentKind.Improvement) {
    llmAnalysisState.value = LlmAnalysisState.PREPARING

    uploadAttachments(serverConfigurator, attachmentsInfo, UploadTarget.Space)
      .then(async (response) => {
        try {
          const llmAnalysisRequest: LlmAnalysisRequest = {
            currentBuildId: `${data.buildId}`,
            currentValue: data.formattedCurrentValue || undefined,
            previousValue: data.formattedPreviousValue || undefined,
            affectedMetric,
            testMethodName: data.description.value?.methodName?.replaceAll("#", "."),
            youtrackIssueReadableId: issueResponse.issue.idReadable,
            youtrackIssueId: issueResponse.issue.id,
            spaceUploadedFiles: response.uploads ?? [],
          }
          llmAnalysisBuildUrl.value = await llmAnalysisClient.sendLlmAnalysisRequest(llmAnalysisRequest)
          llmAnalysisState.value = LlmAnalysisState.DONE
        } catch (error) {
          console.error("LLM Analysis start failed:", error)
          llmAnalysisState.value = LlmAnalysisState.FAILED
        }
      })
      .catch((error: unknown) => {
        console.error("Space attachment upload for LLM analysis failed:", error)
        llmAnalysisState.value = LlmAnalysisState.FAILED
      })
  }
}
</script>

<style scoped>
.pi {
  font-size: 1.5rem;
}

.pi-spinner {
  color: dodgerblue;
}

.pi-times-circle {
  color: red;
}

.pi-verified {
  color: green;
}

.pi-verified.exception {
  color: #e3d716;
}

.icon-wrapper {
  position: relative;
  display: inline-block;
}

.tooltip-text {
  visibility: hidden;
  width: 120px;
  background-color: black;
  color: #fff;
  text-align: center;
  border-radius: 6px;
  padding: 5px 0;

  /* Position the tooltip */
  position: absolute;
  z-index: 1;
  bottom: 125%; /* Position above the icon */
  left: 50%;
  margin-left: -60px;

  opacity: 0;
  transition: opacity 0.3s;
}

.icon-wrapper:hover .tooltip-text {
  visibility: visible;
  opacity: 1;
}

.link-like-text {
  color: dodgerblue;
  text-decoration: underline;
  cursor: pointer;
}

.link-like-text:hover {
  color: darkblue;
}
</style>
