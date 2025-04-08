<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="ruby"
    persistent-id="rubymine_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
      <GroupProjectsWithClientChart
        label="Inspections"
        measure="globalInspections"
        :projects="[
          'diaspora-project-inspections-test/inspection-RubyResolve-app',
          'diaspora-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-RubyResolve-app',
          'gitlab-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-App-RubyResolve',
          'gitlab-project-inspections-test/inspection-App',
          'redmine-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-app',
          'mastodon-project-inspections-test/inspection-RubyResolve-app',
          'mastodon-project-inspections-test/inspection-app',
        ]"
        :aliases="[
          'Unresolved References Inspection (DI)',
          'All Inspections (DI)',
          'Unresolved References Inspection (GL)',
          'All Inspections (GL)',
          'Unresolved References Inspection (GL)',
          'All Inspections (GL)',
          'Unresolved References Inspection (RM)',
          'All Inspections (RM)',
          'Unresolved References Inspection (MN)',
          'All Inspections (MN)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Gitlab Inspections"
        measure="globalInspections"
        :projects="[
          'gitlab-project-inspections-test/inspection-App',
          'gitlab-project-inspections-test/inspection-App-RubyResolve',
          'gitlab-project-inspections-test/inspection-Yaml',
          'gitlab-project-inspections-test/inspection-XML',
          'gitlab-project-inspections-test/inspection-WebStorm-AppSpecFrontend',
          'gitlab-project-inspections-test/inspection-Slim',
          'gitlab-project-inspections-test/inspection-RubyMine-App',
          'gitlab-project-inspections-test/inspection-Others',
          'gitlab-project-inspections-test/inspection-Markdown',
          'gitlab-project-inspections-test/inspection-Liquid',
          'gitlab-project-inspections-test/inspection-Haml',
          'gitlab-project-inspections-test/inspection-Erb',
          'gitlab-project-inspections-test/inspection-DataGrip',
        ]"
        :aliases="[
          'All (app/)',
          'Unresolved Ruby References (app/)',
          'Yaml',
          'XML',
          'All on WebStorm files (app/ & spec/frontend)',
          'Slim',
          'All on RubyMine files (app/)',
          'All on all other files',
          'Markdown',
          'Liquid',
          'Haml',
          'Erb',
          'All on DataGrip files',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="File Analysis on Open"
        measure="firstCodeAnalysis"
        :projects="[
          'GitLab/typing/lambda_body_in_class',
          'RUBY-26170/typing',
          'RUBY-29334/typing',
          'RUBY-29542/typing',
          'GitLab/typing/typing/user/method',
          'GitLab/typing/enter/parser/method',
          'GitLab/typing/enter/structure/inside_query',
          'GitLab/typing/enter/project_spec/describe',
          'GitLab/typing/enter/project_controller/class',
          'GitLab/typing/enter/mr_mail/class',
          'GitLab/typing/enter/user_show_view/before_div',
          'GitLab/typing/enter/routes_project/top',
          'GitLab/typing/enter/emojis_json/map',
          'diaspora-project-test/completion/constant',
          'gitlab-project-test/completion/constant',
          'redmine-project-test/completion/constant',
          'diaspora-project-test/completion/exceptions',
          'gitlab-project-test/completion/exceptions',
          'redmine-project-test/completion/exceptions',
          'diaspora-project-test/completion/exceptions-prefix',
          'gitlab-project-test/completion/exceptions-prefix',
          'redmine-project-test/completion/exceptions-prefix',
          'diaspora-project-test/completion/localization',
          'gitlab-project-test/completion/localization',
          'redmine-project-test/completion/localization',
          'diaspora-project-test/completion/method',
          'gitlab-project-test/completion/method',
          'redmine-project-test/completion/method',
          'diaspora-project-test/completion/qualified',
          'gitlab-project-test/completion/qualified',
          'redmine-project-test/completion/qualified',
          'RUBY-23764-Case1/ruby-23764-findusages-case1',
          'RUBY-23764-Case2/ruby-23764-findusages-case2',
          'gitlab-find-usages/ruby-23764-findusages-case1',
          'gitlab-find-usages/ruby-23764-findusages-case2',
          'RUBY-32357/module',
          'RUBY-32357/method',
          'RUBY-32357/singleton-method',
          'RUBY-32357/instance-variable',
          'RUBY-32357/class-variable',
          'RUBY-32357/global-variable',
          'RUBY-32357/delegate-method',
          'gitlab-find-usages/module',
          'gitlab-find-usages/method',
          'gitlab-find-usages/singleton-method',
          'gitlab-find-usages/instance-variable',
          'gitlab-find-usages/class-variable',
          'gitlab-find-usages/global-variable',
          'gitlab-find-usages/delegate-method',
          'mastodon-find-usages/i18n-key',
        ]"
        :aliases="[
          'Project Model (GL)',
          'swagger_helper.rb',
          'activerecord-generated.rbs',
          'User Model Spec (Rails Sample)',
          'User Model (GL)',
          'Ruby Parser',
          'structure.sql (GL)',
          'Project Model Spec (GL)',
          'Project Controller (GL)',
          'MR Mail (GL)',
          'Users View Haml (GL)',
          'Project Routes (GL)',
          'Emojis.json (GL)',
          'Message Model (DI)',
          'Clusters Controller (GL)',
          'Time Entry Activity Model (RM)',
          'Admins Controller (DI)',
          'Admin App Controller (GL)',
          'Account Controller (RM)',
          'Process Photo Worker (DI)',
          'Build Trace Chunk Model (GL)',
          'App Controller (RM)',
          'Admin Pods View Haml (DI)',
          'Admin Locale View Haml (GL)',
          'Time Entries Import View (RM)',
          'Conversation Visibility Model (DI)',
          'Epic Move List (GL)',
          'Auto Completes Controller (RM)',
          'Contacts Controller Spec (DI)',
          'Environment Entity Serializer (GL)',
          'Admin Info View Erb (RM)',
          'Users Spec Factory (GL)',
          'File Collection Spec (GL)',
          'Users Spec Factory (GL2)',
          'File Collection Spec (GL2)',
          'GitLab Routing Helper (GL)',
          'Todo Model (GL)',
          'GitLab Script (GL)',
          'Labels Finder (GL)',
          'Epics Swimlanes Drag Drop Spec (GL)',
          'Feature Flag Script (GL)',
          'Framework Registry Finder (GL)',
          'GitLab Routing Helper (GL)',
          'Todo Model (GL)',
          'GitLab Script (GL)',
          'Labels Finder (GL)',
          'Epics Swimlanes Drag Drop Spec (GL)',
          'Feature Flag Script (GL)',
          'Framework Registry Finder (GL)',
          'Locale File (MA)',
        ]"
      />
    </section>

    <section>
      <GroupProjectsWithClientChart
        label="GC Pause, ms"
        measure="gcPause"
        :projects="[
          'gitlab-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-App-RubyResolve',
          'gitlab-project-inspections-test/inspection-App',
          'mastodon-project-inspections-test/inspection-RubyResolve-app',
          'mastodon-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-Yaml',
          'gitlab-project-inspections-test/inspection-XML',
          'gitlab-project-inspections-test/inspection-WebStorm-AppSpecFrontend',
          'gitlab-project-inspections-test/inspection-Slim',
          'gitlab-project-inspections-test/inspection-RubyMine-App',
          'gitlab-project-inspections-test/inspection-Others',
          'gitlab-project-inspections-test/inspection-Markdown',
          'gitlab-project-inspections-test/inspection-Liquid',
          'gitlab-project-inspections-test/inspection-Haml',
          'gitlab-project-inspections-test/inspection-Erb',
          'gitlab-project-inspections-test/inspection-DataGrip',
        ]"
        :aliases="[
          'All Inspections (GL)',
          'Unresolved References Inspection (GL)',
          'All Inspections (RM)',
          'Unresolved References Inspection (RM)',
          'Unresolved References Inspection (GL)',
          'All Inspections (GL)',
          'Unresolved References Inspection (MN)',
          'All Inspections (MN)',
          'Yaml Inspections',
          'XML Inspections',
          'Inspections on WebStorm files (app/ & spec/frontend)',
          'Slim Inspections',
          'Inspections on RubyMine files (app/)',
          'Inspections on all other files',
          'Markdown Inspections',
          'Liquid Inspections',
          'Haml Inspections',
          'Erb Inspections',
          'Inspections on DataGrip files',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="GC Memory Collected, Mb"
        measure="freedMemoryByGC"
        :projects="[
          'gitlab-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-App-RubyResolve',
          'gitlab-project-inspections-test/inspection-App',
          'mastodon-project-inspections-test/inspection-RubyResolve-app',
          'mastodon-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-Yaml',
          'gitlab-project-inspections-test/inspection-XML',
          'gitlab-project-inspections-test/inspection-WebStorm-AppSpecFrontend',
          'gitlab-project-inspections-test/inspection-Slim',
          'gitlab-project-inspections-test/inspection-RubyMine-App',
          'gitlab-project-inspections-test/inspection-Others',
          'gitlab-project-inspections-test/inspection-Markdown',
          'gitlab-project-inspections-test/inspection-Liquid',
          'gitlab-project-inspections-test/inspection-Haml',
          'gitlab-project-inspections-test/inspection-Erb',
          'gitlab-project-inspections-test/inspection-DataGrip',
        ]"
        :aliases="[
          'All Inspections (GL)',
          'Unresolved References Inspection (GL)',
          'All Inspections (RM)',
          'Unresolved References Inspection (RM)',
          'Unresolved References Inspection (GL)',
          'All Inspections (GL)',
          'Unresolved References Inspection (MN)',
          'All Inspections (MN)',
          'Yaml Inspections',
          'XML Inspections',
          'Inspections on WebStorm files (app/ & spec/frontend)',
          'Slim Inspections',
          'Inspections on RubyMine files (app/)',
          'Inspections on all other files',
          'Markdown Inspections',
          'Liquid Inspections',
          'Haml Inspections',
          'Erb Inspections',
          'Inspections on DataGrip files',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Routes Highlighting (removed)"
        measure="highlighting"
        :projects="[
          'redmine-project-test/highlight/routes-highlighting',
          'gitlab-project-test/highlight/routes-highlighting',
          'diaspora-project-test/highlight/routes-highlighting',
        ]"
        :aliases="['Redmine', 'Gitlab', 'Diaspora']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { DataQuery, DataQueryExecutorConfiguration } from "../common/dataQuery"

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
