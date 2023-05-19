<template>
  <DashboardPage
    v-slot="{ dashboardConfigurators, averagesConfigurators}"
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-6">
      <div class="flex-1  min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'completion K1'"
        />
      </div>
      <div class="flex-1  min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'completion K2'"
        />
      </div>
      <div class="flex-1  min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'highlighting\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'highlighting K1'"
        />
      </div>
      <div class="flex-1  min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'highlighting\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'highlighting K2'"
        />
      </div>
      <div class="flex-1  min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'findUsage K1'"
        />
      </div>
      <div class="flex-1  min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'findUsage K2'"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1  min-w-0">
        <GroupProjectsChart
          label="'Hello-world' completion mean value on hello-world K1"
          measure="completion#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k1', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k1']"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Hello-world' completion mean value on hello-world K2"
          measure="completion#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k2', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k2']"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Hello-world' firstElementShown mean value on hello-world K1"
          measure="completion#firstElementShown#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k1', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k1']"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Hello-world' firstElementShown mean value on hello-world K2"
          measure="completion#firstElementShown#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k2', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k2']"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij completion mean value K1"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/empty_place_with_library_cache_k1',
            'intellij_commit/completion/after_parameter_with_library_cache_k1',
            'intellij_commit/completion/empty_place_typing_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij completion mean value K2"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/empty_place_with_library_cache_k2',
            'intellij_commit/completion/after_parameter_with_library_cache_k2',
            'intellij_commit/completion/empty_place_typing_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij firstElementShown mean value K1"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'intellij_commit/completion/empty_place_with_library_cache_k1',
            'intellij_commit/completion/after_parameter_with_library_cache_k1',
            'intellij_commit/completion/empty_place_typing_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij firstElementShown mean value K2"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'intellij_commit/completion/empty_place_with_library_cache_k2',
            'intellij_commit/completion/after_parameter_with_library_cache_k2',
            'intellij_commit/completion/empty_place_typing_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suit 2' completion mean value K1"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_typing_with_library_cache_k1',
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_with_library_cache_k1',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_typing_with_library_cache_k1',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suit 2' completion mean value K2"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_typing_with_library_cache_k2',
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_with_library_cache_k2',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_typing_with_library_cache_k2',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suit 2' firstElementShown mean value K1"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_typing_with_library_cache_k1',
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_with_library_cache_k1',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_typing_with_library_cache_k1',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suit 2' firstElementShown mean value K2"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_typing_with_library_cache_k2',
            'intellij_commit/completion/KotlinHighLevelFunctionParameterInfoHandler_emptyPlace_updateUIOrFail_with_library_cache_k2',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_typing_with_library_cache_k2',
            'intellij_commit/completion/KtOCSwiftChangeSignatureTest_emptyPlace_changeReturnType_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij with typing suit 2' completion mean value K1"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/IdeMenuBar_emptyPlace_sout_typing_with_library_cache_k1',
            'intellij_commit/completion/TestModelParser_emptyPlace_if_typing_with_library_cache_k1',
            'intellij_commit/completion/AndroidModuleSystem_emptyPlace_get_typing_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij with typing suit 2' completion mean value K2"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/IdeMenuBar_emptyPlace_sout_typing_with_library_cache_k2',
            'intellij_commit/completion/TestModelParser_emptyPlace_if_typing_with_library_cache_k2',
            'intellij_commit/completion/AndroidModuleSystem_emptyPlace_get_typing_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij with typing suit 2' firstElementShown mean value K1"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'intellij_commit/completion/IdeMenuBar_emptyPlace_sout_typing_with_library_cache_k1',
            'intellij_commit/completion/TestModelParser_emptyPlace_if_typing_with_library_cache_k1',
            'intellij_commit/completion/AndroidModuleSystem_emptyPlace_get_typing_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij with typing suit 2' firstElementShown mean value K2"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'intellij_commit/completion/IdeMenuBar_emptyPlace_sout_typing_with_library_cache_k2',
            'intellij_commit/completion/TestModelParser_emptyPlace_if_typing_with_library_cache_k2',
            'intellij_commit/completion/AndroidModuleSystem_emptyPlace_get_typing_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Kotlin lang Completion mean value K1"
          measure="completion#mean_value"
          :projects="[
            'kotlin_lang/completion/after_parameter_with_library_cache_k1',
            'kotlin_lang/completion/empty_place_with_library_cache_k1',
            'kotlin_lang/completion/empty_place_typing_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Kotlin lang Completion mean value K2"
          measure="completion#mean_value"
          :projects="[
            'kotlin_lang/completion/after_parameter_with_library_cache_k2',
            'kotlin_lang/completion/empty_place_with_library_cache_k2',
            'kotlin_lang/completion/empty_place_typing_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Kotlin lang firstElementShown mean value K1"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'kotlin_lang/completion/after_parameter_with_library_cache_k1',
            'kotlin_lang/completion/empty_place_with_library_cache_k1',
            'kotlin_lang/completion/empty_place_typing_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Kotlin lang firstElementShown mean value K2"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'kotlin_lang/completion/after_parameter_with_library_cache_k2',
            'kotlin_lang/completion/empty_place_with_library_cache_k2',
            'kotlin_lang/completion/empty_place_typing_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' completion mean value K1"
          measure="completion#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k1',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' completion mean value K2"
          measure="completion#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k2',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' firstElementShown mean value K1"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k1',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' firstElementShown mean value K2"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k2',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin empty' highlight mean value with Library cache K1"
          measure="highlighting#mean_value"
          :projects="[
            'kotlin_empty/highlight/Main_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin empty' highlight mean value with Library cache K2"
          measure="highlighting#mean_value"
          :projects="[
            'kotlin_empty/highlight/Main_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij highlight mean value with Library cache K1"
          measure="highlighting#mean_value"
          :projects="[
            'intellij_commit/highlight/KtOCSwiftChangeSignatureTest_with_library_cache_k1',
            'intellij_commit/highlight/KotlinHighLevelFunctionParameterInfoHandler_with_library_cache_k1',
            'intellij_commit/highlight/ContentManagerImpl_with_library_cache_k1',
            'intellij_commit/highlight/JdkList_with_library_cache_k1',
            'intellij_commit/highlight/ComposeCompletionContributorTest_with_library_cache_k1',
            'intellij_commit/highlight/AgpUpgradeRefactoringProcessor_with_library_cache_k1',
            'intellij_commit/highlight/AndroidModelTest_with_library_cache_k1',
            'intellij_commit/highlight/SecureWireOverStreamTransport_with_library_cache_k1',
            'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k1',
            'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k1',
            'intellij_commit/highlight/OraIntrospector_with_library_cache_k1',
            'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k1',
            'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k1',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij highlight mean value with Library cache K2"
          measure="highlighting#mean_value"
          :projects="[
            'intellij_commit/highlight/KtOCSwiftChangeSignatureTest_with_library_cache_k2',
            'intellij_commit/highlight/KotlinHighLevelFunctionParameterInfoHandler_with_library_cache_k2',
            'intellij_commit/highlight/ContentManagerImpl_with_library_cache_k2',
            'intellij_commit/highlight/JdkList_with_library_cache_k2',
            'intellij_commit/highlight/ComposeCompletionContributorTest_with_library_cache_k2',
            'intellij_commit/highlight/AgpUpgradeRefactoringProcessor_with_library_cache_k2',
            'intellij_commit/highlight/AndroidModelTest_with_library_cache_k2',
            'intellij_commit/highlight/SecureWireOverStreamTransport_with_library_cache_k2',
            'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k2',
            'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k2',
            'intellij_commit/highlight/OraIntrospector_with_library_cache_k2',
            'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k2',
            'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k2',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suite 2' highlight mean value with Library cache K1"
          measure="highlighting#mean_value"
          :projects="[
            'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k1',
            'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k1',
            'intellij_commit/highlight/OraIntrospector_with_library_cache_k1',
            'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k1',
            'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k1',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suite 2' highlight mean value with Library cache K2"
          measure="highlighting#mean_value"
          :projects="[
            'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k2',
            'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k2',
            'intellij_commit/highlight/OraIntrospector_with_library_cache_k2',
            'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k2',
            'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k2',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin lang' highlight mean value with Library cache K1"
          measure="highlighting#mean_value"
          :projects="[
            'kotlin_lang/highlight/CommonParser_with_library_cache_k1',
            'kotlin_lang/highlight/FirErrors_with_library_cache_k1',
            'kotlin_lang/highlight/Flag_with_library_cache_k1',
            'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k1',
            'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k1',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin lang' highlight mean value with Library cache K2"
          measure="highlighting#mean_value"
          :projects="[
            'kotlin_lang/highlight/CommonParser_with_library_cache_k2',
            'kotlin_lang/highlight/FirErrors_with_library_cache_k2',
            'kotlin_lang/highlight/Flag_with_library_cache_k2',
            'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k2',
            'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k2',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' highlight mean value with Library cache K1"
          measure="highlighting#mean_value"
          :projects="[
            'kotlin_language_server/highlight/Compiler_with_library_cache_k1',
            'kotlin_language_server/highlight/Completions_with_library_cache_k1',
            'kotlin_language_server/highlight/CompletionsTest_with_library_cache_k1',
            'kotlin_language_server/highlight/JavaElementConverter_with_library_cache_k1',
            'kotlin_language_server/highlight/KotlinTextDocumentService_with_library_cache_k1',
            'kotlin_language_server/highlight/QuickFixesTest_with_library_cache_k1',
            'kotlin_language_server/highlight/SourcePath_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' highlight mean value with Library cache K2"
          measure="highlighting#mean_value"
          :projects="[
            'kotlin_language_server/highlight/Compiler_with_library_cache_k2',
            'kotlin_language_server/highlight/Completions_with_library_cache_k2',
            'kotlin_language_server/highlight/CompletionsTest_with_library_cache_k2',
            'kotlin_language_server/highlight/JavaElementConverter_with_library_cache_k2',
            'kotlin_language_server/highlight/KotlinTextDocumentService_with_library_cache_k2',
            'kotlin_language_server/highlight/QuickFixesTest_with_library_cache_k2',
            'kotlin_language_server/highlight/SourcePath_with_library_cache_k2',


          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'TBE' highlight mean value with Library cache K1"
          measure="highlighting#mean_value"
          :projects="[
            'toolbox_enterprise/highlight/IdeSettingControllerTest_with_library_cache_k1',
            'toolbox_enterprise/highlight/IntelliJPluginSettingTest_with_library_cache_k1',
            'toolbox_enterprise/highlight/LoginTests_with_library_cache_k1',
            'toolbox_enterprise/highlight/PluginAuditLogService_with_library_cache_k1',
            'toolbox_enterprise/highlight/PluginControllerTest_with_library_cache_k1',
            'toolbox_enterprise/highlight/ProfileController_with_library_cache_k1',
            'toolbox_enterprise/highlight/ProfileService_with_library_cache_k1',
            'toolbox_enterprise/highlight/ProfileServiceTest_with_library_cache_k1',
            'toolbox_enterprise/highlight/SecurityTests_with_library_cache_k1',
            'toolbox_enterprise/highlight/UsageDataFlowTests_with_library_cache_k1',
            'toolbox_enterprise/highlight/VmOptionSettingTest_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'TBE' highlight mean value with Library cache K2"
          measure="highlighting#mean_value"
          :projects="[
            'toolbox_enterprise/highlight/IdeSettingControllerTest_with_library_cache_k2',
            'toolbox_enterprise/highlight/IntelliJPluginSettingTest_with_library_cache_k2',
            'toolbox_enterprise/highlight/LoginTests_with_library_cache_k2',
            'toolbox_enterprise/highlight/PluginAuditLogService_with_library_cache_k2',
            'toolbox_enterprise/highlight/PluginControllerTest_with_library_cache_k2',
            'toolbox_enterprise/highlight/ProfileController_with_library_cache_k2',
            'toolbox_enterprise/highlight/ProfileService_with_library_cache_k2',
            'toolbox_enterprise/highlight/ProfileServiceTest_with_library_cache_k2',
            'toolbox_enterprise/highlight/SecurityTests_with_library_cache_k2',
            'toolbox_enterprise/highlight/UsageDataFlowTests_with_library_cache_k2',
            'toolbox_enterprise/highlight/VmOptionSettingTest_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Ktor samples' highlight mean value with Library cache K1"
          measure="highlighting#mean_value"
          :projects="[
            'ktor_samples/highlight/ApplicationTest_with_library_cache_k1',
            'ktor_samples/highlight/HttpBinApplication_with_library_cache_k1',
            'ktor_samples/highlight/Upload_with_library_cache_k1',
            'ktor_samples/highlight/Videos_with_library_cache_k1',
            'ktor_samples/highlight/LocationHeaderApplication_with_library_cache_k1',
            'ktor_samples/highlight/ReverseProxyApplication_with_library_cache_k1',
            'ktor_samples/highlight/SseApplication_with_library_cache_k1',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Ktor samples' highlight mean value with Library cache K2"
          measure="highlighting#mean_value"
          :projects="[
            'ktor_samples/highlight/ApplicationTest_with_library_cache_k2',
            'ktor_samples/highlight/HttpBinApplication_with_library_cache_k2',
            'ktor_samples/highlight/Upload_with_library_cache_k2',
            'ktor_samples/highlight/Videos_with_library_cache_k2',
            'ktor_samples/highlight/LocationHeaderApplication_with_library_cache_k2',
            'ktor_samples/highlight/ReverseProxyApplication_with_library_cache_k2',
            'ktor_samples/highlight/SseApplication_with_library_cache_k2',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij findUsages mean value with Library cache K1"
          measure="findUsages#mean_value"
          :projects="[
            'intellij_commit/findUsages/loadModuleEntity_with_library_cache_k1',
            'intellij_commit/findUsages/setUp_with_library_cache_k1',
            'intellij_commit/findUsages/SolutionModel_with_library_cache_k1',
            'intellij_commit/findUsages/SqlBlock_with_library_cache_k1',
            'intellij_commit/findUsages/UIAutomationInteractionModel_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij findUsages mean value with Library cache K2"
          measure="findUsages#mean_value"
          :projects="[
            'intellij_commit/findUsages/loadModuleEntity_with_library_cache_k2',
            'intellij_commit/findUsages/setUp_with_library_cache_k2',
            'intellij_commit/findUsages/SolutionModel_with_library_cache_k2',
            'intellij_commit/findUsages/SqlBlock_with_library_cache_k2',
            'intellij_commit/findUsages/UIAutomationInteractionModel_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin lang' findUsages mean value with Library cache K1"
          measure="findUsages#mean_value"
          :projects="[
            'kotlin_lang/findUsages/CommonParser_with_library_cache_k1',
            'kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache_k1',
            'kotlin_lang/findUsages/FirErrors_with_library_cache_k1',
            'kotlin_lang/findUsages/Flag_with_library_cache_k1',
            'kotlin_lang/findUsages/ReferenceSymbolTable_with_library_cache_k1',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin lang' findUsages mean value with Library cache K2"
          measure="findUsages#mean_value"
          :projects="[
            'kotlin_lang/findUsages/CommonParser_with_library_cache_k2',
            'kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache_k2',
            'kotlin_lang/findUsages/FirErrors_with_library_cache_k2',
            'kotlin_lang/findUsages/Flag_with_library_cache_k2',
            'kotlin_lang/findUsages/ReferenceSymbolTable_with_library_cache_k2',
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>
        
    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'TBE' findUsages mean value with Library cache K1"
          measure="findUsages#mean_value"
          :projects="[
            'toolbox_enterprise/findUsages/ErrorReport_with_library_cache_k1',
            'toolbox_enterprise/findUsages/genUuid_with_library_cache_k1',
            'toolbox_enterprise/findUsages/getAll_with_library_cache_k1',
            'toolbox_enterprise/findUsages/getTempDirectory_with_library_cache_k1',
            'toolbox_enterprise/findUsages/PrincipalContext_with_library_cache_k1',
            'toolbox_enterprise/findUsages/RequestMapping_with_library_cache_k1',
            'toolbox_enterprise/findUsages/RestController_with_library_cache_k1',
            'toolbox_enterprise/findUsages/ROLE_ADMIN_with_library_cache_k1'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'TBE' findUsages mean value with Library cache K2"
          measure="findUsages#mean_value"
          :projects="[
            'toolbox_enterprise/findUsages/ErrorReport_with_library_cache_k2',
            'toolbox_enterprise/findUsages/genUuid_with_library_cache_k2',
            'toolbox_enterprise/findUsages/getAll_with_library_cache_k2',
            'toolbox_enterprise/findUsages/getTempDirectory_with_library_cache_k2',
            'toolbox_enterprise/findUsages/PrincipalContext_with_library_cache_k2',
            'toolbox_enterprise/findUsages/RequestMapping_with_library_cache_k2',
            'toolbox_enterprise/findUsages/RestController_with_library_cache_k2',
            'toolbox_enterprise/findUsages/ROLE_ADMIN_with_library_cache_k2'
          ]"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'PerformInlineRename mean value on  K1"
          measure="performInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k1', 'kotlin_language_server/insertCode/Rename_renameSymbol_k1',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k1']"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="PerformInlineRename mean value on  K2"
          measure="performInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k2', 'kotlin_language_server/insertCode/Rename_renameSymbol_k2',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k2']"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'StartInlineRename mean value on  K1"
          measure="startInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k1', 'kotlin_language_server/insertCode/Rename_renameSymbol_k1',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k1']"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="StartInlineRename mean value on  K2"
          measure="startInlineRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k2', 'kotlin_language_server/insertCode/Rename_renameSymbol_k2',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k2']"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'PrepareForRename mean value on  K1"
          measure="prepareForRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k1', 'kotlin_language_server/insertCode/Rename_renameSymbol_k1',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k1']"
          :configurators="dashboardConfigurators"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="PrepareForRename mean value on  K2"
          measure="prepareForRename#mean_value"
          :projects="['intellij_commit/rename/SqlBlock_SqlBlockRenamed_k2', 'kotlin_language_server/insertCode/Rename_renameSymbol_k2',
                      'kotlin_language_server/insertCode/SpecialJavaFileForTest_j2k_k2']"
          :configurators="dashboardConfigurators"
        />
      </div>
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

</script>