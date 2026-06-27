<template>
  <ChartAccordion :lazy="true">
    <AccordionPanel value="0">
      <AccordionHeader>Additional metrics</AccordionHeader>
      <AccordionContent>
        <Divider label="GC" />
        <section>
          <GroupProjectsChart
            v-for="chart in gcCharts"
            :key="chart.key"
            :description="chart.description"
            :label="chart.label"
            :measure="chart.measure"
            :projects="projects"
            :value-unit="chart.valueUnit"
          />
        </section>
        <Divider label="Memory" />
        <section>
          <GroupProjectsChart
            v-for="chart in memoryCharts"
            :key="chart.key"
            :description="chart.description"
            :label="chart.label"
            :measure="chart.measure"
            :projects="projects"
            value-unit="counter"
          />
        </section>
      </AccordionContent>
    </AccordionPanel>
  </ChartAccordion>
</template>

<script lang="ts" setup>
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import ChartAccordion from "../charts/ChartAccordion.vue"
import Divider from "../common/Divider.vue"
import AccordionPanel from "primevue/accordionpanel"
import AccordionHeader from "primevue/accordionheader"
import AccordionContent from "primevue/accordioncontent"
import type { ValueUnit } from "../common/chart"

defineProps<{ projects: string[] }>()

interface ChartDef {
  key: string
  label: string
  measure: string
  description: string
  valueUnit?: ValueUnit
}

const gcCharts: ChartDef[] = [
  { key: "freed", label: "Freed", measure: "freedMemoryByGC", valueUnit: "counter", description: "Memory reclaimed by GC — an allocation-churn signal." },
  { key: "pauseTotal", label: "Pause total (ms)", measure: "gcPause", valueUnit: "ms", description: "Total GC pause time — steals from responsiveness." },
  { key: "pauseCount", label: "Pause count", measure: "gcPauseCount", valueUnit: "counter", description: "GC pauses during the run. A rising count signals allocation pressure." },
]

const memoryCharts: ChartDef[] = [
  { key: "avgRam", label: "Avg RAM", measure: "MEM.avgRamMegabytes", description: "Average resident RAM of the IDE process." },
  {
    key: "resident95p",
    label: "Resident 95p",
    measure: "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
    description: "95th-percentile resident set size — near-peak real footprint.",
  },
  { key: "maxHeap", label: "Max heap", measure: "JVM.maxHeapMegabytes", description: "Maximum JVM heap touched. Sustained growth signals heap pressure." },
  { key: "heapPeak", label: "Heap used, peak", measure: "totalHeapUsedMax", description: "Peak used JVM heap. With max heap, shows headroom." },
]
</script>
