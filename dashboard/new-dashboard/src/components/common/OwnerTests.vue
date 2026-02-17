<template>
  <div class="px-7 py-5">
    <PerformanceTests
      v-if="dbName && table"
      :db-name="dbName"
      :table="table"
      :initial-machine="machine"
      :with-installer="withInstaller"
      :machine-group-filter="machine"
    />
  </div>
</template>

<script setup lang="ts">
import { watch } from "vue"
import { useRoute, useRouter } from "vue-router"
import PerformanceTests from "./PerformanceTests.vue"

const route = useRoute()
const router = useRouter()

const dbName = route.query["dbName"] as string
const table = route.query["table"] as string
const machine = (route.query["machine"] as string) ?? null
const withInstaller = dbName === "perfint" || dbName === "ij"

watch(
  () => route.query,
  (newQuery) => {
    if (newQuery["dbName"] !== dbName || newQuery["table"] !== table) {
      void router.replace({
        query: { ...newQuery, dbName, table },
      })
    }
  }
)
</script>
