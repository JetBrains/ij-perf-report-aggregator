<template>
  <!-- One lazy panel per project: opening one renders only its charts. A chart plots one family of inspections on
       one project, so each inspection keeps its color across the projects and between the time and the problems. -->
  <ChartAccordion :lazy="true">
    <AccordionPanel value="0">
      <AccordionHeader>{{ title }}</AccordionHeader>
      <AccordionContent>
        <ChartAccordion
          :lazy="true"
          :multiple="true"
        >
          <AccordionPanel
            v-for="project in projects"
            :key="project"
            :value="project"
          >
            <AccordionHeader>{{ project }}</AccordionHeader>
            <AccordionContent>
              <section
                v-for="family in families"
                :key="family.key"
                class="grid grid-cols-2 gap-x-4"
              >
                <GroupProjectsChart
                  :label="`${family.label}: time`"
                  :measure="family.inspections.map((it) => `${GLOBAL_INSPECTIONS}#${it}.ms`)"
                  :projects="[project]"
                  :legend-formatter="inspectionOf"
                />
                <GroupProjectsChart
                  :label="`${family.label}: problems`"
                  :measure="family.inspections.map((it) => `${GLOBAL_INSPECTIONS}#${it}.problems`)"
                  :projects="[project]"
                  :legend-formatter="inspectionOf"
                />
              </section>
            </AccordionContent>
          </AccordionPanel>
        </ChartAccordion>
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
import { InspectionFamily } from "./inspectionFamilies"

const GLOBAL_INSPECTIONS = "globalInspections"

defineProps<{ title: string; projects: string[]; families: InspectionFamily[] }>()

function inspectionOf(name: string): string {
  return name.replace(`${GLOBAL_INSPECTIONS}#`, "").replace(/\.(ms|problems)$/, "")
}
</script>
