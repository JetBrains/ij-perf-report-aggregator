<template>
  <Message
    v-if="warnings != null && warnings.length > 0"
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
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { encodeRison } from "shared/src/rison"
import { ref } from "vue"


const props = defineProps<{
  table: string
  branchConfigurator: DimensionConfigurator
  timeRangeConfigurator: TimeRangeConfigurator
}>()

class Accident {
  constructor(readonly affectedTest: string, readonly date: string, readonly reason: string, readonly id: number, readonly buildNumber: string) {}
}

const warnings = ref<Array<Accident>>()

refToObservable(props.branchConfigurator.selected).subscribe(data => {
  getWarningFromMetaDb(data, props.table)
})
refToObservable(props.timeRangeConfigurator.value).subscribe(data => {
  getWarningFromMetaDb(props.branchConfigurator.selected.value, props.table)
})

function isDateInsideRange(dateOfAccident: Date, interval: TimeRange): boolean {
  const currentDate = new Date()
  let subtractInterval: number  //5 years ago
  if (interval == "1M") {
    subtractInterval = 24 * 60 * 60 * 1000 * 30
  }
  else if (interval == "3M") {
    subtractInterval = 24 * 60 * 60 * 1000 * 30 * 3
  }
  else if (interval == "1y") {
    subtractInterval = 24 * 60 * 60 * 1000 * 30 * 12
  }
  else {
    subtractInterval = 24 * 60 * 60 * 1000 * 30 * 12 * 5
  }
  const selectedDate = new Date()
  selectedDate.setTime(Date.now() - subtractInterval)
  return dateOfAccident >= selectedDate && dateOfAccident <= currentDate
}

function getWarningFromMetaDb(branches: Array<string> | string | null, table: string) {
  if (branches == null) {
    return
  }
  if (!Array.isArray(branches)) {
    branches = [branches]
  }
  const url = ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/"
  warnings.value = []
  for (const branch of branches) {
    const data = {branch, table}
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
}

</script>

