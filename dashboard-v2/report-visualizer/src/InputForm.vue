<!-- Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <el-row :gutter="16">
    <el-col :span="10">
      <el-input
        v-model="inputData"
        type="textarea"
        :rows="30"
        placeholder="Enter the IntelliJ Platform start-up timeline..."
      />
    </el-col>
    <el-col :span="14">
      <el-form
        :inline="true"
        size="small"
      >
        <el-form-item>
          <el-button
            :loading="isFetching"
            @click="getFromRunningInstance"
          >
            Get from running instance
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-input-number
            v-model="portNumber"
            :min="1024"
            :max="65535"
          />
        </el-form-item>
      </el-form>
      <el-form
        :inline="true"
        size="small"
      >
        <el-form-item>
          <el-button
            :loading="isFetchingDev"
            @click="getFromRunningDevInstance"
          >
            Get from running instance on port 63343
          </el-button>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { DebouncedTask, TaskHandle } from "shared/src/util/debounce"
import { loadJson } from "shared/src/util/httpUtil"
import { defineComponent, ref, watch } from "vue"
import { RouteLocationNormalizedLoaded, useRoute } from "vue-router"
import { recentlyUsedIdePort, reportData } from "./state/state"

export default defineComponent({
  name: "InputForm",

  setup() {
    const route = useRoute()

    // we can set this flag using reference to button, but "[Vue warn]: Avoid mutating a prop directly...",
    // so, it seems that data property it is the only recommended way
    const isFetching = ref(false)
    const isFetchingDev = ref(false)

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
        lastReportUrl = reportUrl as string
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
      }
    }
  }
})

function getIdeaReportUrl(port: number) {
  return `http://127.0.0.1:${port}/api/startUpMeasurement`
}
</script>
