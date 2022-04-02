<!-- Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <div class="grid grid-cols-2 gap-4">
    <Textarea
      v-model="inputData"
      :rows="10"
      placeholder="Enter the IntelliJ Platform start-up timeline..."
      :cols="80"
      class="!font-mono"
    />
    <div class="grid grid-cols-1 gap-2 place-content-start justify-items-start w-fit">
      <Button
        class="p-button-sm"
        :loading="isFetching"
        @click="getFromRunningInstance"
      >
        Get from instance
      </Button>
      <InputNumber
        v-model="portNumber"
        :show-buttons="true"
        class="p-inputtext-sm"
        :min="1024"
        :max="65535"
        size="5"
        :format="false"
      />
      <Button
        :loading="isFetchingDev"
        class="col-span-2 p-button-sm w-full"
        @click="getFromRunningDevInstance"
      >
        Get from instance on port 63343
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { serverUrlKey } from "shared/src/injectionKeys"
import { DebouncedTask, TaskHandle } from "shared/src/util/debounce"
import { loadJson } from "shared/src/util/httpUtil"
import { inject, ref, watch } from "vue"
import { RouteLocationNormalizedLoaded, useRoute } from "vue-router"
import { recentlyUsedIdePort, reportData } from "./state"

const route = useRoute()

// we can set this flag using reference to button, but "[Vue warn]: Avoid mutating a prop directly...",
// so, it seems that data property it is the only recommended way
const isFetching = ref(false)
const isFetchingDev = ref(false)

// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
const serverUrl = inject(serverUrlKey)!

let lastReportUrl = ""

const loadReportDebounced = new DebouncedTask(function (taskHandle: TaskHandle): Promise<unknown> {
  if (lastReportUrl.length === 0) {
    return Promise.resolve()
  }

  // localhost blocked by Firefox, but 127.0.0.1 not.
  // Google Chrome correctly resolves localhost, but Firefox doesn't.
  return loadJson(lastReportUrl, null, taskHandle, data => {
    isFetchingDev.value = false
    isFetching.value = false

    if (data == null) {
      return
    }

    reportData.value = JSON.stringify(data, null, 2)
  })
})

function loadReportUrlIfSpecified(location: RouteLocationNormalizedLoaded) {
  const reportUrl = location.query["reportUrl"]
  if (reportUrl != null && reportUrl.length > 0 && lastReportUrl !== reportUrl) {
    isFetching.value = true
    lastReportUrl = `${serverUrl.value}${(reportUrl as string)}`
    loadReportDebounced.execute()
  }
}

loadReportUrlIfSpecified(route)
watch(route, location => {
  loadReportUrlIfSpecified(location)
})

const inputData = reportData
const portNumber = recentlyUsedIdePort
function getFromRunningInstance() {
  isFetching.value = true
  lastReportUrl = getIdeaReportUrl(recentlyUsedIdePort.value)
  loadReportDebounced.execute()
}
function getFromRunningDevInstance() {
  isFetchingDev.value = true
  lastReportUrl = getIdeaReportUrl(63343)
  loadReportDebounced.execute()
}

function getIdeaReportUrl(port: number) {
  return `http://127.0.0.1:${port}/api/startUpMeasurement`
}
</script>