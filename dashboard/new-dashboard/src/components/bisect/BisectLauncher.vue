<template>
  <div class="flex justify-center">
    <div class="mt-3 border-1 w-1/2">
      <Card>
        <template #title>
          <div class="text-darker">Start Bisect</div>
        </template>
        <template #content>
          <div class="flex flex-col space-y-8 mb-4 mt-6">
            <div class="flex">
              <FloatLabel>
                <InputText
                  id="changes"
                  v-model="model.firstCommit"
                />
                <label for="changes">First commit</label>
              </FloatLabel>
              <FloatLabel>
                <InputText
                  id="changes"
                  v-model="model.lastCommit"
                />
                <label for="changes">Last commit</label>
              </FloatLabel>
            </div>
            <Accordion :value="shouldExpandAccordion ? 0 : -1">
              <AccordionPanel :value="0">
                <AccordionHeader>Additional parameters</AccordionHeader>
                <AccordionContent>
                  <div class="flex flex-col space-y-8 mb-4 mt-4">
                    <FloatLabel>
                      <InputText
                        id="test"
                        v-model="model.test"
                      />
                      <label for="test">Test name</label>
                    </FloatLabel>
                    <FloatLabel>
                      <InputText
                        id="branch"
                        v-model="model.branch"
                      />
                      <label for="metric">Branch</label>
                    </FloatLabel>
                    <FloatLabel>
                      <InputText
                        id="buildType"
                        v-model="model.buildType"
                      />
                      <label for="buildType">Build type</label>
                    </FloatLabel>

                    <FloatLabel>
                      <InputText
                        id="className"
                        v-model="model.className"
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
        </template>
        <template #footer>
          <div class="flex justify-end">
            <Button
              label="Start"
              icon="pi pi-play"
              autofocus
              :loading="loading"
              :disabled="anyFieldIsEmpty"
              @click="startBisect"
            />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed, reactive, ref } from "vue"
import { BisectClient } from "../common/sideBar/BisectClient"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"

const props = withDefaults(
  defineProps<{
    firstCommit?: string
    lastCommit?: string
    test?: string
    branch?: string
    buildType?: string
    className?: string
  }>(),
  {
    firstCommit: "",
    lastCommit: "",
    test: "",
    branch: "",
    buildType: "",
    className: "",
  }
)

const model = reactive({
  firstCommit: props.firstCommit,
  lastCommit: props.lastCommit,
  test: props.test,
  branch: props.branch,
  buildType: props.buildType,
  className: props.className,
})

const shouldExpandAccordion = computed(() => props.test == "" || props.branch == "" || props.buildType == "" || props.className == "")
const anyFieldIsEmpty = computed(
  () =>
    model.firstCommit.trim() == "" ||
    model.lastCommit.trim() == "" ||
    model.test.trim() == "" ||
    model.branch.trim() == "" ||
    model.buildType.trim() == "" ||
    model.className.trim() == ""
)

const error = ref<string | null>(null)
const loading = ref(false)
const serverConfigurator = new ServerWithCompressConfigurator("", "")
const bisectClient = new BisectClient(serverConfigurator)

async function startBisect() {
  error.value = null
  loading.value = true
  try {
    const weburl = await bisectClient.sendBisectRequest({
      changes: model.firstCommit + "^.." + model.lastCommit,
      test: model.test,
      branch: model.branch,
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
</script>
