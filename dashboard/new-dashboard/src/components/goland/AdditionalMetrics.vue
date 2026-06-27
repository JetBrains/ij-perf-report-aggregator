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
            :value-unit="chart.valueUnit"
          />
        </section>
        <Divider label="JVM" />
        <section>
          <GroupProjectsChart
            v-for="chart in jvmCharts"
            :key="chart.key"
            :description="chart.description"
            :label="chart.label"
            :measure="chart.measure"
            :projects="projects"
            :value-unit="chart.valueUnit"
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
  { key: "freedByFullGC", label: "Freed by full GC (MiB)", measure: "freedMemoryByFullGC", valueUnit: "counter", description: "Memory reclaimed by full GC only. High values signal promotion pressure." },
  { key: "fullGCPause", label: "Full GC pause (ms)", measure: "fullGCPause", valueUnit: "ms", description: "Time the IDE was fully paused during full GC. Any spike is a responsiveness killer." },
  { key: "g1MarkCycles", label: "G1 concurrent mark cycles", measure: "g1gcConcurrentMarkCycles", valueUnit: "counter", description: "Number of G1 concurrent marking cycles. Rising count = heap pressure." },
  { key: "g1MarkTime", label: "G1 concurrent mark time (ms)", measure: "g1gcConcurrentMarkTimeMs", valueUnit: "ms", description: "Wall time spent in G1 concurrent marking. Competes with UI threads." },
  { key: "g1ShrinkCount", label: "G1 heap shrinkage count", measure: "g1gcHeapShrinkageCount", valueUnit: "counter", description: "How often G1 returned heap to OS. Zero means heap is pinned at max." },
  { key: "g1ShrinkMb", label: "G1 heap shrinkage (MiB)", measure: "g1gcHeapShrinkageMegabytes", valueUnit: "counter", description: "Total heap memory returned to OS by G1 shrinkage." },
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
  { key: "fileMappings", label: "File-mapped RAM", measure: "MEM.avgFileMappingsRamMegabytes", description: "RAM consumed by memory-mapped files (JARs, indexes, etc.)." },
  { key: "ramMinusMappings", label: "RAM minus mappings", measure: "MEM.avgRamMinusFileMappingsMegabytes", description: "Resident RAM excluding file mappings — the true heap + off-heap footprint." },
  { key: "committedHeap", label: "Committed heap", measure: "JVM.committedHeapMegabytes", description: "Heap memory actually committed from OS. Versus max heap shows headroom." },
]

const jvmCharts: ChartDef[] = [
  { key: "maxThreads", label: "Max threads", measure: "JVM.maxThreadCount", valueUnit: "counter", description: "Peak thread count. Unexpected growth signals thread leaks or misconfiguration." },
  { key: "totalCpu", label: "Total CPU (ms)", measure: "JVM.totalCpuTimeMs", valueUnit: "ms", description: "Cumulative CPU time across all threads. Versus wall-clock time, reveals CPU saturation." },
  { key: "safepointTime", label: "Safepoint time (ms)", measure: "JVM.totalTimeToSafepointsMs", valueUnit: "ms", description: "Time threads spent blocked at JVM safepoints. Above 50 ms indicates contention." },
  { key: "gcCollectionTime", label: "GC collection time (ms)", measure: "JVM.GC.collectionTimesMs", valueUnit: "ms", description: "JVM-reported total GC time. Independent cross-check against gcPause." },
]
</script>
