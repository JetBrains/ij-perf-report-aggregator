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
  { key: "throughput", label: "Throughput (throughput)", measure: "throughput" },
  { key: "pauseTotal", label: "Pause total (gcPause)", measure: "gcPause" },
  { key: "maxPause", label: "Max pause (maxPause)", measure: "maxPause" },
  { key: "fullGCPause", label: "Full GC pause (fullGCPause)", measure: "fullGCPause" },
  { key: "maxFullGCPause", label: "Max full GC pause (maxFullGCPause)", measure: "maxFullGCPause" },
  { key: "accumPause", label: "Accumulated pause (accumPause)", measure: "accumPause" },
  { key: "pauseCount", label: "Pause count (gcPauseCount)", measure: "gcPauseCount" },
  { key: "fullPauseCount", label: "Full GC pause count (fullGcPauseCount)", measure: "fullGcPauseCount" },
  { key: "freed", label: "Freed (freedMemoryByGC)", measure: "freedMemoryByGC" },
  { key: "freedByFullGC", label: "Freed by full GC (freedMemoryByFullGC)", measure: "freedMemoryByFullGC" },
  { key: "g1MarkCycles", label: "G1 concurrent mark cycles (g1gcConcurrentMarkCycles)", measure: "g1gcConcurrentMarkCycles" },
  { key: "g1MarkTime", label: "G1 concurrent mark time (g1gcConcurrentMarkTimeMs)", measure: "g1gcConcurrentMarkTimeMs" },
  { key: "g1ShrinkCount", label: "G1 heap shrinkage count (g1gcHeapShrinkageCount)", measure: "g1gcHeapShrinkageCount" },
  { key: "g1ShrinkMb", label: "G1 heap shrinkage (g1gcHeapShrinkageMegabytes)", measure: "g1gcHeapShrinkageMegabytes" },
]

// Grouped: process RAM → heap sizing → non-heap.
const memoryCharts: ChartDef[] = [
  { key: "avgRam", label: "Avg RAM (MEM.avgRamMegabytes)", measure: "MEM.avgRamMegabytes" },
  { key: "resident95p", label: "Resident 95p (Memory | IDE | RESIDENT SIZE (MB) 95th pctl)", measure: "Memory | IDE | RESIDENT SIZE (MB) 95th pctl" },
  { key: "fileMappings", label: "File-mapped RAM (MEM.avgFileMappingsRamMegabytes)", measure: "MEM.avgFileMappingsRamMegabytes" },
  { key: "ramMinusMappings", label: "RAM minus mappings (MEM.avgRamMinusFileMappingsMegabytes)", measure: "MEM.avgRamMinusFileMappingsMegabytes" },
  { key: "maxHeap", label: "Max heap (JVM.maxHeapMegabytes)", measure: "JVM.maxHeapMegabytes" },
  { key: "committedHeap", label: "Committed heap (JVM.committedHeapMegabytes)", measure: "JVM.committedHeapMegabytes" },
  { key: "usedHeap", label: "Used heap (JVM.usedHeapMegabytes)", measure: "JVM.usedHeapMegabytes" },
  { key: "heapPeak", label: "Heap used, peak (totalHeapUsedMax)", measure: "totalHeapUsedMax" },
  { key: "footprintAfterFullGC", label: "Footprint after full GC (avgfootprintAfterFullGC)", measure: "avgfootprintAfterFullGC" },
  { key: "metaspacePeak", label: "Metaspace peak (totalPermUsedMax)", measure: "totalPermUsedMax" },
  { key: "nativeUsed", label: "Native used (JVM.usedNativeMegabytes)", measure: "JVM.usedNativeMegabytes" },
  { key: "directBuffers", label: "Direct byte buffers (JVM.totalDirectByteBuffersMegabytes)", measure: "JVM.totalDirectByteBuffersMegabytes" },
]

// Grouped: CPU → threads → safepoints → JVM-reported GC → allocation.
const jvmCharts: ChartDef[] = [
  { key: "totalCpu", label: "Total CPU (JVM.totalCpuTimeMs)", measure: "JVM.totalCpuTimeMs" },
  { key: "maxThreads", label: "Max threads (JVM.maxThreadCount)", measure: "JVM.maxThreadCount" },
  { key: "newThreads", label: "New threads (JVM.newThreadsCount)", measure: "JVM.newThreadsCount" },
  { key: "timeToSafepoints", label: "Time to safepoints (JVM.totalTimeToSafepointsMs)", measure: "JVM.totalTimeToSafepointsMs" },
  { key: "timeAtSafepoints", label: "Time at safepoints (JVM.totalTimeAtSafepointsMs)", measure: "JVM.totalTimeAtSafepointsMs" },
  { key: "safepointCount", label: "Safepoint count (JVM.totalSafepointCount)", measure: "JVM.totalSafepointCount" },
  { key: "gcCollectionTime", label: "GC collection time (JVM.GC.collectionTimesMs)", measure: "JVM.GC.collectionTimesMs" },
  { key: "gcCollections", label: "GC collections (JVM.GC.collections)", measure: "JVM.GC.collections" },
  { key: "allocated", label: "Allocated (JVM.totalMegabytesAllocated)", measure: "JVM.totalMegabytesAllocated" },
]

const osCharts: ChartDef[] = [{ key: "osLoad", label: "OS load average (OS.loadAverage)", measure: "OS.loadAverage" }]
</script>
