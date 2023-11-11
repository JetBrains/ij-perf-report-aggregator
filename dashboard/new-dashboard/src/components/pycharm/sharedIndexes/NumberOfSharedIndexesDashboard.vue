<template>
  <DashboardPage
    db-name="perfint"
    table="pycharm"
    persistent-id="py_number_of_shared_indexes_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (Django)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['django-with-project-sharedIndexes', 'django-with-bundled-sharedIndexes', 'django-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (Empty Project)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['empty project-with-project-sharedIndexes', 'empty project-with-bundled-sharedIndexes', 'empty project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (Pandas)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['pandas-with-project-sharedIndexes', 'pandas-with-bundled-sharedIndexes', 'pandas-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of indexing files by shared indexes (Keras)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['keras-with-project-sharedIndexes', 'keras-with-bundled-sharedIndexes', 'keras-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Nova)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['nova-with-project-sharedIndexes', 'nova-with-bundled-sharedIndexes', 'nova-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (PyTorch)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['pytorch-with-project-sharedIndexes', 'pytorch-with-bundled-sharedIndexes', 'pytorch-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (BoTo3)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['boto-with-project-sharedIndexes', 'boto-with-bundled-sharedIndexes', 'boto-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Matplotlib)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['matplotlib-with-project-sharedIndexes', 'matplotlib-with-bundled-sharedIndexes', 'matplotlib-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Dumb Mode Time (Matplotlib)"
        measure="numberOfFilesIndexedByExtensions"
        :projects="['matplotlib-with-project-sharedIndexes', 'matplotlib-with-bundled-sharedIndexes', 'matplotlib-without-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
