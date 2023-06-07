<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing K1"
          measure="indexing"
          :projects="['kotlin_empty/indexing_k1', 'intellij_commit/indexing_k1', 'kotlin_lang/indexing_k1']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing K2"
          measure="indexing"
          :projects="['kotlin_empty/indexing_k2', 'intellij_commit/indexing_k2', 'kotlin_lang/indexing_k2']"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion mean value on hello-world"
          measure="completion#mean_value"
          :projects="['kotlin_empty/completion/empty_place_with_library_cache_k1', 'kotlin_empty/completion/empty_place_with_library_cache_k2']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion mean value K1"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/empty_place_with_library_cache_k1',
            'intellij_commit/completion/after_parameter_with_library_cache_k1',
            'kotlin_lang/completion/after_parameter_with_library_cache_k1',
            'kotlin_lang/completion/empty_place_with_library_cache_k1',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion mean value K2"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/empty_place_with_library_cache_k2',
            'intellij_commit/completion/after_parameter_with_library_cache_k2',
            'kotlin_lang/completion/after_parameter_with_library_cache_k2',
            'kotlin_lang/completion/empty_place_with_library_cache_k2',
          ]"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion mean value with typing"
          measure="completion#mean_value"
          :projects="[
            'intellij_commit/completion/empty_place_typing_with_library_cache_k1',
            'intellij_commit/completion/empty_place_typing_with_library_cache_k2',
            'kotlin_empty/completion/empty_place_typing_with_library_cache_k1',
            'kotlin_empty/completion/empty_place_typing_with_library_cache_k2',
            'kotlin_lang/completion/empty_place_typing_with_library_cache_k1',
            'kotlin_lang/completion/empty_place_typing_with_library_cache_k2',
          ]"
        />
      </div>
    </section>

    <section>
      <GroupProjectsChart
        label="Highlight mean value with Library cache K1 on intellij"
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
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Highlight mean value with Library cache K2 on intellij"
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
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Highlight mean value with Library cache K1"
        measure="highlighting#mean_value"
        :projects="[
          'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k1',
          'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k1',
          'intellij_commit/highlight/OraIntrospector_with_library_cache_k1',
          'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k1',
          'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k1',
          'kotlin_lang/highlight/CommonParser_with_library_cache_k1',
          'kotlin_lang/highlight/FirErrors_with_library_cache_k1',
          'kotlin_lang/highlight/Flag_with_library_cache_k1',
          'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k1',
          'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k1',
        ]"
      />
      <GroupProjectsChart
        label="Highlight mean value with Library cache K2"
        measure="highlighting#mean_value"
        :projects="[
          'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k2',
          'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k2',
          'intellij_commit/highlight/OraIntrospector_with_library_cache_k2',
          'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k2',
          'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k2',
          'kotlin_lang/highlight/CommonParser_with_library_cache_k2',
          'kotlin_lang/highlight/FirErrors_with_library_cache_k2',
          'kotlin_lang/highlight/Flag_with_library_cache_k2',
          'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k2',
          'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k2',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="FindUsages mean value with Library cache K1"
        measure="findUsages#mean_value"
        :projects="[
          'intellij_commit/findUsages/loadModuleEntity_with_library_cache_k1',
          'intellij_commit/findUsages/SqlBlock_with_library_cache_k1',
          'intellij_commit/findUsages/OraIntrospector_with_library_cache_k1',
          'intellij_commit/findUsages/UIAutomationInteractionModel.Generated_with_library_cache_k1',
          'intellij_commit/findUsages/SolutionModel.Generated_with_library_cache_k1',
          'kotlin_lang/findUsages/FirErrors_with_library_cache_k1',
          'kotlin_lang/findUsages/Flag_with_library_cache_k1',
          'kotlin_lang/findUsages/CommonParser_with_library_cache_k1',
          'kotlin_lang/findUsages/SymbolTable_with_library_cache_k1',
          'kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache_k1',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="FindUsages mean value with Library cache K2"
        measure="findUsages#mean_value"
        :projects="[
          'intellij_commit/findUsages/loadModuleEntity_with_library_cache_k2',
          'intellij_commit/findUsages/SqlBlock_with_library_cache_k2',
          'intellij_commit/findUsages/OraIntrospector_with_library_cache_k2',
          'intellij_commit/findUsages/UIAutomationInteractionModel.Generated_with_library_cache_k2',
          'intellij_commit/findUsages/SolutionModel.Generated_with_library_cache_k2',
          'kotlin_lang/findUsages/FirErrors_with_library_cache_k2',
          'kotlin_lang/findUsages/Flag_with_library_cache_k2',
          'kotlin_lang/findUsages/CommonParser_with_library_cache_k2',
          'kotlin_lang/findUsages/SymbolTable_with_library_cache_k2',
          'kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache_k2',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
