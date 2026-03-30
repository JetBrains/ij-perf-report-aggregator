<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_product_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="Typing letter by letter (different tabs)" />
    <Divider title="All results" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in allResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="10 Elements added" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in tenResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="First Elements added" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in firstResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="Inserting the whole word (Symbol)" />
    <Divider title="All results" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in wholeWordAllResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="10 Elements added" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in wholeWordTenResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="First Elements added" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in wholeWordFirstResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="Without Warm Up " />
    <Divider title="All results" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in noWarmupAllResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="10 Elements added" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in noWarmupTenResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
    <Divider title="First Elements added" />
    <section class="flex gap-6 flex-col md:flex-row md:flex-wrap">
      <div
        v-for="chart in noWarmupFirstResults"
        :key="chart.definition.label"
        class="min-w-0 md:basis-[calc(33.333%-1rem)] md:flex-none flex-1"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

const allResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["Wait Until Full or Done - All"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter", "radler/luau/new-se-go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter"],
  },
  {
    labels: ["Wait Until Full or Done - Symbol"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter", "radler/luau/new-se-go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter"],
  },
  {
    labels: ["Wait Until Full or Done - Class"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-class-with-warmup/CompileOptions/typingLetterByLetter", "radler/luau/new-se-go-to-class-with-warmup/CompileOptions/typingLetterByLetter"],
  },
  {
    labels: ["Wait Until Full or Done - File"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-file-with-warmup/TableShape.cpp/typingLetterByLetter", "radler/luau/new-se-go-to-file-with-warmup/TableShape.cpp/typingLetterByLetter"],
  },
  {
    labels: ["Wait Until Full or Done - Action"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-action-with-warmup/RCP/typingLetterByLetter", "radler/luau/new-se-go-to-action-with-warmup/RCP/typingLetterByLetter"],
  },
  {
    labels: ["Wait Until Full or Done - Text"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-text-with-warmup/LUAU_BUILD_TESTS/typingLetterByLetter", "radler/luau/new-se-go-to-text-with-warmup/LUAU_BUILD_TESTS/typingLetterByLetter"],
  },
]

const tenResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["10 Elements Added - All"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter", "radler/luau/new-se-go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter"],
  },
  {
    labels: ["10 Elements Added - Symbol"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter", "radler/luau/new-se-go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter"],
  },
  {
    labels: ["10 Elements Added - Class"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-class-with-warmup/CompileOptions/typingLetterByLetter", "radler/luau/new-se-go-to-class-with-warmup/CompileOptions/typingLetterByLetter"],
  },
  {
    labels: ["10 Elements Added - File"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-file-with-warmup/TableShape.cpp/typingLetterByLetter", "radler/luau/new-se-go-to-file-with-warmup/TableShape.cpp/typingLetterByLetter"],
  },
  {
    labels: ["10 Elements Added - Action"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-action-with-warmup/RCP/typingLetterByLetter", "radler/luau/new-se-go-to-action-with-warmup/RCP/typingLetterByLetter"],
  },
  {
    labels: ["10 Elements Added - Text"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-text-with-warmup/LUAU_BUILD_TESTS/typingLetterByLetter", "radler/luau/new-se-go-to-text-with-warmup/LUAU_BUILD_TESTS/typingLetterByLetter"],
  },
]

const firstResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["First Elements Added - All"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter", "radler/luau/new-se-go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter"],
  },
  {
    labels: ["First Elements Added - Symbol"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter", "radler/luau/new-se-go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter"],
  },
  {
    labels: ["First Elements Added - Class"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-class-with-warmup/CompileOptions/typingLetterByLetter", "radler/luau/new-se-go-to-class-with-warmup/CompileOptions/typingLetterByLetter"],
  },
  {
    labels: ["First Elements Added - File"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-file-with-warmup/TableShape.cpp/typingLetterByLetter", "radler/luau/new-se-go-to-file-with-warmup/TableShape.cpp/typingLetterByLetter"],
  },
  {
    labels: ["First Elements Added - Action"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-action-with-warmup/RCP/typingLetterByLetter", "radler/luau/new-se-go-to-action-with-warmup/RCP/typingLetterByLetter"],
  },
  {
    labels: ["First Elements Added - Text"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-text-with-warmup/LUAU_BUILD_TESTS/typingLetterByLetter", "radler/luau/new-se-go-to-text-with-warmup/LUAU_BUILD_TESTS/typingLetterByLetter"],
  },
]

const wholeWordAllResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["Wait Until Full or Done - Type_boolean - 4 results"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type_Boolean/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/Type_Boolean/insertingTheWholeWord"],
  },
  {
    labels: ["Wait Until Full or Done - compiler - 17 results"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol-with-warmup/compiler/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/compiler/insertingTheWholeWord"],
  },
  {
    labels: ["Wait Until Full or Done - constraint - many results"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol-with-warmup/constraint/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/constraint/insertingTheWholeWord"],
  },
  {
    labels: ["Wait Until Full or Done - Type - a lot of results"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/Type/insertingTheWholeWord"],
  },
]

const wholeWordTenResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["10 Elements Added - Type_boolean - 4 results"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type_Boolean/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/Type_Boolean/insertingTheWholeWord"],
  },
  {
    labels: ["10 Elements Added - compiler - 17 results"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol-with-warmup/compiler/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/compiler/insertingTheWholeWord"],
  },
  {
    labels: ["10 Elements Added - constraint - many results"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol-with-warmup/constraint/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/constraint/insertingTheWholeWord"],
  },
  {
    labels: ["10 Elements Added - Type - a lot of results"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/Type/insertingTheWholeWord"],
  },
]

const wholeWordFirstResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["First Elements Added - Type_boolean - 4 results"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type_Boolean/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/Type_Boolean/insertingTheWholeWord"],
  },
  {
    labels: ["First Elements Added - compiler - 17 results"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol-with-warmup/compiler/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/compiler/insertingTheWholeWord"],
  },
  {
    labels: ["First Elements Added - constraint - many results"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol-with-warmup/constraint/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/constraint/insertingTheWholeWord"],
  },
  {
    labels: ["First Elements Added - Type - a lot of results"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol-with-warmup/Type/insertingTheWholeWord"],
  },
]

const noWarmupAllResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["Wait Until Full or Done - compiler - inserting the whole word"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol/compiler/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol/compiler/insertingTheWholeWord"],
  },
  {
    labels: ["Wait Until Full or Done - Type - inserting the whole word"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol/Type/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol/Type/insertingTheWholeWord"],
  },
  {
    labels: ["Wait Until Full or Done - Type_boolean - typing letter by letter"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol/Type_Boolean/typingLetterByLetter", "radler/luau/new-se-go-to-symbol/Type_Boolean/typingLetterByLetter"],
  },
]

const noWarmupTenResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["10 Elements Added - compiler - inserting the whole word"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol/compiler/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol/compiler/insertingTheWholeWord"],
  },
  {
    labels: ["10 Elements Added - Type - inserting the whole word"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol/Type/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol/Type/insertingTheWholeWord"],
  },
  {
    labels: ["10 Elements Added - Type_boolean - typing letter by letter"],
    measures: ["searchEverywhere_elements_added_10"],
    projects: ["radler/luau/go-to-symbol/Type_Boolean/typingLetterByLetter", "radler/luau/new-se-go-to-symbol/Type_Boolean/typingLetterByLetter"],
  },
]

const noWarmupFirstResultsDeclaration: ChartDefinition[] = [
  {
    labels: ["First Elements Added - compiler - inserting the whole word"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol/compiler/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol/compiler/insertingTheWholeWord"],
  },
  {
    labels: ["First Elements Added - Type - inserting the whole word"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol/Type/insertingTheWholeWord", "radler/luau/new-se-go-to-symbol/Type/insertingTheWholeWord"],
  },
  {
    labels: ["First Elements Added - Type_boolean - typing letter by letter"],
    measures: ["searchEverywhere_first_elements_added"],
    projects: ["radler/luau/go-to-symbol/Type_Boolean/typingLetterByLetter", "radler/luau/new-se-go-to-symbol/Type_Boolean/typingLetterByLetter"],
  },
]

const allResults = combineCharts(allResultsDeclaration)
const tenResults = combineCharts(tenResultsDeclaration)
const firstResults = combineCharts(firstResultsDeclaration)

const wholeWordAllResults = combineCharts(wholeWordAllResultsDeclaration)
const wholeWordTenResults = combineCharts(wholeWordTenResultsDeclaration)
const wholeWordFirstResults = combineCharts(wholeWordFirstResultsDeclaration)

const noWarmupAllResults = combineCharts(noWarmupAllResultsDeclaration)
const noWarmupTenResults = combineCharts(noWarmupTenResultsDeclaration)
const noWarmupFirstResults = combineCharts(noWarmupFirstResultsDeclaration)

// const charts = combineCharts([...seNewVsOldChartsDeclaration])
</script>
