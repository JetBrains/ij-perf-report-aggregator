<template>
  <DashboardPage
    db-name="perfint"
    table="ideaSharedIndices"
    persistent-id="shared_indexes_dashboard"
    initial-machine="Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Indexing (IntelliJ project)"
        measure="indexing"
        :projects="[
          'ijx-intellij-speed/downloading-shared-indexes',
          'ijx-intellij-speed/shared-indexes',
          'ijx-intellij-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (IntelliJ project)"
        measure="scanning"
        :projects="[
          'ijx-intellij-speed/downloading-shared-indexes',
          'ijx-intellij-speed/shared-indexes',
          'ijx-intellij-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by usual indexes (IntelliJ project)"
        measure="numberOfIndexedFiles"
        :projects="[
          'ijx-intellij-speed/downloading-shared-indexes',
          'ijx-intellij-speed/shared-indexes',
          'ijx-intellij-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (IntelliJ project)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="[
          'ijx-intellij-speed/downloading-shared-indexes',
          'ijx-intellij-speed/shared-indexes',
          'ijx-intellij-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing (JDK Maven project)"
        measure="indexing"
        :projects="[
          'ijx-jdk-maven-speed/downloading-shared-indexes',
          'ijx-jdk-maven-speed/shared-indexes',
          'ijx-jdk-maven-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (JDK Maven project)"
        measure="scanning"
        :projects="[
          'ijx-jdk-maven-speed/downloading-shared-indexes',
          'ijx-jdk-maven-speed/shared-indexes',
          'ijx-jdk-maven-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by usual indexes (JDK Maven project)"
        measure="numberOfIndexedFiles"
        :projects="[
          'ijx-jdk-maven-speed/downloading-shared-indexes',
          'ijx-jdk-maven-speed/shared-indexes',
          'ijx-jdk-maven-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (JDK Maven project)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="[
          'ijx-jdk-maven-speed/downloading-shared-indexes',
          'ijx-jdk-maven-speed/shared-indexes',
          'ijx-jdk-maven-speed/usual-indexes',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing (TBE project)"
        measure="indexing"
        :projects="[
          'tbe-downloaded-sharedIndexes',
          'tbe-with-java-sharedIndexes',
          'tbe-with-maven-sharedIndexes',
          'tbe-with-project-sharedIndexes',
          'tbe-without-sharedIndexes'
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (TBE project)"
        measure="scanning"
        :projects="[
          'tbe-downloaded-sharedIndexes',
          'tbe-with-java-sharedIndexes',
          'tbe-with-maven-sharedIndexes',
          'tbe-with-project-sharedIndexes',
          'tbe-without-sharedIndexes'
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by usual indexes (TBE project)"
        measure="numberOfIndexedFiles"
        :projects="[
          'tbe-downloaded-sharedIndexes',
          'tbe-with-java-sharedIndexes',
          'tbe-with-maven-sharedIndexes',
          'tbe-with-project-sharedIndexes',
          'tbe-without-sharedIndexes'
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (TBE project)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="[
          'tbe-downloaded-sharedIndexes',
          'tbe-with-java-sharedIndexes',
          'tbe-with-maven-sharedIndexes',
          'tbe-with-project-sharedIndexes',
          'tbe-without-sharedIndexes'
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion (TBE project)"
        measure="completion"
        :projects="[
          'tbe-downloaded-sharedIndexes',
          'tbe-with-java-sharedIndexes',
          'tbe-with-maven-sharedIndexes',
          'tbe-with-project-sharedIndexes',
          'tbe-without-sharedIndexes'
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (TBE project)"
        measure="firstCodeAnalysis"
        :projects="[
          'tbe-downloaded-sharedIndexes',
          'tbe-with-java-sharedIndexes',
          'tbe-with-maven-sharedIndexes',
          'tbe-with-project-sharedIndexes',
          'tbe-without-sharedIndexes'
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (TBE project)"
        measure="findUsages"
        :projects="[
          'tbe-downloaded-sharedIndexes',
          'tbe-with-java-sharedIndexes',
          'tbe-with-maven-sharedIndexes',
          'tbe-with-project-sharedIndexes',
          'tbe-without-sharedIndexes'
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
