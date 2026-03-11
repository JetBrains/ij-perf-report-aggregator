<template>
  <DashboardPage
    db-name="perfintDev"
    table="ruby"
    persistent-id="ruby_product_dashboard"
    :with-installer="false"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing Time"],
    measures: ["indexingTimeWithoutPauses"],
    projects: [
      "diaspora-project-test/indexing",
      "gem-rbs-collection-indexing-test/indexing",
      "gitlab-project-test/indexing",
      "redmine-project-test/indexing",
      "markus-project-test/indexing",
      "mastodon-project-test/indexing",
      "rubygems-org-project-test/indexing",
    ],
    aliases: ["Diaspora", "RBS Collection", "Gitlab", "Redmine", "Markus", "Mastodon", "RubyGems.org"],
  },
  {
    labels: ["First Code Analysis (GitLab)"],
    measures: ["firstCodeAnalysis#mean_value"],
    projects: [
      "GitLab/firstCodeAnalysis/app_models_user_rb",
      "GitLab/firstCodeAnalysis/app_models_project_rb",
      "GitLab/firstCodeAnalysis/db_structure_sql",
      "GitLab/firstCodeAnalysis/spec_models_project_spec_rb",
      "GitLab/firstCodeAnalysis/app_views_users_show_html_haml",
      "GitLab/firstCodeAnalysis/fixtures_emojis_index_json",
      "GitLab/firstCodeAnalysis/ruby27_parser_rb",
      "GitLab/firstCodeAnalysis/app_controllers_projects_controller_rb",
      "GitLab/firstCodeAnalysis/app_mailers_emails_merge_requests_rb",
      "GitLab/firstCodeAnalysis/config_routes_project_rb",
      "gitlab-project-test/firstCodeAnalysis/ee_app_graphql_mutations_boards_epic_boards_epic_move_list_rb",
      "gitlab-project-test/firstCodeAnalysis/ee_app_serializers_clusters_environment_entity_rb",
      "gitlab-project-test/firstCodeAnalysis/app_controllers_clusters_clusters_controller_rb",
      "gitlab-project-test/firstCodeAnalysis/app_views_admin_locale_html_haml",
      "gitlab-project-test/firstCodeAnalysis/app_controllers_admin_application_controller_rb",
      "gitlab-project-test/firstCodeAnalysis/app_models_ci_build_trace_chunk_rb",
    ],
    aliases: [
      "User Model",
      "Project Model",
      "structure.sql",
      "Project Spec",
      "Users View Haml",
      "Emojis JSON",
      "Ruby Parser",
      "Projects Controller",
      "MR Mail",
      "Routes Project",
      "Epic Move List",
      "Environment Entity",
      "Clusters Controller",
      "Locale Haml",
      "Admin App Controller",
      "Build Trace Chunk",
    ],
  },
  {
    labels: ["First Code Analysis (Diaspora)"],
    measures: ["firstCodeAnalysis#mean_value"],
    projects: [
      "diaspora-project-test/firstCodeAnalysis/app_models_conversation_visibility_rb",
      "diaspora-project-test/firstCodeAnalysis/spec_integration_api_contacts_controller_spec_rb",
      "diaspora-project-test/firstCodeAnalysis/app_models_message_rb",
      "diaspora-project-test/firstCodeAnalysis/app_views_admins_pods_html_haml",
      "diaspora-project-test/firstCodeAnalysis/app_controllers_admins_controller_rb",
      "diaspora-project-test/firstCodeAnalysis/app_workers_process_photo_rb",
    ],
    aliases: ["Conversation Visibility", "Contacts Controller Spec", "Message", "Admin Pods Haml", "Admins Controller", "Process Photo"],
  },
  {
    labels: ["First Code Analysis (Redmine)"],
    measures: ["firstCodeAnalysis#mean_value"],
    projects: [
      "redmine-project-test/firstCodeAnalysis/app_controllers_auto_completes_controller_rb",
      "redmine-project-test/firstCodeAnalysis/app_views_admin_info_html_erb",
      "redmine-project-test/firstCodeAnalysis/app_models_time_entry_activity_rb",
      "redmine-project-test/firstCodeAnalysis/app_views_imports__time_entries_saved_objects_html_erb",
      "redmine-project-test/firstCodeAnalysis/app_controllers_account_controller_rb",
      "redmine-project-test/firstCodeAnalysis/app_controllers_application_controller_rb",
    ],
    aliases: ["Auto Completes Controller", "Admin Info Erb", "Time Entry Activity", "Time Entries Import Erb", "Account Controller", "Application Controller"],
  },
  {
    labels: ["First Code Analysis (Rest)"],
    measures: ["firstCodeAnalysis#mean_value"],
    projects: [
      "RUBY-26170/firstCodeAnalysis/swagger_helper_rb",
      "RBSCollection/firstCodeAnalysis/gems_activerecord_6_0_activerecord-generated_rbs",
      "SampleRailsApp/firstCodeAnalysis/spec_models_user_model_spec_rb",
    ],
    aliases: ["swagger_helper.rb (RUBY-26170)", "activerecord-generated.rbs (RBSCollection)", "User Model Spec (SampleRailsApp)"],
  },
  {
    labels: ["Completion (Diaspora)"],
    measures: ["completion#mean_value"],
    projects: [
      "diaspora-project-test/completion/routes",
      "diaspora-project-test/completion/exceptions",
      "diaspora-project-test/completion/localization",
      "diaspora-project-test/completion/constant",
      "diaspora-project-test/completion/exceptions-prefix",
      "diaspora-project-test/completion/method",
      "diaspora-project-test/completion/qualified",
    ],
    aliases: ["Routes", "Exceptions", "I18n#t", "Constant", "Exceptions (prefix)", "Method", "Qualified"],
  },
  {
    labels: ["Completion (GitLab)"],
    measures: ["completion#mean_value"],
    projects: [
      "gitlab-project-test/completion/routes",
      "gitlab-project-test/completion/exceptions",
      "gitlab-project-test/completion/localization",
      "gitlab-project-test/completion/constant",
      "gitlab-project-test/completion/exceptions-prefix",
      "gitlab-project-test/completion/method",
      "gitlab-project-test/completion/qualified",
    ],
    aliases: ["Routes", "Exceptions", "I18n#t", "Constant", "Exceptions (prefix)", "Method", "Qualified"],
  },
  {
    labels: ["Completion (Redmine)"],
    measures: ["completion#mean_value"],
    projects: [
      "redmine-project-test/completion/routes",
      "redmine-project-test/completion/exceptions",
      "redmine-project-test/completion/localization",
      "redmine-project-test/completion/constant",
      "redmine-project-test/completion/exceptions-prefix",
      "redmine-project-test/completion/method",
      "redmine-project-test/completion/qualified",
    ],
    aliases: ["Routes", "Exceptions", "I18n#t", "Constant", "Exceptions (prefix)", "Method", "Qualified"],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: [],
  },
  {
    labels: ["Typing: Average AWT Delay"],
    measures: ["test#average_awt_delay"],
    projects: [
      "RUBY-26170/typing",
      "RUBY-29334/typing",
      "GitLab/typing/typing/user/method",
      "GitLab/typing/typing/user/class",
      "GitLab/typing/typing/user/lambda",
      "GitLab/typing/typing/parser/method",
      "GitLab/typing/typing/parser/class",
      "GitLab/typing/typing/parser/class_array",
      "GitLab/typing/typing/parser/class_assoc",
      "GitLab/typing/typing/parser/newline_class_body",
      "GitLab/typing/typing/parser/newline_class_array",
      "GitLab/typing/typing/parser/newline_class_method",
    ],
    aliases: [
      "Ruby assoc with map",
      "RBS method",
      "User Model Method (GL)",
      "User Model Class (GL)",
      "User Model Lambda (GL)",
      "Parser Method",
      "Parser Class",
      "Parser Array",
      "Parser Assoc",
      "Parser Class (new line)",
      "Parser Array (new line)",
      "Parser Method (new line)",
    ],
  },
  {
    labels: ["Typing: Total Time"],
    measures: ["typing"],
    projects: [
      "RUBY-26170/typing",
      "RUBY-29334/typing",
      "GitLab/typing/typing/user/method",
      "GitLab/typing/typing/user/class",
      "GitLab/typing/typing/user/lambda",
      "GitLab/typing/typing/parser/method",
      "GitLab/typing/typing/parser/class",
      "GitLab/typing/typing/parser/class_array",
      "GitLab/typing/typing/parser/class_assoc",
      "GitLab/typing/typing/parser/newline_class_body",
      "GitLab/typing/typing/parser/newline_class_array",
      "GitLab/typing/typing/parser/newline_class_method",
    ],
    aliases: [
      "Ruby assoc with map",
      "RBS method",
      "User Model Method (GL)",
      "User Model Class (GL)",
      "User Model Lambda (GL)",
      "Parser Method",
      "Parser Class",
      "Parser Array",
      "Parser Assoc",
      "Parser Class (new line)",
      "Parser Array (new line)",
      "Parser Method (new line)",
    ],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: [
      "diaspora-project-inspections-test/inspection-RubyResolve-app",
      "diaspora-project-inspections-test/inspection-app",
      "gitlab-project-inspections-test/inspection-RubyResolve-app",
      "gitlab-project-inspections-test/inspection-app",
      "gitlab-project-inspections-test/inspection-App-RubyResolve",
      "gitlab-project-inspections-test/inspection-App",
      "redmine-project-inspections-test/inspection-RubyResolve-app",
      "redmine-project-inspections-test/inspection-app",
      "mastodon-project-inspections-test/inspection-RubyResolve-app",
      "mastodon-project-inspections-test/inspection-app",
    ],
    aliases: [
      "Unresolved References Inspection (DI)",
      "All Inspections (DI)",
      "Unresolved References Inspection (GL)",
      "All Inspections (GL)",
      "Unresolved References Inspection (GL)",
      "All Inspections (GL)",
      "Unresolved References Inspection (RM)",
      "All Inspections (RM)",
      "Unresolved References Inspection (MN)",
      "All Inspections (MN)",
    ],
  },
  {
    labels: ["Gitlab Inspections"],
    measures: ["globalInspections"],
    projects: [
      "gitlab-project-inspections-test/inspection-App",
      "gitlab-project-inspections-test/inspection-App-RubyResolve",
      "gitlab-project-inspections-test/inspection-Yaml",
      "gitlab-project-inspections-test/inspection-XML",
      "gitlab-project-inspections-test/inspection-WebStorm-AppSpecFrontend",
      "gitlab-project-inspections-test/inspection-Slim",
      "gitlab-project-inspections-test/inspection-RubyMine-App",
      "gitlab-project-inspections-test/inspection-Others",
      "gitlab-project-inspections-test/inspection-Markdown",
      "gitlab-project-inspections-test/inspection-Liquid",
      "gitlab-project-inspections-test/inspection-Haml",
      "gitlab-project-inspections-test/inspection-Erb",
      "gitlab-project-inspections-test/inspection-DataGrip",
    ],
    aliases: [
      "All (app/)",
      "Unresolved Ruby References (app/)",
      "Yaml",
      "XML",
      "All on WebStorm files (app/ & spec/frontend)",
      "Slim",
      "All on RubyMine files (app/)",
      "All on all other files",
      "Markdown",
      "Liquid",
      "Haml",
      "Erb",
      "All on DataGrip files",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
