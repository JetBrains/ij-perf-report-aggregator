<template>
  <Message
    v-if="warnings?.length > 0"
    severity="warn"
  >
    <li
      v-for="warning in warnings as Array<Accident>"
      :key="warning.id"
    >
      Known degradation in <b>{{ warning.affectedTest }}</b>, in build <b>{{ warning.buildNumber }}</b>. Reason: <b>{{ warning.reason }}</b>.
    </li>
  </Message>
</template>

<script setup lang="ts">
import { combineLatest, of } from "rxjs"
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { encodeRison } from "shared/src/rison"
import { ref } from "vue"


const props = defineProps<{
  table: string
  branchConfigurator: DimensionConfigurator
  scenarioConfigurator?: DimensionConfigurator
  timeRangeConfigurator: TimeRangeConfigurator
}>()

class Accident {
  constructor(readonly affectedTest: string, readonly date: string, readonly reason: string, readonly id: number, readonly buildNumber: string) {
  }
}

const warnings = ref<Array<Accident>>()

const selected = props.scenarioConfigurator == null ? null : refToObservable(props.scenarioConfigurator.selected)
combineLatest([refToObservable(props.branchConfigurator.selected),
  refToObservable(props.timeRangeConfigurator.value),
  selected || of(null),
]).subscribe(data => {
  getWarningFromMetaDb(props.branchConfigurator.selected.value, data[2], props.table)
})

function isDateInsideRange(dateOfAccident: Date, interval: TimeRange): boolean {
  const currentDate = new Date()
  const day = 24 * 60 * 60 * 1000
  const intervalMapping = {
    "1M": day * 30,
    "3M": day * 30 * 3,
    "1y": day * 365,
    "all": day * 365,
  }
  const selectedDate = new Date()
  selectedDate.setTime(Date.now() - intervalMapping[interval])
  return dateOfAccident >= selectedDate && dateOfAccident <= currentDate
}

function getWarningFromMetaDb(branches: Array<string> | string | null, tests: Array<string> | string | null, table: string) {
  if (branches == null) {
    return
  }
  if (!Array.isArray(branches)) {
    branches = [branches]
  }
  if (tests != null && !Array.isArray(tests)) {
    tests = [tests]
  }
  const url = ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/"
  warnings.value = []
  const data = tests == null ? {table, branches} : {table, branches, tests}
  fetch(url + encodeRison(data))
    .then(response => response.json())
    .then((data: Array<Accident>) => {
      for (const datum of data) {
        if (isDateInsideRange(new Date(datum.date), props.timeRangeConfigurator.value.value as TimeRange)) {
          warnings.value?.push(datum)
        }
      }
    })
    .catch(error => console.error(error))
}

</script>

