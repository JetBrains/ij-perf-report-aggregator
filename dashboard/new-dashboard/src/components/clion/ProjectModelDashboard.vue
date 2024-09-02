<template>
  <DashboardPage
    db-name="perfint"
    table="clion"
    persistent-id="clion_project_model_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <Divider title="Workspace Model" />

    <section>
      <CLionVsRadlerWorkspaceUpdateChart
        label="Workspace Model Update"
        :projects="['llvm/indexing', 'opencv/indexing', 'curl/indexing']"
        :names="['LLVM', 'OpenCV', 'curl']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import CLionVsRadlerWorkspaceUpdateChart from "./CLionVsRadlerWorkspaceUpdateChart.vue"
</script>
