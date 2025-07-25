<template>
  <DashboardPage
    db-name="perfintDev"
    :with-installer="false"
    table="webstorm"
    persistent-id="webstorm_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Code Vision (JSReferencesCodeVisionProvider)"
        measure="JSReferencesCodeVisionProvider"
        :projects="[
          'aws_cdk/localInspection/logging',
          'WEB_5976/localInspection/react_mui',
          'toh-pt6/localInspection/hero.service.ts',
          'vue3-admin-vite/localInspection/index.vue',
          'eslint-plugin-jest/localInspection/misc.ts',
          'allure-js/localInspection/JasmineAllureReporter.ts',
          'material-ui-react-admin/localInspection/PostEdit.tsx',
        ]"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Inspections"
        measure="globalInspections"
        :projects="['allure-js/inspection', 'axios/inspection', 'gitlab/inspection-app']"
      />
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing Average Responsiveness"
          measure="test#average_awt_delay"
          :projects="[
            'axios/typing/axios',
            'vue-template/typing/vue-template',
            'react-todo-js/typing/react-todo',
            'toh-pt6/typing/toh-pt6',
            'eslint-plugin-jest/typing/eslintPluginJest',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing Responsiveness"
          measure="test#max_awt_delay"
          :projects="[
            'axios/typing/axios',
            'vue-template/typing/vue-template',
            'react-todo-js/typing/react-todo',
            'toh-pt6/typing/toh-pt6',
            'eslint-plugin-jest/typing/eslintPluginJest',
          ]"
        />
      </div>
    </section>

    <Divider title="TypeScript" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['eslint-plugin-jest/completion/types']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['aws_cdk/localInspection/logging', 'eslint-plugin-jest/localInspection/misc.ts', 'allure-js/localInspection/JasmineAllureReporter.ts']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['eslint-plugin-jest/typing/eslintPluginJest']"
        />
      </div>
    </section>

    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="['aws_cdk/indexing', 'angular/indexing', 'eslint-plugin-jest/indexing', 'dxos/indexing', 'allure-js/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="['aws_cdk/indexing', 'angular/indexing', 'eslint-plugin-jest/indexing', 'dxos/indexing', 'allure-js/indexing']"
        />
      </div>
    </section>

    <Divider title="JavaScript" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['axios/completion/functions']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['axios/localInspection/utils.js']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['axios/typing/axios']"
        />
      </div>
    </section>

    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="['axios/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="['axios/indexing']"
        />
      </div>
    </section>

    <Divider title="Angular" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['toh-pt6/completion/attribute', 'toh-pt6/completion/component']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['toh-pt6/localInspection/hero.service.ts', 'toh-pt6/localInspection/heroes.component.html']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['toh-pt6/typing/toh-pt6']"
        />
      </div>
    </section>

    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="['toh-pt6/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="['toh-pt6/indexing']"
        />
      </div>
    </section>

    <Divider title="React" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="[
            'react-todo-js/completion/attribute',
            'react-todo-js/completion/component',
            'vkui/completion/component',
            'ring-ui/completion/component',
            'material-ui-react-admin/completion/attribute',
          ]"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['react-todo-js/localInspection/App.js', 'WEB_5976/localInspection/react_mui', 'material-ui-react-admin/localInspection/PostEdit.tsx']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['react-todo-js/typing/react-todo']"
        />
      </div>
    </section>

    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="['react-todo-js/indexing', 'ring-ui/indexing', 'vkui/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="['react-todo-js/indexing', 'ring-ui/indexing', 'vkui/indexing']"
        />
      </div>
    </section>

    <Divider title="Vue" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['vue-template/completion/attribute', 'vue-template/completion/component', 'vue3-admin-vite/completion/component', 'vue3-admin-vite/completion/attribute']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['vue-template/localInspection/HelloWorld.vue', 'vue3-admin-vite/localInspection/index.vue']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['vue-template/typing/vue-template']"
        />
      </div>
    </section>

    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="['vue-template/indexing', 'vue3-admin-vite/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="['vue-template/indexing', 'vue3-admin-vite/indexing']"
        />
      </div>
    </section>

    <Divider title="CSS" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['WEB_62578_CSS/completion']"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
</script>
