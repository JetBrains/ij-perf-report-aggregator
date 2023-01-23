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

      <div class="flex gap-1.5 font-medium text-base items-center">
        <span
          class="w-3 h-3 rounded-full"
          :style="{ backgroundColor: vm.data.value?.color }"
        />
        {{ vm.data.value?.projectName }}
      </div>

      <div class="flex flex-col gap-2">
        <span class="flex gap-1.5 text-sm items-center">
          <CalendarIcon class="w-4 h-4" />
          {{ vm.data.value?.date }} build {{ vm.data.value?.build }}
        </span>
        <span class="flex gap-1.5 text-sm items-center">
          <ClockIcon class="w-4 h-4" />
          {{ vm.data.value?.duration }}
        </span>
        <span class="flex gap-1.5 text-sm items-center">
          <ComputerDesktopIcon class="w-4 h-4" />
          {{ vm.data.value?.machineName }}
        </span>
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
          <ArrowDownTrayIcon class="w-4 h-4" /> Installer
        </a>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { inject } from "vue"
import { sidebarVmKey } from "../shared/keys"
import { InfoSidebarVmImpl } from "./InfoSidebarVm"

const vm = inject(sidebarVmKey) || new InfoSidebarVmImpl()

function handleCloseClick() {
  vm.close()
}
</script>
<style>
.infoSidebar {
  min-width: 350px;
}

.infoSidebar_icon::after {
  position: absolute;
  content: '';
  inset: -8px;
  transform: rotate(-45deg);
}
</style>