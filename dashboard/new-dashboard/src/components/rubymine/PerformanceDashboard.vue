<template>
  <DashboardPage
    v-slot="{ dashboardConfigurators}"
    db-name="perfint"
    table="ruby"
    persistent-id="rubymine_dashboard"
    initial-machine="Linux Munich i7-3770, 32 Gb"
  >
    <section>
      <GroupProjectsChart
        label="Indexing"
        measure="indexing"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files"
        measure="numberOfIndexedFiles"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning"
        measure="scanning"
        :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing', 'redmine-project-test/indexing']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inspection"
        measure="globalInspections"
        :projects="['diaspora-project-inspections-test/inspection-RubyResolve-app', 'diaspora-project-inspections-test/inspection-app',
                    'gitlab-project-inspections-test/inspection-RubyResolve-app', 'gitlab-project-inspections-test/inspection-app']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages: execution time"
        measure="findUsages"
        :projects="['RUBY-23764-Case1/ruby-23764-findusages-case1', 'RUBY-23764-Case2/ruby-23764-findusages-case2']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages: number of found usages"
        measure="findUsages#number"
        :projects="['RUBY-23764-Case1/ruby-23764-findusages-case1', 'RUBY-23764-Case2/ruby-23764-findusages-case2']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion Diaspora"
        measure="completion"
        :projects="['diaspora-project-test/completion/routes', 'diaspora-project-test/completion/exceptions', 'diaspora-project-test/completion/localization']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion Gitlab"
        measure="completion"
        :projects="['gitlab-project-test/completion/routes', 'gitlab-project-test/completion/exceptions', 'gitlab-project-test/completion/localization']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion Redmine"
        measure="completion"
        :projects="['redmine-project-test/completion/routes', 'redmine-project-test/completion/exceptions', 'redmine-project-test/completion/localization']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: average delay"
        measure="test#average_awt_delay"
        :projects="['RUBY-26170/typing', 'RUBY-29334/typing', 'RUBY-29542/typing']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: total time"
        measure="typing"
        :projects="['RUBY-26170/typing', 'RUBY-29334/typing', 'RUBY-29542/typing']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Get Symbol Members: execution time"
        measure="getSymbolMembers"
        :projects="['diaspora-project-test/getSymbolMembers-ApplicationController-false', 'diaspora-project-test/getSymbolMembers-ApplicationController-true',
                    'gitlab-project-test/getSymbolMembers-ApplicationController-false', 'gitlab-project-test/getSymbolMembers-ApplicationController-true',
                    'redmine-project-test/getSymbolMembers-ApplicationController-false', 'redmine-project-test/getSymbolMembers-ApplicationController-true']"
        :configurators="dashboardConfigurators"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Get Symbol Members: number"
        measure="getSymbolMembers#number"
        :projects="['diaspora-project-test/getSymbolMembers-ApplicationController-false', 'diaspora-project-test/getSymbolMembers-ApplicationController-true',
                    'gitlab-project-test/getSymbolMembers-ApplicationController-false', 'gitlab-project-test/getSymbolMembers-ApplicationController-true',
                    'redmine-project-test/getSymbolMembers-ApplicationController-false', 'redmine-project-test/getSymbolMembers-ApplicationController-true']"
        :configurators="dashboardConfigurators"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>