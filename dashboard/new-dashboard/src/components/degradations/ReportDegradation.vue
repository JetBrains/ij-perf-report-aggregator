<template>
  <div class="flex justify-center">
    <div
      v-if="!reported"
      class="mt-3 border-1 w-1/2"
    >
      <Card>
        <template #title> <div class="text-darker">Report event</div></template>
        <template #content>
          <div class="flex items-center space-x-4">
            <div class="text-lg font-bold">Date:</div>
            <div>{{ props.date }}</div>
          </div>
          <div class="flex items-center space-x-4">
            <div class="text-lg font-bold">Build:</div>
            <div>{{ props.build }}</div>
          </div>
          <div class="text-lg font-bold">
            Type:
            <Dropdown
              v-model="accidentType"
              placeholder="Event Type"
              :options="getAccidentTypes()"
            />
          </div>
          <div class="flex flex-col space-x-4">
            <div class="text-lg font-bold">Tests:</div>
            <ul
              v-for="test in tests"
              :key="test"
              class="list-disc list-inside"
            >
              <li class="text-sm">{{ test }}</li>
            </ul>
          </div>

          <div class="text-lg mt-4 font-bold">Reason:</div>
          <div class="flex items-center space-x-4 mt-1">
            <Textarea
              id="reason"
              v-model="reason"
              class="w-full"
            />
          </div>
          <div class="flex space-x-2 mt-20">
            <Button
              class="w-1/5"
              label="Report"
              autofocus
              @click="reportRegression"
            />
          </div>
        </template>
      </Card>
    </div>
    <div v-else>
      <div class="text-4xl justify-center">Events reported 👍</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"
import { AccidentKind, AccidentsConfigurator } from "../../configurators/AccidentsConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"

const props = defineProps<{
  date: string
  tests: string
  build: string
}>()
const tests = computed(() => props.tests.split(","))

const accidentType = ref<string>("Regression")
const reason = ref("")
function getAccidentTypes(): string[] {
  const values = Object.values(AccidentKind)
  values.splice(values.indexOf(AccidentKind.InferredRegression), 1)
  values.splice(values.indexOf(AccidentKind.InferredImprovement), 1)
  return values
}

const serverConfigurator = new ServerWithCompressConfigurator("", "")

class AccidentReporter extends AccidentsConfigurator {
  protected getAccidentUrl(): string {
    return serverConfigurator.serverUrl + "/api/meta"
  }
}

const accidentReporter = new AccidentReporter()
const reported = ref(false)
function reportRegression() {
  for (const test of tests.value) {
    accidentReporter.writeAccidentToMetaDb(props.date, test, reason.value, props.build, accidentType.value)
  }
  reported.value = true
}
</script>
<style #scoped>
.p-inputtext {
  font-size: 1rem;
}
.p-dropdown-panel {
  font-size: 1rem;
}
</style>
```
