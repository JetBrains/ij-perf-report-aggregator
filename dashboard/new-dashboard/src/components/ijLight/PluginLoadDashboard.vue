<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="ijLight_plugin_load_dashboard"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="All plugins into empty JetBrains Light"
        :measure="['enablePluginsModality', 'loadPluginsModalProgress', 'enablePlugins']"
        :projects="['jet-brains-light-plugin-load-performance-test/batch-into-empty-light']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Single plugin into loaded JetBrains Light: EDT modality (enablePluginsModality)"
        measure="enablePluginsModality"
        :projects="singlePluginProjects"
        :aliases="singlePluginAliases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Single plugin into loaded JetBrains Light: modal progress (loadPluginsModalProgress)"
        measure="loadPluginsModalProgress"
        :projects="singlePluginProjects"
        :aliases="singlePluginAliases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Single plugin into loaded JetBrains Light: whole apply (enablePlugins)"
        measure="enablePlugins"
        :projects="singlePluginProjects"
        :aliases="singlePluginAliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const singlePluginProjects = [
  "jet-brains-light-plugin-load-performance-test/java-last",
  "jet-brains-light-plugin-load-performance-test/kotlin-last",
  "jet-brains-light-plugin-load-performance-test/python-last",
  "jet-brains-light-plugin-load-performance-test/database-last",
]
const singlePluginAliases = ["Java", "Kotlin", "Python", "Database"]
</script>
