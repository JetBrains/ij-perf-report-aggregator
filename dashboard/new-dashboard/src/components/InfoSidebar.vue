<template>
  <div
    v-show="vm.visible.value"
    class="infoSidebar ml-5 text-gray-500 relative"
  >
    <div class="flex flex-col gap-4 sticky top-2 border border-solid rounded-md border-zinc-200 p-5">
      <span class="flex justify-between uppercase text-xs">
        {{ vm.data.value?.title }}

        <span
          class="infoSidebar_icon text-sm pi pi-plus rotate-45 cursor-pointer hover:text-gray-800 transition duration-150 ease-out relative"
          @click="handleCloseClick"
        />
      </span>

      <VTooltip
        v-if="vm.data.value?.description.value?.description"
        theme="info"
      >
        <div class="flex gap-1.5 font-medium text-base items-center break-all">
          <span
            class="w-3 h-3 rounded-full"
            :style="{ backgroundColor: vm.data.value?.color }"
          />
          <span class="underline decoration-dotted hover:no-underline">{{ vm.data.value?.projectName }}</span>
        </div>
        <template #popper>
          <span class="text-sm">
            {{ vm.data.value?.description.value?.description }}
          </span>
        </template>
      </VTooltip>
      <div
        v-else
        class="flex gap-1.5 font-medium text-base items-center break-all"
      >
        <span
          class="w-3 h-3 rounded-full"
          :style="{ backgroundColor: vm.data.value?.color }"
        />
        {{ vm.data.value?.projectName }}
      </div>

      <SplitButton
        label="Navigate to test"
        :model="getTestActions()"
        text
        plain
        link
        icon="pi pi-chart-line"
        @click="handleNavigateToTest"
      />

      <div class="flex flex-col gap-2">
        <span class="flex gap-1.5 text-sm items-center">
          <CalendarIcon class="w-4 h-4" />
          {{ vm.data.value?.date }}
          <span v-if="vm.data.value?.build">build {{ vm.data.value?.build }}</span>
        </span>
        <span
          v-if="vm.data.value?.metric"
          class="flex gap-1.5 text-sm items-center"
        >
          <BeakerIcon class="w-4 h-4" />
          <span>{{ vm.data.value?.metric }}</span>
        </span>
        <span class="flex gap-1.5 text-sm items-center">
          <ClockIcon class="w-4 h-4" />
          {{ vm.data.value?.value }}
        </span>
        <span class="flex gap-1.5 text-sm items-center">
          <ComputerDesktopIcon class="w-4 h-4" />
          {{ vm.data.value?.machineName }}
        </span>

        <span
          v-if="vm.data.value?.accidents"
          class="flex gap-1.5 text-sm items-center"
        >
          <ExclamationTriangleIcon class="w-4 h-4" />
          Known events:
        </span>
        <ul
          v-if="vm.data.value?.accidents"
          class="gap-1.5 text-sm ml-5 overflow-y-auto max-h-80"
        >
          <li
            v-for="accident in vm.data.value?.accidents"
            :key="accident.id"
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

      <div class="flex gap-4 text-blue-500">
        <a
          :href="vm.data.value?.changesUrl"
          target="_blank"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-blue-600"
        >
          <ArrowPathIcon class="w-4 h-4" />
          Changes
        </a>
        <a
          :href="vm.data.value?.artifactsUrl"
          target="_blank"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-blue-600"
        >
          <ServerStackIcon class="w-4 h-4" />
          Test Artifacts
        </a>
        <a
          v-if="vm.data.value?.installerUrl !== undefined"
          :href="vm.data.value?.installerUrl"
          target="_blank"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-blue-600"
        >
          <ArrowDownTrayIcon class="w-4 h-4" />
          Installer
        </a>
      </div>
      <div class="flex gap-4 text-blue-500 justify-center">
        <a
          v-if="getSpaceUrl()"
          :href="getSpaceUrl()"
          target="_blank"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-blue-600"
        >
          <SpaceIcon class="w-4 h-4" />
          Changes
        </a>
      </div>

      <Button
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
    <div class="flex items-center space-x-4">
      <Dropdown
        v-model="accidentType"
        placeholder="Event Type"
        :options="getAccidentTypes()"
        class="w-[8rem]"
      />
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
    <template #footer>
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
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { inject, ref } from "vue"
import { useRouter } from "vue-router"
import { sidebarVmKey } from "../shared/keys"
import { getAccidentTypes, removeAccidentFromMetaDb, writeAccidentToMetaDb } from "../util/meta"
import { InfoSidebarVmImpl } from "./InfoSidebarVm"
import SpaceIcon from "./common/SpaceIcon.vue"

const vm = inject(sidebarVmKey) ?? new InfoSidebarVmImpl()
const showDialog = ref(false)
const reason = ref("")
const router = useRouter()
const accidentType = ref<string>("Regression")

function reportRegression() {
  showDialog.value = false
  const value = vm.data.value
  if (value == null) {
    console.log("value is zero! This shouldn't happen")
  } else {
    writeAccidentToMetaDb(value.date, value.projectName, reason.value, value.build ?? value.buildId.toString(), accidentType.value)
  }
}

function copyMethodNameToClipboard(methodName: string) {
  void navigator.clipboard.writeText(methodName)
}

function openTestInIDE(methodName: string) {
  const origin = encodeURIComponent("ssh://git@git.jetbrains.team/ij/intellij.git")
  window.open(`jetbrains://idea/navigate/reference?origin=${origin}&fqn=${methodName}`)
}

function handleNavigateToTest() {
  const currentRoute = router.currentRoute.value
  const parts = currentRoute.path.split("/")
  parts[parts.length - 1] = parts.at(-1)?.toLowerCase().endsWith("dev") ? "testsDev" : "tests"
  const testURL = parts.join("/")
  const query: Record<string, string> = { ...currentRoute.query, project: vm.data.value?.projectName ?? "" } as Record<string, string>
  const queryParams: string = new URLSearchParams(query).toString()
  void router.push(testURL + "?" + queryParams)
}

function handleRemove(id: number) {
  removeAccidentFromMetaDb(id)
}

function handleCloseClick() {
  vm.close()
}

function getTestActions() {
  const actions = []
  if (vm.data.value?.description.value) {
    const url = vm.data.value.description.value.url
    if (url && url != "") {
      actions.push({
        label: "Download test project",
        icon: "pi pi-download",
        command() {
          window.open(url)
        },
      })
    }
    const methodName = vm.data.value.description.value.methodName
    if (methodName && methodName != "") {
      actions.push(
        {
          label: "Copy test method name",
          icon: "pi pi-copy",
          command() {
            copyMethodNameToClipboard(methodName)
          },
        },
        {
          label: "Open test method",
          icon: "pi pi-folder-open",
          command() {
            openTestInIDE(methodName)
          },
        }
      )
    }
  }
  return actions
}

function getSpaceUrl() {
  if (vm.data.value?.changes != null) {
    return "https://jetbrains.team/p/ij/repositories/intellij/commits?query=%22" + vm.data.value.changes + "%22&tab=changes"
  }
  return
}
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

.p-splitbutton.p-button-text > .p-button {
  @apply text-gray-600 font-medium text-left border-t border-solid border-neutral-200 relative;
}

.p-tieredmenu .p-menuitem-content {
  @apply text-sm text-gray-600 font-medium text-left relative;
}
</style>
