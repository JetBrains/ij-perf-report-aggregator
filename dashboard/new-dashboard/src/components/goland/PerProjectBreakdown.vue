<template>
  <!-- Small multiples: one row of charts per project, so that the lines of one chart are the measures of a single
       project and each measure keeps its color across the rows. -->
  <ChartAccordion :lazy="true">
    <AccordionPanel value="0">
      <AccordionHeader>{{ title }}</AccordionHeader>
      <AccordionContent>
        <template
          v-for="project in projects"
          :key="project"
        >
          <Divider :label="project" />
          <section>
            <GroupProjectsChart
              v-for="chart in charts"
              :key="chart.key"
              :label="chart.label"
              :measure="chart.measures"
              :projects="[project]"
            />
          </section>
        </template>
      </AccordionContent>
    </AccordionPanel>
  </ChartAccordion>
</template>

<script lang="ts" setup>
import AccordionContent from "openvue/accordioncontent"
import AccordionHeader from "openvue/accordionheader"
import AccordionPanel from "openvue/accordionpanel"
import ChartAccordion from "../charts/ChartAccordion.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import Divider from "../common/Divider.vue"

export interface BreakdownChart {
  key: string
  label: string
  measures: string[]
}

defineProps<{ title: string; projects: string[]; charts: BreakdownChart[] }>()
</script>
