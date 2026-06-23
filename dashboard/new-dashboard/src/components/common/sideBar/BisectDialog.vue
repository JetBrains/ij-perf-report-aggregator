<template>
  <Dialog
    v-model:visible="showDialog"
    modal
    header="Run bisect"
    :style="{ width: '40vw' }"
  >
    <div class="flex flex-col space-y-8 mb-4 mt-6">
      <FloatLabel>
        <InputText
          id="targetValue"
          v-model="targetValue"
          :invalid="!isTargetValueValid()"
        />
        <label for="targetValue">Target value in milliseconds</label>
      </FloatLabel>
      <SelectButton
        v-model="mode"
        :options="modeOptions"
      />
      <div
        v-if="isCommitMode"
        class="flex"
      >
        <FloatLabel class="flex-1 mr-4">
          <InputText
            id="changes"
            v-model="firstCommit"
          />
          <label for="changes">First commit</label>
        </FloatLabel>
        <FloatLabel class="flex-1">
          <InputText
            id="changes"
            v-model="lastCommit"
          />
          <label for="changes">Last commit</label>
        </FloatLabel>
      </div>
      <FloatLabel>
        <Select
          id="direction"
          v-model="direction"
          :options="['OPTIMIZATION', 'DEGRADATION']"
        >
          <template #value="{ value }">
            <div class="group inline-flex justify-center text-sm font-medium">
              {{ value }}
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
        <label for="direction">Direction</label>
      </FloatLabel>
      <FloatLabel class="w-full">
        <InputText
          id="excludedCommits"
          v-model="excludedCommits"
          class="w-full"
        />
        <label for="excludedCommits">List of excluded commits. Ex: 805bfa9758dec2912dcfecba,c7ee80058a9182c3037ee883615</label>
      </FloatLabel>
      <div
        v-if="changesGap?.hasGap"
        class="rounded border border-amber-400 bg-amber-50 p-3 text-sm"
      >
        <div class="flex items-center font-medium text-amber-700 mb-2">
          <ExclamationTriangleIcon class="h-5 w-5 mr-2 shrink-0" />
          Incomplete change range
        </div>
        <p class="text-amber-800 mb-3">
          This build does not include all changes since the previous data point on the graph.
          {{ changesGap.gapCommitCount }} commit{{ changesGap.gapCommitCount === 1 ? "" : "s" }} landed in builds that produced no data point (e.g. failed or timed-out runs) and
          {{ changesGap.gapCommitCount === 1 ? "is" : "are" }} not part of this bisect range.
        </p>
        <div class="flex items-center">
          <Checkbox
            id="acknowledgeGap"
            v-model="acknowledgedGap"
            binary
          />
          <label
            for="acknowledgeGap"
            class="ml-2 text-amber-800"
            >I understand the regression may have been introduced by a change outside this range</label
          >
        </div>
      </div>
      <Accordion>
        <AccordionPanel value="0">
          <AccordionHeader>Additional parameters</AccordionHeader>
          <AccordionContent>
            <div class="flex flex-col space-y-8 mb-4 mt-4">
              <FloatLabel>
                <InputText
                  id="testPatterns"
                  v-model="fullClassName"
                />
                <label for="testPatterns">Test FQN patterns</label>
              </FloatLabel>
              <FloatLabel>
                <InputText
                  id="requester"
                  v-model="requester"
                  :disabled="!!userEmail"
                />
                <label for="requester">Requester</label>
              </FloatLabel>
              <FloatLabel v-if="!isCommitMode">
                <InputText
                  id="buildId"
                  v-model="buildId"
                />
                <label for="buildId">Build ID</label>
              </FloatLabel>
              <FloatLabel>
                <InputText
                  id="test"
                  v-model="test"
                />
                <label for="test">Test name</label>
              </FloatLabel>
              <FloatLabel>
                <InputText
                  id="metric"
                  v-model="metric"
                />
                <label for="metric">Metric name</label>
              </FloatLabel>
              <FloatLabel>
                <InputText
                  id="buildType"
                  v-model="buildType"
                />
                <label for="buildType">Build type</label>
              </FloatLabel>
            </div>
            <div class="flex items-center mb-4 mt-4">
              <Checkbox
                id="targetJpsCompile"
                v-model="targetJpsCompile"
                binary
              />
              <label
                for="targetJpsCompile"
                class="ml-2"
                >JPS compilation</label
              >
            </div>
          </AccordionContent>
        </AccordionPanel>
      </Accordion>
    </div>
    <div
      v-if="error"
      class="text-red-500 mb-4"
    >
      {{ error }}
    </div>
    <!-- Footer buttons -->
    <template #footer>
      <div class="flex justify-end space-x-2">
        <Button
          label="Cancel"
          icon="pi pi-times"
          severity="secondary"
          @click="showDialog = false"
        />
        <Button
          v-tooltip.top="
            reasonOfDisabling === ''
              ? null
              : {
                  value: reasonOfDisabling,
                  autoHide: false,
                }
          "
          label="Start"
          icon="pi pi-play"
          autofocus
          :loading="loading"
          :disabled="reasonOfDisabling !== ''"
          @click="startBisect"
        />
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { getNavigateToTestUrl, InfoData } from "./InfoSidebar"
import { getTeamcityBuildType } from "../../../util/artifacts"
import { injectOrError } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"
import { computedAsync } from "@vueuse/core"
import { computed, onMounted, Ref, ref } from "vue"
import { ChevronDownIcon, ExclamationTriangleIcon } from "@heroicons/vue/20/solid/index"
import { BisectClient } from "./BisectClient"
import { useUserStore } from "../../../shared/useUserStore"
import { getFirstAndLastCommit } from "../../../util/changes"
import { getPersistentLink } from "../../settings/CopyLink"
import { TimeRangeConfigurator } from "../../../configurators/TimeRangeConfigurator"
import { useRouter } from "vue-router"

const { data, timerangeConfigurator } = defineProps<{
  data: InfoData
  timerangeConfigurator: TimeRangeConfigurator
}>()

const serverConfigurator = injectOrError(serverConfiguratorKey)

const mode = ref("Build")
const modeOptions = ref(["Build", "Commits"])
const isCommitMode = computed(() => mode.value === "Commits")

const showDialog = defineModel<boolean>("showDialog")
const metric = ref(data.series[0].metricName ?? "")
const test = ref(data.projectName)
const isDegradation = data.deltaPrevious?.includes("-") ?? false
const direction = ref(isDegradation ? "DEGRADATION" : "OPTIMIZATION")
const buildType = computedAsync(
  () => getTeamcityBuildType(serverConfigurator.db, serverConfigurator.table, data.buildId).then((bt) => (fullClassName && bt ? bt.replace(/_\d+$/, "_1") : bt)),
  null
)
const buildId = ref(data.buildId.toString())
const userEmail = useUserStore().user?.email
const requester = ref(userEmail)
const methodName = data.description.value?.methodName ?? ""
const fullClassName = ref(methodName.slice(0, Math.max(0, methodName.lastIndexOf("#"))))
const targetValue: Ref<string | null> = ref(null)
const excludedCommits = ref("")
const targetJpsCompile = ref(data.branch === "master" && new Date(data.date) <= new Date("2025-10-19T23:59:59.999Z"))

const firstCommit = ref()
const lastCommit = ref()
const { firstCommit: first, lastCommit: last } = await getFirstAndLastCommit(serverConfigurator.db, data.installerId ?? data.buildId)
firstCommit.value = first
lastCommit.value = last

const router = useRouter()
const bisectClient = new BisectClient(serverConfigurator)

// Warn when the current build does not include all changes since the previous
// successful dot: intervening builds that failed/timed out produced no data point,
// so the commits they consumed are absent from this bisect range.
//
// Only meaningful for source-based configurations: when the config runs on an
// installer (installerId is set), the perf build carries no VCS changes of its own
// (they live on the installer build), so a build-level gap check would be moot.
const checkingGap = ref(true)
const changesGap = computedAsync(
  () =>
    data.buildIdPrevious == null || data.installerId != undefined || !firstCommit.value
      ? Promise.resolve(null)
      : bisectClient.fetchChangesGap(data.buildId.toString(), data.buildIdPrevious.toString(), firstCommit.value as string),
  null,
  checkingGap
)
const acknowledgedGap = ref(false)
const dashboardLink = computed(() => window.location.origin + getPersistentLink(getNavigateToTestUrl(data, router), timerangeConfigurator))

onMounted(() => {
  console.log("IJ Perf link for bisect:", dashboardLink.value)
})

const error = ref<string | null>(null)
const loading = ref(false)

function isTargetValueValid() {
  const value = Number(targetValue.value)
  return targetValue.value !== null && targetValue.value !== "" && Number.isInteger(value)
}

const reasonOfDisabling = computed(() => {
  // Bisect checks out individual source commits, but installer-based configurations
  // run on periodically-built installers rather than per-commit, so their data points
  // can't be mapped to a bisectable source range.
  if (data.installerId != undefined) {
    return "Bisect is not supported for configurations running on installers"
  }
  if (firstCommit.value === "" || lastCommit.value === "") {
    return "Build has no changes"
  }
  if (!isTargetValueValid()) {
    return "Target value must be a valid integer"
  }
  if (checkingGap.value) {
    return "Checking whether the build includes all changes since the previous run…"
  }
  if (changesGap.value?.hasGap && !acknowledgedGap.value) {
    return "Please acknowledge that this build doesn't include all changes since the previous successful run"
  }
  return ""
})

async function startBisect() {
  error.value = null
  loading.value = true
  try {
    //todo add validation on all values
    const weburl = await bisectClient.sendBisectRequest({
      targetValue: targetValue.value as string,
      requester: requester.value ?? "",
      changes: (firstCommit.value as string) + "^.." + (lastCommit.value as string),
      buildId: buildId.value,
      mode: isCommitMode.value ? "commit" : "build",
      direction: direction.value,
      test: test.value,
      metric: metric.value,
      buildType: buildType.value as string,
      testPatterns: fullClassName.value,
      excludedCommits: excludedCommits.value
        .split(",")
        .map((commit) => commit.trim())
        .filter((commit) => commit !== "")
        .join(","),
      jpsCompilation: targetJpsCompile.value ? "true" : "false",
      dashboardLink: dashboardLink.value,
    })
    showDialog.value = false // Close dialog on success
    window.open(weburl, "_blank")
  } catch (error_) {
    error.value = error_ instanceof Error ? error_.message : "An unknown error occurred"
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.p-inputtext {
  @apply w-full;
}

.p-float-label {
  @apply w-full;
}
</style>
