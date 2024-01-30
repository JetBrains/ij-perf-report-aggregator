<template>
  <Dialog
    v-model:visible="showDialog"
    modal
    header="Report Event"
    :style="{ width: '30vw' }"
  >
    <div class="flex items-center space-x-4 mb-4">
      <Dropdown
        v-model="accidentType"
        placeholder="Event Type"
        :options="getAccidentTypes()"
      >
        <template #value="{ value }">
          <div class="group inline-flex justify-center text-sm font-medium text-gray-700 hover:text-gray-900">
            {{ value }}
            <ChevronDownIcon
              class="-mr-1 ml-1 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500"
              aria-hidden="true"
            />
          </div>
        </template>
        <template #dropdownicon>
          <!-- empty element to avoid ignoring override of slot -->
          <span />
        </template>
      </Dropdown>
      <span class="p-float-label flex-grow">
        <InputText
          id="reason"
          v-model="reason"
          class="w-full"
        />
        <label
          class="text-sm"
          for="reason"
          >Reason</label
        >
      </span>
    </div>
    <div
      v-if="props.data.series.length == 1"
      class="flex items-center mb-4"
    >
      <InputSwitch
        v-model="reportMetricOnly"
        input-id="reportMetricOnly"
      />
      <label
        for="reportMetricOnly"
        class="text-sm ml-2"
      >
        Report metric only
      </label>
    </div>
    <!-- Footer buttons -->
    <template #footer>
      <div class="flex justify-end space-x-2">
        <Button
          label="Cancel"
          icon="pi pi-times"
          text
          @click="showDialog = false"
        />
        <Button
          label="Report"
          icon="pi pi-check"
          autofocus
          @click="reportRegression"
        />
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { useStorage } from "@vueuse/core/index"
import { ref } from "vue"
import { AccidentKind, AccidentsConfigurator } from "../../../configurators/AccidentsConfigurator"
import { InfoData } from "./InfoSidebar"

const props = defineProps<{
  data: InfoData
  accidentsConfigurator: AccidentsConfigurator
}>()

const showDialog = defineModel<boolean>()

const reportMetricOnly = useStorage("reportMetricOnly", false)
const accidentType = ref<string>("Regression")
const reason = ref("")

function reportRegression() {
  showDialog.value = false
  const value = props.data

  const reportOnlyMetric = reportMetricOnly.value && value.series.length == 1
  props.accidentsConfigurator.writeAccidentToMetaDb(
    value.date,
    value.projectName + (reportOnlyMetric ? "/" + value.series[0].metricName : ""),
    reason.value,
    value.build ?? value.buildId.toString(),
    accidentType.value
  )
}

function getAccidentTypes(): string[] {
  const values = Object.values(AccidentKind)
  //don't report Inferred type manually
  values.splice(values.indexOf(AccidentKind.InferredRegression), 1)
  values.splice(values.indexOf(AccidentKind.InferredImprovement), 1)
  return values
}
</script>
