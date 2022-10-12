<template>
  <TabView
    class="text-sm"
    :active-index="0"
  >
    <TabPanel
      header="Performance dashboard"
    >
      <PerformanceDashboard />
    </TabPanel>
    <TabPanel
      header="Performance tests"
      :disabled="true"
    >
      Content II
    </TabPanel>
    <TabPanel
      header="Compare branches"
      :disabled="true"
    >
      Content III
    </TabPanel>
  </TabView>
</template>

<script setup lang="ts">
import TabPanel from "primevue/tabpanel"
import TabView from "primevue/tabview"
import PerformanceDashboard from "./inteliJ/PerformanceDashboard.vue"

</script>