<template>
  <DashboardPage
    db-name="perfint"
    table="ideaSharedIndices"
    persistent-id="code_analysis_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :project="projects"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="test_project in projects"
          :key="'First Code Analysis (' + test_project.pretty_name + ')'"
          :label="'First Code Analysis (' + test_project.pretty_name + ')'"
          measure="firstCodeAnalysis"
          :projects="[
            test_project.real_name + '-downloaded-sharedIndexes',
            test_project.real_name + '-with-java-sharedIndexes',
            test_project.real_name + '-with-maven-sharedIndexes',
            test_project.real_name + '-with-project-sharedIndexes',
            test_project.real_name + '-without-sharedIndexes',
          ]"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import { SharedIndicesProject } from "../../common/sharedIndices/SharedIndicesProject"

const projects: SharedIndicesProject[] = [
  {
    pretty_name: "ToolboxEnterprise",
    real_name: "tbe",
  },
  {
    pretty_name: "Kotlin Serialization",
    real_name: "serialization",
  },
  {
    pretty_name: "Kotlin Coroutines",
    real_name: "coroutines",
  },
  {
    pretty_name: "Grails",
    real_name: "grails",
  },
  {
    pretty_name: "Java Design Patterns",
    real_name: "javaDesignPatterns",
  },
  {
    pretty_name: "IntelliJ",
    real_name: "intellij",
  },
  {
    pretty_name: "Community",
    real_name: "intellij-community",
  },
  {
    pretty_name: "SpringBoot",
    real_name: "spring-boot",
  },
  {
    pretty_name: "JDK",
    real_name: "jdk-only",
  },
]
</script>
