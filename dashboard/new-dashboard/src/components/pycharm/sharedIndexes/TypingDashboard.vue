<template>
  <DashboardPage
    db-name="perfint"
    table="pycharm"
    persistent-id="py_typing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Typing (Django)"
        measure="test#max_awt_delay"
        :projects="['django-with-project-sharedIndexes', 'django-with-bundled-sharedIndexes', 'django-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (Empty Project)"
        measure="test#max_awt_delay"
        :projects="['empty project-with-project-sharedIndexes', 'empty project-with-bundled-sharedIndexes', 'empty project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (Pandas)"
        measure="test#max_awt_delay"
        :projects="['pandas-with-project-sharedIndexes', 'pandas-with-bundled-sharedIndexes', 'pandas-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing (Keras)"
        measure="test#max_awt_delay"
        :projects="['keras-with-project-sharedIndexes', 'keras-with-bundled-sharedIndexes', 'keras-without-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
