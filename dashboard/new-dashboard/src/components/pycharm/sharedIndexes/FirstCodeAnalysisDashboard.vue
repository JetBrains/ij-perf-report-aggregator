<template>
  <DashboardPage
    db-name="perfint"
    table="pycharm"
    persistent-id="py_code_analysis_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Django)"
        measure="firstCodeAnalysis"
        :projects="['django-with-project-sharedIndexes', 'django-with-bundled-sharedIndexes', 'django-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Empty Project)"
        measure="firstCodeAnalysis"
        :projects="['empty project-with-project-sharedIndexes', 'empty project-with-bundled-sharedIndexes', 'empty project-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Pandas)"
        measure="firstCodeAnalysis"
        :projects="['pandas-with-project-sharedIndexes', 'pandas-with-bundled-sharedIndexes', 'pandas-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis (Keras)"
        measure="firstCodeAnalysis"
        :projects="['keras-with-project-sharedIndexes', 'keras-with-bundled-sharedIndexes', 'keras-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (Nova)"
        measure="firstCodeAnalysis"
        :projects="['nova-with-project-sharedIndexes', 'nova-with-bundled-sharedIndexes', 'nova-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (PyTorch)"
        measure="firstCodeAnalysis"
        :projects="['pytorch-with-project-sharedIndexes', 'pytorch-with-bundled-sharedIndexes', 'pytorch-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (BoTo3)"
        measure="firstCodeAnalysis"
        :projects="['boto-with-project-sharedIndexes', 'boto-with-bundled-sharedIndexes', 'boto-without-sharedIndexes']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning (Matplotlib)"
        measure="firstCodeAnalysis"
        :projects="['matplotlib-with-project-sharedIndexes', 'matplotlib-with-bundled-sharedIndexes', 'matplotlib-without-sharedIndexes']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
