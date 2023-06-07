<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_completion_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'completion K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'completion K2'"
        />
      </div>
    </section>
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
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion mean value Kotlin language server K1"
          measure="completion#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k1',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k1',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k1',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion mean value Kotlin language server K2"
          measure="completion#mean_value"
          :projects="[
            'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k2',
            'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k2',
            'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k2',
          ]"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
