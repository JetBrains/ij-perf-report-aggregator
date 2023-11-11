<template>
  <DashboardPage
    db-name="perfint"
    table="pycharm"
    persistent-id="py_number_of_files_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Number of indexing files (Django)"
        measure="numberOfIndexedFiles"
        :projects="['django-with-project-sharedIndexes', 'django-with-bundled-sharedIndexes', 'django-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files (Empty Project)"
        measure="numberOfIndexedFiles"
        :projects="['empty project-with-project-sharedIndexes', 'empty project-with-bundled-sharedIndexes', 'empty project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files (Pandas)"
        measure="numberOfIndexedFiles"
        :projects="['pandas-with-project-sharedIndexes', 'pandas-with-bundled-sharedIndexes', 'pandas-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files (Keras)"
        measure="numberOfIndexedFiles"
        :projects="['keras-with-project-sharedIndexes', 'keras-with-bundled-sharedIndexes', 'keras-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Nova)"
        measure="numberOfIndexedFiles"
        :projects="['nova-with-project-sharedIndexes', 'nova-with-bundled-sharedIndexes', 'nova-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (PyTorch)"
        measure="numberOfIndexedFiles"
        :projects="['pytorch-with-project-sharedIndexes', 'pytorch-with-bundled-sharedIndexes', 'pytorch-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (BoTo3)"
        measure="numberOfIndexedFiles"
        :projects="['boto-with-project-sharedIndexes', 'boto-with-bundled-sharedIndexes', 'boto-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Matplotlib)"
        measure="numberOfIndexedFiles"
        :projects="['matplotlib-with-project-sharedIndexes', 'matplotlib-with-bundled-sharedIndexes', 'matplotlib-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Matplotlib)"
        measure="numberOfIndexedFiles"
        :projects="['matplotlib-with-project-sharedIndexes', 'matplotlib-with-bundled-sharedIndexes', 'matplotlib-without-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
