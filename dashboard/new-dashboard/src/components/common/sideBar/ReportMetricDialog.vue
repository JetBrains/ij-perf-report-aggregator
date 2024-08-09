<template>
  <Dialog
    v-model:visible="showDialog"
    modal
    header="Report Event"
    :style="{ width: '30vw' }"
  >
    <div class="flex items-center space-x-4 mb-4 mt-4">
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
      <FloatLabel class="w-full">
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
      </FloatLabel>
    </div>
    <FloatLabel
      v-if="accidentType == AccidentKind.Exception"
      class="w-full"
    >
      <Textarea
        id="stacktrace"
        v-model="stacktrace"
        class="w-full"
      />
      <label
        class="text-sm"
        for="stacktrace"
        >Stacktrace</label
      >
    </FloatLabel>
    <div
      v-if="props.data?.series.length == 1"
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
        Report only metric <code>{{ props.data.series[0].metricName }}</code>
      </label>
    </div>
    <div
      v-if="props.data?.series.length == 1"
      class="flex items-center mb-4"
    >
      <InputSwitch
        v-model="reportAllInBuild"
        input-id="reportAllInBuild"
      />
      <label
        for="reportAllInBuild"
        class="text-sm ml-2"
      >
        Report all tests in build <code>{{ build }}</code>
      </label>
    </div>
    <RelatedAccidents
      :data="props.data"
      :accidents-configurator="props.accidentsConfigurator"
      :in-dialog="true"
      @copy-accident="copy"
    />
    <!-- Footer buttons -->
    <template #footer>
      <div v-if="accidentToEdit == null">
        <div class="flex justify-end space-x-2">
          <Button
            label="Cancel"
            icon="pi pi-times"
            severity="secondary"
            @click="showDialog = false"
          />
          <Button
            label="Report"
            icon="pi pi-check"
            autofocus
            @click="reportRegression"
          />
        </div>
      </div>
      <div v-else>
        <div class="flex justify-end space-x-2">
          <Button
            label="Cancel"
            icon="pi pi-times"
            severity="secondary"
            @click="showDialog = false"
          />
          <Button
            label="Delete"
            icon="pi pi-trash"
            severity="danger"
            @click="deleteRegression"
          />
          <Button
            label="Update"
            icon="pi pi-pencil"
            autofocus
            @click="updateRegression"
          />
        </div>
      </div>
    </template>
  </Dialog>
</template>
<script setup lang="ts">
import { ChevronDownIcon } from "@heroicons/vue/20/solid/index"
import { useStorage } from "@vueuse/core/index"
import { computed, ref, watch } from "vue"
import { Accident, AccidentKind, AccidentsConfigurator } from "../../../configurators/AccidentsConfigurator"
import { InfoData } from "./InfoSidebar"
import RelatedAccidents from "./RelatedAccidents.vue"

const props = defineProps<{
  data: InfoData | null
  accidentsConfigurator: AccidentsConfigurator | null
  accidentToEdit: Accident | null
}>()

const showDialog = defineModel<boolean>()

const reportMetricOnly = useStorage("reportMetricOnly", false)
const reportAllInBuild = useStorage("reportAllInBuild", false)
const accidentType = ref(props.accidentToEdit?.kind ?? "Regression")
watch(
  () => props.accidentToEdit,
  (newVal) => {
    accidentType.value = newVal?.kind ?? "Regression"
    reason.value = newVal?.reason ?? ""
  }
)

const reason = ref(props.accidentToEdit?.reason ?? "")
const stacktrace = ref(props.accidentToEdit?.stacktrace ?? "")

const build = computed(() => props.data?.build ?? props.data?.buildId.toString())

function reportRegression() {
  showDialog.value = false
  const value = props.data
  if (value != null && build.value != null) {
    const metricName = value.series[0].metricName
    const reportOnlyMetric = reportMetricOnly.value && value.series.length == 1 && metricName != undefined
    props.accidentsConfigurator?.writeAccidentToMetaDb(
      value.date,
      reportAllInBuild.value ? "" : value.projectName + (reportOnlyMetric ? "/" + metricName : ""),
      reason.value,
      build.value,
      accidentType.value,
      stacktrace.value
    )
  }
}

async function deleteRegression() {
  showDialog.value = false
  if (props.accidentToEdit != null) {
    await props.accidentsConfigurator?.removeAccidentFromMetaDb(props.accidentToEdit.id)
  }
}

function updateRegression() {
  deleteRegression()
    .then(() => {
      reportRegression()
    })
    .catch(() => {
      console.error("Failed to delete accident")
    })
}

watch(reportMetricOnly, (newValue) => {
  if (newValue) {
    reportAllInBuild.value = false
  }
})

watch(reportAllInBuild, (newValue) => {
  if (newValue) {
    reportMetricOnly.value = false
  }
})

function copy(accident: { kind: string; reason: string }) {
  reason.value = accident.reason
  accidentType.value = accident.kind
}

function getAccidentTypes(): string[] {
  const values = Object.values(AccidentKind)
  //don't report Inferred type manually
  values.splice(values.indexOf(AccidentKind.InferredRegression), 1)
  values.splice(values.indexOf(AccidentKind.InferredImprovement), 1)
  return values
}
</script>
