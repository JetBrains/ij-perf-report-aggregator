<template>
  <DashboardPage
    db-name="perfintDev"
    :with-installer="false"
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
import { groupBy3 } from "./utils"

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
      "fp-ts/completion/completion",
      "backstage/completion/completion",
      "backstage/completion/component",
      "vue-naive-ui-admin/completion/built-in-completion",
      "vue-naive-ui-admin/completion/completion",
    ],
  },
  {
    label: "FirstCodeAnalysis",
    measure: "firstCodeAnalysis",
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
      "fp-ts/localInspection/eq",
      "backstage/localInspection/UserProfileCard.tsx",
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
      "fp-ts/localInspection/eq",
      "backstage/localInspection/UserProfileCard.tsx",
      "vue-naive-ui-admin/localInspection/inspections",
    ],
  },
  {
    label: "Global inspections",
    measure: "globalInspections",
    projects: ["eslint-plugin-jest/inspection/whole-project", "ring-ui/inspection/src", "fleetbot/inspection/src", "kibana/inspection/src"],
  },
  {
    label: "Typing",
    measure: "typing",
    projects: ["eslint-plugin-jest/typing/eslintPluginJest", "axios/typing/axios", "toh-pt6/typing/toh-pt6", "react-todo-js/typing/react-todo", "vue-template/typing/vue-template"],
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
      "fp-ts/localInspection/eq",
      "backstage/localInspection/UserProfileCard.tsx",
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
      "backstage/completion/completion",
    ],
  },
  {
    label: "Highlighting - remove symbol",
    measure: "typing_EditorBackSpace_duration",
    projects: [
      "material-ui-react-admin/PostEdit.tsxHighlighting",
      "WEB_5976/Card.jsHighlighting",
      "kibana/alerts_grouping.tsxHighlighting",
      "kibana/project_navigation.tsHighlighting",
      "fp-ts/Foldable.tsHighlighting",
      "backstage/UserProfileCard.tsxHighlighting",
    ],
  },
  {
    label: "Highlighting - type symbol",
    measure: "typing_}_duration",
    projects: [
      "material-ui-react-admin/PostEdit.tsxHighlighting",
      "WEB_5976/Card.jsHighlighting",
      "kibana/alerts_grouping.tsxHighlighting",
      "kibana/project_navigation.tsHighlighting",
      "fp-ts/Foldable.tsHighlighting",
      "backstage/UserProfileCard.tsxHighlighting",
    ],
  },
  {
    label: "Find Usages",
    measure: "findUsages",
    projects: ["pancake-frontend/findUsages/Modal.tsx", "backstage/findUsages/types.ts"],
  },
]
</script>
