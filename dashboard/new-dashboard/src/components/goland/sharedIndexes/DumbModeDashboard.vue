<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="go_dumb_mode_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :project="projects"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="test_project in projects"
          :key="'Dumb Mode Time (' + test_project.pretty_name + ')'"
          :label="'Dumb Mode Time (' + test_project.pretty_name + ')'"
          measure="dumbModeTimeWithPauses"
          :projects="[
            test_project.real_name + '-bundled-sharedIndexes',
            test_project.real_name + '-with-generated-sharedIndexes',
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
    pretty_name: "Empty Project",
    real_name: "go-empty-project",
  },
  {
    pretty_name: "Terraform",
    real_name: "go-terraform",
  },
  {
    pretty_name: "Kratos",
    real_name: "go-kratos",
  },
  {
    pretty_name: "Kubernetes",
    real_name: "kubernetes",
  },
]
</script>
