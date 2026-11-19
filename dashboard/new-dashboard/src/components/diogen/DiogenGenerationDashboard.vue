<template>
  <DashboardPage
    db-name="diogen"
    table="report"
    persistent-id="diogen_generation_dashboard"
    :with-installer="false"
    :initial-machine="null"
  >
    <section
      v-for="project in PROJECTS"
      :key="project"
    >
      <div class="flex flex-wrap gap-x-6 gap-y-1 text-sm text-gray-500 dark:text-gray-400 mb-2 mt-4">
        <a
          v-for="link in diogenLinks(project)"
          :key="link.label"
          :href="link.url"
          target="_blank"
          rel="noopener"
          class="flex items-center gap-1.5 underline decoration-dotted hover:no-underline"
        >
          {{ link.label }}
          <i class="pi pi-external-link text-xs" />
        </a>
      </div>
      <GroupProjectsChart
        :label="project"
        :measure="MEASURES"
        :projects="[project]"
        :value-unit="'counter'"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const PROJECTS = ["ijAllCur", "freezeFlat"]
const MEASURES = ["hosts-last-created-issue", "max-hosts-unsorted"]

const DIOGEN_URL = "https://diogen.labs.jb.gg"
const SORT_PARAMS = "sortBy=uniqueUsers&sortDir=desc&pileSortBy=openReports&pileSortDir=desc&view=flat"

function diogenLinks(project: string) {
  return [
    { label: "Unsorted in Diogen", url: `${DIOGEN_URL}/${project}/Unsorted/index.html?${SORT_PARAMS}` },
    { label: "Next candidates for an issue in Diogen", url: `${DIOGEN_URL}/${project}/?${SORT_PARAMS}&readyToReport=1` },
  ]
}
</script>
