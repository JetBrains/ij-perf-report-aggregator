<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="pycharm"
    persistent-id="pycharm_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="w-1/2">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#Python'"
          :title="'Indexing Python (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        :measure="['indexing', 'indexingTimeWithoutPauses']"
        :projects="['django/indexing', 'empty/indexing', 'flask/indexing', 'keras/indexing', 'mypy/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning"
        :measure="['scanningTimeWithoutPauses', 'scanning']"
        :projects="['django/indexing', 'empty/indexing', 'flask/indexing', 'keras/indexing', 'mypy/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files"
        measure="numberOfIndexedFiles"
        :projects="['django/indexing', 'empty/indexing', 'flask/indexing', 'keras/indexing', 'mypy/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Global Inspections"
        measure="globalInspections"
        :projects="['django/inspection', 'flask/inspection', 'keras/inspection', 'mypy/inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Code Analysis"
        measure="firstCodeAnalysis"
        :projects="[
          'django/findUsages/ForeignKey',
          'django/findUsages/Form',
          'django/findUsages/Model',
          'flask/findUsages/Flask',
          'flask/findUsages/request',
          'keras/findUsages/Sequential',
          'mypy/findUsages/Errors',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages"
        measure="findUsages"
        :projects="[
          'django/findUsages/ForeignKey',
          'django/findUsages/Form',
          'django/findUsages/Model',
          'flask/findUsages/Flask',
          'flask/findUsages/request',
          'keras/findUsages/Sequential',
          'mypy/findUsages/Errors',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
