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
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :is-like="true"
          :title="'Completion'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="[...averagesConfigurators, typingOnlyConfigurator]"
          :aggregated-measure="'test#average_awt_delay'"
          :title="'UI responsiveness during typing'"
          :chart-color="'#F2994A'"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Completion"
        :measure="['fus_time_to_show_90p']"
        :projects="['edx-platform (Django)/completion/model', 'edx-platform (Django)/completion/view']"
      />
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
      <GroupProjectsWithClientChart
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
      <GroupProjectsWithClientChart
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
import { DataQuery, DataQueryExecutorConfiguration } from "../common/dataQuery"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const typingOnlyConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({ f: "project", v: "%typing", o: "like" })
    return true
  },
  createObservable() {
    return null
  },
}
</script>
