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
        <label for="targetValue">Target value</label>
      </FloatLabel>
      <div class="flex">
        <FloatLabel>
          <InputText
            id="changes"
            v-model="firstCommit"
          />
          <label for="changes">First commit</label>
        </FloatLabel>
        <FloatLabel>
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
                class="-mr-1 ml-1 h-5 w-5 flex-shrink-0"
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
      <Accordion>
        <AccordionPanel value="0">
          <AccordionHeader>Additional parameters</AccordionHeader>
          <AccordionContent>
            <div class="flex flex-col space-y-8 mb-4 mt-4">
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

              <FloatLabel>
                <InputText
                  id="className"
                  v-model="className"
                />
                <label for="className">Class name</label>
              </FloatLabel>
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
          label="Start"
          icon="pi pi-play"
          autofocus
          :loading="loading"
          :disabled="!isTargetValueValid()"
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
import { Ref, ref } from "vue"
import { calculateChanges } from "../../../util/changes"
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { BisectClient } from "./BisectClient"

const { data } = defineProps<{
  data: InfoData
}>()

const serverConfigurator = injectOrError(serverConfiguratorKey)

const showDialog = defineModel<boolean>("showDialog")
const metric = ref(data.series[0].metricName ?? "")
const test = ref(data.projectName)
const isDegradation = data.deltaPrevious?.includes("-") ?? false
const direction = ref(isDegradation ? "DEGRADATION" : "OPTIMIZATION")
const buildType = computedAsync(async () => await getTeamcityBuildType(serverConfigurator.db, serverConfigurator.table, data.buildId), null)
const firstCommit = ref()
const lastCommit = ref()
const changesMerged = await calculateChanges(serverConfigurator.db, data.installerId ?? data.buildId)
const changesUnmerged = changesMerged?.split("%2C") as string[] | null
if (Array.isArray(changesUnmerged)) {
  firstCommit.value = changesUnmerged.at(-1) ?? null
  lastCommit.value = changesUnmerged[0] ?? null
}
const methodName = data.description.value?.methodName ?? ""
const fullClassName = methodName.slice(0, Math.max(0, methodName.lastIndexOf("#")))
const className = fullClassName.slice(fullClassName.lastIndexOf(".") + 1)
const targetValue: Ref<string | null> = ref(null)

const bisectClient = new BisectClient(serverConfigurator)

const error = ref<string | null>(null)
const loading = ref(false)

function isTargetValueValid() {
  const value = Number(targetValue.value)
  return targetValue.value !== null && targetValue.value !== "" && Number.isInteger(value)
}

async function startBisect() {
  error.value = null
  loading.value = true
  try {
    //todo add validation on all values
    const weburl = await bisectClient.sendBisectRequest({
      targetValue: targetValue.value as string,
      changes: (firstCommit.value as string) + "^.." + (lastCommit.value as string),
      direction: direction.value,
      test: test.value,
      metric: metric.value,
      buildType: buildType.value as string,
      className,
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
