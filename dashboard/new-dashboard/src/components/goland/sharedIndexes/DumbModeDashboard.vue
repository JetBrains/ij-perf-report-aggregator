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
            'go-' + test_project.real_name + '-bundled-sharedIndexes',
            'go-' + test_project.real_name + '-with-generated-sharedIndexes',
            'go-' + test_project.real_name + '-without-sharedIndexes',
          ]"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"

interface test_project {
  pretty_name: string
  real_name: string
}

const projects: test_project[] = [
  {
    pretty_name: "Empty Project",
    real_name: "empty-project",
  },
  {
    pretty_name: "Terraform",
    real_name: "terraform",
  },
  {
    pretty_name: "Kratos",
    real_name: "kratos",
  },
  {
    pretty_name: "Kubernetes",
    real_name: "kubernetes",
  },
]
</script>
