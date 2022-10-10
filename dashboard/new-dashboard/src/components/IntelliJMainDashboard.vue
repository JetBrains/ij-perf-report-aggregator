<template>
  <TabView
    class="text-sm"
    :active-index="0"
    @tab-click="handleClick"
  >
    <TabPanel
      header="Performance dashboard"
    >
      <PerformanceDashboard />
    </TabPanel>
    <TabPanel
      header="Performance tests"
    >
      Content II
    </TabPanel>
    <TabPanel
      header="Compare branches"
    >
      Content III
    </TabPanel>
  </TabView>
</template>

<script setup lang="ts">
import TabPanel from "primevue/tabpanel"
import TabView from "primevue/tabview"
import PerformanceDashboard from "./inteliJ/PerformanceDashboard.vue"
const path = "/new/ij/tests"
</script>