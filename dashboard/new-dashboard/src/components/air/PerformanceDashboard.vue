<template>
  <DashboardPage
    db-name="perfintDev"
    table="air"
    persistent-id="air_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <template
      v-for="group in longSessionGroups"
      :key="group.title"
    >
      <Divider :title="group.title" />
      <section>
        <GroupProjectsChart
          v-for="chart in group.charts"
          :key="chart.label"
          :label="chart.label"
          :description="chart.description"
          :measure="chart.measure"
          :projects="longSessionProjects"
          :aliases="longSessionAliases"
          :legend-formatter="sizeLegend"
        />
      </section>
    </template>
    <template
      v-for="group in otherGroups"
      :key="group.title"
    >
      <Divider :title="group.title" />
      <section>
        <GroupProjectsChart
          v-for="chart in group.charts"
          :key="chart.label"
          :label="chart.label"
          :description="chart.description"
          :measure="chart.measure"
          :projects="[group.project]"
          :legend-formatter="measureLegend"
        />
      </section>
    </template>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { SERIES_NAME_SEPARATOR } from "../common/DataQueryExecutor"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

interface AirChart {
  label: string
  description: string
  measure: string[]
}

interface AirGroup {
  title: string
  charts: AirChart[]
}

/**
 * The sizes of the long session test. The metric names and the projects keep the size names of the test, M, XL and XXL,
 * so their history and their degradation alerts stay. The legend names each size by what it is.
 */
const sizes = [
  { id: "M", label: "Small session: 30 turns, 9 MB journal" },
  { id: "XL", label: "Medium session: 200 turns, 63 MB journal" },
  { id: "XXL", label: "Large session: 500 turns, 160 MB journal" },
]
const longSessionProjects = sizes.map((size) => `air/long-session/${size.id}`)
const longSessionAliases = sizes.map((size) => size.label)
const sizeLabels = new Set(longSessionAliases)

/** One chart of a long session step with one line for each size. `%s` in [template] is the size of the metric name. */
function bySize(label: string, description: string, template: string): AirChart {
  return { label, description, measure: sizes.map((size) => template.replace("%s", size.id)) }
}

/** A line of a long session chart is its size. The measure part of the series name is the same step for each size. */
function sizeLegend(name: string): string {
  return name.split(SERIES_NAME_SEPARATOR).find((part) => sizeLabels.has(part)) ?? name
}

const jumpPoints = [
  { point: 0, label: "the start of the history" },
  { point: 20, label: "20 % of the history" },
  { point: 40, label: "40 % of the history" },
  { point: 60, label: "60 % of the history" },
  { point: 80, label: "80 % of the history" },
]

const longSessionGroups: AirGroup[] = [
  {
    title: "Long session: open after an IDE restart",
    charts: [
      bySize(
        "Open a session: from the click in Agent Sessions to the open tab",
        "The user clicks the session row in the Agent Sessions tree after an IDE restart. The tab can still be empty.",
        "open.%s.gestureToTabMs"
      ),
      bySize(
        "Open a session: from the click to the first visible text",
        "Until the loading cover of the transcript goes away and the user sees the text of the session.",
        "open.%s.coverRemovedMs"
      ),
      bySize(
        "Open a session: from the click to a settled transcript",
        "Until the transcript is shown, laid out, does not move, and the EDT answers in less than 50 ms.",
        "open.%s.gestureToSettledMs"
      ),
    ],
  },
  {
    title: "Long session: open, the steps that the IDE logs",
    charts: [
      bySize(
        "Open, backend: replay the session journal into the projection",
        "From the request of the projection until it is ready, with the lock waits. The [AIR_SESSION_RESTORE] rich projection ready line of the IDE log.",
        "checkpoint.%s.rich_projection_ready.totalMs"
      ),
      bySize(
        "Open, frontend: apply the session snapshot and render it on the EDT",
        "From the arrival of the snapshot in the frontend until the screen rendered it. The [AIR_SESSION_RESTORE] frontend reload applied line.",
        "checkpoint.%s.frontend_reload_applied.totalMs"
      ),
      bySize(
        "Open, frontend: build the rows of the transcript",
        "Build the rows of the messages, the tool calls and the code blocks. The [AIR_SESSION_RESTORE] frontend transcript rows applied line.",
        "checkpoint.%s.frontend_transcript_rows_applied.elapsedMs"
      ),
      bySize(
        "Open, frontend: from the session panel to the first painted text",
        "From the creation of the session panel until the first paint of the transcript. The [AIR_SESSION_RESTORE] frontend first content paint finished line.",
        "checkpoint.%s.frontend_first_content_paint_finished.totalSincePanelCreationMs"
      ),
    ],
  },
  {
    title: "Long session: continue the old session with a prompt",
    charts: [
      bySize(
        "Continue a session: from Send to the session/load request at the agent",
        "Air starts the agent process and sends session/load. The time includes the start of the agent.",
        "continue.%s.sendToLoadMs"
      ),
      bySize(
        "Continue a session: from Send to the end of the history replay",
        "The agent replays one short turn for each turn of the session before it takes the prompt.",
        "continue.%s.sendToReplayedMs"
      ),
      bySize("Continue a session: from Send to the prompt at the agent", "Until the agent got the prompt that the user sent.", "continue.%s.sendToPromptMs"),
    ],
  },
  {
    title: "Long session: a long streamed answer",
    charts: [
      bySize(
        "Long answer: how long the agent streams it",
        "The test agent sends the chunks of a Markdown answer with code, images and tool calls. The test sets this time.",
        "answer.%s.streamMs"
      ),
      bySize(
        "Long answer: from the stream start to the completed turn in the UI",
        "Until the UI shows the turn as completed and the prompt is ready again. The difference to the stream time is how far the UI lags behind the agent.",
        "answer.%s.streamStartToCompletedMs"
      ),
      bySize("Long answer: from the stream start to a settled transcript", "Until the transcript after the answer does not move any more.", "answer.%s.streamStartToSettledMs"),
    ],
  },
  {
    title: "Long session: typing and scrolling latency",
    charts: [
      bySize(
        "Mouse wheel up from the bottom of the transcript: p95 latency of one wheel event",
        "10 bursts of 400 wheel events, one each 16 ms. The p95 of the worst burst. An event over 2 s fails the test.",
        "scroll.%s.wheelP95Ms"
      ),
      bySize("Mouse wheel up from the bottom of the transcript: the slowest single wheel event", "The longest wait of one wheel event in the 10 bursts.", "scroll.%s.wheelMaxMs"),
      bySize("Mouse wheel right after the long answer: p95 latency of one wheel event", "200 wheel events over the transcript next to the long answer.", "answer.%s.wheelP95Ms"),
      bySize("Typing prose into the prompt: p95 latency of one character", "1000 characters of prose at the pace of a fast typist.", "prompt.%s.proseP95Ms"),
      bySize(
        "Typing code into the prompt: p95 latency of one character",
        "50 lines of code at the pace of a fast typist. A character over 500 ms fails the test.",
        "prompt.%s.codeP95Ms"
      ),
    ],
  },
  {
    title: "Long session: jumps far into the history",
    charts: jumpPoints.map(({ point, label }) =>
      bySize(
        `Jump to ${label}: time until the transcript settles`,
        "The scrollbar jumps to this point of the history, as a user drags it there. The jump builds and lays out rows far from the bottom. Over 60 s fails the test.",
        `scroll.%s.jump${point}.settledMs`
      )
    ),
  },
  {
    title: "Long session: memory",
    charts: [
      bySize("Live heap before the open of the session", "The used heap after a full garbage collection, before the click.", "heap.%s.before_the_open.MiB"),
      bySize("Live heap right after the open of the session", "The difference to the heap before the open is what the open session keeps.", "heap.%s.after_the_open.MiB"),
      bySize(
        "Live heap at the end of the scenario, after a full garbage collection",
        "What the IDE keeps after the open, the scroll, the answer and the typing. Over 3/4 of the max heap fails the test.",
        "heap.%s.heapLive.MiB"
      ),
      bySize("Peak heap during the scenario", "The largest used heap that the watchdog saw during the scenario.", "heap.%s.heapPeak.MiB"),
    ],
  },
]

/** The readable name of each measure of the charts that have one line for each measure. */
const measureLabels: Record<string, string> = {
  "manyAgents.edtP95Ms": "EDT latency, p95",
  "manyAgents.edtMaxMs": "EDT latency, max",
  "manyAgents.roundMaxMs": "Slowest round of one update to each session",
  "manyAgents.focusedTabLagMs": "Focused tab",
  "manyAgents.openTabLagMaxMs": "Slowest of the hidden open tabs",
  "manyAgents.heapLive.MiB": "Live heap after a full garbage collection",
  "manyAgents.heapPeak.MiB": "Peak heap",
  "manyAgents.gcMs": "Garbage collection time since the IDE start",
  "sessionList.firstShowMs": "First show after the restart",
  "sessionList.refreshMedianMs": "Refresh, median of 9",
  "sessionList.wheelP95Ms": "Wheel event latency, p95",
  "sessionList.wheelMaxMs": "Wheel event latency, max",
  "sessionList.heapLive.MiB": "Live heap after a full garbage collection",
  "sessionList.heapPeak.MiB": "Peak heap",
}

function measureLegend(name: string): string {
  const measure = name.split(SERIES_NAME_SEPARATOR).find((part) => part in measureLabels)
  return measure == undefined ? name : measureLabels[measure]
}

const otherGroups: (AirGroup & { project: string })[] = [
  {
    title: "Many agents: 10 sessions stream at once for 5 minutes",
    project: "air/many-agents",
    charts: [
      {
        label: "Ten agents stream at once: EDT latency and the slowest round over all sessions",
        description: "10 sessions of a test agent stream with 5 subagents each. One is in the focused tab, 3 in hidden tabs, 6 without a tab.",
        measure: ["manyAgents.edtP95Ms", "manyAgents.edtMaxMs", "manyAgents.roundMaxMs"],
      },
      {
        label: "Ten agents stream at once: from the end of the stream to the completed turn in a tab",
        description: "How far an open tab lags behind its agent.",
        measure: ["manyAgents.focusedTabLagMs", "manyAgents.openTabLagMaxMs"],
      },
      {
        label: "Ten agents stream at once: memory",
        description: "Over 3/4 of the max heap fails the test.",
        measure: ["manyAgents.heapLive.MiB", "manyAgents.heapPeak.MiB"],
      },
      { label: "Ten agents stream at once: garbage collection time", description: "The time of all collections since the IDE start.", measure: ["manyAgents.gcMs"] },
    ],
  },
  {
    title: "Session list: the Agent Sessions window with 2500 sessions",
    project: "air/session-list",
    charts: [
      {
        label: "Session list of 2500 sessions: first show after a restart, and a refresh",
        description: "10 sessions of 2000 turns, 100 of 150 turns, the rest of one turn, each with a valid summary.json. A journal read fails the test.",
        measure: ["sessionList.firstShowMs", "sessionList.refreshMedianMs"],
      },
      {
        label: "Session list of 2500 sessions: mouse wheel latency over the tree",
        description: "200 wheel events over the Agent Sessions tree.",
        measure: ["sessionList.wheelP95Ms", "sessionList.wheelMaxMs"],
      },
      {
        label: "Session list of 2500 sessions: memory",
        description: "Over 3/4 of the max heap fails the test.",
        measure: ["sessionList.heapLive.MiB", "sessionList.heapPeak.MiB"],
      },
    ],
  },
]
</script>
