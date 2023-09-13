<template>
  <DashboardPage
    db-name="perfint"
    table="ideaSharedIndices"
    persistent-id="code_analysis_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="First Code Analysis (ToolboxEnterprise)"
        measure="firstCodeAnalysis"
        :projects="['tbe-downloaded-sharedIndexes', 'tbe-with-java-sharedIndexes', 'tbe-with-maven-sharedIndexes', 'tbe-with-project-sharedIndexes', 'tbe-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Kotlin Serialization)"
        measure="firstCodeAnalysis"
        :projects="[
          'serialization-downloaded-sharedIndexes',
          'serialization-with-java-sharedIndexes',
          'serialization-with-maven-sharedIndexes',
          'serialization-with-project-sharedIndexes',
          'serialization-without-sharedIndexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Kotlin Coroutines)"
        measure="firstCodeAnalysis"
        :projects="[
          'coroutines-downloaded-sharedIndexes',
          'coroutines-with-java-sharedIndexes',
          'coroutines-with-maven-sharedIndexes',
          'coroutines-with-project-sharedIndexes',
          'coroutines-without-sharedIndexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Grails)"
        measure="firstCodeAnalysis"
        :projects="[
          'grails-downloaded-sharedIndexes',
          'grails-with-java-sharedIndexes',
          'grails-with-maven-sharedIndexes',
          'grails-with-project-sharedIndexes',
          'grails-without-sharedIndexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Java Design Patterns)"
        measure="firstCodeAnalysis"
        :projects="[
          'javaDesignPatterns-downloaded-sharedIndexes',
          'javaDesignPatterns-with-java-sharedIndexes',
          'javaDesignPatterns-with-maven-sharedIndexes',
          'javaDesignPatterns-with-project-sharedIndexes',
          'javaDesignPatterns-without-sharedIndexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (IntelliJ)"
        measure="firstCodeAnalysis"
        :projects="[
          'intellij-downloaded-sharedIndexes',
          'intellij-with-java-sharedIndexes',
          'intellij-with-maven-sharedIndexes',
          'intellij-with-project-sharedIndexes',
          'intellij-without-sharedIndexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Community)"
        measure="firstCodeAnalysis"
        :projects="[
          'intellij-community-downloaded-sharedIndexes',
          'intellij-community-with-java-sharedIndexes',
          'intellij-community-with-maven-sharedIndexes',
          'intellij-community-with-project-sharedIndexes',
          'intellij-community-without-sharedIndexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (SpringBoot)"
        measure="firstCodeAnalysis"
        :projects="[
          'spring-boot-downloaded-sharedIndexes',
          'spring-boot-with-java-sharedIndexes',
          'spring-boot-with-maven-sharedIndexes',
          'spring-boot-with-project-sharedIndexes',
          'spring-boot-without-sharedIndexes',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
