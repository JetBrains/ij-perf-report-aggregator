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
    <section
      v-for="chart in charts"
      :key="chart.label"
    >
      <GroupProjectsWithClientChart
        :label="chart.label"
        :measure="chart.measure"
        :projects="Object.keys(chart.series)"
        :aliases="Object.values(chart.series)"
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

// project → alias shown in the legend, in display order
type Series = Record<string, string>

interface Chart {
  label: string
  measure: string | string[]
  series: Series
}

const findUsagesSeries: Series = {
  "RUBY-23764-Case1/ruby-23764-findusages-case1": "Factory (GL)",
  "RUBY-23764-Case2/ruby-23764-findusages-case2": "Let Variable (GL)",
  "gitlab-find-usages/ruby-23764-findusages-case1": "Factory (GL)",
  "gitlab-find-usages/ruby-23764-findusages-case2": "Let Variable (GL)",
  "RUBY-32357/class": "Class (GL)",
  "RUBY-32357/module": "Module (GL)",
  "RUBY-32357/method": "Method (GL)",
  "RUBY-32357/singleton-method": "Singleton Method (GL)",
  "RUBY-32357/instance-variable": "Instance Variable (GL)",
  "RUBY-32357/class-variable": "Class Variable (GL)",
  "RUBY-32357/global-variable": "Global Variable (GL)",
  "RUBY-32357/delegate-method": "Delegate Method (GL)",
  "RUBY-32357/association": "Association (GL)",
  "gitlab-find-usages/class": "Class (GL)",
  "gitlab-find-usages/module": "Module (GL)",
  "gitlab-find-usages/method": "Method (GL)",
  "gitlab-find-usages/singleton-method": "Singleton Method (GL)",
  "gitlab-find-usages/instance-variable": "Instance Variable (GL)",
  "gitlab-find-usages/class-variable": "Class Variable (GL)",
  "gitlab-find-usages/global-variable": "Global Variable (GL)",
  "gitlab-find-usages/delegate-method": "Delegate Method (GL)",
  "gitlab-find-usages/association": "Association (GL)",
  "mastodon-find-usages/i18n-key": "I18n Key (MA)",
}

// suffix of `<project>/completion/<suffix>-<cold|hot>-cache` → alias
const completionCases: Series = {
  routes: "Routes",
  exceptions: "Exceptions",
  localization: "I18n#t",
  constant: "Constant",
  "exceptions-prefix": "Exceptions (prefix)",
  method: "Method",
  qualified: "Qualified",
}
function completionSeries(project: string, cache: "cold" | "hot"): Series {
  return Object.fromEntries(Object.entries(completionCases).map(([suffix, alias]) => [`${project}/completion/${suffix}-${cache}-cache`, alias]))
}

const typingSeries: Series = {
  "RUBY-26170/typing": "Ruby assoc with map",
  "RUBY-29334/typing": "RBS method",
  "GitLab/typing/typing/user/method": "User Model Method (GL)",
  "GitLab/typing/typing/user/class": "User Model Class (GL)",
  "GitLab/typing/typing/user/lambda": "User Model Lambda (GL)",
  "GitLab/typing/typing/parser/method": "Parser Method",
  "GitLab/typing/typing/parser/class": "Parser Class",
  "GitLab/typing/typing/parser/class_array": "Parser Array",
  "GitLab/typing/typing/parser/class_assoc": "Parser Assoc",
  "GitLab/typing/typing/parser/newline_class_body": "Parser Class (new line)",
  "GitLab/typing/typing/parser/newline_class_array": "Parser Array (new line)",
  "GitLab/typing/typing/parser/newline_class_method": "Parser Method (new line)",
}

const enterSeries: Series = {
  "RUBY-29542/typing": "Do block in spec",
  "GitLab/typing/do_in_method": "Do block in method",
  "GitLab/typing/method": "Method body",
  "GitLab/typing/class": "Class body",
  "GitLab/typing/lambda_body_in_class": "Lambda body in class",
  "GitLab/typing/enter/parser/method": "Ruby Parser Method",
  "GitLab/typing/enter/parser/class": "Ruby Parser Class",
  "GitLab/typing/enter/parser/class_array": "Ruby Parser Array",
  "GitLab/typing/enter/parser/class_assoc": "Ruby Parser Assoc",
  "GitLab/typing/enter/parser/start_file": "Ruby Parser Start File",
  "GitLab/typing/enter/parser/top_level_comment": "Ruby Parser Top Level Comment",
  "GitLab/typing/enter/structure/inside_query": "structure.sql, inside query (GL)",
  "GitLab/typing/enter/structure/after_query": "structure.sql, after query (GL)",
  "GitLab/typing/enter/project_spec/describe": "Project Model Spec (GL)",
  "GitLab/typing/enter/project_controller/class": "Project Controller (GL)",
  "GitLab/typing/enter/mr_mail/class": "MR Mail (GL)",
  "GitLab/typing/enter/user_show_view/before_div": "Users View Haml (GL)",
  "GitLab/typing/enter/routes_project/top": "Project Routes (GL)",
  "GitLab/typing/enter/emojis_json/map": "Emojis.json (GL)",
}

const symbolMembersSeries: Series = {
  "diaspora-project-test/getSymbolMembers-ApplicationController-hot-cache": "ApplicationController (DI, hot cache)",
  "diaspora-project-test/getSymbolMembers-ApplicationController-cold-cache": "ApplicationController (DI, cold cache)",
  "gitlab-project-test/getSymbolMembers-ApplicationController-hot-cache": "ApplicationController (GL, hot cache)",
  "gitlab-project-test/getSymbolMembers-ApplicationController-cold-cache": "ApplicationController (GL, cold cache)",
  "redmine-project-test/getSymbolMembers-ApplicationController-hot-cache": "ApplicationController (RM, hot cache)",
  "redmine-project-test/getSymbolMembers-ApplicationController-cold-cache": "ApplicationController (RM, cold cache)",
}

const gcSeries: Series = {
  "RUBY-23764-Case1/ruby-23764-findusages-case1": "Factory Find Usage (GL)",
  "RUBY-23764-Case2/ruby-23764-findusages-case2": "Let Variable Find Usage (GL)",
  "gitlab-find-usages/ruby-23764-findusages-case1": "Factory Find Usage (GL)",
  "gitlab-find-usages/ruby-23764-findusages-case2": "Let Variable Find Usage (GL)",
  "RUBY-32357/class": "Class Find Usage (GL)",
  "RUBY-32357/module": "Module Find Usage (GL)",
  "RUBY-32357/method": "Method Find Usage (GL)",
  "RUBY-32357/singleton-method": "Singleton Method Find Usage (GL)",
  "RUBY-32357/instance-variable": "Instance Variable Find Usage (GL)",
  "RUBY-32357/class-variable": "Class Variable Find Usage (GL)",
  "RUBY-32357/global-variable": "Global Variable Find Usage (GL)",
  "RUBY-32357/delegate-method": "Delegate Method Find Usage (GL)",
  "RUBY-32357/association": "Association Find Usage (GL)",
  "gitlab-find-usages/class": "Class Find Usage (GL)",
  "gitlab-find-usages/module": "Module Find Usage (GL)",
  "gitlab-find-usages/method": "Method Find Usage (GL)",
  "gitlab-find-usages/singleton-method": "Singleton Method Find Usage (GL)",
  "gitlab-find-usages/instance-variable": "Instance Variable Find Usage (GL)",
  "gitlab-find-usages/class-variable": "Class Variable Find Usage (GL)",
  "gitlab-find-usages/global-variable": "Global Variable Find Usage (GL)",
  "gitlab-find-usages/delegate-method": "Delegate Method Find Usage (GL)",
  "gitlab-find-usages/association": "Association Find Usage (GL)",
  "mastodon-find-usages/i18n-key": "I18n Key (MA)",
  "GitLab/typing/enter/project_spec/describe": "Enter in Project Model Spec (GL)",
}

const charts: Chart[] = [
  {
    label: "First Code Analysis (GitLab)",
    measure: "firstCodeAnalysis#mean_value",
    series: {
      "GitLab/firstCodeAnalysis/app_models_user_rb": "User Model",
      "GitLab/firstCodeAnalysis/app_models_project_rb": "Project Model",
      "GitLab/firstCodeAnalysis/db_structure_sql": "structure.sql",
      "GitLab/firstCodeAnalysis/spec_models_project_spec_rb": "Project Spec",
      "GitLab/firstCodeAnalysis/app_views_users_show_html_haml": "Users View Haml",
      "GitLab/firstCodeAnalysis/fixtures_emojis_index_json": "Emojis JSON",
      "GitLab/firstCodeAnalysis/ruby27_parser_rb": "Ruby Parser",
      "GitLab/firstCodeAnalysis/app_controllers_projects_controller_rb": "Projects Controller",
      "GitLab/firstCodeAnalysis/app_mailers_emails_merge_requests_rb": "MR Mail",
      "GitLab/firstCodeAnalysis/config_routes_project_rb": "Routes Project",
      "gitlab-project-test/firstCodeAnalysis/ee_app_graphql_mutations_boards_epic_boards_epic_move_list_rb": "Epic Move List",
      "gitlab-project-test/firstCodeAnalysis/ee_app_serializers_clusters_environment_entity_rb": "Environment Entity",
      "gitlab-project-test/firstCodeAnalysis/app_controllers_clusters_clusters_controller_rb": "Clusters Controller",
      "gitlab-project-test/firstCodeAnalysis/app_views_admin_locale_html_haml": "Locale Haml",
      "gitlab-project-test/firstCodeAnalysis/app_controllers_admin_application_controller_rb": "Admin App Controller",
      "gitlab-project-test/firstCodeAnalysis/app_models_ci_build_trace_chunk_rb": "Build Trace Chunk",
    },
  },
  {
    label: "First Code Analysis (Diaspora)",
    measure: "firstCodeAnalysis#mean_value",
    series: {
      "diaspora-project-test/firstCodeAnalysis/app_models_conversation_visibility_rb": "Conversation Visibility",
      "diaspora-project-test/firstCodeAnalysis/spec_integration_api_contacts_controller_spec_rb": "Contacts Controller Spec",
      "diaspora-project-test/firstCodeAnalysis/app_models_message_rb": "Message",
      "diaspora-project-test/firstCodeAnalysis/app_views_admins_pods_html_haml": "Admin Pods Haml",
      "diaspora-project-test/firstCodeAnalysis/app_controllers_admins_controller_rb": "Admins Controller",
      "diaspora-project-test/firstCodeAnalysis/app_workers_process_photo_rb": "Process Photo",
    },
  },
  {
    label: "First Code Analysis (Redmine)",
    measure: "firstCodeAnalysis#mean_value",
    series: {
      "redmine-project-test/firstCodeAnalysis/app_controllers_auto_completes_controller_rb": "Auto Completes Controller",
      "redmine-project-test/firstCodeAnalysis/app_views_admin_info_html_erb": "Admin Info Erb",
      "redmine-project-test/firstCodeAnalysis/app_models_time_entry_activity_rb": "Time Entry Activity",
      "redmine-project-test/firstCodeAnalysis/app_views_imports__time_entries_saved_objects_html_erb": "Time Entries Import Erb",
      "redmine-project-test/firstCodeAnalysis/app_controllers_account_controller_rb": "Account Controller",
      "redmine-project-test/firstCodeAnalysis/app_controllers_application_controller_rb": "Application Controller",
    },
  },
  {
    label: "First Code Analysis (Rest)",
    measure: "firstCodeAnalysis#mean_value",
    series: {
      "RUBY-26170/firstCodeAnalysis/swagger_helper_rb": "swagger_helper.rb (RUBY-26170)",
      "RBSCollection/firstCodeAnalysis/gems_activerecord_6_0_activerecord-generated_rbs": "activerecord-generated.rbs (RBSCollection)",
      "SampleRailsApp/firstCodeAnalysis/spec_models_user_model_spec_rb": "User Model Spec (SampleRailsApp)",
    },
  },
  { label: "Find Usages: Execution Time", measure: ["findUsages", "findUsagesInToolWindow"], series: findUsagesSeries },
  { label: "Find Usages: Quantity", measure: ["findUsages#number", "findUsagesInToolWindow#number"], series: findUsagesSeries },
  { label: "Completion Cold Cache (Diaspora)", measure: "completion#mean_value", series: completionSeries("diaspora-project-test", "cold") },
  { label: "Completion Cold Cache (GitLab)", measure: "completion#mean_value", series: completionSeries("gitlab-project-test", "cold") },
  { label: "Completion Cold Cache (Redmine)", measure: "completion#mean_value", series: completionSeries("redmine-project-test", "cold") },
  { label: "Completion Hot Cache (Diaspora)", measure: "completion#mean_value", series: completionSeries("diaspora-project-test", "hot") },
  { label: "Completion Hot Cache (GitLab)", measure: "completion#mean_value", series: completionSeries("gitlab-project-test", "hot") },
  { label: "Completion Hot Cache (Redmine)", measure: "completion#mean_value", series: completionSeries("redmine-project-test", "hot") },
  { label: "Completion First Element Cold Cache (Diaspora)", measure: "completion#firstElementShown#mean_value", series: completionSeries("diaspora-project-test", "cold") },
  { label: "Completion First Element Cold Cache (GitLab)", measure: "completion#firstElementShown#mean_value", series: completionSeries("gitlab-project-test", "cold") },
  { label: "Completion First Element Cold Cache (Redmine)", measure: "completion#firstElementShown#mean_value", series: completionSeries("redmine-project-test", "cold") },
  { label: "Completion First Element Hot Cache (Diaspora)", measure: "completion#firstElementShown#mean_value", series: completionSeries("diaspora-project-test", "hot") },
  { label: "Completion First Element Hot Cache (GitLab)", measure: "completion#firstElementShown#mean_value", series: completionSeries("gitlab-project-test", "hot") },
  { label: "Completion First Element Hot Cache (Redmine)", measure: "completion#firstElementShown#mean_value", series: completionSeries("redmine-project-test", "hot") },
  { label: "Typing: Average AWT Delay", measure: "test#average_awt_delay", series: typingSeries },
  { label: "Typing: Median Time", measure: "typing#median_value", series: typingSeries },
  { label: "Enter Handling: Average AWT Delay", measure: "test#average_awt_delay", series: enterSeries },
  { label: "Enter Handling: Median Time", measure: "typing#median_value", series: enterSeries },
  { label: "Symbol Members: Mean Execution Time", measure: "getSymbolMembers#mean_value", series: symbolMembersSeries },
  { label: "Symbol Members: Quantity", measure: "getSymbolMembers#number#mean_value", series: symbolMembersSeries },
  { label: "GC Pause, ms", measure: "gcPause", series: gcSeries },
  { label: "GC Memory Collected, Mb", measure: "freedMemoryByGC", series: gcSeries },
]
</script>
