<template>
  <DashboardPage
    db-name="perfintDev"
    :with-installer="false"
    table="webstorm"
    persistent-id="webstorm_dashboard_tskotlin"
    initial-machine="linux-blade-hetzner"
  >
    <template
      v-for="operation in operations"
      :key="operation.name"
    >
      <Divider
        :title="`TSKotlin ${operation.name}: Port vs Upstream`"
        :description="operation.description"
      />
      <section class="flex gap-x-6 flex-col md:flex-row">
        <div
          v-for="metric in metrics"
          :key="metric.measure"
          class="flex-1 min-w-0"
        >
          <GroupProjectsChart
            :label="metric.label"
            :measure="`tskotlin.${metric.measure}`"
            :projects="[`tskotlin-${operation.name}`, `tskotlin-${operation.name}UPSTREAM`]"
            :aliases="['Port', 'Upstream']"
            value-unit="counter"
          />
        </div>
      </section>
    </template>
    <template
      v-for="mode in programCreationModes"
      :key="mode.name"
    >
      <Divider
        :title="`TSKotlin Kibana program creation: ${mode.label}`"
        :description="mode.description"
      />
      <section class="flex gap-x-6 flex-col md:flex-row">
        <div class="flex-1 min-w-0">
          <GroupProjectsChart
            label="Total program creation time"
            measure="tskotlin.programCreation"
            :projects="[`tskotlin-kibana-program-creation-${mode.name}`]"
          />
        </div>
        <div class="flex-1 min-w-0">
          <GroupProjectsChart
            label="Allocated memory per operation"
            measure="tskotlin.allocated.kib.per.op"
            :projects="[`tskotlin-kibana-program-creation-${mode.name}`]"
            value-unit="counter"
          />
        </div>
      </section>
    </template>
    <Divider
      title="TSKotlin Kibana find usages of 'category'"
      description="Compare runs with different maximum JVM heap sizes (-Xmx): 4 GB, 10 GB and 30 GB."
    />
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Search duration"
          measure="findUsages"
          :projects="findUsagesProjects"
          :aliases="findUsagesAliases"
          value-unit="auto"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

const operations = [
  { name: "parser", description: "Parse checker.ts" },
  { name: "parser-jsdoc", description: "Parse JSDoc in checker.ts" },
  { name: "binder", description: "Bind a fresh AST of checker.ts" },
  { name: "checker", description: "Typecheck the entire Strada compiler, excluding program creation (parsing, binding, module resolution)" },
]

const metrics = [
  { label: "Time per operation", measure: "time.us.per.op" },
  { label: "Allocated memory per operation", measure: "allocated.bytes.per.op" },
  { label: "Retained memory", measure: "retained.bytes" },
]

const programCreationModes = [
  {
    name: "cold",
    label: "Cold",
    description:
      "Create a Kibana program: load, parse and bind source files and resolve modules, without typechecking. Compiler projects, VFS caches and parsed library caches are cleared before each iteration. Config parsing and cache clearing are outside the measurement.",
  },
  {
    name: "hot",
    label: "Hot",
    description:
      "Create a new Kibana program with warm caches, reusing parsed and bound source files and libraries while rebuilding the program and resolving modules, without typechecking. An unmeasured program creation warms the caches. Config parsing and warm-up are outside the measurement.",
  },
]

const findUsagesHeapSizes = [4, 10, 30]
const findUsagesProjects = findUsagesHeapSizes.map((heap) => `tskotlin-kibana-find-usages-${heap}gb/findUsages/chrome-browser-nav-links-category`)
const findUsagesAliases = findUsagesHeapSizes.map((heap) => `${heap} GB`)
</script>
