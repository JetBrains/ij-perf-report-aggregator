<!-- Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <div class="columns-2">
    <Textarea
      v-model="inputData"
      :rows="10"
      placeholder="Enter the IntelliJ Platform start-up timeline..."
      :cols="80"
      class="w-full"
    />
    <div class="w-full">
      <div class="columns-1">
        <div class="w-full">
          <Button
            :loading="isFetching"
            @click="getFromRunningInstance"
          >
            Get from running instance
          </Button>
          <InputNumber
            v-model="portNumber"
            :min="1024"
            :max="65535"
            :format="false"
          />
        </div>
        <Button
          :loading="isFetchingDev"
          @click="getFromRunningDevInstance"
        >
          Get from running instance on port 63343
        </Button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { serverUrlKey } from "shared/src/injectionKeys"
import { DebouncedTask, TaskHandle } from "shared/src/util/debounce"
import { loadJson } from "shared/src/util/httpUtil"
import { defineComponent, inject, ref, watch } from "vue"
import { RouteLocationNormalizedLoaded, useRoute } from "vue-router"
import { recentlyUsedIdePort, reportData } from "./state"

export default defineComponent({
  name: "InputForm",

  setup() {
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

    return {
      inputData: reportData,
      isFetching,
      isFetchingDev,
      portNumber: recentlyUsedIdePort,
      getFromRunningInstance() {
        isFetching.value = true
        lastReportUrl = getIdeaReportUrl(recentlyUsedIdePort.value)
        loadReportDebounced.execute()
      },
      getFromRunningDevInstance() {
        isFetchingDev.value = true
        lastReportUrl = getIdeaReportUrl(63343)
        loadReportDebounced.execute()
      },
    }
  },
})

function getIdeaReportUrl(port: number) {
  return `http://127.0.0.1:${port}/api/startUpMeasurement`
}
</script>

<style scoped>
/*noinspection CssUnusedSymbol*/
.p-button, .p-inputnumber, .p-inputtextarea {
  margin: 10px;
}
</style>