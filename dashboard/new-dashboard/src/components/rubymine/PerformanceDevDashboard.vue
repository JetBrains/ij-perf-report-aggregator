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
        label="Find Usages: Execution Time"
        :measure="['findUsages', 'findUsagesInToolWindow']"
        :projects="[
          'RUBY-23764-Case1/ruby-23764-findusages-case1',
          'RUBY-23764-Case2/ruby-23764-findusages-case2',
          'gitlab-find-usages/ruby-23764-findusages-case1',
          'gitlab-find-usages/ruby-23764-findusages-case2',
          'RUBY-32357/class',
          'RUBY-32357/module',
          'RUBY-32357/method',
          'RUBY-32357/singleton-method',
          'RUBY-32357/instance-variable',
          'RUBY-32357/class-variable',
          'RUBY-32357/global-variable',
          'RUBY-32357/delegate-method',
          'RUBY-32357/association',
          'gitlab-find-usages/class',
          'gitlab-find-usages/module',
          'gitlab-find-usages/method',
          'gitlab-find-usages/singleton-method',
          'gitlab-find-usages/instance-variable',
          'gitlab-find-usages/class-variable',
          'gitlab-find-usages/global-variable',
          'gitlab-find-usages/delegate-method',
          'gitlab-find-usages/association',
          'mastodon-find-usages/i18n-key',
        ]"
        :aliases="[
          'Factory (GL)',
          'Let Variable (GL)',
          'Factory (GL)',
          'Let Variable (GL)',
          'Class (GL)',
          'Module (GL)',
          'Method (GL)',
          'Singleton Method (GL)',
          'Instance Variable (GL)',
          'Class Variable (GL)',
          'Global Variable (GL)',
          'Delegate Method (GL)',
          'Association (GL)',
          'Class (GL)',
          'Module (GL)',
          'Method (GL)',
          'Singleton Method (GL)',
          'Instance Variable (GL)',
          'Class Variable (GL)',
          'Global Variable (GL)',
          'Delegate Method (GL)',
          'Association (GL)',
          'I18n Key (MA)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Find Usages: Quantity"
        :measure="['findUsages#number', 'findUsagesInToolWindow#number']"
        :projects="[
          'RUBY-23764-Case1/ruby-23764-findusages-case1',
          'RUBY-23764-Case2/ruby-23764-findusages-case2',
          'gitlab-find-usages/ruby-23764-findusages-case1',
          'gitlab-find-usages/ruby-23764-findusages-case2',
          'RUBY-32357/class',
          'RUBY-32357/module',
          'RUBY-32357/method',
          'RUBY-32357/singleton-method',
          'RUBY-32357/instance-variable',
          'RUBY-32357/class-variable',
          'RUBY-32357/global-variable',
          'RUBY-32357/delegate-method',
          'RUBY-32357/association',
          'gitlab-find-usages/class',
          'gitlab-find-usages/module',
          'gitlab-find-usages/method',
          'gitlab-find-usages/singleton-method',
          'gitlab-find-usages/instance-variable',
          'gitlab-find-usages/class-variable',
          'gitlab-find-usages/global-variable',
          'gitlab-find-usages/delegate-method',
          'gitlab-find-usages/association',
          'mastodon-find-usages/i18n-key',
        ]"
        :aliases="[
          'Factory (GL)',
          'Let Variable (GL)',
          'Factory (GL)',
          'Let Variable (GL)',
          'Class (GL)',
          'Module (GL)',
          'Method (GL)',
          'Singleton Method (GL)',
          'Instance Variable (GL)',
          'Class Variable (GL)',
          'Global Variable (GL)',
          'Delegate Method (GL)',
          'Association (GL)',
          'Class (GL)',
          'Module (GL)',
          'Method (GL)',
          'Singleton Method (GL)',
          'Instance Variable (GL)',
          'Class Variable (GL)',
          'Global Variable (GL)',
          'Delegate Method (GL)',
          'Association (GL)',
          'I18n Key (MA)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion (Diaspora)"
        measure="completion#mean_value"
        :projects="[
          'diaspora-project-test/completion/routes',
          'diaspora-project-test/completion/exceptions',
          'diaspora-project-test/completion/localization',
          'diaspora-project-test/completion/constant',
          'diaspora-project-test/completion/exceptions-prefix',
          'diaspora-project-test/completion/method',
          'diaspora-project-test/completion/qualified',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion (GitLab)"
        measure="completion#mean_value"
        :projects="[
          'gitlab-project-test/completion/routes',
          'gitlab-project-test/completion/exceptions',
          'gitlab-project-test/completion/localization',
          'gitlab-project-test/completion/constant',
          'gitlab-project-test/completion/exceptions-prefix',
          'gitlab-project-test/completion/method',
          'gitlab-project-test/completion/qualified',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion (Redmine)"
        measure="completion#mean_value"
        :projects="[
          'redmine-project-test/completion/routes',
          'redmine-project-test/completion/exceptions',
          'redmine-project-test/completion/localization',
          'redmine-project-test/completion/constant',
          'redmine-project-test/completion/exceptions-prefix',
          'redmine-project-test/completion/method',
          'redmine-project-test/completion/qualified',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element (Diaspora)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'diaspora-project-test/completion/routes',
          'diaspora-project-test/completion/exceptions',
          'diaspora-project-test/completion/localization',
          'diaspora-project-test/completion/constant',
          'diaspora-project-test/completion/exceptions-prefix',
          'diaspora-project-test/completion/method',
          'diaspora-project-test/completion/qualified',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element (GitLab)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'gitlab-project-test/completion/routes',
          'gitlab-project-test/completion/exceptions',
          'gitlab-project-test/completion/localization',
          'gitlab-project-test/completion/constant',
          'gitlab-project-test/completion/exceptions-prefix',
          'gitlab-project-test/completion/method',
          'gitlab-project-test/completion/qualified',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element (Redmine)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'redmine-project-test/completion/routes',
          'redmine-project-test/completion/exceptions',
          'redmine-project-test/completion/localization',
          'redmine-project-test/completion/constant',
          'redmine-project-test/completion/exceptions-prefix',
          'redmine-project-test/completion/method',
          'redmine-project-test/completion/qualified',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
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
      <GroupProjectsWithClientChart
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
      <GroupProjectsWithClientChart
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
          'GitLab/typing/enter/parser/start_file',
          'GitLab/typing/enter/parser/top_level_comment',
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
          'Ruby Parser Start File',
          'Ruby Parser Top Level Comment',
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
      <GroupProjectsWithClientChart
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
          'GitLab/typing/enter/parser/start_file',
          'GitLab/typing/enter/parser/top_level_comment',
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
          'Ruby Parser Start File',
          'Ruby Parser Top Level Comment',
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
      <GroupProjectsWithClientChart
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
      <GroupProjectsWithClientChart
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
      <GroupProjectsWithClientChart
        label="GC Pause, ms"
        measure="gcPause"
        :projects="[
          'RUBY-23764-Case1/ruby-23764-findusages-case1',
          'RUBY-23764-Case2/ruby-23764-findusages-case2',
          'gitlab-find-usages/ruby-23764-findusages-case1',
          'gitlab-find-usages/ruby-23764-findusages-case2',
          'RUBY-32357/class',
          'RUBY-32357/module',
          'RUBY-32357/method',
          'RUBY-32357/singleton-method',
          'RUBY-32357/instance-variable',
          'RUBY-32357/class-variable',
          'RUBY-32357/global-variable',
          'RUBY-32357/delegate-method',
          'RUBY-32357/association',
          'gitlab-find-usages/class',
          'gitlab-find-usages/module',
          'gitlab-find-usages/method',
          'gitlab-find-usages/singleton-method',
          'gitlab-find-usages/instance-variable',
          'gitlab-find-usages/class-variable',
          'gitlab-find-usages/global-variable',
          'gitlab-find-usages/delegate-method',
          'gitlab-find-usages/association',
          'mastodon-find-usages/i18n-key',
          'GitLab/typing/enter/project_spec/describe',
        ]"
        :aliases="[
          'Factory Find Usage (GL)',
          'Let Variable Find Usage (GL)',
          'Factory Find Usage (GL)',
          'Let Variable Find Usage (GL)',
          'Class Find Usage (GL)',
          'Module Find Usage (GL)',
          'Method Find Usage (GL)',
          'Singleton Method Find Usage (GL)',
          'Instance Variable Find Usage (GL)',
          'Class Variable Find Usage (GL)',
          'Global Variable Find Usage (GL)',
          'Delegate Method Find Usage (GL)',
          'Association Find Usage (GL)',
          'Class Find Usage (GL)',
          'Module Find Usage (GL)',
          'Method Find Usage (GL)',
          'Singleton Method Find Usage (GL)',
          'Instance Variable Find Usage (GL)',
          'Class Variable Find Usage (GL)',
          'Global Variable Find Usage (GL)',
          'Delegate Method Find Usage (GL)',
          'Association Find Usage (GL)',
          'I18n Key (MA)',
          'Enter in Project Model Spec (GL)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="GC Memory Collected, Mb"
        measure="freedMemoryByGC"
        :projects="[
          'RUBY-23764-Case1/ruby-23764-findusages-case1',
          'RUBY-23764-Case2/ruby-23764-findusages-case2',
          'gitlab-find-usages/ruby-23764-findusages-case1',
          'gitlab-find-usages/ruby-23764-findusages-case2',
          'RUBY-32357/class',
          'RUBY-32357/module',
          'RUBY-32357/method',
          'RUBY-32357/singleton-method',
          'RUBY-32357/instance-variable',
          'RUBY-32357/class-variable',
          'RUBY-32357/global-variable',
          'RUBY-32357/delegate-method',
          'RUBY-32357/association',
          'gitlab-find-usages/class',
          'gitlab-find-usages/module',
          'gitlab-find-usages/method',
          'gitlab-find-usages/singleton-method',
          'gitlab-find-usages/instance-variable',
          'gitlab-find-usages/class-variable',
          'gitlab-find-usages/global-variable',
          'gitlab-find-usages/delegate-method',
          'gitlab-find-usages/association',
          'mastodon-find-usages/i18n-key',
          'GitLab/typing/enter/project_spec/describe',
        ]"
        :aliases="[
          'Factory Find Usage (GL)',
          'Let Variable Find Usage (GL)',
          'Factory Find Usage (GL)',
          'Let Variable Find Usage (GL)',
          'Class Find Usage (GL)',
          'Module Find Usage (GL)',
          'Method Find Usage (GL)',
          'Singleton Method Find Usage (GL)',
          'Instance Variable Find Usage (GL)',
          'Class Variable Find Usage (GL)',
          'Global Variable Find Usage (GL)',
          'Delegate Method Find Usage (GL)',
          'Association Find Usage (GL)',
          'Class Find Usage (GL)',
          'Module Find Usage (GL)',
          'Method Find Usage (GL)',
          'Singleton Method Find Usage (GL)',
          'Instance Variable Find Usage (GL)',
          'Class Variable Find Usage (GL)',
          'Global Variable Find Usage (GL)',
          'Delegate Method Find Usage (GL)',
          'Association Find Usage (GL)',
          'I18n Key (MA)',
          'Enter in Project Model Spec (GL)',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
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
