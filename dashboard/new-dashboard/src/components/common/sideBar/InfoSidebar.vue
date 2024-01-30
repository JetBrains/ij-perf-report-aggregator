<template>
  <div
    v-show="vm.visible.value"
    class="infoSidebar ml-5 text-gray-500 relative"
  >
    <div class="flex flex-col gap-4 sticky top-2 border border-solid rounded-md border-zinc-200 p-5">
      <div
        v-if="useScrollStore().isScrolled"
        class="sticky h-10"
      ></div>
      <span class="flex justify-between uppercase text-xs">
        {{ data?.title }}

        <span
          class="infoSidebar_icon text-sm pi pi-plus rotate-45 cursor-pointer hover:text-gray-800 transition duration-150 ease-out relative"
          @click="handleCloseClick"
        />
      </span>

      <VTooltip
        v-if="data?.description.value?.description"
        theme="info"
      >
        <div class="flex gap-1.5 font-medium text-base items-center break-all">
          <span
            v-if="data?.series.length == 1"
            class="w-3 h-3 rounded-full"
            :style="{ backgroundColor: data?.series[0].color }"
          />
          <span class="underline decoration-dotted hover:no-underline">{{ data?.projectName }}</span>
        </div>
        <template #popper>
          <span class="text-sm">
            {{ data?.description.value?.description }}
          </span>
        </template>
      </VTooltip>
      <div
        v-else
        class="flex gap-1.5 font-medium text-base items-center break-all"
      >
        <span
          v-if="data?.series.length == 1"
          class="w-3 h-3 rounded-full"
          :style="{ backgroundColor: data?.series[0].color }"
        />
        {{ data?.projectName }}
      </div>

      <SplitButton
        v-if="testActions.length > 0"
        :label="testActions[0].label"
        :model="testActions.slice(1)"
        link
        icon="pi pi-chart-line"
        @click="testActions[0].command"
      />

      <div class="flex flex-col gap-2">
        <span class="flex gap-1.5 text-sm items-center">
          <CalendarIcon class="w-4 h-4" />
          {{ data?.date }}
          <span v-if="data?.build">build {{ data?.build }}</span>
        </span>
        <span
          v-if="data?.branch"
          class="flex gap-1.5 text-sm items-center"
        >
          <BranchIcon class="w-4 h-4" />
          <span>{{ data?.branch }}</span>
        </span>
        <div
          v-if="data?.series.length == 1"
          class="flex flex-col gap-2"
        >
          <span
            v-if="data?.series[0].nameToShow"
            class="flex gap-1.5 text-sm items-center"
          >
            <BeakerIcon class="w-4 h-4" />
            <span>{{ data?.series[0].nameToShow }}</span>
          </span>
          <span class="flex gap-1.5 text-sm items-center">
            <ClockIcon class="w-4 h-4" />
            {{ data?.series[0].value }}
          </span>
        </div>
        <div v-else>
          <div class="grid grid-cols-[repeat(3,_max-content)] whitespace-nowrap gap-x-2 items-baseline leading-loose text-sm">
            <template
              v-for="item in data?.series"
              :key="item.metricName"
            >
              <span
                v-if="item.nameToShow"
                class="rounded-lg w-2.5 h-2.5"
                :style="{ 'background-color': item.color }"
              />
              <span v-if="item.nameToShow">{{ item.nameToShow }}</span>
              <span
                v-if="item.nameToShow"
                class="font-mono place-self-end"
                >{{ item.value }}</span
              >
            </template>
          </div>
        </div>

        <span
          v-if="data?.deltaPrevious"
          class="flex gap-1.5 text-sm items-center"
        >
          <ArrowLeftIcon class="w-4 h-4" />
          {{ data?.deltaPrevious }}
        </span>
        <span
          v-if="data?.deltaNext"
          class="flex gap-1.5 text-sm items-center"
        >
          <ArrowRightIcon class="w-4 h-4" />
          {{ data?.deltaNext }}
        </span>

        <span class="flex gap-1.5 text-sm items-center">
          <ComputerDesktopIcon class="w-4 h-4" />
          {{ data?.machineName }}
        </span>

        <span
          v-if="data?.accidents"
          class="flex gap-1.5 text-sm items-center"
        >
          <ExclamationTriangleIcon class="w-4 h-4" />
          Known events:
        </span>
        <ul
          v-if="data?.accidents"
          class="gap-1.5 text-sm ml-5 overflow-y-auto max-h-80"
        >
          <li
            v-for="accident in data?.accidents.value"
            :key="accident?.id"
          >
            <span class="flex items-start justify-between gap-1.5 text-sm">
              &bull;
              <span class="w-full">{{ accident.kind }}: {{ accident.reason }}</span>
              <TrashIcon
                class="w-4 h-4 text-red-500 flex-none"
                @click="handleRemove(accident.id)"
              />
            </span>
          </li>
        </ul>
      </div>

      <div class="flex gap-4 text-primary">
        <a
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-darker cursor-pointer"
          @click="getChangesUrl"
        >
          <ArrowPathIcon class="w-4 h-4" />
          Changes
        </a>
        <a
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-darker cursor-pointer"
          @click="getArtifactsUrl"
        >
          <ServerStackIcon class="w-4 h-4" />
          Test Artifacts
        </a>
        <a
          v-if="data?.installerUrl !== undefined"
          :href="data?.installerUrl"
          target="_blank"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-darker"
        >
          <ArrowDownTrayIcon class="w-4 h-4" />
          Installer
        </a>
      </div>
      <div class="flex gap-4 text-primary justify-center">
        <a
          v-if="data?.installerId || vm.data.value?.buildId"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-darker cursor-pointer"
          @click="getSpaceUrl"
        >
          <SpaceIcon class="w-4 h-4" />
          Changes
        </a>
      </div>

      <Button
        v-if="accidentsConfigurator != null"
        class="text-sm"
        label="Report Event"
        text
        size="small"
        @click="showDialog = true"
      />
    </div>
  </div>
  <Dialog
    v-model:visible="showDialog"
    modal
    header="Report Event"
    :style="{ width: '30vw' }"
  >
    <div class="flex items-center space-x-4 mb-4">
      <Dropdown
        v-model="accidentType"
        placeholder="Event Type"
        :options="getAccidentTypes()"
      >
        <template #value="{ value }">
          <div class="group inline-flex justify-center text-sm font-medium text-gray-700 hover:text-gray-900">
            {{ value }}
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
      <span class="p-float-label flex-grow">
        <InputText
          id="reason"
          v-model="reason"
          class="w-full"
        />
        <label
          class="text-sm"
          for="reason"
          >Reason</label
        >
      </span>
    </div>
    <div
      v-if="vm.data.value?.series.length == 1"
      class="flex items-center mb-4"
    >
      <InputSwitch
        v-model="reportMetricOnly"
        input-id="reportMetricOnly"
      />
      <label
        for="reportMetricOnly"
        class="text-sm ml-2"
      >
        Report metric only
      </label>
    </div>
    <!-- Footer buttons -->
    <template #footer>
      <div class="flex justify-end space-x-2">
        <Button
          label="Cancel"
          icon="pi pi-times"
          text
          @click="showDialog = false"
        />
        <Button
          label="Report"
          icon="pi pi-check"
          autofocus
          @click="reportRegression"
        />
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { useStorage } from "@vueuse/core/index"
import { computed, ref } from "vue"
import { useRouter } from "vue-router"
import { AccidentKind } from "../../../configurators/AccidentsConfigurator"
import { injectOrError, injectOrNull } from "../../../shared/injectionKeys"
import { accidentsConfiguratorKey, serverConfiguratorKey, sidebarVmKey } from "../../../shared/keys"
import { getTeamcityBuildType } from "../../../util/artifacts"
import { calculateChanges } from "../../../util/changes"
import BranchIcon from "../BranchIcon.vue"
import SpaceIcon from "../SpaceIcon.vue"
import { useScrollListeners, useScrollStore } from "../scrollStore"
import { DBType, tcUrl } from "./InfoSidebar"

const vm = injectOrError(sidebarVmKey)
const showDialog = ref(false)
const reason = ref("")
const router = useRouter()
const accidentType = ref<string>("Regression")
const serverConfigurator = injectOrNull(serverConfiguratorKey)

const accidentsConfigurator = injectOrNull(accidentsConfiguratorKey)

const data = computed(() => vm.data.value)

function reportRegression() {
  showDialog.value = false
  const value = vm.data.value
  if (value == null) {
    console.log("value is zero! This shouldn't happen")
  } else {
    const reportOnlyMetric = reportMetricOnly.value && value.series.length == 1
    accidentsConfigurator?.writeAccidentToMetaDb(
      value.date,
      value.projectName + (reportOnlyMetric ? "/" + value.series[0].metricName : ""),
      reason.value,
      value.build ?? value.buildId.toString(),
      accidentType.value
    )
  }
}

const reportMetricOnly = useStorage("reportMetricOnly", false)

function copyMethodNameToClipboard(methodName: string) {
  void navigator.clipboard.writeText(methodName)
}

function openTestInIDE(methodName: string) {
  const origin = encodeURIComponent("ssh://git@git.jetbrains.team/ij/intellij.git")
  window.open(`jetbrains://idea/navigate/reference?origin=${origin}&fqn=${methodName}`)
}

function handleNavigateToTest() {
  const dbType = vm.data.value?.dbType
  const currentRoute = router.currentRoute.value
  let parts = currentRoute.path.split("/")
  if (parts.at(-1) == "startup" || parts.at(1) == "ij") {
    parts = ["", "ij", "explore"]
  } else {
    parts[parts.length - 1] = dbType == DBType.INTELLIJ_DEV ? "testsDev" : "tests"
  }
  const branch = vm.data.value?.branch ?? ""
  const majorBranch = branch.includes(".") ? branch.slice(0, branch.indexOf(".")) : branch
  const testURL = parts.join("/")

  const queryParams: string = new URLSearchParams({
    ...currentRoute.query,
    project: vm.data.value?.projectName ?? "",
    branch: majorBranch,
  }).toString()

  const measures =
    vm.data.value?.series
      .map((s) => s.metricName)
      .map((m) => "&measure=" + m)
      .join("") ?? ""
  void router.push(testURL + "?" + queryParams + measures)
}

function isNavigateToTestSupported(): boolean {
  const currentRoute = router.currentRoute.value
  const parts = currentRoute.path.split("/")
  const pageName = parts.at(-1)?.toLowerCase()
  return pageName != "testsDev" && pageName != "tests" && pageName != "explore"
}

function handleRemove(id: number) {
  accidentsConfigurator?.removeAccidentFromMetaDb(id)
}

function handleCloseClick() {
  vm.close()
}

const testActions = computed(() => getTestActions())

function getTestActions(): {
  label: string
  icon: string
  command: () => void
}[] {
  const actions = []
  if (isNavigateToTestSupported()) {
    actions.push({
      label: "Navigate to test",
      icon: "pi pi-chart-line",
      command() {
        handleNavigateToTest()
      },
    })
  }
  if (vm.data.value?.description) {
    const methodName = vm.data.value.description.value?.methodName
    if (methodName && methodName != "") {
      actions.push(
        {
          label: "Open test method",
          icon: "pi pi-folder-open",
          command() {
            openTestInIDE(methodName)
          },
        },
        {
          label: "Copy test method name",
          icon: "pi pi-copy",
          command() {
            copyMethodNameToClipboard(methodName)
          },
        }
      )
    }

    const url = vm.data.value.description.value?.url
    if (url && url != "") {
      actions.push({
        label: "Download test project",
        icon: "pi pi-download",
        command() {
          window.open(url)
        },
      })
    }
  }
  return actions
}

function getChangesUrl() {
  if (serverConfigurator?.table == null) {
    window.open(vm.data.value?.changesUrl)
  } else if (vm.data.value?.installerId ?? vm.data.value?.buildId) {
    const db = serverConfigurator.db
    if (db == "perfint" || db == "perfintDev") {
      getTeamcityBuildType(db, serverConfigurator.table, vm.data.value.buildId, (type: string | null) => {
        if (vm.data.value) {
          window.open(`${tcUrl}buildConfiguration/${type}/${vm.data.value.buildId}?buildTab=changes`)
        }
      })
    } else {
      window.open(vm.data.value.changesUrl)
    }
  }
}

function getArtifactsUrl() {
  if (serverConfigurator?.table == null) {
    window.open(vm.data.value?.artifactsUrl)
  } else if (vm.data.value?.installerId ?? vm.data.value?.buildId) {
    const db = serverConfigurator.db
    if (db == "perfint" || db == "perfintDev") {
      getTeamcityBuildType(db, serverConfigurator.table, vm.data.value.buildId, (type: string | null) => {
        if (vm.data.value) {
          window.open(`${tcUrl}buildConfiguration/${type}/${vm.data.value.buildId}?buildTab=artifacts#${encodeURIComponent(replaceUnderscore("/" + vm.data.value.projectName))}`)
        }
      })
    } else {
      window.open(vm.data.value.artifactsUrl)
    }
  }
}

function replaceUnderscore(project: string) {
  return project.replaceAll("_", "-")
}

function getSpaceUrl() {
  const db = serverConfigurator?.db
  if (db != null && (vm.data.value?.installerId ?? vm.data.value?.buildId)) {
    calculateChanges(db, vm.data.value.installerId ?? vm.data.value.buildId, (decodedChanges: string | null) => {
      if (decodedChanges == null || decodedChanges.length === 0) {
        console.log("No changes found")
        window.open(vm.data.value?.changesUrl)
      } else {
        window.open(`https://jetbrains.team/p/ij/repositories/ultimate/commits?query=%22${decodedChanges}%22&tab=changes`)
      }
    })
  }
}

function getAccidentTypes(): string[] {
  const values = Object.values(AccidentKind)
  //don't report Inferred type manually
  values.splice(values.indexOf(AccidentKind.InferredRegression), 1)
  values.splice(values.indexOf(AccidentKind.InferredImprovement), 1)
  return values
}

useScrollListeners()
</script>
<style>
.infoSidebar {
  min-width: 350px;
  max-width: 25%;
}

.infoSidebar_icon::after {
  position: absolute;
  content: "";
  inset: -8px;
  transform: rotate(-45deg);
}

.extraMargin {
  margin-top: 7rem; /* Replace [HeightOfYourToolbar] with the actual height */
}

.p-splitbutton.p-button-text > .p-button {
  @apply text-gray-600 font-medium text-left border-t border-solid border-neutral-200 relative;
}

.p-tieredmenu .p-menuitem-content {
  @apply text-sm text-gray-600 font-medium text-left relative;
}
</style>
