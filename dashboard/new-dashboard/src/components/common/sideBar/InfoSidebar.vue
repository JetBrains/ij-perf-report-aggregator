<template>
  <div
    v-show="vm.visible.value"
    class="infoSidebar ml-5 text-gray-500 relative"
  >
    <div class="infoSidebarContent flex flex-col gap-4 sticky top-2 border border-solid rounded-md border-zinc-200 p-5 overflow-y-auto overflow-x-hidden">
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

      <div class="flex gap-1.5 font-medium text-base items-center break-all">
        <span
          v-if="data?.series.length == 1"
          class="w-3 h-3 rounded-full"
          :style="{ backgroundColor: data?.series[0].color }"
        />
        <span v-if="data?.series.length == 1">
          {{ data?.seriesName }}
        </span>
        <span v-else>
          {{ data?.projectName }}
        </span>
      </div>

      <TestActions :data="data" />

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
          <span class="flex gap-1.5 text-sm items-center">
            <ChartBarIcon class="w-4 h-4" />
            <span
              v-tooltip.left="description"
              :class="description != '' ? getURLStyle() : 'break-all'"
            >
              {{ data?.projectName }}
            </span>
          </span>
          <span
            v-if="data?.series[0].metricName"
            class="flex gap-1.5 text-sm items-center"
          >
            <BeakerIcon class="w-4 h-4" />
            <span
              v-tooltip.left="getTooltipForMetric(data?.series[0].metricName)"
              :class="metricDescription != null ? getURLStyle() : ''"
              >{{ data?.series[0].metricName }}</span
            >
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
                v-if="item.metricName"
                class="rounded-lg w-2.5 h-2.5"
                :style="{ 'background-color': item.color }"
              />
              <span v-if="item.metricName">{{ item.metricName }}</span>
              <span
                v-if="item.metricName"
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
              <!-- eslint-disable vue/no-v-html -->
              <span
                class="w-full"
                :class="accident.kind == 'Regression' || accident.kind == 'InferredRegression' ? 'text-red-500' : 'text-green-500'"
                v-html="replaceToLink(accident.reason)"
              />
              <GlobeAltIcon
                v-if="accident.affectedTest == ''"
                class="w-4 h-4 flex-none"
              />
              <!-- eslint-enable -->
              <PencilSquareIcon
                class="w-4 h-4 flex-none"
                @click="editAccident(accident)"
              />
            </span>
          </li>
        </ul>
      </div>

      <RelatedAccidents
        :data="data"
        :accidents-configurator="accidentsConfigurator"
        :in-dialog="false"
      />

      <div class="flex gap-4 text-primary">
        <a
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-darker cursor-pointer"
          @click="getChangesUrl"
          @click.middle="getChangesUrl"
        >
          <ArrowPathIcon class="w-4 h-4" />
          Changes
        </a>
        <a
          class="flex gap-1.5 items-center transition duration-150 ease-out hover:text-darker cursor-pointer"
          @click="getArtifactsUrl"
          @click.middle="getArtifactsUrl"
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
          @click.middle="getSpaceUrl"
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
        @click="createAccident()"
      />
    </div>
  </div>
  <ReportMetricDialog
    v-if="showDialog"
    v-model="showDialog"
    :accidents-configurator="accidentsConfigurator"
    :accident-to-edit="accidentToEdit"
    :data="data"
  />
</template>
<script setup lang="ts">
import { computed, Ref, ref } from "vue"
import { Accident } from "../../../configurators/AccidentsConfigurator"
import { injectOrError, injectOrNull } from "../../../shared/injectionKeys"
import { accidentsConfiguratorKey, serverConfiguratorKey, sidebarVmKey } from "../../../shared/keys"
import { getMetricDescription } from "../../../shared/metricsDescription"
import { getTeamcityBuildType } from "../../../util/artifacts"
import { calculateChanges } from "../../../util/changes"
import { replaceToLink } from "../../../util/linkReplacer"
import BranchIcon from "../BranchIcon.vue"
import SpaceIcon from "../SpaceIcon.vue"
import { useScrollListeners, useScrollStore } from "../scrollStore"
import { tcUrl } from "./InfoSidebar"
import RelatedAccidents from "./RelatedAccidents.vue"
import ReportMetricDialog from "./ReportMetricDialog.vue"
import TestActions from "./TestActions.vue"

const vm = injectOrError(sidebarVmKey)
const showDialog = ref(false)
const accidentToEdit: Ref<Accident | null> = ref(null)

const serverConfigurator = injectOrNull(serverConfiguratorKey)

const accidentsConfigurator = injectOrNull(accidentsConfiguratorKey)

const data = computed(() => vm.data.value)

function editAccident(accident: Accident) {
  showDialog.value = true
  accidentToEdit.value = accident
}
function createAccident() {
  showDialog.value = true
  accidentToEdit.value = null
}

function handleCloseClick() {
  vm.close()
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

useScrollListeners()
const description = computed(() => vm.data.value?.description.value?.description ?? "")

const metricDescription = computed(() => getMetricDescription(data.value?.series[0].metricName))

function getTooltipForMetric(metricName: string | undefined) {
  const metricInfo = getMetricDescription(metricName)

  return metricInfo == null
    ? null
    : {
        value:
          metricInfo.description +
          (metricInfo.url ? "<br/><a class='text-xs underline decoration-dotted hover:no-underline' href='" + metricInfo.url + "' target='_blank'>More info</a>" : ""),
        escape: false,
        hideDelay: metricInfo.url ? 3000 : 0,
        autoHide: !metricInfo.url,
      }
}

function getURLStyle() {
  return "underline decoration-dotted hover:no-underline"
}
</script>
<style>
.infoSidebar {
  min-width: 350px;
  max-width: 25%;
}

.infoSidebarContent {
  height: calc(100vh - 15em);
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
</style>
