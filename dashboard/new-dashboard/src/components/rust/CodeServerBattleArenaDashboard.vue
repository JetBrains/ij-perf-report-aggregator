<template>
  <DashboardPage
    db-name="perfintDev"
    table="rust"
    persistent-id="rust_code_server_battle_arena_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="Startup" />
    <section class="flex gap-x-6">
      <div
        v-for="chart in startupCharts"
        :key="chart.label"
        class="flex-1 min-w-0"
      >
        <GroupProjectsChart
          :label="chart.label"
          :measure="chart.measure"
          :projects="chart.projects"
          :aliases="chart.aliases"
          :value-unit="chart.valueUnit"
          :legend-formatter="chart.legendFormatter"
        />
      </div>
    </section>
    <template
      v-for="section in sections"
      :key="section.title"
    >
      <Divider :title="section.title" />
      <section
        v-for="row in section.timingRows"
        :key="row.key"
        class="flex gap-x-6"
      >
        <div
          v-for="chart in row.charts"
          :key="chart.label"
          class="flex-1 min-w-0"
        >
          <GroupProjectsChart
            :label="chart.label"
            :measure="chart.measure"
            :projects="chart.projects"
            :aliases="chart.aliases"
            :value-unit="chart.valueUnit"
            :legend-formatter="chart.legendFormatter"
          />
        </div>
      </section>
      <ChartAccordion :lazy="true">
        <AccordionPanel value="0">
          <AccordionHeader>{{ section.title }}: memory</AccordionHeader>
          <AccordionContent>
            <section
              v-for="row in section.memoryRows"
              :key="row.key"
              class="flex gap-x-6"
            >
              <div
                v-for="chart in row.charts"
                :key="chart.label"
                class="flex-1 min-w-0"
              >
                <GroupProjectsChart
                  :label="chart.label"
                  :measure="chart.measure"
                  :projects="chart.projects"
                  :aliases="chart.aliases"
                  :value-unit="chart.valueUnit"
                  :legend-formatter="chart.legendFormatter"
                />
              </div>
            </section>
          </AccordionContent>
        </AccordionPanel>
      </ChartAccordion>
    </template>
  </DashboardPage>
</template>

<script setup lang="ts">
import AccordionContent from "openvue/accordioncontent"
import AccordionHeader from "openvue/accordionheader"
import AccordionPanel from "openvue/accordionpanel"
import ChartAccordion from "../charts/ChartAccordion.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ValueUnit } from "../common/chart"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

// The metrics of RustCodeServerBattleArenaTest in ultimate (tests/rustrover). Each test publishes `<project>-<server>/<scenario>`.
// A scenario runs its cold request after a cold start of the server, and its warm request after a restart on the same project.
// Each typing phase of a request runs in a restart of its own, so a chart can show the warm phase and the typing phases together.
interface Scenario {
  id: string
  label: string
  cold: string
  warm: string
  typings: string[]
  // RustRover-only metrics `<part>.<phase>` that are plotted on the same charts as the request metrics
  mergedParts: string[]
  // RustRover-only metrics `<part>.<phase>` that have charts of their own
  rustRoverParts: string[]
}

const PROJECTS = [
  { id: "bevy", label: "Bevy" },
  { id: "ripgrep", label: "ripgrep" },
]

const SCENARIOS: Scenario[] = [
  {
    id: "completion",
    label: "Completion",
    cold: "statementCompletion",
    warm: "memberCompletion",
    typings: ["afterLocalTyping", "afterTopLevelTyping"],
    mergedParts: [],
    rustRoverParts: ["timeToShow"],
  },
  {
    id: "definition",
    label: "Go to Definition",
    cold: "definition",
    warm: "structDefinition",
    typings: ["afterTopLevelTyping"],
    mergedParts: ["resolve"],
    rustRoverParts: [],
  },
  {
    id: "references",
    label: "Find References",
    cold: "references",
    warm: "methodReferences",
    typings: ["afterTopLevelTyping"],
    mergedParts: [],
    rustRoverParts: [],
  },
]

const SERVERS = [
  { id: "rustrover", label: "RustRover" },
  { id: "rust-analyzer", label: "rust-analyzer" },
]

interface Chart {
  label: string
  measure: string[]
  projects: string[]
  aliases: string[]
  valueUnit: ValueUnit
  legendFormatter: (name: string) => string
}

interface Row {
  key: string
  charts: Chart[]
}

interface Section {
  title: string
  timingRows: Row[]
  memoryRows: Row[]
}

/** A chart of the [measures] of both servers. The legend shows the measure without the [request] prefix, and the server. */
function chart(project: string, scenario: Scenario, label: string, measures: string[], valueUnit: ValueUnit, request = "", rustRoverOnly = false): Chart {
  const servers = rustRoverOnly ? SERVERS.slice(0, 1) : SERVERS
  return {
    label,
    measure: measures,
    projects: servers.map((server) => `${project}-${server.id}/${scenario.id}`),
    aliases: servers.map((server) => server.label),
    valueUnit,
    legendFormatter: (name) => (request === "" ? name : name.replaceAll(`${request}.`, "")),
  }
}

/**
 * A row with the cold request on the left and the warm request on the right.
 * Each chart shows the metric `<request>.<prefix><phase>` for each of the [prefixes] and [phases]. The [phases] are the metric name suffixes, such as `warm`.
 */
function requestRow(
  project: string,
  scenario: Scenario,
  prefixes: string[],
  phases: string[],
  valueUnit: ValueUnit,
  rustRoverOnly = false,
  label = `${prefixes[0]}${phases.join(", ")}`
): Row {
  return {
    key: label,
    charts: [scenario.cold, scenario.warm].map((request) =>
      chart(
        project,
        scenario,
        `${request}: ${label}`,
        prefixes.flatMap((prefix) => phases.map((phase) => `${request}.${prefix}${phase}`)),
        valueUnit,
        request,
        rustRoverOnly
      )
    ),
  }
}

/**
 * A row with the cold request on the left and the warm request on the right.
 * Each chart shows the exact [metrics], which are the metric names without the `<request>.` prefix.
 */
function metricsRow(project: string, scenario: Scenario, metrics: string[], valueUnit: ValueUnit, label: string, rustRoverOnly = false): Row {
  return {
    key: label,
    charts: [scenario.cold, scenario.warm].map((request) =>
      chart(
        project,
        scenario,
        `${request}: ${label}`,
        metrics.map((metric) => `${request}.${metric}`),
        valueUnit,
        request,
        rustRoverOnly
      )
    ),
  }
}

/** The peak RSS of both servers, and the peak JVM heap of RustRover, in each of the [phases]. */
function peaks(phases: string[]): string[] {
  return phases.flatMap((phase) => [`rss.peak.${phase}`, `heap.peak.${phase}`])
}

function section(project: { id: string; label: string }, scenario: Scenario): Section {
  // `first` is often much longer than the other phases, so it has a chart of its own
  const phases = ["warm", ...scenario.typings]
  // The request metrics go first, and the merged RustRover metrics follow them
  const requestPrefixes = ["", ...scenario.mergedParts.map((part) => `${part}.`)]
  return {
    title: `${project.label}: ${scenario.label}`,
    timingRows: [
      requestRow(project.id, scenario, requestPrefixes, ["first"], "ms"),
      requestRow(project.id, scenario, requestPrefixes, phases, "ms"),
      ...scenario.rustRoverParts.flatMap((part) => [
        requestRow(project.id, scenario, [`${part}.`], ["first"], "ms", true),
        requestRow(project.id, scenario, [`${part}.`], phases, "ms", true),
      ]),
    ],
    // A chart shows the RSS of both servers, and the JVM heap of RustRover only. `rss.before_exit` is the RSS after the last request of the session.
    memoryRows: [
      metricsRow(project.id, scenario, [...peaks(["first", "warm"]), "rss.before_exit"], "counter", "Peak RSS and JVM heap: first and warm"),
      ...scenario.typings.map((typing) => metricsRow(project.id, scenario, [...peaks([typing]), `rss.before_exit.${typing}`], "counter", `Peak RSS and JVM heap: ${typing}`)),
    ],
  }
}

// Each scenario starts the servers in the same way, so the startup charts take the times of one scenario only
const STARTUP_SCENARIO = SCENARIOS[0]

/** A chart of the [measures] of the startup scenario, with a line for each project, server and measure. */
function startupChart(label: string, measures: string[], valueUnit: ValueUnit, servers = SERVERS): Chart {
  const lines = PROJECTS.flatMap((project) => servers.map((server) => ({ project, server })))
  return {
    label,
    measure: measures,
    projects: lines.map(({ project, server }) => `${project.id}-${server.id}/${STARTUP_SCENARIO.id}`),
    aliases: lines.map(({ project, server }) => `${project.label}: ${server.label}`),
    valueUnit,
    legendFormatter: (name) => name,
  }
}

const startupCharts = [
  startupChart(`ready (cold start, ${STARTUP_SCENARIO.cold})`, ["ready"], "ms"),
  startupChart(`warmReady (restart, ${STARTUP_SCENARIO.warm})`, ["warmReady"], "ms"),
]

const sections = PROJECTS.flatMap((project) => SCENARIOS.map((scenario) => section(project, scenario)))
</script>
