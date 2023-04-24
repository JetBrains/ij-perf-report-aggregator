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

      <div class="flex gap-1.5 font-medium text-base items-center break-all">
        <span
          class="w-3 h-3 rounded-full"
          :style="{ backgroundColor: vm.data.value?.color }"
        />
        {{ vm.data.value?.projectName }}
        <span
          class="infoSidebar_icon text-sm pi pi-external-link cursor-pointer hover:text-gray-800 transition duration-150 ease-out relative"
          @click="handleNavigateToTest"
        />
      </div>

      <div class="flex flex-col gap-2">
        <span class="flex gap-1.5 text-sm items-center">
          <CalendarIcon class="w-4 h-4" />
          {{ vm.data.value?.date }} <span v-if="vm.data.value?.build"> build {{ vm.data.value?.build }} </span>
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
          <ExclamationTriangleIcon class="w-4 h-4 text-red-500" /> Known degradation:
        </span>
        <ul
          v-if="vm.data.value?.accidents"
          class="flex gap-1.5 text-sm"
        >
          <li
            v-for="accident in vm.data.value?.accidents"
            :key="accident.id"
          >
            <span class="flex gap-1.5 text-sm items-center">&bull; {{ accident.reason }} <TrashIcon
              class="w-4 h-4 text-red-500"
              @click="handleRemove(accident.id)"
            /></span>
          </li>
        </ul>
      </div>

      <div class="flex gap-4 text-blue-500">
        <a
          :href="vm.data.value?.changesUrl"
          target="_blank"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-blue-600"
        >
          <ArrowPathIcon class="w-4 h-4" /> Changes
        </a>
        <a
          :href="vm.data.value?.artifactsUrl"
          target="_blank"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-blue-600"
        >
          <ServerStackIcon class="w-4 h-4" /> Test Artifacts
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
      <div
        v-if="vm.data.value?.changes"
        class="text-sm"
      >
        Changes: {{ vm.data.value?.changes }}
      </div>
      <Button
        class="text-sm"
        label="Report regression"
        text
        size="small"
        severity="danger"
        @click="showDialog = true"
      />
    </div>
  </div>
  <Dialog
    v-model:visible="showDialog"
    modal
    header="Report Regression"
    :style="{ width: '30vw' }"
  >
    <span class="p-float-label">
      <InputText
        id="reason"
        v-model="reason"
        class="w-full"
      />
      <label for="reason">Reason</label>
    </span>
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
import { removeAccidentFromMetaDb, writeAccidentToMetaDb } from "shared/src/meta"
import { inject, ref } from "vue"
import { useRouter } from "vue-router"
import { sidebarVmKey } from "../shared/keys"
import { InfoSidebarVmImpl } from "./InfoSidebarVm"

const vm = inject(sidebarVmKey) || new InfoSidebarVmImpl()
const showDialog = ref(false)
const reason = ref("")
const router = useRouter()

function reportRegression(){
  showDialog.value = false
  const value = vm.data.value
  if(value == null){
    console.log("value is zero! This shouldn't happen")
  } else {
    writeAccidentToMetaDb(value.date, value.projectName, reason.value, value.build ?? value.buildId.toString())
  }
}

function handleNavigateToTest(){
  const currentRoute = router.currentRoute.value
  const parts = currentRoute.path.split("/")
  parts[parts.length - 1] = parts[parts.length - 2].toLowerCase().endsWith("dev") ? "testsDev" : "tests"
  const testURL = parts.join("/")
  const query: Record<string, string> = { ...currentRoute.query, project: vm.data.value?.projectName ?? "" }
  const queryParams: string = new URLSearchParams(query).toString()
  void router.push(testURL+"?"+queryParams)
}

function handleRemove(id: number){
  removeAccidentFromMetaDb(id)
}

function handleCloseClick() {
  vm.close()
}
</script>
<style>
.infoSidebar {
  min-width: 350px;
  max-width: 25%;
}

.infoSidebar_icon::after {
  position: absolute;
  content: '';
  inset: -8px;
  transform: rotate(-45deg);
}
</style>