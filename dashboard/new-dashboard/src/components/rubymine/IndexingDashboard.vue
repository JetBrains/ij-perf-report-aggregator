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
        label="Processing Time: Gitlab"
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
        :projects="['gitlab-project-test/indexing']"
        :aliases="['Gitlab']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing Time: Diaspora"
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
        :projects="['diaspora-project-test/indexing']"
        :aliases="['Diaspora']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing Time: RBS Collection"
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
        :projects="['gem-rbs-collection-indexing-test/indexing']"
        :aliases="['RBS Collection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing Time: Redmine"
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
        :projects="['redmine-project-test/indexing']"
        :aliases="['Redmine']"
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
        label="GC Pause, ms"
        measure="gcPause"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Disaspora', 'RBS Collection', 'GitLab', 'Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="GC Memory Collected, Mb"
        measure="freedMemoryByGC"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :aliases="['Disaspora', 'RBS Collection', 'GitLab', 'Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Lexing Speed, Diaspora, kb/sec"
        :measure="[
          'lexingSpeed#CSS',
          'lexingSpeed#CoffeeScript',
          'lexingSpeed#ECMAScript.6',
          'lexingSpeed#Flow.JS',
          'lexingSpeed#GenericSQL',
          'lexingSpeed#HTML',
          'lexingSpeed#Haml',
          'lexingSpeed#JSON',
          'lexingSpeed#Markdown',
          'lexingSpeed#PostCSS',
          'lexingSpeed#RHTML',
          'lexingSpeed#SCSS',
          'lexingSpeed#TypeScript',
          'lexingSpeed#Vue',
          'lexingSpeed#ruby',
          'lexingSpeed#ruby.rbs',
          'lexingSpeed#yaml',
        ]"
        :projects="['diaspora-project-test/indexing']"
        :aliases="['Disaspora']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Lexing Speed, RBS Collection, kb/sec"
        :measure="[
          'lexingSpeed#CSS',
          'lexingSpeed#CoffeeScript',
          'lexingSpeed#ECMAScript.6',
          'lexingSpeed#Flow.JS',
          'lexingSpeed#GenericSQL',
          'lexingSpeed#HTML',
          'lexingSpeed#Haml',
          'lexingSpeed#JSON',
          'lexingSpeed#Markdown',
          'lexingSpeed#PostCSS',
          'lexingSpeed#RHTML',
          'lexingSpeed#SCSS',
          'lexingSpeed#TypeScript',
          'lexingSpeed#Vue',
          'lexingSpeed#ruby',
          'lexingSpeed#ruby.rbs',
          'lexingSpeed#yaml',
        ]"
        :projects="['gem-rbs-collection-indexing-test/indexing']"
        :aliases="['RBS Collection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Lexing Speed, Gitlab, kb/sec"
        :measure="[
          'lexingSpeed#CSS',
          'lexingSpeed#CoffeeScript',
          'lexingSpeed#ECMAScript.6',
          'lexingSpeed#Flow.JS',
          'lexingSpeed#GenericSQL',
          'lexingSpeed#HTML',
          'lexingSpeed#Haml',
          'lexingSpeed#JSON',
          'lexingSpeed#Markdown',
          'lexingSpeed#PostCSS',
          'lexingSpeed#RHTML',
          'lexingSpeed#SCSS',
          'lexingSpeed#TypeScript',
          'lexingSpeed#Vue',
          'lexingSpeed#ruby',
          'lexingSpeed#ruby.rbs',
          'lexingSpeed#yaml',
        ]"
        :projects="['gitlab-project-test/indexing']"
        :aliases="['GitLab']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Lexing Speed, Redmine, kb/sec"
        :measure="[
          'lexingSpeed#CSS',
          'lexingSpeed#CoffeeScript',
          'lexingSpeed#ECMAScript.6',
          'lexingSpeed#Flow.JS',
          'lexingSpeed#GenericSQL',
          'lexingSpeed#HTML',
          'lexingSpeed#Haml',
          'lexingSpeed#JSON',
          'lexingSpeed#Markdown',
          'lexingSpeed#PostCSS',
          'lexingSpeed#RHTML',
          'lexingSpeed#SCSS',
          'lexingSpeed#TypeScript',
          'lexingSpeed#Vue',
          'lexingSpeed#ruby',
          'lexingSpeed#ruby.rbs',
          'lexingSpeed#yaml',
        ]"
        :projects="['redmine-project-test/indexing']"
        :aliases="['Redmine']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Parsing Speed, Diaspora, kb/sec"
        :measure="[
          'parsingSpeed#CSS',
          'parsingSpeed#CoffeeScript',
          'parsingSpeed#ECMAScript.6',
          'parsingSpeed#Flow.JS',
          'parsingSpeed#GenericSQL',
          'parsingSpeed#HTML',
          'parsingSpeed#Haml',
          'parsingSpeed#JSON',
          'parsingSpeed#Markdown',
          'parsingSpeed#PostCSS',
          'parsingSpeed#RHTML',
          'parsingSpeed#SCSS',
          'parsingSpeed#TypeScript',
          'parsingSpeed#Vue',
          'parsingSpeed#ruby',
          'parsingSpeed#ruby.rbs',
          'parsingSpeed#yaml',
        ]"
        :projects="['diaspora-project-test/indexing']"
        :aliases="['Disaspora']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Parsing Speed, RBS Collection, kb/sec"
        :measure="[
          'parsingSpeed#CSS',
          'parsingSpeed#CoffeeScript',
          'parsingSpeed#ECMAScript.6',
          'parsingSpeed#Flow.JS',
          'parsingSpeed#GenericSQL',
          'parsingSpeed#HTML',
          'parsingSpeed#Haml',
          'parsingSpeed#JSON',
          'parsingSpeed#Markdown',
          'parsingSpeed#PostCSS',
          'parsingSpeed#RHTML',
          'parsingSpeed#SCSS',
          'parsingSpeed#TypeScript',
          'parsingSpeed#Vue',
          'parsingSpeed#ruby',
          'parsingSpeed#ruby.rbs',
          'parsingSpeed#yaml',
        ]"
        :projects="['gem-rbs-collection-indexing-test/indexing']"
        :aliases="['RBS Collection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Parsing Speed, Gitlab, kb/sec"
        :measure="[
          'parsingSpeed#CSS',
          'parsingSpeed#CoffeeScript',
          'parsingSpeed#ECMAScript.6',
          'parsingSpeed#Flow.JS',
          'parsingSpeed#GenericSQL',
          'parsingSpeed#HTML',
          'parsingSpeed#Haml',
          'parsingSpeed#JSON',
          'parsingSpeed#Markdown',
          'parsingSpeed#PostCSS',
          'parsingSpeed#RHTML',
          'parsingSpeed#SCSS',
          'parsingSpeed#TypeScript',
          'parsingSpeed#Vue',
          'parsingSpeed#ruby',
          'parsingSpeed#ruby.rbs',
          'parsingSpeed#yaml',
        ]"
        :projects="['gitlab-project-test/indexing']"
        :aliases="['GitLab']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Parsing Speed, Redmine, kb/sec"
        :measure="[
          'parsingSpeed#CSS',
          'parsingSpeed#CoffeeScript',
          'parsingSpeed#ECMAScript.6',
          'parsingSpeed#Flow.JS',
          'parsingSpeed#GenericSQL',
          'parsingSpeed#HTML',
          'parsingSpeed#Haml',
          'parsingSpeed#JSON',
          'parsingSpeed#Markdown',
          'parsingSpeed#PostCSS',
          'parsingSpeed#RHTML',
          'parsingSpeed#SCSS',
          'parsingSpeed#TypeScript',
          'parsingSpeed#Vue',
          'parsingSpeed#ruby',
          'parsingSpeed#ruby.rbs',
          'parsingSpeed#yaml',
        ]"
        :projects="['redmine-project-test/indexing']"
        :aliases="['Redmine']"
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
