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

const gcCharts: ChartDef[] = [
  { key: "freed", label: "Freed", measure: "freedMemoryByGC" },
  { key: "pauseTotal", label: "Pause total (ms)", measure: "gcPause" },
  { key: "pauseCount", label: "Pause count", measure: "gcPauseCount" },
  { key: "freedByFullGC", label: "Freed by full GC (MiB)", measure: "freedMemoryByFullGC" },
  { key: "fullGCPause", label: "Full GC pause (ms)", measure: "fullGCPause" },
  { key: "g1MarkCycles", label: "G1 concurrent mark cycles", measure: "g1gcConcurrentMarkCycles" },
  { key: "g1MarkTime", label: "G1 concurrent mark time (ms)", measure: "g1gcConcurrentMarkTimeMs" },
  { key: "g1ShrinkCount", label: "G1 heap shrinkage count", measure: "g1gcHeapShrinkageCount" },
  { key: "g1ShrinkMb", label: "G1 heap shrinkage (MiB)", measure: "g1gcHeapShrinkageMegabytes" },
]

const memoryCharts: ChartDef[] = [
  { key: "avgRam", label: "Avg RAM", measure: "MEM.avgRamMegabytes" },
  { key: "resident95p", label: "Resident 95p", measure: "Memory | IDE | RESIDENT SIZE (MB) 95th pctl" },
  { key: "maxHeap", label: "Max heap", measure: "JVM.maxHeapMegabytes" },
  { key: "heapPeak", label: "Heap used, peak", measure: "totalHeapUsedMax" },
  { key: "fileMappings", label: "File-mapped RAM", measure: "MEM.avgFileMappingsRamMegabytes" },
  { key: "ramMinusMappings", label: "RAM minus mappings", measure: "MEM.avgRamMinusFileMappingsMegabytes" },
  { key: "committedHeap", label: "Committed heap", measure: "JVM.committedHeapMegabytes" },
]

const jvmCharts: ChartDef[] = [
  { key: "maxThreads", label: "Max threads", measure: "JVM.maxThreadCount" },
  { key: "totalCpu", label: "Total CPU (ms)", measure: "JVM.totalCpuTimeMs" },
  { key: "safepointTime", label: "Safepoint time (ms)", measure: "JVM.totalTimeToSafepointsMs" },
  { key: "gcCollectionTime", label: "GC collection time (ms)", measure: "JVM.GC.collectionTimesMs" },
]
</script>
