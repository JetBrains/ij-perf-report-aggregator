<template>
  <DashboardPage
    db-name="perfint"
    table="ruby"
    persistent-id="rubymine_indexing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Indexing Time"
        :measure="['indexingTimeWithoutPauses']"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Diaspora', 'RBS Collection', 'Gitlab', 'Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files"
        measure="numberOfIndexedFiles"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Diaspora', 'RBS Collection', 'Gitlab', 'Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Index Size"
        measure="indexSize"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Diaspora', 'RBS Collection', 'Gitlab', 'Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing Time"
        :measure="[
          'processingTime#Ruby',
          'processingTime#JavaScript',
          'processingTime#Markdown',
          'processingTime#TypeScript',
          'processingTime#Vue.js',
          'processingTime#YAML',
          'processingTime#Slim',
          'processingTime#SQL',
          'processingTime#RHTML',
          'processingTime#RBS',
          'processingTime#JSON',
          'processingTime#Haml',
          'processingTime#Cucumber',
          'processingTime#CSS',
          'processingTime#SCSS',
        ]"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Diaspora', 'RBS Collection', 'Gitlab', 'Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning"
        :measure="['scanning', 'scanningTimeWithoutPauses']"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Diaspora', 'RBS Collection', 'Gitlab', 'Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing Time (With Data Before 11/07/2023)"
        :measure="['indexing', 'indexingTimeWithoutPauses']"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Diaspora', 'RBS Collection', 'Gitlab', 'Redmine']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
