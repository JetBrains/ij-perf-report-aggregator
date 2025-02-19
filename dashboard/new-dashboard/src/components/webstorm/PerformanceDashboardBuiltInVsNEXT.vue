<template>
  <DashboardPage
    db-name="perfint"
    table="webstorm"
    persistent-id="webstorm_dashboard_builtin_vs_next"
    initial-machine="linux-blade-hetzner"
  >
    <template
      v-for="group in groups"
      :key="group.measure"
    >
      <Divider :title="group.label" />
      <section
        v-for="(groupOf3, groupOf3index) in groupBy3(group.projects)"
        :key="groupOf3index"
        class="flex gap-x-6 flex-col md:flex-row"
      >
        <div
          v-for="project in groupOf3"
          :key="project"
          class="flex-1 min-w-0"
        >
          <GroupProjectsChart
            :label="project"
            :measure="group.measure"
            :projects="[project, project + 'NEXT']"
          />
        </div>
      </section>
    </template>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

const groups = [
  {
    label: "Completion",
    measure: "completion",
    projects: [
      "eslint-plugin-jest/completion/types",
      "axios/completion/functions",
      "toh-pt6/completion/attribute",
      "toh-pt6/completion/component",
      "react-todo-js/completion/attribute",
      "react-todo-js/completion/component",
      "vkui/completion/component",
      "ring-ui/completion/component",
      "material-ui-react-admin/completion/attribute",
      "material-ui-react-admin/completion/component",
      "vue-template/completion/attribute",
      "vue-template/completion/component",
      "vue3-admin-vite/completion/component",
      "vue3-admin-vite/completion/attribute",
      "pancake-frontend/completion/component",
      "pancake-frontend/completion/attribute",
    ],
  },
  {
    label: "Local inspections",
    measure: "localInspections",
    projects: [
      "aws_cdk/localInspection/logging",
      "eslint-plugin-jest/localInspection/misc.ts",
      "allure-js/localInspection/JasmineAllureReporter.ts",
      "axios/localInspection/utils.js",
      "toh-pt6/localInspection/hero.service.ts",
      "toh-pt6/localInspection/heroes.component.html",
      "react-todo-js/localInspection/App.js",
      "WEB_5976/localInspection/react_mui",
      "material-ui-react-admin/localInspection/PostEdit.tsx",
      "vue-template/localInspection/HelloWorld.vue",
      "vue3-admin-vite/localInspection/index.vue",
      "pancake-frontend/localInspection/[tokenId].tsx",
      "kibana/localInspection/alerts_grouping.tsx",
      "kibana/localInspection/project_navigation.ts",
    ],
  },
  {
    label: "Global inspections",
    measure: "globalInspections",
    projects: ["eslint-plugin-jest/globalInspection/whole-project", "ring-ui/globalInspection/src", "fleetbot/globalInspection/src", "kibana/globalInspection/src"],
  },
  {
    label: "Typing",
    measure: "typing",
    projects: [
      "eslint-plugin-jest/typing/eslintPluginJest",
      "axios/typing/axios",
      "toh-pt6/typing/toh-pt6",
      "react-todo-js/typing/react-todo-js",
      "vue-template/typing/vue-template",
    ],
  },
  {
    label: "Code Vision",
    measure: "JSReferencesCodeVisionProvider",
    projects: [
      "aws_cdk/localInspection/logging",
      "WEB_5976/localInspection/react_mui",
      "toh-pt6/localInspection/hero.service.ts",
      "vue3-admin-vite/localInspection/index.vue",
      "eslint-plugin-jest/localInspection/misc.ts",
      "allure-js/localInspection/JasmineAllureReporter.ts",
      "material-ui-react-admin/localInspection/PostEdit.tsx",
    ],
  },
  {
    label: "Completion First Element",
    measure: "completion#firstElementShown#mean_value",
    projects: [
      "material-ui-react-admin/completion/attribute",
      "pancake-frontend/completion/component",
      "pancake-frontend/completion/attribute",
      "eslint-plugin-jest/completion/types",
      "axios/completion/functions",
    ],
  },
]

function groupBy3<T>(array: T[]): T[][] {
  const result = []
  for (let i = 0; i < array.length; i += 3) {
    const component = [array[i]]
    if (i + 1 < array.length) component.push(array[i + 1])
    if (i + 2 < array.length) component.push(array[i + 2])
    result.push(component)
  }
  return result
}
</script>
