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
            :label="chart.label"
            :measure="chart.measure"
            :projects="projects"
          />
        </section>
        <Divider label="Memory" />
        <section>
          <GroupProjectsChart
            v-for="chart in memoryCharts"
            :key="chart.key"
            :label="chart.label"
            :measure="chart.measure"
            :projects="projects"
          />
        </section>
        <Divider label="JVM" />
        <section>
          <GroupProjectsChart
            v-for="chart in jvmCharts"
            :key="chart.key"
            :label="chart.label"
            :measure="chart.measure"
            :projects="projects"
          />
        </section>
        <Divider label="OS" />
        <section>
          <GroupProjectsChart
            v-for="chart in osCharts"
            :key="chart.key"
            :label="chart.label"
            :measure="chart.measure"
            :projects="projects"
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

defineProps<{ projects: string[] }>()

interface ChartDef {
  key: string
  label: string
  measure: string
}

// Grouped: efficiency → pauses → freed → G1 mark → G1 shrink.
const gcCharts: ChartDef[] = [
  { key: "throughput", label: "Throughput", measure: "throughput" },
  { key: "pauseTotal", label: "Pause total", measure: "gcPause" },
  { key: "maxPause", label: "Max pause", measure: "maxPause" },
  { key: "fullGCPause", label: "Full GC pause", measure: "fullGCPause" },
  { key: "maxFullGCPause", label: "Max full GC pause", measure: "maxFullGCPause" },
  { key: "accumPause", label: "Accumulated pause", measure: "accumPause" },
  { key: "pauseCount", label: "Pause count", measure: "gcPauseCount" },
  { key: "fullPauseCount", label: "Full GC pause count", measure: "fullGcPauseCount" },
  { key: "freed", label: "Freed", measure: "freedMemoryByGC" },
  { key: "freedByFullGC", label: "Freed by full GC", measure: "freedMemoryByFullGC" },
  { key: "g1MarkCycles", label: "G1 concurrent mark cycles", measure: "g1gcConcurrentMarkCycles" },
  { key: "g1MarkTime", label: "G1 concurrent mark time", measure: "g1gcConcurrentMarkTimeMs" },
  { key: "g1ShrinkCount", label: "G1 heap shrinkage count", measure: "g1gcHeapShrinkageCount" },
  { key: "g1ShrinkMb", label: "G1 heap shrinkage", measure: "g1gcHeapShrinkageMegabytes" },
]

// Grouped: process RAM → heap sizing → non-heap.
const memoryCharts: ChartDef[] = [
  { key: "avgRam", label: "Avg RAM", measure: "MEM.avgRamMegabytes" },
  { key: "resident95p", label: "Resident 95p", measure: "Memory | IDE | RESIDENT SIZE (MB) 95th pctl" },
  { key: "fileMappings", label: "File-mapped RAM", measure: "MEM.avgFileMappingsRamMegabytes" },
  { key: "ramMinusMappings", label: "RAM minus mappings", measure: "MEM.avgRamMinusFileMappingsMegabytes" },
  { key: "maxHeap", label: "Max heap", measure: "JVM.maxHeapMegabytes" },
  { key: "committedHeap", label: "Committed heap", measure: "JVM.committedHeapMegabytes" },
  { key: "usedHeap", label: "Used heap", measure: "JVM.usedHeapMegabytes" },
  { key: "heapPeak", label: "Heap used, peak", measure: "totalHeapUsedMax" },
  { key: "footprintAfterFullGC", label: "Footprint after full GC", measure: "avgfootprintAfterFullGC" },
  { key: "metaspacePeak", label: "Metaspace peak", measure: "totalPermUsedMax" },
  { key: "nativeUsed", label: "Native used", measure: "JVM.usedNativeMegabytes" },
  { key: "directBuffers", label: "Direct byte buffers", measure: "JVM.totalDirectByteBuffersMegabytes" },
]

// Grouped: CPU → threads → safepoints → JVM-reported GC → allocation.
const jvmCharts: ChartDef[] = [
  { key: "totalCpu", label: "Total CPU", measure: "JVM.totalCpuTimeMs" },
  { key: "maxThreads", label: "Max threads", measure: "JVM.maxThreadCount" },
  { key: "newThreads", label: "New threads", measure: "JVM.newThreadsCount" },
  { key: "timeToSafepoints", label: "Time to safepoints", measure: "JVM.totalTimeToSafepointsMs" },
  { key: "timeAtSafepoints", label: "Time at safepoints", measure: "JVM.totalTimeAtSafepointsMs" },
  { key: "safepointCount", label: "Safepoint count", measure: "JVM.totalSafepointCount" },
  { key: "gcCollectionTime", label: "GC collection time", measure: "JVM.GC.collectionTimesMs" },
  { key: "gcCollections", label: "GC collections", measure: "JVM.GC.collections" },
  { key: "allocated", label: "Allocated", measure: "JVM.totalMegabytesAllocated" },
]

const osCharts: ChartDef[] = [{ key: "osLoad", label: "OS load average", measure: "OS.loadAverage" }]
</script>
