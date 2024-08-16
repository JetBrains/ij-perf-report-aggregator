<template>
  <Dialog
    v-model:visible="showYoutrackDialog"
    modal
    header="Youtrack"
    :style="{ minWidth: '30vw', maxWidth: '50vw' }"
  >
    <div class="flex items-center space-x-4 mb-4 mt-4">
      <FloatLabel class="w-full">
        <InputText
          id="label"
          v-model="reason"
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
      optionLabel="name"
      optionValue="id"
      :disabled="downloadState != DownloadState.NOT_STARTED"
    >
    </Dropdown>
    <!-- Footer buttons -->
    <template #footer>
      <div
        class="flex justify-end space-x-2"
        v-if="downloadState == DownloadState.NOT_STARTED"
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
import { IssueResponse, Project } from "./YoutrackClient"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../../configurators/AccidentsConfigurator"
import { serverConfiguratorKey, youtrackClientKey } from "../../../shared/keys"
import { injectOrError } from "../../../shared/injectionKeys"
import { useRouter } from "vue-router"
import { getTeamcityBuildType } from "../../../util/artifacts"

enum DownloadState {
  NOT_STARTED,
  STARTED,
  FINISHED,
}

const router = useRouter()

const props = defineProps<{
  data: InfoData
  accident: Accident
  accidentConfigurator: AccidentsConfigurator
}>()

const youtrackClient = injectOrError(youtrackClientKey)
const serverConfigurator = injectOrError(serverConfiguratorKey)
const showYoutrackDialog = defineModel<boolean>()
const createdTicket = ref("")
const createException = ref(false)
const attachmentException = ref(false)
const downloadState = ref(DownloadState.NOT_STARTED)
const reason = ref(props.accident?.reason ?? "")
const project = ref("")
const projects: Ref<Project[]> = ref(youtrackClient.getProjects())

async function createTicket() {
  try {
    downloadState.value = DownloadState.STARTED
    const buildId = props.data.buildId!!
    const prevBuildId = props.data.buildIdPrevious!!
    const buildType = await getTeamcityBuildType(serverConfigurator.db, serverConfigurator.table, buildId)
    const chartPng =
      props.accident.kind != AccidentKind.Exception
        ? ((await fetch(props.data.chartDataUrl)
            .then((res) => res.blob())
            .then(
              (blob) =>
                new Promise((resolve, reject) => {
                  const reader = new FileReader()
                  reader.onloadend = () => resolve((reader.result as string).split(",")[1]) // Extract Base64 data
                  reader.onerror = () => reject(new Error("Error reading blob as data URL"))
                  reader.readAsDataURL(blob)
                })
            )) as string)
        : undefined

    const issueInfo = {
      accidentId: `${props.accident.id}`,
      projectId: project.value,
      buildLink: props.data.artifactsUrl,
      changesLink: (await getSpaceUrl(props.data, serverConfigurator)) ?? props.data.changesUrl,
      customFields: [
        {
          name: "Type",
          $type: "SingleEnumIssueCustomField",
          value: {
            name: "Bug",
          },
        },
      ],
      dashboardLink: `${window.location.origin}${getNavigateToTestUrl(props.data, router)}`,
      affectedMetric: props.data.series[0].metricName ?? "",
      delta: props.data.deltaPrevious ?? "",
      testMethodName: props.data.description.value?.methodName?.replaceAll("#", ".") ?? "",
    }

    let issueResponse: IssueResponse
    try {
      issueResponse = await youtrackClient.createIssue(issueInfo)
      createdTicket.value = issueResponse.issue.idReadable
      if (issueResponse.exceptions != null) {
        console.log(`Issue was created, but with some problems:\n ${issueResponse.exceptions.join("\n")}`)
        createException.value = true
      }
    } catch (t: any) {
      console.log(t)
      createException.value = true
      return
    }

    try {
      await props.accidentConfigurator.reloadAccidentData(props.accident.id)
    } catch (t: any) {
      console.log(t)
      createException.value = true
    }

    try {
      const attachmentsInfo = {
        issueId: issueResponse.issue.id,
        teamcityAttachmentInfo: {
          buildTypeId: buildType!!,
          currentBuildId: buildId,
          previousBuildId: prevBuildId,
        },
        affectedTest: props.accident.affectedTest,
        chartPng,
      }
      await youtrackClient.uploadAttachments(attachmentsInfo)
    } catch (t: any) {
      console.log(t)
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
