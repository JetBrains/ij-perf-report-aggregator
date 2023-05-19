<template>
  <DashboardPage
    v-slot="{serverConfigurator, dashboardConfigurators, warnings}"
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_inspection_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin empty' localInspections mean value with Library cache K1"
          measure="localInspections#mean_value"
          :projects="[
            'kotlin_empty/highlight/Main_with_library_cache_k1'
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin empty' localInspections mean value with Library cache K2"
          measure="localInspections#mean_value"
          :projects="[
            'kotlin_empty/highlight/Main_with_library_cache_k2'
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij localInspections mean value with Library cache K1"
          measure="localInspections#mean_value"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij localInspections mean value with Library cache K2"
          measure="localInspections#mean_value"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suite 2' localInspections mean value with Library cache K1"
          measure="localInspections#mean_value"
          :projects="[
            'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k1',
            'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k1',
            'intellij_commit/highlight/OraIntrospector_with_library_cache_k1',
            'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k1',
            'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k1',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Intellij suite 2' localInspections mean value with Library cache K2"
          measure="localInspections#mean_value"
          :projects="[
            'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k2',
            'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k2',
            'intellij_commit/highlight/OraIntrospector_with_library_cache_k2',
            'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k2',
            'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k2',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin lang' localInspections mean value with Library cache K1"
          measure="localInspections#mean_value"
          :projects="[
            'kotlin_lang/highlight/CommonParser_with_library_cache_k1',
            'kotlin_lang/highlight/FirErrors_with_library_cache_k1',
            'kotlin_lang/highlight/Flag_with_library_cache_k1',
            'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k1',
            'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k1',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin lang' localInspections mean value with Library cache K2"
          measure="localInspections#mean_value"
          :projects="[
            'kotlin_lang/highlight/CommonParser_with_library_cache_k2',
            'kotlin_lang/highlight/FirErrors_with_library_cache_k2',
            'kotlin_lang/highlight/Flag_with_library_cache_k2',
            'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k2',
            'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k2',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' localInspections mean value with Library cache K1"
          measure="localInspections#mean_value"
          :projects="[
            'kotlin_language_server/highlight/Compiler_with_library_cache_k1',
            'kotlin_language_server/highlight/Completions_with_library_cache_k1',
            'kotlin_language_server/highlight/CompletionsTest_with_library_cache_k1',
            'kotlin_language_server/highlight/JavaElementConverter_with_library_cache_k1',
            'kotlin_language_server/highlight/KotlinTextDocumentService_with_library_cache_k1',
            'kotlin_language_server/highlight/QuickFixesTest_with_library_cache_k1',
            'kotlin_language_server/highlight/SourcePath_with_library_cache_k1'
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server' localInspections mean value with Library cache K2"
          measure="localInspections#mean_value"
          :projects="[
            'kotlin_language_server/highlight/Compiler_with_library_cache_k2',
            'kotlin_language_server/highlight/Completions_with_library_cache_k2',
            'kotlin_language_server/highlight/CompletionsTest_with_library_cache_k2',
            'kotlin_language_server/highlight/JavaElementConverter_with_library_cache_k2',
            'kotlin_language_server/highlight/KotlinTextDocumentService_with_library_cache_k2',
            'kotlin_language_server/highlight/QuickFixesTest_with_library_cache_k2',
            'kotlin_language_server/highlight/SourcePath_with_library_cache_k2',


          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'TBE' localInspections mean value with Library cache K1"
          measure="localInspections#mean_value"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'TBE' localInspections mean value with Library cache K2"
          measure="localInspections#mean_value"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class=" flex-1 min-w-0">
        <GroupProjectsChart
          label="'Ktor samples' localInspections mean value with Library cache K1"
          measure="localInspections#mean_value"
          :projects="[
            'ktor_samples/highlight/ApplicationTest_with_library_cache_k1',
            'ktor_samples/highlight/HttpBinApplication_with_library_cache_k1',
            'ktor_samples/highlight/Upload_with_library_cache_k1',
            'ktor_samples/highlight/Videos_with_library_cache_k1',
            'ktor_samples/highlight/LocationHeaderApplication_with_library_cache_k1',
            'ktor_samples/highlight/ReverseProxyApplication_with_library_cache_k1',
            'ktor_samples/highlight/SseApplication_with_library_cache_k1',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Ktor samples' localInspections mean value with Library cache K2"
          measure="localInspections#mean_value"
          :projects="[
            'ktor_samples/highlight/ApplicationTest_with_library_cache_k2',
            'ktor_samples/highlight/HttpBinApplication_with_library_cache_k2',
            'ktor_samples/highlight/Upload_with_library_cache_k2',
            'ktor_samples/highlight/Videos_with_library_cache_k2',
            'ktor_samples/highlight/LocationHeaderApplication_with_library_cache_k2',
            'ktor_samples/highlight/ReverseProxyApplication_with_library_cache_k2',
            'ktor_samples/highlight/SseApplication_with_library_cache_k2',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

</script>