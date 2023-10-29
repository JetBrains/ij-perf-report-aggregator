<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="go_find_usages_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Find Usages (Empty Project)"
        measure="findUsages"
        :projects="['go-empty-project-bundled-sharedIndexes', 'go-empty-project-with-generated-sharedIndexes', 'go-empty-project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (Terraform)"
        measure="findUsages"
        :projects="['go-terraform-without-sharedIndexes', 'go-terraform-bundled-sharedIndexes', 'go-terraform-with-generated-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (Kratos)"
        measure="findUsages"
        :projects="['go-kratos-without-sharedIndexes', 'go-kratos-bundled-sharedIndexes', 'go-kratos-with-generated-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
