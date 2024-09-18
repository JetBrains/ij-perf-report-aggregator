<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="vcs_idea_ultimate_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <div>
      <Chip><a href="#index">Vcs indexing</a></Chip>
    </div>

    <Accordion
      :multiple="true"
      :active-index="[0]"
    >
      <AccordionTab header="Indexing">
        <section>
          <GroupProjectsChart
            id="index"
            label="Indexing"
            measure="vcs-log-indexing"
            :projects="indexingProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of collected commits"
            measure="vcs-log-indexing#numberOfCommits"
            :projects="indexingProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Real number of collected commits through git rev-list --count --all"
            measure="realNumberOfCommits"
            :projects="indexingProjects"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="'vcs-log' directory size in gb"
            measure="vcs-log-directory-size-gb"
            :projects="indexingProjects"
          />
        </section>
        <!--<section>-->
        <!--  <GroupProjectsChart-->
        <!--    label="Building of 'git commit-graph write &#45;&#45;reachable &#45;&#45;changed-paths' in minutes"-->
        <!--    measure="git-build-commit-graph"-->
        <!--    :projects="indexingProjects"-->
        <!--  />-->
        <!--</section>-->
      </AccordionTab>
    </Accordion>
  </DashboardPage>
</template>

<script setup lang="ts">
import Accordion from "primevue/accordion"
import Chip from "primevue/chip"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const intellijSpecificProject = "ide_starter/gitLogIndexing"
const indexingProjects = [intellijSpecificProject]
</script>
