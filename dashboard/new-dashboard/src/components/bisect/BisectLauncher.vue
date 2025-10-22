<template>
  <div class="flex justify-center">
    <div class="mt-3 border-1 w-1/2">
      <Card>
        <template #title>
          <div class="text-darker">Start Bisect</div>
        </template>
        <template #content>
          <div class="flex flex-col space-y-8 mb-4 mt-6">
            <div class="grid grid-cols-4">
              <div class="col-span-4 mb-4">
                <FloatLabel class="w-full">
                  <InputText
                    id="errorMessage"
                    v-model="model.errorMessage"
                    class="w-full"
                  />
                  <label for="errorMessage">Error message to look in the build log</label>
                </FloatLabel>
              </div>
              <div class="col-span-4 mb-4">
                <SelectButton
                  v-model="mode"
                  :options="modeOptions"
                />
              </div>
              <div
                v-if="isCommitMode"
                class="col-span-2 mb-4 mt-4 mr-4"
              >
                <FloatLabel class="w-full">
                  <InputText
                    id="firstCommit"
                    v-model="model.firstCommit"
                    class="w-full"
                  />
                  <label for="firstCommit">First commit</label>
                </FloatLabel>
              </div>
              <div
                v-if="isCommitMode"
                class="col-span-2 mb-4 mt-4"
              >
                <FloatLabel class="w-full">
                  <InputText
                    id="lastCommit"
                    v-model="model.lastCommit"
                    class="w-full"
                  />
                  <label for="lastCommit">Last commit</label>
                </FloatLabel>
              </div>
              <FloatLabel class="col-span-4 w-full mt-4">
                <InputText
                  id="excludedCommits"
                  v-model="model.excludedCommits"
                  class="w-full"
                />
                <label for="excludedCommits">List of excluded commits. Ex: 805bfa9758dec2912dcfecba,c7ee80058a9182c3037ee883615</label>
              </FloatLabel>
              <div class="mt-6 col-span-4">
                <FloatLabel class="w-full">
                  <InputText
                    id="className"
                    v-model="model.className"
                    class="w-full"
                  />
                  <label for="className">Short class name. Example: ArrowKtScriptPerformanceTest</label>
                </FloatLabel>
              </div>
            </div>
            <Accordion :value="-1">
              <AccordionPanel :value="0">
                <AccordionHeader>Additional parameters</AccordionHeader>
                <AccordionContent>
                  <div class="mb-2 mt-4">
                    <FloatLabel class="w-full">
                      <InputText
                        id="requester"
                        v-model="email"
                        class="w-full"
                        :disabled="email !== undefined && email !== ''"
                      />
                      <label for="requester">Requester</label>
                    </FloatLabel>
                  </div>
                  <div
                    v-if="!isCommitMode"
                    class="mb-2 mt-6"
                  >
                    <FloatLabel class="w-full">
                      <InputText
                        id="buildID"
                        v-model="model.buildId"
                        class="w-full"
                      />
                      <label for="buildID">Build ID</label>
                    </FloatLabel>
                  </div>
                  <div class="col-span-4 mt-6">
                    <FloatLabel class="w-full">
                      <InputText
                        id="buildType"
                        v-model="model.buildType"
                        class="w-full"
                      />
                      <label for="buildType">Build type</label>
                    </FloatLabel>
                  </div>
                  <div class="flex items-center mb-4 mt-4">
                    <Checkbox
                      id="targetJpsCompile"
                      v-model="model.targetJpsCompile"
                      binary
                    />
                    <label
                      for="targetJpsCompile"
                      class="ml-2"
                      >JPS compilation
                    </label>
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
        </template>
        <template #footer>
          <div class="flex justify-end">
            <Button
              label="Start"
              icon="pi pi-play"
              autofocus
              :loading="loading"
              :disabled="isRequiredFieldEmpty"
              @click="startBisect"
            />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue"
import { BisectClient } from "../common/sideBar/BisectClient"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { useUserStore } from "../../shared/useUserStore"

const {
  buildId = "",
  errorMessage = "",
  className = "",
} = defineProps<{
  buildId?: string
  errorMessage?: string
  className?: string
}>()

const serverConfigurator = new ServerWithCompressConfigurator("", "")
const bisectClient = new BisectClient(serverConfigurator)

const model = reactive({
  errorMessage,
  firstCommit: "",
  lastCommit: "",
  buildType: "",
  excludedCommits: "",
  buildId,
  className,
  targetJpsCompile: false,
})

const mode = ref("Build")
const modeOptions = ref(["Build", "Commits"])
const isCommitMode = computed(() => mode.value === "Commits")

const isRequiredFieldEmpty = computed(() => model.errorMessage.trim() == "" || model.buildType.trim() == "")

const error = ref<string | null>(null)
const loading = ref(false)

async function startBisect() {
  error.value = null
  loading.value = true
  try {
    const weburl = await bisectClient.sendBisectRequest({
      buildId: model.buildId,
      changes: model.firstCommit + "^.." + model.lastCommit,
      requester: email.value ?? "",
      mode: isCommitMode.value ? "commit" : "build",
      errorMessage: model.errorMessage,
      buildType: model.buildType,
      className: model.className,
      excludedCommits: model.excludedCommits
        .split(",")
        .map((commit) => commit.trim())
        .filter((commit) => commit !== "")
        .join(","),
      jpsCompilation: model.targetJpsCompile ? "true" : "false",
    })
    window.open(weburl, "_blank")
  } catch (error_) {
    error.value = error_ instanceof Error ? error_.message : "An unknown error occurred"
  } finally {
    loading.value = false
  }
}

const email = computed(() => useUserStore().user?.email)
onMounted(async () => {
  try {
    const changes = await bisectClient.fetchTeamCityChanges(buildId)
    model.firstCommit = changes.firstCommit
    model.lastCommit = changes.lastCommit
    model.buildType = await bisectClient.fetchBuildType(buildId)

    const buildInfo = await bisectClient.fetchBuildInfo(buildId)
    if (buildInfo !== null) {
      const branch = buildInfo.branchName
      const dateStr = buildInfo.startDate.slice(0, 8)
      const buildDate = new Date(`${dateStr.slice(0, 4)}-${dateStr.slice(4, 6)}-${dateStr.slice(6, 8)}`)
      const cutoffDate = new Date("2025-10-19")
      model.targetJpsCompile = branch === "master" && buildDate <= cutoffDate
    } else model.targetJpsCompile = false
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Failed to fetch TC info"
  } finally {
    loading.value = false
  }
})
</script>
