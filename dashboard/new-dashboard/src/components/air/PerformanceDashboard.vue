<template>
  <DashboardPage
    db-name="perfintDev"
    table="air"
    persistent-id="air_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <Divider title="Long session M" />
    <section>
      <GroupProjectsChart
        v-for="chart in longSessionM"
        :key="chart.label"
        :label="chart.label"
        :measure="chart.measure"
        :projects="['air/long-session/M']"
      />
    </section>
    <Divider title="Long session XL" />
    <section>
      <GroupProjectsChart
        v-for="chart in longSessionXL"
        :key="chart.label"
        :label="chart.label"
        :measure="chart.measure"
        :projects="['air/long-session/XL']"
      />
    </section>
    <Divider title="Long session XXL" />
    <section>
      <GroupProjectsChart
        v-for="chart in longSessionXXL"
        :key="chart.label"
        :label="chart.label"
        :measure="chart.measure"
        :projects="['air/long-session/XXL']"
      />
    </section>
    <Divider title="Many agents" />
    <section>
      <GroupProjectsChart
        v-for="chart in manyAgents"
        :key="chart.label"
        :label="chart.label"
        :measure="chart.measure"
        :projects="['air/many-agents']"
      />
    </section>
    <Divider title="Session list" />
    <section>
      <GroupProjectsChart
        v-for="chart in sessionList"
        :key="chart.label"
        :label="chart.label"
        :measure="chart.measure"
        :projects="['air/session-list']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

interface AirChart {
  label: string
  measure: string[]
}

const longSessionM: AirChart[] = [
  { label: "M: Open session", measure: ["open.M.gestureToTabMs", "open.M.gestureToSettledMs"] },
  {
    label: "M: Checkpoints",
    measure: [
      "checkpoint.M.rich_projection_ready.totalMs",
      "checkpoint.M.frontend_reload_applied.totalMs",
      "checkpoint.M.frontend_transcript_rows_applied.elapsedMs",
      "checkpoint.M.frontend_first_content_paint_finished.totalSincePanelCreationMs",
    ],
  },
  { label: "M: Answer streaming", measure: ["answer.M.streamMs", "answer.M.streamStartToCompletedMs", "answer.M.streamStartToSettledMs"] },
  { label: "M: Continue session", measure: ["continue.M.sendToLoadMs", "continue.M.sendToReplayedMs", "continue.M.sendToPromptMs"] },
  { label: "M: Input and scroll latency (p95)", measure: ["prompt.M.codeP95Ms", "prompt.M.proseP95Ms", "scroll.M.wheelP95Ms", "answer.M.wheelP95Ms"] },
  { label: "M: Memory", measure: ["heap.M.before_the_open.MiB", "heap.M.after_the_open.MiB", "heap.M.heapLive.MiB", "heap.M.heapPeak.MiB"] },
]

const longSessionXL: AirChart[] = [
  { label: "XL: Open session", measure: ["open.XL.gestureToTabMs", "open.XL.gestureToSettledMs"] },
  {
    label: "XL: Checkpoints",
    measure: [
      "checkpoint.XL.rich_projection_ready.totalMs",
      "checkpoint.XL.frontend_reload_applied.totalMs",
      "checkpoint.XL.frontend_transcript_rows_applied.elapsedMs",
      "checkpoint.XL.frontend_first_content_paint_finished.totalSincePanelCreationMs",
    ],
  },
  { label: "XL: Answer streaming", measure: ["answer.XL.streamMs", "answer.XL.streamStartToCompletedMs", "answer.XL.streamStartToSettledMs"] },
  { label: "XL: Continue session", measure: ["continue.XL.sendToLoadMs", "continue.XL.sendToReplayedMs", "continue.XL.sendToPromptMs"] },
  { label: "XL: Input and scroll latency (p95)", measure: ["prompt.XL.codeP95Ms", "prompt.XL.proseP95Ms", "scroll.XL.wheelP95Ms", "answer.XL.wheelP95Ms"] },
  { label: "XL: Memory", measure: ["heap.XL.before_the_open.MiB", "heap.XL.after_the_open.MiB", "heap.XL.heapLive.MiB", "heap.XL.heapPeak.MiB"] },
]

const longSessionXXL: AirChart[] = [
  { label: "XXL: Open session", measure: ["open.XXL.gestureToTabMs", "open.XXL.gestureToSettledMs"] },
  {
    label: "XXL: Checkpoints",
    measure: [
      "checkpoint.XXL.rich_projection_ready.totalMs",
      "checkpoint.XXL.frontend_reload_applied.totalMs",
      "checkpoint.XXL.frontend_transcript_rows_applied.elapsedMs",
      "checkpoint.XXL.frontend_first_content_paint_finished.totalSincePanelCreationMs",
    ],
  },
  { label: "XXL: Answer streaming", measure: ["answer.XXL.streamMs", "answer.XXL.streamStartToCompletedMs", "answer.XXL.streamStartToSettledMs"] },
  { label: "XXL: Continue session", measure: ["continue.XXL.sendToLoadMs", "continue.XXL.sendToReplayedMs", "continue.XXL.sendToPromptMs"] },
  { label: "XXL: Input and scroll latency (p95)", measure: ["prompt.XXL.codeP95Ms", "prompt.XXL.proseP95Ms", "scroll.XXL.wheelP95Ms", "answer.XXL.wheelP95Ms"] },
  { label: "XXL: Memory", measure: ["heap.XXL.before_the_open.MiB", "heap.XXL.after_the_open.MiB", "heap.XXL.heapLive.MiB", "heap.XXL.heapPeak.MiB"] },
]

const manyAgents: AirChart[] = [
  { label: "Many agents: EDT latency", measure: ["manyAgents.edtP95Ms", "manyAgents.edtMaxMs", "manyAgents.roundMaxMs"] },
  { label: "Many agents: Tab lag", measure: ["manyAgents.focusedTabLagMs", "manyAgents.openTabLagMaxMs"] },
  { label: "Many agents: Memory", measure: ["manyAgents.heapLive.MiB", "manyAgents.heapPeak.MiB"] },
  { label: "Many agents: GC time", measure: ["manyAgents.gcMs"] },
]

const sessionList: AirChart[] = [
  { label: "Session list: Show and refresh", measure: ["sessionList.firstShowMs", "sessionList.refreshMedianMs"] },
  { label: "Session list: Scroll latency", measure: ["sessionList.wheelP95Ms", "sessionList.wheelMaxMs"] },
  { label: "Session list: Memory", measure: ["sessionList.heapLive.MiB", "sessionList.heapPeak.MiB"] },
]
</script>
