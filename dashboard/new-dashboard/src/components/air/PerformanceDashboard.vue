<template>
  <DashboardPage
    db-name="perfintDev"
    table="air"
    persistent-id="air_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
    :history-start="historyStart"
  >
    <template #configurator>
      <SelectButton
        v-model="selectedSizes"
        :options="sizes"
        option-label="short"
        option-value="id"
        multiple
        :allow-empty="false"
        aria-label="The long session sizes that the charts show"
      />
    </template>
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
      v-for="group in acpGroups"
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
          :projects="chart.projects"
          :legend-formatter="acpLegend"
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
import { useStorage } from "@vueuse/core"
import { computed } from "vue"
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
/**
 * The first day of the points that the charts show. The tests changed how they measure on that day, and the earlier
 * points are not comparable.
 */
const historyStart = "2026-11-12"

const sizes = [
  {
    id: "M",
    short: "Small: 30 turns",
    label: "Small session: 30 turns, 9 MB journal",
  },
  {
    id: "XL",
    short: "Medium: 200 turns",
    label: "Medium session: 200 turns, 63 MB journal",
  },
  {
    id: "XXL",
    short: "Large: 500 turns",
    label: "Large session: 500 turns, 160 MB journal",
  },
]
const sizeLabels = new Set(sizes.map((size) => size.label))

/**
 * The sizes that every long session chart shows, chosen once for the whole page. A legend entry still hides one line of
 * one chart. The choice stays for the next visit of this browser.
 */
const selectedSizes = useStorage<string[]>(
  "air_dashboard_sizes",
  sizes.map((size) => size.id)
)
const shownSizes = computed(() => sizes.filter((size) => selectedSizes.value.includes(size.id)))
const longSessionProjects = computed(() => shownSizes.value.map((size) => `air/long-session/${size.id}`))
const longSessionAliases = computed(() => shownSizes.value.map((size) => size.label))

/** One chart of a long session step with one line for each size. `%s` in [template] is the size of the metric name. */
function bySize(label: string, description: string, template: string): AirChart {
  return {
    label,
    description,
    measure: sizes.map((size) => template.replace("%s", size.id)),
  }
}

/** A line of a long session chart is its size. The measure part of the series name is the same step for each size. */
function sizeLegend(name: string): string {
  return name.split(SERIES_NAME_SEPARATOR).find((part) => sizeLabels.has(part)) ?? name
}

const jumpPoints = [
  { point: 0, label: "the start of the history" },
  { point: 40, label: "40 % of the history" },
  { point: 80, label: "80 % of the history" },
]

const longSessionGroups: AirGroup[] = [
  {
    title: "Long session: open after an IDE restart",
    charts: [
      bySize(
        "Open a session: from the click in Agent Sessions to the open tab",
        "The user double-clicks the session row after an IDE restart. The time runs from the mouse event to the moment the session screen shows. The tab can still be empty.",
        "open.%s.gestureToTabMs"
      ),
      bySize(
        "Open a session: from the click to the first visible text",
        "Until the loading cover goes away. The page removes it after the text near the viewport is built and laid out, so the next frame shows the transcript. A removal by the timeout of the page fails the test.",
        "open.%s.coverRemovedMs"
      ),
    ],
  },
  {
    title: "Long session: open, the steps that the IDE logs",
    charts: [
      bySize(
        "Open, backend: replay the session journal into the projection",
        "From the request of the projection until it is ready, with the lock waits. It runs off the EDT. The [AIR_SESSION_RESTORE] rich projection ready line of the IDE log.",
        "checkpoint.%s.rich_projection_ready.totalMs"
      ),
      bySize(
        "Open, frontend: apply the session snapshot to the model",
        "The work of the frontend on the snapshot of the session before the render. Its render on the EDT takes about 1 ms, and its wait for the EDT is the other work of the open, which the rows chart shows. The [AIR_SESSION_RESTORE] frontend reload applied line.",
        "checkpoint.%s.frontend_reload_applied.modelApplyMs"
      ),
      bySize(
        "Open, frontend: build the rows of the transcript",
        "One EDT event that builds the rows of the messages, the tool calls and the code blocks. The largest freeze of the open. The [AIR_SESSION_RESTORE] frontend transcript rows applied line.",
        "checkpoint.%s.frontend_transcript_rows_applied.elapsedMs"
      ),
      bySize(
        "Open, frontend: from the session panel to the first painted text",
        "From the creation of the session panel until the paint that removes the cover. The [AIR_SESSION_RESTORE] frontend first content paint finished line.",
        "checkpoint.%s.frontend_first_content_paint_finished.totalSincePanelCreationMs"
      ),
    ],
  },
  {
    title: "Long session: continue the old session with a prompt",
    charts: [
      bySize(
        "Continue a session: from the answer of session/load to the prompt at the agent",
        "After the agent answered session/load, Air drops the replayed updates, reads the tail of the journal and recovers the owner of the session, then sends the prompt. The time ends when the prompt arrives at the agent.",
        "continue.%s.loadAnsweredToPromptMs"
      ),
    ],
  },
  {
    title: "Long session: a long streamed answer",
    charts: [
      bySize(
        "Long answer: from a written chunk to the frame that shows it, p95",
        "The test agent streams a long Markdown answer. The first 16 sections stream at 6000 characters a second, the pace of a fast agent. For each of their section marks, the time runs from the write of its chunk to the end of the first frame that shows it. It shows how far the screen lags behind the agent.",
        "answer.%s.streamLagP95Ms"
      ),
      bySize(
        "Long answer: from the last chunk to a settled transcript",
        "From the write of the last chunk until the end of the first frame with the built text after the last change of the transcript. The EDT backlog at the end of the stream counts.",
        "answer.%s.streamEndToSettledMs"
      ),
    ],
  },
  {
    title: "Long session: scrolling and typing",
    charts: [
      bySize(
        "Mouse wheel up from the bottom of the transcript: p95 latency of one wheel event",
        "One burst of up to 375 wheel notches, one each 16 ms, for 6 s, into history that is not built yet. Animated scrolling is off. Each event runs from its post to the end of the first frame after its dispatch, so the wait in the EDT queue counts. The burst stops at the top. An event over 2 s fails the test.",
        "scroll.%s.wheelP95Ms"
      ),
      bySize(
        "Typing code into the prompt: p95 latency of one character",
        "1740 characters of code. The next key goes after the paint of the last one, the end of its completion and an idle EDT, so each value is the cost of one key. The time runs from the key event to the start of the editor paint. A character over 500 ms fails the test.",
        "prompt.%s.codeP95Ms"
      ),
      bySize(
        "Typing into the prompt: p95 latency of a line break",
        "Shift+Enter in the prose and the code of the prompt, from the key event to the paint of the editor.",
        "prompt.%s.lineBreakP95Ms"
      ),
    ],
  },
  {
    title: "Long session: jumps far into the history",
    charts: jumpPoints.map(({ point, label }) =>
      bySize(
        `Jump to ${label}: time until the transcript settles`,
        "Right after the open, the scrollbar jumps to this point of history that is not built yet, as a user drags it there. The time runs to the end of the first frame with no empty part in the viewport, after the last change. A jump to a built place or to the shown position gives no value. Over 60 s fails the test.",
        `scroll.%s.jump${point}.settledMs`
      )
    ),
  },
  {
    title: "Long session: memory",
    charts: [
      bySize(
        "Live heap right after the open of the session",
        "The used heap after a full garbage collection, after the open. It holds what the open session keeps.",
        "heap.%s.after_the_open.MiB"
      ),
      bySize(
        "Live heap at the end of the scenario, after a full garbage collection",
        "What the IDE keeps after the open, the scroll, the answer and the typing. The lane IDE has the heap of a user IDE, 2048 MiB. Over 3/4 of it fails the test.",
        "heap.%s.heapLive.MiB"
      ),
    ],
  },
]

/** The readable name of each measure of the charts that have one line for each measure. */
const measureLabels: Record<string, string> = {
  "manyAgents.focusedContentLagP95Ms": "p95",
  "manyAgents.focusedContentLagMaxMs": "max",
  "manyAgents.edtMaxMs": "The longest EDT wait",
  "manyAgents.heapLive.MiB": "Live heap after a full garbage collection",
  "sessionList.firstShowMs": "First show: from the activation of the window to the frame with the newest row",
  "sessionList.refreshMedianMs": "Refresh, median of 9",
  "sessionList.heapLive.MiB": "Live heap after a full garbage collection",
}

function measureLegend(name: string): string {
  const measure = name.split(SERIES_NAME_SEPARATOR).find((part) => part in measureLabels)
  return measure == undefined ? name : measureLabels[measure]
}

const otherGroups: (AirGroup & { project: string })[] = [
  {
    title: "Many agents: 10 sessions stream at once for 20 seconds",
    project: "air/many-agents",
    charts: [
      {
        label: "Ten agents stream at once: from a written chunk to the frame that shows it in the focused tab",
        description:
          "10 sessions of a test agent stream, each into one growing answer, with 5 subagents each. One is in the focused tab, 3 in hidden tabs, 6 without a tab. The time runs from the write of a marked chunk to the end of the first frame that shows it.",
        measure: ["manyAgents.focusedContentLagP95Ms", "manyAgents.focusedContentLagMaxMs"],
      },
      {
        label: "Ten agents stream at once: the longest EDT wait",
        description: "A probe posts to the EDT each 10 ms during the stream. The longest time that one probe waited.",
        measure: ["manyAgents.edtMaxMs"],
      },
      {
        label: "Ten agents stream at once: memory",
        description: "Over 3/4 of the max heap fails the test.",
        measure: ["manyAgents.heapLive.MiB"],
      },
    ],
  },
  {
    title: "Session list: the Agent Sessions window with 2500 sessions",
    project: "air/session-list",
    charts: [
      {
        label: "Session list of 2500 sessions: first show after a restart, and a refresh",
        description:
          "The part waits until the list holds all 2500 sessions. The first show runs from the activation of the window to the end of the first frame that shows the newest row. A refresh bypasses the caches. A journal read fails the test.",
        measure: ["sessionList.firstShowMs", "sessionList.refreshMedianMs"],
      },
      {
        label: "Session list of 2500 sessions: memory",
        description: "Over 3/4 of the max heap fails the test.",
        measure: ["sessionList.heapLive.MiB"],
      },
    ],
  },
]

/**
 * The ACP adapter performance tests (lane acp-perf, IJAI-2949). Each night they install the latest claude-acp, codex-acp
 * and junie from the ACP registry and talk to them with a bare JSON-RPC client and a mock model, without Air. A project is
 * one adapter, `air/acp/<agent>`. The big sessions have the sizes of the long session test, and a measure names its size.
 * The test writes the sessions of junie through prompts, so junie has no sizes and no list scenarios.
 */
const acpAgents = [
  { project: "air/acp/claude-acp", label: "claude-acp" },
  { project: "air/acp/codex-acp", label: "codex-acp" },
  { project: "air/acp/junie", label: "junie" },
]
const acpAgentLabels = new Map(acpAgents.map((agent) => [agent.project, agent.label]))
const acpSizedProjects = acpAgents.filter((agent) => agent.label != "junie").map((agent) => agent.project)
const acpAllProjects = acpAgents.map((agent) => agent.project)

const acpMeasureLabels: Record<string, string> = {
  "live.agentMs": "first turn",
  "live.second.agentMs": "second turn",
  "live.chunkLagP50Ms": "p50",
  "live.chunkLagP95Ms": "p95",
  "live.agentRssPeakMiB": "first turn",
  "live.second.agentRssPeakMiB": "second turn",
  "list.1k.activeFirst.warmMs": "1 000 sessions",
  "list.10k.activeFirst.warmMs": "10 000 sessions",
  "list.50k.activeFirst.warmMs": "50 000 sessions",
  "list.initialize.coldMs": "61 000 sessions in the store",
  "archive.ms": "archive",
  "unarchive.ms": "unarchive",
  "start.initializeMs": "initialize",
}

/** A line of an ACP chart is one adapter, with the size or the step that its measure names. */
function acpLegend(name: string): string {
  const parts = name.split(SERIES_NAME_SEPARATOR)
  const agent = parts.map((part) => acpAgentLabels.get(part)).find((label) => label != undefined) ?? name
  const measure = parts.find((part) => !acpAgentLabels.has(part)) ?? ""
  const size = sizes.find((it) => measure.includes(`.${it.id}.`))
  const step = size == undefined ? acpMeasureLabels[measure] : size.short
  return step == undefined ? agent : `${agent}, ${step}`
}

interface AcpChart extends AirChart {
  projects: string[]
}

/** One chart of an ACP step with one line for each adapter and size. `%s` in [template] is the size of the metric name. */
function acpBySize(label: string, description: string, template: string): AcpChart {
  return { label, description, measure: sizes.map((size) => template.replace("%s", size.id)), projects: acpSizedProjects }
}

const acpGroups: { title: string; charts: AcpChart[] }[] = [
  {
    title: "ACP adapters: open and continue a session of each size",
    charts: [
      {
        label: "Start an adapter: the initialize request",
        description: "The start of the adapter process and its answer to initialize, median of 3 starts. Air pays it for each new session.",
        measure: ["start.initializeMs"],
        projects: acpAllProjects,
      },
      acpBySize(
        "Open a session: session/load replays the history",
        "From session/load until its answer, with every replayed update. All user turns must arrive, and the model must not be called.",
        "open.%s.replayMs"
      ),
      acpBySize("Open a session: peak memory of the adapter", "The largest sum of the resident memory of the adapter process tree during the replay.", "open.%s.agentRssPeakMiB"),
      acpBySize(
        "Continue a session: from the prompt to the request to the model",
        "From session/prompt until the headers of the request reach the mock model. The time that the adapter prepares the turn: it reads the history and builds the request.",
        "continue.%s.toModelMs"
      ),
      acpBySize(
        "Continue a session: the size of the request to the model",
        "The body of the request of the continued turn, in bytes. It grows with the history that the adapter sends.",
        "continue.%s.requestBytes"
      ),
    ],
  },
  {
    title: "ACP adapters: a live turn of 241 model steps",
    charts: [
      {
        label: "Live turn: the time of the adapter",
        description:
          "20 000 text parts, 200 file reads, 20 outputs of 300 000 characters and 20 images, then a second turn in the same session. The time of the turn without the time of the mock model.",
        measure: ["live.agentMs", "live.second.agentMs"],
        projects: acpAllProjects,
      },
      {
        label: "Live turn: the delay of a text part",
        description: "From the flush of a stamped text part by the mock model to its arrival at the ACP client. junie sends no text parts, so it has no line.",
        measure: ["live.chunkLagP50Ms", "live.chunkLagP95Ms"],
        projects: acpSizedProjects,
      },
      {
        label: "Live turn: peak memory of the adapter",
        description: "The largest sum of the resident memory of the adapter process tree during the turn.",
        measure: ["live.agentRssPeakMiB", "live.second.agentRssPeakMiB"],
        projects: acpAllProjects,
      },
    ],
  },
  {
    title: "ACP adapters: the session list",
    charts: [
      {
        label: "Session list: the first page of 100 sessions",
        description: "session/list of a project with 1 000, 10 000 and 50 000 sessions, warm. A time that grows with the count means that a page reads the whole history.",
        measure: ["list.1k.activeFirst.warmMs", "list.10k.activeFirst.warmMs", "list.50k.activeFirst.warmMs"],
        projects: acpSizedProjects,
      },
      {
        label: "Session list: the first start after the index is gone",
        description: "The initialize of an adapter on a store of 61 000 sessions without its index, as after an update. Codex builds its state database here.",
        measure: ["list.initialize.coldMs"],
        projects: acpSizedProjects,
      },
      {
        label: "Archive and unarchive a session of the size XXL",
        description: "From _session/archive or _session/unarchive until session/list shows the session in its new place.",
        measure: ["archive.ms", "unarchive.ms"],
        projects: acpSizedProjects,
      },
    ],
  },
]
</script>
