<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="go_number_of_shared_indexes_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (Empty Project)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['go-empty-project-bundled-sharedIndexes', 'go-empty-project-with-generated-sharedIndexes', 'go-empty-project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (Terraform)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['go-terraform-without-sharedIndexes', 'go-terraform-bundled-sharedIndexes', 'go-terraform-with-generated-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (Kratos)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['go-kratos-without-sharedIndexes', 'go-kratos-bundled-sharedIndexes', 'go-kratos-with-generated-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Kubernetes)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['kubernetes-without-sharedIndexes', 'kubernetes-bundled-sharedIndexes', 'kubernetes-with-generated-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
