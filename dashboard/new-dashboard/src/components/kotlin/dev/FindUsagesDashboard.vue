<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_findusage_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all project findUsages K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all project findUsages K2'"
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
            'intellij_commit/findUsages/UIAutomationInteractionModel_with_library_cache_k1',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Intellij findUsages mean value with Library cache K2"
          measure="findUsages#mean_value"
          :projects="[
            'intellij_commit/findUsages/loadModuleEntity_with_library_cache_k2',
            'intellij_commit/findUsages/setUp_with_library_cache_k2',
            'intellij_commit/findUsages/SolutionModel_with_library_cache_k2',
            'intellij_commit/findUsages/SqlBlock_with_library_cache_k2',
            'intellij_commit/findUsages/UIAutomationInteractionModel_with_library_cache_k2',
          ]"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
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
        />
      </div>
      <div class="flex-1 min-w-0">
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
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
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
            'toolbox_enterprise/findUsages/ROLE_ADMIN_with_library_cache_k1',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
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
            'toolbox_enterprise/findUsages/ROLE_ADMIN_with_library_cache_k2',
          ]"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Android canary leak' findUsages mean value with Library cache K1"
          measure="findUsages#mean_value"
          :projects="[
            'leak-canary-android/findUsages/BOOLEAN_with_library_cache_k1',
            'leak-canary-android/findUsages/HeapGraph_with_library_cache_k1',
            'leak-canary-android/findUsages/HeapObject_with_library_cache_k1',
            'leak-canary-android/findUsages/HprofRecordTag_with_library_cache_k1',
            'leak-canary-android/findUsages/INT_with_library_cache_k1',
            'leak-canary-android/findUsages/PrimitiveType_with_library_cache_k1',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="'Android canary leak' findUsages mean value with Library cache K2"
          measure="findUsages#mean_value"
          :projects="[
            'leak-canary-android/findUsages/BOOLEAN_with_library_cache_k2',
            'leak-canary-android/findUsages/HeapGraph_with_library_cache_k2',
            'leak-canary-android/findUsages/HeapObject_with_library_cache_k2',
            'leak-canary-android/findUsages/HprofRecordTag_with_library_cache_k2',
            'leak-canary-android/findUsages/INT_with_library_cache_k2',
            'leak-canary-android/findUsages/PrimitiveType_with_library_cache_k2',
          ]"
        />
      </div>
    </section>
  </DashboardPage>
  >
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
