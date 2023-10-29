<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="go_scanning_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Scanning (Empty Project)"
        :measure="['scanningTimeWithoutPauses']"
        :projects="['go-empty-project-with-project-sharedIndexes', 'go-empty-project-with-bundled-sharedIndexes', 'go-empty-project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (Terraform)"
        :measure="['scanningTimeWithoutPauses']"
        :projects="['go-terraform-with-project-sharedIndexes', 'go-terraform-with-bundled-sharedIndexes', 'go-terraform-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (Kratos)"
        :measure="['scanningTimeWithoutPauses']"
        :projects="['go-kratos-with-project-sharedIndexes', 'go-kratos-with-bundled-sharedIndexes', 'go-kratos-without-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
