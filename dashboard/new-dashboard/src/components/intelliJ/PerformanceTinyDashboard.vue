<template>
  <DashboardPage
    v-slot="{averagesConfigurators}"
    db-name="perfint"
    table="idea"
    persistent-id="idea_tiny_dashboard"
    initial-machine="Linux EC2 C6id.large (2 vCPU Xeon, 4 GB)"
    :charts="charts"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#JAVA'"
          :title="'Indexing Java (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#Kotlin'"
          :title="'Indexing Kotlin (kB/s)'"
          :chart-color="'#9B51E0'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :is-like="true"
          :title="'Completion'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="[...averagesConfigurators, typingOnlyConfigurator]"
          :aggregated-measure="'test#average_awt_delay'"
          :title="'UI responsiveness during typing'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { DataQuery, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import AggregationChart from "../charts/AggregationChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [{
  labels: ["Indexing (Big projects)", "Scanning (Big projects)"],
  measures: ["indexing", "scanning"],
  projects: ["grails/indexing", "kotlin_coroutines/indexing", "kotlin/indexing", "spring_boot/indexing"],
},  {
  labels: ["Indexing", "Scanning"],
  measures: ["indexing", "scanning"],
  projects: ["empty_project/indexing", "java/indexing", "spring_boot_maven/indexing", "kotlin_petclinic/indexing"],
},  {
  labels: ["Rebuild"],
  measures: ["build_compilation_duration"],
  projects: ["grails/rebuild", "java/rebuild", "spring_boot/rebuild"],
},  {
  labels: ["Local Inspection", "First Code Analysis"],
  measures: ["localInspections", "firstCodeAnalysis"],
  projects: ["kotlin/localInspection",
    "kotlin_coroutines/localInspection"],
}, {
  labels: ["Completion"],
  measures: ["completion"],
  projects: ["grails/completion/groovy_file", "grails/completion/java_file"],
},  {
  labels: ["Show Intentions (average awt delay)"],
  measures: ["test#average_awt_delay"],
  projects: ["grails/showIntentions/Find cause", "kotlin/showIntention/Import", "spring_boot/showIntentions"],
}, {
  labels: ["Highlight"],
  measures: ["highlighting"],
  projects: ["kotlin/highlight", "kotlin_coroutines/highlight"],
}]

const charts = combineCharts(chartsDeclaration)

const typingOnlyConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({f: "project", v: "%typing", o: "like"})
    return true
  },
  createObservable() {
    return null
  },
}
</script>