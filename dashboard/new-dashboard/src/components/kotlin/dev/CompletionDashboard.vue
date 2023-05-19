<template>
  <DashboardPage
    v-slot="{serverConfigurator, dashboardConfigurators, averagesConfigurators, warnings}"
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_completion_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all project completion K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all project completion K2'"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Hello-world' completion mean value on  K1"
          measure="completion#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k1', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k1']"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Hello-world' completion mean value on hello-world K2"
          measure="completion#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k2', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k2']"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Hello-world' firstElementShown mean value on hello-world K1"
          measure="completion#firstElementShown#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k1', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k1']"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Hello-world' firstElementShown mean value on hello-world K2"
          measure="completion#firstElementShown#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k2', 'kotlin_empty/completion/empty_place_typing_with_library_cache_k2']"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server ' completion mean value K1"
          measure="completion#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k1',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Kotlin language server ' completion mean value K2"
          measure="completion#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k2',
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
          label="'Kotlin language server' firstElementShown mean value K1"
          measure="completion#firstElementShown#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k1',
          ]"
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
          :accidents="warnings"
        />
      </div>
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"

</script>