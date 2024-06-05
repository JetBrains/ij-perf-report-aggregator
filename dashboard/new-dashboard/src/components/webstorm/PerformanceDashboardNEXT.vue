<template>
  <DashboardPage
    db-name="perfint"
    table="webstorm"
    persistent-id="webstorm_dashboard_next"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Code Vision (JSReferencesCodeVisionProvider)"
        measure="JSReferencesCodeVisionProvider"
        :projects="[
          'aws_cdk/localInspection/loggingNEXT',
          'WEB_5976/localInspection/react_muiNEXT',
          'toh-pt6/localInspection/hero.service.tsNEXT',
          'vue3-admin-vite/localInspection/index.vueNEXT',
          'eslint-plugin-jest/localInspection/misc.tsNEXT',
          'allure-js/localInspection/JasmineAllureReporter.tsNEXT',
          'ts-codec/localInspection/codec.test.tsNEXT',
          'material-ui-react-admin/localInspection/PostEdit.tsxNEXT',
        ]"
      />
    </section>

    <Divider title="TypeScript" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['eslint-plugin-jest/completion/typesNEXT']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="[
            'aws_cdk/localInspection/loggingNEXT',
            'eslint-plugin-jest/localInspection/misc.tsNEXT',
            'allure-js/localInspection/JasmineAllureReporter.tsNEXT',
            'ts-codec/localInspection/codec.test.tsNEXT',
          ]"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['eslint-plugin-jest/typingNEXT']"
        />
      </div>
    </section>

    <Divider title="JavaScript" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['axios/completion/functionsNEXT']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['axios/localInspection/utils.jsNEXT']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['axios/typingNEXT']"
        />
      </div>
    </section>

    <Divider title="Angular" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="['toh-pt6/completion/attributeNEXT', 'toh-pt6/completion/componentNEXT']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['toh-pt6/localInspection/hero.service.tsNEXT', 'toh-pt6/localInspection/heroes.component.htmlNEXT']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['toh-pt6/typing/toh-pt6NEXT']"
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
            'react-todo-js/completion/attributeNEXT',
            'react-todo-js/completion/componentNEXT',
            'vkui/completion/componentNEXT',
            'ring-ui/completion/componentNEXT',
            'material-ui-react-admin/completion/attributeNEXT',
          ]"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['react-todo-js/localInspection/App.jsNEXT', 'WEB_5976/localInspection/react_muiNEXT', 'material-ui-react-admin/localInspection/PostEdit.tsxNEXT']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['react-todo-js/typingNEXT']"
        />
      </div>
    </section>

    <Divider title="Vue" />

    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Completion"
          measure="completion"
          :projects="[
            'vue-template/completion/attributeNEXT',
            'vue-template/completion/componentNEXT',
            'vue3-admin-vite/completion/componentNEXT',
            'vue3-admin-vite/completion/attributeNEXT',
          ]"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local inspections"
          measure="localInspections"
          :projects="['vue-template/localInspection/HelloWorld.vueNEXT', 'vue3-admin-vite/localInspection/index.vueNEXT']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing"
          measure="typing"
          :projects="['vue-template/typingNEXT']"
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
