<template>
  <Message
    v-if="warnings?.length > 0"
    severity="warn"
  >
    <li
      v-for="warning in warnings as Array<Accident>"
      :key="warning.id"
    >
      Known degradation in <b>{{ warning.affectedTest }}</b> <span v-if="warning.buildNumber!=''">, in build <b>{{ warning.buildNumber }}</b></span>.
      Reason: <b>{{ warning.reason }}</b>.
    </li>
  </Message>
</template>

<script setup lang="ts">
import { combineLatest, of } from "rxjs"
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { Accident, getWarningFromMetaDb } from "shared/src/meta"
import { ref } from "vue"


const props = defineProps<{
  table: string
  branchConfigurator: DimensionConfigurator
  scenarioConfigurator?: DimensionConfigurator
  timeRangeConfigurator: TimeRangeConfigurator
}>()

const warnings = ref<Array<Accident>>()

const selected = props.scenarioConfigurator == null ? null : refToObservable(props.scenarioConfigurator.selected)
combineLatest([refToObservable(props.branchConfigurator.selected),
  refToObservable(props.timeRangeConfigurator.value),
  selected || of(null),
]).subscribe(data => {
  getWarningFromMetaDb(warnings, props.branchConfigurator.selected.value, data[2], props.table, props.timeRangeConfigurator.value.value as TimeRange)
})
</script>

