<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="ruby"
    persistent-id="rubymine_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
        label="Inspections"
        measure="globalInspections"
        :projects="[
          'diaspora-project-inspections-test/inspection-RubyResolve-app',
          'diaspora-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-RubyResolve-app',
          'gitlab-project-inspections-test/inspection-app',
          'redmine-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-app',
        ]"
        :aliases="[
          'Unresolved References Inspection (DI)',
          'All Inspections (DI)',
          'Unresolved References Inspection (GL)',
          'All Inspections (GL)',
          'Unresolved References Inspection (RM)',
          'All Inspections (RM)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages: Execution Time"
        measure="findUsages"
        :projects="['RUBY-23764-Case1/ruby-23764-findusages-case1', 'RUBY-23764-Case2/ruby-23764-findusages-case2']"
        :aliases="['Factory (GL)', 'Let Variable (GL)']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages: Quantity"
        measure="findUsages#number"
        :projects="['RUBY-23764-Case1/ruby-23764-findusages-case1', 'RUBY-23764-Case2/ruby-23764-findusages-case2']"
        :aliases="['Factory (GL)', 'Let Variable (GL)']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion (Diaspora)"
        measure="completion"
        :projects="['diaspora-project-test/completion/routes', 'diaspora-project-test/completion/exceptions', 'diaspora-project-test/completion/localization']"
        :aliases="['Routes', 'Exceptions', 'I18n#t']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion (GitLab)"
        measure="completion"
        :projects="['gitlab-project-test/completion/routes', 'gitlab-project-test/completion/exceptions', 'gitlab-project-test/completion/localization']"
        :aliases="['Routes', 'Exceptions', 'I18n#t']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion (Redmine)"
        measure="completion"
        :projects="['redmine-project-test/completion/routes', 'redmine-project-test/completion/exceptions', 'redmine-project-test/completion/localization']"
        :aliases="['Routes', 'Exceptions', 'I18n#t']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: Average AWT Delay"
        measure="test#average_awt_delay"
        :projects="[
          'RUBY-26170/typing',
          'RUBY-29334/typing',
          'GitLab/typing/typing/user/method',
          'GitLab/typing/typing/user/class',
          'GitLab/typing/typing/user/lambda',
          'GitLab/typing/typing/parser/method',
          'GitLab/typing/typing/parser/class',
          'GitLab/typing/typing/parser/class_array',
          'GitLab/typing/typing/parser/class_assoc',
          'GitLab/typing/typing/parser/newline_class_body',
          'GitLab/typing/typing/parser/newline_class_array',
          'GitLab/typing/typing/parser/newline_class_method',
        ]"
        :aliases="[
          'Ruby assoc with map',
          'RBS method',
          'User Model Method (GL)',
          'User Model Class (GL)',
          'User Model Lambda (GL)',
          'Parser Method',
          'Parser Class',
          'Parser Array',
          'Parser Assoc',
          'Parser Class (new line)',
          'Parser Array (new line)',
          'Parser Method (new line)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: Total Time"
        measure="typing"
        :projects="[
          'RUBY-26170/typing',
          'RUBY-29334/typing',
          'GitLab/typing/typing/user/method',
          'GitLab/typing/typing/user/class',
          'GitLab/typing/typing/user/lambda',
          'GitLab/typing/typing/parser/method',
          'GitLab/typing/typing/parser/class',
          'GitLab/typing/typing/parser/class_array',
          'GitLab/typing/typing/parser/class_assoc',
          'GitLab/typing/typing/parser/newline_class_body',
          'GitLab/typing/typing/parser/newline_class_array',
          'GitLab/typing/typing/parser/newline_class_method',
        ]"
        :aliases="[
          'Ruby assoc with map',
          'RBS method',
          'User Model Method (GL)',
          'User Model Class (GL)',
          'User Model Lambda (GL)',
          'Parser Method',
          'Parser Class',
          'Parser Array',
          'Parser Assoc',
          'Parser Class (new line)',
          'Parser Array (new line)',
          'Parser Method (new line)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Enter Handling: Average AWT Delay"
        measure="test#average_awt_delay"
        :projects="[
          'RUBY-29542/typing',
          'GitLab/typing/do_in_method',
          'GitLab/typing/method',
          'GitLab/typing/class',
          'GitLab/typing/lambda_body_in_class',
          'GitLab/typing/enter/parser/method',
          'GitLab/typing/enter/parser/class',
          'GitLab/typing/enter/parser/class_array',
          'GitLab/typing/enter/parser/class_assoc',
          'GitLab/typing/enter/structure/inside_query',
          'GitLab/typing/enter/structure/after_query',
          'GitLab/typing/enter/project_spec/describe',
          'GitLab/typing/enter/project_controller/class',
          'GitLab/typing/enter/mr_mail/class',
          'GitLab/typing/enter/user_show_view/before_div',
          'GitLab/typing/enter/routes_project/top',
          'GitLab/typing/enter/emojis_json/map',
        ]"
        :aliases="[
          'Do block in spec',
          'Do block in method',
          'Method body',
          'Class body',
          'Lambda body in class',
          'Ruby Parser Method',
          'Ruby Parser Class',
          'Ruby Parser Array',
          'Ruby Parser Assoc',
          'structure.sql, inside query (GL)',
          'structure.sql, after query (GL)',
          'Project Model Spec (GL)',
          'Project Controller (GL)',
          'MR Mail (GL)',
          'Users View Haml (GL)',
          'Project Routes (GL)',
          'Emojis.json (GL)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Enter Handling: Total Time"
        measure="typing"
        :projects="[
          'RUBY-29542/typing',
          'GitLab/typing/do_in_method',
          'GitLab/typing/method',
          'GitLab/typing/class',
          'GitLab/typing/lambda_body_in_class',
          'GitLab/typing/enter/parser/method',
          'GitLab/typing/enter/parser/class',
          'GitLab/typing/enter/parser/class_array',
          'GitLab/typing/enter/parser/class_assoc',
          'GitLab/typing/enter/structure/inside_query',
          'GitLab/typing/enter/structure/after_query',
          'GitLab/typing/enter/project_spec/describe',
          'GitLab/typing/enter/project_controller/class',
          'GitLab/typing/enter/mr_mail/class',
          'GitLab/typing/enter/user_show_view/before_div',
          'GitLab/typing/enter/routes_project/top',
          'GitLab/typing/enter/emojis_json/map',
        ]"
        :aliases="[
          'Do block in spec',
          'Do block in method',
          'Method body',
          'Class body',
          'Lambda body in class',
          'Ruby Parser Method',
          'Ruby Parser Class',
          'Ruby Parser Array',
          'Ruby Parser Assoc',
          'structure.sql, inside query (GL)',
          'structure.sql, after query (GL)',
          'Project Model Spec (GL)',
          'Project Controller (GL)',
          'MR Mail (GL)',
          'Users View Haml (GL)',
          'Project Routes (GL)',
          'Emojis.json (GL)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
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
        ]"
        :aliases="[
          'Project Model (GL)',
          'swagger_helper.rb',
          'activerecord-generated.rbs',
          'User Model Spec (GL)',
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
        ]"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Symbol Members: Execution Time"
        measure="getSymbolMembers"
        :projects="[
          'diaspora-project-test/getSymbolMembers-ApplicationController-false',
          'diaspora-project-test/getSymbolMembers-ApplicationController-true',
          'gitlab-project-test/getSymbolMembers-ApplicationController-false',
          'gitlab-project-test/getSymbolMembers-ApplicationController-true',
          'redmine-project-test/getSymbolMembers-ApplicationController-false',
          'redmine-project-test/getSymbolMembers-ApplicationController-true',
        ]"
        :aliases="[
          'ApplicationController (DI)',
          'ApplicationController (DI, no caches)',
          'ApplicationController (GL)',
          'ApplicationController (GL, no caches)',
          'ApplicationController (RM)',
          'ApplicationController (RM, no caches)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Symbol Members: Quantity"
        measure="getSymbolMembers#number"
        :projects="[
          'diaspora-project-test/getSymbolMembers-ApplicationController-false',
          'diaspora-project-test/getSymbolMembers-ApplicationController-true',
          'gitlab-project-test/getSymbolMembers-ApplicationController-false',
          'gitlab-project-test/getSymbolMembers-ApplicationController-true',
          'redmine-project-test/getSymbolMembers-ApplicationController-false',
          'redmine-project-test/getSymbolMembers-ApplicationController-true',
        ]"
        :aliases="[
          'ApplicationController (DI)',
          'ApplicationController (DI, no caches)',
          'ApplicationController (GL)',
          'ApplicationController (GL, no caches)',
          'ApplicationController (RM)',
          'ApplicationController (RM, no caches)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="GC Pause, ms"
        measure="gcPause"
        :projects="[
          'RUBY-23764-Case1/ruby-23764-findusages-case1',
          'RUBY-23764-Case2/ruby-23764-findusages-case2',
          'gitlab-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-app',
          'GitLab/typing/enter/project_spec/describe',
        ]"
        :aliases="[
          'Factory Find Usage (GL)',
          'Let Variable Find Usage (GL)',
          'All Inspections (GL)',
          'Unresolved References Inspection (GL)',
          'All Inspections (RM)',
          'Unresolved References Inspection (RM)',
          'Enter in Project Model Spec (GL)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="GC Memory Collected, Mb"
        measure="freedMemoryByGC"
        :projects="[
          'RUBY-23764-Case1/ruby-23764-findusages-case1',
          'RUBY-23764-Case2/ruby-23764-findusages-case2',
          'gitlab-project-inspections-test/inspection-app',
          'gitlab-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-RubyResolve-app',
          'redmine-project-inspections-test/inspection-app',
          'GitLab/typing/enter/project_spec/describe',
        ]"
        :aliases="[
          'Factory Find Usage (GL)',
          'Let Variable Find Usage (GL)',
          'All Inspections (GL)',
          'Unresolved References Inspection (GL)',
          'All Inspections (RM)',
          'Unresolved References Inspection (RM)',
          'Enter in Project Model Spec (GL)',
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
