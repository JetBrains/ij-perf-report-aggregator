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
    <div class="grid gap-2 place-content-start">
      <Button
        :loading="isFetching"
        icon="pi pi-download"
        label="Get from instance"
        @click="getFromRunningInstance"
      />
      <InputNumber
        v-model="portNumber"
        :show-buttons="true"
        :min="1024"
        :max="65535"
        size="5"
        :format="false"
      />
      <Button
        :loading="isFetchingDev"
        class="col-span-2"
        icon="pi pi-download"
        label="Get from instance on port 63343"
        @click="getFromRunningDevInstance"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { debounceTime, finalize, of, Subject, switchMap } from "rxjs"
import { fromFetchWithRetryAndErrorHandling } from "../shared/configurators/rxjs"
import { shallowRef, watch } from "vue"
import { LocationQuery, useRoute } from "vue-router"
import { recentlyUsedIdePort, reportData } from "./state"

const route = useRoute()

// we can set this flag using reference to button, but "[Vue warn]: Avoid mutating a prop directly...",
// so, it seems that data property it is the only recommended way
const isFetching = shallowRef(false)
const isFetchingDev = shallowRef(false)

const subject = new Subject<{url: string; isDev: boolean}>()
// distinctUntilChanged is not used here because report maybe loaded from the same instance multiple times
subject
  .pipe(
    debounceTime(100),
    switchMap(({url, isDev}) => {
      if (url.length === 0) {
        return of(null)
      }
      const statusRef = isDev ? isFetchingDev : isFetching
      statusRef.value = true
      return fromFetchWithRetryAndErrorHandling<string>(url, {
        summary: "Cannot connect to IDE",
        detail: "Please check that port is correct."
      })
        .pipe(
          finalize(() => {
            statusRef.value = false
          }),
        )
    }),
  )
  .subscribe(data => {
    if (data != null) {
      reportData.value = JSON.stringify(data, null, 2)
    }
  })

function loadReportUrlIfSpecified(query: LocationQuery) {
  const reportUrl = query["reportUrl"]
  if (reportUrl != null && reportUrl.length > 0) {
    subject.next({url: reportUrl as string, isDev: false})
  }
}

loadReportUrlIfSpecified(route.query)
watch(() => route.query, loadReportUrlIfSpecified)

const inputData = reportData
const portNumber = recentlyUsedIdePort

function getFromRunningInstance() {
  subject.next({url: getIdeaReportUrl(recentlyUsedIdePort.value), isDev: false})
}
function getFromRunningDevInstance() {
  subject.next({url: getIdeaReportUrl(63343), isDev: true})
}

function getIdeaReportUrl(port: number) {
  return `http://127.0.0.1:${port}/api/startUpMeasurement`
}
</script>