<template>
  <DashboardPage
    db-name="perfint"
    table="pycharm"
    persistent-id="py_indexing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Indexing (Django)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['django-with-project-sharedIndexes', 'django-with-bundled-sharedIndexes', 'django-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing (Empty Project)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['empty project-with-project-sharedIndexes', 'empty project-with-bundled-sharedIndexes', 'empty project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing (Pandas)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['pandas-with-project-sharedIndexes', 'pandas-with-bundled-sharedIndexes', 'pandas-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing (Keras)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['keras-with-project-sharedIndexes', 'keras-with-bundled-sharedIndexes', 'keras-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Nova)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['nova-with-project-sharedIndexes', 'nova-with-bundled-sharedIndexes', 'nova-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (PyTorch)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['pytorch-with-project-sharedIndexes', 'pytorch-with-bundled-sharedIndexes', 'pytorch-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (BoTo3)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['boto-with-project-sharedIndexes', 'boto-with-bundled-sharedIndexes', 'boto-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Matplotlib)"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['matplotlib-with-project-sharedIndexes', 'matplotlib-with-bundled-sharedIndexes', 'matplotlib-without-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
