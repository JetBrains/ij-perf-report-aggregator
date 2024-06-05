<template>
  <DashboardPage
    db-name="mlEvaluation"
    table="report"
    persistent-id="aiaDashboard"
    initial-machine="Linux EC2 C5ad.xlarge (4 vCPU AMD EPYC 7002, 8 GB)"
    :charts="charts"
    :with-installer="false"
  >
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
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const projectDeclaration = [
  "name-suggestion_css_NativeLLaMA",
  "name-suggestion_css_JBAICloudWithoutControl",
  "name-suggestion_css_JBAICloudControl",
  "name-suggestion_csharp_NativeLLaMA",
  "name-suggestion_csharp_JBAICloudWithoutControl",
  "name-suggestion_csharp_JBAICloudControl",
  "name-suggestion_cpp_NativeLLaMA",
  "name-suggestion_cpp_JBAICloudWithoutControl",
  "name-suggestion_cpp_JBAICloudControl",
  "name-suggestion_go_NativeLLaMA",
  "name-suggestion_go_JBAICloudWithoutControl",
  "name-suggestion_go_JBAICloudControl",
  "name-suggestion_html_NativeLLaMA",
  "name-suggestion_html_JBAICloudWithoutControl",
  "name-suggestion_html_JBAICloudControl",
  "name-suggestion_javascript_NativeLLaMA",
  "name-suggestion_javascript_JBAICloudWithoutControl",
  "name-suggestion_javascript_JBAICloudControl",
  "name-suggestion_java_NativeLLaMA",
  "name-suggestion_java_JBAICloudWithoutControl",
  "name-suggestion_java_JBAICloudControl",
  "name-suggestion_kotlin_NativeLLaMA",
  "name-suggestion_kotlin_JBAICloudWithoutControl",
  "name-suggestion_kotlin_JBAICloudControl",
  "name-suggestion_php_NativeLLaMA",
  "name-suggestion_php_JBAICloudWithoutControl",
  "name-suggestion_php_JBAICloudControl",
  "name-suggestion_python_NativeLLaMA",
  "name-suggestion_python_JBAICloudWithoutControl",
  "name-suggestion_python_JBAICloudControl",
  "name-suggestion_ruby_NativeLLaMA",
  "name-suggestion_ruby_JBAICloudWithoutControl",
  "name-suggestion_ruby_JBAICloudControl",
  "name-suggestion_rust_NativeLLaMA",
  "name-suggestion_rust_JBAICloudWithoutControl",
  "name-suggestion_rust_JBAICloudControl",
  "name-suggestion_scala_NativeLLaMA",
  "name-suggestion_scala_JBAICloudWithoutControl",
  "name-suggestion_scala_JBAICloudControl",
  "name-suggestion_terraform_NativeLLaMA",
  "name-suggestion_terraform_JBAICloudWithoutControl",
  "name-suggestion_terraform_JBAICloudControl",
]

const chartsDeclaration: ChartDefinition[] = projectDeclaration.map((project) => {
  return {
    labels: [project],
    measures: ["EditSimilarity"],
    projects: [project],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
