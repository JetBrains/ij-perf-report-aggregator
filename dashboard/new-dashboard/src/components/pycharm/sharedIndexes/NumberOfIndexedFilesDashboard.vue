<template>
  <DashboardPage
    db-name="perfint"
    table="pycharm"
    persistent-id="py_number_of_files_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :project="projects"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="test_project in projects"
          :key="'Number of indexed files (' + test_project.pretty_name + ')'"
          :label="'Number of indexed files (' + test_project.pretty_name + ')'"
          measure="numberOfIndexedFiles"
          :projects="[
            test_project.real_name + '-with-bundled-sharedIndexes',
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
    pretty_name: "Empty Project",
    real_name: "empty project",
  },
  {
    pretty_name: "Django",
    real_name: "django",
  },
  {
    pretty_name: "Pandas",
    real_name: "pandas",
  },
  {
    pretty_name: "Keras",
    real_name: "keras",
  },
  {
    pretty_name: "Nova",
    real_name: "nova",
  },
  {
    pretty_name: "PyTorch",
    real_name: "pytorch",
  },
  {
    pretty_name: "BoTo3",
    real_name: "boto",
  },
  {
    pretty_name: "Matplotlib",
    real_name: "matplotlib",
  },
]
</script>
