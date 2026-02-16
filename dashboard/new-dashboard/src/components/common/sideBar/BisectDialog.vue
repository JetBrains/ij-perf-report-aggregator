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
                  :disabeld="requester !== undefined && requester !== ''"
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
import { InfoData } from "./InfoSidebar"
import { getTeamcityBuildType } from "../../../util/artifacts"
import { injectOrError } from "../../../shared/injectionKeys"
import { serverConfiguratorKey } from "../../../shared/keys"
import { computedAsync } from "@vueuse/core"
import { computed, Ref, ref } from "vue"
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { BisectClient } from "./BisectClient"
import { useUserStore } from "../../../shared/useUserStore"
import { calculateChanges } from "../../../util/changes"

const { data } = defineProps<{
  data: InfoData
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
const buildType = computedAsync(async () => await getTeamcityBuildType(serverConfigurator.db, serverConfigurator.table, data.buildId), null)
const buildId = ref(data.buildId.toString())
const requester = ref(useUserStore().user?.email)
const methodName = data.description.value?.methodName ?? ""
const fullClassName = methodName.slice(0, Math.max(0, methodName.lastIndexOf("#")))
const targetValue: Ref<string | null> = ref(null)
const excludedCommits = ref("")
const targetJpsCompile = ref(data.branch === "master" && new Date(data.date) <= new Date("2025-10-19T23:59:59.999Z"))

const firstCommit = ref()
const lastCommit = ref()
const changesMerged = await calculateChanges(serverConfigurator.db, data.installerId ?? data.buildId)
const changesUnmerged = changesMerged?.flatMap((chunk) => chunk.split("%2C")) as string[] | null
if (Array.isArray(changesUnmerged)) {
  firstCommit.value = changesUnmerged.at(-1) ?? null
  lastCommit.value = changesUnmerged[0] ?? null
}

const bisectClient = new BisectClient(serverConfigurator)

const error = ref<string | null>(null)
const loading = ref(false)

function isTargetValueValid() {
  const value = Number(targetValue.value)
  return targetValue.value !== null && targetValue.value !== "" && Number.isInteger(value)
}

const reasonOfDisabling = computed(() => {
  if (firstCommit.value === "" || lastCommit.value === "") {
    return "Build has no changes"
  }
  if (!isTargetValueValid()) {
    return "Target value must be a valid integer"
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
      testPatterns: fullClassName,
      excludedCommits: excludedCommits.value
        .split(",")
        .map((commit) => commit.trim())
        .filter((commit) => commit !== "")
        .join(","),
      jpsCompilation: targetJpsCompile.value ? "true" : "false",
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
