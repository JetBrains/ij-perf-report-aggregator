<template>
  <DashboardPage
    db-name="perfintDev"
    table="phpstorm"
    persistent-id="phpstorm_memory_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <Divider title="Indexing: heap after GC (live set) / before GC (as the user sees it)" />

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Heap after GC, Mb"
          :measure="AFTER_GC"
          :projects="indexingProjects"
          value-unit="counter"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Heap before GC, Mb"
          :measure="BEFORE_GC"
          :projects="indexingProjects"
          value-unit="counter"
        />
      </div>
    </section>

    <Divider title="Batch inspections: heap after GC (live set) / before GC (as the user sees it)" />

    <section
      v-for="(group, index) in inspectionGroups"
      :key="index"
      class="flex gap-x-6"
    >
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Heap after GC, Mb"
          :measure="AFTER_GC"
          :projects="group"
          value-unit="counter"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Heap before GC, Mb"
          :measure="BEFORE_GC"
          :projects="group"
          value-unit="counter"
        />
      </div>
    </section>

    <Divider title="Per project: before GC / after GC" />

    <section
      v-for="project in allProjects"
      :key="project"
    >
      <GroupProjectsChart
        :label="project"
        :measure="[BEFORE_GC, AFTER_GC]"
        :projects="[project]"
        value-unit="counter"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

const AFTER_GC = "JVM.heapUsageMb/afterIndexing"
const BEFORE_GC = `${AFTER_GC}/beforeGC`

const indexingProjects = ["Coolify/indexing", "Bagisto/indexing", "Appwrite/indexing"]

const inspectionGroups = [
  [
    "drupal8-master-with-plugin/inspection",
    "shopware/inspection",
    "b2c-demo-shop/inspection",
    "magento/inspection",
    "wordpress/inspection",
    "laravel-io/inspection"
  ],
  ["mediawiki/inspection", "php-cs-fixer/inspection", "proxyManager/inspection"],
  ["akaunting/inspection", "aggregateStitcher/inspection", "prestaShop/inspection", "kunstmaanBundlesCMS/inspection"],
  ["Coolify/inspection", "Bagisto/inspection", "Appwrite/inspection"]
]

const allProjects = [...indexingProjects, ...inspectionGroups.flat()]
</script>
