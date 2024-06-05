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
  "code-generation_css_NativeLLaMA",
  "code-generation_css_JBAICloudWithoutControl",
  "code-generation_css_JBAICloudControl",
  "code-generation_csharp_NativeLLaMA",
  "code-generation_csharp_JBAICloudWithoutControl",
  "code-generation_csharp_JBAICloudControl",
  "code-generation_cpp_NativeLLaMA",
  "code-generation_cpp_JBAICloudWithoutControl",
  "code-generation_cpp_JBAICloudControl",
  "code-generation_go_NativeLLaMA",
  "code-generation_go_JBAICloudWithoutControl",
  "code-generation_go_JBAICloudControl",
  "code-generation_html_NativeLLaMA",
  "code-generation_html_JBAICloudWithoutControl",
  "code-generation_html_JBAICloudControl",
  "code-generation_javascript_NativeLLaMA",
  "code-generation_javascript_JBAICloudWithoutControl",
  "code-generation_javascript_JBAICloudControl",
  "code-generation_java_NativeLLaMA",
  "code-generation_java_JBAICloudWithoutControl",
  "code-generation_java_JBAICloudControl",
  "code-generation_kotlin_NativeLLaMA",
  "code-generation_kotlin_JBAICloudWithoutControl",
  "code-generation_kotlin_JBAICloudControl",
  "code-generation_php_NativeLLaMA",
  "code-generation_php_JBAICloudWithoutControl",
  "code-generation_php_JBAICloudControl",
  "code-generation_python_NativeLLaMA",
  "code-generation_python_JBAICloudWithoutControl",
  "code-generation_python_JBAICloudControl",
  "code-generation_ruby_NativeLLaMA",
  "code-generation_ruby_JBAICloudWithoutControl",
  "code-generation_ruby_JBAICloudControl",
  "code-generation_rust_NativeLLaMA",
  "code-generation_rust_JBAICloudWithoutControl",
  "code-generation_rust_JBAICloudControl",
  "code-generation_scala_NativeLLaMA",
  "code-generation_scala_JBAICloudWithoutControl",
  "code-generation_scala_JBAICloudControl",
  "code-generation_terraform_NativeLLaMA",
  "code-generation_terraform_JBAICloudWithoutControl",
  "code-generation_terraform_JBAICloudControl",
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
