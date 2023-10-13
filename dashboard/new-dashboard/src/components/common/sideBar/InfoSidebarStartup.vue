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
        {{ vm.data.value?.projectName }}
      </div>

      <div class="grid grid-cols-[repeat(3,_max-content)] whitespace-nowrap gap-x-2 items-baseline leading-loose text-sm">
        <template
          v-for="item in vm.data.value?.series"
          :key="item.metricName"
        >
          <span
            v-if="item.metricName"
            class="rounded-lg w-2.5 h-2.5"
            :style="{ 'background-color': item.color }"
          />
          <span v-if="item.metricName">{{ item.metricName.replace("metrics.", "") }}</span>
          <span
            v-if="item.metricName"
            class="font-mono place-self-end"
            >{{ item.value }}</span
          >
        </template>
      </div>

      <div class="flex flex-col gap-2">
        <span class="flex gap-1.5 text-sm items-center">
          <CalendarIcon class="w-4 h-4" />
          {{ vm.data.value?.date }}
          <span v-if="vm.data.value?.build">build {{ vm.data.value?.build }}</span>
        </span>
        <span class="flex gap-1.5 text-sm items-center">
          <BranchIcon class="w-4 h-4" />
          {{ vm.data.value?.branch }}
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
          v-if="vm.data.value?.installerId"
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-blue-600 cursor-pointer"
          @click="getSpaceUrl"
        >
          <SpaceIcon class="w-4 h-4" />
          Changes
        </a>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { injectOrError } from "../../../shared/injectionKeys"
import { sidebarStartupKey } from "../../../shared/keys"
import { calculateChanges } from "../../../util/changes"
import BranchIcon from "../BranchIcon.vue"
import SpaceIcon from "../SpaceIcon.vue"

const vm = injectOrError(sidebarStartupKey)

function getSpaceUrl() {
  if (vm.data.value?.installerId) {
    calculateChanges("ij", vm.data.value.installerId, (decodedChanges: string | null) => {
      window.open(`https://jetbrains.team/p/ij/repositories/intellij/commits?query=%22${decodedChanges}%22&tab=changes`)
    })
  }
}

function handleCloseClick() {
  vm.close()
}
</script>
<style>
.infoSidebar {
  min-width: 350px;
  max-width: 35%;
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
