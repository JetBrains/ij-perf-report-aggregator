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
  "completion_css_NativeLLaMA",
  "completion_css_JBAICloudWithoutControl",
  "completion_css_JBAICloudControl",
  "completion_csharp_NativeLLaMA",
  "completion_csharp_JBAICloudWithoutControl",
  "completion_csharp_JBAICloudControl",
  "completion_cpp_NativeLLaMA",
  "completion_cpp_JBAICloudWithoutControl",
  "completion_cpp_JBAICloudControl",
  "completion_go_NativeLLaMA",
  "completion_go_JBAICloudWithoutControl",
  "completion_go_JBAICloudControl",
  "completion_html_NativeLLaMA",
  "completion_html_JBAICloudWithoutControl",
  "completion_html_JBAICloudControl",
  "completion_javascript_NativeLLaMA",
  "completion_javascript_JBAICloudWithoutControl",
  "completion_javascript_JBAICloudControl",
  "completion_java_NativeLLaMA",
  "completion_java_JBAICloudWithoutControl",
  "completion_java_JBAICloudControl",
  "completion_kotlin_NativeLLaMA",
  "completion_kotlin_JBAICloudWithoutControl",
  "completion_kotlin_JBAICloudControl",
  "completion_php_NativeLLaMA",
  "completion_php_JBAICloudWithoutControl",
  "completion_php_JBAICloudControl",
  "completion_python_NativeLLaMA",
  "completion_python_JBAICloudWithoutControl",
  "completion_python_JBAICloudControl",
  "completion_ruby_NativeLLaMA",
  "completion_ruby_JBAICloudWithoutControl",
  "completion_ruby_JBAICloudControl",
  "completion_rust_NativeLLaMA",
  "completion_rust_JBAICloudWithoutControl",
  "completion_rust_JBAICloudControl",
  "completion_scala_NativeLLaMA",
  "completion_scala_JBAICloudWithoutControl",
  "completion_scala_JBAICloudControl",
  "completion_terraform_NativeLLaMA",
  "completion_terraform_JBAICloudWithoutControl",
  "completion_terraform_JBAICloudControl",
]

const chartsDeclaration: ChartDefinition[] = projectDeclaration.map((project) => {
  return {
    labels: [project],
    measures: ["MatchedRatio"],
    projects: [project],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
