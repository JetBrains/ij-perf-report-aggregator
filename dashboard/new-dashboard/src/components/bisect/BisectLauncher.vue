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
                  <div class="grid grid-cols-4 mt-6">
                    <div class="col-span-4">
                      <FloatLabel class="w-full">
                        <InputText
                          id="className"
                          v-model="model.className"
                          class="w-full"
                        />
                        <label for="className">Class name</label>
                      </FloatLabel>
                    </div>
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

const props = withDefaults(
  defineProps<{
    buildId?: string
    errorMessage?: string
    className?: string
  }>(),
  {
    buildId: "",
    errorMessage: "",
    className: "",
  }
)

const serverConfigurator = new ServerWithCompressConfigurator("", "")
const bisectClient = new BisectClient(serverConfigurator)

const model = reactive({
  errorMessage: props.errorMessage,
  firstCommit: "",
  lastCommit: "",
  buildType: "",
  buildId: props.buildId,
  className: props.className,
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
      mode: isCommitMode.value ? "commits" : "build",
      errorMessage: model.errorMessage,
      buildType: model.buildType,
      className: model.className,
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
    const changes = await bisectClient.fetchTeamCityChanges(props.buildId)
    model.firstCommit = changes.firstCommit
    model.lastCommit = changes.lastCommit
    model.buildType = await bisectClient.fetchBuildType(props.buildId)
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Failed to fetch TC info"
  } finally {
    loading.value = false
  }
})
</script>
