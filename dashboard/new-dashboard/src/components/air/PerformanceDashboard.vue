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
</script>
