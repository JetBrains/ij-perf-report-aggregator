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
        label="First Code Analysis (GitLab)"
        measure="firstCodeAnalysis#mean_value"
        :projects="[
          'GitLab/firstCodeAnalysis/app_models_user_rb',
          'GitLab/firstCodeAnalysis/app_models_project_rb',
          'GitLab/firstCodeAnalysis/db_structure_sql',
          'GitLab/firstCodeAnalysis/spec_models_project_spec_rb',
          'GitLab/firstCodeAnalysis/app_views_users_show_html_haml',
          'GitLab/firstCodeAnalysis/fixtures_emojis_index_json',
          'GitLab/firstCodeAnalysis/ruby27_parser_rb',
          'GitLab/firstCodeAnalysis/app_controllers_projects_controller_rb',
          'GitLab/firstCodeAnalysis/app_mailers_emails_merge_requests_rb',
          'GitLab/firstCodeAnalysis/config_routes_project_rb',
          'gitlab-project-test/firstCodeAnalysis/ee_app_graphql_mutations_boards_epic_boards_epic_move_list_rb',
          'gitlab-project-test/firstCodeAnalysis/ee_app_serializers_clusters_environment_entity_rb',
          'gitlab-project-test/firstCodeAnalysis/app_controllers_clusters_clusters_controller_rb',
          'gitlab-project-test/firstCodeAnalysis/app_views_admin_locale_html_haml',
          'gitlab-project-test/firstCodeAnalysis/app_controllers_admin_application_controller_rb',
          'gitlab-project-test/firstCodeAnalysis/app_models_ci_build_trace_chunk_rb',
        ]"
        :aliases="[
          'User Model',
          'Project Model',
          'structure.sql',
          'Project Spec',
          'Users View Haml',
          'Emojis JSON',
          'Ruby Parser',
          'Projects Controller',
          'MR Mail',
          'Routes Project',
          'Epic Move List',
          'Environment Entity',
          'Clusters Controller',
          'Locale Haml',
          'Admin App Controller',
          'Build Trace Chunk',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="First Code Analysis (Diaspora)"
        measure="firstCodeAnalysis#mean_value"
        :projects="[
          'diaspora-project-test/firstCodeAnalysis/app_models_conversation_visibility_rb',
          'diaspora-project-test/firstCodeAnalysis/spec_integration_api_contacts_controller_spec_rb',
          'diaspora-project-test/firstCodeAnalysis/app_models_message_rb',
          'diaspora-project-test/firstCodeAnalysis/app_views_admins_pods_html_haml',
          'diaspora-project-test/firstCodeAnalysis/app_controllers_admins_controller_rb',
          'diaspora-project-test/firstCodeAnalysis/app_workers_process_photo_rb',
        ]"
        :aliases="['Conversation Visibility', 'Contacts Controller Spec', 'Message', 'Admin Pods Haml', 'Admins Controller', 'Process Photo']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="First Code Analysis (Redmine)"
        measure="firstCodeAnalysis#mean_value"
        :projects="[
          'redmine-project-test/firstCodeAnalysis/app_controllers_auto_completes_controller_rb',
          'redmine-project-test/firstCodeAnalysis/app_views_admin_info_html_erb',
          'redmine-project-test/firstCodeAnalysis/app_models_time_entry_activity_rb',
          'redmine-project-test/firstCodeAnalysis/app_views_imports__time_entries_saved_objects_html_erb',
          'redmine-project-test/firstCodeAnalysis/app_controllers_account_controller_rb',
          'redmine-project-test/firstCodeAnalysis/app_controllers_application_controller_rb',
        ]"
        :aliases="['Auto Completes Controller', 'Admin Info Erb', 'Time Entry Activity', 'Time Entries Import Erb', 'Account Controller', 'Application Controller']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="First Code Analysis (Rest)"
        measure="firstCodeAnalysis#mean_value"
        :projects="[
          'RUBY-26170/firstCodeAnalysis/swagger_helper_rb',
          'RBSCollection/firstCodeAnalysis/gems_activerecord_6_0_activerecord-generated_rbs',
          'SampleRailsApp/firstCodeAnalysis/spec_models_user_model_spec_rb',
        ]"
        :aliases="['swagger_helper.rb (RUBY-26170)', 'activerecord-generated.rbs (RBSCollection)', 'User Model Spec (SampleRailsApp)']"
      />
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
        label="Completion Cold Cache (Diaspora)"
        measure="completion#mean_value"
        :projects="[
          'diaspora-project-test/completion/routes-cold-cache',
          'diaspora-project-test/completion/exceptions-cold-cache',
          'diaspora-project-test/completion/localization-cold-cache',
          'diaspora-project-test/completion/constant-cold-cache',
          'diaspora-project-test/completion/exceptions-prefix-cold-cache',
          'diaspora-project-test/completion/method-cold-cache',
          'diaspora-project-test/completion/qualified-cold-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion Cold Cache (GitLab)"
        measure="completion#mean_value"
        :projects="[
          'gitlab-project-test/completion/routes-cold-cache',
          'gitlab-project-test/completion/exceptions-cold-cache',
          'gitlab-project-test/completion/localization-cold-cache',
          'gitlab-project-test/completion/constant-cold-cache',
          'gitlab-project-test/completion/exceptions-prefix-cold-cache',
          'gitlab-project-test/completion/method-cold-cache',
          'gitlab-project-test/completion/qualified-cold-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion Cold Cache (Redmine)"
        measure="completion#mean_value"
        :projects="[
          'redmine-project-test/completion/routes-cold-cache',
          'redmine-project-test/completion/exceptions-cold-cache',
          'redmine-project-test/completion/localization-cold-cache',
          'redmine-project-test/completion/constant-cold-cache',
          'redmine-project-test/completion/exceptions-prefix-cold-cache',
          'redmine-project-test/completion/method-cold-cache',
          'redmine-project-test/completion/qualified-cold-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion Hot Cache (Diaspora)"
        measure="completion#mean_value"
        :projects="[
          'diaspora-project-test/completion/routes-hot-cache',
          'diaspora-project-test/completion/exceptions-hot-cache',
          'diaspora-project-test/completion/localization-hot-cache',
          'diaspora-project-test/completion/constant-hot-cache',
          'diaspora-project-test/completion/exceptions-prefix-hot-cache',
          'diaspora-project-test/completion/method-hot-cache',
          'diaspora-project-test/completion/qualified-hot-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion Hot Cache (GitLab)"
        measure="completion#mean_value"
        :projects="[
          'gitlab-project-test/completion/routes-hot-cache',
          'gitlab-project-test/completion/exceptions-hot-cache',
          'gitlab-project-test/completion/localization-hot-cache',
          'gitlab-project-test/completion/constant-hot-cache',
          'gitlab-project-test/completion/exceptions-prefix-hot-cache',
          'gitlab-project-test/completion/method-hot-cache',
          'gitlab-project-test/completion/qualified-hot-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion Hot Cache (Redmine)"
        measure="completion#mean_value"
        :projects="[
          'redmine-project-test/completion/routes-hot-cache',
          'redmine-project-test/completion/exceptions-hot-cache',
          'redmine-project-test/completion/localization-hot-cache',
          'redmine-project-test/completion/constant-hot-cache',
          'redmine-project-test/completion/exceptions-prefix-hot-cache',
          'redmine-project-test/completion/method-hot-cache',
          'redmine-project-test/completion/qualified-hot-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element Cold Cache (Diaspora)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'diaspora-project-test/completion/routes-cold-cache',
          'diaspora-project-test/completion/exceptions-cold-cache',
          'diaspora-project-test/completion/localization-cold-cache',
          'diaspora-project-test/completion/constant-cold-cache',
          'diaspora-project-test/completion/exceptions-prefix-cold-cache',
          'diaspora-project-test/completion/method-cold-cache',
          'diaspora-project-test/completion/qualified-cold-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element Cold Cache (GitLab)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'gitlab-project-test/completion/routes-cold-cache',
          'gitlab-project-test/completion/exceptions-cold-cache',
          'gitlab-project-test/completion/localization-cold-cache',
          'gitlab-project-test/completion/constant-cold-cache',
          'gitlab-project-test/completion/exceptions-prefix-cold-cache',
          'gitlab-project-test/completion/method-cold-cache',
          'gitlab-project-test/completion/qualified-cold-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element Cold Cache (Redmine)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'redmine-project-test/completion/routes-cold-cache',
          'redmine-project-test/completion/exceptions-cold-cache',
          'redmine-project-test/completion/localization-cold-cache',
          'redmine-project-test/completion/constant-cold-cache',
          'redmine-project-test/completion/exceptions-prefix-cold-cache',
          'redmine-project-test/completion/method-cold-cache',
          'redmine-project-test/completion/qualified-cold-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element Hot Cache (Diaspora)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'diaspora-project-test/completion/routes-hot-cache',
          'diaspora-project-test/completion/exceptions-hot-cache',
          'diaspora-project-test/completion/localization-hot-cache',
          'diaspora-project-test/completion/constant-hot-cache',
          'diaspora-project-test/completion/exceptions-prefix-hot-cache',
          'diaspora-project-test/completion/method-hot-cache',
          'diaspora-project-test/completion/qualified-hot-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element Hot Cache (GitLab)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'gitlab-project-test/completion/routes-hot-cache',
          'gitlab-project-test/completion/exceptions-hot-cache',
          'gitlab-project-test/completion/localization-hot-cache',
          'gitlab-project-test/completion/constant-hot-cache',
          'gitlab-project-test/completion/exceptions-prefix-hot-cache',
          'gitlab-project-test/completion/method-hot-cache',
          'gitlab-project-test/completion/qualified-hot-cache',
        ]"
        :aliases="['Routes', 'Exceptions', 'I18n#t', 'Constant', 'Exceptions (prefix)', 'Method', 'Qualified']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion First Element Hot Cache (Redmine)"
        measure="completion#firstElementShown#mean_value"
        :projects="[
          'redmine-project-test/completion/routes-hot-cache',
          'redmine-project-test/completion/exceptions-hot-cache',
          'redmine-project-test/completion/localization-hot-cache',
          'redmine-project-test/completion/constant-hot-cache',
          'redmine-project-test/completion/exceptions-prefix-hot-cache',
          'redmine-project-test/completion/method-hot-cache',
          'redmine-project-test/completion/qualified-hot-cache',
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
        label="Typing: Median Time"
        measure="typing#median_value"
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
        label="Enter Handling: Median Time"
        measure="typing#median_value"
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
        label="Symbol Members: Mean Execution Time"
        measure="getSymbolMembers#mean_value"
        :projects="[
          'diaspora-project-test/getSymbolMembers-ApplicationController-hot-cache',
          'diaspora-project-test/getSymbolMembers-ApplicationController-cold-cache',
          'gitlab-project-test/getSymbolMembers-ApplicationController-hot-cache',
          'gitlab-project-test/getSymbolMembers-ApplicationController-cold-cache',
          'redmine-project-test/getSymbolMembers-ApplicationController-hot-cache',
          'redmine-project-test/getSymbolMembers-ApplicationController-cold-cache',
        ]"
        :aliases="[
          'ApplicationController (DI, hot cache)',
          'ApplicationController (DI, cold cache)',
          'ApplicationController (GL, hot cache)',
          'ApplicationController (GL, cold cache)',
          'ApplicationController (RM, hot cache)',
          'ApplicationController (RM, cold cache)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Symbol Members: Quantity"
        measure="getSymbolMembers#number#mean_value"
        :projects="[
          'diaspora-project-test/getSymbolMembers-ApplicationController-hot-cache',
          'diaspora-project-test/getSymbolMembers-ApplicationController-cold-cache',
          'gitlab-project-test/getSymbolMembers-ApplicationController-hot-cache',
          'gitlab-project-test/getSymbolMembers-ApplicationController-cold-cache',
          'redmine-project-test/getSymbolMembers-ApplicationController-hot-cache',
          'redmine-project-test/getSymbolMembers-ApplicationController-cold-cache',
        ]"
        :aliases="[
          'ApplicationController (DI, hot cache)',
          'ApplicationController (DI, cold cache)',
          'ApplicationController (GL, hot cache)',
          'ApplicationController (GL, cold cache)',
          'ApplicationController (RM, hot cache)',
          'ApplicationController (RM, cold cache)',
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
