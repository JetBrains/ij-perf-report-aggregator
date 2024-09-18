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
          :disabled="downloadState != DownloadState.NOT_STARTED"
        />
        <label
          class="text-sm"
          for="label"
          >Label</label
        >
      </FloatLabel>
    </div>
    <Dropdown
      v-model="project"
      placeholder="Project"
      :options="projects"
      option-label="name"
      :disabled="downloadState != DownloadState.NOT_STARTED"
    >
      <template #value="{ value }">
        <div class="group inline-flex justify-center text-sm font-medium text-gray-700 hover:text-gray-900">
          {{ value.name }}
          <ChevronDownIcon
            class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
            aria-hidden="true"
          />
        </div>
      </template>
      <template #dropdownicon>
        <!-- empty element to avoid ignoring override of slot -->
        <span />
      </template>
    </Dropdown>
    <!-- Footer buttons -->
    <template #footer>
      <div
        v-if="downloadState == DownloadState.NOT_STARTED"
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
          <div v-if="downloadState == DownloadState.STARTED">
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
            v-else
            class="icon-wrapper"
          >
            <i class="pi pi-times-circle"></i>
            <span class="tooltip-text">Ticket was not created. See console for details</span>
          </div>
        </div>
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { Ref, ref } from "vue"
import { getNavigateToTestUrl, getSpaceUrl, InfoData } from "../sideBar/InfoSidebar"
import { CreateIssueRequest, IssueResponse, Project, UploadAttachmentsRequest } from "./YoutrackClient"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../../configurators/accidents/AccidentsConfigurator"
import { serverConfiguratorKey, youtrackClientKey } from "../../../shared/keys"
import { injectOrError } from "../../../shared/injectionKeys"
import { useRouter } from "vue-router"
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { getPersistentLink } from "../../settings/CopyLink"
import { TimeRangeConfigurator } from "../../../configurators/TimeRangeConfigurator"
import { dbTypeStore } from "../../../shared/dbTypes"

enum DownloadState {
  NOT_STARTED,
  STARTED,
  FINISHED,
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
const showYoutrackDialog = defineModel<boolean>()
const createdTicket = ref("")
const createException = ref(false)
const attachmentException = ref(false)
const downloadState = ref(DownloadState.NOT_STARTED)
const label = ref(accident?.reason ?? "")
const projects: Ref<Project[]> = ref(youtrackClient.getProjects())
const project = ref(projects.value[0])

async function createTicket() {
  try {
    if (data == null) throw new Error("There is no info data")
    if (accident == null) throw new Error("There is no accident")
    if (accidentConfigurator == null) throw new Error("There is no accidentConfigurator")
    downloadState.value = DownloadState.STARTED
    const buildId = data.buildId
    const affectedMetric = data.series[0].metricName ?? ""

    const issueInfo: CreateIssueRequest = {
      accidentId: `${accident.id}`,
      ticketLabel: label.value,
      projectId: project.value.id,
      buildLink: data.artifactsUrl,
      changesLink: (await getSpaceUrl(data, serverConfigurator)) ?? data.changesUrl,
      customFields: [
        {
          name: "Type",
          $type: "SingleEnumIssueCustomField",
          value: {
            name: "Bug",
          },
        },
      ],
      testMethodName: data.description.value?.methodName?.replaceAll("#", "."),
      dashboardLink: `${window.location.origin}${getPersistentLink(getNavigateToTestUrl(data, router), timerangeConfigurator)}`,
      affectedMetric,
      delta: data.deltaPrevious?.replace(/[+-]/g, (match) => (match === "+" ? "-" : "+")) ?? "",
      testType: dbTypeStore().dbType,
    }

    let issueResponse: IssueResponse
    try {
      issueResponse = await youtrackClient.createIssue(issueInfo)
      createdTicket.value = issueResponse.issue.idReadable
      if (issueResponse.exceptions) {
        console.error(`Issue was created, but with some problems:\n ${issueResponse.exceptions.join("\n")}`)
        createException.value = true
      }
    } catch (error: unknown) {
      console.error(error)
      createException.value = true
      return
    }

    try {
      await accidentConfigurator.reloadAccidentData(accident.id)
    } catch (error: unknown) {
      console.error(error)
      createException.value = true
    }

    try {
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
      }
      await youtrackClient.uploadAttachments(attachmentsInfo)
    } catch (error: unknown) {
      console.error(error)
      attachmentException.value = true
      return
    }
  } finally {
    downloadState.value = DownloadState.FINISHED
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
