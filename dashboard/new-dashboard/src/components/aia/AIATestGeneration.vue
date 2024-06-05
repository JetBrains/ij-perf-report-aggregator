<template>
  <DashboardPage
    db-name="mlEvaluation"
    table="report"
    persistent-id="aiaDashboard"
    :initial-machine="null"
    :charts="charts"
    :with-installer="false"
    :branch="null"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :machines="['Linux EC2 C5ad.xlarge (4 vCPU AMD EPYC 7002, 8 GB)']"
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
  "test-generation_css_NativeLLaMA",
  "test-generation_css_JBAICloudWithoutControl",
  "test-generation_css_JBAICloudControl",
  "test-generation_csharp_NativeLLaMA",
  "test-generation_csharp_JBAICloudWithoutControl",
  "test-generation_csharp_JBAICloudControl",
  "test-generation_cpp_NativeLLaMA",
  "test-generation_cpp_JBAICloudWithoutControl",
  "test-generation_cpp_JBAICloudControl",
  "test-generation_go_NativeLLaMA",
  "test-generation_go_JBAICloudWithoutControl",
  "test-generation_go_JBAICloudControl",
  "test-generation_html_NativeLLaMA",
  "test-generation_html_JBAICloudWithoutControl",
  "test-generation_html_JBAICloudControl",
  "test-generation_javascript_NativeLLaMA",
  "test-generation_javascript_JBAICloudWithoutControl",
  "test-generation_javascript_JBAICloudControl",
  "test-generation_java_NativeLLaMA",
  "test-generation_java_JBAICloudWithoutControl",
  "test-generation_java_JBAICloudControl",
  "test-generation_kotlin_NativeLLaMA",
  "test-generation_kotlin_JBAICloudWithoutControl",
  "test-generation_kotlin_JBAICloudControl",
  "test-generation_php_NativeLLaMA",
  "test-generation_php_JBAICloudWithoutControl",
  "test-generation_php_JBAICloudControl",
  "test-generation_python_NativeLLaMA",
  "test-generation_python_JBAICloudWithoutControl",
  "test-generation_python_JBAICloudControl",
  "test-generation_ruby_NativeLLaMA",
  "test-generation_ruby_JBAICloudWithoutControl",
  "test-generation_ruby_JBAICloudControl",
  "test-generation_rust_NativeLLaMA",
  "test-generation_rust_JBAICloudWithoutControl",
  "test-generation_rust_JBAICloudControl",
  "test-generation_scala_NativeLLaMA",
  "test-generation_scala_JBAICloudWithoutControl",
  "test-generation_scala_JBAICloudControl",
  "test-generation_terraform_NativeLLaMA",
  "test-generation_terraform_JBAICloudWithoutControl",
  "test-generation_terraform_JBAICloudControl",
]

const chartsDeclaration: ChartDefinition[] = projectDeclaration.map((project) => {
  return {
    labels: [project],
    measures: ["SyntaxErrorsSessionRatio"],
    projects: [project],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
