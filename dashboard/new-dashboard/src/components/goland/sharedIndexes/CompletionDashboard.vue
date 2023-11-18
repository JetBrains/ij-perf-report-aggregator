<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="go_completion_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :project="projects"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="project in projects"
          :key="'Completion (' + project.pretty_name + ')'"
          :label="'Completion (' + project.pretty_name + ')'"
          measure="completion"
          :projects="['go-'+ project.real_name +'-bundled-sharedIndexes', 'go-'+ project.real_name +'-with-generated-sharedIndexes', 'go-'+ project.real_name +'-without-sharedIndexes']"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
interface project {
  pretty_name: string
  real_name: string
}
const projects: project[] = [
  {
    pretty_name: "Empty Project",
    real_name: "empty-project"
  },
  {
    pretty_name: "Terraform",
    real_name: "terraform"
  },
  {
    pretty_name: "Kratos",
    real_name: "kratos"
  }
]

</script>
