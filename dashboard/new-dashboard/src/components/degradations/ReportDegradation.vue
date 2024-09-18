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
            <div>{{ date }}</div>
          </div>
          <div class="flex items-center space-x-4">
            <div class="text-lg font-bold">Build:</div>
            <div>{{ build }}</div>
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
              v-for="test in testsArray"
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
import { AccidentKind, AccidentsConfigurator } from "../../configurators/accidents/AccidentsConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"

const { date, tests, build } = defineProps<{
  date: string
  tests: string
  build: string
}>()
const testsArray = computed(() => tests.split(","))

const dateFormatted = computed(() => {
  const parts = date.split("-")
  const day = Number.parseInt(parts[0])
  const month = Number.parseInt(parts[1])
  const year = Number.parseInt(parts[2])
  return year.toString() + "-" + month.toString() + "-" + day.toString()
})

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
    return serverConfigurator.serverUrl + "/api/meta/"
  }
}

const accidentReporter = new AccidentReporter()
const reported = ref(false)
async function reportRegression() {
  for (const test of testsArray.value) {
    await accidentReporter.writeAccidentToMetaDb(dateFormatted.value, test, reason.value, build, accidentType.value)
  }
  reported.value = true
}
</script>
<style #scoped>
.p-select-panel {
  font-size: 1rem;
}
</style>
```
